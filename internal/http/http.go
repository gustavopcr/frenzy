package http

import (
	"bytes"
	"encoding/json"
	"math"
	"net/http"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gustavopcr/frenzy/internal/payloadgen"
)

type HandlerOption func(*httpHandler)

type reqConfig struct {
	reqSeq   int
	duration time.Duration
	timeout  time.Duration
}
type httpConfig struct {
	url       string
	reqConfig reqConfig
	headers   map[string]string
}

func WithUrl(url string) HandlerOption {
	return func(handler *httpHandler) {
		handler.config.url = url
	}
}

func WithReqConfig(reqConfig reqConfig) HandlerOption {
	return func(handler *httpHandler) {
		if reqConfig.duration <= 0 {
			reqConfig.duration = 10
			//panic("reqConfig.duration must be > 0")
		}
		if reqConfig.reqSeq <= 0 {
			reqConfig.reqSeq = 5
			//panic("reqConfig.reqSeq must be > 0")
		}
		if reqConfig.timeout <= 0 {
			reqConfig.timeout = 5 * time.Second // default timeout
		}
	}
}

func WithHeaders(headers map[string]string) HandlerOption {
	return func(handler *httpHandler) {
		handler.config.headers = headers
	}
}

type Result struct {
	Resp *http.Response
	Err  error
}

type httpHandler struct {
	config  *httpConfig
	pg      *payloadgen.PayloadGenerator
	client  *http.Client
	Results chan Result
}

func WithPayloadGenerator(pg *payloadgen.PayloadGenerator) HandlerOption {
	return func(handler *httpHandler) {
		handler.pg = pg
	}
}

func WithHttpClient(client *http.Client) HandlerOption {
	return func(handler *httpHandler) {
		handler.client = client
	}
}

func NewHttpHandler_() *httpHandler {
	httpConfig := &httpConfig{}
	httpClient := &http.Client{}
	resultsChan := make(chan Result)

	return &httpHandler{
		config:  httpConfig,
		pg:      payloadgen.NewPayloadGenerator(),
		client:  httpClient,
		Results: resultsChan,
	}
}

func NewHttpHandler(handlerOptions ...HandlerOption) *httpHandler {
	httpConfig := &httpConfig{url: "http://localhost:8080/hello"}
	httpClient := &http.Client{}

	httpHandler := &httpHandler{
		config: httpConfig,
		pg:     payloadgen.NewPayloadGenerator(),
		client: httpClient,
	}

	for _, opt := range handlerOptions {
		opt(httpHandler)
	}

	totalTime := int(math.Ceil(httpHandler.config.reqConfig.duration.Seconds())) * httpHandler.config.reqConfig.reqSeq
	httpHandler.Results = make(chan Result, totalTime)
	return httpHandler
}

func (httpHandler *httpHandler) TesteGet() {
	//TODO add query string later on
	resp, err := httpHandler.client.Get(httpHandler.config.url)
	if err != nil {
		httpHandler.Results <- Result{nil, err}
		return
	}
	httpHandler.Results <- Result{resp, nil}

}

func (httpHandler *httpHandler) TestePost(schema *openapi3.Schema) {
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
		return
	}
	httpHandler.Results <- Result{resp, nil}
}
