package configs

import (
	"os"
	"strconv"
	"encoding/json"
	"log"
)

var ServicePort, _ = strconv.Atoi(os.Getenv("SERVICE_PORT"))
var DB = os.Getenv("DB")
var CdnDomain = os.Getenv("CDN_DOMAIN")
var CdnHttps, _ = strconv.ParseBool(os.Getenv("CDN_HTTPS"))
var DispatcherServiceUrl = os.Getenv("DISPATCHER_SERVICE_URL")
var StorageServiceUrl = os.Getenv("STORAGE_SERVICE_URL")
var SolrServiceUrl = os.Getenv("SOLR_SERVICE_URL")

var AppConf = AppConfig{}

type AppConfig struct {
	ServicePort          int        `json:"service_port"`
	DbUrl                string     `json:"db_url"`
	CdnUrl               string     `json:"cdn_url"`
	DispatcherServiceUrl string     `json:"dispatcher_service_url"`
	StorageServiceUrl    string     `json:"storage_service_url"`
	SolrServiceUrl       string     `json:"solr_service_url"`
	PubsubConf           PubsubConf `json:"pubsub_conf"`
}

type PubsubConf struct {
	CredsFile string `json:"creds_file"`
	ProjectId string `json:"project_id"`
	SubName   string `json:"sub_name"`
}

func Initialize(confFile string) {
	file, err := os.Open(confFile)
	if err != nil {
		log.Println(err)
	}
	decoder := json.NewDecoder(file)
	conf := AppConfig{}
	err = decoder.Decode(&conf)
	if err != nil {
		log.Println(err)
	}
	AppConf = conf
}
