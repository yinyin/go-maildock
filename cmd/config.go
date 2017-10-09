package cmd

import (
	"github.com/yinyin/go-maildock/database"
	"github.com/yinyin/go-maildock/database/locator"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
)

type databaseSetup struct {
	DatabaseType string
	Config       database.Configuration
}

func (d *databaseSetup) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	var dbType = struct {
		Type string `yaml:"type"`
	}{}
	if err = unmarshal(&dbType); nil != err {
		return err
	}
	d.DatabaseType = dbType.Type
	if cfg, err := locator.NewEmptyConfiguration(d.DatabaseType); nil != err {
		return err
	} else if err = unmarshal(cfg); nil != err {
		return err
	} else {
		d.Config = cfg
	}
	return nil
}

type Configuration struct {
	SMTPListenOn []string `yaml:"smtp-listen"`
	HTTPListenOn string   `yaml:"http-listen"`
	HTTPContent  struct {
		Path    string `yaml:"path"`
		Prefix  string `yaml:"prefix"`
		ProxyTo string `yaml:"proxy-to"`
	} `yaml:"http-content"`
	Database  databaseSetup `yaml:"database"`
	PurgeDays int           `yaml:"purge-days"`
}

func makeDefaultConfiguration() (cfg *Configuration) {
	return &Configuration{
		SMTPListenOn: []string{":25"},
	}
}

func LoadConfigurationFromFile(filePath string) (cfg *Configuration, err error) {
	buf, err := ioutil.ReadFile(filePath)
	if nil != err {
		return nil, err
	}
	cfg = makeDefaultConfiguration()
	if err = yaml.Unmarshal(buf, cfg); nil != err {
		return nil, err
	}
	return cfg, nil
}
