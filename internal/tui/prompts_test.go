package tui

import (
	"testing"

	"github.com/AlecAivazis/survey/v2"
	"github.com/felipevolpatto/genesis/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestPromptForVariables(t *testing.T) {
	tests := []struct {
		name           string
		vars           map[string]config.Variable
		expectedAnswer string
		expectError    bool
	}{
		{
			name: "basic prompt",
			vars: map[string]config.Variable{
				"name": {
					Prompt:  "Enter name:",
					Default: "default",
				},
			},
			expectedAnswer: "test",
			expectError:    false,
		},
		{
			name: "multiple prompts",
			vars: map[string]config.Variable{
				"name": {
					Prompt:  "Enter name:",
					Default: "default",
				},
				"description": {
					Prompt:  "Enter description:",
					Default: "default desc",
				},
			},
			expectedAnswer: "test",
			expectError:    false,
		},
		{
			name: "with regex validation",
			vars: map[string]config.Variable{
				"version": {
					Prompt:  "Enter version:",
					Default: "1.0.0",
					Regex:   "^\\d+\\.\\d+\\.\\d+$",
				},
			},
			expectedAnswer: "1.2.3",
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mock
			oldAskOne := askOne
			askOne = func(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error {
				switch v := response.(type) {
				case *string:
					*v = tt.expectedAnswer
				}
				return nil
			}
			defer func() { askOne = oldAskOne }()

			answers, err := PromptForVariables(tt.vars)
			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, answers)
			assert.Len(t, answers, len(tt.vars))

			for name := range tt.vars {
				assert.Equal(t, tt.expectedAnswer, answers[name])
			}
		})
	}
}

func TestConfirmAction(t *testing.T) {
	tests := []struct {
		name        string
		message     string
		answer      bool
		expectError bool
	}{
		{
			name:        "confirm yes",
			message:     "Proceed?",
			answer:      true,
			expectError: false,
		},
		{
			name:        "confirm no",
			message:     "Proceed?",
			answer:      false,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mock
			oldAskOne := askOne
			askOne = func(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error {
				switch v := response.(type) {
				case *bool:
					*v = tt.answer
				}
				return nil
			}
			defer func() { askOne = oldAskOne }()

			result, err := ConfirmAction(tt.message)
			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.answer, result)
		})
	}
}

func TestSelectOption(t *testing.T) {
	tests := []struct {
		name        string
		message     string
		options     []string
		answer      string
		expectError bool
	}{
		{
			name:        "select first option",
			message:     "Select an option:",
			options:     []string{"option1", "option2", "option3"},
			answer:      "option1",
			expectError: false,
		},
		{
			name:        "select last option",
			message:     "Select an option:",
			options:     []string{"option1", "option2", "option3"},
			answer:      "option3",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mock
			oldAskOne := askOne
			askOne = func(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error {
				switch v := response.(type) {
				case *string:
					*v = tt.answer
				}
				return nil
			}
			defer func() { askOne = oldAskOne }()

			result, err := SelectOption(tt.message, tt.options)
			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.answer, result)
		})
	}
} 