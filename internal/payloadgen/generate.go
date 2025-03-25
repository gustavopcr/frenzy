package payloadgen

import (
	"github.com/getkin/kin-openapi/openapi3"
)

var typeGenerators = map[string]func(*openapi3.Schema) any{}

func init() {
	typeGenerators = map[string]func(*openapi3.Schema) any{
		"string":  wrap(generateString),
		"integer": wrap(generateInteger),
		"number":  wrap(generateNumber),
		"boolean": wrap(generateBoolean),
		"array":   wrap(generateArray),
		"object":  wrap(generateObject),
	}
}
func wrap[T any](fn func(*openapi3.Schema) T) func(*openapi3.Schema) any {
	return func(s *openapi3.Schema) any {
		return fn(s)
	}
}

func generateString(schema *openapi3.Schema) string {
	return "string"
}

func generateInteger(schema *openapi3.Schema) int {
	return 10
}

func generateNumber(schema *openapi3.Schema) float64 {
	return 3.14
}

func generateBoolean(schema *openapi3.Schema) bool {
	return true
}

func generateArray(schema *openapi3.Schema) []any {
	if schema.Items == nil || schema.Items.Value == nil {
		return []any{}
	}
	item := PayloadFromSchema(schema.Items.Value)
	return []any{item}
}

func generateObject(schema *openapi3.Schema) map[string]any {
	obj := make(map[string]any)
	for propName, propSchemaRef := range schema.Properties {
		propSchema := propSchemaRef.Value
		obj[propName] = PayloadFromSchema(propSchema)
	}
	return obj
}

func PayloadFromSchema(schema *openapi3.Schema) any {
	if schema == nil || schema.Type == nil {
		return nil
	}

	switch {
	case schema.Type.Is("object"):
		return typeGenerators["object"](schema)
	case schema.Type.Is("array"):
		return typeGenerators["array"](schema)

	case schema.Type.Is("string"):
		return typeGenerators["string"](schema)

	case schema.Type.Is("integer"):
		return typeGenerators["integer"](schema)

	case schema.Type.Is("number"):
		return typeGenerators["number"](schema)

	case schema.Type.Is("boolean"):
		return typeGenerators["boolean"](schema)
	default:
		return nil
	}

}
