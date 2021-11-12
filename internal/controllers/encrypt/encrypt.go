package encrypt

import (
	"DevIntergTest/internal/common"
	"bytes"
	cryptography "crypto/des"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"strings"
)

//修正cryptography.GetSignDataString bug
func fixGetSignDataString(jsonStruct interface{}) (jsonSignData string, IsSignAdded bool, err error) {

	jsonByte, err := json.Marshal(jsonStruct)
	if err != nil {
		return "", false, err
	}
	jsonSignData = common.BytesToString(jsonByte)
	jsonSignData = strings.ReplaceAll(jsonSignData, `\u0026`, "&")
	jsonSignData = strings.ReplaceAll(jsonSignData, `\u003c`, "<")
	jsonSignData = strings.ReplaceAll(jsonSignData, `\u003e`, ">")
	//json编码时，字符含有  &, <, and >的会被编码成 \u0026, \u003c, and \u003e

	if keyStr := `"Key":"`; strings.Contains(jsonSignData, keyStr) == true {
		jsonStringHalf := jsonSignData[strings.Index(jsonSignData, keyStr):]
		keyJsonString := jsonStringHalf[:strings.LastIndex(jsonStringHalf, `"`)+1]
		jsonSignData = strings.Replace(jsonSignData, keyJsonString, `"Sign":"`+cryptography.MD5(jsonSignData)+`"`, -1)
		IsSignAdded = true
	}

	return

}

//將純文字以3DES UTF8編碼加密 (ECB Mode, PADDING:ZEROS)
func GetEncryptData(inputData string, desKey string) (outputData string, IsSignAdded bool, err error) {
	defer func(err error) {
		if errRecover := recover(); errRecover != nil {
			err = fmt.Errorf("觸發3DES加密Bug")
			return
		}
	}(err)

	if inputData == "" || desKey == "" {
		return "", IsSignAdded, fmt.Errorf("輸入資料不可為空")
	}

	jsonBuffer := new(bytes.Buffer)
	if err = json.Compact(jsonBuffer, common.StringToBytes(inputData)); err == nil {
		inputData = jsonBuffer.String()
		inputData = strings.ReplaceAll(inputData, `\u0026`, "&")
		inputData = strings.ReplaceAll(inputData, `\u003c`, "<")
		inputData = strings.ReplaceAll(inputData, `\u003e`, ">")
		//json编码时，字符含有  &, <, and >的会被编码成 \u0026, \u003c, and \u003e
	}

	if keyStr := `"Key":"`; strings.Contains(inputData, keyStr) == true {
		jsonStringHalf := inputData[strings.Index(inputData, keyStr):]
		keyJsonString := jsonStringHalf[:strings.LastIndex(jsonStringHalf, `"`)+1]
		inputData = strings.Replace(inputData, keyJsonString, `"Sign":"`+cryptography.MD5(inputData)+`"`, -1)
		IsSignAdded = true
	}

	outputData, err = cryptography.TripleEcbDesEncrypt(inputData, desKey, "ZEROS")

	return outputData, IsSignAdded, err
}

//將純文字以3DES UTF8編碼解密 (ECB Mode, PADDING:ZEROS)
func GetDecryptData(inputData string, desKey string) (outputData string, err error) {
	defer func(err error) {
		if errRecover := recover(); errRecover != nil {
			err = fmt.Errorf("觸發3DES加密Bug")
			return
		}
	}(err)

	if inputData == "" || desKey == "" {
		return "", fmt.Errorf("輸入資料不可為空")
	}

	outputData, err = cryptography.TripleEcbDesDecrypt(inputData, desKey, "ZEROS")
	//fmt.Printf("%v(%v)\n", outputData, err)
	return outputData, err
}

//對Response之struct結構轉換成加密文本
func GetEncryptResponseData(jsonStruct interface{}, desKey string) (responseData string, err error) {
	//jsonString := cryptography.GetSignDataString(jsonStruct)
	jsonString, _, err := fixGetSignDataString(jsonStruct)
	if err != nil {
		return "", err
	}

	responseData, err = cryptography.TripleEcbDesEncrypt(jsonString, desKey, "ZEROS")
	return
}

//對Response之struct結構轉換成加密文本,並回傳解密文本
func GetEncryptResponseDataAndJson(jsonStruct interface{}, desKey string) (responseData, jsonData string, err error) {
	//jsonString := cryptography.GetSignDataString(jsonStruct)
	jsonString, _, err := fixGetSignDataString(jsonStruct)
	if err != nil {
		return "", "", err
	}

	responseData, err = cryptography.TripleEcbDesEncrypt(jsonString, desKey, "ZEROS")
	if err == nil {
		var jsonByte bytes.Buffer
		enc := json.NewEncoder(&jsonByte)
		enc.SetEscapeHTML(false) //json编码时，字符含有  &, <, and >的会被编码成 \u0026, \u003c, and \u003e
		enc.Encode(jsonStruct)
		jsonData = jsonByte.String()

	}
	return
}

//確認Request解密文本中json的3DES加密簽名是否符合
func CheckRequestEncryptSign(jsonData string, signKey string) (isMatch bool) {
	// fmt.Println(jsonData)
	// println()
	keySign := `"Sign":"`
	if leftIndex := strings.Index(jsonData, keySign); leftIndex > 0 {
		keySignIndex := leftIndex + len(keySign)
		if rightIndex := strings.Index(jsonData[keySignIndex:], `"`); rightIndex > 0 {
			signJson := jsonData[leftIndex : len(jsonData)-1]
			signValue := jsonData[keySignIndex : keySignIndex+rightIndex]
			checksumSignValue := fmt.Sprintf("%x", md5.Sum([]byte(strings.Replace(jsonData, signJson, `"Key":"`+signKey+`"`, 1)))) //替換Sign成Key後取md5值即為驗正值
			// println()
			// fmt.Println("該帳號資料庫之 Sign:", signKey)
			// fmt.Println("Request簽名    Sign:", signValue)
			// fmt.Println("Request資料庫之Sign:", checksumSignValue)
			// println()
			if signValue == checksumSignValue { //計算出Sign值必須與db中該帳號Sign值相同
				//fmt.Println("正確")
				isMatch = true
			}
		}
	}
	return
}
