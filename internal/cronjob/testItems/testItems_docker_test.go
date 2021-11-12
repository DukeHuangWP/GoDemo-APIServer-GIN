package testItems

import (
	"testing"
)

func TestDockerWatcher(testT *testing.T) {

	TaskType := "DockerWatcher"

	if IsTestTypeExsit(TaskType) == false {
		testT.Errorf("測試類型不存在程式中!")
		return
	}

	scriptFromDB := `
	{
		"AlertLinesUpper": 2,
		"AlertLinesLower": 10,
		"Dockers": {
			"go_API":{"Since":"60m","Keywords":["panic:","panic:"]},
			"go_DS":{"Since":"60m","Keywords":["panic:"]}
		}
	}
	`
	var err error
	testItemsInterface := TestType[TaskType]
	err = testItemsInterface.ValidScript("testName", "testTag", scriptFromDB)
	if err != nil {
		testT.Errorf(" Script設定值 > %v 驗證發生錯誤 > %v", scriptFromDB, err)
		return
	}

	_, err = testItemsInterface.GetScriptFunc()
	if err != nil {
		testT.Errorf("GetScriptFunc() 過程發生錯誤>%v", err)
		return
	}

}
