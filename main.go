package xmlpath

import (
	"encoding/xml"
	"fmt"
	"io"

	"github.com/pkg/errors"
)

// PipeXML is an utility that will get tokens of a specific type, read one by one
// by the callback sent into PipeXml
func PipeXML(xmlTokenDecoder *xml.Decoder, path []string, callback func(start *xml.StartElement)) (int, error) {
	pathLength := len(path)
	if pathLength == 0 {
		return 0, errors.New("Do not call PipeXML with zero path - mistake?")
	}
	var pathDepth, parsedDocuments int

	// Selective path tree parsing
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
			if tokenType.Name.Local == path[pathDepth] {
				if pathDepth == pathLength-1 { // time to harvest!
					callback(&tokenType)
					parsedDocuments++
					// xmlTokenDecoder.Skip()
				} else {
					pathDepth++ // not skip
				}
			} else {
				xmlTokenDecoder.Skip()
			}
		case xml.EndElement:
			if pathDepth == 0 {
				return 0, errors.New(fmt.Sprint(
					"Gone through root level at",
					tokenType.Name.Local))
			}
			pathDepth--
			if path[pathDepth] != tokenType.Name.Local {
				return 0, errors.New(
					fmt.Sprintf("Got end token %v but it does not match expected %v",
						tokenType.Name.Local,
						path[pathDepth]))
			}
		}
	}
}
