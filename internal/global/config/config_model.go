package config

import (
	"DevIntergTest/internal/common/customVar"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

/*
量測並建立SettingMap
@param needClear: 是否強制清除SettingMap
*/
func CreatSettingMap(needClear bool) {
	if needClear || len(SettingMap) <= 0 {
		SettingMap = make(map[string]*AddSetings) //儲存環境變數,可依需求自行擴充
	}
	return
}

/*
	將原先設定好輸入設定值字串寫入變數
	@param inputValue: 輸入設定值字串
	@return error : 錯誤提示
*/
func (setings *AddSetings) ToWrite(inputValue string) error {
	setings.IsWritten = true
	return customVar.SetConfig(inputValue, setings.DefaultValue, setings.OutputPointer, setings.Custom)
}

/*
	將原先設定好輸入預設值字串寫入變數
	@return error : 錯誤提示
*/
func (setings *AddSetings) ToDefault() error {
	return customVar.SetValue(setings.DefaultValue, setings.OutputPointer, setings.Custom)
}

/*
讀取環境變數設定檔
@param configPath: 設定檔路徑
@param isInit: 是否第首次執行(遇錯則關閉整個程序)
*/
func (cfgs *CfgType) SetFromFile(configPath string, isInit bool) {

	if !isInit || configPath == "" {
		configPath = Cfgs.ConfigPath
	}

	configBtye, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Printf("讀取環境設定值失敗: %v (%v)\n", configPath, err)
		if isInit {
			os.Exit(Code_error_exit)
		}
	}

	Cfgs.ConfigPath = configPath

	// →#註解← 一般注解方式
	// →  #註解  ←設定值前後可留半形空白
	// →;註解← 一般注解方式
	// →  ;註解  ←設定值前後可留半形空白
	// →[標題← 一般注解方式
	// →  [標題  ←設定值前後可留半形空白
	// 暫定標題與註解效果相同
	// →settingKey1=value1← 一般設定方式
	// →settingKey2 = value1← 錯誤設定方式

	configStr := string(configBtye)
	parsingChr := "\n"

	for index, value := range strings.Split(configStr, parsingChr) {

		value = strings.TrimRight(value, "\r") //清除尾部多餘換行符號
		value = strings.Trim(value, " ")       //清除前後多餘空白
		if indexComment := strings.Index(value, "#"); indexComment >= 0 && indexComment < 3 {
			continue //此#為設定的註解，故忽略
		}

		if indexComment := strings.Index(value, ";"); indexComment >= 0 && indexComment < 3 {
			continue //此;為設定的註解，故忽略
		}

		if indexComment := strings.Index(value, "["); indexComment >= 0 && indexComment < 3 {
			continue //此[暫定為註解，故忽略
		}

		if indexSetChar := strings.Index(value, "="); indexSetChar < 1 {
			continue //此狀況設定值"="位置錯誤或無"=""
		} else {

			settingKey := value[:indexSetChar]
			settingValue := value[indexSetChar+1:]

			if _, isExsit := SettingMap[settingKey]; !isExsit {
				log.Printf("環境設定值錯誤警告: %v (Line:%v) > %v，檢測到無效設定!", configPath, index, value)
				continue
			}

			if err := SettingMap[settingKey].ToWrite(settingValue); err != nil {
				log.Printf("環境設定值錯誤警告: %v (Line:%v) > %v (%v)，將使用預設值'%v'代替.", configPath, index, value, err, SettingMap[settingKey].DefaultValue)
			}

		}
	}

	for index, value := range SettingMap { //若環境變數無此值則使用預設值
		if value.IsWritten == false {
			if err := value.ToDefault(); err != nil {
				log.Printf("環境設定值錯誤警告: %v 中並未發現'%v', 在設置預設值'%v'發生錯誤'%v', 煩檢查設定值.", configPath, index, value.DefaultValue, err)
			} else {
				log.Printf("環境設定值錯誤警告: %v 中並未發現'%v', 將使用預設值'%v'代替.", configPath, index, value.DefaultValue)
				// if isInit {
				// 	os.Exit(Code_error_exit)
				// }
			}
		}
	}
}

/*
讀取環境變數
@param envTitle: 環境設定值前綴字
@param isInit: 是否第首次執行(遇錯則關閉整個程序)
*/
func (cfgs *CfgType) SetFromEnv(envTitle string) {

	if envTitle == "" {
		envTitle = Cfgs.EnvTitle
	}

	envVarList := os.Environ()
	if len(envVarList) <= 0 {
		log.Printf("讀取環境設定值失敗: 無法抓取到任何環境變數\n")
	}

	envTitleCount := len(envTitle)
	for _, envVar := range envVarList {

		envSlice := strings.SplitN(envVar, "=", 2)
		//log.Printf("'%v'='%v'\n",envSlice[0],envSlice[1])
		settingKey := envSlice[0]
		settingValue := envSlice[1]

		if settingKey == envTitle+"SetFromEnv" {
			continue
		} //此為環境設定值f啟用開關故忽略

		if !strings.HasPrefix(settingKey, envTitle) {
			continue
		}

		settingKey = envSlice[0][envTitleCount:]
		if _, isExsit := SettingMap[settingKey]; !isExsit {
			log.Printf("環境設定值錯誤警告: %v=%v，檢測到無效設定!", settingKey, settingValue)
			continue
		}

		if err := SettingMap[settingKey].ToWrite(settingValue); err != nil {
			log.Printf("環境設定值錯誤警告: %v=%v (%v)，將使用預設值'%v'代替.", settingKey, settingValue, err, SettingMap[settingKey].DefaultValue)
		}
	}

	for index, value := range SettingMap { //若環境變數無此值則使用預設值
		if value.IsWritten == false {
			if err := value.ToDefault(); err != nil {
				log.Printf("環境設定值錯誤警告: 讀取環境變數中並未發現'%v', 在設置預設值'%v'發生錯誤'%v', 煩檢查設定值.", index, value.DefaultValue, err)
			} else {
				log.Printf("環境設定值錯誤警告: 讀取環境變數中並未發現'%v', 將使用預設值'%v'代替.", index, value.DefaultValue)
				// if isInit {
				// 	os.Exit(Code_error_exit)
				// }
			}
		}
	}
}
