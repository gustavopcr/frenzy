package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/getkin/kin-openapi/openapi3"
)

func main() {
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

	teste := generatePayloadFromSchema(userSchema)

	fmt.Println("teste: ", teste)

	jsonBytes, err := json.MarshalIndent(teste, "", " ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsonBytes))

}

func generatePayloadFromSchema(schema *openapi3.Schema) any {
	if schema == nil || schema.Type == nil {
		return nil
	}
	switch {
	case schema.Type.Is("object"):
		obj := make(map[string]any)
		for propName, propSchemaRef := range schema.Properties {
			propSchema := propSchemaRef.Value
			obj[propName] = generatePayloadFromSchema(propSchema)
		}
		return obj
	case schema.Type.Is("array"):
		if schema.Items == nil || schema.Items.Value == nil {
			return []any{}
		}
		item := generatePayloadFromSchema(schema.Items.Value)
		return []any{item}
	case schema.Type.Is("string"):
		return "string"
	case schema.Type.Is("integer"):
		return 10
	case schema.Type.Is("number"):
		return 3.44
	case schema.Type.Is("boolean"):
		return true
	default:
		return nil
	}

}
