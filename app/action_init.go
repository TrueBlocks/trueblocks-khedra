package app

import (
	"fmt"
	"strings"

	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/wizard"
	"github.com/urfave/cli/v2"
)

func (k *KhedraApp) initAction(c *cli.Context) error {
	_ = c // linter
	if _, err := k.ConfigMaker(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	steps := []wizard.Screen{
		wizard.AddScreen(getWelcomeScreen()),
		wizard.AddScreen(getGeneralScreen()),
		wizard.AddScreen(getServicesScreen()),
		wizard.AddScreen(getChainsScreen()),
		wizard.AddScreen(getLoggingScreen()),
		wizard.AddScreen(getSummaryScreen()),
	}

	reloadConfig := func(string) (any, error) {
		if cfg, err := LoadConfig(); err != nil {
			return k.config, err
		} else {
			k.config = &cfg
			return k.config, err
		}
	}

	w := wizard.NewWizard(steps, "", k.config, reloadConfig)
	if err := w.Run(); err != nil {
		return err
	}

	return nil
}

func validWarn(msg, value string) error {
	if strings.Contains(msg, "%s") {
		return fmt.Errorf(msg+"%w", value, wizard.ErrValidateWarn)
	}
	return fmt.Errorf(msg+"%w", wizard.ErrValidateWarn)
}

func validContinue() error {
	return fmt.Errorf("continue %w", wizard.ErrValidateMsg)
}

func validOk(msg, value string) error {
	if strings.Contains(msg, "%s") {
		return fmt.Errorf(msg+"%w", value, wizard.ErrValidateMsg)
	}
	return fmt.Errorf(msg+"%w", wizard.ErrValidateMsg)
}

func validSkipNext(msg, value string) error {
	if strings.Contains(msg, "%s") {
		return fmt.Errorf(msg+"%w", value, wizard.ErrSkipQuestion)
	}
	return fmt.Errorf(msg+"%w", wizard.ErrSkipQuestion)
}
