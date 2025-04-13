package payloadgen

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

func customPayloadGenerator() *PayloadGenerator {
	return NewPayloadGenerator(
		WithGenerateInteger(func(schema *openapi3.Schema) any {
			return 10
		}),
		WithGenerateNumber(func(schema *openapi3.Schema) any {
			return 3.14
		}),
		WithGenerateBoolean(func(schema *openapi3.Schema) any {
			return true
		}),
		WithGenerateEnum(func(enumValues []any) any {
			return enumValues[0]
		}),
	) // Overwriting default enumSelector to make testing predictable
}

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

	pg := NewPayloadGenerator(
		WithGenerateString(func(schema *openapi3.Schema) any {
			return "string"
		}),
		WithGenerateInteger(func(schema *openapi3.Schema) any {
			return 10
		}),
		WithGenerateNumber(func(schema *openapi3.Schema) any {
			return 3.14
		}),
	)

	user := pg.PayloadFromSchema(userSchema)

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
		t.Errorf("Generated payload does not match expected. \nExpected: %s\nGot: %s",
			expectedJSON, resultJSON)
	}
}

func TestPayloadFromSchemaEnum(t *testing.T) {
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile("user_schema_enum.yaml")
	if err != nil {
		t.Errorf("failed to load schema: %v", err)
	}

	if err := doc.Validate(loader.Context); err != nil {
		t.Errorf("Error validating with loader: %v", err)
	}

	userSchemaRef := doc.Components.Schemas["User"]
	userSchema := userSchemaRef.Value

	pg := NewPayloadGenerator(
		WithGenerateString(func(schema *openapi3.Schema) any {
			return "string"
		}),
		WithGenerateInteger(func(schema *openapi3.Schema) any {
			return 10
		}),
		WithGenerateNumber(func(schema *openapi3.Schema) any {
			return 3.14
		}),
		WithGenerateEnum(func(enumValues []any) any {
			return enumValues[0]
		}),
	) // Overwriting default enumSelector to make testing predictable

	user := pg.PayloadFromSchema(userSchema)

	expected := map[string]any{
		"id":       10,
		"name":     "string",
		"userType": "admin",
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
		t.Errorf("Generated payload does not match expected. \nExpected: %s\nGot: %s",
			expectedJSON, resultJSON)
	}
}

func TestPayloadFromSchemaComplete(t *testing.T) { //TODO: finish this test
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile("user_schema_complete.yaml")
	if err != nil {
		t.Errorf("failed to load schema: %v", err)
	}

	if err := doc.Validate(loader.Context); err != nil {
		t.Errorf("Error validating with loader: %v", err)
	}

	userSchemaRef := doc.Components.Schemas["User"]
	userSchema := userSchemaRef.Value
	pg := customPayloadGenerator()
	user := pg.PayloadFromSchema(userSchema)
	if user != nil {
		fmt.Println("name: ", user)
	}
}
