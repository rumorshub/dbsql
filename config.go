package dbsql

import "time"

type ChannelsConfig map[string]Config

type Config struct {
	DriverName      string         `mapstructure:"driver_name" json:"driver_name,omitempty" yaml:"driver_name,omitempty"`
	DataSourceName  string         `mapstructure:"dsn" json:"dsn,omitempty" yaml:"data_source_name,omitempty"`
	Ping            bool           `mapstructure:"ping" json:"ping,omitempty" yaml:"ping,omitempty"`
	MaxIdleConns    *int           `mapstructure:"max_idle_conns" json:"max_idle_conns,omitempty" yaml:"max_idle_conns,omitempty"`
	MaxOpenConns    *int           `mapstructure:"max_open_conns" json:"max_open_conns,omitempty" yaml:"max_open_conns,omitempty"`
	ConnMaxLifetime *time.Duration `mapstructure:"conn_max_lifetime" json:"conn_max_lifetime,omitempty" yaml:"conn_max_lifetime,omitempty"`
	ConnMaxIdleTime *time.Duration `mapstructure:"conn_max_idle_time" json:"conn_max_idle_time,omitempty" yaml:"conn_max_idle_time,omitempty"`
}
