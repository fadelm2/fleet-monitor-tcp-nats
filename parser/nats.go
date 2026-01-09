package parser

import (
	"log"

	"github.com/nats-io/nats.go"
)

var nc *nats.Conn

func InitNATS(url string) {
	var err error
	nc, err = nats.Connect(url)
	if err != nil {
		log.Fatal("❌ NATS connect failed:", err)
	}
	log.Println("✅ Connected to NATS")
}

func PublishGPS(data []byte) error {
	return nc.Publish("gps.parsed", data)
}
