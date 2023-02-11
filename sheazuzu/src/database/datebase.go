package database

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"sheazuzu/sheazuzu/src/entity"
)

func InitDB(conString string) *gorm.DB {
	db, err := gorm.Open("mysql", conString)
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(
		&entity.MatchData{},
		&entity.AdditionalInformation{},
	)
	return db
}

// "root:455279980@/atb?charset=utf8&parseTime=True&loc=Local"
// PATH="$PATH":/usr/local/mysql/bin
// mysql -u root -p
// 455279980
