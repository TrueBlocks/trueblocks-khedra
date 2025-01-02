package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	coreFile "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// Declare a variable for the function that retrieves the config file path so
// we can mock it during testing.
var getConfigFn = mustGetConfigFn

// MustLoadConfig returns a validated Config struct or fails fatally.
func MustLoadConfig(fn string) *Config {
	var err error
	var cfg Config
	if cfg, err = loadConfig(); err != nil {
		log.Fatalf("error loading config: %v", err)
	}
	return &cfg
}

func loadConfig() (Config, error) {
	var fileK = koanf.New(".")
	var envK = koanf.New(".")

	// Load the base configuration file
	fn := getConfigFn()
	if err := fileK.Load(file.Provider(fn), yaml.Parser()); err != nil {
		return Config{}, fmt.Errorf("koanf.Load failed for file %s: %v", fn, err)
	}

	// Unmarshal the configuration from the file
	fileCfg := NewConfig()
	if err := fileK.Unmarshal("", &fileCfg); err != nil {
		return Config{}, fmt.Errorf("koanf.Unmarshal failed for file configuration: %v", err)
	}

	// Create a map of chain and service names to indices
	chainNameToIndex := make(map[string]int)
	for i, chain := range fileCfg.Chains {
		chainNameToIndex[strings.ToLower(chain.Name)] = i
	}

	serviceNameToIndex := make(map[string]int)
	for i, service := range fileCfg.Services {
		serviceNameToIndex[strings.ToLower(service.Name)] = i
	}

	// Build a recursive map of field types from the Config struct
	fieldTypeMap := buildFieldTypeMap(reflect.TypeOf(Config{}), "")

	// Load environment variables with proper key mapping and type handling
	err := envK.Load(env.ProviderWithValue("TB_KHEDRA_", ".", func(key, value string) (string, interface{}) {
		// Transform the environment variable key into a nested configuration key
		transformedKey := strings.ToLower(strings.TrimPrefix(key, "TB_KHEDRA_"))
		transformedKey = strings.ReplaceAll(transformedKey, "_", ".")

		// Check for chains.* and services.*
		parts := strings.Split(transformedKey, ".")
		if len(parts) > 2 && parts[0] == "chains" {
			chainName := parts[1]
			if index, ok := chainNameToIndex[chainName]; ok {
				// Replace chain name with index
				parts[1] = fmt.Sprintf("%d", index)
				transformedKey = strings.Join(parts, ".")
			}
		} else if len(parts) > 2 && parts[0] == "services" {
			serviceName := parts[1]
			if index, ok := serviceNameToIndex[serviceName]; ok {
				// Replace service name with index
				parts[1] = fmt.Sprintf("%d", index)
				transformedKey = strings.Join(parts, ".")
			}
		}

		// Check the field type and handle arrays
		if fieldType, ok := fieldTypeMap[transformedKey]; ok {
			if fieldType.Kind() == reflect.Slice {
				return transformedKey, strings.Split(value, ",")
			} else if fieldType.Kind() == reflect.Bool {
				parsedValue, err := strconv.ParseBool(value)
				if err != nil {
					return "", fmt.Errorf("invalid boolean value for %s: %v", key, err)
				}
				return transformedKey, parsedValue
			}
		}

		// Let Koanf handle the types automatically for other fields
		return transformedKey, value
	}), nil)
	if err != nil {
		return Config{}, fmt.Errorf("koanf.Load failed for environment variables: %v", err)
	}

	// Unmarshal the environment configuration
	envCfg := NewConfig()
	if err := envK.Unmarshal("", &envCfg); err != nil {
		return Config{}, fmt.Errorf("koanf.Unmarshal failed for environment configuration: %v", err)
	}

	// Merge environment configuration into file configuration
	finalCfg := mergeConfigs(fileCfg, envCfg)

	// Ensure directories are established for paths
	configPath := expandPath("~/.khedra")
	coreFile.EstablishFolder(configPath)

	finalCfg.General.DataPath = expandPath(finalCfg.General.DataPath)
	coreFile.EstablishFolder(finalCfg.General.DataPath)

	finalCfg.Logging.Folder = expandPath(finalCfg.Logging.Folder)
	coreFile.EstablishFolder(finalCfg.Logging.Folder)

	// Validate the final configuration
	if err := validate.Struct(finalCfg); err != nil {
		return Config{}, err
	}

	// Return the final configuration
	return finalCfg, nil
}

// Recursively build a map of field types from a struct
func buildFieldTypeMap(t reflect.Type, prefix string) map[string]reflect.Type {
	fieldMap := make(map[string]reflect.Type)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldKey := prefix + strings.ToLower(field.Name)

		// Add the field to the map
		fieldMap[fieldKey] = field.Type

		// Recursively parse nested structs and slices
		if field.Type.Kind() == reflect.Struct {
			for k, v := range buildFieldTypeMap(field.Type, fieldKey+".") {
				fieldMap[k] = v
			}
		} else if field.Type.Kind() == reflect.Slice && field.Type.Elem().Kind() == reflect.Struct {
			for k, v := range buildFieldTypeMap(field.Type.Elem(), fieldKey+".") {
				fieldMap[k] = v
			}
		}
	}

	return fieldMap
}

// Merge the environment configuration into the file configuration
func mergeConfigs(fileCfg, envCfg Config) Config {
	// Merge General
	if envCfg.General.DataPath != NewGeneral().DataPath {
		fileCfg.General.DataPath = envCfg.General.DataPath
	}

	// Merge Logging
	if envCfg.Logging.Folder != NewLogging().Folder {
		fileCfg.Logging.Folder = envCfg.Logging.Folder
	}
	if envCfg.Logging.Filename != NewLogging().Filename {
		fileCfg.Logging.Filename = envCfg.Logging.Filename
	}
	if envCfg.Logging.MaxSizeMb != NewLogging().MaxSizeMb {
		fileCfg.Logging.MaxSizeMb = envCfg.Logging.MaxSizeMb
	}
	if envCfg.Logging.MaxBackups != NewLogging().MaxBackups {
		fileCfg.Logging.MaxBackups = envCfg.Logging.MaxBackups
	}
	if envCfg.Logging.MaxAgeDays != NewLogging().MaxAgeDays {
		fileCfg.Logging.MaxAgeDays = envCfg.Logging.MaxAgeDays
	}
	if envCfg.Logging.Compress != NewLogging().Compress {
		fileCfg.Logging.Compress = envCfg.Logging.Compress
	}
	if envCfg.Logging.LogLevel != NewLogging().LogLevel {
		fileCfg.Logging.LogLevel = envCfg.Logging.LogLevel
	}

	// Merge Chains
	for i, chain := range envCfg.Chains {
		if i < len(fileCfg.Chains) {
			if len(chain.RPCs) > 0 {
				fileCfg.Chains[i].RPCs = chain.RPCs
			}
			if chain.Enabled {
				fileCfg.Chains[i].Enabled = chain.Enabled
			}
		} else {
			// Add new chain from the environment
			fileCfg.Chains = append(fileCfg.Chains, chain)
		}
	}

	// Merge Services
	for i, service := range envCfg.Services {
		if i < len(fileCfg.Services) {
			if service.Port != 0 {
				fileCfg.Services[i].Port = service.Port
			}
			fileCfg.Services[i].Enabled = service.Enabled
		} else {
			// Add new service from the environment
			fileCfg.Services = append(fileCfg.Services, service)
		}
	}

	return fileCfg
}

// mustGetConfigFn returns the path to the config file which must
// be either in the current folder or in the default location. If
// there is no such file, establish it
func mustGetConfigFn() string {
	// current folder
	fn := expandPath("config.yaml")
	if coreFile.FileExists(fn) {
		return fn
	}

	// expanded default config folder
	fn = expandPath(filepath.Join(mustGetConfigDir(), "config.yaml"))
	if coreFile.FileExists(fn) {
		return fn
	}

	_ = establishConfig(fn)
	return fn
}

func mustGetConfigDir() string {
	var err error
	cfgDir := expandPath("~/.khedra")

	if !coreFile.FolderExists(cfgDir) {
		if err = coreFile.EstablishFolder(cfgDir); err != nil {
			log.Fatalf("error establishing log folder %s: %v", cfgDir, err)
		}
	}

	if writable := IsWritable(cfgDir); !writable {
		log.Fatalf("log directory %s is not writable: %v", cfgDir, err)
	}

	return cfgDir
}

// expandPath returns an absolute path expanded for ~, $HOME or other env variables
func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") || path == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("error getting home directory: %v", err)
		}
		path = filepath.Join(home, strings.TrimPrefix(path, "~"))
	}

	var err error
	path, err = filepath.Abs(os.ExpandEnv(path))
	if err != nil {
		log.Fatalf("error making path absolute: %v", err)
	}

	return path
}

// IsWritable checks to see if a folder is writable
func IsWritable(path string) bool {
	tmpFile := filepath.Join(path, ".test")

	if fil, err := os.Create(tmpFile); err != nil {
		fmt.Println(fmt.Errorf("folder %s is not writable: %v", path, err))
		return false
	} else {
		fil.Close()
		if err := os.Remove(tmpFile); err != nil {
			fmt.Println(fmt.Errorf("error cleaning up test file in %s: %v", path, err))
			return false
		}
	}

	return true
}
