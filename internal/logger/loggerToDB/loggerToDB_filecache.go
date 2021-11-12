package loggerToDB

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

const (
	CacheFileSplitTag = "⸽" //log暫存檔案資料間格符號
	CacheFileLineTag  = "" //log暫存檔案資料換行符號
)

var (
	theCacheFile *os.File
)

//獲取暫存檔案並回傳尚未寫入db中的log
func GetMessageFromFile(cacheFilePath string) (thisFile *os.File, lastMessage []*Message, err error) {

	thisFile, err = os.OpenFile(cacheFilePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644) //打開或創建檔案
	fileStats, err := thisFile.Stat()
	if err != nil {
		return
	}

	var fileSize int64 = fileStats.Size()
	cacheFileBtyes := make([]byte, fileSize)

	reader := bufio.NewReader(thisFile)
	_, err = reader.Read(cacheFileBtyes)

	cacheFileText := *(*string)(unsafe.Pointer(&cacheFileBtyes)) //快速轉換成string
	cacheFileBtyes = nil                                         //CG優化

	if strings.Contains(cacheFileText, CacheFileSplitTag) == false {
		return
	}

	messageStruct := reflect.ValueOf(&Message{}).Elem() //反射
	messageStructField := messageStruct.NumField() + 1  //反射數量由0開始算... why?

	cacheFileTextList := strings.Split(cacheFileText, CacheFileLineTag) //解析換行符號
	for index := 0; index < len(cacheFileTextList); index++ {

		cacheFileTextline := strings.Split(cacheFileTextList[index], CacheFileSplitTag) //解析log分段符號

		if len(cacheFileTextline) != messageStructField {
			continue
		}

		apiStatus, _ := strconv.Atoi(cacheFileTextline[8])
		responseHttpStatus, _ := strconv.Atoi(cacheFileTextline[10])
		responseCosttime, _ := strconv.ParseInt(cacheFileTextline[11], 10, 64)
		logTime, _ := strconv.ParseInt(cacheFileTextline[17], 10, 64)

		lastMessage = append(lastMessage, &Message{
			ServerName:          cacheFileTextline[0],
			ServerHost:          cacheFileTextline[1],
			ClientIP:            cacheFileTextline[2],
			RequestURN:          cacheFileTextline[3],
			RequestMethod:       cacheFileTextline[4],
			RequestAccount:      cacheFileTextline[5],
			RequestDataDecrypt:  cacheFileTextline[6],
			RequestBody:         cacheFileTextline[7],
			APIStatus:           apiStatus,
			APIErrorMessage:     cacheFileTextline[9],
			ResponseHttpStatus:  responseHttpStatus,
			ResponseCosttime:    responseCosttime,
			ResponseAccount:     cacheFileTextline[12],
			ResponseMsgCode:     cacheFileTextline[13],
			ResponseMsg:         cacheFileTextline[14],
			ResponseDataDecrypt: cacheFileTextline[15],
			ResponseDataEncrypt: cacheFileTextline[16],
			//Debug資訊取消寫入
			LogTime: logTime,
		})

	}

	return
}

//將log寫入到暫存檔案當中
func AddMessageToFile(thisFile *os.File, rowData gormRow) (err error) {

	var writeText strings.Builder

	// 取消使用反射
	// value := reflect.ValueOf(row)
	// for index := 0; index < value.NumField(); index++ {
	// 	if index == 17 {
	// 		continue //DebugMessage取消不補紀錄
	// 	}
	// 	writeText.WriteString(fmt.Sprint(value.Field(index)) + CacheFileSplitTag)
	// }
	// writeText.WriteString(CacheFileLineTag)

	writeText.WriteString(fmt.Sprint(rowData.ServerName) + CacheFileSplitTag)
	writeText.WriteString(fmt.Sprint(rowData.ServerHost) + CacheFileSplitTag)
	writeText.WriteString(fmt.Sprint(rowData.ClientIP) + CacheFileSplitTag)
	writeText.WriteString(fmt.Sprint(rowData.RequestURN) + CacheFileSplitTag)
	writeText.WriteString(fmt.Sprint(rowData.RequestMethod) + CacheFileSplitTag)
	writeText.WriteString(fmt.Sprint(rowData.RequestAccount.String) + CacheFileSplitTag)
	writeText.WriteString(fmt.Sprint(rowData.RequestDataDecrypt.String) + CacheFileSplitTag)
	writeText.WriteString(fmt.Sprint(rowData.RequestBody) + CacheFileSplitTag)
	writeText.WriteString(fmt.Sprint(rowData.APIStatus) + CacheFileSplitTag)
	writeText.WriteString(fmt.Sprint(rowData.APIErrorMessage.String) + CacheFileSplitTag)
	writeText.WriteString(fmt.Sprint(rowData.ResponseHttpStatus) + CacheFileSplitTag)
	writeText.WriteString(fmt.Sprint(rowData.ResponseCosttime) + CacheFileSplitTag)
	writeText.WriteString(fmt.Sprint(rowData.ResponseAccount) + CacheFileSplitTag)
	writeText.WriteString(fmt.Sprint(rowData.ResponseMsgCode) + CacheFileSplitTag)
	writeText.WriteString(fmt.Sprint(rowData.ResponseMsg) + CacheFileSplitTag)
	writeText.WriteString(fmt.Sprint(rowData.ResponseDataDecrypt.String) + CacheFileSplitTag)
	writeText.WriteString(fmt.Sprint(rowData.ResponseDataEncrypt.String) + CacheFileSplitTag)
	//writeText.WriteString(fmt.Sprint(rowData.DebugMessage.String) + CacheFileSplitTag) //Debug資訊取消寫入
	writeText.WriteString(fmt.Sprint(rowData.LogTime) + CacheFileSplitTag) //最後寫入換行記號
	writeText.WriteString(CacheFileLineTag)                                //最後寫入換行記號

	_, err = thisFile.WriteString(writeText.String())
	if err != nil {
		return
	}

	_, err = thisFile.Seek(0, 0)
	return

}

//抹除暫存檔案內的內容
func TruncateFile(thisFile *os.File) (err error) {
	err = thisFile.Truncate(0)
	if err != nil {
		return
	}
	_, err = thisFile.Seek(0, 0)
	return

}

// 錯誤設計
// func OverwriteMessageToFile(thisFile *os.File, gormRowList []gormRow) (err error) {

// 	err = thisFile.Truncate(0)
// 	if err != nil {
// 		return
// 	}

// 	_, err = thisFile.Seek(0, 0)
// 	if err != nil {
// 		return
// 	}

// 	var writeText strings.Builder
// 	for index := 0; index < len(gormRowList); index++ {

// 		value := reflect.ValueOf(gormRowList[index])
// 		for index := 0; index < value.NumField(); index++ {
// 			if index == 17 {
// 				continue //DebugMessage取消不補紀錄
// 			}
// 			writeText.WriteString(fmt.Sprint(value.Field(index)) + CacheFileSplitTag)
// 		}
// 		writeText.WriteString(CacheFileLineTag)

// 	}

// 	_, err = thisFile.WriteString(writeText.String())

// 	return

// }
