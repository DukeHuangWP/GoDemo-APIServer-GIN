#!/usr/bin/env bash


#用ip來判斷部屬環境
thisIP=`curl ident.me`
case "$thisIP" in
    123.123) thisServer="test"
        ;;
    456.789) thisServer="proc"
        ;;
    *) echo "$thisIP : 並非部屬 Server IP清單" & exit
esac

#暫存目前工作目錄
workingDir=`pwd`

# API設定檔調整
if [ ! -f "$workingDir/configs/config_"$thisServer".conf" ]; then
echo "遺失設定檔： $workingDir/configs/config_"$thisServer".conf"
exit
fi
cp -rf "$workingDir/configs/config_"$thisServer".conf" "$workingDir/configs/config.conf"
cd $workingDir

docker-compose stop IFT-Test

#更正權限,刪除不需要檔案
chmod 755 * -R
find . -type f -name ".DS_Store" -depth -exec rm -rfv {} \;

docker-compose restart IFT-Test 
