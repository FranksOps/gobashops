package main

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var storageCmd = &cobra.Command{
	Use:   "storage",
	Short: "LVM, block storage, and mount triage",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("=== STORAGE TRIAGE ===")

		fmt.Println("\n[LVM Physical Volumes (pvs)]")
		if _, err := exec.LookPath("pvs"); err == nil {
			out, _ := exec.Command("pvs").CombinedOutput()
			fmt.Print(string(out))

			fmt.Println("\n[LVM Volume Groups (vgs)]")
			out, _ = exec.Command("vgs").CombinedOutput()
			fmt.Print(string(out))
		} else {
			fmt.Println("  LVM tools (pvs) not found.")
		}

		fmt.Println("\n[Block Devices (lsblk -f)]")
		out, _ := exec.Command("lsblk", "-f").CombinedOutput()
		fmt.Print(string(out))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(storageCmd)
}
