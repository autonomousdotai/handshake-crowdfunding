package service

import (
	"../dao"
	"../utils/service"
)

var fileUploadService = service.GSService{}
// service
var crowdFundingDao = dao.CrowdFundingDao{}
var crowdFundingImageDao = dao.CrowdFundingImageDao{}
var crowdFundingShakedDao = dao.CrowdFundingShakedDao{}
var ethTxDao = dao.EthTxDao{}
// template
