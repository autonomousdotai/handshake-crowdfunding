package service

import (
	"../dao"
	"../utils/service"
)

var s3Service = service.S3Service{}
// service
var crowdFundingDao = dao.CrowdFundingDao{}
var crowdFundingImageDao = dao.CrowdFundingImageDao{}
var crowdFundingShakedDao = dao.CrowdFundingShakedDao{}
// template
