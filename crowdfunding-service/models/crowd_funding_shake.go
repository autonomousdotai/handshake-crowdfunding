package models

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm"
	_ "encoding/gob"
	"time"
)

type CrowdFundingShake struct {
	DateCreated    time.Time
	DateModified   time.Time
	ID             int64
	UserId         int64
	CrowdFundingId int64
	Price          float64
	Quantity       int
	Amount         float64
	Status         int
}

func (CrowdFundingShake) TableName() string {
	return "crowd_funding_shake"
}
