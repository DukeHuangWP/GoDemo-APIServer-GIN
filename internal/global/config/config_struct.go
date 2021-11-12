package config

import (
	"DevIntergTest/internal/common/customVar"
)

const (
	EnvTitle_Switch     = "SetFromEnv"
	Code_error_exit = 3 //強制關閉code
)

type AddSetings struct {
	IsWritten     bool                  //輸出變數是否已寫入
	DefaultValue  interface{}           //輸出預設值(輸入值錯誤時使用)
	OutputPointer interface{}           //輸出變數(必須為指標)
	Custom        customVar.CustomTypes //自訂設定值類型(工廠方法GetValue判定)
}

var (
	SettingMap map[string]*AddSetings //環境變數設定值暫存

	Cfgs   CfgType    //環境變數設置全域變數
	Public PublicType //公用 ex:本地Redis連接池
)

type CfgType struct {
	ConfigPath string //暫存上次讀取設定值路徑
	EnvTitle string //暫存上次讀取環境設定值前綴字
	CpuCore    int    //cpu核心數
}

type PublicType struct {
	//公用config,如:db連接池...
	PrintExampleIndex int //Demo:打印測試初始值
}
