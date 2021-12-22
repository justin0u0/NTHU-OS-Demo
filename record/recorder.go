package record

import (
	"errors"
	"fmt"

	"github.com/pterm/pterm"
)

type recorder struct {
	Processes []struct {
		Type   recordType `json:"type"`
		Pterm  ptermObj   `json:"pterm"`
		Imgcat imgcatObj  `json:"imgcat"`
		Survey surveyObj  `json:"survey"`
	}

	store map[string]interface{}
}

var (
	ErrInvalidDemoType = errors.New("invalid record type")
)

type recordType string

const (
	recordTypePterm  recordType = "pterm"
	recordTypeImgcat recordType = "imgcat"
	recordTypeSurvey recordType = "survey"
)

func (o *recorder) Execute() error {
	if o.store == nil {
		o.store = make(map[string]interface{})
	}

	for _, p := range o.Processes {
		pterm.Debug.Println(fmt.Sprintf("%+v", p))

		var err error

		switch p.Type {
		case recordTypePterm:
			err = p.Pterm.Execute()
		case recordTypeImgcat:
			err = p.Imgcat.Execute()
		case recordTypeSurvey:
			err = p.Survey.Execute(o.store)
		default:
			err = ErrInvalidDemoType
		}

		if err != nil {
			pterm.Error.Println("Fail to execute record process: ", err)
		}
	}

	return nil
}
