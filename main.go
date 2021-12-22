package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/justin0u0/NTHU-OS-Demo/export"
	"github.com/justin0u0/NTHU-OS-Demo/question"
	"github.com/justin0u0/NTHU-OS-Demo/record"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	cmd := &cobra.Command{
		Use:   "demo [command]",
		Short: "A CLI tool for NTHU OS demo",
		Long: `
A CLI tool written in Golang to make OS TA demo easily, simple to add new demo
and question by just setting new JSON files. A tool that make TAs happy 😀.
`,
	}

	cmd.AddCommand(question.NewQuestionCommand())
	cmd.AddCommand(record.NewRecordCommand())
	cmd.AddCommand(export.NewExportCommand())

	if os.Getenv("PTERM_DEBUG") == "true" {
		pterm.EnableDebugMessages()
	}

	if err := cmd.Execute(); err != nil {
		log.Fatal("fail to run command: " + err.Error())
	}
}
