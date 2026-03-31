# Scripting and automation

## JSON output for parsing

All commands support `-o json`. Combine with `jq`:

```bash
# get connection URL for an instance
cloudamqp instance get --id <id> -o json | jq -r '.url'

# find instances that aren't ready
cloudamqp instance list -o json | jq -r '.[] | select(.ready == false) | "\(.id) \(.name)"'

# get IDs matching a tag
cloudamqp instance list -o json | jq -r '.[] | select(.tags[]? == "staging") | .id'
```

## Create and capture instance ID

```bash
# fetch a valid plan and region first
PLAN=$(cloudamqp plans --backend=rabbitmq -o json | jq -r '.[0].name')
REGION=$(cloudamqp regions -o json | jq -r '.[0].id')

RESULT=$(cloudamqp instance create --name=temp --plan="$PLAN" --region="$REGION" -o json)
INSTANCE_ID=$(echo "$RESULT" | jq -r '.id')
```

## Wait for instance readiness

Prefer the built-in flag:

```bash
cloudamqp instance create --name=my-instance --plan=<plan> --region=<region> --wait --wait-timeout=20m
```

Or poll manually:

```bash
while true; do
  STATUS=$(cloudamqp instance get --id "$INSTANCE_ID" -o json | jq -r '.ready')
  [ "$STATUS" = "true" ] && break
  sleep 30
done
```

## Skip confirmations in scripts

```bash
cloudamqp instance delete --id <id> --force
cloudamqp vpc delete --id <id> --force
```

## Batch operations

```bash
# restart all instances tagged "staging"
for ID in $(cloudamqp instance list -o json | jq -r '.[] | select(.tags[]? == "staging") | .id'); do
  cloudamqp instance restart-rabbitmq --id "$ID"
done
```

## Clone an instance with full config

```bash
cloudamqp instance create \
  --name=staging-copy \
  --plan=<plan> \
  --region=<region> \
  --copy-from-id=<source-id> \
  --copy-settings=alarms,metrics,logs,firewall,config,definitions,plugins \
  --wait
```

Only works between dedicated instances (not shared plans).
