package util

func ReadAllConfigs() map[string]string {
	c := make(map[string]string)

	if err := ReadJson("config.json", &c); err != nil {
		panic(err)
	}
	return c
}

func ReadConfig(key string) string {
	return ReadAllConfigs()[key]
}
