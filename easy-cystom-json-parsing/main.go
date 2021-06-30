package main

import (
	"encoding/json"
	"log"
	"os"
	"sort"
	"unicode/utf8"
)

var jsonStr = []byte(`
{
    "things": [
        {
            "name": "Alice",
            "age": 37
        },
        {
            "city": "Ipoh",
            "country": "Malaysia"
        },
        {
            "name": "Bob",
            "age": 36
        },
        {
            "city": "Northampton",
            "country": "England"
        },
 		{
            "name": "Albert",
            "age": 3
        },
		{
            "city": "Dnipro",
            "country": "Ukraine"
        },
		{
            "name": "Roman",
            "age": 32
        },
		{
            "city": "New York City",
            "country": "US"
        }
    ]
}`)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (p Person) Empty() bool {
	return p.Age == 0 && p.Name == ""
}

type Persons []Person

func (p Persons) Less(i, j int) bool {
	return p[i].Age < p[j].Age
}

func (p Persons) Len() int {
	return len(p)
}

func (p Persons) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type Place struct {
	City    string `json:"city"`
	Country string `json:"country"`
}

func (p Place) Empty() bool {
	return p.City == "" && p.Country == ""
}

type Places []Place

func (p Places) Less(i, j int) bool {
	return utf8.RuneCountInString(p[i].City) < utf8.RuneCountInString(p[j].City)
}

func (p Places) Len() int {
	return len(p)
}

func (p Places) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type Thing struct {
	Person
	Place
}

type Things struct {
	List []Thing `json:"things"`
}

type HumanDecoder interface {
	Decode(data []byte) ([]Person, []Place)
	Sort(dataToSort interface{})
	Print(interface{})
}

type Logger interface {
	Println(v ...interface{})
	Fatalf(format string, v ...interface{})
}

type Service struct {
	log Logger
}

func (s *Service) Decode(data []byte) ([]Person, []Place) {
	var things Things
	err := json.Unmarshal(data, &things)
	if err != nil {
		log.Fatal(err)
	}
	var persons []Person
	var places []Place
	for _, v := range things.List {
		if !v.Place.Empty() {
			places = append(places, v.Place)
		}
		if !v.Person.Empty() {
			persons = append(persons, v.Person)
		}
	}
	return persons, places
}

func (s *Service) Sort(dataToSort interface{}) {
	switch dataToSort.(type) {
	case Persons:
		sort.Sort(dataToSort.(Persons))
	case Places:
		sort.Sort(dataToSort.(Places))
	default:
		s.log.Println("Wrong type for sort")
	}

}
func (s *Service) Print(data interface{}) {
	s.log.Println(data)
}

func main() {
	var s HumanDecoder
	// logger to Inject
	logger := log.New(os.Stdout, "INFO: ", 0)
	s = &Service{log: logger}
	var persons Persons
	var places Places

	persons, places = s.Decode(jsonStr)

	s.Sort(persons)
	s.Sort(places)

	s.Print(persons)
	s.Print(places)
}
