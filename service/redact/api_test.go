package redact_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pangeacyber/go-pangea/internal/pangeatesting"
	"github.com/pangeacyber/go-pangea/pangea"
	"github.com/pangeacyber/go-pangea/service/redact"
	"github.com/stretchr/testify/assert"
)

func TestRedact(t *testing.T) {
	mux, url, teardown := pangeatesting.SetupServer()
	defer teardown()

	mux.HandleFunc("/v1/redact", func(w http.ResponseWriter, r *http.Request) {
		pangeatesting.TestMethod(t, r, "POST")
		pangeatesting.TestBody(t, r, `{"text":"My phone number is: 110303456"}`)
		fmt.Fprint(w,
			`{
				"request_id": "some-id",
				"request_time": "1970-01-01T00:00:00Z",
				"response_time": "1970-01-01T00:00:10Z",
				"status_code": 200,
				"status": "success",
				"result":{
					"redacted_text": "My phone number is: <PHONE_NUMBER>"
				},
				"summary": "success"
			}`)
	})

	client, _ := redact.New(pangeatesting.TestConfig(url))
	input := &redact.TextInput{
		Text: pangea.String("My phone number is: 110303456"),
	}
	ctx := context.Background()
	got, err := client.Redact(ctx, input)

	assert.NoError(t, err)

	want := &redact.TextOutput{
		RedactedText: pangea.String("My phone number is: <PHONE_NUMBER>"),
	}
	assert.Equal(t, want, got.Result)
}

func TestRedactStructured(t *testing.T) {
	mux, url, teardown := pangeatesting.SetupServer()
	defer teardown()

	mux.HandleFunc("/v1/redact_structured", func(w http.ResponseWriter, r *http.Request) {
		pangeatesting.TestMethod(t, r, "POST")
		pangeatesting.TestBody(t, r, `{"data":{"one":{"secret":"(555)-555-5555"}},"jsonp":["$.*.secret"]}`)
		fmt.Fprint(w,
			`{
				"request_id": "some-id",
				"request_time": "1970-01-01T00:00:00Z",
				"response_time": "1970-01-01T00:00:10Z",
				"status_code": 200,
				"status": "success",
				"result": {
					"redacted_data": {
					  "one": { "secret": "<PHONE_NUMBER>" }
					}
				},
				"summary": "success"
			}`)
	})

	client, _ := redact.New(pangeatesting.TestConfig(url))

	type (
		innerType struct {
			Secret string `json:"secret"`
		}
		Payload struct {
			One innerType `json:"one"`
		}
	)

	input := &redact.StructuredInput{
		JSONP: []*string{
			pangea.String("$.*.secret"),
		},
	}
	input.SetData(Payload{One: innerType{Secret: "(555)-555-5555"}})
	ctx := context.Background()
	response, err := client.RedactStructured(ctx, input)

	assert.NoError(t, err)
	assert.NotEmpty(t, response.Result.RedactedData)

	var got Payload
	want := Payload{One: innerType{Secret: "<PHONE_NUMBER>"}}
	assert.NoError(t, response.Result.GetRedactedData(&got))
	assert.Equal(t, want, got)
}

func TestRedactError(t *testing.T) {
	f := func(cfg *pangea.Config) error {
		client, _ := redact.New(cfg)
		_, err := client.Redact(context.Background(), nil)
		return err
	}
	pangeatesting.TestNewRequestAndDoFailure(t, "Redact.Redact", f)
}

func TestRedactStructuredError(t *testing.T) {
	f := func(cfg *pangea.Config) error {
		client, _ := redact.New(cfg)
		_, err := client.RedactStructured(context.Background(), nil)
		return err
	}
	pangeatesting.TestNewRequestAndDoFailure(t, "Redact.RedactStructured", f)
}
