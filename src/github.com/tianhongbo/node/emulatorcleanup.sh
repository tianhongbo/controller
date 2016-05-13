#!/bin/bash
ppid=$1
echo "ppid=$ppid"
for i in `ps -ef| awk '$3 == '${ppid}' { print $2 }'`
do
echo killing $i
kill -9 $i
done