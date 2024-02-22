package deployer

import (
	"bytes"
	"encoding/json"
	"io"
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
	}
}

func Deploy(c *gin.Context) {
	db := dbconfig.GetDb()
	walletAddress := c.GetString(paseto.CTX_WALLET_ADDRES)

	var count int64
	err := db.Model(&models.Sotreus{}).Where("wallet_address = ?", walletAddress).Find(&models.Sotreus{}).Count(&count).Error
	if err != nil {
		logwrapper.Errorf("failed to fetch data from database: %s", err)
		httpo.NewErrorResponse(http.StatusInternalServerError, err.Error()).SendD(c)
		return
	}

	if count >= 1 {
		logwrapper.Error("Can't create more vpn instances, maximum 1 allowed")
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
	deployerRequest := DeployerCreateRequest{SotreusID: req.Name, WalletAddress: walletAddress, Region: regions.Regions[req.Region].Code}
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logwrapper.Errorf("failed to send request: %s", err)
		httpo.NewErrorResponse(http.StatusInternalServerError, err.Error()).SendD(c)
		return
	}

	response := new(SotreusResponse)

	if err := json.Unmarshal(body, response); err != nil {
		logwrapper.Errorf("failed to get response: %s", err)
		httpo.NewErrorResponse(http.StatusInternalServerError, "failed to create VPN").SendD(c)
		return
	}
	contract := models.Sotreus{
		Name:          response.Message.VpnID,
		WalletAddress: walletAddress,
		Region:        req.Region,
	}
	result := db.Create(&contract)
	if result.Error != nil {
		logwrapper.Errorf("failed to create db entry: %s", err)
		httpo.NewErrorResponse(http.StatusInternalServerError, "failed to create VPN").SendD(c)
		return
	}
	payload := DeployResponse{
		VpnID:             response.Message.VpnID,
		VpnEndpoint:       response.Message.VpnEndpoint,
		FirewallEndpoint:  response.Message.FirewallEndpoint,
		DashboardPassword: response.Message.DashboardPassword,
	}
	httpo.NewSuccessResponseP(200, "VPN deployment successful", payload).SendD(c)
}
