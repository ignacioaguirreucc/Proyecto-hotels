package config

import "time"

const (
	MySQLHost     = "mysql"
	MySQLPort     = "3306"
	MySQLDatabase = "users_api"
	MySQLUsername = "root"
	MySQLPassword = "root"

	CacheDuration = 30 * time.Second

	MemcachedHost = "localhost"
	MemcachedPort = "11211"

	JWTKey      = "ThisIsAnExampleJWTKey!"
	JWTDuration = 24 * time.Hour
)
