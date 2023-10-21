package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/distatus/battery"
	"github.com/joho/godotenv"

	"github.com/jaedle/golang-tplink-hs100/pkg/configuration"
	"github.com/jaedle/golang-tplink-hs100/pkg/hs100"
)

func main() {
	// global variables
	chargeLevel := 0.0

	// load .env
	err := godotenv.Load(".env_save_my_macbook_battery")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// get battery info
	batteries, err := battery.GetAll()
	if err != nil {
		fmt.Println("Could not get battery info!")
		return
	}
	for _, battery := range batteries {
		// only get the value of the first battery
		chargeLevel = battery.Current / battery.Full * 100
		break
	}

	fmt.Printf("Current Battery Percentage %.2f%%: ", chargeLevel)

	// load the charge values
	chargeMinimum, err := strconv.ParseFloat(os.Getenv("CHARGE_MINIMUM"), 32)
	if err != nil {
		log.Fatal("CHARGE_MINIMUM is not a valid number float/integer")
	}
	chargeMaximum, err := strconv.ParseFloat(os.Getenv("CHARGE_MAXIMUM"), 32)
	if err != nil {
		log.Fatal("CHARGE_MINIMUM is not a valid number float/integer")
	}
	// exit early if the charge % warrants no action
	if chargeLevel <= chargeMaximum && chargeLevel >= chargeMinimum {
		fmt.Printf("No action needed (Min/Max) = (%.2f%%/%.2f%%) ", chargeMinimum, chargeMaximum)
		os.Exit(0)
	}

	// get the available plugs and act as needed
	networkMask := os.Getenv("NETWORK_MASK")
	devices, err := hs100.Discover(networkMask,
		configuration.Default().WithTimeout(time.Second),
	)
	if err != nil {
		panic(err)
	}

	target_plug := os.Getenv("TARGET_PLUG_NAME_OR_ID")
	// log.Printf("Found devices: %d", len(devices))
	plug_is_found := false
	for _, d := range devices {
		info, _ := d.GetInfo()
		// log.Printf("Found device (name, id): %s, %s", info.Name, info.DeviceId)
		if info.Name == target_plug || info.DeviceId == target_plug {
			plug_is_found = true
			is_on, is_on_err := d.IsOn()
			if is_on_err != nil {
				panic(err)
			}
			// check the plugs current state
			if chargeLevel > chargeMaximum && is_on {
				d.TurnOff()
			}
			if chargeLevel < chargeMinimum && !is_on {
				d.TurnOn()
			}
		}
	}
	// we could not find the plug, show all the available ones
	if !plug_is_found {
		for _, d := range devices {
			info, _ := d.GetInfo()
			log.Printf("Found device (name, id): %s, %s", info.Name, info.DeviceId)
		}
	}
}
