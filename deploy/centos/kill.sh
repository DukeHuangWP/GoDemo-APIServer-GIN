#!/bin/sh

#暫存目前工作目錄
workingDir=$(pwd)

pidIntergTest=$(ps -A | grep "IntergTest" | awk '{print $1}')
if [[ -n "$pidIntergTest" ]]; then
    sudo kill -9 $pidIntergTest
fi