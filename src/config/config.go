package config

type Config struct {
	SessionKey      string `mapstructure:"SESSION_KEY" yaml:"SESSION_KEY"`
	ServerPort      int    `mapstructure:"SERVER_PORT" yaml:"SERVER_PORT"`
	ServerHostname  string `mapstructure:"SERVER_HOSTNAME" yaml:"SERVER_HOSTNAME"`
	OIDCServerPort  int    `mapstructure:"OIDC_SERVER_PORT" yaml:"OIDC_SERVER_PORT"`
	OIDCWellKnown   string `mapstructure:"OIDC_WELL_KNOWN" yaml:"OIDC_WELL_KNOWN"`
	ClientID        string `mapstructure:"CLIENT_ID" yaml:"CLIENT_ID"`
	ClientSecret    string `mapstructure:"CLIENT_SECRET" yaml:"CLIENT_SECRET"`
	Scopes          string `mapstructure:"SCOPES" yaml:"SCOPES"`
	MigaduAttribute string `mapstructure:"MIGADU_SCOPE" yaml:"MIGADU_SCOPE"`
	MigaduAPIuser   string `mapstructure:"MIGADU_API_USER" yaml:"MIGADU_API_USER"`
	MigaduAPIkey    string `mapstructure:"MIGADU_API_KEY" yaml:"MIGADU_API_KEY"`
	OrgName         string `mapstructure:"ORG_NAME" yaml:"ORG_NAME"`
}

var C Config
