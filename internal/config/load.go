package config

import (
	"errors"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const EnvPrefix = "ASCII_MOVIE_"

func Load(cmd *cobra.Command) (*Config, error) {
	c, ok := FromContext(cmd.Context())
	if !ok {
		panic("command missing context")
	}

	var errs []error
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if !f.Changed {
			if val, ok := os.LookupEnv(EnvName(f.Name)); ok {
				if err := f.Value.Set(val); err != nil {
					errs = append(errs, err)
				}
			}
		}
	})
	if len(errs) != 0 {
		return nil, errors.Join(errs...)
	}

	c.InitLog(cmd.ErrOrStderr())

	if flag := cmd.Flags().Lookup(FlagTimeout); flag != nil && flag.Changed {
		d, err := cmd.Flags().GetDuration(FlagTimeout)
		if err == nil {
			c.Server.IdleTimeout = d
			c.Server.MaxTimeout = d
		}
	}

	return c, nil
}

func EnvName(name string) string {
	name = strings.ToUpper(name)
	name = strings.ReplaceAll(name, "-", "_")
	return EnvPrefix + name
}
