// Package xmlpath simplifies stream parsing large xml-files.
package xmlpath

import (
	"encoding/xml"
	"fmt"
	"io"

	"github.com/pkg/errors"
)

// Decoder is the "double" callback type given to Pipe from
// each PathConfig. The decodeInto function is used to extract
// the xml value into a *Type variable of the users choice,
// as with standard xml decoding
type Decoder func(decodeInto func(target interface{}) error)

type matchType int

const (
	without matchType = iota
	within
	exact
)

type pathElements []string

func (pe pathElements) match(ss []string) matchType {
	reference, dynamic := len(pe), len(ss)
	for i := 0; i < reference && i < dynamic; i++ {
		if pe[i] != ss[i] {
			return without
		}
	}
	switch {
	case reference == dynamic:
		return exact
	case reference > dynamic:
		return within
	default:
		return without
	}
}

func (pe *pathElements) add(s string) {
	*pe = append(*pe, s)
}

func (pe *pathElements) pop() {
	*pe = (*pe)[:len(*pe)-1]
}

// PathConfig contains callback and elements of the path
type PathConfig struct {
	Decoder
	pathElements
}

// NewPathConfig sets up a path with a callback
func NewPathConfig(callback Decoder, pathElements ...string) PathConfig {
	return PathConfig{
		callback, pathElements,
	}
}

// Pipe triggers the decoding callback for each path
func Pipe(source io.Reader, paths ...PathConfig) (int, error) {
	decoder := xml.NewDecoder(source)
	xmlTokenDecoder := xml.NewTokenDecoder(decoder)

	// test paths for interference
	if err := testInterference(paths); err != nil {
		return 0, err
	}
	var (
		parsedDocuments int
		currentPath     pathElements
	)
LOOP:
	for {
		token, err := xmlTokenDecoder.Token()
		if err == io.EOF {
			return parsedDocuments, nil
		}
		if err != nil {
			return 0, errors.Wrap(err, "Could not parse XML")
		}
		switch tokenType := token.(type) {
		case xml.StartElement:
			currentPath.add(tokenType.Name.Local)
			var withinAnyPath bool
			for _, path := range paths {
				switch path.match(currentPath) {
				case exact:
					path.Decoder(func(target interface{}) error {
						err := xmlTokenDecoder.DecodeElement(target, &tokenType)
						return err
					})
					parsedDocuments++
					currentPath.pop()
					continue LOOP
				case within:
					withinAnyPath = true
				}
			}
			if !withinAnyPath {
				xmlTokenDecoder.Skip()
				currentPath.pop()
			}
		case xml.EndElement:
			if len(currentPath) == 0 {
				return 0, errors.New(fmt.Sprint(
					"Broken through root level at, XML error",
					tokenType.Name.Local))
			}
			currentPath.pop()
		}
	}
}

func testInterference(paths []PathConfig) error {
	for i := 0; i < len(paths)-1; i++ {
		for j := i + 1; j < len(paths); j++ {
			if paths[i].match(paths[j].pathElements) == within || paths[j].match(paths[i].pathElements) == within {
				return fmt.Errorf("Path %v and path %v - illegal interference", paths[i].pathElements, paths[j].pathElements)
			}
		}
	}
	return nil
}
