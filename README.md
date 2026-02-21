# Gobashops - Senior Admin CLI

This CLI tool implements the Senior Linux Admin mental model. It replaces disjointed bash scripts and mental checklists with an automated triage system written in Go.

## The Mental Model (Core Philosophy)
Always troubleshoot in this order:
1. Scope 
2. Reproduce
3. Logs
4. Resource pressure
5. Network
6. Permissions / SELinux
7. Configuration

## Installation
```bash
go build -o gobashops .
sudo mv gobashops /usr/local/bin/
```

## Usage

### 1. General System Slowness (`gobashops triage`)
Instead of running `uptime`, `free`, `df`, and `top` manually, this command bundles them and interprets the output:
```bash
gobashops triage
```
*Checks Load averages, Memory, Swap, Disk Pressure (and orphaned files), and Top CPU processes.*

### 2. Service Down (`gobashops service <name>`)
Implements the 5-layer service check:
```bash
gobashops service nginx
```
*Checks `systemctl status`, `journalctl`, `ss -tulnp` (listening ports), and `ausearch` (SELinux AVC denials).*

### 3. Server Unreachable (`gobashops net <ip> --port <port>`)
Implements the 4-layer network check natively:
```bash
gobashops net 10.0.0.5 --port 443
```
*Checks ICMP ping, `ip route` tables, and performs a native Go TCP handshake (replacing `nc -zv`).*

### 4. Storage & LVM (`gobashops storage`)
Quick overview of LVM and block devices:
```bash
gobashops storage
```
*Runs `pvs`, `vgs`, and `lsblk -f` to check filesystem mounts and LVM capacity.*
