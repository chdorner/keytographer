package cmd

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/chdorner/keymap-render/internal/keymap"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewRenderCommand() *cobra.Command {
	var ctx context.Context
	var configFile string
	var outFile string

	cmd := &cobra.Command{
		Use:   "render",
		Short: "Render keymap configuration to a SVG file.",

		PreRunE: func(cmd *cobra.Command, args []string) error {
			ctx = createContext(cmd.Flags())
			configureLogging(ctx)

			configFile, _ = cmd.Flags().GetString("config")
			if configFile == "" {
				return errors.New("missing path to keymap configuration file")
			}
			_, err := os.Stat(configFile)
			if err != nil {
				return errors.New("specified keymap configuration file does not exist")
			}

			outFile, _ = cmd.Flags().GetString("out")
			if outFile == "" {
				base := strings.TrimSuffix(configFile, filepath.Ext(configFile))
				if base == "" {
					base = "output"
				}
				outFile = fmt.Sprintf("%s.svg", base)
				logrus.Debugf("output file not set, using %s", outFile)
			}

			if configFile == outFile {
				return errors.New("input and output file are the same")
			}

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := keymap.Parse(configFile)
			if err != nil {
				return err
			}

			renderer := keymap.NewRenderer()
			svg := renderer.Render(config)

			return ioutil.WriteFile(outFile, svg, 0644)
		},
	}

	fl := cmd.Flags()
	fl.StringP("out", "o", "", "Path to output file.")

	return cmd
}