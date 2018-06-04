package service

import (
	"github.com/autonomousdotai/handshake-crowdfunding/dao"
	"github.com/autonomousdotai/handshake-crowdfunding/utils"
)

var fileUploadService = utils.GSService{}
// service
var crowdFundingDao = dao.CrowdFundingDao{}
var crowdFundingImageDao = dao.CrowdFundingImageDao{}
var crowdFundingShakeDao = dao.CrowdFundingShakeDao{}
var ethTxDao = dao.EthTxDao{}
// template
var netUtil = utils.NetUtil{}
