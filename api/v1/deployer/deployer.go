package deployer

import (
	"bytes"
	"encoding/json"
	"fmt"
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
		g.GET("", getMyDeployments)
		g.DELETE("", deleteDeployment)
	}
}

func Deploy(c *gin.Context) {
	db := dbconfig.GetDb()
	walletAddress := c.GetString(paseto.CTX_WALLET_ADDRESS)
	userId := c.GetString(paseto.CTX_USER_ID)
	fmt.Println(walletAddress)

	var count int64
	err := db.Model(&models.Sotreus{}).Where("wallet_address = ?", walletAddress).Find(&models.Sotreus{}).Count(&count).Error
	if err != nil {
		logwrapper.Errorf("failed to fetch data from database: %s", err)
		httpo.NewErrorResponse(http.StatusInternalServerError, err.Error()).SendD(c)
		return
	}
	if count >= 3 {
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
	deployerRequest := DeployerCreateRequest{
		SotreusID:     req.Name,
		WalletAddress: walletAddress,
		Password:      req.Password,
		Region:        regions.Regions[req.Region].Code,
		Firewall:      req.Firewall,
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
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logwrapper.Errorf("failed to send request: %s", err)
		httpo.NewErrorResponse(http.StatusInternalServerError, err.Error()).SendD(c)
		return
	}
	fmt.Println("res: ", string(body))

	var fwEndpoint string

	if req.Firewall == "adguard" {
		fwEndpoint = req.Name + "-firewall." + req.Region + ".sotreus.com"
	} else {
		fwEndpoint = req.Name + "-firewall." + req.Region + ".sotreus.com/admin"
	}
	// Create a new Sotreus instance
	instance := models.Sotreus{
		Name:             req.Name,
		WalletAddress:    walletAddress,
		Region:           req.Region,
		VpnEndpoint:      req.Name + "-vpn." + req.Region + ".sotreus.com",
		FirewallEndpoint: fwEndpoint,
		Password:         string(req.Password),
		Firewall:         string(req.Firewall),
		UserId:           userId,
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
	userId := c.GetString(paseto.CTX_USER_ID)
	var instances []models.Sotreus
	if err := db.Model(&models.Sotreus{}).Where("user_id = ?", userId).Find(&instances).Error; err != nil {
		logwrapper.Errorf("failed to fetch DB : %s", err)
		httpo.NewErrorResponse(http.StatusInternalServerError, "failed to fetch DB").SendD(c)
		return

	}
	httpo.NewSuccessResponseP(200, "VPN's fetched successfully", instances).SendD(c)
}

func deleteDeployment(c *gin.Context) {
	db := dbconfig.GetDb()
	userId := c.GetString(paseto.CTX_USER_ID)
	sotreusName := c.Query("id")

	var instance models.Sotreus
	err := db.Model(&models.Sotreus{}).Where("user_id = ? and name = ?", userId, sotreusName).First(&instance).Error
	if err != nil {
		logwrapper.Errorf("failed to fetch data from database: %s", err)
		httpo.NewErrorResponse(http.StatusInternalServerError, err.Error()).SendD(c)
		return
	}
	ServerLink := regions.Regions[instance.Region].ServerHttp
	deployerRequest := DeployerDeleteRequest{
		SotreusID: sotreusName,
	}
	reqBodyBytes, err := json.Marshal(deployerRequest)
	if err != nil {
		logwrapper.Errorf("failed to encode request: %s", err)
		httpo.NewErrorResponse(http.StatusInternalServerError, "failed to create VPN").SendD(c)
		return
	}
	deleteReq, err := http.NewRequest(http.MethodDelete, ServerLink+"/sotreus", bytes.NewReader(reqBodyBytes))
	if err != nil {
		logwrapper.Errorf("failed to send request: %s", err)
		httpo.NewErrorResponse(http.StatusInternalServerError, "failed to create VPN").SendD(c)
		return
	}
	client := &http.Client{}
	resp, err := client.Do(deleteReq)
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

	result := db.Delete(&instance)
	if result.Error != nil {
		logwrapper.Errorf("failed to create db entry: %s", err)
		httpo.NewErrorResponse(http.StatusInternalServerError, "failed to create VPN").SendD(c)
		return
	}

	httpo.NewSuccessResponse(200, "VPN deleted successfully").SendD(c)

}
