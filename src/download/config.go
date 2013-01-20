package download


type Config struct {
	BaseDir string
}

func readConfig() Config {
	config := Config{}
	readJson("config.json", &config)
	return config
}
