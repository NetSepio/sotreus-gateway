package v1

import (
	"github.com/NetSepio/sotreus-gateway/api/v1/authenticate"
	"github.com/NetSepio/sotreus-gateway/api/v1/deployer"
	"github.com/NetSepio/sotreus-gateway/api/v1/flowid"
	"github.com/gin-gonic/gin"
)

func ApplyRoutes(r *gin.RouterGroup) {
	v1 := r.Group("/v1.0")
	{
		flowid.ApplyRoutes(v1)
		authenticate.ApplyRoutes(v1)
		deployer.ApplyRoutes(v1)
	}
}
