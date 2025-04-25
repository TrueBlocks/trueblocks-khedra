package wizard

import (
	"fmt"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
)

// NavigationBar represents a visual navigation component showing progress through the wizard
type NavigationBar struct {
	CurrentStep    int
	TotalSteps     int
	StepTitles     []string
	Width          int
	ActiveColor    string
	InactiveColor  string
	CompletedColor string
}

// NewNavigationBar creates a new navigation bar
func NewNavigationBar(currentStep, totalSteps int, stepTitles []string) *NavigationBar {
	return &NavigationBar{
		CurrentStep:    currentStep,
		TotalSteps:     totalSteps,
		StepTitles:     stepTitles,
		Width:          78,
		ActiveColor:    colors.BrightBlue,
		InactiveColor:  colors.Black,
		CompletedColor: colors.Green,
	}
}

// Render displays the navigation bar
func (n *NavigationBar) Render() string {
	var output strings.Builder

	// Create the progress indicator
	stepsDisplay := fmt.Sprintf("Step %d/%d", n.CurrentStep+1, n.TotalSteps)
	output.WriteString(colors.BrightYellow + stepsDisplay + colors.Off + "\n")

	// Calculate available width after accounting for non-active steps
	nonActiveStepWidth := 3 // Width of "[N]" format for non-active steps (reducing from 5 to 3)
	totalNonActiveWidth := (n.TotalSteps - 1) * nonActiveStepWidth
	separatorsWidth := n.TotalSteps - 1                                            // One space between each step
	availableWidthForActive := n.Width - totalNonActiveWidth - separatorsWidth - 4 // Subtract 4 more characters for padding safety

	// Draw the navigation bar
	var navBar strings.Builder

	for i := 0; i < n.TotalSteps; i++ {
		if i == n.CurrentStep {
			// Active step - show as much of the title as possible
			title := n.StepTitles[i]
			displayTitle := title

			// Ensure title doesn't exceed available width
			if len(displayTitle) > availableWidthForActive {
				displayTitle = displayTitle[:availableWidthForActive-3] + "..."
			}

			// Add brackets with emphasis and color
			displayTitle = "[ " + displayTitle + " ]"
			navBar.WriteString(colors.BrightWhite + displayTitle + colors.Off)
		} else {
			// Non-active step - just show step number in brackets
			stepDisplay := fmt.Sprintf("[%d]", i+1)

			// Apply color based on step status (completed or upcoming)
			stepColor := n.InactiveColor
			if i < n.CurrentStep {
				stepColor = n.CompletedColor
			}

			navBar.WriteString(stepColor + stepDisplay + colors.Off)
		}

		// Add separator between steps
		if i < n.TotalSteps-1 {
			navBar.WriteString(" ")
		}
	}

	output.WriteString(navBar.String() + "\n")

	// Add a separator line that fits exactly within the width
	separatorWidth := n.Width - 3 // Subtract 2 to account for border characters on each side
	output.WriteString(strings.Repeat("─", separatorWidth) + "\n")

	return output.String()
}

// UpdateWizardWithNavBar adds navigation bars to each screen in the wizard
func UpdateWizardWithNavBar(w *Wizard) {
	if w == nil || len(w.screens) == 0 {
		return
	}

	totalScreens := len(w.screens)

	// Add a navigation bar to each screen
	for i := range w.screens {
		// Create a navigation bar for this screen
		navBar := NewNavigationBar(i+1, totalScreens, []string{w.screens[i].Title})

		// Set this wizard as the screen's wizard (for context)
		w.screens[i].Wizard = w

		// Set the navigation bar for this screen
		w.screens[i].NavigationBar = navBar
	}
}
