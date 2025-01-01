package config

type Service struct {
	Name       string `koanf:"name" validate:"required,oneof=api scraper monitor ipfs"`  // Must be non-empty
	Enabled    bool   `koanf:"enabled"`                                                  // Defaults to false if not specified
	Port       int    `koanf:"port,omitempty" validate:"opt_min=1024,opt_max=65535"`     // Must be between 1024 and 65535
	Sleep      int    `koanf:"sleep,omitempty"`                                          // Must be non-negative
	BatchSize  int    `koanf:"batch_size,omitempty" validate:"opt_min=50,opt_max=10000"` // Must be between 50 and 10000
	RetryCnt   int    `koanf:"retry_cnt,omitempty"`                                      // Must be at least 1
	RetryDelay int    `koanf:"retry_delay,omitempty"`                                    // Must be at least 1
}

func NewScraper() Service {
	return Service{
		Name:       "scraper",
		Enabled:    false,
		Sleep:      10,
		BatchSize:  500,
		RetryCnt:   3,
		RetryDelay: 3,
	}
}

func NewMonitor() Service {
	return Service{
		Name:       "monitor",
		Enabled:    false,
		Sleep:      12,
		BatchSize:  500,
		RetryCnt:   3,
		RetryDelay: 3,
	}
}

func NewApi() Service {
	return Service{
		Name:    "api",
		Enabled: false,
		Port:    8080,
	}
}

func NewIpfs() Service {
	return Service{
		Name:    "ipfs",
		Enabled: false,
		Port:    5001,
	}
}
