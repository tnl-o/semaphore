package cmd

import (
	"fmt"
	"os"

	"github.com/semaphoreui/semaphore/util"
	"github.com/spf13/cobra"
)

func init() {
	runnerCmd.AddCommand(runnerUnregisterCmd)
}

func unregisterRunner() {
	util.ConfigInit(persistentFlags.configPath, persistentFlags.noConfig)

	taskPool := createRunnerJobPool()
	err := taskPool.Unregister()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to unregister runner: %v\n", err)
		os.Exit(1)
	}
}

var runnerUnregisterCmd = &cobra.Command{
	Use:   "unregister",
	Short: "Unregister runner from the server",
	Run: func(cmd *cobra.Command, args []string) {
		unregisterRunner()
	},
}
