# distributed-tracing-sample


./gradlew clean build

docker-compose up

/opt/kafka_2.11-0.10.1.0/bin/kafka-topics.sh --list --zookeeper 127.0.0.1
__consumer_offsets
bar
sleuth



/opt/kafka_2.11-0.10.1.0/bin/kafka-console-consumer.sh --from-beginning --topic sleuth --zookeeper 127.0.0.1


