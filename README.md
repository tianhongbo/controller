# controller
This is a controller written by Go to create, delete a Android emulator based Android SDK.

# How to install
## Step 1: Install Android SDK
- install Java SDK
- install Android Studio
- install Android 2.3.3(API 10)

## Step 2: Install GO

## Step 3: Install adb (even though it has been installed together with Android SDK/Studio)

`$ sudo apt-get install android-tools-adb`

Step 4: Install noVNC
$ sudo git clone git://github.com/kanaka/noVNC

Step 5: Git Clone NODE source code
$ git clone https://github.com/tianhongbo/node.git

Step 6: Modify source code
- install.sh
- deviceinstall.sh
- repo.go

Step 7: Set environment variables
$ sudo vi /etc/environment

- add emulator tools to PATH
/home/ubuntu2/Android/Sdk/tools

- add GOPATH
GOPATH=/home/ubuntu2/controller

Step 8: configure port forwarding for SSH functions
- iptables for Ubuntu
- ip for mac

For example(for ubuntu):
# This is for eth0
- sysctl -w net.ipv4.conf.eth0.route_localnet=1
- iptables -t nat -I PREROUTING -p tcp -i eth0 --dport 5921 -j DNAT --to-destination 127.0.0.1:5921

# This is for wlan
- sysctl -w net.ipv4.conf.wlan0.route_localnet=1
- iptables -t nat -I PREROUTING -p tcp -i wlan0 --dport 5921 -j DNAT --to-destination 127.0.0.1:5921


#Done

#Debugging
1. check adb forward
- $ adb forward --list


2. check iptables NAT
- $ sudo iptables -t nat -L -n -v

#SSH
# user: "test" / password: "test"
ssh test@192.168.1.24 -p 5921

#VNC
It seems the behaviors in MAC and Ubuntu are different. For example, VNC launch.sh only create
listening port on the external interface in MAC. However, it will create listening port on all
the interfaces including localhost.

2. /etc/rc.local
start node automatically when instance booting
su ubuntu - -c /home/ubuntu/node/bin/node
