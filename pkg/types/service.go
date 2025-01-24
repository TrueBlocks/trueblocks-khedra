package types

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/utils"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/validate"
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
		// We may get either a Config...
		if config, ok := fv.Root().(*Config); ok {
			services := config.Services
			for _, s := range services {
				if s.Name == serviceName {
					service = &s
					break
				}
			}
		} else {
			// ...or (when testing) a Service type
			var ok bool
			if service, ok = fv.Root().(*Service); !ok {
				return validate.Failed(fv, "invalid root type", fmt.Sprintf("%T", fv.Root()))
			}
		}

		if service == nil {
			return validate.Failed(fv, "service not found", serviceName)
		}

		if !service.Enabled {
			return validate.Passed(fv, "not-enabled", serviceName)
		}

		switch service.Name {
		case "api", "ipfs":
			if service.Port < 1024 || service.Port > 65535 {
				return validate.Failed(fv, "Port must be between 1024 and 65535 (inclusive)", fmt.Sprintf("Port=%d", service.Port))
			}
		case "scraper", "monitor":
			if service.Sleep <= 0 {
				return validate.Failed(fv, "Sleep must be a positive integer", fmt.Sprintf("Sleep=%d", service.Sleep))
			} else if service.BatchSize < 50 || service.BatchSize > 10000 {
				return validate.Failed(fv, "BatchSize must be between 50 and 10000 (inclusive)", fmt.Sprintf("Port=%d", service.Port))
			}
		default:
			return validate.Failed(fv, "unknown service name", serviceName)
		}

		return validate.Passed(fv, "valid", serviceName)
	}

	validate.RegisterValidator("service_field", serviceFieldValidator)
}
