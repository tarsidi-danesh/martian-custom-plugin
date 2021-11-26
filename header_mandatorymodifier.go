package header

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/martian/parse"
	"github.com/google/uuid"
	"net/http"
)

func init() {
	parse.Register("header.MandatoryModifier", mandatoryRequestConfig)
}

type ValueTypes int

const (
	DYNAMIC ValueTypes = 0
	STATIC             = 1
)

type Generators int

const (
	NONE Generators = 0
	UUID            = 1
)

type HeaderValue struct {
	ValueType   ValueTypes `json:"type"`
	Generator   Generators `json:"generator"`
	StaticValue string     `json:"staticValue"`
}

type Header struct {
	Name  string      `json:"name"`
	Value HeaderValue `json:"value"`
}

type MandatoryRequestConfig struct {
	Headers []Header `json:"headers"`
}

type MandatoryRequestModifier struct {
	Config MandatoryRequestConfig
}

//TODO: Add logic to validate MandatoryRequestConfig
func NewMandatoryRequestModifier(mc MandatoryRequestConfig) (*MandatoryRequestModifier, error) {
	return &MandatoryRequestModifier{
		Config: mc,
	}, nil
}

func (mr *MandatoryRequestModifier) ModifyRequest(req *http.Request) error {
	for _, v := range mr.Config.Headers {
		if v.Value.ValueType == STATIC {
			req.Header[v.Name] = []string{v.Value.StaticValue}
		} else {
			var headerValue string
			switch v.Value.Generator {
			case UUID:
				{
					u := uuid.New()
					headerValue = fmt.Sprintf("%s", u)
				}
			default:
				headerValue = ""
			}
			req.Header[v.Name] = []string{headerValue}
		}
	}
	return nil
}

// mandatoryRequestConfig takes a JSON message as a byte slice and returns
// a mandatoryRequest.Modifier and an error.
//
// Example JSON:
// {
//  "headers": [
//		{
//			"name": "serviceId",
//			"value": {"type": "STATIC", "generator": "NONE", "staticValue": "GATEWAY"}
//		},
//		{
//			"name": "requestId",
//			"value": {"type": "DYNAMIC", "generator": "UUID", "staticValue": null}
//		}
// 	]
// }
func mandatoryRequestConfig(b []byte) (*parse.Result, error) {
	headerConfig := MandatoryRequestConfig{}

	if err := json.Unmarshal(b, &headerConfig); err != nil {
		return nil, err
	}
	modifier, err := NewMandatoryRequestModifier(headerConfig)
	if err != nil {
		return nil, errors.New("invalid configuration")
	}
	return parse.NewResult(modifier, []parse.ModifierType{parse.Request})
}