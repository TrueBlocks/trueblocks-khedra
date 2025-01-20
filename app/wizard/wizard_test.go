package wizard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWizard(t *testing.T) {
	validInitialization := func() {
		screens := []Screen{
			{Title: "Screen 1"},
			{Title: "Screen 2"},
		}
		wizard := NewWizard(screens, "--> ")

		assert.NotNil(t, wizard)
		assert.Equal(t, 0, wizard.current)
		assert.Equal(t, "--> ", wizard.caret)
		assert.False(t, wizard.completed)
	}
	t.Run("Valid Initialization", func(t *testing.T) { validInitialization() })

	panicOnEmptyScreens := func() {
		assert.Panics(t, func() {
			NewWizard([]Screen{}, "--> ")
		})
	}
	t.Run("Panic on Empty Screens", func(t *testing.T) { panicOnEmptyScreens() })
}

func TestWizardNavigation(t *testing.T) {
	nextTest := func() {
		screens := []Screen{
			{Title: "Screen 1"},
			{Title: "Screen 2"},
			{Title: "Screen 3"},
		}
		wizard := NewWizard(screens, "--> ")

		assert.True(t, wizard.Next())
		assert.Equal(t, 1, wizard.current)

		assert.True(t, wizard.Next())
		assert.Equal(t, 2, wizard.current)

		assert.False(t, wizard.Next())
		assert.True(t, wizard.completed)
	}
	t.Run("Next Test", func(t *testing.T) { nextTest() })

	prevTest := func() {
		screens := []Screen{
			{Title: "Screen 1"},
			{Title: "Screen 2"},
			{Title: "Screen 3"},
		}
		wizard := NewWizard(screens, "--> ")
		wizard.Next()
		wizard.Next()

		assert.True(t, wizard.Prev())
		assert.Equal(t, 1, wizard.current)

		assert.True(t, wizard.Prev())
		assert.Equal(t, 0, wizard.current)

		assert.True(t, wizard.Prev())
		assert.Equal(t, 0, wizard.current)
	}
	t.Run("Prev Test", func(t *testing.T) { prevTest() })

	currentTest := func() {
		screens := []Screen{
			{Title: "Screen 1"},
			{Title: "Screen 2"},
		}
		wizard := NewWizard(screens, "--> ")

		current := wizard.Current()
		assert.Equal(t, "Screen 1", current.Title)

		wizard.Next()
		current = wizard.Current()
		assert.Equal(t, "Screen 2", current.Title)
	}
	t.Run("Current Test", func(t *testing.T) { currentTest() })
}

func TestWizardRun(t *testing.T) {
	mockDisplayScreen := func(w *Wizard, screenIndex int) error {
		if screenIndex == 0 {
			return ErrUserQuit
		}
		return nil
	}

	mockScreens := []Screen{
		{Title: "Screen 1"},
		{Title: "Screen 2"},
	}

	wizard := NewWizard(mockScreens, "--> ")
	wizard.displayFn = mockDisplayScreen

	err := wizard.Run()
	assert.NoError(t, err)
	assert.False(t, wizard.completed)
}
