#!/bin/bash

#for num in 5554 5556 5558 5560 5562 5564 5566 5568 5570 5572
for num in 5554 5556 5558 5560 5562 5564
do
  echo $num
  /usr/local/android-sdk-linux/tools/emulator64-arm -avd android-api-10-$num -wipe-data -no-window -no-boot-anim -noskin -port $num &
done

