package models

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm"
	_ "encoding/gob"
	"time"
)

type CrowdFunding struct {
	DateCreated        time.Time
	DateModified       time.Time
	ID                 int64
	Hid                int64
	UserId             int64
	Name               string
	Description        string
	ShortDescription   string
	CrowdDate          time.Time
	DeliverDate        time.Time
	Price              float64
	Goal               float64
	Balance            float64
	ShakeNum           int
	Status             int
	CrowdFundingImages []CrowdFundingImage
	CrowdFundingShake  CrowdFundingShake
}

func (CrowdFunding) TableName() string {
	return "crowd_funding"
}
