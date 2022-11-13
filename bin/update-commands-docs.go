package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/chdorner/keytographer/cmd"
	"github.com/spf13/cobra"
)

var (
	markStartF = "<!-- usage:%s:start -->"
	markStopF  = "<!-- usage:%s:end -->"

	file = "docs/commands.md"
)

func main() {
	readmeb, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	readme := string(readmeb)

	for _, c := range cmd.RootCommand().Commands() {
		readme = replaceUsage(readme, c)
	}

	err = os.WriteFile(file, []byte(readme), 0644)
	if err != nil {
		panic(err)
	}
}

func replaceUsage(readme string, cmd *cobra.Command) string {
	markStart := fmt.Sprintf(markStartF, cmd.Use)
	markStop := fmt.Sprintf(markStopF, cmd.Use)

	startIdx := strings.Index(readme, markStart)
	stopIdx := strings.Index(readme, markStop) + len(markStop) + 1

	if startIdx != -1 && stopIdx != -1 {
		readme = readme[:startIdx] +
			markStart +
			"\n```\n" +
			getHelp(cmd) +
			"```\n" +
			markStop +
			"\n" +
			readme[stopIdx:]
	}

	return readme
}

func getHelp(cmd *cobra.Command) string {
	buf := bytes.NewBuffer(nil)

	cmd.SetOut(buf)
	err := cmd.Help()
	if err != nil {
		return ""
	}

	return buf.String()
}
