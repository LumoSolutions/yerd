package config

type CertInfo struct {
	Name          string `json:"name"`
	Type          string `json:"type"`
	CertPath      string `json:"cert_path"`
	KeyPath       string `json:"key_path"`
	NssRegistered bool   `json:"nss_registered"`
	CaTrusted     bool   `json:"ca_trusted"`
}
