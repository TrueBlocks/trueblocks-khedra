package config

type Logging struct {
	Folder     string `koanf:"folder" validate:"required,dirpath"`
	Filename   string `koanf:"filename" validate:"required,endswith=.log"`
	MaxSizeMb  int    `koanf:"max_size_mb" validate:"required,min=5"`
	MaxBackups int    `koanf:"max_backups" validate:"required,min=1"`
	MaxAgeDays int    `koanf:"max_age_days" validate:"required,min=1"`
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
