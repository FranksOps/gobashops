package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:   "service [name]",
	Short: "Layered troubleshooting for a specific service",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		svc := args[0]
		fmt.Printf("=== LAYERED TROUBLESHOOTING: %s ===\n", svc)

		fmt.Println("\n[Step 1 - Status]")
		out, _ := exec.Command("systemctl", "is-active", svc).CombinedOutput()
		status := strings.TrimSpace(string(out))
		fmt.Printf("  Active state: %s\n", status)
		if status != "active" {
			out, _ := exec.Command("systemctl", "status", svc, "--no-pager").CombinedOutput()
			lines := strings.Split(string(out), "\n")
			for i := 0; i < 5 && i < len(lines); i++ {
				fmt.Println("    " + lines[i])
			}
		}

		fmt.Println("\n[Step 2 - Logs (Last 10)]")
		out, _ = exec.Command("journalctl", "-u", svc, "-n", "10", "--no-pager").CombinedOutput()
		fmt.Println(string(out))

		fmt.Println("\n[Step 3 - Listening Ports]")
		out, _ = exec.Command("ss", "-tulnp").CombinedOutput()
		foundPort := false
		for _, line := range strings.Split(string(out), "\n") {
			if strings.Contains(line, svc) {
				fmt.Println("  " + line)
				foundPort = true
			}
		}
		if !foundPort {
			fmt.Println("  No listening ports detected for this process name.")
		}

		fmt.Println("\n[Step 5 - SELinux Audit]")
		if _, err := os.Stat("/sys/fs/selinux"); err == nil {
			out, err := exec.Command("ausearch", "-m", "AVC", "-ts", "recent").CombinedOutput()
			if err != nil {
				fmt.Println("  No recent AVC denials found.")
			} else {
				fmt.Println("  [!] RECENT SELINUX DENIALS DETECTED:")
				lines := strings.Split(string(out), "\n")
				for i := 0; i < 5 && i < len(lines); i++ {
					fmt.Println("    " + lines[i])
				}
			}
		} else {
			fmt.Println("  SELinux is disabled or not present.")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(serviceCmd)
}
