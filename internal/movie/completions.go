package movie

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func CompleteMovieName(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}

	infos, err := ListEmbedded()
	if err != nil {
		return []string{}, cobra.ShellCompDirectiveDefault
	}

	result := make([]string, 0, len(infos))
	for _, info := range infos {
		if strings.HasPrefix(info.Name, toComplete) {
			result = append(result, fmt.Sprintf(
				"%s\tduration=%s",
				info.Name,
				info.Duration.Round(time.Second),
			))
		}
	}

	if len(result) == 0 {
		return []string{"txt", "txt.gz"}, cobra.ShellCompDirectiveFilterFileExt
	} else {
		return result, cobra.ShellCompDirectiveDefault
	}
}
