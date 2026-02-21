package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gobashops",
	Short: "Senior Admin Triage and Operations CLI",
	Long:  `Gobashops encapsulates the Senior Linux Admin mental model.`,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
