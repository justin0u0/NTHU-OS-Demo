package record

import (
	"os"

	imgcat "github.com/martinlindhe/imgcat/lib"
	"github.com/pterm/pterm"
)

// reference: github.com/martinlindhe/imgcat

type imgcatObj struct {
	FileName string `json:"fileName"`
}

func (o *imgcatObj) Execute() error {
	img, err := recordFS.Open(o.FileName)
	if err != nil {
		pterm.Error.Println("Fail to open image file:", err)
		return err
	}

	if err := imgcat.Cat(img, os.Stdout); err != nil {
		pterm.Error.Println("Fail to print image:", err)
		return err
	}

	return nil
}
