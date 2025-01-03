package config

import (
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/rpc"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	validateStrictURL := func(fl validator.FieldLevel) bool {
		rawURL := fl.Field().String()

		// Access the parent struct (Chain) to check the Enabled field
		parent := fl.Parent().Interface().(Chain)
		if !parent.Enabled {
			// Skip validation if the Chain is not enabled
			return true
		}

		u, err := url.Parse(rawURL)
		// Validate URL only if Chain is enabled
		return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != ""
	}

	// validatePingOne validates that at least one of the given URLs is reachable
	validatePingOne := func(fl validator.FieldLevel) bool {
		isTesting := os.Getenv("TEST_MODE") == "true"
		if isTesting {
			return true
		}

		url := fl.Field().String()
		if err := rpc.PingRpc(url); err == nil {
			return true
		}
		return false
	}

	// validateIsWritable validates that a path exists and is writable
	validateIsWritable := func(fl validator.FieldLevel) bool {
		path := fl.Field().String()
		path = expandPath(path)

		testFile := filepath.Join(path, ".writable_check")
		file, err := os.Create(testFile)
		if err != nil {
			return false
		}
		_ = file.Close()
		_ = os.Remove(testFile)

		return true
	}

	// validatePathExists validates that a path is non-empty and that the path exists
	validatePathExists := func(fl validator.FieldLevel) bool {
		path := fl.Field().String()
		path = expandPath(path)
		_, err := os.Stat(path)
		return err == nil
	}

	// validateOptMin validates that an integer is either zero (unset) or greater than or equal to a given value
	validateOptMin := func(fl validator.FieldLevel) bool {
		param := fl.Param()
		min, err := strconv.Atoi(param)
		if err != nil {
			return false
		}
		value := fl.Field().Int()
		if value == 0 {
			return true
		}
		return value >= int64(min)
	}

	// validateOptMax validates that an integer is either zero (unset) or less than or equal to a given value
	validateOptMax := func(fl validator.FieldLevel) bool {
		param := fl.Param()
		max, err := strconv.Atoi(param)
		if err != nil {
			return false
		}
		value := fl.Field().Int()
		if value == 0 {
			return true
		}
		return value <= int64(max)
	}

	// validateDirPath validates that a string is both non-empty and an existing folder
	validateDirPath := func(fl validator.FieldLevel) bool {
		path := fl.Field().String()
		if path == "" {
			return false
		}
		path = expandPath(path)
		info, err := os.Stat(path)
		return err == nil && info.IsDir()
	}

	// validateService validates the configuration of a service which depends on the service type
	validateService := func(sl validator.StructLevel) {
		service := sl.Current().Interface().(Service)

		switch service.Name {
		case "api":
			// For "api" services, `Port` is required.
			if service.Port == 0 {
				sl.ReportError(service.Port, "Port", "port", "required_api_port", "")
			}
		case "scraper", "monitor":
			// For "scraper" and "monitor" services, required fields must be non-zero.
			if service.Sleep <= 0 {
				sl.ReportError(service.Sleep, "Sleep", "sleep", "required_scraper_monitor_sleep", "")
			}
			if service.BatchSize < 50 || service.BatchSize > 10000 {
				sl.ReportError(service.BatchSize, "BatchSize", "batch_size", "invalid_scraper_monitor_batch_size", "")
			}
		case "ipfs":
			// For "ipfs" services, `Port` is required.
			if service.Port == 0 {
				sl.ReportError(service.Port, "Port", "port", "required_ipfs_port", "")
			}
		}
	}

	// validateServiceField validates fields based on the service type
	validateServiceField := func(fl validator.FieldLevel) bool {
		// Ensure we're validating a Service object
		service, ok := fl.Parent().Interface().(Service)
		if !ok {
			// Return true if it's not a Service (skip validation for non-Service fields)
			return true
		}

		// Get the value of the field being validated
		value := fl.Field().Int()

		// Apply service-specific validation logic
		switch service.Name {
		case "scraper", "monitor":
			if fl.FieldName() == "BatchSize" {
				return value >= 50 && value <= 10000
			}
			if fl.FieldName() == "Sleep" {
				return value >= 0 // Sleep must be non-negative
			}
		case "api", "ipfs":
			if fl.FieldName() == "Port" {
				return value >= 1024 && value <= 65535
			}
		default:
			// Unknown service type
			return false
		}

		return true
	}

	validate.RegisterValidation("strict_url", validateStrictURL)
	validate.RegisterValidation("ping_one", validatePingOne)
	validate.RegisterValidation("is_writable", validateIsWritable)
	validate.RegisterValidation("path_exists", validatePathExists)
	validate.RegisterValidation("opt_min", validateOptMin)
	validate.RegisterValidation("opt_max", validateOptMax)
	validate.RegisterValidation("dirpath", validateDirPath)
	validate.RegisterValidation("service_field", validateServiceField)
	validate.RegisterStructValidation(validateService, Service{})
}
