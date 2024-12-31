package config

type Logging struct {
	Folder     string `koanf:"folder"`
	Filename   string `koanf:"filename"`
	MaxSizeMb  int    `koanf:"max_size_mb"`
	MaxBackups int    `koanf:"max_backups"`
	MaxAgeDays int    `koanf:"max_age_days"`
	Compress   bool   `koanf:"compress"`
}

func NewLogging() Logging {
	return Logging{
		Folder:     "~/.khedra/logs",
		Filename:   "khedra.log",
		MaxSizeMb:  10,
		MaxBackups: 3,
		MaxAgeDays: 10,
		Compress:   true,
	}
}
