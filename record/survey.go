package record

import (
	"errors"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pterm/pterm"
)

// reference: https://github.com/AlecAivazis/survey

type surveyObj struct {
	Type        surveyPromptType        `json:"type"`
	Key         string                  `json:"key"`         // for all types
	ValueType   surveyPromptValueType   `json:"valueType"`   // for type input
	Message     string                  `json:"message"`     // for type input, confirm, select
	Options     []surveyPromptOptionObj `json:"options"`     // for type select
	LoopOptions []surveyPromptOptionObj `json:"loopOptions"` // for type loopSelectInput, loopSelectSelect
}

var (
	ErrInvalidSurveyType                 = errors.New("invalid survey prompt type")
	ErrInvalidSurveyValueType            = errors.New("invalid survey prompt value type")
	ErrInvalidSurveyLoopOptionsValueType = errors.New("invalid survey loop options value type")
)

type surveyPromptOptionObj struct {
	Desc  string      `json:"desc"`
	Value interface{} `json:"value"`
}

type surveyPromptType string

var (
	surveyPromptTypeInput            surveyPromptType = "input"
	surveyPromptTypeConfirm          surveyPromptType = "confirm"
	surveyPromptTypeSelect           surveyPromptType = "select"
	surveyPromptTypeLoopSelectInput  surveyPromptType = "loopSelectInput"
	surveyPromptTypeLoopSelectSelect surveyPromptType = "loopSelectSelect"
)

type surveyPromptValueType string

var (
	surveyPromptValueTypeNumber surveyPromptValueType = "number"
	surveyPromptValueTypeBool   surveyPromptValueType = "bool"
	surveyPromptValueTypeString surveyPromptValueType = "string"
)

func (o *surveyObj) Execute(store map[string]interface{}) error {
	var (
		intValue    int
		numberValue float64
		boolValue   bool
		stringValue string
	)

	var prompt survey.Prompt
	switch o.Type {
	case surveyPromptTypeInput:
		prompt = &survey.Input{Message: o.Message}

		// type surveyPromptTypeInput store ValueType value
		switch o.ValueType {
		case surveyPromptValueTypeNumber:
			store[o.Key] = &numberValue
		case surveyPromptValueTypeBool:
			store[o.Key] = &boolValue
		case surveyPromptValueTypeString:
			store[o.Key] = &stringValue
		default:
			return ErrInvalidSurveyValueType
		}

	case surveyPromptTypeConfirm:
		prompt = &survey.Confirm{Message: o.Message}

		// type surveyPromptTypeConfirm store boolean value
		store[o.Key] = &boolValue

	case surveyPromptTypeSelect:
		options := make([]string, 0, len(o.Options))
		for _, option := range o.Options {
			options = append(options, option.Desc)
		}

		prompt = &survey.Select{Message: o.Message, Options: options, PageSize: 10}

		// type surveyPromptTypeSelect store the chosen option index into `intValue`
		store[o.Key] = &intValue

	case surveyPromptTypeLoopSelectSelect, surveyPromptTypeLoopSelectInput:
		return o.handleLoopTypePrompt(store)

	default:
		return ErrInvalidSurveyType
	}

	if err := survey.AskOne(prompt, store[o.Key]); err != nil {
		return err
	}

	switch o.Type {
	case surveyPromptTypeSelect:
		store[o.Key] = &o.Options[intValue].Value
	}

	return nil
}

var loopTypePromptFinishTag = "*FINISH*"

func (o *surveyObj) handleLoopTypePrompt(store map[string]interface{}) error {
	options := make([]string, 0, len(o.LoopOptions))
	for _, option := range o.LoopOptions {
		options = append(options, option.Desc)
	}
	options = append(options, loopTypePromptFinishTag)

	for {
		selectPrompt := &survey.Select{
			Message:  "Select an option:",
			Options:  options,
			PageSize: 10,
		}

		var optionId int
		if err := survey.AskOne(selectPrompt, &optionId); err != nil {
			return err
		}

		if options[optionId] == loopTypePromptFinishTag {
			break
		}

		subKey, ok := o.LoopOptions[optionId].Value.(string)
		if !ok {
			return ErrInvalidSurveyLoopOptionsValueType
		}

		key := o.Key + "." + subKey

		var innerSurvey *surveyObj
		switch o.Type {
		case surveyPromptTypeLoopSelectInput:
			innerSurvey = &surveyObj{
				Type:      surveyPromptTypeInput,
				Key:       key,
				ValueType: o.ValueType,
				Message:   o.Message,
			}
		case surveyPromptTypeLoopSelectSelect:
			innerSurvey = &surveyObj{
				Type:    surveyPromptTypeSelect,
				Key:     key,
				Message: o.Message,
				Options: o.Options,
			}
		}

		if err := innerSurvey.Execute(store); err != nil {
			return err
		}

		// append a check mark before the select prompt desc
		options[optionId] = "👌 " + o.LoopOptions[optionId].Desc

		pterm.Println("")
	}

	return nil
}
