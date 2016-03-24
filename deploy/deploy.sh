#!/bin/bash

curPath=`dirname $0`
cd $curPath/../
prjHome=`pwd`

if [ $# -lt 1 ]
then
    echo "Usage $0 host1 host2 ..."
    exit 1
fi


echo "Make sure you have sudo permission of online hosts"
echo "Enter your password of online hosts: "
read -s password
sshCmd="sshpass -p $password ssh -o StrictHostKeyChecking=no"
scpCmd="sshpass -p $password scp -o StrictHostKeyChecking=no"


echo "Building binary"
deployTmpDir=$prjHome/tmp/deploy
if [ -d $deployTmpDir ]
then
    rm -rf $deployTmpDir
fi
mkdir -p $deployTmpDir

binName=rigger
cd $deployTmpDir
go build -o $binName $prjHome/main/rigger.go


installDstDir=/usr/local/bin
for host in $*
do
    echo Deploy to $host
    $scpCmd $binName $host:./
    $sshCmd -t $host "echo $password | sudo -S mv $binName $installDstDir"
done


rm -rf $deployTmpDir
