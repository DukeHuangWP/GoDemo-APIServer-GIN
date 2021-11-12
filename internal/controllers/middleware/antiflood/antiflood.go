package antiflood

import (
	"DevIntergTest/internal/common"
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	Service_Name = "Anti-Flood"
	banPageHTML  = `<!DOCTYPE html><html><body><img src="https://stickershop.line-scdn.net/stickershop/v1/sticker/34521489/IOS/sticker.png"></body></html>`
)

type Cfgs struct {
	ProxyHeaderKey    string //請求request Header 判斷Client IP之 key值, 例如 X-Forwarded-For: 42.42.42.42
	MaxCount          uint64 //短時間連線次數上限
	RateSecs          int64  //定義短時間秒數
	ClientIPUnbanSecs int64  //遭禁ClientIP等待解封秒數
	TimerUnbanSecs    int    //檢查解封ClientIP間格秒數
}

type clientIPType struct {
	ConnectCount uint64 //短時間連線次數
	LastTime     int64  //客戶端最後連線時間戳
}

//[clientIP]資訊
var ClientIPRecords = make(map[string]*clientIPType)
var rwMutex = sync.RWMutex{}

var (
	cfgProxyHeaderKey    string        //請求request Header 判斷Client IP之 key值
	cfgMaxCount          uint64 = 255  //短時間連線次數上限
	cfgRateSecs          int64  = 10   //定義短時間秒數
	cfgClientIPUnbanSecs int64  = 600  //遭禁ClientIP等待解封秒數
	cfgTimerUnbanSecs    int    = 1800 //檢查解封ClientIP間格秒數
)

//啟動antiflood定時器,自動封鎖短時間內過度使用api的ip,定時清除閒置clientIP
func StartTimer(cfgs *Cfgs) {
	if cfgs == nil {
		log.Printf("%v : 載入預設值設定 [短時間連線次數上限:%v]&[定義短時間秒數:%v]&[遭禁ClientIP自動解封時間:%v]&[自動解封時間執行周期:%v]", Service_Name, cfgMaxCount, cfgRateSecs, cfgClientIPUnbanSecs, cfgTimerUnbanSecs)
	} else {
		cfgProxyHeaderKey = cfgs.ProxyHeaderKey

		if cfgs.MaxCount > 1 {
			cfgMaxCount = cfgs.MaxCount
		} //不可以小於1

		if cfgs.RateSecs > 1 {
			cfgRateSecs = cfgs.RateSecs
		} //不可以小於1

		if cfgs.ClientIPUnbanSecs > 1 {
			cfgClientIPUnbanSecs = cfgs.ClientIPUnbanSecs
		} //不可以小於1

		if cfgs.TimerUnbanSecs > 1 {
			cfgTimerUnbanSecs = cfgs.TimerUnbanSecs
		} //不可以小於1
		log.Printf("%v : 載入預設值設定 [短時間連線次數上限:%v]&[定義短時間秒數:%v]&[遭禁ClientIP自動解封時間:%v]&[自動解封時間執行周期:%v]", Service_Name, cfgMaxCount, cfgRateSecs, cfgClientIPUnbanSecs, cfgTimerUnbanSecs)
	}

	TimerTicker := time.NewTicker(time.Duration(cfgTimerUnbanSecs) * time.Second)
	defer TimerTicker.Stop()
	go func(ticker *time.Ticker) {
		for {
			<-ticker.C
			clearClientIPRecords(ClientIPRecords)
		}
	}(TimerTicker)
}

//啟用反SYN-Flood機制
func clearClientIPRecords(ClientIPRecords map[string]*clientIPType) {

	rwMutex.Lock()         //寫入上鎖
	defer rwMutex.Unlock() //寫入解鎖
	if ClientIPRecords == nil {
		return
	}

	for index, value := range ClientIPRecords {
		if time.Now().Unix()-value.LastTime > cfgClientIPUnbanSecs {
			log.Printf("ClientIP 監控佔存清除 > %v (Count:'%v', LastTime:'%v')", index, value.ConnectCount, value.LastTime)
			delete(ClientIPRecords, index)
		}
	}
}

//觸發Handler(計算連線次數)
func TriggerHandler(ginCTX *gin.Context) {

	defer common.CatchPanic("antiflood panic")

	rwMutex.Lock()         //寫入上鎖
	defer rwMutex.Unlock() //寫入解鎖

	clientIP := common.GetClientPublicIP(ginCTX.Request.Header.Get(cfgProxyHeaderKey), ginCTX.ClientIP())
	//clientIP := ginCTX.ClientIP()
	if IsIPInAllowList(clientIP) == true {
		return
	} //白名單

	if _, ok := ClientIPRecords[clientIP]; !ok {
		ClientIPRecords[clientIP] = &clientIPType{ConnectCount: 1, LastTime: time.Now().Unix()}
	} else {

		if time.Now().Unix()-ClientIPRecords[clientIP].LastTime < cfgRateSecs {
			ClientIPRecords[clientIP].ConnectCount++
		} else if ClientIPRecords[clientIP].ConnectCount > 0 {
			ClientIPRecords[clientIP].ConnectCount--
		}
		ClientIPRecords[clientIP].LastTime = time.Now().Unix()

		if ClientIPRecords[clientIP].ConnectCount > cfgMaxCount {

			ginCTX.Writer.WriteString(banPageHTML)
			//ginCTX.Writer.WriteHeader(http.StatusGatewayTimeout) //回復timeout
			ginCTX.Abort()
			return
		}
	}
	//fmt.Println(ClientIPRecords[clientIP])
}
