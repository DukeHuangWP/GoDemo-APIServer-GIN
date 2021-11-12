package antiflood

import (
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func Test_AllowIP(testT *testing.T) {

	ginRouter := gin.Default()
	// ginRouter.Use(func(ginCTX *gin.Context) {
	// 	TriggerHandler(ginCTX) //反SYN-Flood
	// })

	// ginRouter.Any("/", func(ginCTX *gin.Context) {
	// }) //此寫法會讓ginRouter接使用該中間鍵

	ginGroup1 := ginRouter.Group("/") 

	ginGroup1.Use(func(ginCTX *gin.Context) {
		TriggerHandler(ginCTX) //反SYN-Flood
	})

	ginGroup1.Any("/", func(ginCTX *gin.Context) {
	})

	err := SetAllowIPRule([]string{"::*", "::1"})
	if err != nil || len(ClientIPAllowList)+len(ClientIPRuleList)+len(clientIPRuleLenList) <= 0 {
		log.Fatalf("請檢查添加IP白名單是有含有錯誤！(%v)\n", err)
	}

	StartTimer(&Cfgs{MaxCount: 2, RateSecs: 10})

	go func() {
		if err := ginRouter.Run(":8081"); err != nil {
			log.Fatal("請檢查localhost:8081是否正常開啟！")
		}
	}()

	for index := 0; index < 10; index++ {
		go func() {
			doPostAllowIPTest(testT)
		}()
	}

	timerTicker := time.NewTicker(time.Duration(2) * time.Second)
	defer timerTicker.Stop()

	<-timerTicker.C
	log.Println("--------測試結束-------")

	return

}

func Test_TryAntiSYNFlood(testT *testing.T) {

	ginRouter := gin.Default()
	// ginRouter.Use(func(ginCTX *gin.Context) {
	// 	TriggerHandler(ginCTX) //反SYN-Flood
	// })

	// ginRouter.Any("/", func(ginCTX *gin.Context) {
	// }) //此寫法會讓ginRouter接使用該中間鍵

	ginGroup1 := ginRouter.Group("/") 

	ginGroup1.Use(func(ginCTX *gin.Context) {
		TriggerHandler(ginCTX) //反SYN-Flood
	})

	ginGroup1.Any("/", func(ginCTX *gin.Context) {
	})

	err := SetAllowIPRule(nil)
	if err != nil {
		log.Fatalf("請檢查添加IP白名單是有含有錯誤！(%v)\n", err)
	}
	StartTimer(&Cfgs{RateSecs: 10})

	go func() {
		if err := ginRouter.Run(":8082"); err != nil {
			log.Fatal("請檢查localhost:8082是否正常開啟！")
		}
	}()

	for index := 0; index < 100; index++ {
		go func() {
			doPostFoodTest(testT)
		}()
	}

	timerTicker := time.NewTicker(time.Duration(5) * time.Second)
	defer timerTicker.Stop()

	<-timerTicker.C
	log.Println("\n--------測試結束-------")
	return

}

func doPostFoodTest(testT *testing.T) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:8082", nil)
	isAntiFoodWorking := false
	for index := uint64(0); index < cfgMaxCount+10; index++ {
		resp, err := client.Do(req)
		if err != nil {
			testT.Errorf("測試連線過程發生錯誤: %v\n", err)
		}
		if resp.StatusCode != http.StatusOK {
			testT.Errorf("Get response 值非 200: %v\n", resp.Status)
		}

		if index > cfgMaxCount {
			respBody, _ := ioutil.ReadAll(resp.Body)
			if string(respBody) != banPageHTML {
				testT.Error("Anti SYN-Flood didnt work!\n")
				return
			}
			isAntiFoodWorking = true
		}
		resp.Body.Close()
	}
	if isAntiFoodWorking == false {
		testT.Errorf("Anti SYN-Flood didnt work! (%v)\n", isAntiFoodWorking)
		return
	}
}

func doPostAllowIPTest(testT *testing.T) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:8081", nil)
	isAllowIPWorking := false
	for index := uint64(0); index < cfgMaxCount+10; index++ {
		resp, err := client.Do(req)
		if err != nil {
			testT.Errorf("測試連線過程發生錯誤: %v\n", err)
		}
		if resp.StatusCode != http.StatusOK {
			testT.Errorf("Get response 值非 200: %v\n", resp.Status)
		}

		if index > cfgMaxCount {
			respBody, _ := ioutil.ReadAll(resp.Body)
			if string(respBody) == banPageHTML {
				testT.Error("Anti SYN-Flood ban it! AllowIP not working\n")
				return
			}
			isAllowIPWorking = true
		}
		resp.Body.Close()
	}

	if isAllowIPWorking == false {
		testT.Errorf("AllowIP didnt work! (%v)\n", isAllowIPWorking)
		return
	}
}
