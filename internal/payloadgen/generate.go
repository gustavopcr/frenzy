package payloadgen

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/getkin/kin-openapi/openapi3"
)

type PayloadGenerator struct {
	generators map[string]func(*openapi3.Schema) any
	enumSelect func(enumValues []any) any
}

func NewPayloadGenerator(enumSelector func(enumValues []any) any) *PayloadGenerator {
	if enumSelector == nil {
		enumSelector = func(enumValues []any) any {
			if len(enumValues) > 0 {
				index := gofakeit.Number(0, len(enumValues)-1)
				return enumValues[index]
			}
			return nil
		}
	}
	pg := PayloadGenerator{
		enumSelect: enumSelector,
	}

	pg.generators = map[string]func(*openapi3.Schema) any{
		"string":  pg.generateString,
		"integer": pg.generateInteger,
		"number":  pg.generateNumber,
		"boolean": pg.generateBoolean,
		"array":   pg.generateArray,
		"object":  pg.generateObject,
	}
	return &pg
}

func (pg *PayloadGenerator) generateString(schema *openapi3.Schema) any {
	return "string"
}

func (pg *PayloadGenerator) generateInteger(schema *openapi3.Schema) any {
	return 10
}

func (pg *PayloadGenerator) generateNumber(schema *openapi3.Schema) any {
	return 3.14
}

func (pg *PayloadGenerator) generateBoolean(schema *openapi3.Schema) any {
	return true
}

// returns []any
func (pg *PayloadGenerator) generateArray(schema *openapi3.Schema) any {
	if schema.Items == nil || schema.Items.Value == nil {
		return []any{}
	}
	item := pg.PayloadFromSchema(schema.Items.Value)
	return []any{item}
}

// returns map[string]any
func (pg *PayloadGenerator) generateObject(schema *openapi3.Schema) any {
	obj := make(map[string]any)
	for propName, propSchemaRef := range schema.Properties {
		propSchema := propSchemaRef.Value
		obj[propName] = pg.PayloadFromSchema(propSchema)
	}
	return obj
}

func (pg *PayloadGenerator) PayloadFromSchema(schema *openapi3.Schema) any {
	if schema == nil || schema.Type == nil {
		return nil
	}

	if len(schema.Enum) > 0 {
		return pg.enumSelect(schema.Enum)
	}
	switch {
	case schema.Type.Is("object"):
		return pg.generators["object"](schema)
	case schema.Type.Is("array"):
		return pg.generators["array"](schema)

	case schema.Type.Is("string"):
		return pg.generators["string"](schema)

	case schema.Type.Is("integer"):
		return pg.generators["integer"](schema)

	case schema.Type.Is("number"):
		return pg.generators["number"](schema)

	case schema.Type.Is("boolean"):
		return pg.generators["boolean"](schema)
	default:
		return nil
	}

}
