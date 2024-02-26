package deployer

import "gorm.io/gorm"

type DeployRequest struct {
	Name     string `json:"name" binding:"required"`
	Region   string `json:"region" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type DeployerCreateRequest struct {
	SotreusID     string `json:"sotreusID" binding:"required"`
	WalletAddress string `json:"walletAddress" binding:"required"`
	Region        string `json:"region" binding:"required"`
	Email         string `json:"email,omitempty"`
	Password      string `json:"password,omitempty"`
}
type SotreusRequest struct {
	VpnId string `json:"vpnId" binding:"required"`
}
type SotreusResponse struct {
	Todo    string `json:"todo"`
	Result  string `json:"result"`
	Message struct {
		VpnID             string `json:"vpn_id"`
		VpnEndpoint       string `json:"vpn_endpoint"`
		VpnAPIPort        int    `json:"vpn_api_port"`
		VpnExternalPort   int    `json:"vpn_external_port"`
		FirewallEndpoint  string `json:"firewall_endpoint"`
		DashboardPassword string `json:"dashboard_password"`
	} `json:"message"`
}

type DeployResponse struct {
	VpnID             string `json:"vpn_id"`
	VpnEndpoint       string `json:"vpn_endpoint"`
	FirewallEndpoint  string `json:"firewall_endpoint"`
	DashboardPassword string `json:"dashboard_password"`
}

type GetDeployments struct {
	Message string     `json:"message"`
	Data    []Instance `json:"data"`
}

type Instance struct {
	gorm.Model
	VpnID             string `json:"vpn_id" gorm:""`
	VpnEndpoint       string `json:"vpn_endpoint"`
	VpnAPIPort        int    `json:"vpn_api_port"`
	VpnExternalPort   int    `json:"vpn_external_port"`
	FirewallEndpoint  string `json:"firewall_endpoint"`
	DashboardPassword string `json:"dashboard_password"`
	Status            string `json:"status"`
	WalletAddress     string `json:"walletAddress,omitempty"`
}

type DeployerDeleteRequest struct {
	SotreusID string `json:"sotreusID"`
}

type ServiceInfoSotreus struct {
	gorm.Model
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
	// Uuid      string                `json:"uuid,omitempty"`
	Category           string `json:"category,omitempty"`
	Status             string `json:"status,omitempty"`
	CreatedAt          int64  `json:"createdAt,omitempty"`
	UpdatedAt          int64  `json:"updatedAt,omitempty"`
	DeletedAt          int64  `json:"deletedAt,omitempty"`
	SotreusID          string `json:"sotreusID,omitempty"`
	PiholeID           string `json:"piholeID,omitempty"`
	SotreusContainerID string `json:"sotreusContainerID,omitempty"`
	PiholeContainerID  string `json:"piholeContainerID,omitempty"`
}
