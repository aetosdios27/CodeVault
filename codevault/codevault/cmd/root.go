package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "codevault",
	Short: "Organize your competitive programming solutions",
}

// Define generateCmd so it isn't undefined
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate and fetch solutions",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: implement logic
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.AddCommand(configureCmd)
	rootCmd.AddCommand(generateCmd) // now defined
}
