package app

// func (k *KhedraApp) initAction(c *cli.Context) error {
// 	logger.Fatal("Should never happen in initAction")
// _ = c // linter
// if _, err := k.ConfigMaker(); err != nil {
// 	return fmt.Errorf("failed to load config: %w", err)
// }

// // Initialize the control service -- we need it for init
// k.controlSvc = k.initializeControlSvc()

// // Register the help handler for context-sensitive help
// registerHelpHandler()

// // Register validation functions for real-time feedback
// registerValidationFunctions()

// steps := getInitScreens()

// reloadConfig := func(string) (any, error) {
// 	if cfg, err := LoadConfig(); err != nil {
// 		return k.config, err
// 	} else {
// 		k.config = &cfg
// 		return k.config, err
// 	}
// }

// w := wizard.NewWizard(steps, "", k.config, reloadConfig)
// if err := w.Run(); err != nil {
// 	return err
// }

// 	return nil
// }

// func validWarn(msg, value string) error {
// 	if strings.Contains(msg, "%s") {
// 		return fmt.Errorf(msg+"%w", value, wizard.ErrValidateWarn)
// 	}
// 	return fmt.Errorf(msg+"%w", wizard.ErrValidateWarn)
// }

// func validContinue() error {
// 	return fmt.Errorf("continue %w", wizard.ErrValidateMsg)
// }

// func validOk(msg, value string) error {
// 	if strings.Contains(msg, "%s") {
// 		return fmt.Errorf(msg+"%w", value, wizard.ErrValidateMsg)
// 	}
// 	return fmt.Errorf(msg+"%w", wizard.ErrValidateMsg)
// }

// func validSkipNext() error {
// 	return fmt.Errorf("skip next %w", wizard.ErrSkipQuestion)
// }

// // --------------------------------------------------------
// type processFn[T any] func(cfg *types.Config) (string, T, error)

// // --------------------------------------------------------
// func prepare[T any](q *wizard.Question, fn processFn[T]) (string, error) {
// 	if cfg, ok := q.Screen.Wizard.Backing.(*types.Config); ok {
// 		input, copy, err := fn(cfg)
// 		bytes, _ := json.Marshal(copy)
// 		q.State = string(bytes)
// 		return input, err
// 	}
// 	return "", validContinue()
// }

// // --------------------------------------------------------
// func confirm[T any](q *wizard.Question, fn processFn[T]) (string, error) {
// 	if cfg, ok := q.Screen.Wizard.Backing.(*types.Config); ok {
// 		input, copy, err := fn(cfg)
// 		if !errors.Is(err, wizard.ErrValidate) {
// 			err1 := cfg.WriteToFile(types.GetConfigFnNoCreate())
// 			if err1 != nil {
// 				fmt.Println(colors.Red+"error writing config file: %v", err, colors.Off)
// 			}
// 		}
// 		bytes, _ := json.Marshal(copy)
// 		q.State = string(bytes)
// 		return input, err
// 	}
// 	return "", validContinue()
// }

// func getInitScreens() []wizard.Screen {
// 	return []wizard.Screen{
// 		getWelcomeScreen(),
// 		getGeneralScreen(),
// 		getChainsScreen(),
// 		getServicesScreen(),
// 		getServicePortsScreen(),
// 		getLoggingScreen(),
// 		getSummaryScreen(),
// 	}
// }
