package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"

	"github.com/ploMP4/kyma/internal/tui"
)

var watch bool

func init() {
	rootCmd.Flags().BoolVarP(&watch, "watch", "w", false, "Watch for changes in the input file")
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

		root, err := parseSlides(string(data))
		if err != nil {
			return err
		}

		p := tea.NewProgram(tui.New(root), tea.WithAltScreen(), tea.WithMouseAllMotion())

		if watch {
			watcher, err := fsnotify.NewWatcher()
			if err != nil {
				return fmt.Errorf("failed to create file watcher: %w", err)
			}
			defer watcher.Close()

			absPath, err := filepath.Abs(filename)
			if err != nil {
				return fmt.Errorf("failed to get absolute path: %w", err)
			}

			if err := watcher.Add(filepath.Dir(absPath)); err != nil {
				return fmt.Errorf("failed to watch directory: %w", err)
			}

			go func() {
				var debounceTimer *time.Timer

				for {
					select {
					case event, ok := <-watcher.Events:
						if !ok {
							return
						}

						if event.Name == absPath || event.Name == filename ||
							strings.HasSuffix(event.Name, "~") ||
							strings.HasPrefix(event.Name, absPath+".") {
							if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
								if debounceTimer != nil {
									debounceTimer.Stop()
								}
								debounceTimer = time.NewTimer(100 * time.Millisecond)

								go func() {
									<-debounceTimer.C
									data, err := os.ReadFile(filename)
									if err != nil {
										fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
										return
									}

									newRoot, err := parseSlides(string(data))
									if err != nil {
										fmt.Fprintf(os.Stderr, "Error parsing slides: %v\n", err)
										return
									}
								}()
							}
						}
					case err, ok := <-watcher.Errors:
						if !ok {
							return
						}
						fmt.Fprintf(os.Stderr, "Error watching file: %v\n", err)
					}
				}
			}()
		}

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

func parseSlides(data string) (*tui.Slide, error) {
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
