#!/bin/bash

# adb_name
adb_name=$1
echo "adb_name=$adb_name"

#waiting for device online
adb -s $adb_name wait-for-device

#waiting for device booting
A=$(adb -s $adb_name shell getprop sys.boot_completed | tr -d '\r')
while [ "$A" != "1" ]; do
        sleep 1
        A=$(adb -s $adb_name shell getprop sys.boot_completed | tr -d '\r')
done

#connect Internet connection
#adb -s $adb_name shell svc data enable
#adb -s $adb_name shell svc wifi enable
#adb -s $adb_name shell setprop net.dns1 192.168.1.1
adb -s $adb_name shell 'su -c "setprop net.dns1 8.8.8.8"'
adb -s $adb_name shell 'su -c "setprop net.dns2 8.8.4.4"'
#Open browser and visit "www.sjsu.edu" website
adb -s $adb_name shell am start -a android.intent.action.VIEW -d http://www.sjsu.edu