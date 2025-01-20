package wizard

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
)

var ErrUserQuit = errors.New("user quit")
var ErrUserBack = errors.New("user back")
var ErrUserHelp = errors.New("user help")

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
	curScreen.wiz = w

	for i := curScreen.Current; i < len(curScreen.Questions); i++ {
		nSkipped := 0
		question := &curScreen.Questions[i]
		if skip := question.Prepare(curScreen); !skip {
			curScreen.Display()
			question.Display(curScreen.GetCaret(w.caret, i, nSkipped))

			reader := bufio.NewReader(os.Stdin)
			if input, err := reader.ReadString('\n'); err != nil {
				return err
			} else {
				err := question.processResponse(input)
				if err != nil {
					if errors.Is(err, ErrValidate) {
						fmt.Println(colors.Red + input + " " + question.ErrorMsg + colors.Off)
						i--
					} else if errors.Is(err, ErrSkipQuestion) {
						i++
					} else if errors.Is(err, ErrValidateWarn) {
						msg := question.Prompt("Response") + err.Error()
						fmt.Println(colors.Yellow + msg + colors.Off)
						if os.Getenv("NO_CLEAR") != "true" {
							time.Sleep(3000 * time.Millisecond)
						}
					} else if errors.Is(err, ErrValidateMsg) {
						msg := question.Prompt("Response") + err.Error()
						fmt.Println(colors.Yellow + msg + colors.Off)
						if os.Getenv("NO_CLEAR") != "true" {
							time.Sleep(1250 * time.Millisecond)
						}
					} else if errors.Is(err, ErrUserHelp) {
						curScreen.OpenHelp()
						i--
					} else if !errors.Is(err, ErrUserBack) || i == 0 {
						return err
					} else {
						i -= 2
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
