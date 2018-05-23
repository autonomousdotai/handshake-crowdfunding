package models

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm"
	_ "encoding/gob"
	"time"
)

type CrowdFundingImage struct {
	DateCreated    time.Time
	DateModified   time.Time
	ID             int64
	CrowdFundingId int64
	CrowdFunding   CrowdFunding
	Image          string
	YoutubeUrl     string
}

func (CrowdFundingImage) TableName() string {
	return "crowd_funding_image"
}
