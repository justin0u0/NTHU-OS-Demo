package question

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

//go:embed assets
var questionsFS embed.FS

type question struct {
	Id   string `json:"id"`
	Desc string `json:"desc"`
}

type questionGroup struct {
	PicksPerStudent int         `json:"picksPerStudent"`
	Questions       []*question `json:"questions"`
}

func NewQuestionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "question [path]",
		Short:   "Generate list of questions for demo",
		Example: "demo question example",
		Args:    cobra.ExactArgs(1),
		Run:     run,
	}

	return cmd
}

func run(_ *cobra.Command, args []string) {
	fileName := "assets/" + args[0] + ".json"

	pterm.Debug.Println("Generating questions from file:", fileName)

	f, err := questionsFS.ReadFile(fileName)
	if err != nil {
		pterm.Fatal.Println("Fail to read questions file:", err)
	}

	var questionGroups []*questionGroup
	if err := json.NewDecoder(bytes.NewReader(f)).Decode(&questionGroups); err != nil {
		pterm.Fatal.Println("Fail to parse questions:", err)
	}

	for _, group := range questionGroups {
		rand.Shuffle(len(group.Questions), func(i, j int) {
			group.Questions[i], group.Questions[j] = group.Questions[j], group.Questions[i]
		})
	}

	for _, groupId := range []int{1, 2} {
		pterm.DefaultSection.Println("Group " + strconv.Itoa(groupId))

		for _, group := range questionGroups {
			offset := groupId * group.PicksPerStudent

			for i := offset; i < group.PicksPerStudent+offset; i++ {
				q := group.Questions[i]

				logger := pterm.PrefixPrinter{
					MessageStyle: &pterm.ThemeDefault.InfoMessageStyle,
					Prefix: pterm.Prefix{
						Style: &pterm.ThemeDefault.InfoPrefixStyle,
						Text:  fmt.Sprintf("%3s", q.Id),
					},
				}
				logger.Println(q.Desc)

				pterm.Println("")
			}
		}
	}
}
