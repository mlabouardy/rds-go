#!/bin/sh
yum update
yum install -y golang
mkdir -p /home/ec2-user/go/src/github.com/
cd /home/ec2-user/go/src/github.com/
git clone
go get ./...
