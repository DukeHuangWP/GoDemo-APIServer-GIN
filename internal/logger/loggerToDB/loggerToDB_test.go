package loggerToDB

import (
	"log"
	"testing"
)

func Test_count(testT *testing.T) {

	log.Println(countLoggerTableNames("aaaa"))

}

// func Test_TryLoogerToDB(testT *testing.T) {

// 	ctx, cancal := context.WithCancel(context.Background())
// 	defer cancal()

// 	saveLastLogFilePath := "./.UniTest.txt"
// 	saveLastLogText := `⸽模擬殘留資料1⸽⸽⸽⸽X1⸽X1⸽⸽0⸽X1⸽0⸽0⸽⸽⸽⸽X1⸽X1⸽0⸽⸽模擬殘留資料2⸽⸽⸽⸽X2⸽X2⸽⸽0⸽X2⸽0⸽0⸽⸽⸽⸽X2⸽X2⸽0⸽`
// 	err := ioutil.WriteFile(saveLastLogFilePath, []byte(saveLastLogText), 0644)
// 	if err != nil {
// 		testT.Errorf("寫入模擬殘留log資料失敗 > '%v'\n", saveLastLogFilePath)
// 	}

// 	InitLogger(ctx, &MySQLConfig{Account: "iftadmin", Password: "qweRTY123", DBHost: "35.241.104.139:3306", DBName: "system_server"}, 5, saveLastLogFilePath)
// 	StartKeepDBAlive(nil, 2)
// 	NowTableName = "api_log"

// 	log.Println("開始模擬DB連線與計時器寫入測試...,測試過程不影響DB,起始Goroutine數量:", runtime.NumGoroutine())
// 	go func() {
// 		for {
// 			//log.Println("idle.... ", runtime.NumGoroutine())
// 			time.Sleep(time.Second)
// 		}
// 	}()

// 	go func() {
// 		index := 0
// 		for {

// 			AddMessage(&Message{ClientIP: "", ServerHost: "x1"})
// 			time.Sleep(time.Second)
// 			index++
// 			if index > 50 {
// 				break
// 			}
// 		}
// 	}()

// 	go func() {
// 		index := 0
// 		for {

// 			AddMessage(&Message{ClientIP: "", ServerHost: "x2"})
// 			time.Sleep(time.Second)
// 			index++
// 			if index > 50 {
// 				break
// 			}
// 		}
// 	}()

// 	go func() {
// 		index := 0
// 		for {

// 			AddMessage(&Message{ClientIP: "", ServerHost: "x3"})
// 			time.Sleep(time.Second)
// 			index++
// 			if index > 50 {
// 				break
// 			}
// 		}
// 	}()

// 	go func() {
// 		index := 0
// 		for {

// 			AddMessage(&Message{ClientIP: "", ServerHost: "x4"})
// 			time.Sleep(time.Second)
// 			index++
// 			if index > 50 {
// 				break
// 			}
// 		}
// 	}()

// 	go func() {
// 		index := 0
// 		for {

// 			AddMessage(&Message{ClientIP: "", ServerHost: "x5"})
// 			time.Sleep(time.Second)
// 			index++
// 			if index > 50 {
// 				break
// 			}
// 		}
// 	}()

// 	TimerTicker := time.NewTicker(time.Duration(20) * time.Second)
// 	defer TimerTicker.Stop()
// 	for {
// 		<-TimerTicker.C
// 		err := os.Remove(saveLastLogFilePath)
// 		if err != nil {
// 			testT.Errorf("刪除模擬殘留log資料失敗 > '%v'\n", saveLastLogFilePath)
// 		}
// 		return
// 	}

// }
