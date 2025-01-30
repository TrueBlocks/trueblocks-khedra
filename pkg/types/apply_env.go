package types

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

var (
	ErrEmptyEnvValue   = errors.New("environment variable is empty")
	ErrInvalidEnvValue = errors.New("environment variable has an invalid value")
)

const (
	// Prefixes
	PrefixChains   = "TB_KHEDRA_CHAINS_"
	PrefixServices = "TB_KHEDRA_SERVICES_"

	// General Keys
	KeyDataFolder = "TB_KHEDRA_GENERAL_DATAFOLDER"
	KeyStrategy   = "TB_KHEDRA_GENERAL_STRATEGY"
	KeyDetail     = "TB_KHEDRA_GENERAL_DETAIL"

	// Logging Keys
	KeyLoggingFolder     = "TB_KHEDRA_LOGGING_FOLDER"
	KeyLoggingFilename   = "TB_KHEDRA_LOGGING_FILENAME"
	KeyLoggingToFile     = "TB_KHEDRA_LOGGING_TOFILE"
	KeyLoggingMaxSize    = "TB_KHEDRA_LOGGING_MAXSIZE"
	KeyLoggingMaxBackups = "TB_KHEDRA_LOGGING_MAXBACKUPS"
	KeyLoggingMaxAge     = "TB_KHEDRA_LOGGING_MAXAGE"
	KeyLoggingCompress   = "TB_KHEDRA_LOGGING_COMPRESS"
)

const (
	// Chain Keys
	ChainKeyRPCs    = "rpcs"
	ChainKeyEnabled = "enabled"
	ChainKeyChainID = "chainid"

	// Service Keys
	ServiceKeyEnabled   = "enabled"
	ServiceKeyPort      = "port"
	ServiceKeySleep     = "sleep"
	ServiceKeyBatchSize = "batchsize"
)

func ApplyEnv(keys []string, receiver *Config) error {
	return applyEnv(keys, receiver)
}

// applyEnv updates the provided Config by applying environment variable values for the given keys.
// It validates values for correctness (non-empty, parsable types) and assigns them to matching fields
// in the Config struct. Returns an error if any value is invalid or empty.
func applyEnv(keys []string, receiver *Config) error {
	wrapError := func(base error, key, val string) error {
		return errors.New(base.Error() + ": key=[" + key + "], value=[" + val + "]")
	}

	validateNonEmptyEnv := func(key string, envValue string) error {
		if envValue == "" {
			return wrapError(ErrEmptyEnvValue, key, envValue)
		}
		return nil
	}

	validateValueParsing := func(key string, parseErr error) error {
		if parseErr != nil {
			return wrapError(ErrInvalidEnvValue, key, "")
		}
		return nil
	}

	isValidKey := func(prefix string, key string) error {
		if !strings.HasPrefix(key, prefix) {
			return wrapError(ErrInvalidEnvValue, "invalid key prefix: "+key, "")
		}
		parts := strings.Split(strings.TrimPrefix(key, prefix), "_")
		if len(parts) < 2 {
			return wrapError(ErrInvalidEnvValue, "invalid key structure: "+key, "")
		}
		return nil
	}

	// Define handlers for Chains
	chainHandlers := map[string]func(*Chain, string) error{
		ChainKeyRPCs: func(chain *Chain, value string) error {
			rpcs := strings.Split(value, ",")
			if len(rpcs) == 0 {
				return wrapError(ErrInvalidEnvValue, ChainKeyRPCs, value)
			}
			chain.RPCs = rpcs
			return nil
		},
		ChainKeyEnabled: func(chain *Chain, value string) error {
			enabled, err := strconv.ParseBool(value)
			if err := validateValueParsing(ChainKeyEnabled, err); err != nil {
				return err
			}
			chain.Enabled = enabled
			return nil
		},
		ChainKeyChainID: func(chain *Chain, value string) error {
			chainId, err := strconv.Atoi(value)
			if err := validateValueParsing(ChainKeyChainID, err); err != nil {
				return err
			}
			chain.ChainID = chainId
			return nil
		},
	}

	// Define handlers for Services
	serviceHandlers := map[string]func(*Service, string) error{
		ServiceKeyEnabled: func(service *Service, value string) error {
			enabled, err := strconv.ParseBool(value)
			if err := validateValueParsing(ServiceKeyEnabled, err); err != nil {
				return err
			}
			service.Enabled = enabled
			return nil
		},
		ServiceKeyPort: func(service *Service, value string) error {
			port, err := strconv.Atoi(value)
			if err := validateValueParsing(ServiceKeyPort, err); err != nil {
				return err
			}
			service.Port = port
			return nil
		},
		ServiceKeySleep: func(service *Service, value string) error {
			sleep, err := strconv.Atoi(value)
			if err := validateValueParsing(ServiceKeySleep, err); err != nil {
				return err
			}
			service.Sleep = sleep
			return nil
		},
		ServiceKeyBatchSize: func(service *Service, value string) error {
			batchSize, err := strconv.Atoi(value)
			if err := validateValueParsing(ServiceKeyBatchSize, err); err != nil {
				return err
			}
			service.BatchSize = batchSize
			return nil
		},
	}

	for _, key := range keys {
		envValue := os.Getenv(key)
		if err := validateNonEmptyEnv(key, envValue); err != nil {
			return err
		}

		switch {
		// General settings
		case key == KeyDataFolder:
			receiver.General.DataFolder = envValue
		case key == KeyStrategy:
			receiver.General.Strategy = envValue
		case key == KeyDetail:
			receiver.General.Detail = envValue

		// Logging settings
		case key == KeyLoggingFolder:
			receiver.Logging.Folder = envValue
		case key == KeyLoggingFilename:
			receiver.Logging.Filename = envValue
		case key == KeyLoggingToFile:
			toFile, err := strconv.ParseBool(envValue)
			if err := validateValueParsing(key, err); err != nil {
				return err
			}
			receiver.Logging.ToFile = toFile
		case key == KeyLoggingMaxSize:
			size, err := strconv.Atoi(envValue)
			if err := validateValueParsing(key, err); err != nil {
				return err
			}
			receiver.Logging.MaxSize = size
		case key == KeyLoggingMaxBackups:
			backups, err := strconv.Atoi(envValue)
			if err := validateValueParsing(key, err); err != nil {
				return err
			}
			receiver.Logging.MaxBackups = backups
		case key == KeyLoggingMaxAge:
			age, err := strconv.Atoi(envValue)
			if err := validateValueParsing(key, err); err != nil {
				return err
			}
			receiver.Logging.MaxAge = age
		case key == KeyLoggingCompress:
			compress, err := strconv.ParseBool(envValue)
			if err := validateValueParsing(key, err); err != nil {
				return err
			}
			receiver.Logging.Compress = compress

		// Chains
		case strings.HasPrefix(key, PrefixChains):
			if err := isValidKey(PrefixChains, key); err != nil {
				return err
			}
			if err := processMap(PrefixChains, receiver.Chains, key, envValue, chainHandlers); err != nil {
				return err
			}

		// Services
		case strings.HasPrefix(key, PrefixServices):
			if err := isValidKey(PrefixServices, key); err != nil {
				return err
			}
			if err := processMap(PrefixServices, receiver.Services, key, envValue, serviceHandlers); err != nil {
				return err
			}
		}
	}

	return nil
}

// processMap searches for a handler in the provided handlers map based on the item's key, applies the
// environment variable value to the corresponding item in target, and updates the map if successful.
// Returns an error if the key is invalid or the value cannot be applied correctly.
func processMap[T any](prefix string, target map[string]T, key, envValue string, handlers map[string]func(*T, string) error) error {
	parts := strings.Split(strings.TrimPrefix(key, prefix), "_")
	if len(parts) < 2 {
		return nil
	}
	itemName, itemKey := strings.ToLower(parts[0]), strings.ToLower(parts[1])

	handler, ok := handlers[itemKey]
	if !ok {
		// Nothing to do if there's no handler for this key.
		return nil
	}

	item, exists := target[itemName]
	if !exists {
		var zeroValue T
		item = zeroValue
	}

	if err := handler(&item, envValue); err != nil {
		return err
	}

	target[itemName] = item
	return nil
}
