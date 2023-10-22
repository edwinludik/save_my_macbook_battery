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

	fmt.Print("*** ")
	fmt.Printf("Checking at %s | ", time.Now().String())

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
		fmt.Printf("No action needed (Current/Min/Max) = (%.2f%%/%.2f%%/%.2f%%) \n", chargeLevel, chargeMinimum, chargeMaximum)
		// clean exit
		os.Exit(0)
	}

	// get the available plugs and switch on/off
	targetPlug := os.Getenv("TARGET_PLUG_NAME_OR_ID")
	if debuggingEnabled {
		fmt.Printf("Searching the Target plug (%s) ...\n", targetPlug)
	}
	networkMask := os.Getenv("NETWORK_MASK")
	devices, err := hs100.Discover(networkMask,
		configuration.Default().WithTimeout(time.Second),
	)
	if err != nil {
		panic(err)
	}
	// log.Printf("Found devices: %d", len(devices))
	plugIsFound := false
	for _, d := range devices {
		info, _ := d.GetInfo()
		// log.Printf("Found device (name, id): %s, %s", info.Name, info.DeviceId)
		if (info != nil) &&
			(info.Name == targetPlug || info.DeviceId == targetPlug) {
			plugIsFound = true
			isPlugOn, _ := d.IsOn()
			// check the plugs current state
			if (chargeLevel < chargeMinimum) && !isPlugOn {
				fmt.Printf("Action: Turn On\n")
				err = d.TurnOn()
				if err != nil {
					panic(err)
				}
			} else if (chargeLevel > chargeMaximum) && isPlugOn {
				fmt.Printf("Action: Turn Off\n")
				err = d.TurnOff()
				if err != nil {
					panic(err)
				}
			} else {
				fmt.Printf("Action: None\n")
				if debuggingEnabled {
					fmt.Printf("Status (chargeLevel/Min/Max/isPlugOn) = (%.2f%% / %.2f%% / %.2f%% / %t)", chargeLevel, chargeMinimum, chargeMaximum, isPlugOn)
				}
			}
		} else {
			fmt.Printf(" info is null | (chargeLevel/Min/Max) = (%.2f%% / %.2f%% / %.2f%%) \n", chargeLevel, chargeMinimum, chargeMaximum)
		}
	}
	// we could not find the plug, show all the available ones
	if !plugIsFound {
		fmt.Printf("*** ERR *** Could not find plug: %s", targetPlug)
		log.Fatal("Could not find plug: ", targetPlug)
		for _, d := range devices {
			info, _ := d.GetInfo()
			fmt.Printf("But Found device (name, id): %s, %s", info.Name, info.DeviceId)
			log.Printf("But Found device (name, id): %s, %s", info.Name, info.DeviceId)
		}
	}
}
