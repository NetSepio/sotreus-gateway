package regions

import "github.com/NetSepio/sotreus-gateway/config/envconfig"

type Region struct {
	Name       string
	Code       string
	ServerHttp string
}

var Regions map[string]Region
var ErebrusRegions map[string]Region

//todo
//Store regions in persistent DB
//

func InitRegions() {
	Regions = map[string]Region{
		"us-east-2": {
			Name:       "US east 2",
			Code:       "us-east-2",
			ServerHttp: envconfig.EnvVars.VPN_DEPLOYER_API_US_EAST,
		},
		"ap-southeast-1": {
			Name:       "Asia Pacific",
			Code:       "ap-southeast-1",
			ServerHttp: envconfig.EnvVars.VPN_DEPLOYER_API_SG,
		},
		"us": {
			Name:       "US",
			Code:       "us",
			ServerHttp: envconfig.EnvVars.SOTREUS_US,
		},
		"sg": {
			Name:       "Singapore",
			Code:       "sg",
			ServerHttp: envconfig.EnvVars.SOTREUS_SG,
		},
	}
}
