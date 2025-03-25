package main

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

func TestGenereatePayloadFromSchema(t *testing.T) {
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile("user_schema.yaml")
	if err != nil {
		log.Fatalf("failed to load schema: %v", err)
	}

	if err := doc.Validate(loader.Context); err != nil {
		log.Fatalf("Error validating with loader: %v", err)
	}

	userSchemaRef := doc.Components.Schemas["User"]
	userSchema := userSchemaRef.Value
	user := generatePayloadFromSchema(userSchema)

	expected := map[string]any{
		"id":   10,
		"name": "string",
	}

	resultJSON, _ := json.Marshal(user)
	expectedJSON, _ := json.Marshal(expected)

	if string(resultJSON) != string(expectedJSON) {
		t.Errorf("Geenerated paylod does not match expected. \nExpected: %s\nGot: %s",
			expectedJSON, resultJSON)
	}
}
