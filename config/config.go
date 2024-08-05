package config

const (
	PathEnv     = "GINSON_CONFIG"
	DefaultPath = "./config/config.yaml"
)

var (
	Config config
)

type config struct {
	Service Service `yaml:"service"`
	Log     Log     `yaml:"log"`
	DB      DB      `yaml:"db"`
	Cache   Cache   `yaml:"cache"`
	Auth    Auth    `yaml:"auth"`
}
type Service struct {
	Name        string `yaml:"name"`
	Mode        string `yaml:"mode"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
	Host        string `yaml:"host"`
	Port        string `yaml:"port"`
	Scheme      string `yaml:"scheme"`
	Prefix      string `yaml:"prefix"`
}

type Log struct {
	Level      string `yaml:"level"`
	Path       string `yaml:"path"`
	MaxSize    int    `yaml:"maxSize"`
	MaxAge     int    `yaml:"maxAge"`
	MaxBackups int    `yaml:"maxBackups"`
}
type DB struct {
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Schema   string `yaml:"schema"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Conn     Conn   `yaml:"conn"`
}
type Conn struct {
	Min uint64 `yaml:"min"`
	Max uint64 `yaml:"max"`
}
type Cache struct {
	Type     string `yaml:"type"`
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type Auth struct {
	JWT JWT `yaml:"jwt"`
}
type JWT struct {
	SigningKey  string `yaml:"signingKey"`
	RefreshTime string `yaml:"refreshTime"`
	ExpiresTime string `yaml:"expiresTime"`
}
