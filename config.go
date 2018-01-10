package begonia

type Config struct {
	*ApplicationConfig
	*DatabaseConfig
}

type ApplicationConfig struct {
	Env            string
	Debug          bool
	URL            string
	Timezone       string
	Locale         string
	Languages      []string
	FallbackLocale string
	Key            string
	Cipher         string
	Log            string
	Providers      Providers
	Aliases        map[string]interface{}
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
