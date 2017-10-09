package cmd

import (
	"flag"
	"errors"
	"log"
	"gopkg.in/natefinch/lumberjack.v2"
)

func setupLogging(logFilePath string) {
	log.SetOutput(&lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    1,
		MaxBackups: 10,
		MaxAge:     5,
	})
}

func LoadConfigurationWithFlags() (cfg * Configuration, err error) {
	const defaultConfigFilePath = "/opt/maildock-1/etc/config.yaml"
	const usageConfigFilePath = "Path of configuration file."
	var configFilePath string
	var logFilePath string
	flag.StringVar(&configFilePath, "conf", defaultConfigFilePath, usageConfigFilePath)
	flag.StringVar(&configFilePath, "C", defaultConfigFilePath, usageConfigFilePath)
	flag.StringVar(&logFilePath, "log", "", "Path of log file")
	flag.Parse()
	if "" == configFilePath {
		return nil, errors.New("Path of configuration is required.")
	}
	if "" != logFilePath {
		setupLogging(logFilePath)
	}
	return LoadConfigurationFromFile(configFilePath)
}
