package webapp

import (
	"github.com/NetSepio/sotreus-gateway/api/middleware/auth/paseto"
	"github.com/TheLazarusNetwork/go-helpers/httpo"
	"github.com/gin-gonic/gin"
)

func ApplyRoutes(r *gin.RouterGroup) {
	g := r.Group("/webapp")
	{
		g.Use(paseto.PASETO(false))
		g.GET("/auth", WebappAuth)
	}
}

func WebappAuth(c *gin.Context) {
	walletAddress := c.GetString(paseto.CTX_WALLET_ADDRES)

	res := webappResponse{
		WalletAddress: walletAddress,
	}
	httpo.NewSuccessResponseP(200, "Created new charge", res).SendD(c)

}
