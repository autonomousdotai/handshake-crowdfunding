package models

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm"
	_ "encoding/gob"
	"time"
)

type CrowdFundingShaked struct {
	DateCreated    time.Time
	DateModified   time.Time
	ID             int64
	UserId         int64
	CrowdFundingId int64
	CrowdFunding   CrowdFunding
	Price          float64
	Quantity       int
	Amount         float64
	Status         int
}

func (CrowdFundingShaked) TableName() string {
	return "crowd_funding_shaked"
}
