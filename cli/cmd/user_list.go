package cmd

import (
	"fmt"
	"os"

	"github.com/semaphoreui/semaphore/db"
	"github.com/spf13/cobra"
)

func init() {
	userCmd.AddCommand(userListCmd)
}

var userListCmd = &cobra.Command{
	Use:   "list",
	Short: "Print all users",
	Run: func(cmd *cobra.Command, args []string) {
		store := createStore("")
		defer store.Close("")

		users, err := store.GetUsers(db.RetrieveQueryParams{})

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to get users: %v\n", err)
			os.Exit(1)
		}

		for _, user := range users {
			fmt.Println(user.Username)
		}
	},
}
