# Instance upgrades and maintenance

## Check available upgrade versions

```bash
cloudamqp instance upgrade-versions --id 1234
cloudamqp instance nodes versions --id 1234
```

## Upgrade paths

### Upgrade everything (recommended)

Upgrades both Erlang and RabbitMQ/LavinMQ to the latest compatible versions:

```bash
cloudamqp instance upgrade-all --id 1234
```

### Upgrade individually

```bash
# erlang first, then RabbitMQ
cloudamqp instance upgrade-erlang --id 1234
cloudamqp instance upgrade-rabbitmq --id 1234 --version=3.13.0
```

Always upgrade Erlang before RabbitMQ when doing both separately.

## Restart operations

```bash
# restart RabbitMQ process (rolling for HA clusters)
cloudamqp instance restart-rabbitmq --id 1234

# restart specific nodes only
cloudamqp instance restart-rabbitmq --id 1234 --nodes=node1,node2

# full cluster restart
cloudamqp instance restart-cluster --id 1234

# restart management interface only
cloudamqp instance restart-management --id 1234
```

## Start, stop, and reboot

```bash
cloudamqp instance stop --id 1234
cloudamqp instance start --id 1234
cloudamqp instance reboot --id 1234

# cluster-level
cloudamqp instance stop-cluster --id 1234
cloudamqp instance start-cluster --id 1234
```

## Disk resize

```bash
# resize disk (may require downtime)
cloudamqp instance resize-disk --id 1234 --disk-size=100 --allow-downtime
```

## Maintenance workflow

A typical maintenance sequence:

```bash
# 1. check current state
cloudamqp instance get --id 1234
cloudamqp instance nodes list --id 1234

# 2. check available upgrades
cloudamqp instance upgrade-versions --id 1234

# 3. perform upgrade
cloudamqp instance upgrade-all --id 1234

# 4. verify after upgrade
cloudamqp instance nodes list --id 1234
cloudamqp instance nodes versions --id 1234
```

## Important notes

- Upgrade operations are **async**: they return immediately but run in the background.
- Poll `instance get --id <id>` to check when the instance is ready again.
- Plan changes via `instance update --id <id> --plan=<plan>` may cause brief downtime.
- `restart-rabbitmq` does a rolling restart on HA clusters (minimal disruption).
- `restart-cluster` restarts all nodes simultaneously (causes downtime).
