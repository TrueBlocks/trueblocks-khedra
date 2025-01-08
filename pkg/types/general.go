package types

type General struct {
	DataFolder string `koanf:"dataFolder" yaml:"dataFolder" validate:"required,folder_exists"`
}

func NewGeneral() General {
	return General{
		DataFolder: "~/.khedra/data",
	}
}
