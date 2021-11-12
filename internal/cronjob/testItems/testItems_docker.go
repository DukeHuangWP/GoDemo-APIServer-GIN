package testItems

import (
	"DevIntergTest/internal/alert"
	"DevIntergTest/internal/common"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

const (
	TaskType_DockerWatcher = "DockerWatcher"
	DockerKeyWords_Key     = "Keywords"
	DockerSinces_Key       = "Since"
)

type DockerWatcher struct {
	TestName        string
	IsValid         bool                   //手否經過驗證
	Dockers         map[string]interface{} `json:"Dockers"`         //API發起帳號 表account.login_account
	AlertLinesUpper uint8                  `json:"AlertLinesUpper"` //收款閒置多久後開始出款
	AlertLinesLower uint8                  `json:"AlertLinesLower"` //收款閒置多久後開始出款
	DockerKeyWords  map[string][]string    //[容器名稱]檢索關鍵字
	DockerSinces    map[string]string      //[容器名稱]檢索關鍵字
}

func (scriptData *DockerWatcher) ValidScript(testName, testTag, testScript string) (err error) {

	if testName == "" {
		return fmt.Errorf("TestName不可為空")
	}

	err = json.Unmarshal([]byte(testScript), scriptData)
	if err != nil {
		return err
	}

	if scriptData.AlertLinesUpper+scriptData.AlertLinesUpper == 0 {
		return fmt.Errorf("AlertLines_Upper+AlertLines_Lower 不可等於0")
	}

	if len(scriptData.Dockers) <= 0 {
		return fmt.Errorf("Dockers內無關鍵字無任何設定!")
	}

	scriptData.TestName = testName
	scriptData.DockerKeyWords = make(map[string][]string)
	scriptData.DockerSinces = make(map[string]string)

	for dockerName, value := range scriptData.Dockers {
		if settings, ok := value.(map[string]interface{}); ok == true {
			scriptData.DockerSinces[dockerName] = fmt.Sprint(settings[DockerSinces_Key])
			if keyWords, ok := settings[DockerKeyWords_Key].([]interface{}); ok == true {
				for _, aKeyword := range keyWords {
					scriptData.DockerKeyWords[dockerName] = append(scriptData.DockerKeyWords[dockerName], fmt.Sprint(aKeyword))
				}
			}
		}
	}

	if len(scriptData.DockerKeyWords) <= 0 {
		return fmt.Errorf("'%v' 內無關鍵字無任何設定!", DockerKeyWords_Key)
	}

	if len(scriptData.DockerSinces) <= 0 {
		return fmt.Errorf("'%v' ƒ內無關鍵字無任何設定!", DockerSinces_Key)
	}

	scriptData.IsValid = true
	return nil
}

func (scriptData *DockerWatcher) GetScriptFunc() (scriptFunc func(), err error) {
	if scriptData.IsValid == false {
		return nil, fmt.Errorf("尚未ValidScript")
	}

	//fmt.Println(scriptData.DockerKeyWords)

	scriptFunc = func() {
		var systemCMD *exec.Cmd
		var output []byte
		systemCMD = exec.Command("docker", "ps")
		output, _ = systemCMD.CombinedOutput()
		dockerRunningList := common.BytesToString(output)
		// CONTAINER ID   IMAGE                              COMMAND                  CREATED        STATUS                 PORTS                                                                              NAMES
		// 9957657388f6   gcr.io/cadvisor/cadvisor:v0.39.2   "/usr/bin/cadvisor -…"   2 months ago   Up 2 weeks (healthy)   0.0.0.0:3000->8080/tcp, :::3000->8080/tcp                                          cadvisor
		// 937b967f090a   prom/node-exporter:v1.1.2          "/bin/node_exporter …"   2 months ago   Up 2 weeks             0.0.0.0:9100->9100/tcp, :::9100->9100/tcp
		//fmt.Println(dockerRunningList)
		textLines := strings.Split(dockerRunningList, "\n")
		if len(textLines) < 2 {
			message := fmt.Sprintf("%v : '%v' 無法監控系統上的docker,請檢查權限!\n", SERVICE_NAME, scriptData.TestName)
			log.Printf(message)
			return
		} else if strings.Contains(textLines[0], "CONTAINER ID") == false && strings.Contains(textLines[0], "NAMES") == false {
			message := fmt.Sprintf("%v : '%v' 無法監控系統上的docker,請檢查權限!\n", SERVICE_NAME, scriptData.TestName)
			log.Printf(message)
			return
		}

		DockerList := make(map[string]string)
		for index := 1; index < len(textLines); index++ {
			for dockerName := range scriptData.DockerKeyWords {
				if strings.Index(textLines[index], dockerName) == len(textLines[index])-len(dockerName) { //找到末尾docker容器名稱
					dockerID := strings.ReplaceAll(textLines[index][:strings.Index(textLines[index], " ")], " ", "")
					DockerList[dockerName] = dockerID
				}
			}
		}

		if len(DockerList) < 0 {
			message := fmt.Sprintf("%v : '%v' 並未發現有任何Docker容器啟動!\n", SERVICE_NAME, scriptData.TestName)
			log.Printf(message)
			return
		}

		// message := fmt.Sprintf("%v : '%v'\n", SERVICE_NAME, scriptData.TestName)
		// alert.SendMessageWithTime(message)
		//fmt.Println(DockerList)

		for dockerName, dockerID := range DockerList {
			//fmt.Println(dockerName, dockerID)
			systemCMD = exec.Command("docker", "logs", dockerID, "--since", scriptData.DockerSinces[dockerName])
			output, _ = systemCMD.CombinedOutput()
			textLines := strings.SplitAfter(common.BytesToString(output), "\n")
			keyWords := scriptData.DockerKeyWords[dockerName]
			var alertMessageBuilder strings.Builder
			textLineNum := len(textLines)
			for index := 0; index < textLineNum; index++ {
				lineIndex := index
				for index := 0; index < len(keyWords); index++ {
					if strings.Contains(textLines[lineIndex], keyWords[index]) == true {

						if lineIndex > 0 {
							count := int(scriptData.AlertLinesUpper)
							if lineIndex-count < 0 {
								count += lineIndex - count
							}
							for index := count; index > 0; index-- {
								alertMessageBuilder.WriteString(textLines[lineIndex-index])
							}
						} //顯示上行文字

						alertMessageBuilder.WriteString(textLines[lineIndex]) //顯示本行文字

						if lineIndex+1 < textLineNum {
							count := int(scriptData.AlertLinesLower)
							if textLineNum-lineIndex-1 < count {
								count = textLineNum - lineIndex - 1
							}

							for index := 1; index <= count; index++ {
								alertMessageBuilder.WriteString(textLines[lineIndex+index])
							}
						} //顯示下行文字

						message := fmt.Sprintf("%v : '%v'發現[%v]內包含關鍵字:'%v' >\n%v", SERVICE_NAME, scriptData.TestName, dockerName, keyWords[index], alertMessageBuilder.String())
						log.Printf(message)
						alert.SendMessageWithTime(message)

					}
				}
			}

		}

	}

	//log.Println(scriptData)
	return scriptFunc, nil
}
