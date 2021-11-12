package common

import (
	"fmt"
	"testing"
)

func Test_CatchPanic(t *testing.T) {
	defer CatchPanic("單元測試")
	panic("這是單元測試") //此panic會遭攔截
}

func Test_ToCoverLeftToRightIP(t *testing.T) {

	ipv4 := "192.168.0.1"
	testCase1 := ToCoverLeftToRightIP(ipv4, 2)
	if testCase1 != "*.*.0.1" {
		t.Errorf("非預料輸出結果: %v!=%v\n", testCase1, "*.*.0.1")
	}

	testCase2 := ToCoverLeftToRightIP(ipv4, 0)
	if testCase2 != ipv4 {
		t.Errorf("非預料輸出結果: %v!=%v\n", testCase2, ipv4)
	}

	testCase3 := ToCoverLeftToRightIP(ipv4, 4)
	if testCase3 != ipv4 {
		t.Errorf("非預料輸出結果: %v!=%v\n", testCase3, ipv4)
	}

	ipv6 := "2001:db8:2de:0000:0000:0000:0000:e13"
	testCase4 := ToCoverLeftToRightIP(ipv6, 2)
	if testCase4 != "*:*:2de:0000:0000:0000:0000:e13" {
		t.Errorf("非預料輸出結果: %v!=%v\n", testCase1, "*:*:2de:0000:0000:0000:0000:e13")
	}

	testCase5 := ToCoverLeftToRightIP(ipv6, 0)
	if testCase5 != ipv6 {
		t.Errorf("非預料輸出結果: %v!=%v\n", testCase2, ipv6)
	}

	testCase6 := ToCoverLeftToRightIP(ipv6, 8)
	if testCase6 != ipv6 {
		t.Errorf("非預料輸出結果: %v!=%v\n", testCase3, ipv6)
	}
}

func Test_GetUnixNow(t *testing.T) {

	testCase1 := fmt.Sprint(GetUnixNowSec())
	if len(testCase1) != 10 {
		t.Errorf("非預料輸出結果: len(%v)!=%v UnixTime位數不正確\n", testCase1, 10)
	}

	testCase2 := fmt.Sprint(GetUnixNowMS())
	if len(testCase2) != 13 {
		t.Errorf("非預料輸出結果: len(%v)!=%v UnixTime位數不正確\n", testCase2, 13)
	}

	testCase3 := fmt.Sprint(GetUnixNowSinceNStoMS(GetUnixNowMS()))
	if len(testCase3) != 13 {
		t.Errorf("非預料輸出結果: len(%v)!=%v UnixTime位數不正確\n", testCase3, 13)
	}

}

func Test_StringAndBytes(t *testing.T) {

	testCase1 := string(StringToBytes("Test"))
	if testCase1 != "Test" {
		t.Errorf("非預料輸出結果: %v!=%v 字串轉換不正確\n", testCase1, "Test")
	}

	testCase2 := BytesToString([]byte("9999"))
	if testCase2 != "9999" {
		t.Errorf("非預料輸出結果: %v!=%v 字串轉換不正確\n", testCase2, "9999")
	}

}

func Test_IsValidUrlAndIP(t *testing.T) {

	testCase1 := IsValidUrl("https://192.168.0.1:2560")
	if testCase1 != true {
		t.Errorf("非預料輸出結果: %v!=%v URL判斷不正確\n", testCase1, true)
	}

	testCase2 := IsValidUrl("https://192.168.0.1:0")
	if testCase2 != true {
		t.Errorf("非預料輸出結果: %v!=%v URL判斷不正確\n", testCase2, true)
	}

	testCase3 := IsValidUrl("192.168.0.1")
	if testCase3 != false {
		t.Errorf("非預料輸出結果: %v!=%v URL判斷不正確\n", testCase3, false)
	}

	testCase4 := IsValidUrl("2001:db8:2de:0000:0000:0000:0000:e13")
	if testCase4 != false {
		t.Errorf("非預料輸出結果: %v!=%v URL判斷不正確\n", testCase4, false)
	}

	testCase5 := IsValidUrl("")
	if testCase5 != false {
		t.Errorf("非預料輸出結果: %v!=%v URL判斷不正確\n", testCase5, false)
	}

	testCase6 := IsValidIP("192.168.0.1")
	if testCase6 != true {
		t.Errorf("非預料輸出結果: %v!=%v IP判斷不正確\n", testCase6, true)
	}

	testCase7 := IsValidIP("192.168.0.1/99")
	if testCase7 != false {
		t.Errorf("非預料輸出結果: %v!=%v IP判斷不正確\n", testCase7, false)
	}

	testCase8 := IsValidIP("192.168.0.999")
	if testCase8 != false {
		t.Errorf("非預料輸出結果: %v!=%v IP判斷不正確\n", testCase8, false)
	}

}

func Test_IsSafeSQLValue(t *testing.T) {

	testCase1 := IsSafeSQLValue("\n")
	if testCase1 != false {
		t.Errorf("非預料輸出結果: %v!=%v SQL注入判斷不正確\n", testCase1, false)
	}

	testCase2 := IsSafeSQLValue("\t")
	if testCase2 != false {
		t.Errorf("非預料輸出結果: %v!=%v SQL注入判斷不正確\n", testCase2, false)
	}

	testCase3 := IsSafeSQLValue("\r")
	if testCase3 != false {
		t.Errorf("非預料輸出結果: %v!=%v SQL注入判斷不正確\n", testCase3, false)
	}

	testCase4 := IsSafeSQLValue(";")
	if testCase4 != false {
		t.Errorf("非預料輸出結果: %v!=%v SQL注入判斷不正確\n", testCase4, false)
	}

	testCase5 := IsSafeSQLValue("'")
	if testCase5 != false {
		t.Errorf("非預料輸出結果: %v!=%v SQL注入判斷不正確\n", testCase5, false)
	}

	testCase6 := IsSafeSQLValue(`"`)
	if testCase6 != false {
		t.Errorf("非預料輸出結果: %v!=%v SQL注入判斷不正確\n", testCase6, false)
	}

	testCase7 := IsSafeSQLValue("`")
	if testCase7 != false {
		t.Errorf("非預料輸出結果: %v!=%v SQL注入判斷不正確\n", testCase7, false)
	}

	testCase8 := IsSafeSQLValue(" SELECT DISTINCT server_uid FROM server_info WHERE status = 1")
	if testCase8 != true {
		t.Errorf("非預料輸出結果: %v!=%v SQL注入判斷不正確\n", testCase8, true)
	}

	testCase9 := IsSafeSQLValue(" SELECT DISTINCT server_uid FROM server_info WHERE status = 1\n;")
	if testCase9 != false {
		t.Errorf("非預料輸出結果: %v!=%v SQL注入判斷不正確\n", testCase9, false)
	}
}

func Test_RoundFloat(t *testing.T) {

	testCase1 := RoundFloat64(99.999, 2)
	if testCase1 != 100.0 {
		t.Errorf("非預料輸出結果: %v!=%v 小數進行四捨五入判斷不正確\n", testCase1, 100.0)
	}

	testCase2 := RoundFloat64(33.33, 0)
	if testCase2 != 33.0 {
		t.Errorf("非預料輸出結果: %v!=%v 小數進行四捨五入判斷不正確\n", testCase2, 33.0)
	}

	testCase3 := RoundFloat64(-99.0, 1)
	if testCase3 != -99 {
		t.Errorf("非預料輸出結果: %v!=%v 小數進行四捨五入判斷不正確\n", testCase3, -99)
	}

	testCase4 := RoundUpInt(99.999)
	if testCase4 != 100 {
		t.Errorf("非預料輸出結果: %v!=%v 小數進行四捨五入判斷不正確\n", testCase4, 100.0)
	}

	testCase5 := RoundUpInt(33.33)
	if testCase5 != 34 {
		t.Errorf("非預料輸出結果: %v!=%v 小數進行四捨五入判斷不正確\n", testCase5, 34)
	}

	testCase6 := RoundUpInt(-99.0)
	if testCase6 != -99 {
		t.Errorf("非預料輸出結果: %v!=%v 小數進行四捨五入判斷不正確\n", testCase5, -99)
	}
}


func Test_GetClientPublicIP(t *testing.T) {

	testCase1 := GetClientPublicIP(" 123.456.111.61 ", "")
	if testCase1 != "123.456.111.61" {
		t.Errorf("非預料輸出結果: %v!=%v ClientIP判斷不正確\n", testCase1, "123.456.111.61")
	}
	testCase2 := GetClientPublicIP(" 192.168.0.1 , , [::1], 10.0.0.1 , 172.17.0.1, 123.456.111.61 ", "")
	if testCase2 != "123.456.111.61" {
		t.Errorf("非預料輸出結果: %v!=%v ClientIP判斷不正確\n", testCase2, "10.0.0.1")
	}

	testCase3 := GetClientPublicIP(" 192.168.0.1 , , [::1], 10.0.0.1 , 172.17.0.1", "123.456.111.61")
	if testCase3 != "123.456.111.61" {
		t.Errorf("非預料輸出結果: %v!=%v ClientIP判斷不正確\n", testCase3, "123.456.111.61")
	}

	testCase4 := GetClientPublicIP("", "123.456.111.61")
	if testCase4 != "123.456.111.61" {
		t.Errorf("非預料輸出結果: %v!=%v ClientIP判斷不正確\n", testCase4, "123.456.111.61")
	}

	testCase5 := GetClientPublicIP(",", "")
	if testCase5 != "" {
		t.Errorf("非預料輸出結果: %v!=%v ClientIP判斷不正確\n", testCase5, "")
	}


}

