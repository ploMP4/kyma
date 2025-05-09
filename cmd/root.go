package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/ploMP4/kyma/internal/tui"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var rootCmd = &cobra.Command{
	Use: "kyma <filename>",
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}

		if filepath.Ext(args[0]) != ".md" {
			return fmt.Errorf("expected markdown file got: %v", args[0])
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := args[0]

		data, err := os.ReadFile(filename)
		if err != nil {
			return err
		}

		root, err := parseSildes(string(data))
		if err != nil {
			return err
		}

		p := tea.NewProgram(tui.New(root), tea.WithAltScreen(), tea.WithMouseAllMotion())
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

func parseSildes(data string) (*tui.Slide, error) {
	slides := strings.Split(string(data), "----\n")

	rootSlide, properties := parseSlide(slides[0])
	p, err := tui.NewProperties(properties)
	if err != nil {
		return nil, err
	}

	root := &tui.Slide{
		Data:       rootSlide,
		Properties: p,
	}

	curr := root
	for _, slide := range slides[1:] {
		slide, properties := parseSlide(slide)
		p, err := tui.NewProperties(properties)
		if err != nil {
			return nil, err
		}

		curr.Next = &tui.Slide{
			Data:       slide,
			Prev:       curr,
			Properties: p,
		}
		curr = curr.Next
	}

	return root, nil
}

func parseSlide(s string) (slide, properties string) {
	slide = s

	if strings.HasPrefix(strings.TrimSpace(s), "---\n") {
		parts := strings.Split(s, "---\n")
		properties = parts[1]
		slide = parts[2]
	}

	return slide, properties
}
