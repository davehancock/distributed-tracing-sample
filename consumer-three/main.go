package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"log"
)

const port = ":8080"

func main() {

	//configureService()

	fmt.Println("Starting consumer-three... on port " + port)

	http.HandleFunc("/gofoo", handleStuff)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
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

	log.Print("Gophered - " + string(payload))
}


// Setup / Configuration for the service, i.e Tracing in this case

//func configureService(){
//
//	var (
//		zipkinAddr = flag.String("zipkin.addr", "", "Enable Zipkin tracing via a Zipkin HTTP Collector endpoint")
//		zipkinKafkaAddr = flag.String("zipkin.kafka.addr", "", "Enable Zipkin tracing via a Kafka server host:port")
//	)
//	flag.Parse()
//
//	// Tracing domain.
//	var tracer stdopentracing.Tracer
//	{
//		if *zipkinAddr != "" {
//			logger := log.With(logger, "tracer", "ZipkinHTTP")
//			logger.Log("addr", *zipkinAddr)
//
//			// endpoint typically looks like: http://zipkinhost:9411/api/v1/spans
//			collector, err := zipkin.NewHTTPCollector(*zipkinAddr)
//			if err != nil {
//				logger.Log("err", err)
//				os.Exit(1)
//			}
//			defer collector.Close()
//
//			tracer, err = zipkin.NewTracer(
//				zipkin.NewRecorder(collector, false, "localhost:80", "addsvc"),
//			)
//			if err != nil {
//				logger.Log("err", err)
//				os.Exit(1)
//			}
//		}
//	}
//
//}