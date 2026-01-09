package main

import (
	"net"

	"fleet-monitor/logger"
	"fleet-monitor/parser"

	"github.com/nats-io/nats.go"
)

func main() {
	logger.Init()
	parser.InitNATS(nats.DefaultURL)

	addr := ":9000"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Log.Fatal(err)
	}

	logger.Log.Infof("🚀 TCP Server listening on %s", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 2048)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			return
		}

		// COPY buffer (anti corruption)
		data := make([]byte, n)
		copy(data, buf[:n])

		parser.ParseAndPublish(data)
	}
}
