package kafka

import (
	"github.com/Shopify/sarama"
)

// NewSyncProducer new sync producer
func NewSyncProducer(conf Config) (sarama.SyncProducer, error) {
	config := pre(conf)
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(conf.Broker, config)
	if err != nil {
		return nil, err
	}

	return producer, nil
}

// NewAsyncProducer new async producer
func NewAsyncProducer(conf Config) (sarama.AsyncProducer, error) {
	config := pre(conf)
	config.Producer.Return.Successes = conf.Sarama.Producer.Return.Successes

	producer, err := sarama.NewAsyncProducer(conf.Broker, config)
	if err != nil {
		return nil, err
	}

	return producer, nil
}
