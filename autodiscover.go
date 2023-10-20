package main

import (
	"log"
	"time"

	"github.com/jaedle/golang-tplink-hs100/pkg/configuration"
	"github.com/jaedle/golang-tplink-hs100/pkg/hs100"
)

func main() {
	devices, err := hs100.Discover("192.168.1.0/24",
		configuration.Default().WithTimeout(time.Second),
	)

	if err != nil {
		panic(err)
	}

	log.Printf("Found devices: %d", len(devices))
	for _, d := range devices {
		info, _ := d.GetInfo()
		log.Printf("Found device (name, id): %s, %s", info.Name, info.DeviceId)
	}
}
