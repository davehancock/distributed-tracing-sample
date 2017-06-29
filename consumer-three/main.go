package main

import (
	"net/http"
	"fmt"
	"io/ioutil"

	"github.com/go-kit/kit/tracing/opentracing"

	httptransport "github.com/go-kit/kit/transport/http"
	//stdopentracing "github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	"github.com/go-kit/kit/endpoint"
	"context"
	"errors"
	"os"
	"github.com/go-kit/kit/log"
)

const port = ":8081"

func main() {

	// Tracing domain.
	var (
		debugMode = false
		serviceName = "consumer-three"
		serviceHostPort = "localhost:8081"
		zipkinHTTPEndpoint = "http://localhost:9100/api/v1/spans"
	)

	// Logging domain.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	logger = log.With(logger, "tracer", "ZipkinHTTP")
	logger.Log("addr", zipkinHTTPEndpoint)


	// endpoint typically looks like: http://zipkinhost:9411/api/v1/spans
	collector, err := zipkin.NewHTTPCollector(zipkinHTTPEndpoint)
	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}
	defer collector.Close()

	tracer, err := zipkin.NewTracer(
		zipkin.NewRecorder(collector, debugMode, serviceHostPort, serviceName),
	)
	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}

	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
	}

	svc := fooService{}


	fooHandler := httptransport.NewServer(
		makeFooEndpoint(svc),
		decodeFooRequest,
		encodeFooRequest,
		append(options, httptransport.ServerBefore(opentracing.FromHTTPRequest(tracer, "gofoo", logger)))...
	)

	http.Handle("/foo", handleForMethod(http.MethodPut, fooHandler))
	fmt.Println(http.ListenAndServe(port, nil))
}

// Just some HTTP Method verification, delegating to the underlying function if applicable
func handleForMethod(method string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == method {
			h.ServeHTTP(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func decodeFooRequest(_ context.Context, r *http.Request) (interface{}, error) {

	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	return fooRequest{string(payload)}, nil
}

func encodeFooRequest(_ context.Context, w http.ResponseWriter, response interface{}) error {

	r := response.(fooResponse)
	_, err := w.Write([]byte(r.payload))
	if err != nil {
		return errors.New("Error writing response: " + err.Error())
	}

	return nil
}

// Some service definition
type FooService interface {
	ConsumeMessage(string) (string, error)
}

// Some service implementation
type fooService struct{}

func (fooService) ConsumeMessage(message string) (string, error) {
	return "Gophered: " + message, nil
}

// Some Request
type fooRequest struct {
	payload string
}

// Some Response
type fooResponse struct {
	payload string
}

// Some Endpoint
func makeFooEndpoint(service FooService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(fooRequest)

		v, err := service.ConsumeMessage(req.payload)
		if err != nil {
			return nil, errors.New("Error Consuming message")
		}
		return fooResponse{v}, nil
	}
}
