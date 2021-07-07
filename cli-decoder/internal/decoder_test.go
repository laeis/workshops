package internal

import (
	"errors"
	"testing"
)

type TestNewCliDecoderCases struct {
	name        string
	decoderType string
	testData    string
	success     bool
}

func TestNewCliDecoder(t *testing.T) {
	cases := []TestNewCliDecoderCases{
		{
			name:        "Valid xml",
			decoderType: "xml",
			testData:    "<note><to>Tove</to><from>Jani</from><heading>Reminder</heading><body>Don't forget me this weekend!</body></note>",
			success:     true,
		},
		{
			name:        "Invalid xml",
			decoderType: "xml",
			testData:    "some string",
			success:     false,
		},
		{
			name:        "Valid json",
			decoderType: "json",
			testData:    "{\"alias\":\"go-dms-workshop\",\"desc\":\"Create app and try it with different DMS\", \"type\":\"important\", \"ts\":1473837996,\"tags\":[\"Golang\",\"Workshop\",\"DMS\"],\"etime\":\"4h\",\"rtime\":\"8h\",\"reminders\":[\"3h\", \"15m\"]}",
			success:     true,
		},
		{
			name:        "Invalid json",
			decoderType: "json",
			testData:    "some string",
			success:     false,
		},
		{
			name:        "Invalid decoder",
			decoderType: "undefined",
			testData:    "some string",
			success:     false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var err error
			d := NewCliDecoder(c.decoderType)
			if d != nil {
				err = d.Decode([]byte(c.testData))
			} else {
				err = errors.New("Empty decoder type")
			}

			if c.success && err != nil {
				t.Errorf("Assert succec got error %v", err)
			}
			if !c.success && err == nil {
				t.Error("Assert error and dont have it")
			}
		})
	}
}
