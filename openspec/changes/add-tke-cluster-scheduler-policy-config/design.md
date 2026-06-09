# Design: tencentcloud_kubernetes_cluster_scheduler_policy_config

## Architecture

Follows `tencentcloud_config_deliver_config` style (Create = SetId + Update):

```
provider.go
    └─ resource_tc_kubernetes_cluster_scheduler_policy_config.go (Create=SetId+Update, Read, Update, Delete=no-op)
           └─ service_tencentcloud_tke.go (DescribeKubernetesClusterSchedulerPolicy)
```

## Async Pattern

`ModifyClusterSchedulerPolicy` is async. After calling it, poll `DescribeTasks` with filter `ClusterId=<id>` and `Latest=true` until `Tasks[0].LifeState == "done"`. Use `resource.Retry` with `WriteRetryTimeout`.

## Schema

### Required

| Field | Type | Description |
|---|---|---|
| `cluster_id` | String, ForceNew | Cluster ID |

### Optional

| Field | Type | Description |
|---|---|---|
| `scheduler_policy_config` | List (object) | List of scheduler policy configs |
| `extenders` | List (object) | Extender scheduler configs |
| `client_connection` | List, MaxItems=1 | Client connection config (QPS, Burst) |
| `high_performance` | Bool | High performance mode switch |

### Nested: scheduler_policy_config

`scheduler_name`, `plugin_configs` (list: `name`+`args`), `plugin_set` (enabled/disabled list: `name`+`weight`)

### Nested: extenders

`filter_verb`, `prioritize_verb`, `weight`, `preempt_verb`, `node_cache_capable`, `extender_client_config` → `service` → `namespace`/`name`/`port`/`path`/`scheme`

### Nested: client_connection

`qps` (Float), `burst` (Int)

### Computed

`policy` (String) — raw JSON policy string returned by DescribeClusterSchedulerPolicy
