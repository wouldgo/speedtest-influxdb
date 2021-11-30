package speedtest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type SpeedTestServer struct {
	// The right side is the name of the JSON variable
	Country     string  `json:"country"`
	City        string  `json:"city"`
	Lat         float64 `json:"latitude"`
	Lon         float64 `json:"longitude"`
	Roundrobin  bool    `json:"roundrobin"`
	Site        string  `json:"site"`
	UplinkSpeed string  `json:"uplink_speed"`
}

func ListServer() ([]SpeedTestServer, error) {
	var servers []SpeedTestServer

	response, err := http.Get("https://siteinfo.mlab-oti.measurementlab.net/v1/sites/locations.json")
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {

		return nil, err
	}

	err = json.Unmarshal(body, &servers)
	if err != nil {

		return nil, err
	}

	return servers, nil
}
