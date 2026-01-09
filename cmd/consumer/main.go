package main

import (
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
)

type GPSLocation struct {
	IMEI  string  `json:"imei"`
	Lat   float64 `json:"lat"`
	Lon   float64 `json:"lon"`
	Speed float64 `json:"speed"`
}

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)
	defer nc.Close()

	log.Println("📥 GPS Consumer started")

	nc.Subscribe("gps.parsed", func(msg *nats.Msg) {
		var gps GPSLocation
		json.Unmarshal(msg.Data, &gps)

		log.Printf(
			"💾 SAVE DB IMEI=%s LAT=%f LON=%f SPEED=%f\n",
			gps.IMEI, gps.Lat, gps.Lon, gps.Speed,
		)

		// saveToDB(gps)
	})

	select {}
}
