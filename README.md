# controller
This is a controller written by Go to create, delete a Android emulator based Android SDK on Mac OS, Linux Desktop, or Linux Server.

# setup
AWS EC2: free tier (t2.micro)
Linux Ubuntu: Linux ip-172-31-16-251 4.4.0-53-generic #74-Ubuntu SMP Fri Dec 2 15:59:10 UTC 2016 x86_64 x86_64 x86_64 GNU/Linux

# How can I install it?
## 1. Install Android SDK
- install Java SDK (java-8-openjdk-amd64)
- install Android SDK Tool (25.2.3)
- install Android 2.3.3(API 10)

## 2. Install GO

## 3. Install adb
even though it has been installed together with Android SDK/Studio
`$ sudo apt-get install android-tools-adb`

## 4. Install noVNC
$ sudo git clone git://github.com/kanaka/noVNC

## Step 5: Git Clone NODE source code
$ git clone https://github.com/tianhongbo/controller.git

# How can I configure it?
## 1. Modify source code
- install.sh
- deviceinstall.sh
- repo.go

## 2. Set environment variables
$ sudo vi /etc/environment

- add emulator tools to PATH
/home/ubuntu2/Android/Sdk/tools

- add GOPATH
GOPATH=/home/ubuntu2/controller

## 3. configure port forwarding for SSH functions
- iptables for Ubuntu
- ip for mac

For example(for ubuntu):
### This is for eth0
- sysctl -w net.ipv4.conf.eth0.route_localnet=1
- iptables -t nat -I PREROUTING -p tcp -i eth0 --dport 5921 -j DNAT --to-destination 127.0.0.1:5921

### This is for wlan
- sysctl -w net.ipv4.conf.wlan0.route_localnet=1
- iptables -t nat -I PREROUTING -p tcp -i wlan0 --dport 5921 -j DNAT --to-destination 127.0.0.1:5921

# How can I debug?
## 1. check adb forward
- $ adb forward --list


## 2. check iptables NAT
- $ sudo iptables -t nat -L -n -v

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
