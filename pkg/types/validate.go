package types

import (
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	coreFile "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/rpc"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/utils"
	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New()
	Validate.RegisterValidation("req_if_enabled", validateReqIfEnabled)
	Validate.RegisterValidation("strict_url", validateStrictURL)
	Validate.RegisterValidation("strict_url", validateStrictURL)
	Validate.RegisterValidation("ping_one", validatePingOne)
	Validate.RegisterValidation("is_writable", validateIsWritable)
	Validate.RegisterValidation("file_exists", validateFileExists)
	Validate.RegisterValidation("folder_exists", validateFolderExists)
	Validate.RegisterValidation("opt_min", validateOptMin)
	Validate.RegisterValidation("opt_max", validateOptMax)
	Validate.RegisterValidation("service_field", validateServiceField)
	Validate.RegisterStructValidation(validateService, Service{})
}

// reqIfEnabled checks if the parent Chain is enabled and ensures the RPCs field is non-nil and has at least one element.
func validateReqIfEnabled(fl validator.FieldLevel) bool {
	parent, ok := fl.Parent().Interface().(Chain)
	if !ok {
		return false
	}
	if !parent.Enabled {
		return true
	}

	rpcs := fl.Field().Interface().([]string)
	return len(rpcs) > 0
}

// validateStrictURL validates the given URL is well-formed and has a valid scheme
func validateStrictURL(fl validator.FieldLevel) bool {
	rawURL := fl.Field().String()
	parent := fl.Parent().Interface().(Chain)
	if !parent.Enabled {
		return true
	}

	u, err := url.Parse(rawURL)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != ""
}

// validatePingOne validates that at least one of the given URLs is reachable
func validatePingOne(fl validator.FieldLevel) bool {
	isTesting := os.Getenv("TEST_MODE") == "true"
	if isTesting {
		return true
	}
	parent := fl.Parent().Interface().(Chain)
	if !parent.Enabled {
		return true
	}
	url := fl.Field().String()
	if err := rpc.PingRpc(url); err == nil {
		return true
	}
	return false
}

// validateIsWritable validates that a path exists and is writable
func validateIsWritable(fl validator.FieldLevel) bool {
	path := fl.Field().String()
	path = utils.ResolvePath(path)

	testFile := filepath.Join(path, ".writable_check")
	file, err := os.Create(testFile)
	if err != nil {
		return false
	}
	_ = file.Close()
	_ = os.Remove(testFile)

	return true
}

// validateFileExists validates that a path is non-empty and that the path exists
func validateFileExists(fl validator.FieldLevel) bool {
	isTesting := os.Getenv("TEST_MODE") == "true"
	if isTesting {
		return true
	}
	path := fl.Field().String()
	if path == "" {
		return false
	}
	path = utils.ResolvePath(path)
	return coreFile.FileExists(path)
}

// validateFolderExists validates that a path is non-empty and that the path exists
func validateFolderExists(fl validator.FieldLevel) bool {
	isTesting := os.Getenv("TEST_MODE") == "true"
	if isTesting {
		return true
	}
	path := fl.Field().String()
	if path == "" {
		return false
	}
	path = utils.ResolvePath(path)
	return coreFile.FolderExists(path)
}

// validateOptMin validates that an integer is either zero (unset) or greater than or equal to a given value
func validateOptMin(fl validator.FieldLevel) bool {
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
func validateOptMax(fl validator.FieldLevel) bool {
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

// validateServiceField validates fields based on the service type
func validateServiceField(fl validator.FieldLevel) bool {
	service, ok := fl.Parent().Interface().(Service)
	if !ok {
		return true
	}

	value := fl.Field().Int()
	switch service.Name {
	case "scraper", "monitor":
		if fl.FieldName() == "BatchSize" {
			return value >= 50 && value <= 10000
		}
		if fl.FieldName() == "Sleep" {
			return value >= 0
		}
	case "api", "ipfs":
		if fl.FieldName() == "Port" {
			return value >= 1024 && value <= 65535
		}
	default:
		return false
	}

	return true
}

// validateService validates the configuration of a service which depends on the service type
func validateService(sl validator.StructLevel) {
	service := sl.Current().Interface().(Service)
	switch service.Name {
	case "api", "ipfs":
		if service.Port == 0 {
			sl.ReportError(service.Port, "Port", "port", "required_"+service.Name+"_port", "")
		}
	case "scraper", "monitor":
		if service.Sleep <= 0 {
			sl.ReportError(service.Sleep, "Sleep", "sleep", "required_scraper_monitor_sleep", "")
		}
		if service.BatchSize < 50 || service.BatchSize > 10000 {
			sl.ReportError(service.BatchSize, "BatchSize", "batchSize", "invalid_scraper_monitor_batchSize", "")
		}
	}
}
