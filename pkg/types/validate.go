package types

import (
	"fmt"
	"net/url"
	"strings"
)

// Validate is the main entry point for configuration validation.
// It validates a Config, Chain, Service, General, or Logging object and returns any validation errors.
func Validate(input interface{}) error {
	switch v := input.(type) {
	case *Config:
		return v.Validate()
	case *Chain:
		return v.validate("")
	case *Service:
		return v.validate()
	case *General:
		return v.validate()
	case *Logging:
		return v.validate()
	default:
		return fmt.Errorf("Validate expects a *Config, *Chain, *Service, *General, or *Logging, got %T", input)
	}
}

// Validate validates the entire configuration object.
func (c *Config) Validate() error {
	var errs []string

	// Validate chains
	for name, chain := range c.Chains {
		if err := chain.validate(name); err != nil {
			errs = append(errs, err.Error())
		}
	}

	// Validate services
	for _, service := range c.Services {
		if err := service.validate(); err != nil {
			errs = append(errs, err.Error())
		}
	}

	// Validate general settings
	if err := c.General.validate(); err != nil {
		errs = append(errs, err.Error())
	}

	// Validate logging settings
	if err := c.Logging.validate(); err != nil {
		errs = append(errs, err.Error())
	}

	if len(errs) > 0 {
		return newValidationError(errs)
	}
	return nil
}

// validate validates a single chain within a config context.
func (ch Chain) validate(name string) error {
	var errs []string

	// ChainID must always be non-zero, even for disabled chains
	if ch.ChainID == 0 {
		errs = append(errs, fmt.Sprintf("Chain[%s].ChainID must be non-zero", name))
	}

	// If not enabled, skip other validations
	if !ch.Enabled {
		if len(errs) > 0 {
			return newValidationError(errs)
		}
		return nil
	}

	// Enabled chain validation
	if ch.Name == "" {
		errs = append(errs, fmt.Sprintf("Chain[%s].Name is required when enabled", name))
	}

	if len(ch.RPCs) == 0 {
		errs = append(errs, fmt.Sprintf("Chain[%s].RPCs cannot be empty when enabled", name))
	} else {
		// Validate each RPC is a valid URL
		for i, rpc := range ch.RPCs {
			if !isValidURL(rpc) {
				errs = append(errs, fmt.Sprintf("Chain[%s].RPCs[%d] is not a valid URL: %q", name, i, rpc))
			}
		}
	}

	if len(errs) > 0 {
		return newValidationError(errs)
	}
	return nil
}

// validate validates a single service within a config context.
func (s Service) validate() error {
	var errs []string

	// If not enabled, no field validation required
	if !s.Enabled {
		return nil
	}

	// Enabled service validation - type-specific
	switch s.Name {
	case "api", "ipfs":
		if s.Port < 1024 || s.Port > 65535 {
			errs = append(errs, fmt.Sprintf("Service[%s].Port must be between 1024 and 65535, got %d", s.Name, s.Port))
		}

	case "scraper":
		if s.Sleep <= 0 {
			errs = append(errs, fmt.Sprintf("Service[%s].Sleep must be positive, got %d", s.Name, s.Sleep))
		}
		if s.BatchSize < 50 || s.BatchSize > 10000 {
			errs = append(errs, fmt.Sprintf("Service[%s].BatchSize must be between 50 and 10000, got %d", s.Name, s.BatchSize))
		}

	case "monitor":
		if s.Sleep <= 0 {
			errs = append(errs, fmt.Sprintf("Service[%s].Sleep must be positive, got %d", s.Name, s.Sleep))
		}
		if s.BatchSize < 1 || s.BatchSize > 1000 {
			errs = append(errs, fmt.Sprintf("Service[%s].BatchSize must be between 1 and 1000, got %d", s.Name, s.BatchSize))
		}

	default:
		errs = append(errs, fmt.Sprintf("[service_field] FAILED for Service.Name unknown service name (got %s)", s.Name))
	}

	if len(errs) > 0 {
		return newValidationError(errs)
	}
	return nil
}

// validate validates a General configuration object.
func (g *General) validate() error {
	var errs []string

	if g.DataFolder == "" {
		errs = append(errs, "General.DataFolder is required")
	}

	// Check for invalid characters in dataFolder
	if strings.ContainsAny(g.DataFolder, "\x00") {
		errs = append(errs, "General.DataFolder contains invalid characters")
	}

	// Validate Strategy
	if g.Strategy == "" {
		errs = append(errs, "General.Strategy is required")
	} else if g.Strategy != "download" && g.Strategy != "scratch" {
		errs = append(errs, fmt.Sprintf("General.Strategy must be 'download' or 'scratch', got %q", g.Strategy))
	}

	// Validate Detail
	if g.Detail == "" {
		errs = append(errs, "General.Detail is required")
	} else if g.Detail != "index" && g.Detail != "bloom" {
		errs = append(errs, fmt.Sprintf("General.Detail must be 'index' or 'bloom', got %q", g.Detail))
	}

	if len(errs) > 0 {
		return newValidationError(errs)
	}
	return nil
}

// validate validates a Logging configuration object.
func (l *Logging) validate() error {
	var errs []string

	if l.Folder == "" {
		errs = append(errs, "Logging.Folder is required")
	}
	if l.Filename == "" {
		errs = append(errs, "Logging.Filename is required")
	} else if !strings.HasSuffix(l.Filename, ".log") {
		errs = append(errs, fmt.Sprintf("Logging.Filename must end with '.log', got %q", l.Filename))
	}

	// Validate Level
	if l.Level != "debug" && l.Level != "info" && l.Level != "warn" && l.Level != "error" {
		errs = append(errs, fmt.Sprintf("Logging.Level must be 'debug', 'info', 'warn', or 'error', got %q", l.Level))
	}

	// Validate MaxSize
	if l.MaxSize <= 0 {
		errs = append(errs, fmt.Sprintf("Logging.MaxSize must be greater than 0, got %d", l.MaxSize))
	}

	// Validate MaxBackups
	if l.MaxBackups < 0 {
		errs = append(errs, fmt.Sprintf("Logging.MaxBackups must be non-negative, got %d", l.MaxBackups))
	}

	// Validate MaxAge
	if l.MaxAge < 0 {
		errs = append(errs, fmt.Sprintf("Logging.MaxAge must be non-negative, got %d", l.MaxAge))
	}

	if len(errs) > 0 {
		return newValidationError(errs)
	}
	return nil
}

// Helper functions

func isValidURL(urlStr string) bool {
	u, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return false
	}
	return u.Scheme != "" && u.Host != ""
}

// validationError aggregates multiple validation errors.
type validationError struct {
	errors []string
}

func newValidationError(errors []string) error {
	if len(errors) == 0 {
		return nil
	}
	return &validationError{errors: errors}
}

func (ve *validationError) Error() string {
	var msg strings.Builder
	for i, err := range ve.errors {
		if i > 0 {
			msg.WriteString("\n")
		}
		msg.WriteString(fmt.Sprintf("Error %d: %s", i+1, err))
	}
	return msg.String()
}
