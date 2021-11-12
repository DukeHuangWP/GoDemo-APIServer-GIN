package models

import (
	"log"
	"sync"
	"time"
)

//const dbCacheTimerSec = 300 //db儲存快取重新讀取秒數

var (
	dbCache   map[string]interface{} //db儲存快取
	dbRWMutex = &sync.RWMutex{}      //db快取專用鎖
)

func init() {
	dbCache = make(map[string]interface{})
}

//獲取db快取變數值
func GetDBCacheValue(key string) interface{} {
	dbRWMutex.RLock()
	defer dbRWMutex.RUnlock() //讀取解鎖
	if value, isExsit := dbCache[key]; isExsit {
		return value
	}
	return nil
}

//寫入db快取變數值
func SetDBCacheValue(key string, value interface{}) {
	dbRWMutex.Lock()
	defer dbRWMutex.Unlock() //讀寫入解鎖
	dbCache[key] = value
}

//啟動定時器定期更新db快取
func StartBuildDBCache(timerSec uint32) {
	timer := time.NewTicker(time.Duration(int64(timerSec)) * time.Second)
	for {
		dbCache = make(map[string]interface{}) //刪除舊map
		saveDBCache()
		<-timer.C
	}
}

//建立db快取變數值
func saveDBCache() {

	//執行許要建立db快取之models
	DoSQLs.GetServerUID()
	///////////////////////////
	log.Printf("%v : 獲取DB快取變數完成 > %v\n", ServiceName, dbCache)
	return
}
