package validators

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type stringListCase struct {
	Name     string
	Needle   string
	Haystack stringList
	Want     bool
}

func TestStringList_Contains(t *testing.T) {
	cases := []stringListCase{
		{
			Name:     "Needle exists",
			Needle:   "exists",
			Haystack: stringList{"first", "exists", "third"},
			Want:     true,
		},
		{
			Name:     "Needle not exists",
			Needle:   "not exists",
			Haystack: stringList{"first", "exists", "third"},
			Want:     false,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			got := c.Haystack.Contains(c.Needle)
			assert.Equal(t, got, c.Want, "Got and result must be the same")
		})
	}

}

type validatorTestCase struct {
	Name  string
	Field string
	Want  string
	Input []string
}

func TestTaskValidator_ValidateParameter(t *testing.T) {
	var validator TaskValidator
	cases := []validatorTestCase{
		{
			Name:  "Category exists",
			Field: "category",
			Want:  NOTE,
			Input: []string{NOTE, "category"},
		},
		{
			Name:  "Category not exists",
			Field: "category",
			Want:  "",
			Input: []string{"wrong", "category"},
		},
		{
			Name:  "Category empty input",
			Field: "category",
			Want:  "",
			Input: []string{},
		},
		//period
		{
			Name:  "Period exists",
			Field: "period",
			Want:  DAY,
			Input: []string{DAY, "category"},
		},
		{
			Name:  "Period not exists",
			Field: "period",
			Want:  "",
			Input: []string{"wrong", "period"},
		},
		{
			Name:  "Period empty input",
			Field: "period",
			Want:  "",
			Input: []string{},
		},
		//order
		{
			Name:  "Order exists",
			Field: "order",
			Want:  DESC,
			Input: []string{DESC},
		},
		{
			Name:  "Order not exists",
			Field: "order",
			Want:  "asc",
			Input: []string{"wrong", "category"},
		},
		{
			Name:  "Order empty input",
			Field: "order",
			Want:  "asc",
			Input: []string{},
		},
		//order_by
		{
			Name:  "Order by exists",
			Field: "order_by",
			Want:  "title",
			Input: []string{"title", "category"},
		},
		{
			Name:  "Order by exists",
			Field: "order_by",
			Want:  "id",
			Input: []string{"wrong", "category"},
		},
		{
			Name:  "Order by input",
			Field: "order_by",
			Want:  "id",
			Input: []string{},
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			got := validator.ValidateParameter(c.Field, c.Input)
			assert.Equal(t, got, c.Want, "Got and want must be the same")
		})
	}

}
