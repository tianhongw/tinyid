package config

import (
	"errors"

	"github.com/spf13/viper"
)

var gConfig *Config

func GetConfig() *Config {
	return gConfig
}

type Config struct {
	Log *LogOptions `mapstructure:"log"`
	DB  *DBOptions  `mapstructure:"db"`
}

func (c *Config) Validate() error {
	if c.Log == nil {
		return errors.New("nil log config")
	}

	if c.DB == nil {
		return errors.New("nil database config")
	}

	if len(c.DB.Host) == 0 ||
		len(c.DB.Port) == 0 ||
		len(c.DB.Username) == 0 ||
		len(c.DB.Password) == 0 ||
		len(c.DB.Name) == 0 {

		return errors.New("invalid database config")
	}

	if len(c.DB.Host) != len(c.DB.Port) ||
		len(c.DB.Host) != len(c.DB.Username) ||
		len(c.DB.Host) != len(c.DB.Password) ||
		len(c.DB.Host) != len(c.DB.Name) {

		return errors.New("invalid database config")
	}

	return nil
}

type DBOptions struct {
	Host         []string `mapstructure:"host"`
	Port         []int    `mapstructure:"port"`
	Username     []string `mapstructure:"username"`
	Password     []string `mapstructure:"password"`
	Name         []string `mapstructure:"name"`
	MaxIdleConns int      `mapstructure:"maxIdleConns"`
	MaxOpenConns int      `mapstructure:"maxOpenConns"`
}

type LogOptions struct {
	Type         string   `mapstructure:"type"`
	Level        string   `mapstructure:"level"`
	Format       string   `mapstructure:"format"`
	Outputs      []string `mapstructure:"outputs"`
	ErrorOutputs []string `mapstructure:"errorOutputs"`
}

func Init(cfgFile, cfgType string) (string, error) {
	v := viper.New()

	v.SetConfigFile(cfgFile)
	v.SetConfigType(cfgType)

	if err := v.ReadInConfig(); err != nil {
		return "", err
	}

	c := new(Config)

	if err := v.Unmarshal(c); err != nil {
		return "", err
	}

	if err := c.Validate(); err != nil {
		return "", err
	}

	gConfig = c

	return v.ConfigFileUsed(), nil
}
