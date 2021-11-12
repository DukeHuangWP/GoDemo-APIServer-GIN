package common

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	"unsafe"
)

//攔截Panic，防止非預期錯誤導致當機
func CatchPanic(errorTitle string) {
	if err := recover(); err != nil {
		var logText string
		for index := 0; index < 8; index++ { //最多捕捉8層log
			ptr, filename, line, ok := runtime.Caller(index)
			if !ok {
				break
			}
			logText = logText + fmt.Sprintf(" %v:%d,%v > %v\n", filename, line, runtime.FuncForPC(ptr).Name(), err)
		}
		log.Printf("%v : 發生嚴重錯誤 %v\n", errorTitle, logText)
		return
	}
}

//獲取Unix時間戳ms
func GetUnixNowSec() int64 {
	return time.Now().Unix()
}

func GetUnixNowMS() int64 {
	return time.Now().UnixNano() / 1e6
}

func GetUnixNowNS() int64 {
	return time.Now().UnixNano()
}

func GetUnixNowSinceNStoMS(unixNanoTime int64) int64 {
	return (time.Now().UnixNano() - unixNanoTime) / 1e6
}

func GetTimeNowFormat() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func GetTimeNowDateFormat() string {
	return time.Now().Format("2006-01-02")
}

// StringToBytes converts string to byte slice without a memory allocation.
func StringToBytes(s string) (b []byte) {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

// BytesToString converts byte slice to string without a memory allocation.
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

//解析網址格式是否正確
func IsValidUrl(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

//解析網址格式是否正確
func IsValidIP(toTest string) bool {

	ip := net.ParseIP(toTest)
	if ip == nil {
		return false
	}

	if ip.To4() == nil || ip.To16() == nil { //ipv4 or ipv6
		return false
	}
	return true

}

//檢查該字串是否包含SQL注入字串
func IsSafeSQLValue(SQLValue string) bool {

	if strings.Contains(SQLValue, "\n") ||
		strings.Contains(SQLValue, "\t") ||
		strings.Contains(SQLValue, "\r") ||
		strings.Contains(SQLValue, ";") ||
		strings.Contains(SQLValue, "'") ||
		strings.Contains(SQLValue, `"`) ||
		strings.Contains(SQLValue, "`") {
		return false
	} //初步防止遭受SQL注入字串

	return true

}

//檢測資料夾目錄是否存在
func IsFolderExist(folderPath string) bool {
	info, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		return false
	}
	if info.IsDir() {
		return true
	}
	return false
}

//檢測檔案目錄是否存在
func IsFileExist(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	if !info.IsDir() {
		return true
	}
	return false
}

//清空檔案內的資料
func TruncateFile(filePath string, writeBytes []byte) error {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755) //危險此指令會先將檔案內資料刪除
	if err != nil {
		return err
	}
	_, err = file.Write(writeBytes)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

//獲取目錄內所有檔案名稱清單
func GetFileList(FilePathInput string) (fileList []string, err error) {
	pty_OpenFile, err := os.Open(FilePathInput)
	if err != nil {
		return fileList, err
	}

	ty_FileInfos, err := pty_OpenFile.Readdir(-1)
	pty_OpenFile.Close()
	if err != nil {
		return fileList, err
	}
	for _, value := range ty_FileInfos {
		if !value.IsDir() {
			fileList = append(fileList, value.Name())
		}
	}
	return fileList, nil
}

//獲取目錄內所有資料夾名稱清單
func GetFolderList(FilePathInput string) (folderList []string, err error) {
	pty_OpenFile, err := os.Open(FilePathInput)
	if err != nil {
		return folderList, err
	}

	ty_FileInfos, err := pty_OpenFile.Readdir(-1)
	pty_OpenFile.Close()
	if err != nil {
		return folderList, err
	}
	for _, value := range ty_FileInfos {
		if value.IsDir() {
			folderList = append(folderList, value.Name())
		}
	}
	return folderList, nil
}

//複製檔案
func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

//複製目錄
func CopyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	// if err == nil {
	// 	return fmt.Errorf("destination already exists")
	// }

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}

//甚用：清空目錄內所有檔案
func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

//小數進行四捨五入
func RoundFloat64(input float64, decimals uint32) float64 {
	return math.Round(input*math.Pow(10, float64(decimals))) / math.Pow(10, float64(decimals))
}

//小數無條件升成整數
func RoundUpInt(input float64) int {
	return int(math.Ceil(input))
}

//遮蔽前段部分IP,例如: 192.168.0.1 > *.*.0.1
func ToCoverLeftToRightIP(originalIP string, coverNum int) (coveredIP string) {

	if coverNum <= 0 {
		return originalIP //不作任何變更
	}

	if strings.Count(originalIP, ".") > coverNum { //IPv4
		cacheList := strings.Split(originalIP, ".")
		cacheList = cacheList[coverNum:]
		for index := 0; index < coverNum-1; index++ {
			coveredIP = coveredIP + "*."
		}
		coveredIP = coveredIP + "*"

		for index := 0; index < len(cacheList); index++ {
			coveredIP = coveredIP + "." + cacheList[index]
		}

		return coveredIP
	}

	if strings.Count(originalIP, ":") > coverNum == true && strings.Contains(originalIP, ".") == false { //IPv6
		cacheList := strings.Split(originalIP, ":")
		cacheList = cacheList[coverNum:]
		for index := 0; index < coverNum-1; index++ {
			coveredIP = coveredIP + "*:"
		}
		coveredIP = coveredIP + "*"

		for index := 0; index < len(cacheList); index++ {
			coveredIP = coveredIP + ":" + cacheList[index]
		}

		return coveredIP
	}

	return originalIP //不作任何變更

}

func ComplementTWO(input string) (output string) {

	if dotIndex := strings.LastIndex(input, "."); dotIndex == -1 {
		return input + ".00"
	} else if count := len(input[dotIndex+1:]); count < 2 { //補上小數後兩位
		for index := 0; index < 2-count; index++ { //補上小數後兩位
			input = input + "0"
		}
	}
	return input
} //字串最後小數補0,待優化

//判斷代理服務器所傳遞header中的ip清單中取出外網ip,否則回傳finalIP| ipCommaList:代理服務器ip清單 192.168.0.1,123.456.0.1 returnIP:最終回傳IP
func GetClientPublicIP(ipCommaList, returnIP string) (clientIP string) {
	ipList := strings.Split(ipCommaList, ",")
	for index := 0; index < len(ipList); index++ {
		if theIP := strings.TrimSpace(ipList[index]); len(theIP) > 7 {
			indeNumber := theIP[:3]
			if indeNumber != "10." && indeNumber != "127" && indeNumber != "172" && indeNumber != "192" && indeNumber != "[::" {
				return theIP
			}
		}
	}
	return returnIP
}

//判斷ip是否為localIP, 0.0開頭,127開頭,172開頭,192開頭,10開頭,[::開頭
func IsALocalIP(ip string) (isALocalIP bool) {

	if theIP := strings.TrimSpace(ip); len(theIP) >= 3 {
		indexIP := theIP[:3]

		if indexIP == "0.0" ||
			indexIP == "10." ||
			indexIP == "127" ||
			indexIP == "172" ||
			indexIP == "192" ||
			indexIP == "::1" ||
			indexIP == "[::" {
			return true
		}
	}

	return false
}
