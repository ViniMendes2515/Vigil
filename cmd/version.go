package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version = "development"
	vers    bool
)

var cmdVersion = &cobra.Command{
	Use:   "version",
	Short: "Exibe a versão do Vigil",
	Long:  `Exibe a versão do Vigil`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}

func init() {
	rootCmd.Flags().BoolVarP(&vers, "version", "v", false, "Exibe versão do Vigil")

	originalRunE := rootCmd.RunE

	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		if vers {
			fmt.Println(Version)
			return nil
		}

		if originalRunE != nil {
			return originalRunE(cmd, args)
		}
		return nil
	}

	rootCmd.AddCommand(cmdVersion)
}
