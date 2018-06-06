package service

import (
	"github.com/ninjadotorg/handshake-crowdfunding/dao"
	"github.com/ninjadotorg/handshake-crowdfunding/utils"
)

var (
	fileUploadService = utils.GSService{}
	// service
	crowdFundingDao       = dao.CrowdFundingDao{}
	crowdFundingImageDao  = dao.CrowdFundingImageDao{}
	crowdFundingShakeDao  = dao.CrowdFundingShakeDao{}
	crowdFundingFaqDao    = dao.CrowdFundingFaqDao{}
	crowdFundingUpdateDao = dao.CrowdFundingUpdateDao{}
	// template
	netUtil = utils.NetUtil{}
)
