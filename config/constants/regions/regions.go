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
		"us02": {
			Name:       "US 2",
			Code:       "us02",
			ServerHttp: envconfig.EnvVars.VPN_DEPLOYER_US02,
		},
	}
}
