package mongo

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestNew(t *testing.T) {
	client, err := New(&Config{
		Direct: true,
		Hosts:  []string{"192.168.200.19:27017"},
		Credential: struct {
			AuthMechanism           string
			AuthMechanismProperties map[string]string
			AuthSource              string
			Username                string
			Password                string
			PasswordSet             bool
		}{
			AuthMechanism: "SCRAM-SHA-1",
			AuthSource:    "admin",
			Username:      "root",
			Password:      "uyWxtvt6gCOy3VPLB3rTpa0rQ",
			PasswordSet:   false,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	type demo struct {
		ID   string
		Name string
	}

	db := client.Database("test")
	ans, err := db.Collection("test").InsertOne(context.TODO(), &demo{
		ID:   "1",
		Name: "haha",
	})

	if err != nil {
		t.Fatal(err)
	}

	mock := new(demo)
	result := db.Collection("test").FindOne(context.TODO(), bson.M{"_id": ans.InsertedID})

	err = result.Decode(mock)
	if err != nil {
		t.Fatal(err)
	}
}
