package models

import (
	"fmt"

	"gorm.io/gorm"
)

type doSQLs struct {
	DB *gorm.DB
}

//獲取所有Test
func (doSQL *doSQLs) GetTestCasesUpdateTime() (updateTime int64) {

	query := []struct {
		UpdateTime int64 `gorm:"column:update_time"`
	}{}

	sql := fmt.Sprintf(`SELECT UNIX_TIMESTAMP(UPDATE_TIME) AS update_time FROM information_schema.TABLES WHERE information_schema.TABLES.TABLE_SCHEMA = 'system_server' AND information_schema.TABLES.TABLE_NAME = 'test_case';
	`)
	doSQL.DB.Raw(sql).Take(&query)
	if len(query) > 0 {
		return query[0].UpdateTime
	} else {
		return 0
	}

}

//查看該裝置deviceUid資訊
func (doSQL *doSQLs) GetFromDeviceInfo(selectCols []string, deviceUid string) (query map[string]interface{}, err error) {
	query = make(map[string]interface{})
	result := doSQL.DB.Table(Table_ServerInfo{}.TableName()).Select(selectCols).Where("device_uid", deviceUid).Take(&query)
	//例如: SELECT deviceinfo,status FROM `server_uid` WHERE `device_uid`='123465' LIMIT 1

	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return query, result.Error
}

//獲取ServerUID(名稱)
func (doSQL *doSQLs) GetServerUID() (serverUid string) {

	if value := GetDBCacheValue("server_uid"); value != nil {
		return fmt.Sprint(value)
	} //確認並讀取快取資料

	queryUID := []struct {
		UID string `gorm:"column:server_uid"` //'ip',
	}{}

	doSQL.DB.Table(Table_ServerInfo{}.TableName()).Distinct("server_uid").Where("type = ? AND status = ?", 1, 1).Find(&queryUID)
	// SELECT DISTINCT `ip` FROM `server_uid` WHERE `type` = 1 AND status = 1; //API設定
	if len(queryUID) <= 0 { //server設置或未online
		doSQL.DB.Table(Table_ServerInfo{}.TableName()).Distinct("server_uid").Where("type = ? AND status = ?", 0, 1).Find(&queryUID)
		// SELECT DISTINCT `ip` FROM `server_uid` WHERE `type` = 0 AND status = 1; //通用設定
		if len(queryUID) <= 0 { //server並未online
			return ""
		}
	}
	if len(queryUID) <= 0 { //server並未online
		return ""
	} // SELECT DISTINCT `server_uid` FROM `server_info` WHERE status = '1';

	serverUid = queryUID[0].UID
	SetDBCacheValue("server_uid", serverUid) //存入快取
	return serverUid                         //僅擷取第一個,若未來有需要在行擴充
}

//獲取所有Test
func (doSQL *doSQLs) GetAllTestCases() (query [][9]string) {

	queryCache := []struct {
		TestName          string `gorm:"column:test_name"`
		TestTag           string `gorm:"column:test_tag"`
		TestType          string `gorm:"column:test_type"`
		TestScript        string `gorm:"column:test_script"`
		TestScriptDefault string `gorm:"column:test_script_default"`
		TestScriptExample string `gorm:"column:test_script_example"`
		CrontabSpec       string `gorm:"column:crontab"`
		RunAtOnce         string `gorm:"column:run_at_once"`
		RunAtATime        string `gorm:"column:run_at_atime"`
	}{}

	sql := fmt.Sprintf(`SELECT test_name, test_tag, test_type, test_script,crontab, run_at_once ,run_at_atime,test_script_default,test_script_example FROM test_case WHERE enable = 1;`)
	doSQL.DB.Raw(sql).Take(&queryCache)
	for index := 0; index < len(queryCache); index++ {
		query = append(query, [9]string{
			queryCache[index].TestName,
			queryCache[index].TestTag,
			queryCache[index].TestType,
			queryCache[index].TestScript,
			queryCache[index].CrontabSpec,
			queryCache[index].RunAtOnce,
			queryCache[index].RunAtATime,
			queryCache[index].TestScriptDefault,
			queryCache[index].TestScriptExample,
		})
	}

	return
}

//獲取所有Test
func (doSQL *doSQLs) UpdateTestCases(testName, testScript string) (err error) {
	sql := fmt.Sprintf(`UPDATE test_case SET test_script = '%v' WHERE test_name = '%v';`,
		testScript, testName)
	result := doSQL.DB.Exec(sql)
	return result.Error
}
