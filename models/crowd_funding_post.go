package models

import (
	"time"
	_ "time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type CrowdFundingPost struct {
	ID               int64
	CrowdFundingId   int64
	CrowdFunding     CrowdFunding
	DateCreated      time.Time
	DateModified     time.Time
	ShortDescription string
	Description      string
	Title            string
	UserId           int64
	User             User
	Status           int
}

func (CrowdFundingPost) TableName() string {
	return "crowd_funding_post"
}
