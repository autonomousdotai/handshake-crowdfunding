package models

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "time"
	"time"
)

type CrowdFundingFaq struct {
	ID             int64
	DateCreated    time.Time
	DateModified   time.Time
	UserId         int64
	CrowdFundingId int64
	Question       string
	Answer         string
	Status         int
	User           User
}

func (CrowdFundingFaq) TableName() string {
	return "crowd_funding_faq"
}
