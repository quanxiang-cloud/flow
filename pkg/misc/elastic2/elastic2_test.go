package elastic2

import (
	"context"
	"testing"
)

func TestElastic2(t *testing.T) {
	conf := &Config{
		Host: []string{"http://192.168.200.18:9200",
			"http://192.168.200.19:9200",
			"http://192.168.200.20:9200"},
	}

	client, err := NewClient(conf, nil)
	if err != nil {
		t.Fatal(err)
		return
	}

	mock := struct {
		ID       string `json:"id,omitempty"`
		UserID   string `json:"userID,omitempty"`
		UserName string `json:"userName,omitempty"`
	}{
		ID:       "1",
		UserID:   "1",
		UserName: "mock",
	}

	ctx := context.Background()
	_, err = client.Index().
		Index("mock").
		Type("mock").
		Id(mock.ID).
		BodyJson(mock).
		Do(ctx)
	if err != nil {
		t.Fatal(err)
		return
	}
}
