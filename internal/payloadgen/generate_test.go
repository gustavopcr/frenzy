package payloadgen

import (
	"encoding/json"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

func TestPayloadFromSchema(t *testing.T) {
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile("user_schema.yaml")
	if err != nil {
		t.Errorf("failed to load schema: %v", err)
	}

	if err := doc.Validate(loader.Context); err != nil {
		t.Errorf("Error validating with loader: %v", err)
	}

	userSchemaRef := doc.Components.Schemas["User"]
	userSchema := userSchemaRef.Value
	user := PayloadFromSchema(userSchema)

	expected := map[string]any{
		"id":   10,
		"name": "string",
	}

	resultJSON, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("failed to marshal generated payload: %v", err)
	}

	expectedJSON, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("failed to marshal expected payload: %v", err)
	}

	if string(resultJSON) != string(expectedJSON) {
		t.Errorf("Geenerated paylod does not match expected. \nExpected: %s\nGot: %s",
			expectedJSON, resultJSON)
	}
}
