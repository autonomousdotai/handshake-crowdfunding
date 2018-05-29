package service

import (
	"github.com/autonomousdotai/handshake-crowdfunding/crowdfunding-service/dao"
	"github.com/autonomousdotai/handshake-crowdfunding/crowdfunding-service/utils"
)

var fileUploadService = utils.GSService{}
// service
var crowdFundingDao = dao.CrowdFundingDao{}
var crowdFundingImageDao = dao.CrowdFundingImageDao{}
var crowdFundingShakedDao = dao.CrowdFundingShakedDao{}
var ethTxDao = dao.EthTxDao{}
// template
var netUtil = utils.NetUtil{}
