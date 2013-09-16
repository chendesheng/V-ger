package util

import (
	"strconv"
	// "time"
)

var configCache map[string]string

func ReadAllConfigs() map[string]string {
	if configCache == nil {
		configCache = make(map[string]string)
		// go func() {
		// 	time.Sleep(time.Second * 20)
		// 	configCache = nil
		// }()
	}

	if err := ReadJson("config.json", &configCache); err != nil {
		panic(err)
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

func ReadBoolConfig(key string) (res bool) {
	return ReadConfig(key) == "true"
}

func SaveConfig(name, value string) {
	ReadAllConfigs()
	configCache[name] = value

	WriteJson("config.json", configCache)
}
