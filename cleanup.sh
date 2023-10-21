#!/bin/bash
set -e
# Any subsequent(*) commands which fail will cause the shell script to exit immediately
sudo launchctl disable system/com.edwinludik.save_my_macbook_battery
sudo launchctl remove system/com.edwinludik.save_my_macbook_battery.plist
sudo rm /Library/LaunchAgents/com.edwinludik.save_my_macbook_battery.plist
sudo rm /usr/local/bin/save_my_macbook_battery
sudo rm /usr/local/bin/.env_save_my_macbook_battery
echo Done
# launchctl list | grep save_my_macbook_battery
# /Library/LaunchAgents
# rm ~/Library/LaunchAgents/com.edwinludik.save_my_macbook_battery.plist