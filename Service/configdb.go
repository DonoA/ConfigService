package main

import (
	"maps"
	"slices"

	"github.com/pkg/errors"
)

type ConfigOverrides map[string]Override

type ConfigDb struct {
	Config ConfigDbConfig

	Configs   map[string]Config
	Overrides map[string]ConfigOverrides
}

type ConfigDbConfig struct {
	User     string
	Password string
	Database string
}

func GetConfigPathStr(config *ConfigPath) string {
	return config.Service + "/" +
		config.Name
}

func GetOverridePathStr(override *OverrideKey) string {
	return override.EntityType + "/" +
		override.EntityId
}

func (db *ConfigDb) GetConfigs() ([]Config, error) {
	return slices.Collect(maps.Values(db.Configs)), nil
}

func (db *ConfigDb) AddConfig(config *Config) error {
	strPath := GetConfigPathStr(&config.ConfigPath)
	db.Configs[strPath] = *config
	db.Overrides[strPath] = make(ConfigOverrides)

	return nil
}

func (db *ConfigDb) GetConfig(path *ConfigPath) (Config, error) {
	strPath := GetConfigPathStr(path)
	if value, found := db.Configs[strPath]; found {
		return value, nil
	} else {
		return Config{}, errors.New("Config not found")
	}
}

func (db *ConfigDb) DeleteConfig(path *ConfigPath) error {
	strPath := GetConfigPathStr(path)
	delete(db.Configs, strPath)
	return nil
}

func (db *ConfigDb) GetOverrides(config *ConfigPath) ([]Override, error) {
	configStr := GetConfigPathStr(config)
	configOverrides, found := db.Overrides[configStr]
	if !found {
		return []Override{}, errors.New("Config not found")
	}

	values := slices.Collect(maps.Values(configOverrides))
	return values, nil
}

func (db *ConfigDb) AddOverride(config *ConfigPath, override *Override) error {
	configStr := GetConfigPathStr(config)
	configOverrides, found := db.Overrides[configStr]
	if !found {
		return errors.New("Config not found")
	}

	overrideStr := GetOverridePathStr(&override.OverrideKey)
	configOverrides[overrideStr] = *override
	return nil
}

func (db *ConfigDb) GetOverride(config *ConfigPath, overrideKey *OverrideKey) (Override, bool, error) {
	configStr := GetConfigPathStr(config)
	configOverrides, found := db.Overrides[configStr]
	if !found {
		return Override{}, false, errors.New("Config not found")
	}

	overrideStr := GetOverridePathStr(overrideKey)
	override, found := configOverrides[overrideStr]
	if !found {
		return Override{}, false, nil
	}

	return override, true, nil
}

func (db *ConfigDb) DeleteOverride(config *ConfigPath, overrideKey *OverrideKey) error {
	configStr := GetConfigPathStr(config)
	configOverrides, found := db.Overrides[configStr]
	if !found {
		return nil
	}

	overrideStr := GetOverridePathStr(overrideKey)
	delete(configOverrides, overrideStr)
	return nil
}
