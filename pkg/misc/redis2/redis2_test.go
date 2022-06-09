package redis2

import (
	"context"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	conf := Config{
		Addrs: []string{"192.168.200.18:6379", "192.168.200.19:6379", "192.168.200.20:6379"},
	}
	client, err := NewClient(conf)
	if err != nil {
		t.Fatal(err)
		return
	}
	err = client.Set(context.TODO(), "AA", "BB", time.Second*10).Err()
	if err != nil {
		t.Fatal(err)
		return
	}
	val, err := client.Get(context.TODO(), "AA").Result()
	if err != nil {
		t.Fatal(err)
		return
	}
	if val != "BB" {
		t.Fatal("val not equal")
	}
}
