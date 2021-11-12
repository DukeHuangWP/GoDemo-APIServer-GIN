package models

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const ServiceName = "GORM.DB"

var (
	gormDB   *gorm.DB     //gormDB連接服務
	DoSQLs   doSQLs       //model主要CRUD執行清單
	DBConfig *MySQLConfig //MySQL設定值
)

//MySQL設定值
type MySQLConfig struct {
	Account  string //DB帳號
	Password string //DB密碼
	DBHost   string //DB的IP：Port
	DBName   string //DB的database名稱

	TimeoutSec      int //DB預設連線timeout時間
	ReadTimeoutSec  int //DB預設讀取timeout時間
	WriteTimeoutSec int //DB預設寫入timeout時間

	MaxIdleConns       int //设置空闲连接池中连接的最大数量
	MaxOpenConns       int //设置打开数据库连接的最大数量。
	ConnMaxLifeTimeSec int // SetConnMaxLifetime 设置了连接可复用的最大时间。

	InsertBatchMaxNum int //使用Gorm批次Insert資料的最大寫入筆數(每批)

	DBCacheTimerSec int //db儲存快取重新讀取秒數
}

/*
初始化gormDB
@param conf: db連線設定值
@param timerSec: 定時器log寫入db間格時間,低於1秒將會由改使用預設值10秒
@return error : 錯誤提示
*/
func InitDatabase(conf *MySQLConfig) (err error) {

	err = conectMySQL(conf) //連線MySql
	if err != nil {
		return err
	}
	DBConfig = conf
	DoSQLs.DB = gormDB
	go StartBuildDBCache(uint32(conf.DBCacheTimerSec)) //啟動獲取db快取
	return nil
}

/*
連線MySQL建立連接池
@param conf: db連線設定值
@return error : 錯誤提示
*/
func conectMySQL(conf *MySQLConfig) (err error) {

	if conf.Account == "" || conf.DBHost == "" || conf.DBName == "" {
		return fmt.Errorf("MySQL連接參數不應為空[Account='%v',DBHost='%v',DBName='%v']", DBConfig.Account, DBConfig.DBHost, DBConfig.DBName)
	}

	DBConfig = conf //存入全域變數

	//初始值校正
	if DBConfig.TimeoutSec <= 0 {
		DBConfig.TimeoutSec = 10
	}

	if DBConfig.ReadTimeoutSec <= 0 {
		DBConfig.ReadTimeoutSec = 30
	}

	if DBConfig.WriteTimeoutSec <= 0 {
		DBConfig.WriteTimeoutSec = 60
	}

	if DBConfig.MaxIdleConns <= 0 {
		DBConfig.MaxIdleConns = 1
	}

	if DBConfig.MaxOpenConns <= 0 {
		DBConfig.MaxOpenConns = 2
	}

	if DBConfig.ConnMaxLifeTimeSec <= 0 {
		DBConfig.ConnMaxLifeTimeSec = 60
	}

	if DBConfig.InsertBatchMaxNum <= 1 {
		DBConfig.InsertBatchMaxNum = 1000
	}

	if DBConfig.DBCacheTimerSec <= 1 {
		DBConfig.DBCacheTimerSec = 10
	}

	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	gormDB, err = gorm.Open(mysql.New(mysql.Config{
		//iftadmin:qweRTY123@tcp(34.80.149.81:3306)/system_server?charset=utf8&parseTime=True&loc=Local&TimeoutSec=10s&readTimeoutSec=30s&writeTimeoutSec=60s
		DSN: fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8&parseTime=True&loc=Local&timeout=%vs&readTimeout=%vs&writeTimeout=%vs", // DSN data source name
			DBConfig.Account, DBConfig.Password, DBConfig.DBHost, DBConfig.DBName, DBConfig.TimeoutSec, DBConfig.ReadTimeoutSec, DBConfig.WriteTimeoutSec),

		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		PrepareStmt: true,
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Duration(2) * time.Second, // 慢 SQL 阈值
				LogLevel:      logger.Info,                    // Log level
				Colorful:      true,                           // 彩色打印
			},
		),
	})

	if err != nil {
		return err
	}

	sqlDB, err := gormDB.DB() //此函數為sql.DB限定使用
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(DBConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(DBConfig.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(DBConfig.ConnMaxLifeTimeSec) * time.Second)

	log.Printf("%v : ORM連接完成 > [%v:%v@tcp(%v)/%v?charset=utf8&parseTime=True&loc=Local&TimeoutSec=%vs&readTimeoutSec=%vs&writeTimeoutSec=%vs]\n", ServiceName,
		DBConfig.Account, "******", DBConfig.DBHost, DBConfig.DBName, DBConfig.TimeoutSec, DBConfig.ReadTimeoutSec, DBConfig.WriteTimeoutSec)

	return nil
}

/*
每隔'%v'秒對db進行ping連線測試
@param ctx: context.Context用於批次終止子函式所產生的gorotuine的開關, 若不使用值可為nil
@param timerSec:  計數器執行間隔秒數
*/
func StartKeepDBAlive(ctx context.Context, timerSec int) {

	if DBConfig == nil {
		log.Printf("%v : DB尚未進行初始化連線!\n", ServiceName)
		return
	}

	if gormDB == nil {
		log.Printf("%v : DB尚未進行初始化連線失敗 [Account='%v',DBHost='%v',DBName='%v']\n", ServiceName, DBConfig.Account, DBConfig.DBHost, DBConfig.DBName)
		return
	}

	var err error
	sqlDB, err := gormDB.DB() //此函數為sql.DB限定使用
	if err != nil {
		log.Printf("%v : DB尚未進行初始化連線時發生錯誤 [Account='%v',DBHost='%v',DBName='%v'] > %v\n", ServiceName, DBConfig.Account, DBConfig.DBHost, DBConfig.DBName, err)
		return
	}

	if timerSec < 1 {
		timerSec = 120
	}

	log.Printf("%v : DB監聽連線將於每隔'%v'秒進行ping連線測試！\n", ServiceName, timerSec)
	go func(ctx context.Context) {
		timer := time.NewTicker(time.Duration(timerSec) * time.Second)
		for {
			select {
			case <-timer.C:
				err = sqlDB.Ping()
				if err != nil {
					err = conectMySQL(DBConfig)
					if err != nil {
						log.Printf("%v : DB斷線後嘗試重新連線錯誤 [Account='%v',DBHost='%v',DBName='%v'] > %v\n", ServiceName, DBConfig.Account, DBConfig.DBHost, DBConfig.DBName, err)
					}
				}
			}
		}
	}(ctx)

}
