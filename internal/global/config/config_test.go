package config

import (
	"DevIntergTest/internal/common/customVar"
	"testing"
)

func Test_InitLoad(t *testing.T) {
	testCase(t)
	t.Logf("SettingMap : %v\n", SettingMap)
}

func Test_Reload(t *testing.T) {
	testCase(t)
	t.Logf("SettingMap : %v\n", SettingMap)
}

func testCase(t *testing.T) {
	var exInt uint8
	SettingMap["unit_test_uint8"] = &AddSetings{DefaultValue: 0, OutputPointer: &exInt, Custom: &customVar.Uint8Type{}}
	SettingMap["unit_test_uint8"].ToWrite("9")
	if exInt != 9 { //輸出結果應為: 9 (uint8)
		t.Errorf("非預料輸出結果: %v (%T)\n", exInt, exInt)
	} else {
		t.Logf("輸出結果: %v (%T)\n", exInt, exInt)
	}

	var exFloat float64
	SettingMap["unit_test_float64"] = &AddSetings{DefaultValue: 11.1, OutputPointer: &exFloat, Custom: &customVar.Float64Type{}}
	SettingMap["unit_test_float64"].ToDefault()
	if exFloat != 11.1 { //輸出結果應為: 11.1 (float64)
		t.Errorf("非預料輸出結果: %v (%T)\n", exFloat, exFloat)
	} else {
		t.Logf("輸出結果: %v (%T)\n", exFloat, exFloat)
	}
}
