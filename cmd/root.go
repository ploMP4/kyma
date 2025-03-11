package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/ploMP4/orama/internal/tui"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var rootCmd = &cobra.Command{
	Use: "orama <filename>",
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}

		if filepath.Ext(args[0]) != ".md" {
			return fmt.Errorf("Expected markdown file got: %v", args[0])
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := args[0]

		data, err := os.ReadFile(filename)
		if err != nil {
			return err
		}

		slides := strings.Split(string(data), "---")
		for _, slide := range slides {
			fmt.Println(strings.TrimSpace(slide))
		}

		p := tea.NewProgram(tui.New(slides), tea.WithAltScreen(), tea.WithMouseAllMotion())
		if _, err := p.Run(); err != nil {
			return err
		}

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
