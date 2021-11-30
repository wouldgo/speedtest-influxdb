package options

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

var (
	eigthHours, _ = time.ParseDuration("8h")
	oneMinute, _  = time.ParseDuration("1m")

	clientNameEnv     = os.Getenv("CLIENT_NAME")
	clientVersionEnv  = os.Getenv("CLIENT_VERSION")
	defaultTimeoutEnv = os.Getenv("DEFAULT_TIMEOUT")

	cacheSizeEnv     = os.Getenv("CACHE_SIZE")
	cacheEvictionEnv = os.Getenv("CACHE_EVICTION")

	influxDbUrlEnv    = os.Getenv("INFLUXDB_URL")
	influxDbTokenEnv  = os.Getenv("INFLUXDB_TOKEN")
	influxDbOrgEnv    = os.Getenv("INFLUXDB_ORG")
	influxDbBucketEnv = os.Getenv("INFLUXDB_BUCKET")

	clientName     = flag.String("client-name", "", "Speed test client name")
	clientVersion  = flag.String("client-version", "", "Speed test client version")
	defaultTimeout = flag.Duration("default-timeout", oneMinute, "Default timeout for speed test")

	cacheSize     = flag.Int("cache-size", 256, "LRU cache for meta gathering")
	cacheEviction = flag.Duration("cache-eviction", eigthHours, "LRU cache cache duration")

	influxDbUrl    = flag.String("influxdb-url", "", "Influxdb server url")
	influxDbToken  = flag.String("influxdb-token", "", "Influxdb api token")
	influxDbOrg    = flag.String("influxdb-organization", "", "Influxdb organization")
	influxDbBucket = flag.String("influxdb-bucket", "", "Influxb bucket")

	autocomplete = flag.Bool("zsh-autocomplete", false, "Print zsh autocomplete")
)

func ParseOptions(log *zap.SugaredLogger) (*Options, error) {
	flag.Parse()

	if *autocomplete {
		printCompletions("speedtest")
		return nil, nil
	}

	if *clientName == "" && clientNameEnv != "" {
		clientName = &clientNameEnv
	}

	if *clientVersion == "" && clientVersionEnv != "" {
		clientVersion = &clientVersionEnv
	}

	if *defaultTimeout == oneMinute && defaultTimeoutEnv != "" {
		clientVersion = &defaultTimeoutEnv
	}

	if *cacheSize == 256 && cacheSizeEnv != "" {
		buffSize, err := strconv.ParseInt(cacheSizeEnv, 10, 32)
		if err != nil {
			return nil, err
		}

		cacheSizeFromEnv := int(buffSize)
		cacheSize = &cacheSizeFromEnv
	}

	if *cacheEviction == eigthHours && cacheEvictionEnv != "" {
		cacheEvictionFromEnv, cacheEvictionFromEnvErr := time.ParseDuration(cacheEvictionEnv)
		if cacheEvictionFromEnvErr != nil {

			return nil, cacheEvictionFromEnvErr
		}

		cacheEviction = &cacheEvictionFromEnv
	}

	if *influxDbUrl == "" && influxDbUrlEnv != "" {
		influxDbUrl = &influxDbUrlEnv
	}

	if *influxDbToken == "" && influxDbTokenEnv != "" {
		influxDbToken = &influxDbTokenEnv
	}

	if *influxDbOrg == "" && influxDbOrgEnv != "" {
		influxDbOrg = &influxDbOrgEnv
	}

	if *influxDbBucket == "" && influxDbBucketEnv != "" {
		influxDbBucket = &influxDbBucketEnv
	}

	if *clientName == "" {
		return nil, errors.New("client name is required")
	}

	if *clientVersion == "" {
		return nil, errors.New("client version is required")
	}

	if *influxDbUrl == "" {

		return nil, errors.New("InfluxDb Url must be present")
	}

	if *influxDbToken == "" {

		return nil, errors.New("InfluxDb Token must be present")
	}

	if *influxDbOrg == "" {

		return nil, errors.New("InfluxDb Organitazion must be present")
	}

	if *influxDbBucket == "" {

		return nil, errors.New("InfluxDb Bucket must be present")
	}

	opts := Options{
		Autocomplete: autocomplete,
		Log:          log,
		SpeedTestConfiguration: &SpeedTestConfiguration{
			ClientName:     clientName,
			ClientVersion:  clientVersion,
			DefaultTimeout: defaultTimeout,
		},
		Cache: &CacheConfiguration{
			Size:     cacheSize,
			Eviction: cacheEviction,
		},
		InfluxDb: &InfluxDbConfigurations{
			Url:    influxDbUrl,
			Token:  influxDbToken,
			Org:    influxDbOrg,
			Bucket: influxDbBucket,
		},
	}

	return &opts, nil
}

type InfluxDbConfigurations struct {
	Url    *string
	Token  *string
	Org    *string
	Bucket *string
}

type GeoIpConfigurations struct {
	Folder *string
	Token  *string
}

type CacheConfiguration struct {
	Size     *int
	Eviction *time.Duration
}

type SpeedTestConfiguration struct {
	ClientName     *string
	ClientVersion  *string
	DefaultTimeout *time.Duration
}

type Options struct {
	Autocomplete           *bool
	SpeedTestConfiguration *SpeedTestConfiguration
	Log                    *zap.SugaredLogger
	Cache                  *CacheConfiguration
	InfluxDb               *InfluxDbConfigurations
}

func printCompletions(name string) {
	var cmpl []string
	flag.VisitAll(func(f *flag.Flag) {
		cmpl = append(
			cmpl, fmt.Sprintf("\t'-%s[%s]' \\\n", f.Name, f.Usage))
	})

	args := fmt.Sprintf("#compdef %s\n\n_arguments -s \\\n%s\n\n",
		name, strings.Join(cmpl, " "))
	fmt.Print(args)
}
