package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Config Config
type Config struct {
	Hosts      []string
	Direct     bool
	Credential struct {
		AuthMechanism           string
		AuthMechanismProperties map[string]string
		AuthSource              string
		Username                string
		Password                string
		PasswordSet             bool
	}
}

// New new
func New(conf *Config) (*mongo.Client, error) {
	clientOptions := options.Client().
		SetDirect(conf.Direct).
		SetHosts(conf.Hosts).SetAuth(options.Credential{
		AuthMechanism:           conf.Credential.AuthMechanism,
		AuthMechanismProperties: conf.Credential.AuthMechanismProperties,
		AuthSource:              conf.Credential.AuthSource,
		Username:                conf.Credential.Username,
		Password:                conf.Credential.Password,
		PasswordSet:             conf.Credential.PasswordSet,
	})

	return mongo.Connect(context.TODO(), clientOptions)
}
