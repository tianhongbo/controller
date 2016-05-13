#!/usr/bin/env bash
sudo sysctl -w net.ipv4.conf.wlan0.route_localnet=1

for ssh_port in {5921..5940}
do
    echo "sudo iptables -t nat -I PREROUTING -p tcp -i wlan0 --dport $ssh_port -j DNAT --to-destination 127.0.0.1:$ssh_port"
    sudo iptables -t nat -I PREROUTING -p tcp -i wlan0 --dport $ssh_port -j DNAT --to-destination 127.0.0.1:$ssh_port
done
