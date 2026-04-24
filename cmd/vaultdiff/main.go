package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/youorg/vaultdiff/internal/vault"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: vaultdiff <diff|watch|snapshot> [options]")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "diff":
		runDiff(os.Args[2:])
	case "watch":
		runWatch(os.Args[2:])
	case "snapshot":
		runSnapshot(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func runDiff(args []string) {
	fs := flag.NewFlagSet("diff", flag.ExitOnError)
	configFile := fs.String("config", "vault.yaml", "path to config file")
	envA := fs.String("env-a", "", "first environment")
	envB := fs.String("env-b", "", "second environment")
	secretPath := fs.String("path", "", "secret path")
	filter := fs.String("filter", "", "comma-separated change types to include")
	_ = fs.Parse(args)

	cfg, err := vault.LoadConfig(*configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}

	envCfgA, err := vault.GetEnvironment(cfg, *envA)
	if err != nil {
		fmt.Fprintf(os.Stderr, "env-a error: %v\n", err)
		os.Exit(1)
	}
	envCfgB, err := vault.GetEnvironment(cfg, *envB)
	if err != nil {
		fmt.Fprintf(os.Stderr, "env-b error: %v\n", err)
		os.Exit(1)
	}

	clientA, _ := vault.NewClient(envCfgA.Address, envCfgA.Token)
	clientB, _ := vault.NewClient(envCfgB.Address, envCfgB.Token)

	dataA, _ := vault.GetSecretVersion(clientA, *secretPath, 0)
	dataB, _ := vault.GetSecretVersion(clientB, *secretPath, 0)

	changes := vault.DiffSecrets(dataA, dataB)
	if *filter != "" {
		changes = vault.FilterDiff(changes, vault.FilterOptions{ChangeTypes: splitCSV(*filter)})
	}
	vault.RenderDiff(os.Stdout, changes)
}

func runWatch(args []string) {
	fs := flag.NewFlagSet("watch", flag.ExitOnError)
	addr := fs.String("addr", "", "vault address")
	token := fs.String("token", "", "vault token")
	path := fs.String("path", "", "secret path")
	interval := fs.Duration("interval", 30*time.Second, "poll interval")
	format := fs.String("format", "text", "output format: text|json")
	_ = fs.Parse(args)

	client, err := vault.NewClient(*addr, *token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "client error: %v\n", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	events := vault.WatchSecret(ctx, client, *path, *interval)
	vault.HandleWatchEvents(os.Stdout, events, *format)
}

func runSnapshot(args []string) {
	fs := flag.NewFlagSet("snapshot", flag.ExitOnError)
	addr := fs.String("addr", "", "vault address")
	token := fs.String("token", "", "vault token")
	path := fs.String("path", "", "secret path")
	env := fs.String("env", "default", "environment label")
	output := fs.String("output", "snapshot.json", "output file path")
	version := fs.Int("version", 0, "secret version (0 = latest)")
	_ = fs.Parse(args)

	client, err := vault.NewClient(*addr, *token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "client error: %v\n", err)
		os.Exit(1)
	}

	data, err := vault.GetSecretVersion(client, *path, *version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "secret fetch error: %v\n", err)
		os.Exit(1)
	}

	snap := vault.TakeSnapshot(*env, *path, *version, data)
	if err := vault.SaveSnapshot(snap, *output); err != nil {
		fmt.Fprintf(os.Stderr, "save snapshot error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("snapshot saved to %s\n", *output)
}

func splitCSV(s string) []string {
	var out []string
	for _, part := range strings.Split(s, ",") {
		if t := strings.TrimSpace(part); t != "" {
			out = append(out, t)
		}
	}
	return out
}
