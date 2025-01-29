package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Set up judge credentials",
	Run: func(cmd *cobra.Command, args []string) {
		judge, _ := cmd.Flags().GetString("judge")
		var handle, cookie string

		fmt.Printf("Configuring %s:\n", judge)
		fmt.Print("Handle/Username: ")
		fmt.Scanln(&handle)
		fmt.Print("Session Cookie: ")
		fmt.Scanln(&cookie)

		viper.Set(fmt.Sprintf("judges.%s.handle", judge), handle)
		viper.Set(fmt.Sprintf("judges.%s.cookie", judge), cookie)
		viper.SafeWriteConfig()
		fmt.Println("✅ Configuration saved!")
	},
}

func init() {
	configureCmd.Flags().StringP("judge", "j", "", "Judge to configure")
	configureCmd.MarkFlagRequired("judge")
}