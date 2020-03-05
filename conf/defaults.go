package conf

import (
	config "github.com/spf13/viper"
)

func init() {
	// Logger Defaults
	config.SetDefault("logger.level", "info")
	config.SetDefault("logger.encoding", "console")
	config.SetDefault("logger.color", true)
	config.SetDefault("logger.dev_mode", true)

	// Server Configuration
	config.SetDefault("server.host", "")
	config.SetDefault("server.port", "8080")
	config.SetDefault("server.tls", false)
	config.SetDefault("server.devcert", false)
	config.SetDefault("server.certfile", "server.crt")
	config.SetDefault("server.keyfile", "server.key")
	config.SetDefault("server.jwt.key", "myjwtsecret")
	config.SetDefault("server.jwt.token_age", 3600)
	config.SetDefault("server.log_requests", true)
	config.SetDefault("server.log_requests_body", false)
	config.SetDefault("server.log_disabled_http", []string{"/version"})
	config.SetDefault("server.profiler_enabled", false)
	config.SetDefault("server.profiler_path", "/debug")

	// Database Settings
	config.SetDefault("storage.type", "postgres")
	config.SetDefault("storage.username", "postgres")
	config.SetDefault("storage.password", "mysecretpassword")
	config.SetDefault("storage.host", "localhost")
	config.SetDefault("storage.port", 5432)
	config.SetDefault("storage.database", "gorestapi")
	config.SetDefault("storage.sslmode", "disable")
	config.SetDefault("storage.retries", 5)
	config.SetDefault("storage.sleep_between_retries", "7s")
	config.SetDefault("storage.max_connections", 80)
	config.SetDefault("storage.wipe_confirm", false)

}
