package internal

import (
	"encoding/json"
	"encoding/xml"
)

const (
	XML  = "xml"
	JSON = "json"
)

type CliDecoder interface {
	Decode(data []byte) error
}

type JSONDecoder struct {
	i interface{}
}

func (j *JSONDecoder) Decode(data []byte) error {
	return json.Unmarshal(data, &j.i)
}

type XMLDecoder struct {
	i interface{}
}

func (x *XMLDecoder) Decode(data []byte) error {
	return xml.Unmarshal(data, &x.i)
}

func NewCliDecoder(decoder string) CliDecoder {
	switch decoder {
	case JSON:
		return &JSONDecoder{}
	case XML:
		return &XMLDecoder{}
	default:
		return nil
	}
}
