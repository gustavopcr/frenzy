package main

import (
	"fmt"
	"io"
	"log"
	"time"

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

	httpHandler := http.NewHttpHandler()

	go httpHandler.TestePost(userSchema)
	go httpHandler.TestePost(userSchema)
	go httpHandler.TesteGet()

	for i := 0; i < 3; i++ {
		select {
		case res, ok := <-httpHandler.Results:
			if !ok {
				return
			}
			if res.Err != nil {
				fmt.Println("ERROR FOUND")
				continue
			}
			fmt.Println("alo")
			jsonBytes, err := io.ReadAll(res.Resp.Body)
			if err != nil {
				panic(err)
			}
			defer res.Resp.Body.Close()
			fmt.Println(string(jsonBytes))
		case <-time.After(5 * time.Second):
			fmt.Println("Timeout waiting for response")
		}
	}

}
