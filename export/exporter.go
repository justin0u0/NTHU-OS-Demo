package export

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"

	"github.com/pterm/pterm"
)

type exporter struct {
	GroupSize int            `json:"groupSize"`
	Titles    []*exportTitle `json:"titles"`
	Rules     []*exportRule  `json:"rules"`

	// KeyColumns is the list of titles index to identify that 2 files belongs
	// to the same row and should merge information into the same row
	KeyColumns []int `json:"keyColumns"`
	// SortColumns is the list of titles index to give the order of the
	// `detailRows` and `summaryRows`
	SortColumns []int `json:"sortColumns"`

	detailRows  [][]string
	summaryRows [][]string
}

var (
	ErrInvalidExportRuleType      = errors.New("invalid export rule type")
	ErrRusultTypeMismatchRuleType = errors.New("result type mismatch rule type")
)

type exportRule struct {
	Regexp string         `json:"regexp"`
	Type   exportRuleType `json:"type"`
	Value  int            `json:"value"`
	For    string         `json:"for"`

	regexp *regexp.Regexp
}

type exportRuleType string

var (
	exportRuleTypePlainText        exportRuleType = "plaintext"
	exportRuleTypeValuableBoolean  exportRuleType = "valuable_boolean"
	exportRuleTypeValuableComplete exportRuleType = "valuable_complete"
	exportRuleTypeVaulablePartial  exportRuleType = "valuable_partial"
)

type exportTitle struct {
	Title   string `json:"title"`
	Regexp  string `json:"regexp"`
	Default string `json:"default"`

	regexp *regexp.Regexp
	index  int
}

func (e *exporter) compile() error {
	for _, rule := range e.Rules {
		regexp, err := regexp.Compile(rule.Regexp)
		if err != nil {
			return fmt.Errorf("fail to compile rule %s: %w", rule.Regexp, err)
		}

		rule.regexp = regexp
	}

	for i, title := range e.Titles {
		regexp, err := regexp.Compile(title.Regexp)
		if err != nil {
			return fmt.Errorf("fail to compile title %s: %w", title.Regexp, err)
		}

		title.regexp = regexp
		title.index = i
	}

	return nil
}

func (e *exporter) evaluateDetail(result map[string]interface{}) error {
	detail := make(map[string]interface{})

	for k, v := range result {
		rule := e.getRule(k)

		if rule == nil {
			pterm.Warning.Println("No match rule, skipping key:", k)
			continue
		}

		switch rule.Type {
		case exportRuleTypePlainText:
			value, ok := v.(string)
			if !ok {
				return fmt.Errorf("rule %s expect type string on key %s: %w", rule.Type, k, ErrRusultTypeMismatchRuleType)
			}

			detail[k] = value

		case exportRuleTypeValuableBoolean:
			value, ok := v.(bool)
			if !ok {
				return fmt.Errorf("rule %s expect type bool on key %s: %w", rule.Type, k, ErrRusultTypeMismatchRuleType)
			}

			detail[k] = 0
			if value {
				detail[k] = rule.Value
			}

		case exportRuleTypeValuableComplete:
			value, ok := v.(float64)
			if !ok {
				return fmt.Errorf("rule %s expect type float64 on key %s: %w", rule.Type, k, ErrRusultTypeMismatchRuleType)
			}

			detail[k] = value

		case exportRuleTypeVaulablePartial:
			value, ok := v.(float64)
			if !ok {
				return fmt.Errorf("rule %s expect type float64 on key %s: %w", rule.Type, k, ErrRusultTypeMismatchRuleType)
			}

			detail[k] = value * float64(rule.Value)

		default:
			return ErrInvalidExportRuleType
		}
	}

	pterm.Debug.Println("Evaluate detail done:", detail)

	if err := e.insertDetailRows(detail); err != nil {
		return fmt.Errorf("fail to insert detail rows: %w", err)
	}

	return nil
}

func (e *exporter) insertDetailRows(detail map[string]interface{}) error {
	detailRows, err := e.getDetailRows(detail)
	if err != nil {
		return fmt.Errorf("fail to get detail rows: %w", err)
	}

	for _, newRow := range detailRows {
		isNewRow := true

		newRowKey := e.getRowKey(newRow)
		pterm.Debug.Println("inserting detail rows with row key:", newRowKey)

		for _, existsRow := range e.detailRows {
			existsRowKey := e.getRowKey(existsRow)

			if newRowKey == existsRowKey {
				isNewRow = false
				pterm.Debug.Println("merging into exists row with row key:", newRowKey)

				// merge into exists row
				for i := range newRow {
					if newRow[i] != "" {
						existsRow[i] = newRow[i]
					}
				}
			}
		}

		if isNewRow {
			e.detailRows = append(e.detailRows, newRow)
		}
	}

	return nil
}

func (e *exporter) getRowKey(row []string) string {
	key := ""

	for i := range row {
		for _, keyColumn := range e.KeyColumns {
			if i == keyColumn {
				key += fmt.Sprintf("%d:%s;", i, row[i])
			}
		}
	}

	return key
}

func (e *exporter) getDetailRows(detail map[string]interface{}) ([][]string, error) {
	rows := make([][]string, e.GroupSize)
	for i := range rows {
		rows[i] = make([]string, len(e.Titles))
		for j := range rows[i] {
			rows[i][j] = e.Titles[j].Default
		}
	}

	for k, v := range detail {
		rule := e.getRule(k)
		if rule == nil {
			pterm.Warning.Println("No match rule, skipping key:", k)
			continue
		}

		title := e.getTitle(k)
		if title == nil {
			pterm.Warning.Println("No match title, skipping key:", k)
			continue
		}

		rowIndexes, err := e.getForRowIndexes(k, rule)
		if err != nil {
			return nil, fmt.Errorf("fail to get row indexes: %w", err)
		}

		for _, idx := range rowIndexes {
			rows[idx][title.index] = fmt.Sprintf("%v", v)
		}
	}

	return rows, nil
}

func (e *exporter) getRule(key string) *exportRule {
	for _, rule := range e.Rules {
		if rule.regexp.MatchString(key) {
			return rule
		}
	}

	return nil
}

func (e *exporter) getTitle(key string) *exportTitle {
	for _, title := range e.Titles {
		if title.regexp.MatchString(key) {
			return title
		}
	}

	return nil
}

func (e *exporter) getForRowIndexes(k string, rule *exportRule) ([]int, error) {
	var indexes []int

	if rule.For == "all" {
		for i := 0; i < e.GroupSize; i++ {
			indexes = append(indexes, i)
		}

		return indexes, nil
	}

	for i, name := range rule.regexp.SubexpNames() {
		if i > 0 && name == rule.For {
			matches := rule.regexp.FindStringSubmatch(k)

			groupId, err := strconv.Atoi(matches[i])
			if err != nil {
				return nil, fmt.Errorf("expect rule.for matches an integer: %w", err)
			}

			if groupId < 1 || groupId > e.GroupSize {
				return nil, fmt.Errorf("expect groupId is in range [1, groupSize]")
			}

			indexes = []int{groupId - 1}
			return indexes, nil
		}
	}

	return indexes, nil
}

func (e *exporter) exportDetailRowsCSV(fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("fail to create csv file: %w", err)
	}

	writer := csv.NewWriter(f)

	titleRow := make([]string, 0, len(e.Titles))
	for _, title := range e.Titles {
		titleRow = append(titleRow, title.Title)
	}

	if err := writer.Write(titleRow); err != nil {
		return fmt.Errorf("fail to write header row: %w", err)
	}

	detailRowSlice := &detailRowSlice{
		detailRows:  e.detailRows,
		sortColumns: e.SortColumns,
	}
	sort.Sort(detailRowSlice)

	if err := writer.WriteAll(detailRowSlice.detailRows); err != nil {
		return fmt.Errorf("fail to write detail rows: %w", err)
	}

	return nil
}
