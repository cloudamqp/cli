# CloudAMQP CLI Agent Reference

This document provides a comprehensive reference for AI agents using the CloudAMQP CLI tool.

## Configuration

### API Key Setup
The CLI uses a single API key for all operations in this priority order:
1. `CLOUDAMQP_APIKEY` environment variable
2. `~/.cloudamqprc` plain text config file
3. Interactive prompt

### Base URL
Default: `https://customer.cloudamqp.com/api` (unified API endpoint)

## Command Structure

```
cloudamqp <category> <action> [--id <instance_id>] [other flags]
```

All instance-specific operations use the `--id` flag to specify the instance.

## Main API Commands

### Instance Management

#### List Instances
```bash
cloudamqp instance list
```
- Returns: Array of instances with id, name, plan, region, ready status

#### Get Instance Details
```bash
cloudamqp instance get --id <id>
```
- Returns: Full instance details including API key, URLs, hostnames

#### Create Instance
```bash
cloudamqp instance create --name=<name> --plan=<plan> --region=<region> [--tags=<tag1> --tags=<tag2>] [--vpc-subnet=<subnet>] [--vpc-id=<id>]
```
- Required: name, plan, region
- Optional: tags (multiple allowed), vpc-subnet, vpc-id
- Returns: Instance creation response with id, url, apikey

#### Update Instance
```bash
cloudamqp instance update --id <id> --name=<new_name> --plan=<new_plan>
```
- Updates instance name and/or plan
- Use for upgrading/downgrading plans

#### Delete Instance
```bash
cloudamqp instance delete --id <id>
```
- Permanently deletes the instance

#### Resize Instance Disk
```bash
cloudamqp instance resize-disk --id <id> --disk-size=<gb> [--allow-downtime]
```
- Required: disk-size (in GB)
- Optional: allow-downtime flag

### VPC Management

#### List VPCs
```bash
cloudamqp vpc list
```

#### Get VPC Details
```bash
cloudamqp vpc get --id <id>
```

#### Create VPC
```bash
cloudamqp vpc create --name=<name> --region=<region> --subnet=<subnet> [--tags=<tag>]
```

#### Update VPC
```bash
cloudamqp vpc update --id <id> --name=<new_name>
```

#### Delete VPC
```bash
cloudamqp vpc delete --id <id>
```

### Team Management

#### List Team Members
```bash
cloudamqp team list
```

#### Invite Team Member
```bash
cloudamqp team invite --email=<email> [--role=<role>] [--tags=<tag>]
```

#### Update Team Member
```bash
cloudamqp team update --user-id <id> --role=<role>
```

#### Remove Team Member
```bash
cloudamqp team remove --email=<email>
```

### Billing & Plans

#### List Available Plans
```bash
cloudamqp plans [--backend=<rabbitmq|lavinmq>]
```
- Returns: Array of plans with name, price, backend, shared status

#### List Available Regions
```bash
cloudamqp regions [--provider=<provider>]
```

### Audit & Security

#### Export Audit Logs
```bash
cloudamqp audit [--timestamp=<timestamp>]
```


## Instance-Specific Operations

All instance-specific commands use the unified API and `--id` flag pattern.

### Node Management

#### List Nodes
```bash
cloudamqp instance nodes list --id <id>
```

#### Get Available Versions
```bash
cloudamqp instance nodes versions --id <id>
```

### Plugin Management

#### List Plugins
```bash
cloudamqp instance plugins list --id <id>
```
- Returns: Array of plugins with name, version, description, enabled status

### RabbitMQ Configuration

#### List All Configuration Settings
```bash
cloudamqp instance config list --id <id>
```

#### Get Specific Configuration Setting
```bash
cloudamqp instance config get --id <id> --key <config_key>
```

#### Set Configuration Setting
```bash
cloudamqp instance config set --id <id> --key <config_key> --value <config_value>
```

### Account Operations


### Instance Actions

#### Restart Operations
```bash
cloudamqp instance restart-rabbitmq --id <id> [--nodes=node1,node2]
cloudamqp instance restart-cluster --id <id>
cloudamqp instance restart-management --id <id> [--nodes=node1,node2]
```

#### Start/Stop Operations
```bash
cloudamqp instance start --id <id> [--nodes=node1,node2]
cloudamqp instance stop --id <id> [--nodes=node1,node2]
cloudamqp instance reboot --id <id> [--nodes=node1,node2]
cloudamqp instance start-cluster --id <id>
cloudamqp instance stop-cluster --id <id>
```

#### Upgrade Operations
```bash
cloudamqp instance upgrade-erlang --id <id>
cloudamqp instance upgrade-rabbitmq --id <id> --version=<version>
cloudamqp instance upgrade-all --id <id>
cloudamqp instance upgrade-versions --id <id>  # Check available versions
```


## Common Usage Patterns

### 1. Create and Wait for Instance
```bash
# Create instance
cloudamqp instance create --name="my-instance" --plan="bunny-1" --region="amazon-web-services::us-east-1"

# Poll until ready
while true; do
  cloudamqp instance get --id <id> | grep '"ready": true' && break
  sleep 30
done
```

### 2. Upgrade Instance Plan
```bash
cloudamqp instance update --id <id> --plan="rabbit-3"
```

### 3. Complete Instance Management Workflow
```bash
# Get instance details
cloudamqp instance get --id <id>

# Check nodes and configuration
cloudamqp instance nodes list --id <id>
cloudamqp instance config list --id <id>

# Perform maintenance
cloudamqp instance restart-rabbitmq --id <id>
cloudamqp instance upgrade-all --id <id>
```

### 4. Configuration Management
```bash
# List all configuration
cloudamqp instance config list --id <id>

# Get specific setting
cloudamqp instance config get --id <id> --key tcp_listen_options

# Update configuration
cloudamqp instance config set --id <id> --key tcp_listen_options --value '[{"port": 5672}]'
```

## Common Plans

### Free Tier
- `lemur` - Free RabbitMQ shared
- `lemming` - Free LavinMQ shared

### Paid Tiers (RabbitMQ)
- `bunny-1` - $99
- `bunny-3` - $297  
- `rabbit-1` - $299
- `rabbit-3` - $897
- `rabbit-5` - $1495

### Paid Tiers (LavinMQ)
- `penguin-1` - $99
- `penguin-3` - $297
- `penguin-5` - $495

## Common Regions
- `amazon-web-services::us-east-1`
- `amazon-web-services::us-west-1`
- `amazon-web-services::us-west-2`
- `amazon-web-services::eu-west-1`
- `google-compute-engine::us-central1-a`

## Error Handling

- API errors return non-zero exit codes
- Error messages are printed to stderr
- Most commands return JSON output on success
- Use environment variables for API keys to avoid exposing them in command history

## Notes for AI Agents

1. **Unified API**: All operations now use a single API key and the unified customer API endpoint
2. **Flag-based Commands**: All instance-specific operations use `--id <instance_id>` instead of positional arguments
3. **Instance Creation**: Instance creation is async - poll with `get --id <id>` until `ready: true`
4. **Plan Upgrades**: Plan upgrades are immediate but may cause brief downtime
5. **Async Operations**: Some actions (upgrades, restarts) are asynchronous - they return immediately but run in background
6. **Multiple Tags**: The `--tags` flag can be used multiple times: `--tags=prod --tags=web`
7. **VPC Requirements**: VPC operations require the instance to be in the same region as the VPC
8. **Configuration Management**: Use the config commands to manage RabbitMQ settings directly
9. **OpenAPI Specs**: When asked for OpenAPI specs, use make targets `make openapi.yaml` and `make openapi-instance.yaml` to download the latest versions

## Configuration File Format

The `~/.cloudamqprc` file contains only your API key in plain text:
```
your-api-key-here
```

No JSON formatting or multiple keys are needed - the unified API handles all operations with a single key.