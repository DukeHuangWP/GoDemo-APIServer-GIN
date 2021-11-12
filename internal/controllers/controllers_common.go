package controllers

import (
	"DevIntergTest/internal/alert"
	"DevIntergTest/internal/logger/loggerToDB"
	"fmt"
	"net/http"
	"runtime"
	"strconv"

	"github.com/gin-gonic/gin"
)

//Panic捕捉,並且喜入紀錄
func finalNodePanicLog(ginCTX *gin.Context, loggerMsg *loggerToDB.Message) {
	if err := recover(); err != nil {
		// var logText string
		// for index := 1; index < 11; index++ { //最多捕捉10層log
		// 	ptr, filename, line, ok := runtime.Caller(index)
		// 	if !ok {
		// 		break
		// 	}
		// 	logText = logText + fmt.Sprintf(" %v:%d,%v > %v\n", filename, line, runtime.FuncForPC(ptr).Name(), err)
		// }

		// 可以捕捉函式名稱及行數
		// /Users/duke/+WorkDir/Dev-API/internal/controllers/controllers_for_task.go:211 (DevIntergTest/internal/controllers.EndTaskInquireTask.func1)
		// /Users/duke/+WorkDir/Dev-API/internal/controllers/controllers_for_task.go:361 (DevIntergTest/internal/controllers.EndTaskInquireTask)
		// /Users/duke/+WorkDir/Dev-API/pkg/src/github.com/gin-gonic/gin@v1.7.2/context.go:165 (github.com/gin-gonic/gin.(*Context).Next)
		// /Users/duke/+WorkDir/Dev-API/internal/controllers/controllers_for_task.go:167 (DevIntergTest/internal/controllers.MidSignOnly)
		// /Users/duke/+WorkDir/Dev-API/pkg/src/github.com/gin-gonic/gin@v1.7.2/context.go:165 (github.com/gin-gonic/gin.(*Context).Next)
		// /Users/duke/+WorkDir/Dev-API/pkg/src/github.com/gin-gonic/gin@v1.7.2/recovery.go:99 (github.com/gin-gonic/gin.CustomRecoveryWithWriter.func1)

		buf := make([]byte, 8192)
		buf = buf[:runtime.Stack(buf, true)]
		logText := string(buf)
		//可以捕捉紀錄非常多的訊息 路例如
		//goroutine 28 [running]:
		// DevIntergTest/internal/logger/loggerToDB.(*Message).TransToGorm(0xc000510300, 0xc05d70bf9a)
		// /Users/duke/+WorkDir/Dev-API/internal/controllers/loggerToDB/loggerToDB.go:94 +0x5fb
		// DevIntergTest/internal/logger/loggerToDB.AddMessage(0xc000510300, 0x108bba6, 0x60eee0e2)
		// /Users/duke/+WorkDir/Dev-API/internal/controllers/loggerToDB/loggerToDB.go:208 +0x5b

		loggerMsg.APIStatus = -2
		loggerMsg.APIErrorMessage = fmt.Sprintf("發生嚴重錯誤 >> \n%v", logText)
		loggerToDB.AddMessage(loggerMsg) //寫入logger
		alert.SendMessageWithTime(logText)

		ginCTX.AbortWithStatus(http.StatusInternalServerError) //後面程式碼將不會被執行
		return
	}

}

//檢查是否為合法外網http網址, 0.0開頭,127開頭,172開頭,192開頭,10開頭,::開頭,[::開頭 都不行
func isValidInternetHTTPURL(url string) bool {
	if url == "" {
		return false
	} else if len(url) > 8 {

		var urn string //網域與路徑
		var isSchemeValid bool
		if url[:7] == "http://" {
			isSchemeValid = true
			urn = url[7:]
		}

		if url[:8] == "https://" {
			isSchemeValid = true
			urn = url[8:]
		}

		if isSchemeValid == false {
			return false
		}

		if len(urn) >= 9 { //localhost
			if urn[:9] == "localhost" { //localhost.com也不行....
				return false
			}
		}

		if len(urn) >= 3 { //IPs
			urnTitle := urn[:3]
			if urnTitle == "10." {
				return false
			}

			if urnTitle == "0.0" {
				return false
			}

			if urnTitle == "::1" {
				return false
			}

			if urnTitle == "[::" {
				return false
			}
		}

		if len(urn) >= 4 { //IPs
			urnTitle := urn[:4]
			if urnTitle == "127." {
				return false
			}

			if urnTitle == "172." {
				return false
			}

			if urnTitle == "192." {
				return false
			}

			if urnTitle == "172." {
				return false
			}

			if urnTitle == "192." {
				return false
			}

		}

		return true

	}

	return false
}

func ConvertInterfaceToFloat64(input interface{}) (output float64, isFloat64 bool) {

	switch value := input.(type) {

	case float32:
		return float64(value), true
	case float64:
		return value, true
	case string:
		output, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return 0.0, false
		}
		return output, true
	case int:
		return float64(value), true
	case int16:
		return float64(value), true
	case int32:
		return float64(value), true
	case int64:
		return float64(value), true
	case int8:
		return float64(value), true
	case uint:
		return float64(value), true
	case uint16:
		return float64(value), true
	case uint32:
		return float64(value), true
	case uint64:
		return float64(value), true
	case uint8:
		return float64(value), true
	default:
		return 0.0, false
	}
}
