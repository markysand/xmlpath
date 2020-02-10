package xmlpath_test

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/markysand/xmlpath"
	"github.com/stretchr/testify/assert"
)

func Test_Pipe(t *testing.T) {
	t.Run(
		"Multiple path parsing",
		func(t *testing.T) {
			type Person struct {
				Name      string `xml:"name"`
				SchoolRef string `xml:"school-ref"`
			}

			type School struct {
				Text string `xml:",chardata"`
				ID   string `xml:"id,attr"`
			}

			file, err := os.Open("./testdata/data.xml")
			defer file.Close()
			if err != nil {
				t.Fatal("Could not read file " + err.Error())
			}
			defer file.Close()

			var persons []Person
			pathConfig1 := xmlpath.NewPathConfig(
				func(decodeInto func(interface{})) {
					person := new(Person)
					decodeInto(person)
					persons = append(persons, *person)
				},
				"test-education-register", "persons", "person",
			)

			var schools []School
			pathConfig2 := xmlpath.NewPathConfig(func(decodeInto func(interface{})) {
				school := new(School)
				decodeInto(school)
				schools = append(schools, *school)
			}, "test-education-register", "educations", "education-sites", "school")

			_, err = xmlpath.Pipe(file, pathConfig1, pathConfig2)
			if err != nil {
				t.Fatal(err)
			}
			refPersons := []Person{
				Person{"Kit", "1"}, Person{"Cool", "2"}, Person{"Calm", "1"},
			}
			assert.Equal(t, refPersons, persons)

			refSchools := []School{
				School{
					"Hogwarts", "1",
				},
				School{
					"School of Life", "2",
				},
			}
			assert.Equal(t, refSchools, schools)
		})
	t.Run("Single path parsing", func(t *testing.T) {
		file, err := os.Open("./testdata/data.xml")
		defer file.Close()
		if err != nil {
			t.Fatal("Could not read file " + err.Error())
		}
		defer file.Close()

		type Person struct {
			Name      string `xml:"name"`
			SchoolRef string `xml:"school-ref"`
		}

		persons := new([]Person)

		p1 := xmlpath.NewPathConfig(func(decodeInto func(interface{})) { decodeInto(persons) },
			"test-education-register", "persons", "person")

		_, err = xmlpath.Pipe(file, p1)

		refPersons := []Person{
			Person{"Kit", "1"}, Person{"Cool", "2"}, Person{"Calm", "1"},
		}
		assert.Equal(t, refPersons, *persons)
	})
}

func ExamplePipe() {
	xmlContent := `<?xml version="1.0" encoding="UTF-8" ?>
<document>
	<basket>
		<apple>Astrakan</apple>
		<apple>Fuji</apple>
		<pear>Concorde</pear>
		<apple>Gala</apple>
	</basket>
</document>`

	xmlFile := strings.NewReader(xmlContent)

	type Apple struct {
		XMLName xml.Name `xml:"apple"`
		Text    string   `xml:",chardata"`
	}

	pathConfig1 := xmlpath.NewPathConfig(func(decodeInto func(interface{})) {

		apple := new(Apple)
		decodeInto(apple)
		fmt.Println(apple.Text)
	},
		"document", "basket", "apple")
	xmlpath.Pipe(xmlFile, pathConfig1)

	// Output:
	// Astrakan
	// Fuji
	// Gala
}
