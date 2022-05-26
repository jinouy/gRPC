package util

import (
	"gorm.io/gorm"
)

const (
	user     = "postgres"
	password = "aa2122822"
	dbname   = "gorm_project"
	port     = 5432
)

var (
	Db  *gorm.DB
	err error
)

//func init() {
//
//	psqlInfo := fmt.Sprintf("user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
//		user, password, dbname, port)
//	Db, err = gorm.Open(postgres.New(postgres.Config{
//		DSN:                  psqlInfo,
//		PreferSimpleProtocol: true,
//	}), &gorm.Config{})
//	if err != nil {
//		panic(err.Error())
//	}
//
//}
