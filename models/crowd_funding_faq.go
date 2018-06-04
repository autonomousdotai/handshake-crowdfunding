package models

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "time"
	"time"
)

type CrowdFundingFaq struct {
	ID             int
	CrowdFundingId int
	CrowdFunding   CrowdFunding
	DateCreated    time.Time
	DateModified   time.Time
	Question       string
	Answer         string
	UserId         int
}

func (CrowdFundingFaq) TableName() string {
	return "crowd_funding_faq"
}
