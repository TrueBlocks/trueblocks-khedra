package types

type Service struct {
	Name      string `koanf:"name" validate:"required,oneof=api scraper monitor ipfs"`
	Enabled   bool   `koanf:"enabled"`
	Port      int    `koanf:"port,omitempty" validate:"service_field"`
	Sleep     int    `koanf:"sleep,omitempty" validate:"service_field"`
	BatchSize int    `koanf:"batchSize,omitempty" yaml:"batchSize,omitempty" validate:"service_field"`
}

func NewService(serviceType string) Service {
	switch serviceType {
	case "scraper":
		return Service{
			Name:      "scraper",
			Enabled:   false,
			Sleep:     10,
			BatchSize: 500,
		}
	case "monitor":
		return Service{
			Name:      "monitor",
			Enabled:   false,
			Sleep:     12,
			BatchSize: 500,
		}
	case "api":
		return Service{
			Name:    "api",
			Enabled: false,
			Port:    8080,
		}
	case "ipfs":
		return Service{
			Name:    "ipfs",
			Enabled: false,
			Port:    5001,
		}
	}

	if true {
		// The `true` above avoids a linter warning...
		panic("Unknown service type: " + serviceType)
	}

	return Service{}
}
