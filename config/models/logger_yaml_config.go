package models

type LoggerYAMLConfig struct {
	LogLevel      int    `yaml:"log_level"`
	LogPath       string `yaml:"log_path"`
	LogName       string `yaml:"log_name"`
	LogMaxSize    int    `yaml:"log_max_size"`
	LogMaxBackups int    `yaml:"log_max_backups"`
	LogMaxAge     int    `yaml:"log_max_age"`
	LogCompress   bool   `yaml:"log_compress"`
}
