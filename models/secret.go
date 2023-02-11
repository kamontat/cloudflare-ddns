package models

type Secrets struct {
	ApiToken string `json:"api-token" yaml:"api-token"`
	ZoneName string `json:"zone-name" yaml:"zone-name"`
}
