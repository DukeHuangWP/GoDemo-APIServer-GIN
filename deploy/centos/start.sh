#!/bin/sh

#暫存目前工作目錄
workingDir=$(pwd)
log_file=$workingDir/file.log

pidIntergTest=$(ps -A | grep "IntergTest" | awk '{print $1}')
if [[ -n "$pidIntergTest" ]]; then
    sudo kill -9 $pidIntergTest
fi

count=$(ps -ef | grep "IntergTest" | grep -v "grep" | wc -l)
if [ ${count} -lt 1 ]; then
    #echo "$(date +'%Y/%m/%d %H:%M:%S') start IntergTest....." >>$log_file
    nohup ./IntergTest &>$log_file &
fi

: <<BLOCK
BLOCK

#sudo ps -ef | grep "./BFSystem-IntergTest"

# nohup /path/my_program >my.out 2>my.err & 

# nohup "./IntergTest" 2>&1 my_log.txt &