package record

import (
	"errors"

	"github.com/pterm/pterm"
)

// reference: https://github.com/pterm/pterm

type ptermObj struct {
	Type    ptermType       `json:"type"`
	Section ptermSectionObj `json:"section"`
	Prefix  ptermPrefixObj  `json:"prefix"`
}

var (
	ErrInvalidPtermType        = errors.New("invalid pterm type")
	ErrInvalidPtermPrefixLevel = errors.New("invalid pterm prefix level")
)

type ptermType string

const (
	ptermTypeSection ptermType = "section"
	ptermTypePrefix  ptermType = "prefix"
)

// ref: https://github.com/pterm/pterm#section
type ptermSectionObj struct {
	Level   int    `json:"level"`
	Println string `json:"println"`
}

// ref: https://github.com/pterm/pterm#prefix
type ptermPrefixLevel string

const (
	ptermPrefixDebug   ptermPrefixLevel = "debug"
	ptermPrefixInfo    ptermPrefixLevel = "info"
	ptermPrefixSuccess ptermPrefixLevel = "success"
	ptermPrefixWarning ptermPrefixLevel = "warning"
	ptermPrefixError   ptermPrefixLevel = "error"
)

type ptermPrefixObj struct {
	Level   ptermPrefixLevel `json:"level"`
	Println string           `json:"println"`
}

func (o *ptermObj) Execute() error {
	switch o.Type {
	case ptermTypeSection:
		if o.Section.Level == 0 {
			o.Section.Level = 1
		}

		pterm.DefaultSection.WithLevel(o.Section.Level).Println(o.Section.Println)
	case ptermTypePrefix:
		var logger pterm.PrefixPrinter

		switch o.Prefix.Level {
		case ptermPrefixDebug:
			logger = pterm.Debug
		case ptermPrefixInfo:
			logger = pterm.Info
		case ptermPrefixSuccess:
			logger = pterm.Success
		case ptermPrefixWarning:
			logger = pterm.Warning
		case ptermPrefixError:
			logger = pterm.Error
		default:
			return ErrInvalidPtermPrefixLevel
		}

		logger.Println(o.Prefix.Println)
	default:
		return ErrInvalidPtermType
	}

	return nil
}
