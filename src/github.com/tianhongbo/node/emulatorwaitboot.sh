#!/bin/bash
#exit when go main() send ok.Interrupt signal to this process
#trap 'echo "Exit 2(os.Interrupt signal detected... vncserver_id=$vncserver_pid, novnc_pid=$novnc_pid"; kill -9 $vncserver_pid; kill -9 $novnc_pid; exit 0' 2

# ADB name
adb_name=$1

echo "adb_name=$adb_name"

#waiting for device online
adb -s $adb_name wait-for-device

#waiting for device booting
A=$(adb -s $adb_name shell getprop sys.boot_completed | tr -d '\r')
while [ "$A" != "1" ]; do
        sleep 3
        A=$(adb -s $adb_name shell getprop sys.boot_completed | tr -d '\r')
done
