# save_my_macbook_battery
Impliment the recommended Macbook (2015) charging cycle with TP link plugs for macBooks that are plugged in permanently.
Note: The idea is to simplify the daily tasks of extend the battery life a bit, you will have to do research to ensure the proper setup for you battery.
USE AT YOUR OWN RISK

This is a small app that monitors the battery level and switches the a plug on/off based on how full the battery is charged.

## Requirements

- admin rights on your macBook
- TPLink Plug 100/110

## Install/update service

sh release.sh

## Remove the service

sh cleanup.sh

## Notes

- Logs are not rotated, (can grow substantially over time)
- No tests in the code (did not seem "warranted" at such a small scale)
- no GUI to monitor if it is working, but you will notice if when macOS tells you the battery is at 5% :-P
- For personal use, no guarantees, warrantees, etc
- Use at your own risk :)
