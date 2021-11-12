package config

import (
	"DevIntergTest/internal/common/customVar"
	"log"
	"os"
)

func init() {
	CreatSettingMap(false)
	SettingMap["cpu_core"] = &AddSetings{DefaultValue: 0, OutputPointer: &Cfgs.CpuCore, Custom: &customVar.Uint8Type{}}
	//SettingMap["XXX"] = &AddSetings{初始值, 輸出變數指標, 自訂設定值類型}
	return
}

//初始化環境變數設定值
func InitLoad(configPath string, envTitle string) {
	Cfgs.EnvTitle = envTitle
	var isEnvMode bool
	err := customVar.SetConfig(os.Getenv(envTitle+EnvTitle_Switch), false, &isEnvMode, &customVar.SwitchType{})
	if err != nil || isEnvMode == false { //此狀況為關閉使用環境變數
		log.Printf("直接由檔案讀取環境設定值: %v", configPath)
		Cfgs.SetFromFile(configPath, true)
	} else {
		log.Printf("系統將擷取環境設定值開頭為: '%v'", envTitle)
		Cfgs.SetFromEnv(envTitle)
	}

	return
}

//重新讀取環境變數設定值
func Reload() {

	CreatSettingMap(false)
	var isEnvMode bool
	err := customVar.SetConfig(os.Getenv(Cfgs.EnvTitle+EnvTitle_Switch), false, &isEnvMode, &customVar.SwitchType{})
	if err != nil || isEnvMode == false { //此狀況為關閉使用環境變數
		Cfgs.SetFromFile("", true)
		log.Printf("直接由檔案讀取環境設定值 (%v)", Cfgs.ConfigPath)
	} else {
		Cfgs.SetFromEnv("")
		log.Printf("重新讀取環境設定值完畢 ('%v')", Cfgs.EnvTitle)
	}

	return
}
