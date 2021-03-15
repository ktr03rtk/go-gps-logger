package dataconverter

import (
	"strconv"

	gpsd "github.com/koppacetic/go-gpsd"
)

type ConvertedDataSet map[string]convertedData

type convertedData struct {
	timestamp string
	value1    string
}

func Convert(r gpsd.TPVReport) ConvertedDataSet {
	time := r.Time.Format(("2006-01-02T15:04:05.000000Z"))
	dataSet := make(ConvertedDataSet)

	dataSet["mode"] = convertedData{
		time,
		string(r.Mode),
	}

	if r.Mode == 1 {
		return dataSet
	}

	dataSet["latitude"] = convertedData{
		time,
		strconv.FormatFloat(r.Lat, 'f', -1, 64),
	}
	dataSet["longitude"] = convertedData{
		time,
		strconv.FormatFloat(r.Lon, 'f', -1, 64),
	}

	if r.Mode == 2 {
		return dataSet
	}

	dataSet["altitude"] = convertedData{
		time,
		strconv.FormatFloat(r.Alt, 'f', -1, 64),
	}

	return dataSet
}
