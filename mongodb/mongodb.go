package mongodb

import (
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/outputs"
	"github.com/elastic/beats/libbeat/outputs/outil"
)

type mongodbOut struct {
	beat beat.Info
}

var debugf = logp.MakeDebug("mongodb")


func init() {
	outputs.RegisterType("mongodb", makeMongodb)
}

func makeMongodb(
	beat beat.Info,
	stats *outputs.Stats,
	cfg *common.Config,
) (outputs.Group, error) {
	config := defaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return outputs.Fail(err)
	}

	if !cfg.HasField("collection") {
		cfg.SetString("collection", -1, beat.Beat)
	}

	collection, err := outil.BuildSelectorFromConfig(cfg, outil.Settings{
		Key:              "collection",
		MultiKey:         "collections",
		EnableSingleOnly: true,
		FailEmpty:        true,
	})
	if err != nil {
		return outputs.Fail(err)
	}

	hosts, err := outputs.ReadHostList(cfg)
	if err != nil {
		return outputs.Fail(err)
	}

	clients := make([]outputs.NetworkClient, len(hosts))
	for i, host := range hosts {
		clients[i] = newClient(host, stats, config.Timeout, config.Db, collection)
	}

	return outputs.SuccessNet(config.LoadBalance, config.BulkMaxSize, config.MaxRetries, clients)
}
