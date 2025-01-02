package config

type General struct {
	DataPath string `koanf:"data_dir" validate:"required"`
}

func NewGeneral() General {
	return General{
		DataPath: "~/.khedra/data",
	}
}
