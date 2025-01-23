package wizard

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/types"
)

var ErrUserQuit = errors.New("user quit")
var ErrUserBack = errors.New("user back")
var ErrUserHelp = errors.New("user help")
var ErrUserEdit = errors.New("user edit")
var ErrUserChains = errors.New("user chains")

// the following are wrapped, so we can check with errors.Is
var ErrValidate = errors.New("")
var ErrValidateWarn = errors.New("")
var ErrValidateMsg = errors.New("")
var ErrSkipQuestion = errors.New("")

var screenWidth = 80

func displayScreen(w *Wizard, screenIndex int) error {
	if screenIndex < 0 || screenIndex >= len(w.screens) {
		return fmt.Errorf("invalid screen index")
	}

	curScreen := w.Current()
	curScreen.Wizard = w

	for i := curScreen.Current; i < len(curScreen.Questions); i++ {
		nSkipped := 0
		question := &curScreen.Questions[i]
		if skip := question.Prepare(curScreen); !skip {
			caret := curScreen.GetCaret(w.caret, i, nSkipped)
			curScreen.Display(question, caret)

			reader := bufio.NewReader(os.Stdin)
			if input, err := reader.ReadString('\n'); err != nil {
				return err
			} else {
				err := question.processResponse(input)
				if err != nil {
					if errors.Is(err, ErrValidate) {
						fmt.Println(colors.Red + input + " " + question.ErrorStr + colors.Off)
						i--
					} else if errors.Is(err, ErrSkipQuestion) {
						i++
					} else if errors.Is(err, ErrValidateWarn) {
						curScreen.Display(question, caret)
						if os.Getenv("NO_CLEAR") != "true" {
							time.Sleep(2000 * time.Millisecond)
						}
					} else if errors.Is(err, ErrValidateMsg) {
						curScreen.Display(question, caret)
						if os.Getenv("NO_CLEAR") != "true" {
							time.Sleep(500 * time.Millisecond)
						}
					} else if errors.Is(err, ErrUserHelp) {
						curScreen.OpenHelp()
						i--
					} else if errors.Is(err, ErrUserEdit) {
						configPath := types.GetConfigFn()
						curScreen.EditFile(configPath)
						if err := curScreen.Reload(configPath); err != nil {
							return err
						}
						i--
					} else if errors.Is(err, ErrUserChains) {
						chainsPath := strings.ReplaceAll(types.GetConfigFn(), "config.yaml", "chains.json")
						curScreen.EditFile(chainsPath)
						i--
					} else if !errors.Is(err, ErrUserBack) || i == 0 {
						return err
					} else {
						prev := &curScreen.Questions[i-1]
						skip := prev.Prepare(curScreen)
						if skip {
							curScreen.Questions[i-1].ErrorStr = ""
							curScreen.Questions[i-2].ErrorStr = ""
							curScreen.Questions[i-1].Response = ""
							curScreen.Questions[i-2].Response = ""
							i -= 3
						} else {
							curScreen.Questions[i-1].ErrorStr = ""
							curScreen.Questions[i-1].Response = ""
							i -= 2
						}
					}
				}
			}
		} else {
			nSkipped++
		}
	}

	w.Next()
	return nil
}
