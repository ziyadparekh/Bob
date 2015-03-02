package builder

var Services = []string{
	"identity-service",
	"inventory-service",
	"ecom-service",
	"boss-service",
	"logistics-service",
	"message-center",
	"1stdibs-admin-v2",
}

var JenkinsServices = map[string]string{
	"identity":  Services[0],
	"inventory": Services[1],
	"ecom":      Services[2],
	"boss":      Services[3],
	"logistics": Services[4],
	"msgcenter": Services[5],
	"adminv2":   Services[6],
}

var Jobs = map[string]string{
	"boss-service":        "JAVA-Boss-Offline (Custom)",
	"solr-config-service": "JAVA-DBServiceConfig-Solr (Custom)",
	"solr-ecom":           "JAVA-DBServiceConfig-SolrEcom (Custom)",
	"ecom-offline":        "JAVA-Ecom-Offline (Custom)",
	"ecom-service":        "JAVA-EcomService (Custom)",
	"identity-service":    "JAVA-IdentityService (Custom)",
	"inventory-service":   "JAVA-InventoryService (Custom)",
	"logistics-service":   "JAVA-LogisticsService (Custom)",
	"message-center":      "JAVA-MessageCenterService(Custom)",
	"query-service":       "JAVA-QueryServiceQueries (Custom)",
	"solr-service":        "JAVA-SolrService (Custom)",
	"1stdibs-admin-v2":    "Admin-v2 Deploy ANY FRUIT SERVER",
}

var DibsyJenkins = "http://jenkins.1stdibs.com"

var JenkinsHostKey = "SERVER_HOSTNAME"
var JenkinsBranchKey = "BRANCH_NAME"
var JenkinsServerName = "deathstar.1stdibs.com"
