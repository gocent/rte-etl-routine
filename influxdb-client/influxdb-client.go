package influxdb_client

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"rte-etl-routine/config"
)

type InfluxdbClient struct{}

func NewClient() influxdb2.Client {
	url := config.GetEnv().Influxdb.Host + ":" + config.GetEnv().Influxdb.Port
	return influxdb2.NewClient(url, config.GetEnv().Influxdb.Token)
}
