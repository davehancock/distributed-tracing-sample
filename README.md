# distributed-tracing-sample


<h2> Build / Run </h2>

env GOOS=linux GOARCH=386 go build -v -o ./consumer-three/build/app ./consumer-three

./gradlew clean build

docker-compose build && docker-compose up




<h2> Kafka Inspection of  </h2>

/opt/kafka_2.11-0.10.1.0/bin/kafka-topics.sh --list --zookeeper 127.0.0.1
__consumer_offsets
bar
sleuth


/opt/kafka_2.11-0.10.1.0/bin/kafka-console-consumer.sh --from-beginning --topic sleuth --zookeeper 127.0.0.1





Things to add:

- A Go service - using gokit? for tracing
- gliffy diagram

- A service that adds baggage
- A downstream service of a kafka stream consumer
