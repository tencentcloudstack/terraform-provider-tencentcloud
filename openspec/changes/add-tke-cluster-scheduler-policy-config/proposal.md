# Add tencentcloud_kubernetes_cluster_scheduler_policy_config Resource

## What

Add a new Terraform resource `tencentcloud_kubernetes_cluster_scheduler_policy_config` for managing TKE cluster scheduler policy configuration. This is a singleton config resource per cluster — the resource ID is the `cluster_id`.

## Why

TKE cluster scheduler policies (plugin configs, extenders, client connection) are not yet manageable via Terraform. Users need to configure custom scheduling strategies as infrastructure code.

## APIs Used

| Operation | API | Notes |
|---|---|---|
| Create (initial set) | `ModifyClusterSchedulerPolicy` | Async — poll `DescribeTasks` until `LifeState == "done"` |
| Read | `DescribeClusterSchedulerPolicy` | Returns current policy config |
| Update | `ModifyClusterSchedulerPolicy` | Same as Create, async |
| Delete | No-op | No delete API; config exists as long as cluster exists |

## Resource ID

`cluster_id` (e.g. `cls-5e7wsn94`).
