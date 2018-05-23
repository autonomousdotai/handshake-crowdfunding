package utils

import (
	"../setting"
)

func CdnUrlFor(fileUrl string) string {
	if fileUrl == "" {
		return ""
	}
	configuration := setting.CurrentConfig()
	result := ""
	if configuration.CdnHttps == true {
		result += "https://"
	} else {
		result += "http://"
	}
	result += configuration.CdnDomain + "/" + fileUrl
	return result
}

func CdnUrlFor2(filePath string, fileUrl string) string {
	if fileUrl == "" {
		return ""
	}
	configuration := setting.CurrentConfig()
	result := ""
	if configuration.CdnHttps == true {
		result += "https://"
	} else {
		result += "http://"
	}
	result += configuration.CdnDomain + "/" + filePath + fileUrl
	return result
}
