package models

import (
	"time"
)

type Metrics struct {
	Id             int       `json:"id" db:"id" redis:"id"`
	MediaId        string    `json:"media_id" db:"media_id" redis:"media_id"`
	Ip             string    `json:"ip" db:"ip" redis:"ip"`
	Browser        string    `json:"browser" db:"browser" redis:"browser"`
	BrowserVersion string    `json:"browser_version" db:"browser_version" redis:"browser_version"`
	Os             string    `json:"os" db:"os" redis:"os"`
	OsVersion      string    `json:"os_version" db:"os_version" redis:"os_version"`
	Device         string    `json:"device" db:"device" redis:"device"`
	DeviceBrand    string    `json:"device_brand" db:"device_brand" redis:"device_brand"`
	DeviceModel    string    `json:"device_model" db:"device_model" redis:"device_model"`
	Country        string    `json:"country" db:"country" redis:"country"`
	Referrer       string    `json:"referrer" db:"referrer" redis:"referrer"`
	CapturedAt     time.Time `json:"captured_at" db:"captured_at" redis:"captured_at"`
}
