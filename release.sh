#!/bin/bash
set -e
# Any subsequent(*) commands which fail will cause the shell script to exit immediately
go build
sudo cp save_my_macbook_battery /usr/local/bin/
sudo cp .env_save_my_macbook_battery /usr/local/bin/
sudo chmod 777 /usr/local/bin/save_my_macbook_battery
sudo cp macOS/com.edwinludik.save_my_macbook_battery.plist  ~/Library/LaunchAgents/
sudo launchctl load ~/Library/LaunchAgents/com.edwinludik.save_my_macbook_battery.plist
