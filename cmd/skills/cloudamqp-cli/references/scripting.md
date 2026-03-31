# Scripting and automation

## JSON output for parsing

Read commands (`list`, `get`, `plans`, `regions`) support `-o json`. All values come out as strings.

```bash
# get connection URL for an instance (masked; add --show-url for full URL)
cloudamqp instance get --id <id> -o json | jq -r '.url'

# find instances that aren't ready (requires --details; ready is "Yes"/"No" string)
cloudamqp instance list --details -o json | jq -r '.[] | select(.ready == "No") | "\(.id) \(.name)"'

# get IDs matching a tag (requires --details; tags is a comma-joined string)
cloudamqp instance list --details -o json | jq -r '.[] | select(.tags | split(",") | map(ltrimstr(" ")) | contains(["staging"])) | .id'
```

## Create and capture instance ID

`instance create` prints a human-readable prefix before the JSON, so pipe through `tail -n +2`:

```bash
# fetch a valid plan and region first
PLAN=$(cloudamqp plans --backend=rabbitmq -o json | jq -r '.[0].name')
REGION=$(cloudamqp regions -o json | jq -r '.[0].id')

RESULT=$(cloudamqp instance create --name=temp --plan="$PLAN" --region="$REGION" | tail -n +2)
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
  [ "$STATUS" = "Yes" ] && break
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
# restart all instances tagged "staging" (--details required for tags field)
for ID in $(cloudamqp instance list --details -o json | jq -r '.[] | select(.tags | split(",") | map(ltrimstr(" ")) | contains(["staging"])) | .id'); do
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
