package v1

import (
	"github.com/NetSepio/sotreus-gateway/api/v1/deployer"
	"github.com/gin-gonic/gin"
)

func ApplyRoutes(r *gin.RouterGroup) {
	v1 := r.Group("/v1.0")
	{
		deployer.ApplyRoutes(v1)
	}
}
