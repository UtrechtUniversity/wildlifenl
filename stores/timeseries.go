package stores

import (
	"context"
	"errors"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type Timeseries struct {
	client       influxdb2.Client
	organization string
}

func NewTimeseries(url string, organization string, token string) *Timeseries {
	t := Timeseries{
		client:       influxdb2.NewClient(url, token),
		organization: organization,
	}
	return &t
}

func (t *Timeseries) Ping() error {
	ok, err := t.client.Ping(context.Background())
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("server is down")
	}
	return nil
}

func (t *Timeseries) Writer(bucket string) api.WriteAPIBlocking {
	return t.client.WriteAPIBlocking(t.organization, bucket)
}

func (t *Timeseries) Reader() api.QueryAPI {
	return t.client.QueryAPI(t.organization)
}
