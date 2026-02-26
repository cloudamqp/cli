# CloudAMQP CLI

A command line interface for the CloudAMQP API that provides complete management of CloudAMQP instances, VPCs, and instance-specific operations.

## Features

- **Unified API**: Single API key manages all operations through the customer API
- **Simple Configuration**: Plain text API key storage in `~/.cloudamqprc`
- **Positional IDs**: Clean command structure using positional IDs for instance and VPC operations
- **Copy Settings**: Clone configuration from existing instances (metrics, firewall, alarms, etc.)
- **Wait for Ready**: Optional `--wait` flag for long-running operations (create, resize-disk, upgrades)
- **User-Friendly**: Clear help messages, examples, and safety confirmations
- **Error Handling**: Proper API error extraction and display

## Installation

### Pre-built binaries

Pre-built binaries are available for each release here: https://github.com/cloudamqp/cli/releases

### Build from Source

```bash
go mod download
go build -o cloudamqp
```

### Usage

```bash
./cloudamqp --help
```

## Configuration

### API Key Configuration

The CLI looks for your API key in the following order:

1. `CLOUDAMQP_APIKEY` environment variable
2. `~/.cloudamqprc` file (plain text format)
3. If neither exists, you will be prompted to enter it

### Config File Format

The configuration file `~/.cloudamqprc` contains only your API key in plain text:

```
your-api-key-here
```

### Environment Variables

- `CLOUDAMQP_APIKEY` - Your CloudAMQP API key

### Shell Completion

The CLI supports shell completion for zsh, providing:
- Command and subcommand completion
- Flag completion
- Dynamic completion for instance IDs, VPC IDs, plan names, and regions (fetched from the API)

#### Zsh Completion Setup

**Option 1: Source in your shell session**

Add to your `~/.zshrc`:
```bash
source <(cloudamqp completion zsh)
```

**Option 2: Install to completion directory**

```bash
cloudamqp completion zsh > "${fpath[1]}/_cloudamqp"
```

After installation, restart your shell or reload your configuration:
```bash
exec zsh
```

#### Testing Completion

After setup, you can test completion by typing:
```bash
cloudamqp instance <TAB>          # Lists instance subcommands
cloudamqp instance get <TAB>      # Lists your instance IDs
cloudamqp instance create --plan <TAB>   # Lists available plans
cloudamqp instance create --region <TAB> # Lists available regions
```

Note: Dynamic completions (instance IDs, plans, regions) require a configured API key. Completion data is cached in `~/.cache/cloudamqp/` (clear with `rm -rf ~/.cache/cloudamqp/` if needed).

## Commands

#### Output
You can output either as JSON via `-o json` or Table format using `-o table`

### Instance Management

Manage CloudAMQP instances using your main API key.

```bash
# Create a new instance
cloudamqp instance create --name=my-instance --plan=bunny-1 --region=amazon-web-services::us-east-1

# Create instance with copy_settings (dedicated instances only)
cloudamqp instance create --name=my-copy --plan=bunny-1 --region=amazon-web-services::us-east-1 \
  --copy-from-id=1234 --copy-settings=metrics,firewall,config

# Create instance and wait for it to be ready (default timeout: 15m)
cloudamqp instance create --name=my-instance --plan=bunny-1 --region=amazon-web-services::us-east-1 --wait

# Create instance with custom wait timeout
cloudamqp instance create --name=my-instance --plan=bunny-1 --region=amazon-web-services::us-east-1 --wait --wait-timeout=20m

# List all instances
cloudamqp instance list

# List all instances with more details
cloudamqp instance list --details

# Get instance details
cloudamqp instance get 1234

# Update instance properties
cloudamqp instance update 1234 --name=new-name --plan=rabbit-1

# Resize instance disk
cloudamqp instance resize-disk 1234 --disk-size=100 --allow-downtime

# Delete instance (with confirmation)
cloudamqp instance delete 1234
```

### VPC Management

Manage Virtual Private Clouds.

```bash
# Create a VPC
cloudamqp vpc create --name=my-vpc --region=amazon-web-services::us-east-1 --subnet=10.56.72.0/24

# List all VPCs
cloudamqp vpc list

# Get VPC details
cloudamqp vpc get 5678

# Update VPC
cloudamqp vpc update 5678 --name=new-vpc-name

# Delete VPC (with confirmation)
cloudamqp vpc delete 5678
```

### Instance-Specific Management

Manage specific instances using the unified API. All commands take the instance ID as a positional argument.

#### Node Management

```bash
# List nodes in an instance
cloudamqp instance nodes list 1234

# Get available versions for upgrade
cloudamqp instance nodes versions 1234
```

#### Plugin Management

```bash
# List available RabbitMQ plugins
cloudamqp instance plugins list 1234
```

#### RabbitMQ Configuration

```bash
# List all configuration settings
cloudamqp instance config list 1234

# Get specific configuration setting
cloudamqp instance config get 1234 tcp_listen_options

# Set configuration setting
cloudamqp instance config set 1234 tcp_listen_options '[{"port": 5672}]'
```

#### Instance Actions

```bash
# Restart RabbitMQ
cloudamqp instance restart-rabbitmq 1234
cloudamqp instance restart-rabbitmq 1234 --nodes=node1,node2

# Cluster operations
cloudamqp instance restart-cluster 1234
cloudamqp instance stop-cluster 1234
cloudamqp instance start-cluster 1234

# Instance lifecycle
cloudamqp instance stop 1234
cloudamqp instance start 1234
cloudamqp instance reboot 1234

# Management interface
cloudamqp instance restart-management 1234

# Upgrades (asynchronous operations)
cloudamqp instance upgrade-erlang 1234
cloudamqp instance upgrade-rabbitmq 1234 --version=3.10.7
cloudamqp instance upgrade-all 1234

# Get target upgrade versions
cloudamqp instance upgrade-versions 1234

```

### Informational Commands

```bash
# List available regions
cloudamqp regions
cloudamqp regions --provider=amazon-web-services

# List available plans
cloudamqp plans
cloudamqp plans --backend=rabbitmq
```

### Team Management

```bash
# List team members
cloudamqp team list

# Invite new team member
cloudamqp team invite --email=user@example.com --role=admin --tags=production

# Update team member
cloudamqp team update --user-id=uuid-here --role=devops

# Remove team member
cloudamqp team remove --email=user@example.com
```

### Administrative Commands

```bash
# Export audit log
cloudamqp audit
cloudamqp audit --timestamp=2024-01
```

## Examples

### Complete Workflow

```bash
# 1. Create an instance and wait for it to be ready
cloudamqp instance create --name=production --plan=bunny-1 --region=amazon-web-services::us-east-1 --wait

# 2. Get instance details
cloudamqp instance get 1234

# 3. Check instance nodes
cloudamqp instance nodes list 1234

# 4. List RabbitMQ configuration
cloudamqp instance config list 1234

# 5. Install plugins (if needed)
cloudamqp instance plugins list 1234

# 6. Restart RabbitMQ
cloudamqp instance restart-rabbitmq 1234

# 7. Upgrade when needed
cloudamqp instance upgrade-all 1234
```

### Copy Settings Between Instances

Copy configuration from an existing dedicated instance to a new one:

```bash
# 1. Create original instance and wait for it to be ready
cloudamqp instance create --name=production --plan=bunny-1 --region=amazon-web-services::us-east-1 --wait

# 2. Configure the original instance (alarms, metrics, firewall, etc.)
# ... perform your configuration ...

# 3. Create a new instance copying specific settings and wait for it to be ready
cloudamqp instance create --name=staging --plan=bunny-1 --region=amazon-web-services::us-east-1 \
  --copy-from-id=1234 --copy-settings=metrics,firewall,config --wait

# Available settings to copy:
# - alarms: Copy alarm configurations and recipients
# - metrics: Copy metrics configuration
# - logs: Copy log settings
# - firewall: Copy firewall rules
# - config: Copy RabbitMQ configuration
# - definitions: Copy RabbitMQ definitions (queues, exchanges, etc.)
# - plugins: Copy plugin configurations

# Note: Only works between dedicated instances (not shared plans)
```

### Team Setup

```bash
# Invite team members
cloudamqp team invite --email=dev1@company.com --role=devops
cloudamqp team invite --email=dev2@company.com --role=member

# List current team
cloudamqp team list
```

### Monitoring Setup

```bash
# Check available regions and plans
cloudamqp regions --provider=amazon-web-services
cloudamqp plans --backend=rabbitmq

# Create VPC for isolation
cloudamqp vpc create --name=prod-vpc --region=amazon-web-services::us-east-1 --subnet=10.0.0.0/24

# Create instance in VPC
cloudamqp instance create --name=prod-instance --plan=rabbit-1 --region=amazon-web-services::us-east-1 --vpc-id=5678
```

## Error Handling

The CLI provides clear error messages for common issues:

- **401 Unauthorized**: Check your API key configuration
- **404 Not Found**: Verify instance/VPC IDs are correct
- **400 Bad Request**: Check required parameters and formats

## Advanced Usage

### Using Environment Variables

```bash
# Set API key
export CLOUDAMQP_APIKEY="your-api-key"

# Use the CLI without prompts
cloudamqp instance list
```

### Scripting

The CLI is designed for scripting with:

- JSON output for structured data
- Exit codes for success/failure
- `--force` flags to skip confirmations
- Environment variable support

```bash
#!/bin/bash

# Create instance and capture ID
RESULT=$(cloudamqp instance create --name=temp-instance --plan=lemming --region=amazon-web-services::us-east-1)
INSTANCE_ID=$(echo "$RESULT" | jq -r '.id')

# Get instance details
cloudamqp instance get "$INSTANCE_ID"

# Perform operations
cloudamqp instance restart-rabbitmq "$INSTANCE_ID"

# Cleanup
cloudamqp instance delete "$INSTANCE_ID" --force
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For issues and questions:

1. Check the CLI help: `cloudamqp --help`
2. Verify your API key configuration
3. Check the CloudAMQP API documentation
4. Create an issue in this repository

## API Documentation

- [CloudAMQP API](https://docs.cloudamqp.com/api.html)
- [Terraform Provider](https://registry.terraform.io/providers/cloudamqp/cloudamqp/latest/docs)
