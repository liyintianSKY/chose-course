package config

type Config struct {
	Server struct {
		HttpListen string `mapstructure:"http_listen"`
	}
	Database struct {
		Host         string
		Port         int
		User         string
		Password     string
		DBName       string
		MaxOpenConns int `mapstructure:"max_open_conns"`
		MaxIdleConns int `mapstructure:"max_idle_conns"`
	}
	Redis struct {
		Addr     string
		Password string
		DB       int
	}
	NATS struct {
		URL string
	}
	Logging struct {
		Level string
		File  string
	}
	Metrics struct {
		PrometheusPort int `mapstructure:"prometheus_port"`
	}
}
