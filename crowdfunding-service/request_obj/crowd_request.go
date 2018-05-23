package request_obj

import "time"

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
