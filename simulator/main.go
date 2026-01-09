package main

import (
	"fmt"
	"net"
	"time"
)

type Device struct {
	IMEI string
	Lat  string
	Lon  string
}

func main() {
	devices := []Device{
		{"044400001", "0610.2215S", "10643.9911E"},
		{"044400002", "0610.3215S", "10644.0911E"},
		{"044400003", "0610.4215S", "10644.1911E"},
		{"044400004", "0610.5215S", "10644.2911E"},
		{"044400005", "0610.6215S", "10644.3911E"},
	}

	for _, d := range devices {
		go simulateDevice(d)
	}

	select {} // block forever
}

func simulateDevice(d Device) {
	conn, err := net.Dial("tcp", "127.0.0.1:9000")
	if err != nil {
		fmt.Println("❌ connect error:", err)
		return
	}
	defer conn.Close()

	fmt.Println("🚗 device connected:", d.IMEI)

	for {
		packet := buildGT02APacket(d)
		conn.Write([]byte(packet))
		fmt.Printf("📤 %s -> %s\n", d.IMEI, packet)
		time.Sleep(5 * time.Second)
	}
}

func buildGT02APacket(d Device) string {
	/*
	   Format disederhanakan (sesuai parser kamu):
	   (02 8 IMEI BR 00 DATE A LAT LON SPEED ...)
	*/

	speed := "045" // 45 km/h
	packet := fmt.Sprintf(
		"(028%sBR00260105A%s%s%s000000000L00000000)",
		d.IMEI,
		d.Lat,
		d.Lon,
		speed,
	)
	return packet
}
