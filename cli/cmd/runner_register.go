package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/semaphoreui/semaphore/util"
	"github.com/spf13/cobra"
)

var runnerRegisterArgs struct {
	stdinRegistrationToken bool
}

func init() {
	runnerRegisterCmd.PersistentFlags().BoolVar(&runnerRegisterArgs.stdinRegistrationToken, "stdin-registration-token", false, "Read registration token from stdin")
	runnerCmd.AddCommand(runnerRegisterCmd)
}

func initRunnerRegistrationToken() {
	if !runnerRegisterArgs.stdinRegistrationToken {
		return
	}

	tokenBytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to read token from stdin: %v\n", err)
		os.Exit(1)
	}

	if len(tokenBytes) == 0 {
		fmt.Fprintf(os.Stderr, "Error: Empty token\n")
		os.Exit(1)
	}

	util.Config.Runner.RegistrationToken = strings.TrimSpace(string(tokenBytes))
}

func registerRunner() {

	configFile := util.ConfigInit(persistentFlags.configPath, persistentFlags.noConfig)

	initRunnerRegistrationToken()

	taskPool := createRunnerJobPool()

	err := taskPool.Register(configFile)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to register runner: %v\n", err)
		os.Exit(1)
	}
}

var runnerRegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "Register runner on the server",
	Run: func(cmd *cobra.Command, args []string) {
		registerRunner()
	},
}
