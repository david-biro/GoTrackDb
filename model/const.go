package model

import (
	"GoTrackDb/config"
	"time"
)

const (
	Header       = `<?xml version="1.0" encoding="UTF-8"?>` + "\n"
	Markersbegin = `<markers>` + "\n"
	Markersend   = `</markers>`
)

var (
	Dbhost           = ""
	Dbport           = ""
	Dbname           = ""
	Dbuser           = ""
	Dbpass           = ""
	Precision        = "s"
	MeasurementPoint = "locationdata"
	Timeout          = ConvertTimeout("500ms")
)

func SetDbParams() {
	Dbhost = config.Cfg.Dbhost
	Dbport = config.Cfg.Dbport
	Dbname = config.Cfg.Dbname
	Dbuser = config.Cfg.Dbuser
	Dbpass = config.Cfg.Dbpass

}

// GetInfluxFieldNames ...
func GetInfluxFieldNames() []string {
	return []string{"Car", "Lat", "Lon", "Speed", "Alt", "Mnw", "Time"}
}

func ConvertTimeout(duration string) time.Duration {
	result, err := time.ParseDuration(duration)
	if err != nil {
		panic(err)
	}
	return result
}
