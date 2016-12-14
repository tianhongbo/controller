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
### Install Java SDK (java-8-openjdk-amd64)

`$ sudo apt-update`

/* check whether java is installed or not */

`$ java -version`

/* Install java sdk */

`$ sudo apt-get install default-jdk`

/* check java path*/
```
$ sudo update-alternatives --config java
There is only one alternative in link group java (providing /usr/bin/java): /usr/lib/jvm/java-8-openjdk-amd64/jre/bin/java
Nothing to configure
```
/* set JAVA_HOME */
```
$ sudo vi /etc/environment
JAVA_HOME="/usr/lib/jvm/java-8-openjdk-amd64"
```

### Install Android SDK Tool (25.2.3)

/* install Android SDK tool */
```
$ mkdir /usr/local/android-sdk-linux
$ cd /usr/local/android-sdk-linux
$ wget https://dl.google.com/android/repository/tools_r25.2.3-linux.zip
```
/* install unzip */

`$ sudo apt install unzip`

/* unzip Android SDK tool to android-sdk-linux/ */

```
$ unzip tools_r25.2.3-linux.zip android-sdk-linux/
$ ls /usr/local/android-sdk-linux/
tools
```

/* set environment */
```
$ sudo vi /etc/environment
"/usr/local/games:/usr/local/android-sdk-linux/tools:/usr/local/android-sdk-linux/platform-tools"
```

### Install Android 2.3.3(API 10)

/*
 * install SDK packages via /tools/bin/sdkmanager
 * The sdkmanager is a command line tool that allows you to view,
 * install, update, and uninstall packages for the Android SDK.
 * If you're using Android Studio, then you do not need to use
 * this tool and you can instead manage your SDK packages from the IDE.
 * link: https://developer.android.com/studio/command-line/sdkmanager.html#usage
 */

/* list all available SDK package */
```
$ cd /usr/local/android-sdk-linux/tools/bin
$ ./sdkmanager --list
```
/* install the platform tools */

`$ ./sdkmanager "platform-tools"`

/* install the Android 2.3.3(API 10)
 * Why API 10?
 * No special reason, just because we tested it from beginning
 */

`$ ./sdkmanager "platforms;android-10"`

Create AVD
```
$ android -s create avd -n android-api-10-5555 -t android-10 --abi default/armeabi
Android 2.3.3 is a basic Android platform.
Do you wish to create a custom hardware profile [no]no
Created AVD 'android-api-10-5555' based on Android 2.3.3, ARM (armeabi) processor,
with the following hardware config:
hw.lcd.density=240
hw.ramSize=256
vm.heapSize=24
```

/*
 * Add lib to PATH
 * Why?
 * some dynamically linked libraries were moved around with the new Android emulator.
 * All you need to do is:
 *   add the folder with the libraries to the search path
 *   before you launch the emulator from command line.
 */
/* set environment */
```
$ sudo vi /etc/environment
export LD_LIBRARY_PATH="$ANDROID_SDK_HOME/tools/lib64:$ANDROID_SDK_HOME/tools/lib64/qt/lib:$LD_LIBRARY_PATH"

$ source /etc/environment

$ echo $LD_LIBRARY_PATH
/usr/local/android-sdk-linux/tools/lib64:/usr/local/android-sdk-linux/tools/lib64/qt/lib:
```

/*
 * create emulator
 */
```
$ emulator64-arm -avd android-api-10-5555 -wipe-data -no-window -no-boot-anim -noskin -port 5556
emulator: WARNING: the -no-skin flag is obsolete. to have a non-skinned virtual device, create one through the AVD manager
emulator: WARNING: Classic qemu does not support SMP. The hw.cpu.ncore option from your config file is ignored.
emulator: warning: opening audio output failed

emulator: Requested console port 5556: Inferring adb port 5557.
emulator: Listening for console connections on port: 5556
emulator: Serial number of this emulator (for ADB): emulator-5556
```
/*
 * connect emulator
 */

`$ adb devices`

## 2. Install GO
In the EC2, AWS has already installed go, so nothing to do here.
```
$ go version

go version go1.6.2 linux/amd64
```

## 3. Install noVNC (optional)
`$ sudo git clone git://github.com/kanaka/noVNC`

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
$ cat /etc/environment
PATH="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games:/usr/local/games:/usr/local/android-sdk-linux/tools:/usr/local/android-sdk-linux/platform-tools"
JAVA_HOME="/usr/lib/jvm/java-8-openjdk-amd64"
LD_LIBRARY_PATH="/usr/local/android-sdk-linux/tools/lib64:/usr/local/android-sdk-linux/tools/lib64/qt/lib"
GOPATH="/home/ubuntu/controller"
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
