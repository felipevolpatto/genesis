package tui

import (
	"errors"
	"regexp"

	"github.com/AlecAivazis/survey/v2"
	"github.com/felipevolpatto/genesis/internal/config"
)

var ErrInvalidRegexp = errors.New("invalid regular expression")

// askOne is a variable that holds the survey.AskOne function
var askOne = survey.AskOne

// PromptForVariables asks the user for values for each template variable
func PromptForVariables(vars map[string]config.Variable) (map[string]string, error) {
	answers := make(map[string]string)

	for name, v := range vars {
		prompt := &survey.Input{
			Message: v.Prompt,
			Default: v.Default,
		}

		var answer string
		var err error

		if v.Regex != "" {
			regex := regexp.MustCompile(v.Regex)
			err = askOne(prompt, &answer, survey.WithValidator(func(val interface{}) error {
				str, ok := val.(string)
				if !ok {
					return nil
				}
				if !regex.MatchString(str) {
					return ErrInvalidRegexp
				}
				return nil
			}))
		} else {
			err = askOne(prompt, &answer)
		}

		if err != nil {
			return nil, err
		}

		answers[name] = answer
	}

	return answers, nil
}

// ConfirmAction asks the user to confirm an action
func ConfirmAction(message string) (bool, error) {
	var confirm bool
	prompt := &survey.Confirm{
		Message: message,
	}

	err := askOne(prompt, &confirm)
	if err != nil {
		return false, err
	}

	return confirm, nil
}

// SelectOption asks the user to select an option from a list
func SelectOption(message string, options []string) (string, error) {
	var selected string
	prompt := &survey.Select{
		Message: message,
		Options: options,
	}

	err := askOne(prompt, &selected)
	if err != nil {
		return "", err
	}

	return selected, nil
} 