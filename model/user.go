package model

import (
	"time"
)

type User struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `json:"name" grom:"name"`
}

type Username struct {
	OldName string `json:"oldname" grom:"oldname"`
	NewName string `json:"newname" grom:"newname"`
}

type Page struct {
	PageNum  int `json:"page_num"  grom:"page_num"`
	PageSize int `json:"page_size" grom:"page_size"`
}

func init() {
	//util.Db.AutoMigrate(&User{})
}
