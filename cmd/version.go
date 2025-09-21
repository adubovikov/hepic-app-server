package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long: `Show detailed version information for HEPIC App Server.

This command displays:
- Application version
- Go version
- Build information
- Runtime information
- Feature flags

Examples:
  hepic-app-server version
  hepic-app-server version --verbose
  hepic-app-server version --json`,
	Run: runVersion,
}

var (
	versionVerbose bool
	versionJSON    bool
)

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVar(&versionVerbose, "verbose", false, "Show verbose version information")
	versionCmd.Flags().BoolVar(&versionJSON, "json", false, "Output in JSON format")
}

func runVersion(cmd *cobra.Command, args []string) {
	version := VersionInfo{
		Application: "HEPIC App Server",
		Version:     "2.0.0",
		BuildTime:   time.Now().Format(time.RFC3339),
		GitCommit:   "unknown",
		GoVersion:    runtime.Version(),
		Platform:   fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		Compiler:   runtime.Compiler,
	}

	if versionVerbose {
		version.Features = []string{
			"ClickHouse Integration",
			"Structured JSON Logging",
			"JWT Authentication",
			"REST API",
			"Docker Support",
			"Health Monitoring",
			"Graceful Shutdown",
			"Request ID Tracking",
			"Error Recovery",
			"Performance Metrics",
		}
		version.Dependencies = []string{
			"github.com/labstack/echo/v4",
			"github.com/ClickHouse/clickhouse-go/v2",
			"github.com/spf13/cobra",
			"github.com/spf13/viper",
			"log/slog",
		}
	}

	if versionJSON {
		// Output JSON format
		jsonData, err := json.MarshalIndent(version, "", "  ")
		if err != nil {
			fmt.Printf("Error marshaling version: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(jsonData))
	} else {
		// Output human-readable format
		fmt.Printf("HEPIC App Server v%s\n", version.Version)
		fmt.Printf("Build: %s\n", version.BuildTime)
		fmt.Printf("Go: %s\n", version.GoVersion)
		fmt.Printf("Platform: %s\n", version.Platform)
		fmt.Printf("Compiler: %s\n", version.Compiler)

		if verbose {
			fmt.Println("\nFeatures:")
			for _, feature := range version.Features {
				fmt.Printf("  ✓ %s\n", feature)
			}

			fmt.Println("\nDependencies:")
			for _, dep := range version.Dependencies {
				fmt.Printf("  - %s\n", dep)
			}
		}
	}
}

type VersionInfo struct {
	Application  string   `json:"application"`
	Version      string   `json:"version"`
	BuildTime    string   `json:"build_time"`
	GitCommit    string   `json:"git_commit"`
	GoVersion    string   `json:"go_version"`
	Platform     string   `json:"platform"`
	Compiler     string   `json:"compiler"`
	Features     []string `json:"features,omitempty"`
	Dependencies []string `json:"dependencies,omitempty"`
}
