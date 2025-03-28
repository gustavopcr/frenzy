package payloadgen

import "github.com/getkin/kin-openapi/openapi3"

type GeneratorOptions func(*PayloadGenerator)

func WithGenerateString(fn func(schema *openapi3.Schema) any) GeneratorOptions {
	return func(pg *PayloadGenerator) {
		pg.generateString = fn
	}
}
func WithGenerateInteger(fn func(schema *openapi3.Schema) any) GeneratorOptions {
	return func(pg *PayloadGenerator) {
		pg.generateInteger = fn
	}
}
func WithGenerateNumber(fn func(schema *openapi3.Schema) any) GeneratorOptions {
	return func(pg *PayloadGenerator) {
		pg.generateNumber = fn
	}
}
func WithGenerateBoolean(fn func(schema *openapi3.Schema) any) GeneratorOptions {
	return func(pg *PayloadGenerator) {
		pg.generateBoolean = fn
	}
}

func WithGenerateEnum(fn func(enumValues []any) any) GeneratorOptions {
	return func(pg *PayloadGenerator) {
		pg.generateEnum = fn
	}
}
