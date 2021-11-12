package loggerToDB

import (
	"context"
	"database/sql"
	"runtime"

	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	ServiceName    = "Logger-To-DB"
	TableNameTitle = "api_log"
	MaxQueueBuff   = 10000 //gortuine單次處理資料筆數上限,避免cpu突然標高
)

var (
	theDB         *gorm.DB //DB服務入口
	NowTableName  string   //logger寫入表名稱(現在)
	NextTableName string   //logger寫入表名稱(未來)

	Config         *MySQLConfig //MySQL設定值
	SwitchDubugLog bool         //開啟寫入DubugLog至db

	QueueTimerSec int           //定時器寫入log到db的間隔時間
	QueueChan     chan *gormRow //定時器寫入log通道
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
}

type Message struct {
	ServerName          string // 服務名稱
	ServerHost          string // 伺服器Host名稱
	ClientIP            string // 客戶端IP
	RequestURN          string // 客戶端請求URN
	RequestMethod       string // 客戶端請求Method
	RequestAccount      string // 客戶端請求body參數{"Account"} //空字串將會被轉成Null
	RequestDataDecrypt  string // 客戶端請求body參數{"Data"}(解密文本) //空字串將會被轉成Null
	RequestBody         string // 客戶端請求body(完整)
	APIStatus           int    // API處理狀態: -3系統錯誤,-2設計錯誤,-1未知錯誤,0成功完成,1成功(除理中斷或未達期望結果),2成功(請求格式不符)
	APIErrorMessage     string // API錯誤訊息相關紀錄 //空字串將會被轉成Null
	TaskNo              string // API所產生任務代碼 //空字串將會被轉成Null
	ResponseHttpStatus  int    // API回傳HttpCode,例如:404
	ResponseCosttime    int64  // API處理耗費時間(單位:ms毫秒)
	ResponseAccount     string // API回傳body參數{"Account"}(使用者帳號名稱)
	ResponseMsgCode     string // API回傳body參數{"MsgCode"}(回傳提示碼))
	ResponseMsg         string // API回傳body參數{"Msg"}(回傳提示訊息)
	ResponseDataDecrypt string // API回傳body參數{"Data"}(解密文本) //空字串將會被轉成Null
	ResponseDataEncrypt string // API回傳body參數{"Data"}(加密文本) //空字串將會被轉成Null
	LogTime             int64  // log紀錄時間
}

func (msg *Message) TransToGorm() *gormRow {

	isRequestAccountNull := (msg.RequestAccount != "") //Gorm的Null寫入跟判斷方式真的很奇怪...
	isRequestDataDecryptNull := (msg.RequestDataDecrypt != "")
	isAPIErrorMessageNull := (msg.APIErrorMessage != "")
	isTaskNoNull := (msg.TaskNo != "")
	isResponseDataDecryptNull := (msg.ResponseDataDecrypt != "")
	isResponseDataEncryptNull := (msg.ResponseDataEncrypt != "")

	debugMessage := &sql.NullString{}
	if SwitchDubugLog == true {
		var logText string
		for index := 2; index < 8; index++ { //紀錄最後5行原始碼執行順序
			pc, filename, line, ok := runtime.Caller(index)
			if !ok {
				break
			}
			logText = logText + fmt.Sprintf(" %v:%d (%v)\n", filename, line, runtime.FuncForPC(pc).Name())
		}
		debugMessage = &sql.NullString{String: logText, Valid: true}
	}

	return &gormRow{
		ServerName:          msg.ServerName,
		ServerHost:          msg.ServerHost,
		ClientIP:            msg.ClientIP,
		RequestURN:          msg.RequestURN,
		RequestMethod:       msg.RequestMethod,
		RequestAccount:      &sql.NullString{String: msg.RequestAccount, Valid: isRequestAccountNull},
		RequestDataDecrypt:  &sql.NullString{String: msg.RequestDataDecrypt, Valid: isRequestDataDecryptNull},
		RequestBody:         msg.RequestBody,
		APIStatus:           msg.APIStatus,
		APIErrorMessage:     &sql.NullString{String: msg.APIErrorMessage, Valid: isAPIErrorMessageNull},
		TaskNo:              &sql.NullString{String: msg.TaskNo, Valid: isTaskNoNull},
		ResponseHttpStatus:  msg.ResponseHttpStatus,
		ResponseCosttime:    msg.ResponseCosttime,
		ResponseAccount:     msg.ResponseAccount,
		ResponseMsgCode:     msg.ResponseMsgCode,
		ResponseMsg:         msg.ResponseMsg,
		ResponseDataDecrypt: &sql.NullString{String: msg.ResponseDataDecrypt, Valid: isResponseDataDecryptNull},
		ResponseDataEncrypt: &sql.NullString{String: msg.ResponseDataEncrypt, Valid: isResponseDataEncryptNull},
		DebugMessage:        debugMessage,
		LogTime:             msg.LogTime,
	}
}

//Null值無法完美被轉換....只好強轉
type gormRow struct {
	//Id           uint      `gorm:"column:id"` //暫時用不到
	ServerName          string          `gorm:"column:server_name"`           // 服務名稱
	ServerHost          string          `gorm:"column:server_host"`           // 伺服器Host名稱
	ClientIP            string          `gorm:"column:client_ip"`             // 客戶端IP
	RequestURN          string          `gorm:"column:request_urn"`           // 客戶端請求URN
	RequestMethod       string          `gorm:"column:request_method"`        // 客戶端請求Method
	RequestAccount      *sql.NullString `gorm:"column:request_account"`       // 客戶端請求body參數{"Account"} //准許Null
	RequestDataDecrypt  *sql.NullString `gorm:"column:request_data"`          // 客戶端請求body參數{"Data"}(解密文本) //准許Null
	RequestBody         string          `gorm:"column:request_body"`          // 客戶端請求body(完整)
	APIStatus           int             `gorm:"column:api_status"`            // API處理狀態: -3系統錯誤,-2設計錯誤,-1未知錯誤,0成功完成,1成功(除理中斷或未達期望結果),2成功(請求格式不符)
	APIErrorMessage     *sql.NullString `gorm:"column:api_error_message"`     // API錯誤訊息相關紀錄 //准許Null
	TaskNo              *sql.NullString `gorm:"column:task_no"`               // API所產生任務代碼 //准許Null
	ResponseHttpStatus  int             `gorm:"column:response_http_status"`  // API回傳HttpCode,例如:404
	ResponseCosttime    int64           `gorm:"column:response_costtime"`     // API處理耗費時間(單位:ms毫秒)
	ResponseAccount     string          `gorm:"column:response_account"`      // API回傳body參數{"Account"}(使用者帳號名稱)
	ResponseMsgCode     string          `gorm:"column:response_msg_code"`     // API回傳body參數{"MsgCode"}(回傳提示碼))
	ResponseMsg         string          `gorm:"column:response_msg"`          // API回傳body參數{"Msg"}(回傳提示訊息)
	ResponseDataDecrypt *sql.NullString `gorm:"column:response_data_decrypt"` // API回傳body參數{"Data"}(解密文本) //准許Null
	ResponseDataEncrypt *sql.NullString `gorm:"column:response_data_encrypt"` // API回傳body參數{"Data"}(加密文本)//准許Null
	DebugMessage        *sql.NullString `gorm:"column:debug_message"`         // API紀錄debug訊息
	LogTime             int64           `gorm:"column:log_time"`              // log紀錄時間

	// CreateTime time.Time `gorm:"column:create_time"` // 建立時間
	// UpdateTime time.Time `gorm:"column:update_time"` // 最後更新時間
}

func (_ gormRow) TableName() string {
	return NowTableName //設定Table實際名稱,注意:此值僅會被讀取一次,若要改變名稱需用theDB.Table("新名稱")...
}

/*
初始化Logger
@param conf: db連線設定值
@param timerSec: 定時器log寫入db間格時間,低於1秒將會由改使用預設值10秒
@return error : 錯誤提示
*/
func InitLogger(ctx context.Context, conf *MySQLConfig, timerSec int, cacheFilePath string) (err error) {

	QueueChan = make(chan *gormRow, MaxQueueBuff)
	err = conectMySQL(conf)
	if err != nil {
		return err
	}
	Config = conf

	if timerSec < 1 {
		QueueTimerSec = 10 //預設值10秒
	} else {
		QueueTimerSec = timerSec
	}

	//------------開始讀取尚未寫入db暫存檔案------------
	if cacheFilePath != "" {
		cacheFile, lastMessage, err := GetMessageFromFile(cacheFilePath)
		if err != nil {
			log.Printf("%v : 讀取'%v'內容失敗,請檢查剛檔案寫入全系是否正常？", ServiceName, cacheFilePath)
		} else if len(lastMessage) > 0 {
			for index := 0; index < len(lastMessage); index++ {
				AddMessage(lastMessage[index])
			} //將上次暫存未完成之log取出

			log.Printf("%v : 讀取尚未寫入db暫存檔案 : '%v' , 已準備寫入至db當中！", ServiceName, cacheFilePath)
		} else {
			log.Printf("%v : 讀取'%v'內容後並未發現有任何殘存log", ServiceName, cacheFilePath)
		}
		theCacheFile = cacheFile //將開啟檔案交給全域變數
	}
	//--------------------END------------------------

	startLogger(ctx, QueueTimerSec) //開啟logger gorotuine
	go startTableCreator(ctx)       //開啟定時器用於創建表格

	return nil
}

/*
新增log至queue序列中等待被寫入db
@param message: log資訊
@return error : 錯誤提示
*/
func AddMessage(message *Message) (err error) {
	if len(QueueChan) == MaxQueueBuff {
		err = fmt.Errorf("%v : len(QueueChan)==%v 已滿,為何這麼多佇列未成功消費,請檢查logger與DB連線是否正常？", ServiceName, MaxQueueBuff)
	}

	QueueChan <- message.TransToGorm()
	return err
}

/*
啟動logger計時器,注意一旦啟用後goroutine將無法被關閉
@param timerSec: 定時器log寫入db間格時間
*/
func startLogger(ctx context.Context, timerSec int) {
	go func(ctx context.Context) {
		msglist := []gormRow{} //將每筆資料拼接成slice
		timer := time.NewTicker(time.Duration(timerSec) * time.Second)

		//定時器執行策略：當QueueChan有訊息進來即將該筆資料拼接成slice,若未發現有消息進來則將佇列中所有資料寫入db後休息至下一輪
		for {
			select {
			case <-ctx.Done():
				timer = time.NewTicker(time.Duration(1)) //立即命令timer馬上觸發
				if len(msglist) > 0 || len(QueueChan) > 0 {
					continue //加入寫入速度
				} else {
					return //終止gorotine
				}
			case message := <-QueueChan:
				msglist = append(msglist, *message)
				if theCacheFile != nil {
					if err := AddMessageToFile(theCacheFile, *message); err != nil {
						log.Printf("%v : 寫入備份快取時發生錯誤 > %v", ServiceName, err)
					}
				}
			case <-timer.C:
				if msgCount := len(msglist); msgCount > 0 {
					nowTableName := NowTableName //可能存在Race問題,待觀察
					if result := theDB.Table(nowTableName).CreateInBatches(msglist, Config.InsertBatchMaxNum); result.Error != nil {
						log.Printf("%v : 寫入log時發生錯誤[ %v ], 遺失的寫入訊息為 > %v", ServiceName, result.Error, msglist)
					} else {
						if leftCount := len(QueueChan); leftCount > 0 {
							log.Printf("%v : [%v.%v] 成功寫入訊息 [%v筆] ,尚有[%v筆]訊息在佇列當中", Config.DBName, NowTableName, ServiceName, msgCount, leftCount)
						} else {
							log.Printf("%v : [%v.%v] 成功寫入訊息 [%v筆] ", Config.DBName, NowTableName, ServiceName, msgCount)
						}
					}

					if theCacheFile != nil {
						if err := TruncateFile(theCacheFile); err != nil { //確認寫入DB完成後,將暫存檔案內資料清除
							log.Printf("%v : 清除備份快取時發生錯誤 > %v", ServiceName, err)
						}
					}

					msglist = nil //紀錄完成後需移除slice

				}
			}
		}
	}(ctx)
}

/*
連線MySQL建立連接池
@param conf: db連線設定值
@return error : 錯誤提示
*/
func conectMySQL(conf *MySQLConfig) (err error) {

	if conf.Account == "" || conf.DBHost == "" || conf.DBName == "" {
		return fmt.Errorf("MySQL連接參數不應為空[Account='%v',DBHost='%v',DBName='%v']", Config.Account, Config.DBHost, Config.DBName)
	}

	Config = conf //存入全域變數

	if Config.TimeoutSec <= 0 {
		Config.TimeoutSec = 10
	}

	if Config.ReadTimeoutSec <= 0 {
		Config.ReadTimeoutSec = 30
	}

	if Config.WriteTimeoutSec <= 0 {
		Config.WriteTimeoutSec = 60
	}

	if Config.MaxIdleConns <= 0 {
		Config.MaxIdleConns = 1
	}

	if Config.MaxOpenConns <= 0 {
		Config.MaxOpenConns = 2
	}

	if Config.ConnMaxLifeTimeSec <= 0 {
		Config.ConnMaxLifeTimeSec = 60
	}

	if Config.InsertBatchMaxNum <= 1 {
		Config.InsertBatchMaxNum = 1000
	}

	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	theDB, err = gorm.Open(mysql.New(mysql.Config{
		//iftadmin:qweRTY123@tcp(34.80.149.81:3306)/system_server?charset=utf8&parseTime=True&loc=Local&TimeoutSec=10s&readTimeoutSec=30s&writeTimeoutSec=60s
		DSN: fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8&parseTime=True&loc=Local&timeout=%vs&readTimeout=%vs&writeTimeout=%vs", // DSN data source name
			Config.Account, Config.Password, Config.DBHost, Config.DBName, Config.TimeoutSec, Config.ReadTimeoutSec, Config.WriteTimeoutSec),

		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Duration(2) * time.Second, // 慢 SQL 阈值
				LogLevel:      logger.Silent,                  // Log level
				Colorful:      true,                           // 彩色打印
			},
		),
	})

	if err != nil {
		return err
	}

	sqlDB, err := theDB.DB() //此函數為sql.DB限定使用
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(Config.MaxIdleConns)                                       //最小閒置連接數
	sqlDB.SetMaxOpenConns(Config.MaxOpenConns)                                       //最大連接數
	sqlDB.SetConnMaxLifetime(time.Duration(Config.ConnMaxLifeTimeSec) * time.Second) //連接持最多存活時間

	log.Printf("%v : ORM連接完成 > [%v:%v@tcp(%v)/%v?charset=utf8&parseTime=True&loc=Local&TimeoutSec=%vs&readTimeoutSec=%vs&writeTimeoutSec=%vs]\n", ServiceName,
		Config.Account, "******", Config.DBHost, Config.DBName, Config.TimeoutSec, Config.ReadTimeoutSec, Config.WriteTimeoutSec)

	err = creatTables()
	if err != nil {
		log.Printf("%v : 初始化表 [%v.%v] 時發生問題 %v\n", ServiceName, Config.DBName, NowTableName, err)
		return err
	} else {
		log.Printf("%v : 初始化表 [%v.%v] 完成可以開始準備寫入訊息!\n", ServiceName, Config.DBName, NowTableName)
	}

	return nil
}

/*
每隔'%v'秒對db進行ping連線測試
@param ctx: context.Context用於批次終止子函式所產生的gorotuine的開關, 若不使用值可為nil
@param timerSec:  計數器執行間隔秒數
*/
func StartKeepDBAlive(ctx context.Context, timerSec int) {

	if Config == nil {
		log.Printf("%v : DB尚未進行初始化連線!\n", ServiceName)
		return
	}

	if theDB == nil {
		log.Printf("%v : DB尚未進行初始化連線失敗 [Account='%v',DBHost='%v',DBName='%v']\n", ServiceName, Config.Account, Config.DBHost, Config.DBName)
		return
	}

	var err error
	sqlDB, err := theDB.DB() //此函數為sql.DB限定使用
	if err != nil {
		log.Printf("%v : DB尚未進行初始化連線時發生錯誤 [Account='%v',DBHost='%v',DBName='%v'] > %v\n", ServiceName, Config.Account, Config.DBHost, Config.DBName, err)
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
					err = conectMySQL(Config)
					if err != nil {
						log.Printf("%v : DB斷線後嘗試重新連線錯誤 [Account='%v',DBHost='%v',DBName='%v'] > %v\n", ServiceName, Config.Account, Config.DBHost, Config.DBName, err)
					}
				}
			}
		}
	}(ctx)

}
