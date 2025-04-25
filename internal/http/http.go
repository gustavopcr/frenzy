package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gustavopcr/frenzy/internal/payloadgen"
)

type httpConfig struct {
	url     string
	Timeout time.Duration
	reqSeq  int
	headers map[string]string
}

type Result struct {
	Resp *http.Response
	err  error
}

type httpHandler struct {
	config  *httpConfig
	client  *http.Client
	Results chan Result
	Wg      *sync.WaitGroup
}

func NewHttpHandler() *httpHandler {
	httpConfig := &httpConfig{url: "http://localhost:8080/hello"}
	httpClient := &http.Client{}
	resultsChan := make(chan Result)

	return &httpHandler{
		config:  httpConfig,
		client:  httpClient,
		Results: resultsChan,
		Wg:      &sync.WaitGroup{},
	}
}

func (httpHandler *httpHandler) TesteGet() {
	defer httpHandler.Wg.Done()
	//TODO add query string later on
	resp, err := httpHandler.client.Get(httpHandler.config.url)
	if err != nil {
		httpHandler.Results <- Result{nil, err}

	}
	httpHandler.Results <- Result{resp, nil}

}

func (httpHandler *httpHandler) TestePost(schema *openapi3.Schema) {
	defer httpHandler.Wg.Done()
	pg := payloadgen.NewPayloadGenerator()
	teste := pg.PayloadFromSchema(schema)

	jsonBytes, err := json.MarshalIndent(teste, "", " ")
	if err != nil {
		panic(err)
	}

	body := bytes.NewReader(jsonBytes)
	resp, err := httpHandler.client.Post(httpHandler.config.url, "application/json", body)
	if err != nil {
		httpHandler.Results <- Result{nil, err}
	}
	httpHandler.Results <- Result{resp, nil}
}
