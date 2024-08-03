package config

import "github.com/segmentio/kafka-go"

const (
	kafkaTopic     = "temporal-topic"
	kafkaBrokerURL = "localhost:9092"
	groupID        = "consumer-group-1"
)

func NewKafkaWriter() *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaBrokerURL),
		Topic:    kafkaTopic,
		Balancer: &kafka.LeastBytes{},
	}
}

func NewKafkaReader() *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBrokerURL},
		Topic:   kafkaTopic,
		GroupID: groupID,
	})
}
