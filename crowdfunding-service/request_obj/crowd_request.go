package request_obj

import (
	"time"
)

type CrowdFundingRequest struct {
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	ShortDescription string    `json:"short_description"`
	YoutubeUrl       string    `json:"youtube_url"`
	CrowdDate        time.Time `json:"crowd_date"`
	DeliverDate      time.Time `json:"deliver_date"`
	Price            float64   `json:"price"`
	Goal             float64   `json:"goal"`
	Status           int       `json:"status"`
}

type CrowdFundingFaqRequest struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type CrowdFundingUpdateRequest struct {
	Title            string `json:"title"`
	ShortDescription string `json:"short_description"`
	Description      string `json:"description"`
}
