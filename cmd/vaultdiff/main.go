// main.go is the entry point for the vaultdiff CLI tool.
// It wires together configuration, Vault clients, diffing, filtering,
// auditing, and output rendering into a cohesive command-line interface.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/youorg/vaultdiff/internal/vault"
)

const defaultInterval = 30 * time.Second

func main() {
	// Subcommands
	diffCmd := flag.NewFlagSet("diff", flag.ExitOnError)
	watchCmd := flag.NewFlagSet("watch", flag.ExitOnError)

	// diff flags
	diffConfig := diffCmd.String("config", "vaultdiff.yaml", "Path to config file")
	diffEnvA := diffCmd.String("env-a", "", "First environment name (required)")
	diffEnvB := diffCmd.String("env-b", "", "Second environment name (required)")
	diffPath := diffCmd.String("path", "", "Secret path to diff (required)")
	diffVersionA := diffCmd.Int("version-a", 0, "Version of secret in env-a (0 = latest)")
	diffVersionB := diffCmd.Int("version-b", 0, "Version of secret in env-b (0 = latest)")
	diffFormat := diffCmd.String("format", "text", "Output format: text or json")
	diffExport := diffCmd.String("export", "", "Export audit log to file (optional)")
	diffExportFmt := diffCmd.String("export-format", "json", "Export format: json or csv")
	diffFilter := diffCmd.String("filter", "", "Filter by change type: added, removed, changed, unchanged")
	diffPrefix := diffCmd.String("prefix", "", "Filter keys by prefix")
	diffExclude := diffCmd.String("exclude", "", "Comma-separated list of keys to exclude")

	// watch flags
	watchConfig := watchCmd.String("config", "vaultdiff.yaml", "Path to config file")
	watchEnv := watchCmd.String("env", "", "Environment name to watch (required)")
	watchPath := watchCmd.String("path", "", "Secret path to watch (required)")
	watchInterval := watchCmd.Duration("interval", defaultInterval, "Polling interval")
	watchFormat := watchCmd.String("format", "text", "Notification format: text or json")

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: vaultdiff <command> [flags]\n")
		fmt.Fprintf(os.Stderr, "Commands: diff, watch\n")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "diff":
		_ = diffCmd.Parse(os.Args[2:])
		runDiff(*diffConfig, *diffEnvA, *diffEnvB, *diffPath,
			*diffVersionA, *diffVersionB, *diffFormat,
			*diffExport, *diffExportFmt, *diffFilter, *diffPrefix, *diffExclude)

	case "watch":
		_ = watchCmd.Parse(os.Args[2:])
		runWatch(*watchConfig, *watchEnv, *watchPath, *watchInterval, *watchFormat)

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func runDiff(configPath, envA, envB, path string, versionA, versionB int,
	format, exportPath, exportFormat, filterType, prefix, exclude string) {

	if envA == "" || envB == "" || path == "" {
		fmt.Fprintln(os.Stderr, "Error: --env-a, --env-b, and --path are required for diff")
		os.Exit(1)
	}

	cfg, err := vault.LoadConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	clientA, err := vault.NewClient(cfg, envA)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create client for %s: %v\n", envA, err)
		os.Exit(1)
	}

	clientB, err := vault.NewClient(cfg, envB)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create client for %s: %v\n", envB, err)
		os.Exit(1)
	}

	secretsA, err := vault.GetSecretVersion(clientA, path, versionA)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to fetch secret from %s: %v\n", envA, err)
		os.Exit(1)
	}

	secretsB, err := vault.GetSecretVersion(clientB, path, versionB)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to fetch secret from %s: %v\n", envB, err)
		os.Exit(1)
	}

	changes := vault.DiffSecrets(secretsA, secretsB)

	// Apply optional filters
	if filterType != "" || prefix != "" || exclude != "" {
		changes = vault.FilterDiff(changes, vault.FilterOptions{
			ChangeType: filterType,
			KeyPrefix:  prefix,
			ExcludeKeys: splitCSV(exclude),
		})
	}

	vault.RenderDiff(os.Stdout, changes, format)

	// Record audit log
	log := vault.NewAuditLog()
	log.Record(path, envA, envB, changes)

	if exportPath != "" {
		if err := vault.ExportAuditLog(log, exportPath, exportFormat); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to export audit log: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "Audit log exported to %s\n", exportPath)
	}
}

func runWatch(configPath, env, path string, interval time.Duration, format string) {
	if env == "" || path == "" {
		fmt.Fprintln(os.Stderr, "Error: --env and --path are required for watch")
		os.Exit(1)
	}

	cfg, err := vault.LoadConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	client, err := vault.NewClient(cfg, env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create client for %s: %v\n", env, err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle OS signals for graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Fprintln(os.Stdout, "\nStopping watch...")
		cancel()
	}()

	events := vault.WatchSecret(ctx, client, path, interval)
	vault.HandleWatchEvents(ctx, events, os.Stdout, format)
}

// splitCSV splits a comma-separated string into a slice, ignoring empty entries.
func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	var out []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ',' {
			part := s[start:i]
			if part != "" {
				out = append(out, part)
			}
			start = i + 1
		}
	}
	return out
}
