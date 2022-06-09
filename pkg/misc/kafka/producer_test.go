package kafka

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/Shopify/sarama"
)

func TestSyncProducer(t *testing.T) {
	const (
		topic = "VerificationCode"
	)

	conf := Config{
		Broker: []string{"192.168.200.20:9092", "192.168.200.19:9092", "192.168.200.18:9092"},
	}
	conf.Sarama.Version = sarama.V2_0_0_0
	producer, err := NewSyncProducer(conf)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer producer.Close()

	_, _, err = producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder("this is test message"),
	})

	if err != nil {
		t.Fatal(err)
		return
	}

}

func TestASyncProducer(t *testing.T) {
	const (
		topic = "test"
	)

	conf := Config{
		Broker: []string{"192.168.200.20:9092", "192.168.200.19:9092", "192.168.200.18:9092"},
	}
	conf.Sarama.Version = sarama.V2_0_0_0
	producer, err := NewAsyncProducer(conf)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer producer.Close()

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	wait := sync.WaitGroup{}
	wait.Add(1)
	go func(ctx context.Context, producer sarama.AsyncProducer, wait *sync.WaitGroup) {
		select {
		case <-producer.Successes():
		case <-producer.Errors():
		case <-ctx.Done():
		}
		wait.Done()
	}(ctx, producer, &wait)

	input := producer.Input()
	input <- &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder("this is test message(async)"),
	}

	wait.Wait()
}
