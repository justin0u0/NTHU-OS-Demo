package export

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	exportDir    string
	storeDir     string
	filterRegexp string
)

//go:embed assets
var exportFS embed.FS

func NewExportCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "export [path]",
		Short:   "Export stored results with custom rules",
		Example: "demo export example",
		Args:    cobra.ExactArgs(1),
		Run:     run,
	}

	cmd.Flags().StringVarP(&exportDir, "export", "e", "export/store", "The directory to store the exported csv files")
	cmd.Flags().StringVarP(&storeDir, "store", "s", "record/store", "The directory to load all result files")
	cmd.Flags().StringVarP(&filterRegexp, "filter", "f", ".*\\.json", "The regex pattern to filter files")

	return cmd
}

func run(_ *cobra.Command, args []string) {
	resultFileNames, err := loadResultFiles()
	if err != nil {
		pterm.Fatal.Println("Fail to load result file:", err)
	}

	exp, err := loadExportFile("assets/" + args[0] + ".json")
	if err != nil {
		pterm.Fatal.Println("Fail to load export file:", err)
	}

	for _, resultFileName := range resultFileNames {
		handleResultFile(resultFileName, exp)
	}

	if err := os.MkdirAll(exportDir, 0755); err != nil {
		pterm.Fatal.Println("Fail to mkdir export directory:", err)
	}

	exportDetailFileName := fmt.Sprintf("%s/%s_detail_%d.csv", exportDir, args[0], time.Now().Unix())
	if err := exp.exportDetailRowsCSV(exportDetailFileName); err != nil {
		pterm.Fatal.Println("Fail to export detail file:", err)
	}
}

func loadResultFiles() ([]string, error) {
	files, err := os.ReadDir(storeDir)
	if err != nil {
		return nil, fmt.Errorf("fail to read the store directory %s: %w", storeDir, err)
	}

	fileNames := make([]string, 0)

	re, err := regexp.Compile(filterRegexp)
	if err != nil {
		return nil, fmt.Errorf("fail to compile regexp %s: %w", filterRegexp, err)
	}

	for _, file := range files {
		if !re.MatchString(file.Name()) {
			continue
		}

		fileNames = append(fileNames, storeDir+"/"+file.Name())
	}

	return fileNames, nil
}

func loadExportFile(fileName string) (*exporter, error) {
	pterm.Debug.Println("Running export from file:", fileName)

	f, err := exportFS.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("fail to read export file %s: %w", fileName, err)
	}

	var exp exporter
	if err := json.NewDecoder(bytes.NewReader(f)).Decode(&exp); err != nil {
		return nil, fmt.Errorf("fail to decode export file %s: %w", fileName, err)
	}

	pterm.Debug.Println("Load export file content:", exp)

	if err := exp.compile(); err != nil {
		return nil, fmt.Errorf("fail to compile exporter %s: %w", fileName, err)
	}

	return &exp, nil
}

func handleResultFile(fileName string, exp *exporter) error {
	pterm.Debug.Println("Handling result file:", fileName)

	f, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("fail to read result file %s: %w", fileName, err)
	}

	result := make(map[string]interface{})
	if err := json.NewDecoder(bytes.NewReader(f)).Decode(&result); err != nil {
		return fmt.Errorf("fail to decode reuslt file %s: %w", fileName, err)
	}

	pterm.Debug.Println("Load result file content:", result)

	if err := exp.evaluateDetail(result); err != nil {
		return fmt.Errorf("fail to evaluate detail: %w", err)
	}

	return nil
}
