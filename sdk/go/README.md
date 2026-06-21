# Neburst Go SDK

Official Go SDK for the [Neburst](https://neburst.com) OpenAPI.

## Install

```bash
go get github.com/neburstnetworks/openapi/sdk/go/neburst@latest
```

## Quick Start

```go
package main

import (
	"fmt"
	"log"

	"github.com/neburstnetworks/openapi/sdk/go/neburst"
)

func main() {
	client := neburst.NewClient(
		"https://api.neburst.com",
		"your-key-id",
		"your-api-secret",
	)

	// List instances with pagination
	result, err := client.ListInstances(neburst.WithPage(1), neburst.WithPageSize(10))
	if err != nil {
		log.Fatal(err)
	}
	for _, inst := range result.Items {
		fmt.Printf("%s  %s  %s  %s\n", inst.UUID, inst.Name, inst.PrimaryIPv4, inst.Status)
	}

	// Power off a bare-metal instance
	err = client.BareMetalPowerAction(result.Items[0].UUID, neburst.PowerOff)
	if err != nil {
		log.Fatal(err)
	}

	// Check balance
	balance, err := client.GetBalance()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Balance: %.2f %s\n", balance.Available, balance.Currency)
}
```

## Authentication

The SDK handles HMAC-SHA256 request signing automatically. You need an API key pair (key ID + secret) which can be created in the Neburst dashboard.

Each request is signed with:
- `X-Nb-Key` - your key ID
- `X-Nb-Timestamp` - unix timestamp (seconds)
- `X-Nb-Nonce` - unique UUID v4 per request
- `X-Nb-Signature` - HMAC-SHA256 signature

## API Reference

### Compute (Common)

| Method | Description |
|--------|-------------|
| `ListInstances(opts...)` | List instances (paginated). Options: `WithPage()`, `WithPageSize()` |
| `GetInstance(id)` | Get instance details |
| `GetInstanceStatus(id)` | Get power status |
| `GetInstanceTraffic(id)` | Get traffic usage |

### Compute (Bare Metal)

| Method | Description |
|--------|-------------|
| `BareMetalPowerAction(id, action)` | Power action (PowerOn/PowerOff/PowerCycle/PowerReset) |
| `RebuildInstance(id, profileID, opts...)` | Rebuild OS. Options: `WithHostname()`, `WithPublicKeys()` |
| `RescueInstance(id, profileID)` | Boot into rescue mode |
| `GetReinstallProfiles(id)` | List available OS profiles |
| `GetRescueProfiles(id)` | List rescue profiles |
| `GetBareMetalMetrics(id)` | Get metrics (network, bandwidth) |

### Compute (Cloud)

| Method | Description |
|--------|-------------|
| `CloudPowerAction(id, action)` | Power action (start/stop/restart) |
| `GetCloudMetrics(id)` | Get metrics (CPU, memory, disk, network, bandwidth) |

### Billing

| Method | Description |
|--------|-------------|
| `GetBalance()` | Get account balance |
| `ListInvoices(opts...)` | List invoices (paginated) |
| `GetInvoice(id)` | Get invoice details |

## Error Handling

API errors are returned as `*neburst.APIError`:

```go
instance, err := client.GetInstance("invalid-uuid")
if err != nil {
	var apiErr *neburst.APIError
	if errors.As(err, &apiErr) {
		fmt.Printf("API error %d: %s\n", apiErr.Code, apiErr.Message)
	}
}
```

## Custom HTTP Client

```go
client := neburst.NewClient(baseURL, keyID, secret)
client.SetHTTPClient(&http.Client{
	Timeout: 60 * time.Second,
})
```
