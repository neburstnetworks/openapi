package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/neburstnetworks/openapi/sdk/go/neburst"
)

func runCommand(args []string) {
	if len(args) == 0 {
		printUsage()
		return
	}

	switch args[0] {
	case "account":
		cmdAccount(args[1:])
	case "instance":
		withClient(func(c *neburst.Client) { cmdInstance(c, args[1:]) })
	case "bare-metal":
		withClient(func(c *neburst.Client) { cmdBareMetal(c, args[1:]) })
	case "balance":
		withClient(func(c *neburst.Client) { cmdBalance(c) })
	case "invoices":
		page := 1
		if len(args) > 1 {
			page, _ = strconv.Atoi(args[1])
		}
		withClient(func(c *neburst.Client) { cmdInvoices(c, page) })
	case "invoice":
		if len(args) < 2 {
			fatal("usage: neburst-cli invoice <uuid>")
		}
		withClient(func(c *neburst.Client) { cmdInvoice(c, args[1]) })
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", args[0])
		printUsage()
		os.Exit(1)
	}
}

func withClient(fn func(*neburst.Client)) {
	cfg := loadConfig()
	acct, ok := cfg.activeAccount()
	if !ok {
		fatal("no active account. Run: neburst-cli account add <alias>")
	}
	c := neburst.NewClient(acct.BaseURL, acct.KeyID, acct.Secret)
	fn(c)
}

func cmdAccount(args []string) {
	if len(args) == 0 {
		fmt.Println("usage: neburst-cli account <list|add|remove|switch|current>")
		return
	}
	cfg := loadConfig()

	switch args[0] {
	case "list", "ls":
		if len(cfg.Accounts) == 0 {
			fmt.Println("No accounts configured.")
			return
		}
		for _, name := range cfg.sortedNames() {
			acct := cfg.Accounts[name]
			marker := "  "
			if name == cfg.Active {
				marker = "* "
			}
			fmt.Printf("%s%-16s  %s  %s\n", marker, name, acct.BaseURL, acct.KeyID)
		}

	case "add":
		if len(args) < 2 {
			fatal("usage: neburst-cli account add <alias>")
		}
		name := args[1]
		reader := bufio.NewReader(os.Stdin)
		baseURL := promptLine(reader, "API Base URL [https://api.neburst.com]: ")
		if baseURL == "" {
			baseURL = "https://api.neburst.com"
		}
		keyInput := promptLine(reader, "API Key (base64 combined key or Key ID): ")
		keyID, secret, ok := parseCombinedKey(keyInput)
		if !ok {
			keyID = keyInput
			secret = promptLine(reader, "API Secret: ")
		} else {
			fmt.Printf("  Parsed combined key: %s\n", keyID)
		}
		if keyID == "" || secret == "" {
			fatal("key_id and secret are required")
		}
		cfg.addAccount(name, Account{BaseURL: baseURL, KeyID: keyID, Secret: secret})
		if cfg.Active == "" || len(cfg.Accounts) == 1 {
			cfg.Active = name
		}
		must(cfg.save())
		fmt.Printf("✓ Account '%s' added", name)
		if cfg.Active == name {
			fmt.Print(" (active)")
		}
		fmt.Println()

	case "remove", "rm":
		if len(args) < 2 {
			fatal("usage: neburst-cli account remove <alias>")
		}
		name := args[1]
		if _, ok := cfg.Accounts[name]; !ok {
			fatal("account not found: " + name)
		}
		cfg.removeAccount(name)
		must(cfg.save())
		fmt.Printf("✓ Account '%s' removed\n", name)

	case "switch", "use":
		if len(args) < 2 {
			fatal("usage: neburst-cli account switch <alias>")
		}
		name := args[1]
		if !cfg.setActive(name) {
			fatal("account not found: " + name)
		}
		must(cfg.save())
		fmt.Printf("✓ Switched to '%s'\n", name)

	case "current":
		if cfg.Active == "" {
			fmt.Println("No active account.")
			return
		}
		acct := cfg.Accounts[cfg.Active]
		fmt.Printf("  Alias:    %s\n", cfg.Active)
		fmt.Printf("  Base URL: %s\n", acct.BaseURL)
		fmt.Printf("  Key ID:   %s\n", acct.KeyID)
		fmt.Printf("  Secret:   %s\n", maskSecret(acct.Secret))

	default:
		fmt.Fprintf(os.Stderr, "unknown account subcommand: %s\n", args[0])
	}
}

func cmdInstance(c *neburst.Client, args []string) {
	if len(args) == 0 {
		fmt.Println("usage: neburst-cli instance <list|get|status|traffic|power|metrics> [args...]")
		return
	}
	switch args[0] {
	case "list", "ls":
		page := 1
		if len(args) > 1 {
			page, _ = strconv.Atoi(args[1])
		}
		result, err := c.ListInstances(neburst.WithPage(page))
		exitOnErr(err)
		printJSON(result)
	case "get":
		requireArg(args, 1, "usage: neburst-cli instance get <uuid>")
		inst, err := c.GetInstance(args[1])
		exitOnErr(err)
		printJSON(inst)
	case "status":
		requireArg(args, 1, "usage: neburst-cli instance status <uuid>")
		s, err := c.GetInstanceStatus(args[1])
		exitOnErr(err)
		printJSON(s)
	case "traffic":
		requireArg(args, 1, "usage: neburst-cli instance traffic <uuid>")
		t, err := c.GetInstanceTraffic(args[1])
		exitOnErr(err)
		printJSON(t)
	case "power":
		requireArg(args, 2, "usage: neburst-cli instance power <uuid> <start|stop|restart>")
		exitOnErr(c.CloudPowerAction(args[1], args[2]))
		fmt.Println("✓ Success")
	case "metrics":
		requireArg(args, 1, "usage: neburst-cli instance metrics <uuid>")
		m, err := c.GetCloudMetrics(args[1])
		exitOnErr(err)
		printJSON(m)
	default:
		fmt.Fprintf(os.Stderr, "unknown instance subcommand: %s\n", args[0])
	}
}

func cmdBareMetal(c *neburst.Client, args []string) {
	if len(args) == 0 {
		fmt.Println("usage: neburst-cli bare-metal <list|get|status|traffic|power|metrics|profiles|rescue-profiles|rebuild|rescue> [args...]")
		return
	}
	switch args[0] {
	case "list", "ls":
		page := 1
		if len(args) > 1 {
			page, _ = strconv.Atoi(args[1])
		}
		result, err := c.ListBareMetalInstances(neburst.WithPage(page))
		exitOnErr(err)
		printJSON(result)
	case "get":
		requireArg(args, 1, "usage: neburst-cli bare-metal get <uuid>")
		inst, err := c.GetBareMetalInstance(args[1])
		exitOnErr(err)
		printJSON(inst)
	case "status":
		requireArg(args, 1, "usage: neburst-cli bare-metal status <uuid>")
		s, err := c.GetBareMetalStatus(args[1])
		exitOnErr(err)
		printJSON(s)
	case "traffic":
		requireArg(args, 1, "usage: neburst-cli bare-metal traffic <uuid>")
		t, err := c.GetBareMetalTraffic(args[1])
		exitOnErr(err)
		printJSON(t)
	case "power":
		requireArg(args, 2, "usage: neburst-cli bare-metal power <uuid> <power-on|power-off|power-cycle|power-reset>")
		exitOnErr(c.BareMetalPowerAction(args[1], args[2]))
		fmt.Println("✓ Success")
	case "metrics":
		requireArg(args, 1, "usage: neburst-cli bare-metal metrics <uuid>")
		m, err := c.GetBareMetalMetrics(args[1])
		exitOnErr(err)
		printJSON(m)
	case "profiles":
		requireArg(args, 1, "usage: neburst-cli bare-metal profiles <uuid>")
		p, err := c.GetReinstallProfiles(args[1])
		exitOnErr(err)
		printJSON(p)
	case "rescue-profiles":
		requireArg(args, 1, "usage: neburst-cli bare-metal rescue-profiles <uuid>")
		p, err := c.GetRescueProfiles(args[1])
		exitOnErr(err)
		printJSON(p)
	case "rebuild":
		requireArg(args, 2, "usage: neburst-cli bare-metal rebuild <uuid> <profile_id>")
		pid, _ := strconv.Atoi(args[2])
		exitOnErr(c.RebuildInstance(args[1], pid))
		fmt.Println("✓ Rebuild initiated")
	case "rescue":
		requireArg(args, 2, "usage: neburst-cli bare-metal rescue <uuid> <profile_id>")
		pid, _ := strconv.Atoi(args[2])
		exitOnErr(c.RescueInstance(args[1], pid))
		fmt.Println("✓ Rescue initiated")
	default:
		fmt.Fprintf(os.Stderr, "unknown bare-metal subcommand: %s\n", args[0])
	}
}

func cmdBalance(c *neburst.Client) {
	b, err := c.GetBalance()
	exitOnErr(err)
	fmt.Printf("Available: %.2f %s\nLocked:    %.2f %s\n", b.Available, b.Currency, b.Locked, b.Currency)
}

func cmdInvoices(c *neburst.Client, page int) {
	if page < 1 {
		page = 1
	}
	result, err := c.ListInvoices(neburst.WithPage(page))
	exitOnErr(err)
	printJSON(result)
}

func cmdInvoice(c *neburst.Client, id string) {
	inv, err := c.GetInvoice(id)
	exitOnErr(err)
	printJSON(inv)
}

func requireArg(args []string, n int, msg string) {
	if len(args) <= n {
		fatal(msg)
	}
}

func exitOnErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "✗ Error: %v\n", err)
		os.Exit(1)
	}
}

func fatal(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func must(err error) {
	if err != nil {
		fatal(err.Error())
	}
}

func promptLine(reader *bufio.Reader, msg string) string {
	fmt.Print(msg)
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func printUsage() {
	fmt.Print(`Neburst OpenAPI CLI

Usage:
  neburst-cli                                     Interactive mode
  neburst-cli <command> [args...]                  Command mode

Account:
  account list                                     List all accounts
  account add <alias>                              Add a new account
  account remove <alias>                           Remove an account
  account switch <alias>                           Set active account
  account current                                  Show active account

Cloud Instance:
  instance list [page]                             List instances
  instance get <uuid>                              Get instance details
  instance status <uuid>                           Get power status
  instance traffic <uuid>                          Get traffic usage
  instance power <uuid> <start|stop|restart>       Power action
  instance metrics <uuid>                          Get metrics

Bare Metal:
  bare-metal list [page]                           List instances
  bare-metal get <uuid>                            Get instance details
  bare-metal status <uuid>                         Get power status
  bare-metal traffic <uuid>                        Get traffic usage
  bare-metal power <uuid> <action>                 Power action
  bare-metal metrics <uuid>                        Get metrics
  bare-metal profiles <uuid>                       List OS profiles
  bare-metal rescue-profiles <uuid>                List rescue profiles
  bare-metal rebuild <uuid> <profile_id>           Rebuild instance
  bare-metal rescue <uuid> <profile_id>            Enter rescue mode

Billing:
  balance                                          Get account balance
  invoices [page]                                  List invoices
  invoice <uuid>                                   Get invoice details
`)
}
