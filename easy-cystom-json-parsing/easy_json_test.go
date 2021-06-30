package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type personEmptyCases struct {
	name   string
	input  Person
	result bool
}

type placeEmptyCases struct {
	name   string
	input  Place
	result bool
}

func TestPerson_Empty(t *testing.T) {
	casess := []personEmptyCases{
		{
			name:   "Empty Struct",
			input:  Person{},
			result: true,
		},
		{
			name:   "Not Empty Struct",
			input:  Person{Name: "Test", Age: 1},
			result: false,
		},
	}
	for _, c := range casess {
		t.Run(c.name, func(t *testing.T) {
			got := c.input.Empty()
			want := c.result
			assert.Equal(t, got, want, "Should be equal")
		})
	}

}

func TestPlace_Emptye(t *testing.T) {
	casess := []placeEmptyCases{
		{
			name:   "Empty Struct",
			input:  Place{},
			result: true,
		},
		{
			name:   "Not empty Struct",
			input:  Place{City: "Test", Country: "test"},
			result: false,
		},
	}
	for _, c := range casess {
		t.Run(c.name, func(t *testing.T) {
			got := c.input.Empty()
			want := c.result
			assert.Equal(t, got, want, "Should be equal")
		})
	}
}

type testLogger struct {
}

func (l testLogger) Println(data ...interface{}) {

}

func (l testLogger) Fatalf(format string, v ...interface{}) {

}

func TestService(t *testing.T) {
	t.Run("Should implement HumanDecoder Interface", func(t *testing.T) {
		var h HumanDecoder
		h = &Service{log: &testLogger{}}
		var s interface{} = h
		switch s.(type) {
		case HumanDecoder:
		default:
			t.Fatal("Should implement HumanDecoder and doesnt")
		}
	})
}

func TestService_Decode(t *testing.T) {
	testData := []byte(`
		{
			"things": [
				{
					"name": "Alice",
					"age": 37
				},
				{
					"city": "Ipoh",
					"country": "Malaysia"
				}
			]
		}`)
	s := &Service{log: &testLogger{}}
	persons := []Person{
		{
			Name: "Alice",
			Age:  37,
		},
	}
	places := []Place{
		{
			City:    "Ipoh",
			Country: "Malaysia",
		},
	}
	var pn []Person
	var pl []Place
	pn, pl = s.Decode(testData)
	assert.Equal(t, pn, persons, "")
	assert.Equal(t, pl, places, "")
}

func TestService_PersonSort(t *testing.T) {
	persons := Persons{
		Person{
			Name: "Alice",
			Age:  37,
		},
		Person{
			Name: "Roman",
			Age:  3,
		},
	}
	sortedPersons := Persons{
		Person{
			Name: "Roman",
			Age:  3,
		},
		Person{
			Name: "Alice",
			Age:  37,
		},
	}

	s := &Service{log: &testLogger{}}
	s.Sort(persons)
	assert.Equal(t, persons, sortedPersons, "Should be equal")
}

func TestService_PlaceSort(t *testing.T) {
	places := Places{
		Place{
			City:    "Longer",
			Country: "Malaysia",
		},
		Place{
			City:    "Short",
			Country: "Malaysia",
		},
	}
	sortedPlaces := Places{
		Place{
			City:    "Short",
			Country: "Malaysia",
		},
		Place{
			City:    "Longer",
			Country: "Malaysia",
		},
	}

	s := &Service{log: &testLogger{}}
	s.Sort(places)
	assert.Equal(t, places, sortedPlaces, "Should be equal")
}
