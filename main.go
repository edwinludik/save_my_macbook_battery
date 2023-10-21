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

	debuggingEnabled, err := strconv.ParseBool(os.Getenv("DEBUGGING_ENABLED"))
	if err != nil {
		log.Fatal("DEBUGGING_ENABLED is not a valid boolean: True/False")
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
	if debuggingEnabled {
		fmt.Printf("Current Battery Percentage %.2f%%\n", chargeLevel)
	}

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
		if debuggingEnabled {
			fmt.Printf("No action needed (Min/Max) = (%.2f%%/%.2f%%)", chargeMinimum, chargeMaximum)
		}
		// clean exit
		os.Exit(0)
	}

	// get the available plugs and switch on/off
	if debuggingEnabled {
		fmt.Printf("Searching the Target plug...\n")
	}
	networkMask := os.Getenv("NETWORK_MASK")
	devices, err := hs100.Discover(networkMask,
		configuration.Default().WithTimeout(time.Second),
	)
	if err != nil {
		panic(err)
	}
	targetPlug := os.Getenv("TARGET_PLUG_NAME_OR_ID")
	// log.Printf("Found devices: %d", len(devices))
	plugIsFound := false
	for _, d := range devices {
		info, _ := d.GetInfo()
		// log.Printf("Found device (name, id): %s, %s", info.Name, info.DeviceId)
		if info.Name == targetPlug || info.DeviceId == targetPlug {
			plugIsFound = true
			isPlugOn, is_on_err := d.IsOn()
			if is_on_err != nil {
				panic(is_on_err)
			}
			// check the plugs current state
			if (chargeLevel < chargeMinimum) && !isPlugOn {
				if debuggingEnabled {
					fmt.Printf("Action: Turn On\n")
				}
				d.TurnOn()
			} else if (chargeLevel > chargeMaximum) && isPlugOn {
				if debuggingEnabled {
					fmt.Printf("Action: Turn Off\n")
				}
				d.TurnOff()
			} else {
				if debuggingEnabled {
					fmt.Printf("Action: None\n")
					fmt.Printf("Status (chargeLevel/Min/Max/isPlugOn) = (%.2f%% / %.2f%% / %.2f%% / %t)", chargeLevel, chargeMinimum, chargeMaximum, isPlugOn)
				}
			}
		}
	}
	// we could not find the plug, show all the available ones
	if !plugIsFound {
		for _, d := range devices {
			info, _ := d.GetInfo()
			log.Printf("Found device (name, id): %s, %s", info.Name, info.DeviceId)
		}
		log.Fatal("Could not find plug: ", targetPlug)
	}
}
