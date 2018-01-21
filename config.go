package hemlock

type Config struct {
	Name               string
	Env                string
	URL                string
	TemplatesDirectory string
	PublicDirectory    string
	Database           *DatabaseConfig
	HTTP               *HTTPConfig
	Extra              map[string]interface{}
}

func (conf *Config) GetExtraStr(key string) string {
	v, ok := conf.Extra[key]
	if !ok {
		return ""
	}
	return v.(string)
}

type HTTPConfig struct {
	Host string
	Port string
}

type DatabaseConfig struct {
	Default     string // 'postgres'
	Migrations  string // 'migrations'
	Connections []DatabaseConnectionConfig
}

type DatabaseConnectionConfig struct {
	Driver    string
	Host      string
	Database  string
	Username  string
	Password  string
	Charset   string
	Collation string
	Prefix    string
	Schema    string
	SSLMode   bool
}
