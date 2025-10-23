# CloudAMQP CLI Agent Reference

This document provides a comprehensive reference for AI agents using the CloudAMQP CLI tool.

## Configuration

### API Key Setup
The CLI uses API keys in this priority order:
1. `CLOUDAMQP_APIKEY` environment variable (for main API)
2. `~/.cloudamqprc` JSON config file
3. Interactive prompt

For instance-specific operations, use:
- `CLOUDAMQP_INSTANCE_{ID}_APIKEY` environment variable
- Stored in `~/.cloudamqprc` under `instance_keys`

### Base URL
Default: `https://customer.cloudamqp.com/api` (can be changed for development)

## Command Structure

```
cloudamqp <category> <action> [arguments] [flags]
```

## Main API Commands (using main API key)

### Instance Management

#### List Instances
```bash
cloudamqp instance list
```
- Returns: Array of instances with id, name, plan, region, ready status

#### Get Instance Details
```bash
cloudamqp instance get <id>
```
- Automatically saves instance API key for later use
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
cloudamqp instance update <id> --name=<new_name> --plan=<new_plan>
```
- Updates instance name and/or plan
- Use for upgrading/downgrading plans

#### Delete Instance
```bash
cloudamqp instance delete <id>
```
- Permanently deletes the instance

#### Resize Instance Disk
```bash
cloudamqp instance resize <id> --disk-size=<gb> [--allow-downtime]
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
cloudamqp vpc get <id>
```

#### Create VPC
```bash
cloudamqp vpc create --name=<name> --region=<region> --subnet=<subnet> [--tags=<tag>]
```

#### Update VPC
```bash
cloudamqp vpc update <id> --name=<new_name>
```

#### Delete VPC
```bash
cloudamqp vpc delete <id>
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
cloudamqp team update <id> --role=<role>
```

#### Remove Team Member
```bash
cloudamqp team remove <id>
```

### Billing & Plans

#### List Available Plans
```bash
cloudamqp plans
```
- Returns: Array of plans with name, price, backend, shared status

#### List Available Regions
```bash
cloudamqp regions
```


### Audit & Security

#### List Audit Logs
```bash
cloudamqp audit
```

#### Rotate API Key
```bash
cloudamqp rotate-key
```

## Instance-Specific API Commands (using instance API key)

These commands use the pattern: `cloudamqp instance manage <instance_id> <category> <action>`

### Node Management

#### List Nodes
```bash
cloudamqp instance manage <id> nodes
```

### Plugin Management

#### List Plugins
```bash
cloudamqp instance manage <id> plugins list
```
- Returns: Array of plugins with name, version, description, enabled status

### Account Operations

#### Rotate Instance Password
```bash
cloudamqp instance manage <id> account rotate-password
```

#### Rotate Instance API Key
```bash
cloudamqp instance manage <id> account rotate-apikey
```

### Instance Actions

#### Restart Operations
```bash
cloudamqp instance manage <id> actions restart-rabbitmq [--nodes=node1,node2]
cloudamqp instance manage <id> actions restart-cluster
cloudamqp instance manage <id> actions restart-management [--nodes=node1,node2]
```

#### Start/Stop Operations
```bash
cloudamqp instance manage <id> actions start [--nodes=node1,node2]
cloudamqp instance manage <id> actions stop [--nodes=node1,node2]
cloudamqp instance manage <id> actions reboot [--nodes=node1,node2]
cloudamqp instance manage <id> actions start-cluster
cloudamqp instance manage <id> actions stop-cluster
```

#### Upgrade Operations
```bash
cloudamqp instance manage <id> actions upgrade-erlang
cloudamqp instance manage <id> actions upgrade-rabbitmq --version=<version>
cloudamqp instance manage <id> actions upgrade-all
cloudamqp instance manage <id> actions upgrade-versions  # Check available versions
```

#### Feature Toggle Operations
```bash
cloudamqp instance manage <id> actions toggle-hipe --enable=true/false [--nodes=node1,node2]
cloudamqp instance manage <id> actions toggle-firehose --enable=true/false --vhost=<vhost>
```

## Common Usage Patterns

### 1. Create and Wait for Instance
```bash
# Create instance
cloudamqp instance create --name="my-instance" --plan="bunny-1" --region="amazon-web-services::us-east-1"

# Poll until ready
while true; do
  cloudamqp instance get <id> | grep '"ready": true' && break
  sleep 30
done
```

### 2. Upgrade Instance Plan
```bash
cloudamqp instance update <id> --plan="rabbit-3"
```

### 3. Instance Management Workflow
```bash
# Get instance details (saves API key automatically)
cloudamqp instance get <id>

# Now you can use instance-specific commands
cloudamqp instance manage <id> nodes
cloudamqp instance manage <id> plugins list
cloudamqp instance manage <id> actions restart-rabbitmq
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

1. Always use `cloudamqp instance get <id>` before instance-specific operations to ensure the instance API key is saved
2. Instance creation is async - poll with `get` until `ready: true`
3. Plan upgrades are immediate but may cause brief downtime
4. Some actions (upgrades, restarts) are asynchronous - they return immediately but run in background
5. The `--tags` flag can be used multiple times: `--tags=prod --tags=web`
6. VPC operations require the instance to be in the same region as the VPC
7. When asked for the openapi spec there are make targets `make openapi.yaml` and `make openapi-instance.yaml` to download the latest version.
