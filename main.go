package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"DevIntergTest/internal/controllers"
	"DevIntergTest/internal/controllers/middleware/antiflood"
	"DevIntergTest/internal/cronjob"
	"DevIntergTest/internal/global"
	"DevIntergTest/internal/models"

	"DevIntergTest/internal/logger/loggerToDB"

	"github.com/gin-gonic/gin"
)

func main() {

	var err error

	//--------初始化model 資料庫相關模組---------------{
	err = models.InitDatabase(&models.MySQLConfig{Account: global.DBUserName,
		Password:           global.DBPassword,
		DBHost:             global.DBHost,
		DBName:             global.DBName,
		MaxIdleConns:       1,
		MaxOpenConns:       5,
		ConnMaxLifeTimeSec: 300,
		DBCacheTimerSec:    300,
	})
	if err != nil {
		log.Fatalf("%v 初始化失敗 > %v", models.ServiceName, err)
	}
	models.StartKeepDBAlive(nil, 300) //每隔'%v'秒對db進行檢查測試
	//------------END--------------------------------}

	//--------初始化anti SYN-FOOD模組------------------{
	antiflood.StartTimer(&antiflood.Cfgs{
		MaxCount:          global.APILimiterMaxCount,
		RateSecs:          global.APILimiterRateSecs,
		ClientIPUnbanSecs: global.APILimiterClientIPUnbanSecs,
		TimerUnbanSecs:    global.APILimiterTimerUnbanSecs,
		ProxyHeaderKey:    global.APIProxyHeaderClientIP,
	}) //啟動反SYN-Flood計時器

	if antiflood.SetAllowIPRule(global.APILimiterAllowIPList) != nil {
		log.Fatal("IP白名設定失敗 > ", err.Error())
	} else if len(global.APILimiterAllowIPList) > 0 && (len(antiflood.ClientIPAllowList)+len(antiflood.ClientIPRuleList)) > 0 {
		log.Printf("%v : 載入開放無限制使用API > 准許IP清單:%v 准許IP規則:%v\n", antiflood.Service_Name, antiflood.ClientIPAllowList, antiflood.ClientIPRuleList)
	}
	//------------END--------------------------------}

	//--------初始化API Router(gin)模組----------------{
	ginRouter := gin.Default()
	//---- gin框架範例
	ginRouter.Use(func(ginCTX *gin.Context) {
		//ginRouter.Use中每個路徑都會先加載func相當於Middleware中間鍵的功能
		antiflood.TriggerHandler(ginCTX) //反SYN-Flood
	})

	//---- gin框架範例
	//ginRouter.Static("BoswerPath", "./Public") //Static()會將目錄內檔案公開出去，較不安全
	//ginRouter.StaticFS("BoswerPath", http.Dir("./Public")) //StaticFS()可指定目錄內哪些檔案可以被公開
	//ginRouter.StaticFile(global.FolderPath_Images, "./Public/gabo.png") //只能公開一個檔案

	ginRouter.LoadHTMLGlob(global.FolderPath_Public + "/" + global.FolderPath_Templs + "/*.html") //.LoadHTMLGlob()僅會在板模後傳給客戶端，該html檔案並不會被公開
	ginGroup1 := ginRouter.Group(controllers.URN_DevToolRootPath)                                 //表示http://{URI}/DevTool
	{
		//StaticFS()可指定目錄內哪些檔案可以被公開
		//ginRouter.Static("BoswerPath", "./Public") //Static()會將目錄內檔案公開出去，較不安全
		//ginRouter.StaticFS("BoswerPath", http.Dir("./Public")) //StaticFS()可指定目錄內哪些檔案可以被公開
		//ginRouter.StaticFile(global.FolderPath_Images, "./Public/gabo.png") //只能公開一個檔案

		ginGroup1.StaticFS(controllers.URN_DevToolRootPath, http.Dir(global.FolderPath_Public+"/"+global.FolderPath_Pages)) //公開靜態頁面
		ginGroup1.StaticFS(controllers.FolderPath_Images, http.Dir(global.FolderPath_Public+"/"+global.FolderPath_Images))
		ginGroup1.StaticFS(controllers.FolderPath_Others, http.Dir(global.FolderPath_Public+"/"+global.FolderPath_Others))
		ginGroup1.StaticFS(controllers.FolderPath_Downloads, http.Dir(global.FolderPath_Public+"/"+global.FolderPath_Downloads))
	}

	ginGroup2 := ginRouter.Group(controllers.URN_ValidRootPath) //表示http://{URI}/Valid
	{
		ginGroup2.Any(controllers.URN_Valid_fundlog, controllers.MidValid, controllers.EndNodeValidTask)
	}

	ginGroup3 := ginRouter.Group(controllers.URN_SettingsPath) //表示http://{URI}/Settings
	{
		ginGroup3.GET("", controllers.LoadSettrings)
		//ginGroup3.POST(":testName", controllers.SaveSettrings)
	}
	//------------END---------------------------------}

	//--------初始化API Router(gin)模組-----------------{
	gin.SetMode("release") //debug mode會有較多的性能資訊
	apiServ := &http.Server{
		Addr:           ":" + global.ServicePort, //host port
		Handler:        ginRouter,                //ginRouter
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20, //計算原理不明,限制Request接收長度
	}

	go func() {
		if err := apiServ.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("HTTP端口監聽失敗 > ", err.Error())
		}
	}()
	//------------END---------------------------------}

	//--------初始化corn 定時器模組----------------------{
	cronjob.StartTimer()

	//------------END---------------------------------}

	//alert.SendMessageWithTime(fmt.Sprintf("服務已上線 > %v", common.GetUnixNowSec())) //發出告警說明系統已上線成功

	//========================往下為關機後觸發func========================================{
	quitChan := make(chan os.Signal, 10) //Notify：系統訊號轉將發至channel
	signal.Notify(quitChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-quitChan
	//永久等待收到系統發出訊號類型,若符合則往下執行
	close(quitChan)

	shutdownCTX, shutdownClose := context.WithTimeout(context.Background(), 8*time.Second)
	defer shutdownClose() //透過context.WithTimeout產生一個新的子context，它的特性是有生命週期，只要超過8秒就會自動發出Done()的訊息
	if err := apiServ.Shutdown(shutdownCTX); err != nil {
		log.Fatal("API關閉服務過程中發生錯誤:", err)
	}

	loggerToDB.AddMessage(&loggerToDB.Message{
		ServerName:      global.ServiceName,
		APIErrorMessage: "收到關閉訊號...",
	}) //若api意外結束則回傳此訊息
	println()
	log.Printf("正在關閉API服務... \n")
	//loggerClose()                                                                  //催速logger快點寫入
	//alert.SendMessageWithTime(fmt.Sprintf("服務已遭受關閉 > %v", common.GetUnixNowSec())) //TG發出告警

	<-shutdownCTX.Done() //當子context發出Done()的訊號才繼續向下走
	log.Printf("API服務已關閉！遺失[%v]筆紀錄.\n", len(loggerToDB.QueueChan))
	//====================================================================================}

}
