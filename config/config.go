package config

// Configurations has all configurations in this app.
type Configurations struct {
	Database DatabaseConfigurations
}

// DatabaseConfigurations is configurations about db.
type DatabaseConfigurations struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}
