package httpSender

import (
	"DevIntergTest/internal/common"
	"io/ioutil"
	"net/http"
	"strings"
)

func PostJson(sendURL, requestBody string) (responseStatus int, responseBody string, err error) {

	client := &http.Client{}
	request, err := http.NewRequest("POST", sendURL, strings.NewReader(requestBody))
	if err != nil {
		return -1, "", err
	}
	request.Header.Add("Content-Type", "application/json")
	request.Close = true //避免EOF錯誤
	response, err := client.Do(request)
	if err != nil {
		return -1, "", err
	}
	defer response.Body.Close()

	responseBodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return -1, "", err
	}

	return response.StatusCode, common.BytesToString(responseBodyBytes), nil

}
