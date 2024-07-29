package cmd

import (
	"errors"
	"os"
	"strings"

	"github.com/gabe565/ascii-movie/cmd/get"
	"github.com/gabe565/ascii-movie/cmd/ls"
	"github.com/gabe565/ascii-movie/cmd/play"
	"github.com/gabe565/ascii-movie/cmd/serve"
	"github.com/gabe565/ascii-movie/internal/config"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

var version = "beta"

func NewCommand() *cobra.Command {
	cmdVersion, commit := buildVersion(version)
	cmd := &cobra.Command{
		Use:     "ascii-movie",
		Short:   "Command line ASCII movie player.",
		Version: cmdVersion,

		PersistentPreRunE: preRun,
		DisableAutoGenTag: true,
	}
	cmd.AddCommand(
		play.NewCommand(),
		serve.NewCommand(version, commit),
		ls.NewCommand(),
		get.NewCommand(),
	)
	config.RegisterLogFlags(cmd)
	return cmd
}

const EnvPrefix = "ASCII_MOVIE_"

func preRun(cmd *cobra.Command, _ []string) error {
	if err := loadFlagEnvs(cmd.Flags()); err != nil {
		return err
	}
	config.InitLog(cmd)
	return nil
}

func loadFlagEnvs(flags *flag.FlagSet) error {
	var errs []error
	flags.VisitAll(func(f *flag.Flag) {
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
