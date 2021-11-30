package speedtest

import (
	"errors"

	"github.com/m-lab/ndt7-client-go"
	"github.com/m-lab/ndt7-client-go/spec"
)

// ValueUnitPair represents a {"Value": ..., "Unit": ...} pair.
type ValueUnitPair struct {
	Value float64
	Unit  string
}

// Summary is a struct containing the values displayed to the user at
// the end of an speedtest test.
type Summary struct {
	// ServerFQDN is the FQDN of the server used for this test.
	ServerFQDN string

	// ServerIP is the (v4 or v6) IP address of the server.
	ServerIP string

	// ClientIP is the (v4 or v6) IP address of the Client.
	ClientIP string

	// DownloadUUID is the UUID of the download test.
	DownloadUUID string

	// Download is the download speed, in Mbit/s. This is measured at the
	// receiver.
	Download ValueUnitPair

	// Upload is the upload speed, in Mbit/s. This is measured at the sender.
	Upload ValueUnitPair

	// DownloadRetrans is the retransmission rate. This is based on the TCPInfo
	// values provided by the server during a download test.
	DownloadRetrans ValueUnitPair

	// RTT is the round-trip time of the latest measurement, in milliseconds.
	// This is provided by the server during a download test.
	MinRTT ValueUnitPair
}

// NewSummary returns a new Summary struct for a given FQDN.
func NewSummary(FQDN string, result map[spec.TestKind]*ndt7.LatestMeasurements) (*Summary, error) {
	summary := &Summary{
		ServerFQDN: FQDN,
	}
	download, downloadOk := result[spec.TestDownload]
	upload, uploadOk := result[spec.TestUpload]

	if downloadOk {
		if download.Client.AppInfo != nil &&
			download.Client.AppInfo.ElapsedTime > 0 {
			elapsed := float64(download.Client.AppInfo.ElapsedTime) / 1e06
			summary.Download = ValueUnitPair{
				Value: (8.0 * float64(download.Client.AppInfo.NumBytes)) /
					elapsed / (1000.0 * 1000.0),
				Unit: "Mbit/s",
			}
		}

		if download.Server.TCPInfo != nil {
			if download.Server.TCPInfo.BytesSent > 0 {
				summary.DownloadRetrans = ValueUnitPair{
					Value: float64(download.Server.TCPInfo.BytesRetrans) / float64(download.Server.TCPInfo.BytesSent) * 100,
					Unit:  "%",
				}
			}
			summary.MinRTT = ValueUnitPair{
				Value: float64(download.Server.TCPInfo.MinRTT) / 1000,
				Unit:  "ms",
			}
		}
	}

	// Upload comes from the client-side Measurement during the upload test.
	if uploadOk &&
		upload.Client.AppInfo != nil &&
		upload.Client.AppInfo.ElapsedTime > 0 {
		elapsed := float64(upload.Client.AppInfo.ElapsedTime) / 1e06
		summary.Upload = ValueUnitPair{
			Value: (8.0 * float64(upload.Client.AppInfo.NumBytes)) /
				elapsed / (1000.0 * 1000.0),
			Unit: "Mbit/s",
		}
	}

	if !downloadOk || !uploadOk {

		return nil, errors.New("download or upload failed")
	}

	return summary, nil
}
