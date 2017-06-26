package main

import (
	"net/http"
	"fmt"
	"io/ioutil"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"

	httptransport "github.com/go-kit/kit/transport/http"
	stdopentracing "github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	"github.com/go-kit/kit/endpoint"
	"context"
	"os"
)

const port = ":8080"

func main() {

	configureService()

	logger := log.NewLogfmtLogger(os.Stderr)

	logger.Log("Starting consumer-three... on port " + port)

	svc := fooService{}

	fooHandler := httptransport.NewServer(
		makeFooEndpoint(svc),
		decodeFooRequest,
		nil,
		append(options, httptransport.ServerBefore(opentracing.FromHTTPRequest(tracer, "gofoo", logger)))...,
	)

	http.Handle("/foo", fooHandler)
	fmt.Println(http.ListenAndServe(":8080", nil))
}

func decodeFooRequest(_ context.Context, r *http.Request) (interface{}, error) {

	var request fooRequest

	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	request = string(payload);

	return request, nil
}

// Some service definition
type FooService interface {
	ConsumeMessage(string) (string, error)
}

// Some service implementation
type fooService struct{}

func (fooService) ConsumeMessage(message string) {
	fmt.Println(message)
}

// Some Request
type fooRequest struct {
	S string
}

// Some Response
type fooResponse struct {
	S string
}

// Some Endpoint
func makeFooEndpoint(service FooService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(fooRequest)
		v, err := service.ConsumeMessage(req.S)
		if err != nil {
			return fooResponse{v, err.Error()}, nil
		}
		return fooResponse{v, ""}, nil
	}
}

func handleStuff(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Print("Gophered - " + string(payload))
}

func MakeHTTPHandler(tracer stdopentracing.Tracer, logger log.Logger) http.Handler {

	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
	}

	m := http.NewServeMux()
	m.Handle("/gofoo", httptransport.NewServer(
		MakeSumEndpoint,
		DecodeHTTPSumRequest,
		EncodeHTTPGenericResponse,
		append(options, httptransport.ServerBefore(opentracing.FromHTTPRequest(tracer, "gofoo", logger)))...,
	))

	return m
}



// Setup / Configuration for the service, i.e Tracing in this case

func configureService(logger log.Logger) {

	// Tracing domain.
	var (
		debugMode = false
		serviceName = "consumer-three"
		serviceHostPort = "localhost:8080"
		zipkinHTTPEndpoint = "zipkin-server:8080"
	)

	collector, err := zipkin.NewHTTPCollector(zipkinHTTPEndpoint)
	if err != nil {
		logger.Log("Error initializing zipkin collector", err)
	}

	tracer, err := zipkin.NewTracer(
		zipkin.NewRecorder(collector, debugMode, serviceHostPort, serviceName),
	)
	if err != nil {
		logger.Log("Error initializing zipkin recorder", err)
	}

	// HTTP transport.
	go func() {
		logger := log.With(logger, "transport", "HTTP")
		h := MakeHTTPHandler(endpoints, tracer, logger)
		logger.Log("addr", *httpAddr)
		errc <- http.ListenAndServe(*httpAddr, h)
	}()

}

