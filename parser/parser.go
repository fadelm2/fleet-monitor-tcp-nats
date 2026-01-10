package parser

import (
	"encoding/json"
	"strconv"
	"strings"

	"fleet-monitor/logger"
)

func ParseAndPublish(data []byte) {
	logger.Log.Infof("PARSER HIT raw=%q", string(data))
	logger.Log.Info("PARSER WILL PUBLISH gps.parsed")

	if len(data) < 20 {
		return
	}

	payload := string(data)

	if !strings.HasPrefix(payload, "(") || !strings.HasSuffix(payload, ")") {
		return
	}

	proto := payload[1:3]
	imei := payload[5:14]

	if proto != "02" {
		return
	}

	parseLocation(payload, imei)
}

func parseLocation(p, imei string) {

	if p[23:24] != "A" {
		return
	}

	latRaw := p[24:33]
	latDir := p[33:34]
	lonRaw := p[34:44]
	lonDir := p[44:45]

	speedRaw := p[45:48]
	speed, _ := strconv.ParseFloat(speedRaw, 64)

	lat := convertCoord(latRaw, latDir)
	lon := convertCoord(lonRaw, lonDir)

	gps := GPSLocation{
		IMEI:  imei,
		Lat:   lat,
		Lon:   lon,
		Speed: speed,
	}

	payload, err := json.Marshal(gps)
	if err != nil {
		return
	}

	err = PublishGPS(payload)
	if err != nil {
		logger.Log.Error("NATS publish error:", err)
	}
}

func convertCoord(raw, dir string) float64 {
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0
	}

	deg := float64(int(v / 100))
	min := v - deg*100
	dec := deg + min/60

	if dir == "S" || dir == "W" {
		dec = -dec
	}
	return dec
}
