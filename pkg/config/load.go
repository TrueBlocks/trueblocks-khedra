package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	coreFile "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/types"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/utils"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

func LoadConfig() (types.Config, error) {
	var fileK = koanf.New(".")
	var envK = koanf.New(".")

	fn := types.GetConfigFn()
	if err := fileK.Load(file.Provider(fn), yaml.Parser()); err != nil {
		return types.Config{}, fmt.Errorf("koanf.Load failed for file %s: %v", fn, err)
	}

	fileCfg := types.NewConfig()
	if err := fileK.Unmarshal("", &fileCfg); err != nil {
		return types.Config{}, fmt.Errorf("koanf.Unmarshal failed for file configuration: %v", err)
	}

	for key, chain := range fileCfg.Chains {
		chain.Name = key
		fileCfg.Chains[key] = chain
	}

	for key, service := range fileCfg.Services {
		service.Name = key
		fileCfg.Services[key] = service
	}

	fieldTypeMap := buildFieldTypeMap(reflect.TypeOf(types.Config{}), "")

	err := envK.Load(env.ProviderWithValue("TB_KHEDRA_", ".", func(key, value string) (string, interface{}) {
		transformedKey := strings.ToLower(strings.TrimPrefix(key, "TB_KHEDRA_"))
		transformedKey = strings.ReplaceAll(transformedKey, "_", ".")

		if strings.HasSuffix(transformedKey, ".rpcs") {
			return transformedKey, strings.Split(value, ",")
		}

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

		return transformedKey, value
	}), nil)
	if err != nil {
		return types.Config{}, fmt.Errorf("koanf.Load failed for environment variables: %v", err)
	}

	envCfg := types.Config{} // Empty config to unmarshal into
	if err := envK.Unmarshal("", &envCfg); err != nil {
		return types.Config{}, fmt.Errorf("koanf.Unmarshal failed for environment configuration: %v", err)
	}

	for key, chain := range envCfg.Chains {
		if existingChain, exists := fileCfg.Chains[key]; exists {
			if len(chain.RPCs) > 0 {
				existingChain.RPCs = chain.RPCs
			}
			existingChain.Enabled = chain.Enabled
			fileCfg.Chains[key] = existingChain
		} else {
			return types.Config{}, fmt.Errorf("chain %s found in the environment but not in the configuration file", key)
		}
	}

	finalCfg := fileCfg

	for key, chain := range fileCfg.Chains {
		if len(chain.RPCs) == 0 {
			return types.Config{}, fmt.Errorf("chain %s has an empty RPCs field, which is not allowed", key)
		}
	}

	configPath := utils.ExpandPath("~/.khedra")
	coreFile.EstablishFolder(configPath)

	finalCfg.General.DataDir = utils.ExpandPath(finalCfg.General.DataDir)
	coreFile.EstablishFolder(finalCfg.General.DataDir)

	finalCfg.Logging.Folder = utils.ExpandPath(finalCfg.Logging.Folder)
	coreFile.EstablishFolder(finalCfg.Logging.Folder)

	if err := types.Validate.Struct(finalCfg); err != nil {
		return types.Config{}, err
	}

	for name, service := range finalCfg.Services {
		envEnabled := os.Getenv(fmt.Sprintf("TB_KHEDRA_SERVICES_%s_ENABLED", strings.ToUpper(name)))
		if envEnabled != "" {
			service.Enabled, _ = strconv.ParseBool(envEnabled)
		}

		envPort := os.Getenv(fmt.Sprintf("TB_KHEDRA_SERVICES_%s_PORT", strings.ToUpper(name)))
		if envPort != "" {
			port, err := strconv.Atoi(envPort)
			if err == nil {
				service.Port = port
			}
		}

		finalCfg.Services[name] = service
	}

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
