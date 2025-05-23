package wizard

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-khedra/v5/pkg/boxes"
	"github.com/TrueBlocks/trueblocks-khedra/v5/pkg/types"
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

	// Get the current screen
	curScreen := w.Current()
	curScreen.Wizard = w

	// Ensure we're using a Single border style for all box displays
	// We'll pass this when creating any boxes in the Display method
	// borderStyle := boxes.Single | boxes.All

	for i := curScreen.Current; i < len(curScreen.Questions); i++ {
		nSkipped := 0
		question := curScreen.Questions[i]

		if skip := question.Prepare(curScreen); !skip {
			caret := curScreen.GetCaret("-->", i, nSkipped)

			// Always display the screen, regardless of question content
			curScreen.Display(question, caret)

			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				return err
			}

			err = question.HandleResponse(strings.TrimSpace(input))
			if err != nil {
				switch {
				case errors.Is(err, ErrValidate):
					fmt.Println(colors.Red + input + " " + question.GetError() + colors.Off)
					i--
				case errors.Is(err, ErrSkipQuestion):
					i++
				case errors.Is(err, ErrValidateWarn):
					curScreen.Display(question, caret)
					if os.Getenv("NO_WAIT") != "true" {
						time.Sleep(2000 * time.Millisecond)
					}
				case errors.Is(err, ErrValidateMsg):
					curScreen.Display(question, caret)
					if os.Getenv("NO_WAIT") != "true" {
						time.Sleep(500 * time.Millisecond)
					}
				case errors.Is(err, ErrUserHelp):
					helpText := GetHelp(curScreen, question.(*Question))
					if helpText == "" {
						helpText = "No help is available for this item."
					}
					displayHelpScreen(helpText)
					i--
				case errors.Is(err, ErrUserEdit):
					configPath := types.GetConfigFn()
					_ = curScreen.EditFile(configPath)
					if err := curScreen.Reload(configPath); err != nil {
						return err
					}
					i--
				case errors.Is(err, ErrUserChains):
					chainsPath := strings.ReplaceAll(types.GetConfigFn(), "config.yaml", "chains.json")
					_ = curScreen.EditFile(chainsPath)
					i--
				case errors.Is(err, ErrUserBack):
					if i == 0 {
						return err
					}
					prevQuestion := curScreen.Questions[i-1]
					skip := prevQuestion.Prepare(curScreen)
					if skip {
						curScreen.Questions[i-2].Clear()
						prevQuestion.Clear()
						i -= 3
					} else {
						prevQuestion.Clear()
						i -= 2
					}
				default:
					return err
				}
			}
		} else {
			nSkipped++
		}
	}

	w.Next()
	return nil
}

// displayHelpScreen shows a help screen with formatted content
func displayHelpScreen(helpText string) {
	fmt.Print(clearScreen)

	// Create a styled box for the help content
	style := boxes.NewStyle()
	style.Width = screenWidth - 4

	// Format the help text for display
	contentLines := strings.Split(helpText, "\n")
	boxContent := boxes.Box(contentLines, style.Width, boxes.Single|boxes.All, boxes.Left)

	// Print the formatted box content
	fmt.Println(boxContent)

	// Prompt user to continue
	fmt.Print("\n" + colors.BrightBlue + "Press Enter to continue..." + colors.Off)
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
	fmt.Print(clearScreen)
}
