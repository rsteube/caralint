package cmd

import (
	"os"

	"github.com/owenrumney/go-sarif/sarif"
	"github.com/rsteube/caralint/internal/rules"
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "caralint",
	Short: "",
	Args:  cobra.MinimumNArgs(1),
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	Run: func(cmd *cobra.Command, args []string) {
		exitCode := 0

		report, err := sarif.New(sarif.Version210)
		if err != nil {
			panic(err)
		}
		run := sarif.NewRun("caralint", "https://github.com/rsteube/caralint")

		run.AddRule(rules.C001.Id).WithDescription(rules.C001.Desciption)

		for _, arg := range args {
			// TODO verify go file
			result, err := rules.C001.Check(arg)
			if err != nil {
				panic(err)
			}
			// TODO handle error

			run.AddResult(rules.C001.Id).
				WithMessage(sarif.NewTextMessage(result.Message)).
				WithLocation(
					sarif.NewLocationWithPhysicalLocation(
						sarif.NewPhysicalLocation().
							WithArtifactLocation(sarif.NewSimpleArtifactLocation(arg)).
							WithRegion(result.Region)))
		}
        report.AddRun(run)
		report.PrettyWrite(os.Stdout)
		os.Exit(exitCode)
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	carapace.Gen(rootCmd).PositionalAnyCompletion(
		carapace.ActionFiles(".go"),
	)
}
