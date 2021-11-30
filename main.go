package main

import (
	"github.com/m-lab/ndt7-client-go/spec"
	"go.uber.org/zap"

	"github.com/wouldgo/speedtest-influxdb/src/model"
	cliOptions "github.com/wouldgo/speedtest-influxdb/src/options"
	"github.com/wouldgo/speedtest-influxdb/src/speedtest"
)

func main() {
	log, logErr := configureLog()
	if logErr != nil {
		panic(logErr)
	}

	options, err := cliOptions.ParseOptions(log)
	if err != nil &&
		options == nil {

		log.Fatal(err)
	} else if err == nil &&
		options == nil {

		return
	}

	testSuite, err := speedtest.New(options)

	if err != nil {

		log.Fatal(err)
		return
	}

	theModel, err := model.New(options)
	if err != nil {

		log.Fatal(err)
		return
	}

	go traffic(log, testSuite.Data)

	fqdn, result, err := testSuite.Run()

	if err != nil {
		log.Fatal(err)
	}

	summary, err := speedtest.NewSummary(*fqdn, result)
	if err != nil {
		log.Fatal(err)
	}

	theModel.Write(summary)

	testSuite.Dispose()
	theModel.Dispose()
	log.Infof("Resources are disposed. Bye")
}

func configureLog() (*zap.SugaredLogger, error) {
	config := zap.Config{
		Level:            zap.NewDevelopmentConfig().Level,
		Encoding:         "console",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig:    zap.NewDevelopmentConfig().EncoderConfig,
	}

	logger, err := config.Build()

	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	return sugar, err
}

func traffic(log *zap.SugaredLogger, data chan *spec.Measurement) {
	for range data {

		log.Debug(".")
	}
}
