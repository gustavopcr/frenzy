package http

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gustavopcr/frenzy/internal/payloadgen"
)

func TesteGet(url string, schema *openapi3.Schema) (*http.Response, error) {
	//TODO add query string later on
	resp, err := http.Get(url)
	if err != nil {
		return &http.Response{Status: "400", StatusCode: http.StatusBadRequest}, nil
	}
	return resp, nil
}

func TestePost(url string, schema *openapi3.Schema) (*http.Response, error) {
	pg := payloadgen.NewPayloadGenerator()
	teste := pg.PayloadFromSchema(schema)

	jsonBytes, err := json.MarshalIndent(teste, "", " ")
	if err != nil {
		panic(err)
	}

	body := bytes.NewReader(jsonBytes)
	resp, err := http.Post(url, "application/json", body)
	if err != nil {
		return &http.Response{Status: "400", StatusCode: http.StatusBadRequest}, nil
	}
	return resp, nil
}
