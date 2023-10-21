#!/bin/bash
set -e
# Any subsequent(*) commands which fail will cause the shell script to exit immediately

# build and deploy app
go build
cp save_my_macbook_battery /usr/local/bin/
sudo chmod 775 /usr/local/bin/save_my_macbook_battery
sudo chown $USER /usr/local/bin/save_my_macbook_battery
cp .env_save_my_macbook_battery /usr/local/bin/
sudo chmod 775 /usr/local/bin/.env_save_my_macbook_battery
sudo chown $USER /usr/local/bin/.env_save_my_macbook_battery
# create logging/output directors
sudo mkdir -p /var/log/save_my_macbook_battery
# macos service
sudo launchctl stop system/com.edwinludik.save_my_macbook_battery.plist
sudo launchctl disable system/com.edwinludik.save_my_macbook_battery
sudo launchctl remove system/com.edwinludik.save_my_macbook_battery.plist
sudo cp macOS/com.edwinludik.save_my_macbook_battery.plist /Library/LaunchAgents/
sudo chmod 644 /Library/LaunchAgents/com.edwinludik.save_my_macbook_battery.plist
sudo chown $USER /Library/LaunchAgents/com.edwinludik.save_my_macbook_battery.plist
sudo launchctl load com.edwinludik.save_my_macbook_battery.plist
sudo launchctl enable system/com.edwinludik.save_my_macbook_battery.plist
sudo launchctl start system/com.edwinludik.save_my_macbook_battery.plist
echo Done

# random commands used during testing/debugging
# launchctl list | grep save_my_macbook_battery
# /Library/LaunchAgents
# rm ~/Library/LaunchAgents/com.edwinludik.save_my_macbook_battery.plist
# plutil /Library/LaunchAgents/com.edwinludik.save_my_macbook_battery.plist
# sudo launchctl status system/com.edwinludik.save_my_macbook_battery.plist