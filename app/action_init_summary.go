package app

// // screen|---------|---------|---------|---------|---------|---------|---|74
// func getSummaryScreen() wizard.Screen {
// 	sTitle := `Configuration Summary`
// 	sSubtitle := ``
// 	sInstructions := `Review your configuration and press enter to finish.`
// 	sBody := `
// Please review your configuration settings. If everything looks correct,
// press enter to save your configuration and exit the wizard.

// If you need to make changes, use "b" or "back" to return to previous screens.
// `
// 	sReplacements := []wizard.Replacement{
// 		{Color: colors.Yellow, Values: []string{sTitle}},
// 		{Color: colors.Green, Values: []string{
// 			"\"b\"", "\"back\"",
// 		}},
// 	}
// 	sQuestions := []wizard.Questioner{&sum0}
// 	sStyle := wizard.NewStyle()

// 	return wizard.Screen{
// 		Title:        sTitle,
// 		Subtitle:     sSubtitle,
// 		Body:         sBody,
// 		Instructions: sInstructions,
// 		Replacements: sReplacements,
// 		Questions:    sQuestions,
// 		Style:        sStyle,
// 	}
// }

// // --------------------------------------------------------
// var sum0 = wizard.Question{
// 	//.....question-|---------|---------|---------|---------|---------|----|65
// 	Question: `Press Enter to save your configuration, or "b" to go back.`,
// 	Hint:     `Configuration preview is shown above.`,
// 	PrepareFn: func(input string, q *wizard.Question) (string, error) {
// 		_ = input // delint
// 		cfg, ok := q.Screen.Wizard.Backing.(*types.Config)
// 		if !ok {
// 			return "", fmt.Errorf("could not cast backing data")
// 		}

// 		// Display configuration preview
// 		displayConfigPreview(cfg, q)

// 		return "", nil
// 	},
// 	Validate: func(input string, q *wizard.Question) (string, error) {
// 		// Only accept empty input or "b"/"back" otherwise
// 		if input != "" && input != "b" && input != "back" {
// 			return input, fmt.Errorf(`press Enter to continue or "b" to go back %w`, wizard.ErrValidate)
// 		}

// 		// Write the configuration to file when the user confirms
// 		if input == "" {
// 			cfg, ok := q.Screen.Wizard.Backing.(*types.Config)
// 			if !ok {
// 				return input, fmt.Errorf("could not access configuration %w", wizard.ErrValidate)
// 			}

// 			if err := cfg.WriteToFile(types.GetConfigFn()); err != nil {
// 				return input, fmt.Errorf("failed to save configuration: %s %w", err.Error(), wizard.ErrValidate)
// 			}

// 			return input, validOk("Configuration saved successfully", "")
// 		}

// 		return input, nil
// 	},
// }

// // displayConfigPreview shows a preview of the configuration
// func displayConfigPreview(cfg *types.Config, q *wizard.Question) {
// 	preview := strings.Builder{}

// 	// Add header
// 	preview.WriteString(colors.Yellow + "🔍 Configuration Preview:" + colors.Off + "\n\n")

// 	// General section
// 	preview.WriteString(colors.BrightBlue + "📂 General Settings:" + colors.Off + "\n")
// 	preview.WriteString(fmt.Sprintf("  Data Folder: %s\n", cfg.General.DataFolder))
// 	preview.WriteString(fmt.Sprintf("  Strategy: %s\n", cfg.General.Strategy))
// 	preview.WriteString(fmt.Sprintf("  Detail: %s\n\n", cfg.General.Detail))

// 	// Services section
// 	preview.WriteString(colors.BrightBlue + "🔧 Services:" + colors.Off + "\n")
// 	titleCase := cases.Title(language.English)
// 	for name, service := range cfg.Services {
// 		status := "Disabled"
// 		if service.Enabled {
// 			status = fmt.Sprintf("Enabled (Port: %d)", service.Port)
// 		}
// 		preview.WriteString(fmt.Sprintf("  %s: %s\n", titleCase.String(name), status))
// 	}
// 	preview.WriteString("\n")

// 	// Chains section
// 	preview.WriteString(colors.BrightBlue + "⛓️ Chains:" + colors.Off + "\n")

// 	// Sort chains for consistent display
// 	var chainNames = []string{}
// 	for name := range cfg.Chains {
// 		chainNames = append(chainNames, name)
// 	}
// 	sort.Strings(chainNames)

// 	for _, name := range chainNames {
// 		chain := cfg.Chains[name]
// 		status := "Disabled"
// 		if chain.Enabled {
// 			status = "Enabled"
// 		}

// 		preview.WriteString(fmt.Sprintf("  %s (%s):\n", titleCase.String(name), status))
// 		for i, rpc := range chain.RPCs {
// 			preview.WriteString(fmt.Sprintf("    RPC #%d: %s\n", i+1, rpc))
// 		}
// 	}
// 	preview.WriteString("\n")

// 	// Logging section
// 	preview.WriteString(colors.BrightBlue + "📝 Logging:" + colors.Off + "\n")
// 	preview.WriteString(fmt.Sprintf("  Log to File: %t\n", cfg.Logging.ToFile))
// 	preview.WriteString(fmt.Sprintf("  Log Folder: %s\n", cfg.Logging.Folder))
// 	preview.WriteString(fmt.Sprintf("  Log Level: %s\n", cfg.Logging.Level))

// 	// Create a box for the preview
// 	style := boxes.NewStyle()
// 	style.Width = 75
// 	style.Justify = boxes.Left
// 	style.BorderStyle = boxes.BorderStyle{
// 		TopLeft:     "┌",
// 		Top:         "─",
// 		TopRight:    "┐",
// 		Right:       "│",
// 		BottomRight: "┘",
// 		Bottom:      "─",
// 		BottomLeft:  "└",
// 		Left:        "│",
// 	}

// 	// Print the preview
// 	fmt.Println()
// 	box := boxes.NewBox("", preview.String(), style)
// 	box.Display()
// 	fmt.Println()

// 	// Adapt the question based on configuration validity
// 	q.Hint = "Review the preview above and press Enter to save your configuration."
// }
