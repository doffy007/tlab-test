package config

import (
	"time"
)

type DefaultConfig struct {
	Apps      Apps      `mapstructure:"apps"`
	Server    Server    `mapstructure:"server"`
	Database  Database  `mapstructure:"database"`
	Snowflake Snowflake `mapstructure:"snowflake"`
}

type Apps struct {
	Name           string        `mapstructure:"name"`
	Version        string        `mapstructure:"version"`
	RequestTimeout time.Duration `mapstructure:"requestTimeout"`
}

type Server struct {
	Port    int `mapstructure:"port"`
	Zerolog int `mapstructure:"zerolog"`
}

type Datasource struct {
	Url               string        `mapstructure:"url"`
	Port              int           `mapstructure:"port"`
	DatabaseName      string        `mapstructure:"databaseName"`
	Username          string        `mapstructure:"username"`
	Password          string        `mapstructure:"password"`
	Schema            string        `mapstructure:"schema"`
	ConnectionTimeout time.Duration `mapstructure:"connectionTimeout"`
	MaxIdleConnection int           `mapstructure:"maxIdleConnection"`
	MaxOpenConnection int           `mapstructure:"maxOpenConnection"`
	DebugMode         bool          `mapstructure:"debugMode"`
	DSN               string        `mapstructure:"databaseDSN"`
	PingInterval      int64         `mapstructure:"pingInterval"`
}

type Database struct {
	Postgres Datasource `mapstructure:"postgres"`
}

type Snowflake struct {
	IP     string `mapstructure:"snowflake_ip"`
	NodeID int64  `mapstructute:"snowflake_node_id"`
}
