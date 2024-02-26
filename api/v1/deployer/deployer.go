package deployer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/NetSepio/sotreus-gateway/api/middleware/auth/paseto"
	"github.com/NetSepio/sotreus-gateway/config/constants/regions"
	"github.com/NetSepio/sotreus-gateway/config/dbconfig"
	"github.com/NetSepio/sotreus-gateway/models"
	"github.com/NetSepio/sotreus-gateway/util/pkg/logwrapper"
	"github.com/TheLazarusNetwork/go-helpers/httpo"
	"github.com/gin-gonic/gin"
)

func ApplyRoutes(r *gin.RouterGroup) {
	g := r.Group("/vpn")
	{
		g.Use(paseto.PASETO(false))
		g.POST("", Deploy)
		g.GET("", getMyDeployments)
	}
}

func Deploy(c *gin.Context) {
	db := dbconfig.GetDb()
	walletAddress := c.GetString(paseto.CTX_WALLET_ADDRES)
	fmt.Println(walletAddress)

	var count int64
	err := db.Model(&models.Sotreus{}).Where("wallet_address = ?", walletAddress).Find(&models.Sotreus{}).Count(&count).Error
	if err != nil {
		logwrapper.Errorf("failed to fetch data from database: %s", err)
		httpo.NewErrorResponse(http.StatusInternalServerError, err.Error()).SendD(c)
		return
	}
	fmt.Println(count)
	if count >= 1 {
		httpo.NewErrorResponse(http.StatusBadRequest, "Can't create more vpn instances, maximum 1 allowed").SendD(c)
		return
	}

	var req DeployRequest
	err = c.BindJSON(&req)
	if err != nil {
		logwrapper.Errorf("failed to bind JSON: %s", err)
		httpo.NewErrorResponse(http.StatusInternalServerError, err.Error()).SendD(c)
		return
	}
	ServerLink := regions.Regions[req.Region].ServerHttp
	deployerRequest := DeployerCreateRequest{
		SotreusID:     req.Name,
		WalletAddress: walletAddress,
		Password:      req.Password,
		Region:        regions.Regions[req.Region].Code,
	}
	reqBodyBytes, err := json.Marshal(deployerRequest)
	if err != nil {
		logwrapper.Errorf("failed to encode request: %s", err)
		httpo.NewErrorResponse(http.StatusInternalServerError, "failed to create VPN").SendD(c)
		return
	}
	contractReq, err := http.NewRequest(http.MethodPost, ServerLink+"/sotreus", bytes.NewReader(reqBodyBytes))
	if err != nil {
		logwrapper.Errorf("failed to send request: %s", err)
		httpo.NewErrorResponse(http.StatusInternalServerError, "failed to create VPN").SendD(c)
		return
	}
	client := &http.Client{}
	resp, err := client.Do(contractReq)
	if err != nil {
		logwrapper.Errorf("failed to send request: %s", err)
		httpo.NewErrorResponse(http.StatusInternalServerError, "failed to create VPN").SendD(c)
		return
	}
	if resp.StatusCode != 200 {
		logwrapper.Errorf("Error in response: %s", err)
		httpo.NewErrorResponse(http.StatusInternalServerError, "Error in response").SendD(c)
		return
	}
	defer resp.Body.Close()

	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	logwrapper.Errorf("failed to send request: %s", err)
	// 	httpo.NewErrorResponse(http.StatusInternalServerError, err.Error()).SendD(c)
	// 	return
	// }

	// response := new(SotreusResponse)

	// if err := json.Unmarshal(body, response); err != nil {
	// 	logwrapper.Errorf("failed to get response: %s", err)
	// 	httpo.NewErrorResponse(http.StatusInternalServerError, "failed to create VPN").SendD(c)
	// 	return
	// }
	instance := models.Sotreus{
		Name:             req.Name,
		WalletAddress:    walletAddress,
		Region:           req.Region,
		VpnEndpoint:      req.Name + "-vpn." + req.Region + ".sotreus.com",
		FirewallEndpoint: req.Name + "-firewall." + req.Region + ".sotreus.com",
		Password:         req.Password,
	}
	result := db.Create(&instance)
	if result.Error != nil {
		logwrapper.Errorf("failed to create db entry: %s", err)
		httpo.NewErrorResponse(http.StatusInternalServerError, "failed to create VPN").SendD(c)
		return
	}
	payload := DeployResponse{
		VpnID: req.Name,
	}
	httpo.NewSuccessResponseP(200, "VPN deployment successful", payload).SendD(c)
}

func getMyDeployments(c *gin.Context) {
	db := dbconfig.GetDb()
	walletAddress := c.GetString(paseto.CTX_WALLET_ADDRES)
	var instances []models.Sotreus
	if err := db.Model(&models.Sotreus{}).Where("wallet_address = ?", walletAddress).Find(&instances).Error; err != nil {
		logwrapper.Errorf("failed to fetch DB : %s", err)
		httpo.NewErrorResponse(http.StatusInternalServerError, "failed to fetch DB").SendD(c)
		return

	}
	httpo.NewSuccessResponseP(200, "VPN's fetched successfully", instances).SendD(c)

}
