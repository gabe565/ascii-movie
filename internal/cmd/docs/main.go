package main

import (
	"fmt"
	"github.com/gabe565/ascii-movie/cmd"
	"github.com/spf13/cobra/doc"
	"log"
	"os"
)

func main() {
	output := "./docs"

	if err := os.RemoveAll(output); err != nil {
		log.Fatal(fmt.Errorf("failed to remove existing dia: %w", err))
	}

	if err := os.MkdirAll(output, 0755); err != nil {
		log.Fatal(fmt.Errorf("failed to mkdir: %w", err))
	}

	rootCmd := cmd.NewCommand()
	if err := doc.GenMarkdownTree(rootCmd, output); err != nil {
		log.Fatal(fmt.Errorf("failed to generate markdown: %w", err))
	}
}
