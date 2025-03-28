package payloadgen

import (
	"log"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/getkin/kin-openapi/openapi3"
)

type PayloadGenerator struct {
	generateString  func(schema *openapi3.Schema) any
	generateInteger func(schema *openapi3.Schema) any
	generateNumber  func(schema *openapi3.Schema) any
	generateBoolean func(schema *openapi3.Schema) any
	generateEnum    func(enumValues []any) any
}

func NewPayloadGenerator(opts ...GeneratorOptions) *PayloadGenerator {
	pg := &PayloadGenerator{
		generateString:  defaultGenerateString,
		generateInteger: defaultGenerateInteger,
		generateNumber:  defaultGenerateNumber,
		generateBoolean: defaultGenerateBoolean,
		generateEnum: func(enumValues []any) any {
			if len(enumValues) > 0 {
				index := gofakeit.Number(0, len(enumValues)-1)
				return enumValues[index]
			}
			return nil
		},
	}

	for _, opt := range opts {
		if opt != nil {
			opt(pg)
		}
	}

	return pg
}

func defaultGenerateString(schema *openapi3.Schema) any {
	switch schema.Format {
	case "uuid":
		return gofakeit.UUID()
	case "email":
		return gofakeit.Email()
	case "uri":
		return gofakeit.URL()
	case "hostname":
		return gofakeit.DomainName()
	case "ipv4":
		return gofakeit.IPv4Address()
	case "ipv6":
		return gofakeit.IPv6Address()
	case "date-time":
		return gofakeit.Date().Format(time.RFC3339)
	case "phone":
		return gofakeit.Phone()
	default:
		min := int(schema.MinLength)
		max := gofakeit.Number(min, min+10) // default value in case schema.MaxLength wasn't setted
		if schema.MaxLength != nil && *schema.MaxLength > schema.MinLength {
			max = int(*schema.MaxLength)
		}
		length := gofakeit.Number(min, max)
		str, err := gofakeit.Generate(strings.Repeat("?", length))
		if err != nil {
			log.Fatalf("failed to generate random string: %v", err)
		}
		return str
	}
}

func defaultGenerateInteger(schema *openapi3.Schema) any {
	return 10
}

func defaultGenerateNumber(schema *openapi3.Schema) any {
	return 3.14
}

func defaultGenerateBoolean(schema *openapi3.Schema) any {
	return true
}

func (pg *PayloadGenerator) generateArray(schema *openapi3.Schema) []any {
	if schema.Items == nil || schema.Items.Value == nil {
		return []any{}
	}
	item := pg.PayloadFromSchema(schema.Items.Value)
	return []any{item}
}

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
		return pg.generateEnum(schema.Enum)
	}
	switch {
	case schema.Type.Is("object"):
		return pg.generateObject(schema)
	case schema.Type.Is("array"):
		return pg.generateArray(schema)

	case schema.Type.Is("string"):
		return pg.generateString(schema)

	case schema.Type.Is("integer"):
		return pg.generateInteger(schema)

	case schema.Type.Is("number"):
		return pg.generateNumber(schema)

	case schema.Type.Is("boolean"):
		return pg.generateBoolean(schema)
	default:
		return nil
	}

}
