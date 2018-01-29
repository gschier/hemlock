package hemlock

// Config contains the static configuration for a Hemlock application.  There is
// an optional 'Extra' field for storing configuration not directly related to
// the core Hemlock functionality.
type Config struct {
	Name               string
	Env                string
	URL                string
	TemplatesDirectory string
	PublicDirectory    string
	PublicPrefix       string
	Database           *DatabaseConfig
	HTTP               *HTTPConfig
	Extra              []interface{}
}

// HTTPConfig contains server configuration.
type HTTPConfig struct {
	Host string
	Port string
}

// DatabaseConfig contains database settings.
type DatabaseConfig struct {
	Default     string // 'postgres'
	Migrations  string // 'migrations'
	Connections []DatabaseConnectionConfig
}

// DatabaseConnectionConfig contains settings for connecting to DB instances.
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
