# Neburst API Java SDK

Official Java SDK for the [Neburst OpenAPI](https://neburst.com).

## Requirements

- Java 11+
- No external HTTP client required (uses `java.net.http.HttpClient`)

## Installation

### Maven

```xml
<dependency>
    <groupId>com.neburst</groupId>
    <artifactId>neburst-api</artifactId>
    <version>0.1.0</version>
</dependency>
```

## Quick Start

```java
import com.neburst.api.*;
import com.neburst.api.model.*;
import java.util.List;

public class Example {
    public static void main(String[] args) throws Exception {
        // Create client
        NeburstClient client = new NeburstClient(
            "https://api.neburst.com",
            "your-key-id",
            "your-api-secret"
        );

        // Compute operations
        ComputeApi compute = new ComputeApi(client);

        // List all instances
        List<Instance> instances = compute.listInstances();
        for (Instance inst : instances) {
            System.out.println(inst.getName() + " - " + inst.getStatus());
        }

        // Get instance details
        Instance instance = compute.getInstance("instance-uuid");

        // Get power status
        PowerStatus status = compute.getInstanceStatus("instance-uuid");
        System.out.println("Power: " + status.getStatus());

        // Power actions: "power-on", "power-off", "power-cycle", "power-reset"
        compute.powerAction("instance-uuid", "power-on");

        // Rebuild instance
        compute.rebuildInstance("instance-uuid", 1, "my-host", null);

        // Rescue mode
        compute.rescueInstance("instance-uuid", 1);

        // Get traffic
        Traffic traffic = compute.getInstanceTraffic("instance-uuid");

        // Billing operations
        BillingApi billing = new BillingApi(client);

        Balance balance = billing.getBalance();
        System.out.printf("Balance: %.2f %s%n", balance.getAvailable(), balance.getCurrency());

        List<Invoice> invoices = billing.listInvoices();
        Invoice invoice = billing.getInvoice("invoice-uuid");
    }
}
```

## Error Handling

```java
try {
    Instance inst = compute.getInstance("invalid-id");
} catch (NeburstApiException e) {
    System.err.println("API error code: " + e.getCode());
    System.err.println("Message: " + e.getMessage());
}
```

## Authentication

All requests are signed with HMAC-SHA256. The signing process:

1. Construct the string to sign:
   ```
   StringToSign = timestamp + "\n" + method + "\n" + path + "\n" + sorted_query + "\n" + SHA256(body)
   ```
2. Compute `HMAC-SHA256(api_secret, StringToSign)`
3. Send the hex-encoded signature in the `X-Nb-Signature` header

Required headers on every request:
- `X-Nb-Key` - Your API key ID
- `X-Nb-Timestamp` - Current Unix timestamp (seconds)
- `X-Nb-Nonce` - Unique UUID v4 per request
- `X-Nb-Signature` - HMAC-SHA256 signature (hex)

## License

MIT
