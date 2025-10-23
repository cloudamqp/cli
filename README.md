# CloudAMQP CLI

A command line interface for the CloudAMQP API that provides complete management of CloudAMQP instances, VPCs, and instance-specific operations.

## Features

- **Unified API**: Single API key manages all operations through the customer API
- **Simple Configuration**: Plain text API key storage in `~/.cloudamqprc`
- **Flag-Based Commands**: Clean command structure with `--id` flags for instance operations
- **User-Friendly**: Clear help messages, examples, and safety confirmations
- **Error Handling**: Proper API error extraction and display

## Installation

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

## Commands

### Instance Management

Manage CloudAMQP instances using your main API key.

```bash
# Create a new instance
cloudamqp instance create --name=my-instance --plan=bunny-1 --region=amazon-web-services::us-east-1

# List all instances
cloudamqp instance list

# Get instance details
cloudamqp instance get --id 1234

# Update instance properties
cloudamqp instance update --id 1234 --name=new-name --plan=rabbit-1

# Resize instance disk
cloudamqp instance resize --id 1234 --disk-size=100 --allow-downtime

# Delete instance (with confirmation)
cloudamqp instance delete --id 1234
```

### VPC Management

Manage Virtual Private Clouds.

```bash
# Create a VPC
cloudamqp vpc create --name=my-vpc --region=amazon-web-services::us-east-1 --subnet=10.56.72.0/24

# List all VPCs
cloudamqp vpc list

# Get VPC details
cloudamqp vpc get --id 5678

# Update VPC
cloudamqp vpc update --id 5678 --name=new-vpc-name

# Delete VPC (with confirmation)
cloudamqp vpc delete --id 5678
```

### Instance-Specific Management

Manage specific instances using the unified API. All commands use `--id` flag to specify the instance.

#### Node Management

```bash
# List nodes in an instance
cloudamqp instance nodes list --id 1234

# Get available versions for upgrade
cloudamqp instance nodes versions --id 1234
```

#### Plugin Management

```bash
# List available RabbitMQ plugins
cloudamqp instance plugins list --id 1234
```

#### RabbitMQ Configuration

```bash
# List all configuration settings
cloudamqp instance config list --id 1234

# Get specific configuration setting
cloudamqp instance config get --id 1234 --key tcp_listen_options

# Set configuration setting
cloudamqp instance config set --id 1234 --key tcp_listen_options --value '[{"port": 5672}]'
```

#### Instance Actions

```bash
# Restart RabbitMQ
cloudamqp instance actions restart-rabbitmq --id 1234
cloudamqp instance actions restart-rabbitmq --id 1234 --nodes=node1,node2

# Cluster operations
cloudamqp instance actions restart-cluster --id 1234
cloudamqp instance actions stop-cluster --id 1234
cloudamqp instance actions start-cluster --id 1234

# Instance lifecycle
cloudamqp instance actions stop --id 1234
cloudamqp instance actions start --id 1234
cloudamqp instance actions reboot --id 1234

# Management interface
cloudamqp instance actions restart-management --id 1234

# Upgrades (asynchronous operations)
cloudamqp instance actions upgrade-erlang --id 1234
cloudamqp instance actions upgrade-rabbitmq --id 1234 --version=3.10.7
cloudamqp instance actions upgrade-all --id 1234

# Get target upgrade versions
cloudamqp instance actions upgrade-versions --id 1234

# Toggle features
cloudamqp instance actions toggle-hipe --id 1234 --enable=true
cloudamqp instance actions toggle-firehose --id 1234 --enable=true --vhost=/
```

#### Account Operations

```bash
# Rotate instance password
cloudamqp instance account rotate-password --id 1234

# Rotate instance API key
cloudamqp instance account rotate-apikey --id 1234
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

# Rotate main API key
cloudamqp rotate-key
```

## Examples

### Complete Workflow

```bash
# 1. Create an instance
cloudamqp instance create --name=production --plan=bunny-1 --region=amazon-web-services::us-east-1

# 2. Get instance details
cloudamqp instance get --id 1234

# 3. Check instance nodes
cloudamqp instance nodes list --id 1234

# 4. List RabbitMQ configuration
cloudamqp instance config list --id 1234

# 5. Install plugins (if needed)
cloudamqp instance plugins list --id 1234

# 6. Restart RabbitMQ
cloudamqp instance actions restart-rabbitmq --id 1234

# 7. Upgrade when needed
cloudamqp instance actions upgrade-all --id 1234
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
cloudamqp instance get --id "$INSTANCE_ID"

# Perform operations
cloudamqp instance actions restart-rabbitmq --id "$INSTANCE_ID"

# Cleanup
cloudamqp instance delete --id "$INSTANCE_ID" --force
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