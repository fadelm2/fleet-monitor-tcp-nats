package parser

type GPSLocation struct {
	IMEI  string  `json:"imei"`
	Lat   float64 `json:"lat"`
	Lon   float64 `json:"lon"`
	Speed float64 `json:"speed"`
}
