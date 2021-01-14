package config

type Config struct {
	ImagesPath   string `json:"images_path"`
	DatabasePath string `json:"database_path"`
	Port         string `json:"port"`
}

func New() *Config {
	return &Config{
		Port:         "9876",
		ImagesPath:   "./phi-store/images",
		DatabasePath: "./phi-store/database.sqlite3",
	}
}
