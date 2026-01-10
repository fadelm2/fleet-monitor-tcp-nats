package parser

import (
	"fleet-monitor/logger"
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
	logger.Log.Info("NATS PUBLISH EXECUTED")

	return nc.Publish("gps.parsed", data)
}
