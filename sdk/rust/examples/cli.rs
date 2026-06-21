use neburst_api::NeburstClient;
use std::io::{self, Write};

fn prompt(msg: &str) -> String {
    print!("{}", msg);
    io::stdout().flush().unwrap();
    let mut buf = String::new();
    io::stdin().read_line(&mut buf).unwrap();
    buf.trim().to_string()
}

#[tokio::main]
async fn main() {
    println!("╔══════════════════════════════════════╗");
    println!("║     Neburst OpenAPI CLI (Rust)       ║");
    println!("╚══════════════════════════════════════╝");
    println!();

    let base_url = {
        let v = prompt("API Base URL [https://api.neburst.com]: ");
        if v.is_empty() { "https://api.neburst.com".to_string() } else { v }
    };

    let api_key = prompt("API Key (base64 combined key or Key ID): ");
    let client = if api_key.starts_with("nb_key_") {
        let secret = prompt("API Secret: ");
        NeburstClient::new(&base_url, &api_key, &secret)
    } else {
        NeburstClient::from_combined_key(&base_url, &api_key).expect("Invalid combined key")
    };

    println!("\n✓ Client initialized\n");

    loop {
        print_menu();
        let choice = prompt("\n> ");
        println!();

        let result: Result<(), Box<dyn std::error::Error>> = async {
            match choice.as_str() {
                "1" => {
                    let page: i32 = prompt("Page [1]: ").parse().unwrap_or(1);
                    let r = client.list_instances(page, 20).await?;
                    println!("Total: {}, Page: {}", r.total, r.page);
                    for i in &r.items { println!("  {}  {:<10}  {:<10}  {}", i.uuid, i.instance_type, i.status, i.name); }
                }
                "2" => {
                    let id = prompt("Instance UUID: ");
                    println!("{}", serde_json::to_string_pretty(&client.get_instance(&id).await?)?);
                }
                "3" => {
                    let id = prompt("Instance UUID: ");
                    println!("{}", serde_json::to_string_pretty(&client.get_instance_status(&id).await?)?);
                }
                "4" => {
                    let id = prompt("Instance UUID: ");
                    println!("{}", serde_json::to_string_pretty(&client.get_instance_traffic(&id).await?)?);
                }
                "5" => {
                    let id = prompt("Instance UUID: ");
                    let action = prompt("Action (start/stop/restart): ");
                    client.cloud_power_action(&id, &action).await?;
                    println!("✓ Success");
                }
                "6" => {
                    let id = prompt("Instance UUID: ");
                    println!("{}", serde_json::to_string_pretty(&client.get_cloud_metrics(&id).await?)?);
                }
                "11" => {
                    let page: i32 = prompt("Page [1]: ").parse().unwrap_or(1);
                    let r = client.list_bare_metal_instances(page, 20).await?;
                    println!("Total: {}, Page: {}", r.total, r.page);
                    for i in &r.items { println!("  {}  {:<10}  {}", i.uuid, i.status, i.name); }
                }
                "12" => {
                    let id = prompt("Instance UUID: ");
                    println!("{}", serde_json::to_string_pretty(&client.get_bare_metal_instance(&id).await?)?);
                }
                "13" => {
                    let id = prompt("Instance UUID: ");
                    println!("{}", serde_json::to_string_pretty(&client.get_bare_metal_status(&id).await?)?);
                }
                "14" => {
                    let id = prompt("Instance UUID: ");
                    println!("{}", serde_json::to_string_pretty(&client.get_bare_metal_traffic(&id).await?)?);
                }
                "15" => {
                    let id = prompt("Instance UUID: ");
                    let action = prompt("Action (power-on/power-off/power-cycle/power-reset): ");
                    client.bare_metal_power_action(&id, &action).await?;
                    println!("✓ Success");
                }
                "16" => {
                    let id = prompt("Instance UUID: ");
                    println!("{}", serde_json::to_string_pretty(&client.get_bare_metal_metrics(&id).await?)?);
                }
                "17" => {
                    let id = prompt("Instance UUID: ");
                    for p in &client.get_reinstall_profiles(&id).await? { println!("  [{}] {:<30}  {}", p.id, p.name, p.category); }
                }
                "18" => {
                    let id = prompt("Instance UUID: ");
                    for p in &client.get_rescue_profiles(&id).await? { println!("  [{}] {:<30}  {}", p.id, p.name, p.category); }
                }
                "31" => {
                    let b = client.get_balance().await?;
                    println!("  Available: {:.2} {}", b.available, b.currency);
                    println!("  Locked:    {:.2} {}", b.locked, b.currency);
                }
                "32" => {
                    let page: i32 = prompt("Page [1]: ").parse().unwrap_or(1);
                    let r = client.list_invoices(page, 20).await?;
                    println!("Total: {}, Page: {}", r.total, r.page);
                    for inv in &r.items { println!("  {}  ${:<8.2}  {:<12}  {}", inv.uuid, inv.amount, inv.status, inv.category); }
                }
                "33" => {
                    let id = prompt("Invoice UUID: ");
                    println!("{}", serde_json::to_string_pretty(&client.get_invoice(&id).await?)?);
                }
                "0" | "q" | "quit" | "exit" => {
                    println!("Bye!");
                    std::process::exit(0);
                }
                _ => println!("Unknown command"),
            }
            Ok(())
        }.await;

        if let Err(e) = result {
            println!("✗ Error: {}", e);
        }
        println!();
    }
}

fn print_menu() {
    println!("─── Cloud Instance ───────────────────");
    println!("  1  List instances");
    println!("  2  Get instance");
    println!("  3  Get status");
    println!("  4  Get traffic");
    println!("  5  Power action (start/stop/restart)");
    println!("  6  Metrics");
    println!("─── Bare Metal ───────────────────────");
    println!("  11 List instances");
    println!("  12 Get instance");
    println!("  13 Get status");
    println!("  14 Get traffic");
    println!("  15 Power action");
    println!("  16 Metrics");
    println!("  17 OS profiles");
    println!("  18 Rescue profiles");
    println!("─── Billing ──────────────────────────");
    println!("  31 Balance");
    println!("  32 List invoices");
    println!("  33 Get invoice");
    println!("──────────────────────────────────────");
    println!("  0  Exit");
}
