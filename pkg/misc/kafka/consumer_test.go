package kafka

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Shopify/sarama"
)

type demo struct{}

// Setup Setup
func (d *demo) Setup(_ sarama.ConsumerGroupSession) error {
	fmt.Println("Setup")
	return nil
}

// Cleanup Cleanup
func (d *demo) Cleanup(_ sarama.ConsumerGroupSession) error {
	fmt.Println("Cleanup")
	return nil
}

// ConsumeClaim ConsumeClaim
func (d *demo) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	fmt.Printf("ConsumeClaim\n")
	for msg := range claim.Messages() {
		fmt.Printf("Message topic:%q partition:%d offset:%d\n", msg.Topic, msg.Partition, msg.Offset)
		fmt.Printf("%s:%s\n", string(msg.Key), string(msg.Value))
		sess.MarkMessage(msg, "")
	}
	return nil
}
func TestConsumerGroup(t *testing.T) {
	const (
		topic = "test"
		group = "test"
	)

	conf := Config{
		Broker: []string{"192.168.200.20:9092", "192.168.200.19:9092", "192.168.200.18:9092"},
	}
	conf.Sarama.Version = sarama.V2_0_0_0
	consumer, err := NewConsumerGroup(conf, group)
	if err != nil {
		t.Fatal(err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	err = consumer.Consume(ctx, []string{topic}, new(demo))
	if err != nil {
		t.Fatal(err)
		return
	}
	select {
	case err := <-consumer.Errors():
		t.Fatal(err)
	case <-ctx.Done():

	}
}
