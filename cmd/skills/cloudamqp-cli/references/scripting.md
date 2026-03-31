# Scripting and automation

## JSON output for parsing

All commands support `-o json` for structured output. Combine with `jq` for extraction:

```bash
# get connection URL for an instance
cloudamqp instance get --id 1234 -o json | jq -r '.url'

# list instance IDs matching a tag
cloudamqp instance list -o json | jq -r '.[] | select(.tags[]? == "production") | .id'

# get all instance names and plans
cloudamqp instance list -o json | jq -r '.[] | "\(.id) \(.name) \(.plan)"'
```

## Create and capture instance ID

```bash
RESULT=$(cloudamqp instance create --name=temp --plan=lemming --region=amazon-web-services::us-east-1 -o json)
INSTANCE_ID=$(echo "$RESULT" | jq -r '.id')
echo "Created instance: $INSTANCE_ID"
```

## Wait for instance readiness

Use the built-in `--wait` flag (default timeout: 15 minutes):

```bash
cloudamqp instance create --name=my-instance --plan=bunny-1 \
  --region=amazon-web-services::us-east-1 --wait --wait-timeout=20m
```

Or poll manually:

```bash
while true; do
  STATUS=$(cloudamqp instance get --id "$INSTANCE_ID" -o json | jq -r '.ready')
  [ "$STATUS" = "true" ] && break
  sleep 30
done
```

## Skip confirmations

Use `--force` on destructive commands:

```bash
cloudamqp instance delete --id 1234 --force
```

## Environment-based configuration

```bash
export CLOUDAMQP_APIKEY="your-api-key"
cloudamqp instance list  # no prompts
```

## Batch operations

```bash
# restart all instances tagged "staging"
for ID in $(cloudamqp instance list -o json | jq -r '.[] | select(.tags[]? == "staging") | .id'); do
  echo "Restarting instance $ID"
  cloudamqp instance restart-rabbitmq --id "$ID"
done
```

## Clone an instance with full config

```bash
cloudamqp instance create \
  --name=staging-copy \
  --plan=bunny-1 \
  --region=amazon-web-services::us-east-1 \
  --copy-from-id=1234 \
  --copy-settings=alarms,metrics,logs,firewall,config,definitions,plugins \
  --wait
```
