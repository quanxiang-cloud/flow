package redis2

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"time"

	"github.com/go-redis/redis/v8"
)

// Config config
type Config struct {
	Addrs           []string
	Username        string
	Password        string
	MaxRetries      int
	MinRetryBackoff time.Duration
	MaxRetryBackoff time.Duration

	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	PoolSize           int
	MinIdleConns       int
	MaxConnAge         time.Duration
	PoolTimeout        time.Duration
	IdleTimeout        time.Duration
	IdleCheckFrequency time.Duration

	TLS *TLS
}

// TLS tls
type TLS struct {
	ClientCertFile string
	clientKeyFile  string
	CACertFile     string
}

// NewClient new redis cluster client
func NewClient(conf Config) (*redis.ClusterClient, error) {
	var (
		tlsConfig *tls.Config
		err       error
	)
	if conf.TLS != nil {
		tlsConfig, err = NewTLSConfig(conf.TLS.ClientCertFile, conf.TLS.clientKeyFile, conf.TLS.CACertFile)
	}
	if err != nil {
		return nil, err
	}

	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:              conf.Addrs,
		Username:           conf.Username,
		Password:           conf.Password,
		MaxRetries:         conf.MaxRetries,
		MinRetryBackoff:    conf.MinRetryBackoff,
		MaxRetryBackoff:    conf.MaxRetryBackoff,
		DialTimeout:        conf.DialTimeout,
		ReadTimeout:        conf.ReadTimeout,
		WriteTimeout:       conf.WriteTimeout,
		PoolSize:           conf.PoolSize,
		MinIdleConns:       conf.MinIdleConns,
		MaxConnAge:         conf.MaxConnAge,
		PoolTimeout:        conf.PoolTimeout,
		IdleTimeout:        conf.IdleTimeout,
		IdleCheckFrequency: conf.IdleCheckFrequency,
		TLSConfig:          tlsConfig,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return client, nil
}

// NewTLSConfig generates a TLS configuration used to authenticate on server with
// certificates.
// Parameters are the three pem files path we need to authenticate: client cert, client key and CA cert.
func NewTLSConfig(clientCertFile, clientKeyFile, caCertFile string) (*tls.Config, error) {
	tlsConfig := tls.Config{}

	// Load client cert
	cert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
	if err != nil {
		return &tlsConfig, err
	}
	tlsConfig.Certificates = []tls.Certificate{cert}

	// Load CA cert
	caCert, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		return &tlsConfig, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tlsConfig.RootCAs = caCertPool

	tlsConfig.BuildNameToCertificate()
	return &tlsConfig, err
}
