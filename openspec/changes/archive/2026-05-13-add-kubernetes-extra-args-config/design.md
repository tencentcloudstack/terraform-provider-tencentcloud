# Design: tencentcloud_kubernetes_cluster_extra_args_config

## Architecture

Follows `tencentcloud_waf_owasp_rule_status_config` style (Create = SetId + Update):

```
provider.go
    └─ resource_tc_kubernetes_cluster_extra_args_config.go (Create=SetId+Update, Read, Update, Delete=no-op)
           └─ service_tencentcloud_tke.go (DescribeKubernetesClusterExtraArgsConfig)
```

## SDK Extension

`ModifyClusterExtraArgs` does not yet exist in the vendored SDK. Add it to:
- `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525/models.go` (request/response structs)
- `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525/client.go` (client methods)

`DescribeClusterExtraArgs` already exists in the SDK.

## Schema

### Required

| Field | Type | Description |
|---|---|---|
| `cluster_id` | String, ForceNew | Cluster ID. Only managed (托管) clusters are supported. |

### Optional (ClusterExtraArgs)

| Field | Type | Description |
|---|---|---|
| `kube_apiserver` | List(String) | Custom args for kube-apiserver, format: ["k1=v1", "k2=v2"] |
| `kube_controller_manager` | List(String) | Custom args for kube-controller-manager |
| `kube_scheduler` | List(String) | Custom args for kube-scheduler |
| `etcd` | List(String) | Custom args for etcd (standalone clusters only) |

All four fields map directly to `ClusterExtraArgs` in the `ModifyClusterExtraArgs` request.

## Sync Pattern

`ModifyClusterExtraArgs` is synchronous. No polling required.

## Retry Pattern

All API calls wrapped with `resource.Retry(tccommon.WriteRetryTimeout, ...)` for writes and `resource.Retry(tccommon.ReadRetryTimeout, ...)` for reads.

## File Naming Convention

Resource files follow the pattern of `resource_tc_config_compliance_pack.md` / `resource_tc_config_compliance_pack_test.go`:
- `resource_tc_kubernetes_cluster_extra_args_config.go`
- `resource_tc_kubernetes_cluster_extra_args_config.md`
- `resource_tc_kubernetes_cluster_extra_args_config_test.go`
