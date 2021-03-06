package util

import (
	"os/user"
	"path"

	// "log"

	"log"
	"time"
	// "path/filepath"
	"strconv"
	"strings"
)

var configCache map[string]string = make(map[string]string)
var ConfigPath string

func getConfigPath() string {
	if ConfigPath == "" {
		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		ConfigPath = path.Join(usr.HomeDir, ".vger/config.json")
	}

	return ConfigPath
}
func ReadAllConfigs() map[string]string {
	if err := ReadJson(getConfigPath(), &configCache); err != nil {
		log.Print(getConfigPath())
		log.Fatal(err)
	}
	return configCache
}

func ReadConfig(key string) string {
	return ReadAllConfigs()[key]
}

func ReadIntConfig(key string) (res int) {
	res, err := strconv.Atoi(ReadConfig(key))
	if err != nil {
		panic(err)
	}
	return
}

func ReadStringSliceConfig(key string) (res []string) {
	res = strings.Split(ReadAllConfigs()[key], ",")
	for i, s := range res {
		res[i] = strings.TrimSpace(s)
	}
	return
}

func ReadSecondsConfig(key string) time.Duration {
	res, err := strconv.Atoi(ReadConfig(key))
	if err != nil {
		panic(err)
	}
	return time.Duration(res) * time.Second
}

func ReadBoolConfig(key string) (res bool) {
	return ReadConfig(key) == "true"
}

func SaveConfig(name, value string) {
	ReadAllConfigs()
	configCache[name] = value

	WriteJson(getConfigPath(), configCache)
}

func ToggleBoolConfig(name string) bool {
	if ReadBoolConfig(name) {
		SaveConfig(name, "false")
		return false
	} else {
		SaveConfig(name, "true")
		return true
	}
}
