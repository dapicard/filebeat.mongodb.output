package mongodb

import (
	"time"
)

type mongoDbConfig struct {
	Collection  string                `config:"collection"`
	LoadBalance bool                  `config:"loadbalance"`
	Timeout     time.Duration         `config:"timeout"`
	BulkMaxSize int                   `config:"bulk_max_size"`
	MaxRetries  int                   `config:"max_retries"`
	Db          string                `config:"db"`
}

var (
	defaultConfig = mongoDbConfig{
		LoadBalance: true,
		Timeout:     5 * time.Second,
		BulkMaxSize: 2048,
		MaxRetries:  3,
		Db:          "test",
		Collection:    "lines",
	}
)
