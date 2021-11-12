package global

import (
	"DevIntergTest/internal/common"
	"DevIntergTest/internal/common/customVar"
	"DevIntergTest/internal/global/config"
	"DevIntergTest/internal/system"
	"fmt"
	"log"
	"reflect"
)

const (
	ConfigFilePath = "./configs/config.conf" //環境變數設定檔
	ConfigEnvTitle = "DevIntergTest_"        //環境設定值前綴字

	URNRoot              = ""          //網頁根目錄 http://localhost:8080/xxx
	FolderPath_Public    = "./website" //網頁公用檔案目錄
	FolderPath_Images    = "images"    //網頁公用目錄名稱 - 圖片
	FolderPath_Pages     = "pages"     //網頁公用目錄名稱 - 純html頁面
	FolderPath_Templs    = "templates" //網頁公用目錄名稱 - 版模
	FolderPath_Downloads = "downloads" //網頁公用目錄名稱 - 下載目錄
	FolderPath_Uploads   = "uploads"   //網頁公用目錄名稱 - 上傳目錄
	FolderPath_Others    = "others"    //網頁公用目錄名稱 - 其他檔案分類

	// URNPath_APITest = "/api_test"

	// URNPath_Login = "/login" //登陸Path名稱
	// URNPath_Upload = "/upload" //上傳檔案Path名稱
)

var (
	ServiceName string
	ServicePort string              //= "8080"
	SettingMap  = config.SettingMap //參照package config

	DBName     string
	DBHost     string
	DBPort     uint
	DBUserName string
	DBPassword string

	OutboundIPChecker        string //獲取對外ip測試網址
	ServiceOutboundIP        string //對外ip
	LocalHostUrl             string //本機網址
	TelegramSendURL          string //TG發送地址
	TelegramTitle_Format     string //TG發送前格式
	TelegramSendMessageLimit int    //TG發送文本長度現住

	TestCaseReloadSecs int //資料庫測項重新讀取間格時間(秒)

	APIProxyHeaderClientIP string //獲取代理服務器之ClientIP之header名稱
	APICallbackServiceName string //API回調服務名稱
	APICallbackMaxNum      int    // API傳送回調最高上限次數

	APILimiterMaxCount          uint64   //短時間連線次數上限
	APILimiterRateSecs          int64    //定義短時間秒數
	APILimiterClientIPUnbanSecs int64    //遭禁ClientIP自動解封時間
	APILimiterTimerUnbanSecs    int      //自動解封時間執行周期
	APILimiterAllowIPList       []string //限制器白名單IP

)

func init() {

	//------------設置環境設定值到全域變數中
	//環境變數設定值解析 SettingMap["XXX"] = &config.AddSetings{初始值, 輸出變數指標, 自訂設定值類型}
	SettingMap["service_name"] = &config.AddSetings{DefaultValue: "IntergTest", OutputPointer: &ServiceName, Custom: &customVar.StringType{}}
	SettingMap["service_port"] = &config.AddSetings{DefaultValue: "9090", OutputPointer: &ServicePort, Custom: &customVar.Uint16Type{}} //設定變數 ,預設值: "8080"
	SettingMap["db_name"] = &config.AddSetings{DefaultValue: "", OutputPointer: &DBName, Custom: &customVar.StringType{}}
	SettingMap["db_host"] = &config.AddSetings{DefaultValue: "0.0.0.0", OutputPointer: &DBHost, Custom: &customVar.ValidIPnPort{}}
	SettingMap["db_port"] = &config.AddSetings{DefaultValue: 3306, OutputPointer: &DBPort, Custom: &customVar.Uint16Type{}}
	SettingMap["db_username"] = &config.AddSetings{DefaultValue: "", OutputPointer: &DBUserName, Custom: &customVar.StringType{}}
	SettingMap["db_password"] = &config.AddSetings{DefaultValue: "", OutputPointer: &DBPassword, Custom: &customVar.StringType{}}
	SettingMap["outbound_ip_checker"] = &config.AddSetings{DefaultValue: "ifconfig.me", OutputPointer: &OutboundIPChecker, Custom: &customVar.ValidWebURL{}}
	SettingMap["localhost_url"] = &config.AddSetings{DefaultValue: "http://127.0.0.1", OutputPointer: &LocalHostUrl, Custom: &customVar.ValidWebURL{}}
	SettingMap["telegram_send_url"] = &config.AddSetings{DefaultValue: "", OutputPointer: &TelegramSendURL, Custom: &customVar.StringType{}}
	SettingMap["telegram_send_message_limit"] = &config.AddSetings{DefaultValue: 1024, OutputPointer: &TelegramSendMessageLimit, Custom: &customVar.Uint32Type{}}

	SettingMap["test_case_reload_secs"] = &config.AddSetings{DefaultValue: 300, OutputPointer: &TestCaseReloadSecs, Custom: &customVar.SecondsInADay{}}

	// SettingMap["device_server_url"] = &config.AddSetings{DefaultValue: "http://127.0.0.1", OutputPointer: &DeviceServerURL, Custom: &customVar.ValidWebURL{}}
	// SettingMap["device_server_apiname_collectiontask"] = &config.AddSetings{DefaultValue: "/Task/CollectionTask", OutputPointer: &DeviceServerAPINameCollectionTask, Custom: &customVar.StringType{}}
	// SettingMap["device_server_apiname_addtask"] = &config.AddSetings{DefaultValue: "/Task/CollectionTask", OutputPointer: &DeviceServerAPINameAddTask, Custom: &customVar.StringType{}}
	// SettingMap["device_server_apiname_callback_bankcard"] = &config.AddSetings{DefaultValue: "/Task/CollectionTask", OutputPointer: &DeviceServerAPINameCallbackBankCard, Custom: &customVar.StringType{}}
	// SettingMap["device_server_apiname_API_notify_statementtask"] = &config.AddSetings{DefaultValue: "/API/NotifyStatementTask", OutputPointer: &DeviceServerAPINameAPINotifyStatementTask, Custom: &customVar.StringType{}}
	SettingMap["api_proxy_header_clientip"] = &config.AddSetings{DefaultValue: "Cf-Connecting-Ip", OutputPointer: &APIProxyHeaderClientIP, Custom: &customVar.StringType{}}
	SettingMap["api_limter_max_count"] = &config.AddSetings{DefaultValue: 3000, OutputPointer: &APILimiterMaxCount, Custom: &customVar.Uint64Type{}}
	SettingMap["api_limter_rate_secs"] = &config.AddSetings{DefaultValue: 10, OutputPointer: &APILimiterRateSecs, Custom: &customVar.SecondsInADay{}}
	SettingMap["api_limter_client_unban_sec"] = &config.AddSetings{DefaultValue: 600, OutputPointer: &APILimiterClientIPUnbanSecs, Custom: &customVar.SecondsInADay{}}
	SettingMap["api_limter_timer_unban_sec"] = &config.AddSetings{DefaultValue: 1800, OutputPointer: &APILimiterTimerUnbanSecs, Custom: &customVar.SecondsInADay{}}
	SettingMap["api_limter_allow_IPList_add"] = &config.AddSetings{DefaultValue: nil, OutputPointer: &APILimiterAllowIPList, Custom: &customVar.BaseSlice{}}
	// SettingMap["api_cron_timer_statement_task_short_min"] = &config.AddSetings{DefaultValue: 3, OutputPointer: &APICronTimerStatemenTaskShortMins, Custom: &customVar.Int16Type{}}
	// SettingMap["api_cron_timer_statement_task_long_min"] = &config.AddSetings{DefaultValue: 5, OutputPointer: &APICronTimerStatemenTaskLongMins, Custom: &customVar.Int16Type{}}
	// SettingMap["loggerToDB_TimerSec"] = &config.AddSetings{DefaultValue: 15, OutputPointer: &LoggerToDBTimerSec, Custom: &customVar.SecondsInADay{}}
	// SettingMap["loggerToDB_DebugSwitch"] = &config.AddSetings{DefaultValue: false, OutputPointer: &LoggerToDBDebug, Custom: &customVar.SwitchType{}}
	// SettingMap["loggerToDB_FilePath"] = &config.AddSetings{DefaultValue: "", OutputPointer: &LoggerToDBCacheFilePath, Custom: &customVar.StringType{}}

	//os.Setenv(ConfigEnvTitle+"SetFromEnv", "on") //本機測試用強制打開環境變數模式
	config.InitLoad(ConfigFilePath, ConfigEnvTitle)   //載入環境設定檔案
	SettingMapDisplay := make(map[string]interface{}) //作為初始化顯示值用
	for key, value := range SettingMap {
		SettingMapDisplay[key] = fmt.Sprint(reflect.ValueOf(value.OutputPointer).Elem()) //反射該指標之值
	} //顯示所有最終環境設定值
	log.Printf("%v\n", SettingMapDisplay)
	SettingMapDisplay = nil //CG優化
	//============

	//------------設置Telegram抬頭名稱格式
	ServiceOutboundIP = system.GetOutboundIP(OutboundIPChecker)    //獲取對外ip
	displayIP := common.ToCoverLeftToRightIP(ServiceOutboundIP, 2) //遮蔽部分IP
	TelegramTitle_Format = "[🅰️PI][%v](" + displayIP + ")"
	//===========

	// 目錄靜態資源目錄
	// if err := common.CreatFolderIfNotExist(FolderPath_Public + "/" + FolderPath_Uploads); err != nil {
	// 	log.Fatal(err)
	// } //檢查並建立webserver上傳目錄

}
