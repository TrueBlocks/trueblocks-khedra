package config

type General struct {
	DataDir string `koanf:"data_dir" yaml:"data_dir" validate:"required"`
}

func NewGeneral() General {
	return General{
		DataDir: "~/.khedra/data",
	}
}
