package hemlock

type Config struct {
	Env       string
	URL       string
	Providers Providers
	Database  *DatabaseConfig
	Server    *ServerConfig
}

type ServerConfig struct {
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
