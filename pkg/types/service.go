package types

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/utils"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/validate"
)

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
	default:
		panic("Unknown service type: " + serviceType)
	}
}

func (s *Service) IsEnabled() bool {
	return s.Enabled
}

func init() {
	serviceFieldValidator := func(fv validate.FieldValidator) error {
		serviceName := utils.RemoveAny(fv.Context(), "[]") // Assuming the context has the service name

		var service *Service
		if config, ok := fv.Root().(*Config); ok {
			services := config.Services
			for _, s := range services {
				if s.Name == serviceName {
					service = &s
					break
				}
			}
		} else {
			var ok bool
			if service, ok = fv.Root().(*Service); !ok {
				return validate.Failed(fv, "service_field", "invalid root type", fmt.Sprintf("%T", fv.Root()))
			}
		}

		if service == nil {
			return validate.Failed(fv, "service_field", "service not found", serviceName)
		}

		// Check if the service is enabled
		if !service.Enabled {
			return validate.Passed(fv, "service_field", "not-enabled", serviceName)
		}

		// Validate fields based on service name
		switch service.Name {
		case "api", "ipfs":
			if service.Port < 1024 || service.Port > 65535 {
				return validate.Failed(fv, "service_field", "Port must be between 1024 and 65535 (inclusive)", fmt.Sprintf("Port=%d", service.Port))
			}
		case "scraper", "monitor":
			if service.Sleep <= 0 {
				return validate.Failed(fv, "service_field", "Sleep must be a positive integer", fmt.Sprintf("Sleep=%d", service.Sleep))
			}
			if service.BatchSize < 50 || service.BatchSize > 10000 {
				return validate.Failed(fv, "service_field", "BatchSize must be between 50 and 10000 (inclusive)", fmt.Sprintf("Port=%d", service.Port))
			}
		default:
			return validate.Failed(fv, "service_field", "unknown service name", serviceName)
		}

		return validate.Passed(fv, "service_field", "valid", serviceName)
	}

	validate.RegisterValidator("service_field", serviceFieldValidator)
}
