# controller
This is a controller written by Go to create, delete a Android emulator based Android SDK on Mac OS, Linux Desktop, or Linux Server.

# setup
AWS EC2: free tier (t2.micro)

Linux Ubuntu: Linux ip-172-31-16-251 4.4.0-53-generic #74-Ubuntu SMP Fri Dec 2 15:59:10 UTC 2016 x86_64 x86_64 x86_64 GNU/Linux

# Must read before starting
This doc instructs how to install the controller on a Ubuntu Server.
The big chanllenge in this installation is the Android SDK part since it is usually installed in a machine with graphic device via GUI. But in this case, all the installation need to be done via command line becasue a server does not provide GUI interface.

Follow this post to install Android SDK on headleass Ubuntu server
http://sblackwell.com/blog/2014/06/installing-the-android-sdk-on-a-headless-server/

The sdkmanager is a command line tool that allows you to view, install, update, and uninstall packages for the Android SDK. If you're using Android Studio, then you do not need to use this tool and you can instead manage your SDK packages from the IDE.

https://developer.android.com/studio/command-line/sdkmanager.html

# How can I install it?
## 1. Install Android SDK
- install Java SDK (java-8-openjdk-amd64)
- install Android SDK Tool (25.2.3)
- install Android 2.3.3(API 10)

## 2. Install GO
In the EC2, AWS has already installed go, so nothing to do here.
```
$ go version

go version go1.6.2 linux/amd64
```

## 3. Install adb
even though it has been installed together with Android SDK/Studio
`$ sudo apt-get install android-tools-adb`

## 4. Install noVNC
$ sudo git clone git://github.com/kanaka/noVNC

## 5: Git Clone NODE source code
`$ git clone https://github.com/tianhongbo/controller.git`

# How can I configure it?
## 1. Modify source code
- install.sh
- deviceinstall.sh
- repo.go

## 2. Build go executable file
`$ cd /home/ubuntu/controller/src/github.com/tianhongbo/node`

get dependency packages

`$ go get`

build bin file

`$ go install`

check the bin file
```
$ ls -l ~/controller/bin
total 8640
-rwxrwxr-x 1 ubuntu ubuntu 8846224 Dec 14 19:15 node
```

## 2. Set environment variables
$ sudo vi /etc/environment

- add emulator tools to PATH
- add dynamic lib path for Android SDK
- add GOPATH

Here is one sample
```
export ANDROID_SDK_HOME="/usr/local/android-sdk-linux"
PATH="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games:/usr/local/games:$ANDROID_SDK_HOME/tools:$ANDROID_SDK_HOME/platform-tools"
JAVA_HOME="/usr/lib/jvm/java-8-openjdk-amd64"
export LD_LIBRARY_PATH="$ANDROID_SDK_HOME/tools/lib64:$ANDROID_SDK_HOME/tools/lib64/qt/lib:$LD_LIBRARY_PATH"
export GOPATH="/home/ubuntu/controller"
```

## 3. configure port forwarding for SSH functions
- iptables for Ubuntu
- ip for mac

For example(for ubuntu):
### This is for eth0
`sysctl -w net.ipv4.conf.eth0.route_localnet=1`
`iptables -t nat -I PREROUTING -p tcp -i eth0 --dport 5921 -j DNAT --to-destination 127.0.0.1:5921`

### This is for wlan
`sysctl -w net.ipv4.conf.wlan0.route_localnet=1`
`iptables -t nat -I PREROUTING -p tcp -i wlan0 --dport 5921 -j DNAT --to-destination 127.0.0.1:5921`

# How can I debug?
## 1. check adb forward
`$ adb forward --list`


## 2. check iptables NAT
`$ sudo iptables -t nat -L -n -v`

## 3. SSH
on the mobile phone
user: "test" / password: "test"
ssh test@192.168.1.24 -p 5921

## 4. VNC
It seems the behaviors in MAC and Ubuntu are different. For example, VNC launch.sh only create
listening port on the external interface in MAC. However, it will create listening port on all
the interfaces including localhost.

2. /etc/rc.local
start node automatically when instance booting
su ubuntu - -c /home/ubuntu/node/bin/node
