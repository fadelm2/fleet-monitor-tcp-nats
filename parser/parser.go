package parser

import (
	"encoding/json"
	"strconv"
	"strings"

	"fleet-monitor/logger"
)

func ParseAndPublish(data []byte) {
	payload := strings.Trim(string(data), "()")
	logger.Log.Infof("RAW=%s", payload)

	// split by known tokens
	if !strings.Contains(payload, "A") {
		logger.Log.Warn("NO GPS STATUS")
		return
	}

	// IMEI: setelah "02" + 1 digit length
	if len(payload) < 15 {
		return
	}
	imei := payload[3:12]

	// cari posisi status A
	idx := strings.Index(payload, "A")
	if idx == -1 || len(payload) < idx+20 {
		return
	}

	coord := payload[idx+1:]

	latRaw := coord[:9]
	latDir := coord[9:10]
	lonRaw := coord[10:20]
	lonDir := coord[20:21]

	speedRaw := coord[21:24]

	speed, _ := strconv.ParseFloat(speedRaw, 64)

	lat := convertCoord(latRaw, latDir)
	lon := convertCoord(lonRaw, lonDir)

	gps := GPSLocation{
		IMEI:  imei,
		Lat:   lat,
		Lon:   lon,
		Speed: speed,
	}

	b, _ := json.Marshal(gps)
	PublishGPS(b)

	logger.Log.Infof("✅ PUBLISHED IMEI=%s LAT=%f LON=%f", imei, lat, lon)
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
