package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gustavopcr/frenzy/internal/payloadgen"
)

func main() {
	gofakeit.Seed(0)

	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile("internal/payloadgen/user_schema.yaml")
	if err != nil {
		log.Fatalf("failed to load schema: %v", err)
	}

	if err := doc.Validate(loader.Context); err != nil {
		log.Fatalf("Error validating with loader: %v", err)
	}

	userSchemaRef := doc.Components.Schemas["User"]
	userSchema := userSchemaRef.Value

	pg := payloadgen.NewPayloadGenerator()
	teste := pg.PayloadFromSchema(userSchema)

	jsonBytes, err := json.MarshalIndent(teste, "", " ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsonBytes))

}
