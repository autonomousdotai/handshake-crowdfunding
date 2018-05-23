package response_obj

import (
	"time"
	"../models"
	"../utils"
)

type CrowdFundingResponse struct {
	ID               int64                       `json:"id"`
	UserId           int64                       `json:"user_id"`
	Name             string                      `json:"name"`
	Description      string                      `json:"description"`
	ShortDescription string                      `json:"short_description"`
	Image            string                      `json:"image"`
	CrowdDate        time.Time                   `json:"crowd_date"`
	DeliverDate      time.Time                   `json:"deliver_date"`
	Price            float64                     `json:"price"`
	Goal             float64                     `json:"goal"`
	Status           int                         `json:"status"`
	Images           []CrowdFundingImageResponse `json:"images"`
}

type CrowdFundingImageResponse struct {
	ID             int64                `json:"id"`
	CrowdFundingId int64                `json:"crowd_funding_id"`
	CrowdFunding   CrowdFundingResponse `json:"crowd_funding"`
	Image          string               `json:"name"`
}

type CrowdFundingShakedResponse struct {
	ID           int64                `json:"id"`
	UserId       int64                `json:"user_id"`
	Price        float64              `json:"price"`
	Quantity     int                  `json:"quantity"`
	Amount       float64              `json:"amount"`
	CrowdFunding CrowdFundingResponse `json:"crowd_funding"`
}

func MakeCrowdFundingResponse(model models.CrowdFunding) CrowdFundingResponse {
	result := CrowdFundingResponse{}
	result.ID = model.ID
	result.UserId = model.UserId
	result.Name = model.Name
	result.Description = model.Description
	result.ShortDescription = model.ShortDescription
	result.Image = utils.CdnUrlFor2("upload/images/crowd-funding/", model.Image)
	result.CrowdDate = model.CrowdDate
	result.DeliverDate = model.DeliverDate
	result.Price = model.Price
	result.Goal = model.Goal
	result.Status = model.Status
	result.Images = MakeArrayCrowdFundingImageResponse(model.CrowdFundingImages)
	return result
}

func MakeCrowdFundingImageResponse(model models.CrowdFundingImage) CrowdFundingImageResponse {
	result := CrowdFundingImageResponse{}
	result.ID = model.ID
	result.Image = model.Image
	return result
}

func MakeArrayCrowdFundingImageResponse(models []models.CrowdFundingImage) []CrowdFundingImageResponse {
	results := []CrowdFundingImageResponse{}
	for _, model := range models {
		result := MakeCrowdFundingImageResponse(model)
		results = append(results, result)
	}
	return results
}
