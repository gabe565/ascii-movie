package cmd

import (
	"errors"
	"github.com/gabe565/ascii-movie/cmd/play"
	"github.com/gabe565/ascii-movie/cmd/serve"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"os"
	"strings"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ascii-movie",
		Short: "Command line ASCII movie player.",

		PersistentPreRunE: FlagEnvs,
		DisableAutoGenTag: true,
	}
	cmd.AddCommand(
		play.NewCommand(),
		serve.NewCommand(),
	)
	return cmd
}

const EnvPrefix = "ASCII_MOVIE_"

func FlagEnvs(cmd *cobra.Command, args []string) error {
	var errs []error
	cmd.Flags().VisitAll(func(f *flag.Flag) {
		optName := strings.ToUpper(f.Name)
		optName = strings.ReplaceAll(optName, "-", "_")
		varName := EnvPrefix + optName
		if val, ok := os.LookupEnv(varName); !f.Changed && ok {
			if err := f.Value.Set(val); err != nil {
				errs = append(errs, err)
			}
		}
	})
	return errors.Join(errs...)
}
