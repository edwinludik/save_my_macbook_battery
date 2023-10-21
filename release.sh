#!/bin/bash
set -e
# Any subsequent(*) commands which fail will cause the shell script to exit immediately
go build
cp save_my_macbook_battery /usr/local/bin/
sudo chmod 775 /usr/local/bin/save_my_macbook_battery
sudo chown $USER /usr/local/bin/save_my_macbook_battery
cp .env_save_my_macbook_battery /usr/local/bin/
sudo chmod 775 /usr/local/bin/.env_save_my_macbook_battery
sudo chown $USER /usr/local/bin/.env_save_my_macbook_battery

sudo cp macOS/com.edwinludik.save_my_macbook_battery.plist /Library/LaunchAgents/
sudo chmod 644 /Library/LaunchAgents/com.edwinludik.save_my_macbook_battery.plist
sudo chown $USER /Library/LaunchAgents/com.edwinludik.save_my_macbook_battery.plist
sudo launchctl load /Library/LaunchAgents/com.edwinludik.save_my_macbook_battery.plist
echo Done
# launchctl list | grep save_my_macbook_battery
# /Library/LaunchAgents
# rm ~/Library/LaunchAgents/com.edwinludik.save_my_macbook_battery.plist
# plutil /Library/LaunchAgents/com.edwinludik.save_my_macbook_battery.plist
