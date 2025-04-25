package main

import (
	"fmt"
	"io"
	"log"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/getkin/kin-openapi/openapi3"

	"github.com/gustavopcr/frenzy/internal/http"
)

func main() {
	gofakeit.Seed(0)
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile("internal/payloadgen/spring_boot_app_schema.yaml")
	if err != nil {
		log.Fatalf("failed to load schema: %v", err)
	}

	if err := doc.Validate(loader.Context); err != nil {
		log.Fatalf("Error validating with loader: %v", err)
	}

	userSchemaRef := doc.Components.Schemas["User"]
	userSchema := userSchemaRef.Value

	resp, err := http.TestePost("http://localhost:8080/hello", userSchema)
	if err != nil {
		panic(err)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println("resp: ", string(bodyBytes))
}
