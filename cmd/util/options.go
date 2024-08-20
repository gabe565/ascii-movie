package util

import "github.com/spf13/cobra"

type Option func(cmd *cobra.Command)

const (
	VersioonKey = "version"
	CommitKey   = "commit"
)

func WithVersion(version string) Option {
	return func(cmd *cobra.Command) {
		if cmd.Annotations == nil {
			cmd.Annotations = make(map[string]string, 2)
		}
		cmd.Annotations[VersioonKey] = version
		cmd.Version, cmd.Annotations[CommitKey] = buildVersion(version)
		cmd.InitDefaultVersionFlag()
	}
}
