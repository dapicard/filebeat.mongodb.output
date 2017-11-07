package mongodb

import (
	"time"

	"gopkg.in/mgo.v2"

	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/outputs"
	"github.com/elastic/beats/libbeat/outputs/outil"
	"github.com/elastic/beats/libbeat/publisher"
)

type publishFn func(
	keys outil.Selector,
	data []publisher.Event,
) ([]publisher.Event, error)

type client struct {
	outputs.NetworkClient
	url			string
	stats    	*outputs.Stats
	db       	string
	collection  outil.Selector
	publish  	publishFn
	timeout  	time.Duration
}

type bulkInfo struct {
	data	[]publisher.Event
	bulk	*mgo.Bulk
}

func newClient(
	url	string,
	stats *outputs.Stats,
	timeout time.Duration,
	db string, 
	collection outil.Selector,
) *client {
	return &client{
		url:		url,
		stats:    	stats,
		timeout:  	timeout,
		db:       	db,
		collection: collection,
	}
}

func (c *client) Connect() error {
	debugf("connect to %s", c.url)
	session, err := mgo.DialWithTimeout(c.url, c.timeout)
	database := session.DB(c.db)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			session.Close()
		}
	}()

	c.publish, err = c.makePublish(database)

	return err
}

func (c *client) Publish(batch publisher.Batch) error {
	if c == nil {
		panic("no client")
	}
	if batch == nil {
		panic("no batch")
	}

	events := batch.Events()
	c.stats.NewBatch(len(events))
	rest, err := c.publish(c.collection, events)
	if rest != nil {
		c.stats.Failed(len(rest))
		batch.RetryEvents(rest)
		return err
	}

	batch.ACK()
	return err
}

func (c *client) makePublish(
	database *mgo.Database,
) (publishFn, error) {
	return func(collection outil.Selector, data []publisher.Event) ([]publisher.Event, error) {

		var bulks map[string]bulkInfo
		bulks = map[string]bulkInfo{}

		// data = okEvents[:0]
		dropped := 0
		for _, datum := range data {
			//Find where to push the event
			collection, err := collection.Select(&datum.Content)
			debugf("Event key: %s", collection)
			if err != nil {
				logp.Err("Failed to set mongodb key: %v", err)
				dropped++
				continue
			}
			bulk, exists := bulks[collection]
			if (!exists) {
				bulk = bulkInfo{
					bulk: database.C(collection).Bulk(),
					data: []publisher.Event{},
				}
				bulks[collection] = bulk
			}

			debugf("Add document to bulk: %s", datum)
			bulk.bulk.Insert(datum)
			bulk.data = append(bulk.data, datum)
		}
		c.stats.Dropped(dropped)

		failed := []publisher.Event{}
		var lastError error
		for collection, bulk := range bulks {
			debugf("Run %s bulk, containing %d documents", collection, len(bulk.data))
			if r, err := bulk.bulk.Run(); err != nil {
				logp.Warn("Bulk failed: %v", err)
				failed = append(failed, bulk.data...)
				lastError = err
			} else {
				debugf("Bulk return: %s", r)
				c.stats.Acked(len(bulk.data))
			}
		}

		return failed, lastError
	}, nil
}
