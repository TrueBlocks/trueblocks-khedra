package config

type General struct {
	DataPath string `koanf:"data_dir" validate:"required"`
	LogLevel string `koanf:"log_level" validate:"oneof=debug info warn error"`
}

func NewGeneral() General {
	return General{
		DataPath: "~/.khedra/data",
		LogLevel: "info",
	}
}
