package timeseries

import (
	"context"
	"errors"
	"fmt"

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

func (t *Timeseries) CreateBucketIfNotExists(name string) error {
	buckets, err := t.client.BucketsAPI().FindBucketsByOrgName(context.Background(), t.organization)
	if err != nil {
		return err
	}
	bucketExists := false
	for _, bucket := range *buckets {
		if bucket.Name == name {
			bucketExists = true
			break
		}
	}
	if !bucketExists {
		// TODO: Does not find organization???????
		organization, err := t.client.OrganizationsAPI().FindOrganizationByName(context.Background(), t.organization)
		if err != nil {
			return err
		}
		if organization == nil {
			return fmt.Errorf("organization %v not found", t.organization)
		}
		if _, err := t.client.BucketsAPI().CreateBucketWithName(context.Background(), organization, name); err != nil {
			return err
		}
	}
	return nil
}
