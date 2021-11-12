
# Dev-API
DEMO gin框架API

## 檔案樹說明(File tree)
```bash
.
│
├── build #專案編譯完成Server執行檔放置處
│   ├── BFSystem-API #執行檔
│   ├── assets
│   │   └── views #靜態網頁目錄
│   │       ├── CheckOutCounterError.html
│   │       └── CheckOutCounterSucces.html
│   └── configs #專案環境設定檔放置處
│       └── config.conf
│
├── configs #專案環境設定檔放置處(本機測試)
│   └── config.conf
├── docs #專案相關文件放置處
│   └── gorm_struct_tool.xlsx #gorm轉換struct方便工具
│
├── internal #golang main() 私有pkg
│   ├── alert #告警程序 - 將訊息發送至telegream上
│   │   └── alert.go
│   ├── common #常用函式放置處
│   │   ├── common.go
│   │   ├── common_test.go
│   │   └── customVar #自訂變數模組 - 可將字串轉換成任何基本類型
│   │       ├── customVar.go
│   │       ├── customVar_struct.go
│   │       └── customVar_test.go
│   ├── controllers #控制器 - 負責API處理客戶端請求/回傳
│   │   ├── controllers_for_common.go
│   │   ├── controllers_for_device.go
│   │   ├── controllers_for_task.go
│   │   ├── controllers_struct.go
│   │   ├── loggerToDB #logger模組 - 將log寫入至DB當中
│   │   │   ├── loggerToDB.go
│   │   │   ├── loggerToDB_filecache.go
│   │   │   ├── loggerToDB_table_creator.go
│   │   │   └── loggerToDB_test.go
│   │   └── task.go #舊專案亡魂等待刪除
│   ├── global #全域變數 - 放置/處理 
│   │   ├── config #環境設定值模組
│   │   │   ├── config.go
│   │   │   ├── config_model.go
│   │   │   ├── config_struct.go
│   │   │   └── config_test.go
│   │   └── global.go
│   ├── middleware #中間鍵 - 放置控制器請求與回傳之間的添加服務
│   │   └── antiflood #反SYN-Food模組
│   │       ├── antiflood.go
│   │       └── antiflood_test.go
│   ├── models #模型 - 負責控制器與DB之間的服務
│   │   ├── db.go #就行亡魂等待刪除
│   │   ├── define_protocol.go #舊行亡魂等待刪除
│   │   ├── encrypt #加密/解密模組
│   │   │   ├── encrypt.go
│   │   │   └── encrypt_test.go
│   │   ├── models_database.go #DB連接函式放置處
│   │   ├── models_sqls.go #DB SQL執行函式放置處
│   │   └── models_tables.go #DB Talbe定義放處
│   └── system #系統設定相關
│       ├── system.go
│       └── system_test.go
│
├── pkg #golang main() 公用pkg 
│   ├── src #專案外部引用pkg(相當於SET gopath)
│
├── main.go #golang main()整體專案入口
├── go.mod #golang版控相關
├── go.sum #golang版控相關
├── runBuild.sh #build執行檔腳本（編譯完成自動放置至./build）
└── runUnitTest.sh #執行所以pkg單元測試

```

## 建置設定(Build Setup)
1. 建議使用Golang ``1.15``以上版本
2. 需強制開啟go.mod使用 ``SET GO111MODULE="on"``

## 建置步驟(Build)
1. 使用腳本 ``./runBuild.sh``

## 注意事項(Reference)
1. 單元測試腳本  ``./runUnitTest.sh``

## 系統需求(System requirement)
 * OS: linux base
 * Container: up Docker-CE 3.19 ,Docker-Compose 3.8 (安裝步驟詳見後面)

 ## 前置設定(Server Setup)
1. docker-compose使用預設docker0網卡:
> 其網卡設置需要對外連接``DB``,
> 預設准許外部連接port為 ``8080``.

2. docker與使用者帳號確認
> ``sudo groupadd docker``

> ``sudo usermod -aG docker $USER``

> ``sudo systemctl restart docker``

3. docker 設定虛擬網卡``outside``指令如下:
> ``docker network create outside``


## Docker首次安裝筆記

```bash
sudo yum remove docker docker-common container-selinux docker-selinux docker-engine
sudo yum install -y yum-utils
sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo ##添加yum源
sudo yum clean all ##推荐先清空索引，特别是新版本docker需要换成老版本docker的时候
sudo yum makecache fast
yum list docker-ce --showduplicates | sort -r ##查看下自己能安装的版本都有哪些
sudo yum install -y docker-ce-18.09.9-3.el7 ##此处也可以安装指定版本的docker....
sudo yum install docker-ce ##不指定版本，安装最新版本的docker
sudo groupadd docker 
sudo usermod -aG docker $USER #將user加入可執行docker權限
sudo systemctl start docker 
sudo systemctl enable docker #设置开机启动
docker info ##验证是否安装成功

curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose
docker-compose --version
```

## git clone 設定方式

1. ssh必須設置本機金鑰 
`` nano ~/.ssh/config`` :
```bash
Host github.com
    HostName github.com
    User icsdev帳號
    IdentityFile ~/.ssh/git私密金鑰
```

2. 設置`` ~/.ssh/config``權限:
```bash
chmod 755 ~/.ssh/config
chmod 644 ~/.ssh/known_hosts
chmod 600 ~/.ssh/git私密金鑰
```

3. ``git clone git@github.com:帳號/XXX.git``

## 注意事項(Reference)
1. Docker log記得限制大小,避免log無止境肥大:
> 需要設定 ``nano  /etc/docker/daemon.json``
>```
>{
>  "log-driver": "local",
>  "log-opts": {
>    "max-size": "100m"
>  }
>}
>```
> 儲存後並重新啟用Docker ``systemctl restart docker.service``

2. 一些查看服務運作狀態指令:
>```bash
>docker ps #查看容器服務運行時間
>docker-compose ps #查看容器組啟用狀態
>docker-compose logs --tail 50 #查看容器最後的50行打印訊息
>docker-compose logs | grep "嚴重錯誤" #查看容器log "嚴重錯誤" 訊息
>```

3. 重啟所有docker容器方式:
>```bash
>docker-compose stop #暫停所有容器組(注意:docker-compose down會刪除容器與log)
>docker-compose up -d #背景執行,避免CLI影響容器運作
>```

4. 安裝docker-compose碰上python版本問題:
```
sudo yum remove python2 #移除舊版python
python #確認是否正常移除

sudo yum install -y python3 #安裝python3
sudo mv /usr/bin/python /usr/bin/python2.7 #若無法正確移除python2,則強制修改路徑
sudo ln -s /usr/bin/python3.6 /usr/bin/python #將python3設為預設python
sudo ln -s /usr/bin/python2.7 /usr/bin/python #由於yum必須使用python2.7,所以之後需要改回來
