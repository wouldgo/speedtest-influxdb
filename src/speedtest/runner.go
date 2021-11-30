package speedtest

import (
	"context"
	"time"

	"github.com/m-lab/ndt7-client-go"
	"github.com/m-lab/ndt7-client-go/spec"
	"github.com/wouldgo/speedtest-influxdb/src/options"
)

type Runner struct {
	Data           chan *spec.Measurement
	clientVersion  string
	clientName     string
	defaultTimeout time.Duration
	ctx            context.Context
	Dispose        func()
}

func New(options *options.Options) (*Runner, error) {
	clientName := *options.SpeedTestConfiguration.ClientName
	clientVersion := *options.SpeedTestConfiguration.ClientVersion
	defaultTimeout := *options.SpeedTestConfiguration.DefaultTimeout
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	runner := &Runner{
		Data: make(chan *spec.Measurement),
		Dispose: func() {
			cancel()
		},
		clientVersion:  clientVersion,
		clientName:     clientName,
		defaultTimeout: defaultTimeout,
		ctx:            ctx,
	}

	return runner, nil
}

func (r Runner) Run() (*string, map[spec.TestKind]*ndt7.LatestMeasurements, error) {
	client := ndt7.NewClient(r.clientName, r.clientVersion)

	var (
		channel <-chan spec.Measurement
		err     error
	)

	channel, err = client.StartDownload(r.ctx)
	if err != nil {
		return nil, nil, err
	}

	for ev := range channel {
		r.Data <- &ev
	}

	channel, err = client.StartUpload(r.ctx)
	if err != nil {
		return nil, nil, err
	}

	for ev := range channel {
		r.Data <- &ev
	}

	fqdn := client.FQDN
	return &fqdn, client.Results(), nil
}
