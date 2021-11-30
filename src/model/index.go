package model

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/wouldgo/speedtest-influxdb/src/options"
	"github.com/wouldgo/speedtest-influxdb/src/speedtest"
	"go.uber.org/zap"
)

type Model struct {
	writeApi api.WriteAPI
	client   influxdb2.Client
	log      *zap.SugaredLogger
}

func (model *Model) Dispose() error {
	model.writeApi.Flush()
	model.client.Close()

	return nil
}

func New(options *options.Options) (*Model, error) {
	influxdbConfigurations := options.InfluxDb
	httpClient := &http.Client{
		Timeout: 1 * time.Minute,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout: 5 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
		},
	}
	influxDbOptions := influxdb2.DefaultOptions().SetHTTPClient(httpClient)

	client := influxdb2.NewClientWithOptions(*influxdbConfigurations.Url, *influxdbConfigurations.Token, influxDbOptions)
	writeApi := client.WriteAPI(*influxdbConfigurations.Org, *influxdbConfigurations.Bucket)

	toReturn := &Model{
		writeApi: writeApi,
		client:   client,
		log:      options.Log,
	}

	return toReturn, nil
}

func (m Model) Write(aData *speedtest.Summary) {
	now := time.Now()

	m.log.Info("Writing data to InfluxDB:", now)
	bandwidthDownloadPoint := influxdb2.NewPointWithMeasurement("download").
		AddTag("unit", aData.Download.Unit).
		AddField("value", aData.Download.Value).
		SetTime(now)

	bandwidthUploadPoint := influxdb2.NewPointWithMeasurement("upload").
		AddTag("unit", aData.Upload.Unit).
		AddField("value", aData.Upload.Value).
		SetTime(now)

	bandwidthDownloadRetransmissionPoint := influxdb2.NewPointWithMeasurement("download-retransmission").
		AddTag("unit", aData.DownloadRetrans.Unit).
		AddField("value", aData.DownloadRetrans.Value).
		SetTime(now)

	bandwidthMinRoundTripTimePoint := influxdb2.NewPointWithMeasurement("min-round-trip-time").
		AddTag("unit", aData.MinRTT.Unit).
		AddField("value", aData.MinRTT.Value).
		SetTime(now)

	m.writeApi.WritePoint(bandwidthDownloadPoint)
	m.writeApi.WritePoint(bandwidthUploadPoint)
	m.writeApi.WritePoint(bandwidthDownloadRetransmissionPoint)
	m.writeApi.WritePoint(bandwidthMinRoundTripTimePoint)

	m.log.Info("Data written to InfluxDB:", now)
}
