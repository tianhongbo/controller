#!/bin/bash

# adb name
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

#disconnect Internet connection
#adb -s $adb_name shell svc data disable
#adb -s $adb_name shell svc wifi disable

adb -s $adb_name shell 'su -c "setprop net.dns1 0.0.0.0"'
adb -s $adb_name shell 'su -c "setprop net.dns2 0.0.0.0"'
