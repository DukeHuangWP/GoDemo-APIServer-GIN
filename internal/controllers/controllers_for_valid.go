package controllers

import (
	"DevIntergTest/internal/common"
	"DevIntergTest/internal/models"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	URN_ValidRootPath = "/Valid"

	URN_Valid_Task_Tag = "ValidTag"
	URN_Valid_Task_ID  = "ValidID"
	URN_Valid_TaskPath = "/:" + URN_Valid_Task_Tag + "/*" + URN_Valid_Task_ID

	URN_Valid_fundlog = "/fundlog" + URN_Valid_TaskPath

	Valid_Task_Tag_TaskNo  = "TaskNo"
	Valid_Task_Tag_TaskUID = "TaskUID"

	TaskNoTitle_StartString = "_"
	TaskNoTitle_Title1      = "Title1_"
	TaskNoTitle_Title2      = "Title2_"

	TaskStatus_ManualRequired = -2 //-2:人工确认
	TaskStatus_Fail           = -1 //-1: 失败
	TaskStatus_Untreated      = 0  //0:未处理
	TaskStatus_Processing     = 1  //1:处理中
	TaskStatus_Success        = 2  //2:成功
	TaskStatus_PartialSuccess = 3  // 3:部份成功

	MsgCode_Err  = "500"
	MsgCode_Warn = "200"
	MsgCode_Pass = "210"
	MsgCode_Fail = "220"
)

type ResponseValid struct {
	MsgCode string `json:"MsgCode"` //回傳處理結果代碼
	Msg     string `json:"Msg"`     //回傳處理結果訊息
	Data    string `json:"Data"`    //加密後的請求資料
}

func MidValid(ginCTX *gin.Context) {
	validType := ginCTX.Param(URN_Valid_Task_Tag)
	validTag := ginCTX.Param(URN_Valid_Task_ID)
	if len(validTag) > 0 {
		validTag = validTag[1:]
	}

	switch validType {
	case Valid_Task_Tag_TaskNo:
		TaskNo := validTag
		ginCTX.Set(Valid_Task_Tag_TaskNo, TaskNo)
	case Valid_Task_Tag_TaskUID:
		TaskUID := validTag
		ginCTX.Set(Valid_Task_Tag_TaskNo, "models.DoSQLs.GetTaskNoFromTaskUID(TaskUID)"+TaskUID)
	default:
		ginCTX.Abort()
		return
	}

	ginCTX.Next() // 將Set內之參數傳給下一個gin.HandlerFunc

}

func EndNodeValidTask(ginCTX *gin.Context) {

	interfaceKey, isExsit := ginCTX.Get(Valid_Task_Tag_TaskNo)
	if isExsit == false || interfaceKey == nil {
		ginCTX.JSON(http.StatusBadRequest, &ResponseValid{
			MsgCode: MsgCode_Warn,
			Msg:     fmt.Sprintf("未正確輸入 %v", Valid_Task_Tag_TaskNo),
		})
		return
	}

	taskNo := fmt.Sprint(interfaceKey)
	taskNoTypeStartIndex := strings.Index(taskNo, TaskNoTitle_StartString)
	fmt.Println(taskNo)

	if taskNoTypeStartIndex < 1 || len(taskNo) < taskNoTypeStartIndex+2 {
		ginCTX.JSON(http.StatusBadRequest, &ResponseValid{
			MsgCode: MsgCode_Warn,
			Msg:     fmt.Sprintf("taskNo格式無法判別 > %v", Valid_Task_Tag_TaskNo),
		})
		return
	}

	switch taskNo[:taskNoTypeStartIndex+1] {
	case TaskNoTitle_Title1:
		ginCTX.Next()
		TaskNoTitleTitle1(ginCTX)
		return
	case TaskNoTitle_Title2:
		ginCTX.Next()
		TaskNoTitleTitle2(ginCTX)
		return
	}
	ginCTX.JSON(http.StatusBadRequest, &ResponseValid{
		Msg: fmt.Sprintf("無法驗證此類型 taskNo  > %v", Valid_Task_Tag_TaskNo),
	})
	return

}

func TaskNoTitleTitle1(ginCTX *gin.Context) {

	var err error
	var taskNo string
	var isAPIDone bool
	defer func() {
		if isAPIDone == false {
			ginCTX.AbortWithStatusJSON(http.StatusInternalServerError, &ResponseValid{
				MsgCode: MsgCode_Err,
				Msg:     fmt.Sprintf("執行過程中發生錯誤 > 設計錯誤!!!"),
			})
			return
		} else if err != nil {
			ginCTX.AbortWithStatusJSON(http.StatusInternalServerError, &ResponseValid{
				MsgCode: MsgCode_Err,
				Msg:     fmt.Sprintf("執行過程中發生錯誤 > %v", err),
			})
			return
		}

	}()

	interfaceKey, isExsit := ginCTX.Get(Valid_Task_Tag_TaskNo)
	if isExsit == false || interfaceKey == nil {
		ginCTX.JSON(http.StatusBadRequest, &ResponseValid{
			MsgCode: MsgCode_Warn,
			Msg:     fmt.Sprintf("未正確輸入 %v", Valid_Task_Tag_TaskNo),
		})
		isAPIDone = true
		return
	}
	taskNo = fmt.Sprint(interfaceKey)
	query, err := models.DoSQLs.GetFromDeviceInfo([]string{"taskStatus", "creatTime"}, taskNo)
	if err != nil {
		isAPIDone = true
		return
	}

	taskStatus, ok := query["taskStatus"].(int)
	if ok == false {
		ginCTX.JSON(http.StatusBadRequest, &ResponseValid{
			MsgCode: MsgCode_Warn,
			Msg:     fmt.Sprintf("未正確輸入 %v", Valid_Task_Tag_TaskNo),
		})
		isAPIDone = true
		return
	}
	creatTime, ok := query["creatTime"].(int64)
	if ok == false {
		ginCTX.JSON(http.StatusBadRequest, &ResponseValid{
			MsgCode: MsgCode_Warn,
			Msg:     fmt.Sprintf("未正確輸入 %v", Valid_Task_Tag_TaskNo),
		})
		isAPIDone = true
		return
	}

	if creatTime == 0 {
		ginCTX.JSON(http.StatusOK, &ResponseValid{
			MsgCode: MsgCode_Warn,
			Msg:     fmt.Sprintf("驗證失敗 : %v 此單並不存在", taskNo),
		})
		isAPIDone = true
		return
	}

	nowTime := common.GetUnixNowSec()
	countTime := nowTime - creatTime
	switch taskStatus {
	case TaskStatus_Untreated:

		if countTime > 300 {
			ginCTX.JSON(http.StatusOK, &ResponseValid{
				MsgCode: MsgCode_Fail,
				Msg:     fmt.Sprintf("驗證失敗 : %v 此單 task_status = '%v' creatTime > '%v - %v = %vs' 不應未完成狀態!", Valid_Task_Tag_TaskNo, taskStatus, nowTime, creatTime, countTime),
			})
		} else {
			ginCTX.JSON(http.StatusOK, &ResponseValid{
				MsgCode: MsgCode_Warn,
				Msg:     fmt.Sprintf("驗證失敗 : %v 此單 task_status = '%v' creatTime > '%v - %v = %vs' 尚未未完成狀態!", Valid_Task_Tag_TaskNo, taskStatus, nowTime, creatTime, countTime),
			})
		}
		isAPIDone = true
		return
	case TaskStatus_Processing:
		ginCTX.JSON(http.StatusOK, &ResponseValid{
			MsgCode: MsgCode_Warn,
			Msg:     fmt.Sprintf("驗證失敗 : %v 此單 task_status = '%v' creatTime > '%v - %v = %vs' 處理中狀態!", Valid_Task_Tag_TaskNo, taskStatus, nowTime, creatTime, countTime),
		})
		isAPIDone = true
		return

	case TaskStatus_ManualRequired:

		ginCTX.JSON(http.StatusOK, &ResponseValid{
			MsgCode: MsgCode_Warn,
			Msg:     fmt.Sprintf("驗證失敗 : %v 此單 task_status = '%v' creatTime > '%v - %v = %vs' 需要人工確認!", Valid_Task_Tag_TaskNo, taskStatus, nowTime, creatTime, countTime),
		})
		isAPIDone = true
		return
	}

	ginCTX.JSON(http.StatusOK, &ResponseValid{
		MsgCode: MsgCode_Pass,
		Msg:     "測試正常 " + taskNo,
	})
	isAPIDone = true
	return
}


func TaskNoTitleTitle2(ginCTX *gin.Context) {

	var err error
	var isAPIDone bool
	defer func() {
		if isAPIDone == false {
			ginCTX.AbortWithStatusJSON(http.StatusInternalServerError, &ResponseValid{
				MsgCode: MsgCode_Err,
				Msg:     fmt.Sprintf("執行過程中發生錯誤 > 設計錯誤!!!"),
			})
			return
		} else if err != nil {
			ginCTX.AbortWithStatusJSON(http.StatusInternalServerError, &ResponseValid{
				MsgCode: MsgCode_Err,
				Msg:     fmt.Sprintf("執行過程中發生錯誤 > %v", err),
			})
			return
		}

	}()


	isAPIDone = true
	return
}
