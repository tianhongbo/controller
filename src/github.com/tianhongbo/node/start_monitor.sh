#!/usr/bin/env bash
# 1.start ssh
 ssh -i /Users/hongbotian/Downloads/cloud.key -NL 5554:localhost:5554 -L 5555:localhost:5555 ubuntu@8.21.28.162 &
# 2.adb devices
echo /Users/hongbotian/Downloads/android-sdk-macosx/platform-tools/adb devices
/Users/hongbotian/Downloads/android-sdk-macosx/platform-tools/adb devices
# 3.adb kill-server
/Users/hongbotian/Downloads/android-sdk-macosx/platform-tools/adb kill-server
# 4. ./platoform-tools/adb start-server
/Users/hongbotian/Downloads/android-sdk-macosx/platform-tools/adb start-server
# 5. ./platoform-tools/adb devices
/Users/hongbotian/Downloads/android-sdk-macosx/platform-tools/adb devices
# 6. ./tools/monitor
/Users/hongbotian/Downloads/android-sdk-macosx/tools/monitor &
~
