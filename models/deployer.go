package models

type Sotreus struct {
	Name             string `gorm:"primary_key" json:"name"`
	WalletAddress    string `json:"walletAddress"`
	Region           string `json:"region"`
	VpnEndpoint      string `json:"VpnEndpoint"`
	FirewallEndpoint string `json:"firewallEndpoint"`
	Password         string `json:"password"`
	Firewall         string `json:"firewall"`
}
