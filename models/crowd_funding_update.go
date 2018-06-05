package models

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "time"
	"time"
)

type CrowdFundingUpdate struct {
	ID               int
	CrowdFundingId   int
	CrowdFunding     CrowdFunding
	DateCreated      time.Time
	DateModified     time.Time
	ShortDescription string
	Description      string
	Title            string
	UserId           int
	User             User
}

func (CrowdFundingUpdate) TableName() string {
	return "crowd_funding_update"
}
