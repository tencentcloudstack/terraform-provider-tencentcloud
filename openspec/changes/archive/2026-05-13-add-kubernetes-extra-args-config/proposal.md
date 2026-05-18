# Add tencentcloud_kubernetes_cluster_extra_args_config Resource

## What

Add a new Terraform resource `tencentcloud_kubernetes_cluster_extra_args_config` for managing TKE cluster custom component extra arguments. This is a singleton config resource per cluster — the resource ID is the `cluster_id`.

## Why

TKE cluster custom extra args for kube-apiserver, kube-controller-manager, kube-scheduler, and etcd are not yet manageable via Terraform. Users need to configure these runtime parameters as infrastructure code to achieve repeatable and auditable cluster configurations.

## APIs Used

| Operation | API | Notes |
|---|---|---|
| Create (initial set) | `ModifyClusterExtraArgs` | Synchronous |
| Read | `DescribeClusterExtraArgs` | Returns current extra args config |
| Update | `ModifyClusterExtraArgs` | Same as Create |
| Delete | No-op | No delete API; config exists as long as cluster exists |

## Resource ID

`cluster_id` (e.g. `cls-5e7wsn94`).
