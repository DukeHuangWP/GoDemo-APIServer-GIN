package models

import (
	"time"
)

// func getFieldName(tag string, field interface{}) (fieldname string) {
// 	rtType := reflect.TypeOf(field)
// 	for index := 0; index < rtType.NumField(); index++ {
// 		rtField := rtType.Field(index)
// 		if rtField.Tag.Get("gorm") == tag {
// 			return rtField.Name
// 		}
// 	}
// 	return ""
// }

type Table_ServerInfo struct {
	//Id           uint64      `gorm:"column:id"` //暫時用不到
	ServerUid  string    `gorm:"server_uid"`  //'server代号',
	IP         string    `gorm:"ip"`          //'ip',
	Type       int8      `gorm:"type"`        //' 0: 通用',
	Attach     string    `gorm:"attach"`      //'备注',
	Status     int8      `gorm:"status"`      //'0:停用 1:啟用 -1刪除',
	DeleteTime time.Time `gorm:"delete_time"` //'刪除時間',
	CreateTime time.Time `gorm:"create_time"` //'建立時間',
	UpdateTime time.Time `gorm:"update_time"` //'最後更新時間',
}

func (_ Table_ServerInfo) TableName() string {
	return "server_info" //設定Table實際名稱
}
