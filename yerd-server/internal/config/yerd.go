package config

type YerdConfig struct {
	YerdPort  int `json:"yerd_port"`
	HttpPort  int `json:"http_port"`
	HttpsPort int `json:"https_port"`
}
