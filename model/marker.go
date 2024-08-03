package model

type Markers struct {
	M map[string]Marker `xml:"marker"`
}

type Marker struct {
	Sn       string `xml:"id,attr"`
	Car      string `xml:"car,attr"`
	Lat      string `xml:"lat,attr"`
	Lon      string `xml:"lon,attr"`
	Speed    string `xml:"speed,attr"`
	Alt      string `xml:"alt,attr"`
	Brg      string `xml:"brg,attr"`
	Hdg	 string `xml:"hdg,attr"`
	Mnw      string `xml:"mnw,attr"`
	Time     string `xml:"time,attr"`
	UnixTime int64  `xml:"unixtime,attr"`
}

