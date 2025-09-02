package config

import (
	"errors"

	"github.com/go-viper/mapstructure/v2"
	"github.com/lumosolutions/yerd/server/internal/constants"
	"github.com/lumosolutions/yerd/server/internal/utils"

	"github.com/spf13/viper"
)

type Config struct {
	Yerd  YerdConfig
	Php   PhpConfig
	Nginx NginxConfig
}

var appConfig *Config

func LoadConfig() error {
	path := getConfigPath()

	v := viper.New()
	v.SetConfigName(constants.ConfigName)
	v.AddConfigPath(path)
	v.SetConfigType("json")

	setDefaults(v)

	if err := v.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return err
		}
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return err
	}

	appConfig = &config

	return nil
}

func GetConfig() (*Config, error) {
	if appConfig != nil {
		return appConfig, nil
	}

	if err := LoadConfig(); err != nil {
		return nil, err
	}

	return appConfig, nil
}

func WriteConfig() error {
	path := getConfigPath()

	v := viper.New()
	v.SetConfigName(constants.ConfigName)
	v.AddConfigPath(path)
	v.SetConfigType("json")

	var configMap map[string]interface{}
	err := mapstructure.Decode(appConfig, &configMap)
	if err != nil {
		return err
	}

	err = v.MergeConfigMap(configMap)
	if err != nil {
		return err
	}

	err = v.WriteConfig()
	if err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			return v.SafeWriteConfig()
		}
		return err
	}

	return nil
}

func getConfigPath() string {
	path := constants.ConfigPath
	isElevated := utils.IsRunningElevated()

	if isElevated {
		path = constants.ConfigElevatedPath
	}

	return constants.ExpandPath(path)
}

func setDefaults(v *viper.Viper) {
	elevated := utils.IsRunningElevated()

	v.SetDefault("yerd.elevated", elevated)
	v.SetDefault("yerd.port", constants.YerdPort)
	v.SetDefault("nginx.installed", false)
	v.SetDefault("nginx.elevated", elevated)
	v.SetDefault("nginx.servicename", constants.NginxServiceName)
	v.SetDefault("php", make(map[string]interface{}))

	if elevated {
		v.SetDefault("yerd.path", constants.WorkingElevatedPath)
		v.SetDefault("nginx.httpport", constants.SiteElevatedHttpPort)
		v.SetDefault("nginx.httpsport", constants.SiteElevatedHttpsPort)
	} else {
		v.SetDefault("yerd.path", constants.WorkingPath)
		v.SetDefault("nginx.httpport", constants.SiteHttpPort)
		v.SetDefault("nginx.httpsport", constants.SiteHttpsPort)
	}
}
