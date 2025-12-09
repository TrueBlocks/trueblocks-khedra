package types

import (
	"github.com/TrueBlocks/trueblocks-chifra/v6/pkg/logger"
)

type Service struct {
	Name      string `koanf:"name" json:"name" validate:"required,oneof=api scraper monitor ipfs"`
	Enabled   bool   `koanf:"enabled" json:"enabled"`
	Port      int    `koanf:"port,omitempty" yaml:"port,omitempty" json:"port,omitempty" validate:"service_field"`
	Sleep     int    `koanf:"sleep,omitempty" yaml:"sleep,omitempty" json:"sleep,omitempty" validate:"service_field"`
	BatchSize int    `koanf:"batchSize,omitempty" yaml:"batchSize,omitempty" json:"batchSize,omitempty" validate:"service_field"`
}

func NewService(serviceType string) Service {
	switch serviceType {
	case "scraper":
		return Service{
			Name:      "scraper",
			Enabled:   true,
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
			Enabled: true,
			Port:    8080,
		}
	case "ipfs":
		return Service{
			Name:    "ipfs",
			Enabled: true,
			Port:    5001,
		}
	default:
		logger.Panic("Unknown service type: " + serviceType)
	}
	return Service{}
}

func (s *Service) IsEnabled() bool {
	return s.Enabled
}
