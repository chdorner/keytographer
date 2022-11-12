package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chdorner/keytographer/config"
	"github.com/chdorner/keytographer/renderer"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewRenderCommand() *cobra.Command {
	var ctx context.Context
	var configFile string
	var outDir string

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

			outDir, _ = cmd.Flags().GetString("out")
			if outDir == "" {
				outDir = strings.TrimSuffix(configFile, filepath.Ext(configFile))
				if outDir == "" {
					outDir = "output"
				}
				logrus.Infof("output directory not set, using %s", outDir)
			}

			return nil
		},

		Run: func(cmd *cobra.Command, args []string) {
			data, err := config.Load(configFile)
			if err != nil {
				logrus.WithField("error", err).Error("failed to load config file")
				os.Exit(1)
			}

			err = config.Validate(data)
			if err != nil {
				logrus.WithField("error", err).Error("configuration is invalid")
				os.Exit(1)
			}

			config, err := config.Parse(data)
			if err != nil {
				logrus.WithField("error", err).Error("failed to parse config")
				os.Exit(1)
			}

			renderer := renderer.NewRenderer()
			layers, err := renderer.RenderAllLayers(config)
			if err != nil {
				logrus.WithField("error", err).Error("failed to render layers")
				os.Exit(1)
			}

			err = os.MkdirAll(outDir, 0755)
			if err != nil {
				logrus.WithField("error", err).Errorf("failed to create output directory %s", outDir)
				os.Exit(1)
			}

			for idx, layer := range layers {
				path := filepath.Join(outDir, fmt.Sprintf("%d-%s.svg", (idx+1), layer.Name))
				err = os.WriteFile(path, layer.Svg, 0644)
				if err != nil {
					logrus.WithField("error", err).Errorf("failed to write layer svg to %s", path)
				}
			}
		},
	}

	fl := cmd.Flags()
	fl.StringP("out", "o", "", "Path to the output directory.")

	return cmd
}