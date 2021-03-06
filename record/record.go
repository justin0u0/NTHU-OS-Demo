package record

import (
	"bytes"
	"embed"
	"encoding/json"
	"os"
	"os/user"
	"strconv"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/justin0u0/NTHU-OS-Demo/version"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

//go:embed assets
var recordFS embed.FS

var storeDir string

func NewRecordCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "record [path]",
		Short:   "Start to record information and save the result into the store",
		Example: "demo record example",
		Args:    cobra.ExactArgs(1),
		Run:     run,
	}

	cmd.Flags().StringVarP(&storeDir, "store", "s", "record/store", "The directory to store the result file")

	return cmd
}

var (
	storeKeyCreatedAt = "createdAt"
	storeKeyCreatedBy = "createdBy"
	storeKeyVersion   = "version"
)

func run(_ *cobra.Command, args []string) {
	fileName := "assets/" + args[0] + ".json"
	storeKeyCreatedAt = args[0] + "." + storeKeyCreatedAt
	storeKeyCreatedBy = args[0] + "." + storeKeyCreatedBy
	storeKeyVersion = args[0] + "." + storeKeyVersion

	pterm.Debug.Println("Running record from file:", fileName)

	f, err := recordFS.ReadFile(fileName)
	if err != nil {
		pterm.Fatal.Println("Fail to read record file:", err)
	}

	var rec recorder
	if err := json.NewDecoder(bytes.NewReader(f)).Decode(&rec.Processes); err != nil {
		pterm.Fatal.Println("Fail to parse record object:", err)
	}

	if err := rec.Execute(); err != nil {
		pterm.Fatal.Println("Fail to execute record process:", err)
	}

	// add additional informations
	rec.store[storeKeyCreatedAt] = time.Now().Format(time.RFC3339)
	if user, err := user.Current(); err != nil {
		pterm.Error.Println("Fail to get current username:", err)
		rec.store[storeKeyCreatedBy] = "unknown"
	} else {
		rec.store[storeKeyCreatedBy] = user.Name
	}
	rec.store[storeKeyVersion] = version.Version

	// marshal result into json bytes
	result, err := json.Marshal(rec.store)
	if err != nil {
		pterm.Fatal.Println("Fail to marshal result store:", err)
	}

	pterm.Println("")
	pterm.Success.Println("result: ", string(result))
	pterm.Println("")

	var store bool
	if err := survey.AskOne(&survey.Confirm{Message: "Do you want to store the result?"}, &store); err != nil {
		pterm.Fatal.Println("Fail to confirm should store:", err)
	}

	if store {
		if err := os.MkdirAll(storeDir, 0755); err != nil {
			pterm.Fatal.Println("Fail to mkdir store directory:", err)
		}

		fileName := storeDir + "/" + args[0] + "_" + strconv.FormatInt(time.Now().Unix(), 10) + ".json"

		if err := os.WriteFile(fileName, result, 0644); err != nil {
			pterm.Fatal.Println("Fail to store result to file:", err)
		}
	}

	pterm.Println("")
	pterm.Success.Println("done.")
}
