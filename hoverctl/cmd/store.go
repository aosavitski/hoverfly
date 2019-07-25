package cmd

import (
	"fmt"

	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
	"github.com/spf13/cobra"
)

var storeCmd = &cobra.Command{
	Use:   "store [path to simulation]",
	Short: "Store a simulation on the Hoverfly server",
	Long: `
Store a simulation on the Hoverfly server. The simulation JSON
will be written to the file with name provided.
	`,

	Run: func(cmd *cobra.Command, args []string) {
		checkTargetAndExit(target)

		checkArgAndExit(args, "You have not provided a file name of simulation", "store")

		_, err := wrapper.StoreSimulation(*target, args[0])
		handleIfError(err)

		fmt.Println("Successfully stored simulation on Hoverfly server to", args[0])
	},
}

func init() {
	RootCmd.AddCommand(storeCmd)
}
