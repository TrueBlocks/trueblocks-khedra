package app

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/types"
)

func LoadConfig() (types.Config, error) {
	cfg, err := loadFileConfig()
	if err != nil {
		return types.Config{}, fmt.Errorf("failed to load file configuration: %w", err)
	}
	keys := types.GetEnvironmentKeys(cfg, types.InEnv)
	if err := types.ApplyEnv(keys, &cfg); err != nil {
		return types.Config{}, fmt.Errorf("failed to apply environment configuration: %w", err)
	}

	if err := validateConfig(cfg); err != nil {
		return types.Config{}, fmt.Errorf("configuration validation failed: %w", err)
	}

	if err := initializeFolders(cfg); err != nil {
		return types.Config{}, fmt.Errorf("failed to initialize folders: %w", err)
	}

	// if err := types.Validate.Struct(mergedCfg); err != nil {
	// 	return types.Config{}, fmt.Errorf("validation error: %v", err)
	// }

	return cfg, nil
}

// func LoadConfig() (types.Config, bool, error) {
// 	var fileK = koanf.New(".")
// 	var envK = koanf.New(".")

// 	fn := types.GetConfigFn()
// 	if err := fileK.Load(file.Provider(fn), yaml.Parser()); err != nil {
// 		return types.Config{}, true, fmt.Errorf("koanf.Load failed for file %s: %v", fn, err)
// 	}

// 	fileCfg := types.NewConfig()
// 	if err := fileK.Unmarshal("", &fileCfg); err != nil {
// 		return types.Config{}, true, fmt.Errorf("koanf.Unmarshal failed for file configuration: %v", err)
// 	}

// 	for key, chain := range fileCfg.Chains {
// 		chain.Name = key
// 		fileCfg.Chains[key] = chain
// 	}

// 	for key, service := range fileCfg.Services {
// 		service.Name = key
// 		fileCfg.Services[key] = service
// 	}

// 	fieldTypeMap := buildFieldTypeMap(reflect.TypeOf(types.Config{}), "")

// 	err := envK.Load(env.ProviderWithValue("TB_KHEDRA_", ".", func(key, value string) (string, interface{}) {
// 		transformedKey := strings.ToLower(strings.TrimPrefix(key, "TB_KHEDRA_"))
// 		transformedKey = strings.ReplaceAll(transformedKey, "_", ".")

// 		if strings.HasSuffix(transformedKey, ".rpcs") {
// 			return transformedKey, strings.Split(value, ",")
// 		}

// 		if fieldType, ok := fieldTypeMap[transformedKey]; ok {
// 			if fieldType.Kind() == reflect.Slice {
// 				return transformedKey, strings.Split(value, ",")
// 			} else if fieldType.Kind() == reflect.Bool {
// 				parsedValue, err := strconv.ParseBool(value)
// 				if err != nil {
// 					return "", fmt.Errorf("invalid boolean value for %s: %v", key, err)
// 				}
// 				return transformedKey, parsedValue
// 			}
// 		}

// 		return transformedKey, value
// 	}), nil)
// 	if err != nil {
// 		return types.Config{}, true, fmt.Errorf("koanf.Load failed for environment variables: %v", err)
// 	}

// 	envCfg := types.Config{} // Empty config to unmarshal into
// 	if err := envK.Unmarshal("", &envCfg); err != nil {
// 		return types.Config{}, true, fmt.Errorf("koanf.Unmarshal failed for environment configuration: %v", err)
// 	}

// 	for key, chain := range envCfg.Chains {
// 		if existingChain, exists := fileCfg.Chains[key]; exists {
// 			if len(chain.RPCs) > 0 {
// 				existingChain.RPCs = chain.RPCs
// 			}
// 			existingChain.Enabled = chain.Enabled
// 			fileCfg.Chains[key] = existingChain
// 		} else {
// 			return types.Config{}, true, fmt.Errorf("chain %s found in the environment but not in the configuration file", key)
// 		}
// 	}

// 	finalCfg := fileCfg

// 	for key, chain := range fileCfg.Chains {
// 		if len(chain.RPCs) == 0 {
// 			return types.Config{}, true, fmt.Errorf("chain %s has an empty RPCs field, which is not allowed", key)
// 		}
// 	}

// 	configPath := utils.ExpandPath("~/.khedra")
// 	coreFile.EstablishFolder(configPath)

// 	finalCfg.General.DataFolder = utils.ExpandPath(finalCfg.General.DataFolder)
// 	coreFile.EstablishFolder(finalCfg.General.DataFolder)

// 	finalCfg.Logging.Folder = utils.ExpandPath(finalCfg.Logging.Folder)
// 	coreFile.EstablishFolder(finalCfg.Logging.Folder)

// 	if err := types.Validate.Struct(finalCfg); err != nil {
// 		return types.Config{}, true, err
// 	}

// 	for name, service := range finalCfg.Services {
// 		envEnabled := os.Getenv(fmt.Sprintf("TB_KHEDRA_SERVICES_%s_ENABLED", strings.ToUpper(name)))
// 		if envEnabled != "" {
// 			service.Enabled, _ = strconv.ParseBool(envEnabled)
// 		}

// 		envPort := os.Getenv(fmt.Sprintf("TB_KHEDRA_SERVICES_%s_PORT", strings.ToUpper(name)))
// 		if envPort != "" {
// 			port, err := strconv.Atoi(envPort)
// 			if err == nil {
// 				service.Port = port
// 			}
// 		}

// 		finalCfg.Services[name] = service
// 	}

// 	return finalCfg, true, nil
// }

// // Recursively build a map of field types from a struct
// func buildFieldTypeMap(t reflect.Type, prefix string) map[string]reflect.Type {
// 	fieldMap := make(map[string]reflect.Type)

// 	for i := 0; i < t.NumField(); i++ {
// 		field := t.Field(i)
// 		fieldKey := prefix + strings.ToLower(field.Name)

// 		// Add the field to the map
// 		fieldMap[fieldKey] = field.Type

// 		// Recursively parse nested structs and slices
// 		if field.Type.Kind() == reflect.Struct {
// 			for k, v := range buildFieldTypeMap(field.Type, fieldKey+".") {
// 				fieldMap[k] = v
// 			}
// 		} else if field.Type.Kind() == reflect.Slice && field.Type.Elem().Kind() == reflect.Struct {
// 			for k, v := range buildFieldTypeMap(field.Type.Elem(), fieldKey+".") {
// 				fieldMap[k] = v
// 			}
// 		}
// 	}

// 	return fieldMap
// }

/*
// EstablishConfig either reads an existing configuration file or creates it if it doesn't exist.
func (a *App) EstablishConfig() error {
	for _, arg := range os.Args {
		if arg == "--help" || arg == "-h" || arg == "--version" {
			return nil
		}
	}
	var ok bool
	var err error
	if a.Config.ConfigPath, ok = os.LookupEnv("TB_NODE_DATAFOLDER"); !ok {
		return errors.New("environment variable `TB_NODE_DATAFOLDER` is required but not found")
	} else {
		if a.Config.ConfigPath, err = cleanDataPath(a.Config.ConfigPath); err != nil {
			return err
		}
	}
	a.Logger.Info("data directory", "dataFolder", a.Config.ConfigPath)

	var targets string
	chainStr, ok := os.LookupEnv("TB_NODE_CHAINS")
	if !ok {
		chainStr, targets = "mainnet", "mainnet"
	} else {
		if chainStr, targets, err = cleanChainString(chainStr); err != nil {
			return err
		}
	}
	a.Logger.Info("configured chains", "chainStr", chainStr, "targets", targets)
	a.Config.Targets = strings.Split(targets, ",")

	chains := strings.Split(chainStr, ",")
	for _, chain := range chains {
		key := "TB_NODE_" + strings.ToUpper(chain) + "RPC"
		if providerUrl, ok := os.LookupEnv(key); !ok {
			msg := fmt.Sprintf("environment variable `%s` is required but not found (implied by TB_NODE_CHAINS=%s)", key, chainStr)
			return errors.New(msg)
		} else {
			providerUrl = strings.Trim(providerUrl, "/")
			if !isValidURL(providerUrl) {
				return fmt.Errorf("invalid URL for %s: %s", key, providerUrl)
			}
			a.Config.ProviderMap[chain] = providerUrl
		}
	}

	// // Set the environment trueblocks-core needs
	os.Setenv("XDG_CONFIG_HOME", a.Config.ConfigPath)
	os.Setenv("TB_SETTINGS_DEFAULTCHAIN", "mainnet")
	os.Setenv("TB_SETTINGS_INDEXPATH", a.Config.IndexPath())
	os.Setenv("TB_SETTINGS_CACHEPATH", a.Config.CachePath())
	for chain, providerUrl := range a.Config.ProviderMap {
		envKey := "TB_CHAINS_" + strings.ToUpper(chain) + "_RPCPROVIDER"
		os.Setenv(envKey, providerUrl)
	}

	for _, env := range os.Environ() {
		if (strings.HasPrefix(env, "TB_") || strings.HasPrefix(env, "XDG_")) && strings.Contains(env, "=") {
			parts := strings.Split(env, "=")
			if len(parts) > 1 {
				a.Logger.Info("environment", parts[0], parts[1])
			} else {
				a.Logger.Info("environment", parts[0], "<empty>")
			}
		}
	}

	for _, chain := range chains {
		providerUrl := a.Config.ProviderMap[chain]
		if err := a.tryConnect(chain, providerUrl, 5); err != nil {
			return err
		} else {
			a.Logger.Info("test connection", "result", "okay", "chain", chain, "providerUrl", providerUrl)
		}
	}

	configFn := filepath.Join(a.Config.ConfigPath, "trueBlocks.toml")
	if file.FileExists(configFn) {
		a.Logger.Info("config loaded", "configFile", configFn, "nChains", len(a.Config.ProviderMap))
		// check to make sure the config file has all the chains
		contents := file.AsciiFileToString(configFn)
		for chain := range a.Config.ProviderMap {
			search := "[chains." + chain + "]"
			if !strings.Contains(contents, search) {
				msg := fmt.Sprintf("config file {%s} does not contain {%s}", configFn, search)
				msg = colors.ColoredWith(msg, colors.Red)
				return errors.New(msg)
			}
		}
		return nil
	}

	if err := file.EstablishFolder(a.Config.ConfigPath); err != nil {
		return err
	}
	for _, chain := range chains {
		chainConfig := filepath.Join(a.Config.ConfigPath, "config", chain)
		if err := file.EstablishFolder(chainConfig); err != nil {
			return err
		}
	}

	tmpl, err := template.New("tmpl").Parse(configTmpl)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, &a.Config); err != nil {
		return err
	}

	_ = file.StringToAsciiFile(configFn, buf.String())
	a.Logger.Info("Created config file", "configFile", configFn, "nChains", len(a.Config.ProviderMap))

	return nil
}

func cleanDataPath(in string) (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return in, err
	}
	out := strings.ReplaceAll(in, "PWD", pwd)

	home, err := os.UserHomeDir()
	if err != nil {
		return in, err
	}
	out = strings.ReplaceAll(out, "~", home)
	out = strings.ReplaceAll(out, "HOME", home)
	ret := filepath.Clean(out)
	if strings.HasSuffix(ret, "/unchained") {
		ret = strings.ReplaceAll(ret, "/unchained", "")
	}
	return ret, nil
}

var configTmpl string = `[version]
  current = "v4.0.0"

[settings]
  cachePath = "{{.CachePath}}"
  defaultChain = "mainnet"
  indexPath = "{{.IndexPath}}"

[keys]
  [keys.etherscan]
    apiKey = ""

[chains]{{.ChainDescriptors}}
`
*/
