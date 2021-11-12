package loggerToDB

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
	"time"
)


var rwMutex = sync.RWMutex{}

//計算出 logger 的表名稱
func countLoggerTableNames(tableNameTitie string) (nowPartName, nextPartName string) {

	nowTime := time.Now()
	nowYear, nowMonthTime, nowDay := nowTime.Date()
	nowPartCount := nowDay / 16
	nowPartName = fmt.Sprintf("%v_%v_%02d_%02dh", tableNameTitie, nowYear, int(nowMonthTime), nowPartCount)

	nextPartTime := nowTime.Add((time.Hour * 24 * 15)) //取得下個月
	nextPartYear, nextPartMonthTime, nextPartDay := nextPartTime.Date()
	nextPartCount := nextPartDay / 16
	nextPartName = fmt.Sprintf("%v_%v_%02d_%02dh", tableNameTitie, nextPartYear, int(nextPartMonthTime), nextPartCount)

	//當月01~16日為上半,當月17~31日為下半
	return
}

//創建 NowTableName 和 NextTableName的表
func creatTables() (err error) {

	rows, err := theDB.Raw("SHOW TABLES;").Rows()
	if err != nil {
		return err
	}
	defer rows.Close()
	tableNameList := []string{}
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return err
		}

		if strings.Contains(tableName, TableNameTitle) {
			tableNameList = append(tableNameList, tableName)
		} //將含有TableNameTitle字串的表格名稱記入
	}

	if len(tableNameList) <= 0 {
		return fmt.Errorf("不存在舊表%v", TableNameTitle)
	}

	sort.Strings(tableNameList)                                        //強制排序確保日期
	latestTableName := tableNameList[len(tableNameList)-1]             //取最新的表名稱
	nowPartName, nextPartName := countLoggerTableNames(TableNameTitle) //計算表名稱
	rwMutex.Lock()                                                     //寫入上鎖
	NowTableName = nowPartName
	NextTableName = nextPartName
	rwMutex.Unlock()        //寫入解鎖
	var isNowPartDone bool  //資料庫內是否已存在本次時間表
	var isNextPartDone bool //資料庫內是否已存在下是時間表
	for _, value := range tableNameList {
		if value == NowTableName {
			isNowPartDone = true
		}
		if value == NextTableName {
			isNextPartDone = true
		}
	}

	if isNowPartDone && isNextPartDone {
		return nil
	}

	if isNowPartDone == false {
		sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%v` LIKE `%v`;", NowTableName, latestTableName)
		err := theDB.Exec(sql).Error
		if err != nil {
			return err
		}

	}

	if isNextPartDone == false {
		sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%v` LIKE `%v`;", NextTableName, latestTableName)
		err := theDB.Exec(sql).Error
		if err != nil {
			return err
		}
	}

	return
}

//開始執行定時器創建logger表,每日00：00執行
func startTableCreator(ctx context.Context) {

	for {
		now := time.Now()                                                                    //获取当前时间，放到now里面，要给next用
		next := now.Add(time.Hour * 24)                                                      //獲取明天時間(+24小時)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location()) //計算下一個執行時間點
		log.Printf("%v: 使用表[%v.%v],定時器程序開始將於下個時間執行任務 [%v] > 檢查並創建表[%v.%v]", ServiceName, Config.DBName, NowTableName, next, Config.DBName, NextTableName)
		timer := time.NewTimer(next.Sub(now)) //计算当前时间到凌晨的时间间隔，设置一个定时器
		<-timer.C                             //等待下一個執行時間到來後往下執行

		//以下为定时执行的操作
		creatTables()
	}

}
