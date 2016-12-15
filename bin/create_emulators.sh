#!/bin/bash

# set environment for MTaaS
export PATH="$PATH:/usr/local/android-sdk-linux/tools:/usr/local/android-sdk-linux/platform-tools"
export JAVA_HOME="/usr/lib/jvm/java-8-openjdk-amd64"
export LD_LIBRARY_PATH="$LD_LIBRARY_PATH:/usr/local/android-sdk-linux/tools/lib64:/usr/local/android-sdk-linux/tools/lib64/qt/lib"
export GOPATH="/home/ubuntu/controller"

#for num in 5554 5556 5558 5560 5562 5564 5566 5568 5570 5572
env > /tmp/scott2.out

for num in 5554 5556 5558
do
  echo $num
  emulator64-arm -avd android-api-10-$num -wipe-data -no-window -no-boot-anim -noskin -port $num &
done

