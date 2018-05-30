package setting

import (
	"os"
	"log"
	"encoding/json"
	"strconv"
)

type Configuration struct {
	ServicePort          int
	DB                   string
	CdnDomain            string
	CdnHttps             bool
	SolrServiceUrl       string
	DispatcherServiceUrl string
	StorageServiceUrl    string
}

func (configuration Configuration) String() string {
	temp, err := json.Marshal(configuration)
	if err != nil {
		return ""
	}
	return string(temp)
}

var configuration *Configuration = nil

func LoadConfig() (*Configuration, error) {
	configuration := Configuration{}

	configuration.ServicePort, _ = strconv.Atoi(os.Getenv("SERVICE_PORT"))
	configuration.DB = os.Getenv("DB")
	configuration.CdnDomain = os.Getenv("CDN_DOMAIN")
	configuration.CdnHttps, _ = strconv.ParseBool(os.Getenv("CDN_HTTPS"))
	configuration.DispatcherServiceUrl = os.Getenv("DISPATCHER_SERVICE_URL")
	configuration.StorageServiceUrl = os.Getenv("STORAGE_SERVICE_URL")

	return &configuration, nil
}

func CurrentConfig() *Configuration {
	if configuration == nil {
		configurationVal, err := LoadConfig()
		if err != nil {
			log.Print(err)
		}
		configuration = configurationVal
	}
	return configuration
}
