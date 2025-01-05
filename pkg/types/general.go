package types

type General struct {
	DataDir string `koanf:"data_dir" yaml:"data_dir" validate:"required,folder_exists"`
}

func NewGeneral() General {
	return General{
		DataDir: "~/.khedra/data",
	}
}
