package dataconverter

import (
	"strconv"

	gpsd "github.com/koppacetic/go-gpsd"
)

type ConvertedDataSet map[string]ConvertedData

type ConvertedData struct {
	Timestamp string `json:"timestamp"`
	Value1    string `json:"value1"`
}

func Convert(r gpsd.TPVReport) ConvertedDataSet {
	time := r.Time.Format(("2006-01-02T15:04:05.000000Z"))
	dataSet := make(ConvertedDataSet)

	dataSet["mode"] = ConvertedData{
		time,
		string(r.Mode),
	}

	if r.Mode == 1 {
		return dataSet
	}

	dataSet["latitude"] = ConvertedData{
		time,
		strconv.FormatFloat(r.Lat, 'f', -1, 64),
	}
	dataSet["longitude"] = ConvertedData{
		time,
		strconv.FormatFloat(r.Lon, 'f', -1, 64),
	}

	if r.Mode == 2 {
		return dataSet
	}

	dataSet["altitude"] = ConvertedData{
		time,
		strconv.FormatFloat(r.Alt, 'f', -1, 64),
	}

	return dataSet
}
