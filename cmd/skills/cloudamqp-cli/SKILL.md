---
name: cloudamqp-cli
description: Manage CloudAMQP instances, VPCs, teams, and RabbitMQ/LavinMQ configuration from the command line. Use when the user needs to create, configure, monitor, upgrade, or troubleshoot CloudAMQP message broker instances.
allowed-tools: Bash(cloudamqp:*)
---

# CloudAMQP CLI

## Quick start

```bash
# list all instances
cloudamqp instance list

# get instance details (includes connection URL and API key)
cloudamqp instance get --id 1234

# create an instance and wait for it to be ready
cloudamqp instance create --name=my-instance --plan=bunny-1 --region=amazon-web-services::us-east-1 --wait

# restart RabbitMQ on an instance
cloudamqp instance restart-rabbitmq --id 1234

# delete an instance
cloudamqp instance delete --id 1234
```

## Authentication

Before running any commands, check that auth is configured — the CLI won't work without it and can't prompt interactively when run by an agent.

```bash
# check if already configured
cat ~/.cloudamqprc 2>/dev/null || echo "not configured"

# if not set up, ask the user for their API key, then write it:
echo "YOUR_API_KEY" > ~/.cloudamqprc
chmod 600 ~/.cloudamqprc
```

The CLI checks in this order:

1. `CLOUDAMQP_APIKEY` environment variable
2. `~/.cloudamqprc` file (plain text, just the key)
3. Interactive prompt (won't work in agent context — use one of the above)

Base URL defaults to `https://customer.cloudamqp.com/api` (override with `CLOUDAMQP_URL`).

## Output

All commands support `-o json` for machine-readable output and `-o table` (default) for human-readable output. Use `-fields` to select specific columns.

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

### Copy settings between instances

```bash
cloudamqp instance create --name=staging --plan=bunny-1 --region=amazon-web-services::us-east-1 \
  --copy-from-id=1234 --copy-settings=metrics,firewall,config,alarms,logs,definitions,plugins --wait
```

Only works between dedicated instances (not shared plans).

### Node management

```bash
cloudamqp instance nodes list --id <id>
cloudamqp instance nodes versions --id <id>
```

### Plugin management

```bash
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
# restart
cloudamqp instance restart-rabbitmq --id <id> [--nodes=node1,node2]
cloudamqp instance restart-cluster --id <id>
cloudamqp instance restart-management --id <id> [--nodes=node1,node2]

# start/stop
cloudamqp instance start --id <id> [--nodes=node1,node2]
cloudamqp instance stop --id <id> [--nodes=node1,node2]
cloudamqp instance reboot --id <id> [--nodes=node1,node2]
cloudamqp instance start-cluster --id <id>
cloudamqp instance stop-cluster --id <id>

# upgrades (async, return immediately)
cloudamqp instance upgrade-erlang --id <id>
cloudamqp instance upgrade-rabbitmq --id <id> --version=<version>
cloudamqp instance upgrade-all --id <id>
cloudamqp instance upgrade-versions --id <id>
```

### VPC management

```bash
cloudamqp vpc create --name=<name> --region=<region> --subnet=<cidr> [--tags=<tag>]
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

### Audit log

```bash
cloudamqp audit [--timestamp=2024-01]
```

### API key rotation

```bash
cloudamqp rotate-key
```

## Important behavior

- **Async operations**: Instance creation, disk resizes, and upgrades are async. Use `--wait` on create, or poll with `instance get --id <id>` until `ready: true`.
- **Destructive commands** (delete, stop) prompt for confirmation. Use `--force` to skip in scripts.
- **Multiple tags**: Use `--tags` multiple times: `--tags=prod --tags=web`.
- **Shell completion**: Run `source <(cloudamqp completion zsh)` for tab completion of commands, instance IDs, plans, and regions.

## Plans and regions

Always fetch live data — don't guess plan names or regions:

```bash
cloudamqp plans [--backend=rabbitmq|lavinmq]
cloudamqp regions [--provider=amazon-web-services]
```

## Specific tasks

* **Scripting and automation** [references/scripting.md](references/scripting.md)
* **Instance upgrades and maintenance** [references/upgrades.md](references/upgrades.md)
* **VPC and network setup** [references/vpc-setup.md](references/vpc-setup.md)
