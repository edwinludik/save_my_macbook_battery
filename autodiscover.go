package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/distatus/battery"
	"github.com/joho/godotenv"

	"github.com/jaedle/golang-tplink-hs100/pkg/configuration"
	"github.com/jaedle/golang-tplink-hs100/pkg/hs100"
)

func main() {
	// load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// get battery info
	batteries, err := battery.GetAll()
	if err != nil {
		fmt.Println("Could not get battery info!")
		return
	}
	for i, battery := range batteries {
		fmt.Printf("Bat%d: ", i)
		fmt.Printf("state: %s, ", battery.State.String())
		fmt.Printf("current capacity: %f mWh, ", battery.Current)
		fmt.Printf("last full capacity: %f mWh, ", battery.Full)
		fmt.Printf("design capacity: %f mWh, ", battery.Design)
		fmt.Printf("charge rate: %f mW, ", battery.ChargeRate)
		fmt.Printf("voltage: %f V, ", battery.Voltage)
		fmt.Printf("design voltage: %f V\n", battery.DesignVoltage)
	}

	// get the available plugs
	networkMask := os.Getenv("NETWORK_MASK")
	devices, err := hs100.Discover(networkMask,
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
