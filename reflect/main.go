package main

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type User struct {
	Name string `json:"test"`
	Age  int64
	Has  bool
	Sl   []string
}

type City struct {
	Name       string
	Population int64 `json:"test"`
	GDP        int64
	Mayor      string
}

func main() {
	var u User = User{"", 10, false, []string{"test", "one more"}}

	res, err := JSONEncode(u)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(res))

	c := City{"sf", 5000000, 567896, "mr jones"}
	res, err = JSONEncode(c)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(res))
}

func JSONEncode(v interface{}) ([]byte, error) {
	buf := bytes.Buffer{}
	buf.WriteString("{")

	s := reflect.ValueOf(v)
	typeOfT := s.Type()
	// TODO: check if v is a struct else return error
	if s.Kind() != reflect.Struct {
		return nil, errors.New("Encode target is not structure")
	}
	// TODO: iterate over v`s reflect value using NumField()
	var c int
	for i := 0; i < s.NumField(); i++ {
		name := typeOfT.Field(i).Name
		jsonTag := typeOfT.Field(i).Tag.Get("json")
		if jsonTag != "" {
			tags := strings.Split(jsonTag, ",")
			if len(tags) == 2 {
				if tags[1] == "-" {
					continue
				}
				if tags[1] == "omitempty" && s.Field(i).IsZero() {
					continue
				}
			}
			if len(tags) >= 1 {
				name = tags[0]
			}
		}

		if c != 0 {
			buf.WriteString(", ")
		}
		buf.WriteRune('"')

		buf.WriteString(name)
		buf.WriteRune('"')
		buf.WriteRune(':')
		buf.WriteRune(' ')
		buf.WriteRune('"')
		switch s.Field(i).Kind() {
		case reflect.String:
			buf.WriteString(s.Field(i).String())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			buf.WriteString(strconv.Itoa(int(s.Field(i).Int())))
		case reflect.Bool:
			buf.WriteString(strconv.FormatBool(s.Field(i).Bool()))
		case reflect.Array, reflect.Slice:
			fieldLen := s.Field(i).Len()
			for y := 0; y < fieldLen; y++ {
				//Do something with slice element
			}

		}
		buf.WriteRune('"')
		c++
	}
	// use type switch to create string result of "{field}" + ": " + "{value}"
	// start with just 2 types - reflect.String and reflect.Int64

	buf.WriteString("}")
	return buf.Bytes(), nil
}
