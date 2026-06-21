#!/usr/bin/env python3
"""Neburst OpenAPI interactive CLI (Python)."""

import json
import sys

from neburst import NeburstClient


def main():
    print("╔══════════════════════════════════════╗")
    print("║     Neburst OpenAPI CLI (Python)     ║")
    print("╚══════════════════════════════════════╝")
    print()

    base_url = input("API Base URL [https://api.neburst.com]: ").strip() or "https://api.neburst.com"
    api_key = input("API Key (base64 combined key or Key ID): ").strip()

    if api_key.startswith("nb_key_"):
        secret = input("API Secret: ").strip()
        client = NeburstClient(base_url, api_key, secret)
    else:
        client = NeburstClient(base_url, api_key)

    print("\n✓ Client initialized\n")

    while True:
        print_menu()
        choice = input("\n> ").strip()
        print()

        try:
            if choice == "1":
                page = int(input("Page [1]: ").strip() or "1")
                result = client.list_instances(page=page)
                print(f"Total: {result.total}, Page: {result.page}")
                for inst in result.items:
                    print(f"  {inst.uuid}  {inst.type:<10}  {inst.status:<10}  {inst.name}")

            elif choice == "2":
                uid = input("Instance UUID: ").strip()
                print(pjson(client.get_instance(uid)))

            elif choice == "3":
                uid = input("Instance UUID: ").strip()
                print(pjson(client.get_instance_status(uid)))

            elif choice == "4":
                uid = input("Instance UUID: ").strip()
                print(pjson(client.get_instance_traffic(uid)))

            elif choice == "5":
                uid = input("Instance UUID: ").strip()
                action = input("Action (start/stop/restart): ").strip()
                client.cloud_power_action(uid, action)
                print("✓ Success")

            elif choice == "6":
                uid = input("Instance UUID: ").strip()
                print(pjson(client.get_cloud_metrics(uid)))

            elif choice == "11":
                page = int(input("Page [1]: ").strip() or "1")
                result = client.list_bare_metal_instances(page=page)
                print(f"Total: {result.total}, Page: {result.page}")
                for inst in result.items:
                    print(f"  {inst.uuid}  {inst.status:<10}  {inst.name}")

            elif choice == "12":
                uid = input("Instance UUID: ").strip()
                print(pjson(client.get_bare_metal_instance(uid)))

            elif choice == "13":
                uid = input("Instance UUID: ").strip()
                print(pjson(client.get_bare_metal_status(uid)))

            elif choice == "14":
                uid = input("Instance UUID: ").strip()
                print(pjson(client.get_bare_metal_traffic(uid)))

            elif choice == "15":
                uid = input("Instance UUID: ").strip()
                action = input("Action (power-on/power-off/power-cycle/power-reset): ").strip()
                client.bare_metal_power_action(uid, action)
                print("✓ Success")

            elif choice == "16":
                uid = input("Instance UUID: ").strip()
                print(pjson(client.get_bare_metal_metrics(uid)))

            elif choice == "17":
                uid = input("Instance UUID: ").strip()
                for p in client.get_reinstall_profiles(uid):
                    print(f"  [{p.id}] {p.name:<30}  {p.category}")

            elif choice == "18":
                uid = input("Instance UUID: ").strip()
                for p in client.get_rescue_profiles(uid):
                    print(f"  [{p.id}] {p.name:<30}  {p.category}")

            elif choice == "19":
                uid = input("Instance UUID: ").strip()
                pid = int(input("Profile ID: ").strip())
                hostname = input("Hostname (optional): ").strip() or None
                client.rebuild_instance(uid, pid, hostname=hostname)
                print("✓ Rebuild initiated")

            elif choice == "20":
                uid = input("Instance UUID: ").strip()
                pid = int(input("Profile ID: ").strip())
                client.rescue_instance(uid, pid)
                print("✓ Rescue initiated")

            elif choice == "31":
                b = client.get_balance()
                print(f"  Available: {b.available:.2f} {b.currency}")
                print(f"  Locked:    {b.locked:.2f} {b.currency}")

            elif choice == "32":
                page = int(input("Page [1]: ").strip() or "1")
                result = client.list_invoices(page=page)
                print(f"Total: {result.total}, Page: {result.page}")
                for inv in result.items:
                    print(f"  {inv.uuid}  ${inv.amount:<8.2f}  {inv.status:<12}  {inv.category}")

            elif choice == "33":
                uid = input("Invoice UUID: ").strip()
                print(pjson(client.get_invoice(uid)))

            elif choice == "41":
                print(pjson(client.get_user_info()))

            elif choice in ("0", "q", "quit", "exit"):
                print("Bye!")
                sys.exit(0)

            else:
                print("Unknown command")

        except Exception as e:
            print(f"✗ Error: {e}")

        print()


def print_menu():
    print("─── Cloud Instance ───────────────────")
    print("  1  List instances")
    print("  2  Get instance")
    print("  3  Get status")
    print("  4  Get traffic")
    print("  5  Power action (start/stop/restart)")
    print("  6  Metrics")
    print("─── Bare Metal ───────────────────────")
    print("  11 List instances")
    print("  12 Get instance")
    print("  13 Get status")
    print("  14 Get traffic")
    print("  15 Power action")
    print("  16 Metrics")
    print("  17 OS profiles")
    print("  18 Rescue profiles")
    print("  19 Rebuild")
    print("  20 Rescue")
    print("─── Billing ──────────────────────────")
    print("  31 Balance")
    print("  32 List invoices")
    print("  33 Get invoice")
    print("─── User ─────────────────────────────")
    print("  41 User info")
    print("──────────────────────────────────────")
    print("  0  Exit")


def pjson(obj):
    if hasattr(obj, "__dict__"):
        return json.dumps(obj.__dict__, indent=2, default=str)
    return json.dumps(obj, indent=2, default=str)


if __name__ == "__main__":
    main()
