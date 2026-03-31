---
name: cloudamqp-cli
description: Manage CloudAMQP instances, VPCs, teams, and RabbitMQ/LavinMQ configuration using the cloudamqp CLI. Use this skill whenever the user wants to create, list, inspect, update, delete, upgrade, restart, or troubleshoot CloudAMQP instances — even if they just say "spin up a RabbitMQ", "check my instances", or "upgrade my broker". Also use it for VPC setup, team management, and RabbitMQ config changes.
allowed-tools: Bash(cloudamqp:*), Bash(jq:*), Bash(cat:*), Bash(echo:*), Bash(chmod:*), Bash(grep:*), Bash(sleep:*)
---

# CloudAMQP CLI

## Quick start

```bash
cloudamqp instance list
cloudamqp instance get --id <id>
cloudamqp instance create --name=<name> --plan=<plan> --region=<region> --wait
cloudamqp instance restart-rabbitmq --id <id>
cloudamqp instance delete --id <id> --force
```

## Authentication

Check auth before running anything — interactive prompts don't work in agent context:

```bash
cat ~/.cloudamqprc 2>/dev/null || echo "not configured"
```

If not configured, ask the user for their API key (from https://customer.cloudamqp.com/apikeys), then:

```bash
echo "<api-key>" > ~/.cloudamqprc
chmod 600 ~/.cloudamqprc
```

Alternatively, set `CLOUDAMQP_APIKEY` in the environment. If neither is set, all commands will fail.

## Output

Use `-o json` for parsing, `-o table` (default) for display. Use `--fields` to select columns.

## Commands

### Instance lifecycle

```bash
cloudamqp instance create --name=<name> --plan=<plan> --region=<region> [--tags=<tag>...] [--vpc-id=<id>] [--wait] [--wait-timeout=20m]
cloudamqp instance list [--details]
cloudamqp instance get --id <id>
cloudamqp instance update --id <id> [--name=<name>] [--plan=<plan>]
cloudamqp instance delete --id <id> [--force]
cloudamqp instance resize-disk --id <id> --disk-size=<gb> [--allow-downtime]
```

### Copy settings between instances (dedicated only)

```bash
cloudamqp instance create --name=staging --plan=<plan> --region=<region> \
  --copy-from-id=<id> --copy-settings=metrics,firewall,config,alarms,logs,definitions,plugins --wait
```

### Node and plugin management

```bash
cloudamqp instance nodes list --id <id>
cloudamqp instance nodes versions --id <id>
cloudamqp instance plugins list --id <id>
```

### RabbitMQ configuration

```bash
cloudamqp instance config list --id <id>
cloudamqp instance config get --id <id> --key <key>
cloudamqp instance config set --id <id> --key <key> --value <value>
```

### Instance actions

```bash
# restarts (rolling for HA clusters)
cloudamqp instance restart-rabbitmq --id <id> [--nodes=node1,node2]
cloudamqp instance restart-cluster --id <id>          # full restart, causes downtime
cloudamqp instance restart-management --id <id>

# start/stop
cloudamqp instance start --id <id>
cloudamqp instance stop --id <id>
cloudamqp instance reboot --id <id>
cloudamqp instance start-cluster --id <id>
cloudamqp instance stop-cluster --id <id>

# upgrades — async, return immediately, poll until ready
cloudamqp instance upgrade-erlang --id <id>
cloudamqp instance upgrade-rabbitmq --id <id> --version=<version>
cloudamqp instance upgrade-all --id <id>
cloudamqp instance upgrade-versions --id <id>
```

### VPC management

```bash
cloudamqp vpc create --name=<name> --region=<region> --subnet=<cidr>
cloudamqp vpc list
cloudamqp vpc get --id <id>
cloudamqp vpc update --id <id> --name=<name>
cloudamqp vpc delete --id <id>
```

### Team management

```bash
cloudamqp team list
cloudamqp team invite --email=<email> [--role=<role>] [--tags=<tag>]
cloudamqp team update --user-id=<id> --role=<role>
cloudamqp team remove --email=<email>
```

### Plans, regions, audit

```bash
cloudamqp plans [--backend=rabbitmq|lavinmq]   # always fetch, never guess
cloudamqp regions [--provider=<provider>]       # always fetch, never guess
cloudamqp audit [--timestamp=2024-01]
cloudamqp rotate-key
```

## Key behaviors

- **Async**: creation, resizes, upgrades return immediately. Use `--wait` on create, or poll `instance get --id <id> -o json | jq -r '.ready'` until `"true"`.
- **Destructive commands** prompt for confirmation — use `--force` to skip.
- **Multiple tags**: repeat the flag: `--tags=prod --tags=web`.
- **Plan/region names**: always run `cloudamqp plans` / `cloudamqp regions` first — never hardcode them.

## Reference guides

Read these before tackling the relevant task:

- **Scripting, JSON parsing, batch ops** → [references/scripting.md](references/scripting.md)
- **Upgrades, restarts, maintenance workflows** → [references/upgrades.md](references/upgrades.md)
- **VPC creation and network setup** → [references/vpc-setup.md](references/vpc-setup.md)
