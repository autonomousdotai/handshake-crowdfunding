package models

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm"
	_ "encoding/gob"
)

type JsonUserResponse struct {
	Status  int    `json:"status"`
	Data    User   `json:"data"`
}

type User struct {
	ID     int64  `json:"id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Status int    `json:"status"`
}
