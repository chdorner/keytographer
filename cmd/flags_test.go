package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestParseConfigFlag(t *testing.T) {
	cmd := &cobra.Command{}
	addConfigFlag(cmd)

	err := cmd.ParseFlags([]string{})
	require.NoError(t, err)
	_, err = parseConfigFlag(cmd)
	require.ErrorContains(t, err, "missing path")

	err = cmd.ParseFlags([]string{"-c nothing"})
	require.NoError(t, err)
	_, err = parseConfigFlag(cmd)
	require.ErrorContains(t, err, "file does not exist")

	err = cmd.ParseFlags([]string{"-c", "testdata/config.yaml"})
	require.NoError(t, err)
	actual, err := parseConfigFlag(cmd)
	require.NoError(t, err)
	require.Equal(t, "testdata/config.yaml", actual)
}
