package service

import (
	"github.com/autonomousdotai/handshake-crowdfunding/crowdfunding-service/dao"
	"github.com/autonomousdotai/handshake-crowdfunding/crowdfunding-service/utils/service"
)

var fileUploadService = service.GSService{}
// service
var crowdFundingDao = dao.CrowdFundingDao{}
var crowdFundingImageDao = dao.CrowdFundingImageDao{}
var crowdFundingShakedDao = dao.CrowdFundingShakedDao{}
var ethTxDao = dao.EthTxDao{}
// template
