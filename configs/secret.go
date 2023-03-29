package configs

type Secret struct {
	ApiToken    *string `json:"api-token" yaml:"api-token"`
	ZoneName    *string `json:"zone-name" yaml:"zone-name"`
	AccountName *string `json:"account-name" yaml:"account-name"`
}

func (s *Secret) GetApiToken() string {
	if s.ApiToken == nil {
		return ""
	}
	return *s.ApiToken
}

func (s *Secret) GetZoneName() string {
	if s.ZoneName == nil {
		return ""
	}
	return *s.ZoneName
}

func (s *Secret) GetAccountName() string {
	if s.AccountName == nil {
		return ""
	}
	return *s.AccountName
}
