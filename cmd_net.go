package main

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var netCmd = &cobra.Command{
	Use:   "net [host]",
	Short: "Layered troubleshooting for server reachability",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		host := args[0]
		fmt.Printf("=== NETWORK TRIAGE: %s ===\n", host)

		fmt.Println("\n[Layer 1 - Host Alive?]")
		err := exec.Command("ping", "-c", "2", "-W", "1", host).Run()
		if err == nil {
			fmt.Println("  [SUCCESS] Host responded to ICMP Ping.")
		} else {
			fmt.Println("  [FAILED] Host did not respond to ICMP Ping.")
		}

		fmt.Println("\n[Layer 2 - Routing Table]")
		out, _ := exec.Command("ip", "route", "get", host).CombinedOutput()
		fmt.Printf("  Route to host: %s\n", strings.TrimSpace(string(out)))

		port, _ := cmd.Flags().GetString("port")
		if port != "" {
			fmt.Printf("\n[Layer 3 - Port Reachability (%s)]\n", port)

			target := net.JoinHostPort(host, port)
			conn, err := net.DialTimeout("tcp", target, 2*time.Second)
			if err != nil {
				fmt.Printf("  [FAILED] Connection refused or timed out: %v\n", err)
			} else {
				fmt.Printf("  [SUCCESS] TCP Handshake completed. Port is open.\n")
				conn.Close()
			}
		} else {
			fmt.Println("\n[Layer 3 - Port Reachability]")
			fmt.Println("  Skip: No --port provided.")
		}

		return nil
	},
}

func init() {
	netCmd.Flags().StringP("port", "p", "", "Port to test (e.g., 22, 443)")
	rootCmd.AddCommand(netCmd)
}
