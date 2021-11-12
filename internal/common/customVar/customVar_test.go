package customVar

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func TestPkg(t *testing.T) {

	var err error

	var defBool bool
	err = SetValue("true", &defBool, &SwitchType{})
	if err != nil || defBool != true { //輸出結果應為: true (bool)
		t.Errorf("非預料輸出結果: %v (%T), 錯誤提示: %v\n", defBool, defBool, err)
	} else {
		t.Logf("輸出結果: %v (%T)\n", defBool, defBool)
	}

	var exBool bool
	err = SetConfig("ON", false, &exBool, &SwitchType{})
	if err != nil || exBool != true { //輸出結果應為: true (bool)
		t.Errorf("非預料輸出結果: %v (%T), 錯誤提示: %v\n", exBool, exBool, err)
	} else {
		t.Logf("輸出結果: %v (%T)\n", exBool, exBool)
	}

	var exInt int8
	err = SetConfig("60.1", "90", &exInt, &SecondsInADay{})
	if err != nil || exInt != 60 { //輸出結果應為: 60 (int8)
		t.Errorf("非預料輸出結果: %v (%T), 錯誤提示: %v\n", exInt, exInt, err)
	} else {
		t.Logf("輸出結果: %v (%T)\n", exInt, exInt)
	}

	var stringInt string
	err = SetConfig("xx60", "xx90", &stringInt, &StringType{})
	if err != nil || stringInt != "xx60" { //輸出結果應為: xx60 (string)
		t.Errorf("非預料輸出結果: %v (%T), 錯誤提示: %v\n", stringInt, stringInt, err)
	} else {
		t.Logf("輸出結果: %v (%T)\n", stringInt, stringInt)
	}

	var stringWrongInt string
	err = SetConfig("60.1", "90", &stringWrongInt, &Int8Type{})
	if err != nil || stringWrongInt != "90" { //輸出結果應為: 60 (string)
		t.Logf("輸出結果: %v (%T), 錯誤提示: %v\n", stringWrongInt, stringWrongInt, err)
	} else {
		t.Errorf("非預料輸出結果: %v (%T)\n", stringWrongInt, stringWrongInt)
	}

	var exUint uint64
	err = SetConfig("-99.9", "11", &exUint, &customStruct{})
	if err != nil && exUint == 11 { //輸出結果應為: 11 (uint64), 錯誤提示: input value > '-99.9' > 不可低於10 strconv.ParseUint: parsing "-99": invalid syntax ; default value > <nil>
		t.Logf("輸出結果: %v (%T), 錯誤提示: %v\n", exUint, exUint, err)
	} else {
		t.Errorf("非預料輸出結果: %v (%T)\n", exUint, exUint)
	}

	var exSlice []uint32
	exSlice = append(exSlice, 123)
	err = SetConfig("456", "111", &exSlice, &BaseSlice{})
	if err != nil || exSlice[0] != 123 || exSlice[1] != 456 { //輸出結果應為:[123 456] ([]uint32)
		t.Errorf("非預料輸出結果: %v (%T), 錯誤提示: %v\n", exSlice, exSlice, err)
	} else {
		t.Logf("輸出結果: %v (%T)\n", exSlice, exSlice)
	}

	var exCustomTypeSlice []interface{}
	exCustomTypeSlice = append(exCustomTypeSlice, 456)
	err = SetConfig("Test", "Def", &exCustomTypeSlice, &BaseSlice{})
	if err != nil || exCustomTypeSlice[0] != 456 || exCustomTypeSlice[1] != "Test" { //輸出結果應為: [456 Test] ([]interface {})
		t.Errorf("非預料輸出結果: %v (%T), 錯誤提示: %v\n", exCustomTypeSlice, exCustomTypeSlice, err)
	} else {
		t.Logf("輸出結果: %v (%T)\n", exCustomTypeSlice, exCustomTypeSlice)
	}

	exStringMap := make(map[string]interface{})
	exStringMap["原始Key"] = "原始Value"
	err = SetConfig("VVVVVVVVVVVVV", exStringMap, &exStringMap, &stringMapType{SetKey: "KeyTest"})
	if err != nil || exStringMap["KeyTest"] != "VVVVVVVVVVVVV" || exStringMap["原始Key"] != "原始Value" { //輸出結果應為: map[KeyTest:VVVVVVVVVVVVV 原始Key:原始Value] (map[string]interface {})
		t.Errorf("非預料輸出結果: %v (%T), 錯誤提示: %v\n", exStringMap, exStringMap, err)
	} else {
		t.Logf("輸出結果: %v (%T)\n", exStringMap, exStringMap)
	}

	var exCustomTypeMap interface{}
	err = SetConfig("testKey|t1|t2|t3|t4", nil, &exCustomTypeMap, &customTypeMap{})
	exCustomTypeMapAssertion := exCustomTypeMap.(map[string]*CustomTypeValue)
	if err != nil || exCustomTypeMapAssertion["testKey"].f1 != "t1" || exCustomTypeMapAssertion["testKey"].f2 != "t2" || exCustomTypeMapAssertion["testKey"].f3 != "t3" || exCustomTypeMapAssertion["testKey"].f4 != "t4" { //輸出結果應為: testKey &{t1 t2 t3 t4} (*main.CustomTypeValue)
		t.Errorf("非預料輸出結果: %v (%T), 錯誤提示: %v\n", exCustomTypeMapAssertion, exCustomTypeMapAssertion, err)
	} else {
		t.Logf("輸出結果: %v (%T)\n", exCustomTypeMap, exCustomTypeMap)
	}

}

/*
	工廠方法函式說明:
	func SetConfig(inputValue string, defaultValue, outputPointer interface{}, Custom customTypes) (err error)
	> 將輸入字串轉換任意變數 (若轉換失敗則輸出設定預設值)
	@param inputValue: 外部輸入變數(需先轉成字串)
	@param defaultValue: 預設值(當自定型別轉換發生錯誤時則輸出此值)
	@param outputPointer:  寫入輸出變數的指標
	@param Custom:  自訂型別，可OOP繼承與方法重載
	@return error : 錯誤提示
*/
type customStruct struct{} //使用者可以自訂型別,藉由OOP繼承與方法重載擴增功能

func (_ customStruct) GetValue(inputValue string) (output interface{}, err error) { //使用者可依需求調整功能

	var cache int
	if index := strings.Index(inputValue, "."); index < 1 { //無條件捨去小數，直接直接取整數
		cache, err = strconv.Atoi(inputValue)
	} else {
		cache, err = strconv.Atoi(inputValue[:index])
	}

	if err != nil {
		return cache, err
	} else if cache < 10 {
		return cache, fmt.Errorf("'%v' > 不可低於10", inputValue)
	} else {
		return cache, nil
	}
}

type stringMapType struct {
	SetKey string
}

func (mp *stringMapType) GetValue(inputValue string) (output interface{}, err error) {
	if len(inputValue) < 10 {
		return nil, fmt.Errorf("'%v' > 字串長度不可低於10", inputValue)
	}

	cache := make(map[string]interface{})
	cache[mp.SetKey] = inputValue
	return cache, nil
}

type customTypeMap map[string]*CustomTypeValue

type CustomTypeValue struct {
	f1 string
	f2 string
	f3 string
	f4 string
}

func (mp customTypeMap) GetValue(inputValue string) (output interface{}, err error) {
	cacheSlice := strings.Split(inputValue, "|")
	if len(cacheSlice) < 5 {
		return nil, fmt.Errorf("'%v' > ' | ' 分隔號數量不可低於5", inputValue)
	} else if cacheSlice[1] == "" || cacheSlice[2] == "" || cacheSlice[3] == "" || cacheSlice[4] == "" {
		return nil, fmt.Errorf("'%v' > 部份數值不可為空", inputValue)
	}

	cacheMap := make(map[string]*CustomTypeValue)
	cacheMap[cacheSlice[0]] = &CustomTypeValue{f1: cacheSlice[1], f2: cacheSlice[2], f3: cacheSlice[3], f4: cacheSlice[4]}
	return cacheMap, nil
}
