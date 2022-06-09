package kafka

import (
	"github.com/Shopify/sarama"
)

// NewConsumerGroup new kafka consumer group
func NewConsumerGroup(conf Config, group string) (sarama.ConsumerGroup, error) {
	config := pre(conf)

	consumer, err := sarama.NewConsumerGroup(conf.Broker, group, config)
	if err != nil {
		return nil, err
	}

	return consumer, nil
}

// NewConsumer new kafka consumer
func NewConsumer(conf Config) (sarama.Consumer, error) {
	config := pre(conf)

	consumer, err := sarama.NewConsumer(conf.Broker, config)
	if err != nil {
		return nil, err
	}

	return consumer, nil
}
