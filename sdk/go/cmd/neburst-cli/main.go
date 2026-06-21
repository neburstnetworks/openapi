package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/neburstnetworks/openapi/sdk/go/neburst"
)

var (
	client    *neburst.Client
	cfg       *Config
	activeName string
)

func main() {
	if len(os.Args) > 1 {
		runCommand(os.Args[1:])
		return
	}
	interactiveMode()
}

func interactiveMode() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("╔══════════════════════════════════════╗")
	fmt.Println("║       Neburst OpenAPI CLI (Go)       ║")
	fmt.Println("╚══════════════════════════════════════╝")
	fmt.Println()

	cfg = loadConfig()

	if acct, ok := cfg.activeAccount(); ok {
		activeName = cfg.Active
		client = neburst.NewClient(acct.BaseURL, acct.KeyID, acct.Secret)
		fmt.Printf("✓ Connected as [%s] (%s)\n\n", activeName, acct.BaseURL)
	} else if len(cfg.Accounts) > 0 {
		fmt.Println("No active account. Use '94' to switch.")
		fmt.Println()
	} else {
		fmt.Println("No accounts found. Let's set up your first account.")
		fmt.Println()
		interactiveAddAccount(reader)
	}

	for {
		printMenu()
		choice := prompt(reader, fmt.Sprintf("\n[%s] > ", activeName))
		fmt.Println()

		switch choice {
		case "1":
			requireClient(func() { cmdListCloudInstances(reader) })
		case "2":
			requireClient(func() { cmdGetCloudInstance(reader) })
		case "3":
			requireClient(func() { cmdGetCloudStatus(reader) })
		case "4":
			requireClient(func() { cmdGetCloudTraffic(reader) })
		case "5":
			requireClient(func() { cmdCloudPower(reader) })
		case "6":
			requireClient(func() { cmdCloudMetrics(reader) })
		case "11":
			requireClient(func() { cmdListBareMetalInstances(reader) })
		case "12":
			requireClient(func() { cmdGetBareMetalInstance(reader) })
		case "13":
			requireClient(func() { cmdGetBareMetalStatus(reader) })
		case "14":
			requireClient(func() { cmdGetBareMetalTraffic(reader) })
		case "15":
			requireClient(func() { cmdBareMetalPower(reader) })
		case "16":
			requireClient(func() { cmdBareMetalMetrics(reader) })
		case "17":
			requireClient(func() { cmdBareMetalProfiles(reader) })
		case "18":
			requireClient(func() { cmdBareMetalRescueProfiles(reader) })
		case "19":
			requireClient(func() { cmdBareMetalRebuild(reader) })
		case "20":
			requireClient(func() { cmdBareMetalRescue(reader) })
		case "31":
			requireClient(func() { cmdGetBalance() })
		case "32":
			requireClient(func() { cmdListInvoicesInteractive(reader) })
		case "33":
			requireClient(func() { cmdGetInvoiceInteractive(reader) })
		case "91":
			interactiveListAccounts()
		case "92":
			interactiveAddAccount(reader)
		case "93":
			interactiveRemoveAccount(reader)
		case "94":
			interactiveSwitchAccount(reader)
		case "95":
			interactiveCurrentAccount()
		case "0", "q", "quit", "exit":
			fmt.Println("Bye!")
			return
		default:
			fmt.Println("Unknown command")
		}
		fmt.Println()
	}
}

func requireClient(fn func()) {
	if client == nil {
		fmt.Println("✗ No active account. Use '92' to add or '94' to switch.")
		return
	}
	fn()
}

func switchToAccount(name string) {
	acct := cfg.Accounts[name]
	cfg.Active = name
	activeName = name
	client = neburst.NewClient(acct.BaseURL, acct.KeyID, acct.Secret)
	_ = cfg.save()
}

func interactiveAddAccount(r *bufio.Reader) {
	name := prompt(r, "Account alias: ")
	if name == "" {
		fmt.Println("✗ Alias cannot be empty")
		return
	}
	baseURL := prompt(r, "API Base URL [https://api.neburst.com]: ")
	if baseURL == "" {
		baseURL = "https://api.neburst.com"
	}
	keyInput := prompt(r, "API Key (base64 combined key or Key ID): ")
	keyID, secret, ok := parseCombinedKey(keyInput)
	if !ok {
		keyID = keyInput
		secret = prompt(r, "API Secret: ")
	} else {
		fmt.Printf("  Parsed combined key: %s\n", keyID)
	}
	if keyID == "" || secret == "" {
		fmt.Println("✗ Key ID and Secret are required")
		return
	}
	cfg.addAccount(name, Account{BaseURL: baseURL, KeyID: keyID, Secret: secret})
	switchToAccount(name)
	fmt.Printf("✓ Account '%s' added and activated\n", name)
}

func interactiveRemoveAccount(r *bufio.Reader) {
	if len(cfg.Accounts) == 0 {
		fmt.Println("No accounts to remove.")
		return
	}
	interactiveListAccounts()
	name := prompt(r, "Account alias to remove: ")
	if _, ok := cfg.Accounts[name]; !ok {
		fmt.Println("✗ Account not found")
		return
	}
	cfg.removeAccount(name)
	if name == activeName {
		activeName = ""
		client = nil
	}
	_ = cfg.save()
	fmt.Printf("✓ Account '%s' removed\n", name)
}

func interactiveSwitchAccount(r *bufio.Reader) {
	if len(cfg.Accounts) == 0 {
		fmt.Println("No accounts. Use '92' to add one.")
		return
	}
	interactiveListAccounts()
	name := prompt(r, "Switch to: ")
	if _, ok := cfg.Accounts[name]; !ok {
		fmt.Println("✗ Account not found")
		return
	}
	switchToAccount(name)
	fmt.Printf("✓ Switched to '%s'\n", name)
}

func interactiveListAccounts() {
	if len(cfg.Accounts) == 0 {
		fmt.Println("No accounts configured.")
		return
	}
	for _, name := range cfg.sortedNames() {
		acct := cfg.Accounts[name]
		marker := "  "
		if name == activeName {
			marker = "* "
		}
		fmt.Printf("%s%-16s  %s  %s\n", marker, name, acct.BaseURL, acct.KeyID)
	}
}

func interactiveCurrentAccount() {
	if activeName == "" {
		fmt.Println("No active account.")
		return
	}
	acct := cfg.Accounts[activeName]
	fmt.Printf("  Alias:    %s\n", activeName)
	fmt.Printf("  Base URL: %s\n", acct.BaseURL)
	fmt.Printf("  Key ID:   %s\n", acct.KeyID)
	fmt.Printf("  Secret:   %s\n", maskSecret(acct.Secret))
}

func printMenu() {
	fmt.Println("─── Cloud Instance ───────────────────")
	fmt.Println("  1  List instances")
	fmt.Println("  2  Get instance")
	fmt.Println("  3  Get status")
	fmt.Println("  4  Get traffic")
	fmt.Println("  5  Power action (start/stop/restart)")
	fmt.Println("  6  Metrics")
	fmt.Println("─── Bare Metal ───────────────────────")
	fmt.Println("  11 List instances")
	fmt.Println("  12 Get instance")
	fmt.Println("  13 Get status")
	fmt.Println("  14 Get traffic")
	fmt.Println("  15 Power action")
	fmt.Println("  16 Metrics")
	fmt.Println("  17 OS profiles")
	fmt.Println("  18 Rescue profiles")
	fmt.Println("  19 Rebuild")
	fmt.Println("  20 Rescue")
	fmt.Println("─── Billing ──────────────────────────")
	fmt.Println("  31 Balance")
	fmt.Println("  32 List invoices")
	fmt.Println("  33 Get invoice")
	fmt.Println("─── Account ──────────────────────────")
	fmt.Println("  91 List accounts")
	fmt.Println("  92 Add account")
	fmt.Println("  93 Remove account")
	fmt.Println("  94 Switch account")
	fmt.Println("  95 Current account info")
	fmt.Println("──────────────────────────────────────")
	fmt.Println("  0  Exit")
}

func prompt(reader *bufio.Reader, msg string) string {
	fmt.Print(msg)
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func printJSON(v any) {
	data, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(data))
}

func printErr(err error) {
	fmt.Printf("✗ Error: %v\n", err)
}

// ── Cloud Instance ──

func cmdListCloudInstances(r *bufio.Reader) {
	page, _ := strconv.Atoi(prompt(r, "Page [1]: "))
	if page < 1 {
		page = 1
	}
	result, err := client.ListInstances(neburst.WithPage(page))
	if err != nil {
		printErr(err)
		return
	}
	fmt.Printf("Total: %d, Page: %d/%d\n\n", result.Total, result.Page, (result.Total+result.PageSize-1)/maxInt(result.PageSize, 1))
	for _, inst := range result.Items {
		fmt.Printf("  %-36s  %-6s  %-10s  %s\n", inst.UUID, inst.Type, inst.Status, inst.Name)
	}
}

func cmdGetCloudInstance(r *bufio.Reader) {
	id := prompt(r, "Instance UUID: ")
	inst, err := client.GetInstance(id)
	if err != nil {
		printErr(err)
		return
	}
	printJSON(inst)
}

func cmdGetCloudStatus(r *bufio.Reader) {
	id := prompt(r, "Instance UUID: ")
	status, err := client.GetInstanceStatus(id)
	if err != nil {
		printErr(err)
		return
	}
	printJSON(status)
}

func cmdGetCloudTraffic(r *bufio.Reader) {
	id := prompt(r, "Instance UUID: ")
	traffic, err := client.GetInstanceTraffic(id)
	if err != nil {
		printErr(err)
		return
	}
	printJSON(traffic)
}

func cmdCloudPower(r *bufio.Reader) {
	id := prompt(r, "Instance UUID: ")
	action := prompt(r, "Action (start/stop/restart): ")
	err := client.CloudPowerAction(id, action)
	if err != nil {
		printErr(err)
		return
	}
	fmt.Println("✓ Success")
}

func cmdCloudMetrics(r *bufio.Reader) {
	id := prompt(r, "Instance UUID: ")
	m, err := client.GetCloudMetrics(id)
	if err != nil {
		printErr(err)
		return
	}
	printJSON(m)
}

// ── Bare Metal ──

func cmdListBareMetalInstances(r *bufio.Reader) {
	page, _ := strconv.Atoi(prompt(r, "Page [1]: "))
	if page < 1 {
		page = 1
	}
	result, err := client.ListBareMetalInstances(neburst.WithPage(page))
	if err != nil {
		printErr(err)
		return
	}
	fmt.Printf("Total: %d, Page: %d\n\n", result.Total, result.Page)
	for _, inst := range result.Items {
		fmt.Printf("  %-36s  %-10s  %s\n", inst.UUID, inst.Status, inst.Name)
	}
}

func cmdGetBareMetalInstance(r *bufio.Reader) {
	id := prompt(r, "Instance UUID: ")
	inst, err := client.GetBareMetalInstance(id)
	if err != nil {
		printErr(err)
		return
	}
	printJSON(inst)
}

func cmdGetBareMetalStatus(r *bufio.Reader) {
	id := prompt(r, "Instance UUID: ")
	status, err := client.GetBareMetalStatus(id)
	if err != nil {
		printErr(err)
		return
	}
	printJSON(status)
}

func cmdGetBareMetalTraffic(r *bufio.Reader) {
	id := prompt(r, "Instance UUID: ")
	traffic, err := client.GetBareMetalTraffic(id)
	if err != nil {
		printErr(err)
		return
	}
	printJSON(traffic)
}

func cmdBareMetalPower(r *bufio.Reader) {
	id := prompt(r, "Instance UUID: ")
	action := prompt(r, "Action (power-on/power-off/power-cycle/power-reset): ")
	err := client.BareMetalPowerAction(id, action)
	if err != nil {
		printErr(err)
		return
	}
	fmt.Println("✓ Success")
}

func cmdBareMetalMetrics(r *bufio.Reader) {
	id := prompt(r, "Instance UUID: ")
	m, err := client.GetBareMetalMetrics(id)
	if err != nil {
		printErr(err)
		return
	}
	printJSON(m)
}

func cmdBareMetalProfiles(r *bufio.Reader) {
	id := prompt(r, "Instance UUID: ")
	profiles, err := client.GetReinstallProfiles(id)
	if err != nil {
		printErr(err)
		return
	}
	for _, p := range profiles {
		fmt.Printf("  [%d] %-30s  %s\n", p.ID, p.Name, p.Category)
	}
}

func cmdBareMetalRescueProfiles(r *bufio.Reader) {
	id := prompt(r, "Instance UUID: ")
	profiles, err := client.GetRescueProfiles(id)
	if err != nil {
		printErr(err)
		return
	}
	for _, p := range profiles {
		fmt.Printf("  [%d] %-30s  %s\n", p.ID, p.Name, p.Category)
	}
}

func cmdBareMetalRebuild(r *bufio.Reader) {
	id := prompt(r, "Instance UUID: ")
	pid, _ := strconv.Atoi(prompt(r, "Profile ID: "))
	hostname := prompt(r, "Hostname (optional): ")
	var opts []neburst.RebuildOption
	if hostname != "" {
		opts = append(opts, neburst.WithHostname(hostname))
	}
	err := client.RebuildInstance(id, pid, opts...)
	if err != nil {
		printErr(err)
		return
	}
	fmt.Println("✓ Rebuild initiated")
}

func cmdBareMetalRescue(r *bufio.Reader) {
	id := prompt(r, "Instance UUID: ")
	pid, _ := strconv.Atoi(prompt(r, "Profile ID: "))
	err := client.RescueInstance(id, pid)
	if err != nil {
		printErr(err)
		return
	}
	fmt.Println("✓ Rescue initiated")
}

// ── Billing ──

func cmdGetBalance() {
	balance, err := client.GetBalance()
	if err != nil {
		printErr(err)
		return
	}
	fmt.Printf("  Available: %.2f %s\n", balance.Available, balance.Currency)
	fmt.Printf("  Locked:    %.2f %s\n", balance.Locked, balance.Currency)
}

func cmdListInvoicesInteractive(r *bufio.Reader) {
	page, _ := strconv.Atoi(prompt(r, "Page [1]: "))
	if page < 1 {
		page = 1
	}
	result, err := client.ListInvoices(neburst.WithPage(page))
	if err != nil {
		printErr(err)
		return
	}
	fmt.Printf("Total: %d, Page: %d\n\n", result.Total, result.Page)
	for _, inv := range result.Items {
		fmt.Printf("  %-36s  $%-8.2f  %-12s  %s\n", inv.UUID, inv.Amount, inv.Status, inv.Category)
	}
}

func cmdGetInvoiceInteractive(r *bufio.Reader) {
	id := prompt(r, "Invoice UUID: ")
	inv, err := client.GetInvoice(id)
	if err != nil {
		printErr(err)
		return
	}
	printJSON(inv)
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
