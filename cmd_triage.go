package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
)

var triageCmd = &cobra.Command{
	Use:   "triage",
	Short: "General system slowness triage bundle",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("=== SYSTEM TRIAGE (Senior Admin Bundle) ===")

		load, _ := os.ReadFile("/proc/loadavg")
		fmt.Printf("[Load Average] %s\n", strings.TrimSpace(string(load)))

		meminfo, _ := os.ReadFile("/proc/meminfo")
		lines := strings.Split(string(meminfo), "\n")
		var memTotal, memAvail, swapTotal, swapFree string
		for _, line := range lines {
			if strings.HasPrefix(line, "MemTotal:") {
				memTotal = line
			}
			if strings.HasPrefix(line, "MemAvailable:") {
				memAvail = line
			}
			if strings.HasPrefix(line, "SwapTotal:") {
				swapTotal = line
			}
			if strings.HasPrefix(line, "SwapFree:") {
				swapFree = line
			}
		}
		fmt.Println("\n[Memory]")
		fmt.Printf("  %s\n  %s\n", strings.TrimSpace(memTotal), strings.TrimSpace(memAvail))
		fmt.Printf("  %s\n  %s\n", strings.TrimSpace(swapTotal), strings.TrimSpace(swapFree))

		fmt.Println("\n[Disk Space Pressure]")
		var stat syscall.Statfs_t
		syscall.Statfs("/", &stat)
		totalGB := (stat.Blocks * uint64(stat.Bsize)) / (1024 * 1024 * 1024)
		freeGB := (stat.Bavail * uint64(stat.Bsize)) / (1024 * 1024 * 1024)
		pct := float64(totalGB-freeGB) / float64(totalGB) * 100

		fmt.Printf("  Root FS (/): %.1f%% used (%d GB free)\n", pct, freeGB)
		if pct > 85.0 {
			fmt.Println("  [!] WARNING: Root partition is highly constrained.")
			fmt.Println("  [!] Checking for held-open deleted files (lsof | grep deleted)...")
			out, _ := exec.Command("bash", "-c", "lsof +L1").CombinedOutput()
			if len(out) > 0 {
				fmt.Println(string(out))
			}
		}

		fmt.Println("\n[Top 3 CPU Hogs]")
		out, _ := exec.Command("ps", "-eo", "pid,ppid,cmd,%mem,%cpu", "--sort=-%cpu").CombinedOutput()
		lines = strings.Split(string(out), "\n")
		for i := 0; i < 4 && i < len(lines); i++ {
			fmt.Println("  " + lines[i])
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(triageCmd)
}
