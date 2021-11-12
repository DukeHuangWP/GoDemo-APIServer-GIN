package alert

import (
	"DevIntergTest/internal/common"
	"DevIntergTest/internal/global"
	"DevIntergTest/internal/models"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
)

const (
	ServiceName = "The-Alert"
)

var IsOnline bool //alert是否啟用

func init() {

	if global.TelegramSendURL != "" {
		client := &http.Client{}
		request, err := http.NewRequest("GET", global.TelegramSendURL, &strings.Reader{})
		if err != nil {
			log.Printf("%v : 發送telegream時發生錯誤(http) > %v", ServiceName, err)
			return
		}
		response, err := client.Do(request)
		if err != nil {
			log.Printf("%v : 發送telegream時發生錯誤(client) > %v", ServiceName, err)
			return
		}
		defer response.Body.Close()

		_, err = ioutil.ReadAll(response.Body)
		if err != nil {
			log.Printf("%v : 發送telegream時發生錯誤(ioutil) > %v", ServiceName, err)
			return
		}

		if response.StatusCode == http.StatusBadRequest {
			IsOnline = true //回傳400表示tg服務器正常
		}
	}

}

func SendMessageWithTime(message string) {
	SendMessageCustom(fmt.Sprintf("("+common.GetTimeNowFormat()+")"+global.TelegramTitle_Format+" : %v", models.DoSQLs.GetServerUID(), message))
}

func SendMessageCustom(message string) {
	if IsOnline == false || message == "" {
		log.Printf("%v : TG並未開啟或收到空訊息", ServiceName)
	} else {

		payload := &bytes.Buffer{}
		writer := multipart.NewWriter(payload)
		_ = writer.WriteField("chat_id", "-392839279")
		if len(message) > global.TelegramSendMessageLimit { //擷取長度限制
			message = message[:global.TelegramSendMessageLimit]
		}
		_ = writer.WriteField("text", message)
		err := writer.Close()
		if err != nil {
			log.Printf("%v : 發送telegream時發生錯誤(WriteField) > %v", ServiceName, err)
			return
		}

		client := &http.Client{}
		req, err := http.NewRequest("POST", global.TelegramSendURL, payload)

		if err != nil {
			log.Printf("%v : 發送telegream時發生錯誤(http) > %v", ServiceName, err)
			return
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())
		res, err := client.Do(req)
		if err != nil {
			log.Printf("%v : 發送telegream時發生錯誤(client) > %v", ServiceName, err)
			return
		}
		defer res.Body.Close()

		// body, err := ioutil.ReadAll(res.Body)
		// if err != nil {
		//   fmt.Println(err)
		//   return
		// }
		// fmt.Println(string(body))
	}
}
