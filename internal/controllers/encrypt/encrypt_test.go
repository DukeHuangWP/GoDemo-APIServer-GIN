package encrypt

import (
	"strings"
	"testing"
)

func Test_TripleEcb(t *testing.T) {

	desKey := "TheLengthOfKeyOf3desIs24"
	inputText := `{
"Account":"admintest",
"Key":"202CB962AC59075B964B07152D234B70"
}`
	encryptText, _, err := GetEncryptData(inputText, desKey)
	if err != nil {
		t.Errorf("非預料輸出結果: 加密('%v')過程出現F錯誤提示: %v\n", inputText, err)
	}

	decryptText, err := GetDecryptData(encryptText, desKey)
	if err != nil {
		t.Errorf("非預料輸出結果: 解密('%v')過程出現F錯誤提示: %v\n", encryptText, err)
	}

	if decryptText != `{"Account":"admintest","Sign":"30e7a2c9c6a78fe8df84b6edeaa6ecdf"}` || inputText == "" || encryptText == "" {
		t.Errorf("非預料輸出結果: '%v' != '%v' , 錯誤提示: %v\n", encryptText, decryptText, err)
	} else {
		t.Logf("輸出加解密結果(正常): %v (%T)\n", encryptText, encryptText)
	}

	jsonStruct := struct {
		TestKey1 int    `json:"test_key1"`
		TestKey2 string `json:"test_key2"`
		Key      string `json:"Key"` //簽名專用
	}{
		TestKey1: 1,
		TestKey2: "test?&<>Value",
		Key:      "60EDFF7FCD537F17398F22173FD2EB18",
	}

	responseData1, _ := GetEncryptResponseData(jsonStruct, desKey)
	responseData2, jsonData, _ := GetEncryptResponseDataAndJson(jsonStruct, desKey)
	if responseData1 != responseData2 || responseData1 == "" {
		t.Errorf("非預料輸出結果: '%v' != '%v' , 錯誤提示: %v\n", responseData1, responseData2, err)
	} else {
		t.Logf("輸出加密前結果(正常): %v (%T)\n", jsonData, jsonData)
		t.Logf("輸出加密結果(正常): %v (%T)\n", responseData2, responseData2)
	}

	jsonSignStruct := struct {
		TaskUID            string `json:"TaskUID"` //客戶訂單編號
		TaskNo             string `json:"TaskNo"`  //工作任務代號
		CollectionCardInfo struct {
			BankName    string `json:"BankName"`    //銀行名稱
			BankCode    string `json:"BankCode"`    //銀行代碼
			AccountNo   string `json:"AccountNo"`   //帳號
			AccountName string `json:"AccountName"` //戶名
		}
		PayURL string `json:"PayURL"` //任務類型 2:轉帳 3:收款
		Key    string `json:"Key"`    //Response.Data必須加入簽名
	}{
		TaskUID: "DL1630388624209",
		TaskNo:  "COLL_ift00017_1630388627346",
		CollectionCardInfo: struct {
			BankName    string `json:"BankName"`    //銀行名稱
			BankCode    string `json:"BankCode"`    //銀行代碼
			AccountNo   string `json:"AccountNo"`   //帳號
			AccountName string `json:"AccountName"` //戶名
		}{
			BankName:    "招商银行",
			BankCode:    "CMB",
			AccountNo:   "6214830016359647",
			AccountName: "叶泰佑",
		},
		PayURL: "http://106.14.19.111:81/CheckOutCounterSucces.php?TaskUid=DL1630388624209&TaskNo=COLL_ift00017_1630388627346&BankCode=CMB&BankBranch=123&AccountNo=6214830016359647&AccountName=叶泰佑&RemitterAccountName=辛建芳&Amount=3&CreateTime=2021/08/31 13:58:47",
		Key:    "60EDFF7FCD537F17398F22173FD2EB18",
	}

	responseData3, _, _ := GetEncryptResponseDataAndJson(jsonSignStruct, desKey)
	decryptText, _ = GetDecryptData(responseData3, desKey)
	if strings.Contains(decryptText, `"Sign":"ad4bb7312705987c0c15f73ba13569d2"`) == false || decryptText == "" {
		t.Errorf("非預料輸出結果: '%v' != '%v' , Sign md5值不符ad4bb7312705987c0c15f73ba13569d2\n", decryptText, decryptText)
	} else {
		//t.Logf("輸出解密結果(Sign md5值正常): %v (%T)\n", decryptText, decryptText)
	}

	request := `{"TaskUid":"11","TaskNo":"22","TaskStage":"33","TaskType":"44","Content01":{"PayeeBankCode":"123123456489455","PayeeAccountNo":"ift00007","PayeeAccountName":"工商银行","Amount":"10"},"Content02":"66","Content03":"77","Content04":"88","Sign":"b959958ad26cbca80db1d3981030b578"}`
	signKey := "402CB962AC59075B964B07152D234B70"
	if CheckRequestEncryptSign(request, signKey) == false {
		t.Errorf("非預料輸出結果: request='%v' ,signKey='%v' md5值不相符\n", request, signKey)
	}
}
