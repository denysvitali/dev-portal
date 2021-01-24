package app

type OauthConfig struct {
	AuthURL      string `yaml:"auth_url"`
	TokenURL     string `yaml:"token_url"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`

	UserInfoURL string `yaml:"user_info_url"`
}

type Config struct {
	Oidc       OauthConfig `yaml:"oidc"`
	BaseUrl    string      `yaml:"base_url"`
	ListenAddr string      `yaml:"listen_addr" default:":8080"`
}
