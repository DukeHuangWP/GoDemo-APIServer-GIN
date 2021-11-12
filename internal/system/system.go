package system

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
)

//設定CPU核心數 0=開啟所有核心,4=限制使用四顆核心,-1=所有核心數減一
func SetCpuCore(setNum int) {
	cpuNum := runtime.NumCPU()

	if setNum > 0 {

		if setNum > cpuNum { //超過實體核心數者，改用最大實體核心數
			setNum = cpuNum
		}
		runtime.GOMAXPROCS(setNum)

	} else if setNum < 0 {

		if sumNum := cpuNum + setNum; sumNum > 0 {
			setNum = sumNum
		} else { //低過核心數0者，改用單核心數
			setNum = 1
		}
		runtime.GOMAXPROCS(setNum)

	} else { //setNum=0不做任何設置
		setNum = runtime.NumCPU()
	}

	log.Printf("設定使用CPU核心數為 : %v (核心總數:%v)\n", setNum, cpuNum)
}

//獲取對外ip
func GetOutboundIP(checkerUrl string) (outboundIP string) {

	fmt.Println(checkerUrl)
	request, err := http.NewRequest("GET", checkerUrl, nil)

	if err != nil {
		//fmt.Println(err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		//fmt.Println(err)
		return
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		//fmt.Println(err)
		return
	}

	return string(responseBody)
}
