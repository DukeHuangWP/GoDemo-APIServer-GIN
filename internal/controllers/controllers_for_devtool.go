package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	URN_DevToolRootPath = "/DevTool"
	URN_DevToolWebPath  = "/"

	URNPath_APITest = "/api_test"

	FolderPath_Images    = "Images"    //網頁公用目錄名稱 - 圖片
	FolderPath_Pages     = "Pages"     //網頁公用目錄名稱 - 純html頁面
	FolderPath_Templs    = "Templates" //網頁公用目錄名稱 - 版模
	FolderPath_Downloads = "Downloads" //網頁公用目錄名稱 - 下載目錄
	FolderPath_Uploads   = "Uploads"   //網頁公用目錄名稱 - 上傳目錄
	FolderPath_Others    = "Others"    //網頁公用目錄名稱 - 其他檔案分類

	Template_IndexMenu = "index-Menu.html"
	Template_APITest   = "System-APITest.html"

	// Template_PostTest     = "Post.html"
	// Template_RunSh        = "RunSh.html"
	// Template_DownloadList = "Download-List.html"
	// Template_UploadMenu   = "Upload-Menu.html"

	H5Page_Menu        = "HTMLPage_Menu"
	H5Value_Title      = "HTMLValue_Title"
	H5Value_LastUpdate = "HTMLValue_LastUpdate"
	H5Value_Return     = "HTMLValue_Return"

	H5Value_APITest_InputURL       = "HTMLInput_URL"
	H5Value_APITest_InputURN       = "HTMLInput_URN"
	H5Value_APITest_LoginAccount   = "HTMLInput_LoginAccount"
	H5Value_APITest_DesKey         = "HTMLInput_DesKey"
	H5Value_APITest_SignKey        = "HTMLInput_SignKey"
	H5Value_APITest_PostmanBody    = "HTMLInput_PostmanBody"
	H5Value_APITest_PostWebPath    = "HTMLValue_PostWebPath"
	H5Value_APITest_Response       = "HTMLValue_Response"
	H5Value_APITest_DecodeResponse = "HTMLValue_Decode_Response"
	H5Value_APITest_ReviewPayload  = "HTMLValue_ReviewPayload"
	H5Value_APITest_ReviewHTTP     = "HTMLValue_ReviewHTTP"
	H5Value_APITest_ReviewCurl     = "HTMLValue_ReviewCurl"

	APITest_ContentHeader = "application/json"

	PayloadFormat = `{
		"Account": "%v",
		"Data": "%v"
	}`

	CurlScriptFormat = `curl --location --request POST '%v' \
	--header 'Content-Type: %v' \
	--data-raw '{
		"Account": "%v",
		"Data": "%v"
	}'`

	HTTPScriptFormat = `POST %v HTTP/1.1
	Host: %v
	Content-Type: %v
	Content-Length: %v
	{
		"Account": "%v",
		"Data": "%v"
	}`
)

var (
	LastOnLineTime string //伺服器啟動時間

)


/*
URN:/
Response: 307
說明: 將根目錄導向首頁選單
*/
func HandleIndexMenu(ginCTX *gin.Context) {

	ginCTX.Redirect(http.StatusTemporaryRedirect, URN_DevToolRootPath[1:])
	ginCTX.Abort() //首頁

}

/*
URN:/
Response: HTML
說明: API測試頁面ssssssssssss
*/
func HandleGetAPITest(ginCTX *gin.Context) {

	vs_RequestURL := ginCTX.Request.URL.String()
	if strings.Index(vs_RequestURL, "?") > 1 {
		vs_RequestURL = vs_RequestURL[:strings.Index(vs_RequestURL, "?")]
	}

	ginCTX.HTML(http.StatusOK, Template_APITest, gin.H{
		H5Page_Menu:                 "./" + FolderPath_Pages + "/" + Template_IndexMenu,
		H5Value_Title:               "BFSystem : APITest 接口測試",
		H5Value_LastUpdate:          LastOnLineTime,
		H5Value_APITest_PostWebPath: URNPath_APITest[1:],
	})

}

/*
URN:/
Response: HTML
說明: API測試頁面
*/
// func HandlePostAPITest(ginCTX *gin.Context) {

// 	var err error
// 	var errorResult string
// 	var H5Return string
// 	var apiTestPostmanBody string
// 	var apiTestInputURL string
// 	var apiTestInputURN string
// 	var apiTestLoginAccount string
// 	var apiTestSignKey string
// 	var apiTestDesKey string
// 	var apiTestResponse string
// 	var apiTestDecodeResponse string
// 	var apiTestReviewPayload string
// 	var apiTestReviewCurl string
// 	var apiTestReviewHTTP string

// 	defer func() {
// 		if err_recover := recover(); err_recover != nil { //攔截伺服器POST所產無法預期錯誤
// 			errorResult = errorResult + "\n\n伺服器發生錯誤: " + fmt.Sprint(err_recover) + "\n\n" + string(debug.Stack())
// 			ginCTX.HTML(http.StatusOK, Template_APITest, gin.H{
// 				H5Page_Menu:        "./" + FolderPath_Pages + "/" + Template_IndexMenu,
// 				H5Value_Title:      "BFSystem : APITest 接口測試",
// 				H5Value_LastUpdate: LastOnLineTime,
// 				H5Value_Return:     H5Return,

// 				H5Value_APITest_InputURL:     apiTestInputURL,
// 				H5Value_APITest_InputURN:     apiTestInputURN,
// 				H5Value_APITest_LoginAccount: apiTestLoginAccount,
// 				H5Value_APITest_SignKey:      apiTestSignKey,
// 				H5Value_APITest_DesKey:       apiTestDesKey,
// 				H5Value_APITest_PostWebPath:  URNPath_APITest[1:],
// 				H5Value_APITest_PostmanBody:  apiTestPostmanBody,
// 				H5Value_APITest_Response:     errorResult,
// 			})
// 			ginCTX.Abort()
// 		} else { //正常傳送結果
// 			ginCTX.HTML(http.StatusOK, Template_APITest, gin.H{
// 				H5Page_Menu:        "./" + FolderPath_Pages + "/" + Template_IndexMenu,
// 				H5Value_Title:      "BFSystem : APITest 接口測試",
// 				H5Value_LastUpdate: LastOnLineTime,
// 				H5Value_Return:     H5Return,

// 				H5Value_APITest_InputURL:       apiTestInputURL,
// 				H5Value_APITest_InputURN:       apiTestInputURN,
// 				H5Value_APITest_LoginAccount:   apiTestLoginAccount,
// 				H5Value_APITest_SignKey:        apiTestSignKey,
// 				H5Value_APITest_DesKey:         apiTestDesKey,
// 				H5Value_APITest_PostWebPath:    URNPath_APITest[1:],
// 				H5Value_APITest_PostmanBody:    apiTestPostmanBody,
// 				H5Value_APITest_Response:       apiTestResponse,
// 				H5Value_APITest_DecodeResponse: apiTestDecodeResponse,
// 				H5Value_APITest_ReviewPayload:  apiTestReviewPayload,
// 				H5Value_APITest_ReviewHTTP:     apiTestReviewHTTP,
// 				H5Value_APITest_ReviewCurl:     apiTestReviewCurl,
// 			})
// 		}
// 	}()

// 	H5Return = fmt.Sprint(time.Now().Format("2006-01-02_15:04:05"))
// 	apiTestInputURL = ginCTX.PostForm(H5Value_APITest_InputURL)
// 	apiTestInputURN = ginCTX.PostForm(H5Value_APITest_InputURN)
// 	apiTestPostmanBody = ginCTX.PostForm(H5Value_APITest_PostmanBody)
// 	apiTestLoginAccount = ginCTX.PostForm(H5Value_APITest_LoginAccount)
// 	apiTestSignKey = ginCTX.PostForm(H5Value_APITest_SignKey)
// 	apiTestDesKey = ginCTX.PostForm(H5Value_APITest_DesKey)

// 	// vs_RequestURL := ginCTX.Request.URL.String()
// 	// if strings.Index(vs_RequestURL, "?") > 1 {
// 	// 	vs_RequestURL = vs_RequestURL[:strings.Index(vs_RequestURL, "?")]
// 	// }

// 	//步驟1 :  POST url檢查
// 	inputURLStr := ginCTX.PostForm(H5Value_APITest_InputURL) + ginCTX.PostForm(H5Value_APITest_InputURN)
// 	if !common.IsValidUrl(inputURLStr) {
// 		apiTestResponse = errorResult + "接口測試伺服器url 格式錯誤!\n"
// 		return
// 	}

// 	if inputURLStr[len(inputURLStr)-1:] == "/" { //自動處理尾部/
// 		inputURLStr = inputURLStr[:len(inputURLStr)-1]
// 	}

// 	//步驟2 : Postman body解析
// 	var postDataMap = make(map[string]string)                                            //用於判斷是否重複
// 	var postDataKeys = []string{}                                                        //紀錄key順序
// 	postmanBody := strings.Split(strings.ReplaceAll(apiTestPostmanBody, "\r", ""), "\n") //將換行符號政規化
// 	for _, value := range postmanBody {                                                  //解析每行postmanBody
// 		if value == "" { //空白行數
// 			continue //忽略
// 		} else if value[:2] == "//" { //Postman註解方式
// 			continue //忽略
// 		} else {
// 			postKey := url.QueryEscape(value[:strings.Index(value, ":")])
// 			postValue := value[strings.Index(value, ":")+1:]
// 			if _, isExsit := postDataMap[postKey]; isExsit {
// 				continue
// 			}
// 			postDataMap[postKey] = postValue
// 			postDataKeys = append(postDataKeys, postKey)
// 		}
// 	}
// 	// if len(postDataMap) <= 0 {
// 	// 	apiTestResponse = errorResult + "POST內容不能為空！\n"
// 	// 	return
// 	// }

// 	//步驟3 : POST postBody 生成
// 	cachePostBody := &bytes.Buffer{}
// 	cachePostBody.WriteString("{")
// 	for _, postKey := range postDataKeys {
// 		cachePostBody.WriteString(`"` + postKey + `":` + postDataMap[postKey] + `,`)
// 	}
// 	signJson := `"Key":"` + apiTestSignKey + `"` + "}"
// 	cachePostBody.WriteString(signJson)
// 	postBody := cachePostBody.String()
// 	//log.Println(postBody)

// 	//步驟4 : POST Payload 生成
// 	var cacheEncrypt string
// 	apiTestReviewPayload = postBody
// 	cacheEncrypt, err = encrypt.GetEncryptData(postBody, apiTestDesKey, signJson)
// 	if err != nil {
// 		apiTestResponse = errorResult + "加密錯誤:" + err.Error() + "\n"
// 		return
// 	}
// 	postPayload := fmt.Sprintf(`{"Account":"%v","Data":"%v"}`, apiTestLoginAccount, cacheEncrypt)

// 	//HTTPReview
// 	apiTestReviewHTTP = fmt.Sprintf(HTTPScriptFormat, apiTestInputURN, apiTestInputURL, APITest_ContentHeader, 33+len(apiTestLoginAccount)+len(cacheEncrypt), apiTestLoginAccount, cacheEncrypt)
// 	//Curl
// 	apiTestReviewCurl = fmt.Sprintf(CurlScriptFormat, inputURLStr, APITest_ContentHeader, apiTestLoginAccount, cacheEncrypt)

// 	//步驟5 : 執行POST測試
// 	postPayloadReader := strings.NewReader(postPayload)
// 	httpClient := &http.Client{}
// 	httpRequest, err := http.NewRequest("POST", inputURLStr, postPayloadReader)
// 	if err != nil {
// 		errorResult = errorResult + "\n" + err.Error()
// 		log.Println(err)
// 	}

// 	httpRequest.Header.Add("Content-Type", APITest_ContentHeader)
// 	httpResponse, err := httpClient.Do(httpRequest)
// 	if err != nil {
// 		errorResult = errorResult + "\n" + err.Error()
// 		log.Println(err)
// 	}

// 	cacheBytes, err := ioutil.ReadAll(httpResponse.Body)
// 	if err != nil {
// 		errorResult = errorResult + "\n" + err.Error()
// 		log.Println(err)
// 	}
// 	defer httpResponse.Body.Close()

// 	//步驟6 : POST結果處理
// 	if errorResult != "" {
// 		apiTestResponse = errorResult
// 	} else {
// 		cacheResponse := &bytes.Buffer{}
// 		if err := json.Indent(cacheResponse, cacheBytes, "", "  "); err != nil {
// 			apiTestResponse = string(cacheBytes) //非Json格式
// 		} else {
// 			apiTestResponse = cacheResponse.String() //Json格式

// 			cacheJson := make(map[string]string)
// 			err = json.Unmarshal([]byte(apiTestResponse), &cacheJson)
// 			if err != nil {
// 				apiTestDecodeResponse = err.Error()
// 				return
// 			}
// 			//fmt.Printf("%s", cacheJson["Data"])
// 			if data, isExsit := cacheJson["Data"]; isExsit && data != "" {
// 				apiTestDecodeResponse, err = encrypt.GetDecryptData(data, apiTestDesKey)
// 				if err != nil {
// 					apiTestDecodeResponse = err.Error()
// 					return
// 				} //對客戶API解碼方式使用商戶的DESKey

// 				if strings.Contains(apiTestDecodeResponse, "{") == false || strings.Contains(apiTestDecodeResponse, "}") == false {
// 					domanName := strings.Replace(inputURLStr, "http://", "", 1)
// 					domanName = strings.Replace(domanName, "https://", "", 1)

// 					if index := strings.Index(domanName, "/"); index > 1 {
// 						domanName = domanName[:index]
// 					}

// 					if len(domanName) >= 24 {
// 						apiTestDesKey = domanName[:24]
// 					} else {
// 						apiTestDesKey = domanName
// 						for index := 0; index < 24-len(domanName); index++ {
// 							apiTestDesKey = apiTestDesKey + "0"
// 						} //低於24字元則補0
// 					}

// 					apiTestDecodeResponse, err = encrypt.GetDecryptData(data, apiTestDesKey)
// 					if err != nil {
// 						apiTestDecodeResponse = err.Error()
// 						return
// 					}

// 				} //Android內部專用API採用網域前24字作為DESKey

// 				cacheDecode := &bytes.Buffer{}
// 				if err := json.Indent(cacheDecode, []byte(apiTestDecodeResponse), "", "  "); err != nil {
// 					return
// 				} //Data解碼後排序
// 				apiTestDecodeResponse = cacheDecode.String()

// 			}

// 		}
// 	}
// 	return

// }
