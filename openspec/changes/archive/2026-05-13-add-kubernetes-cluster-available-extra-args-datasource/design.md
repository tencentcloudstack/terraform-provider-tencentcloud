# Design: tencentcloud_kubernetes_cluster_available_extra_args Data Source

## API

- **API**: `DescribeClusterAvailableExtraArgs` (TKE v20180525)
- **SDK struct**: `DescribeClusterAvailableExtraArgsRequest` / `DescribeClusterAvailableExtraArgsResponse`
- **SDK already available**: Yes — no SDK changes needed.

### Request Parameters

| Field           | Type   | Required | Description                                      |
|-----------------|--------|----------|--------------------------------------------------|
| ClusterVersion  | string | Yes      | Cluster version, e.g. `1.28.3`                   |
| ClusterType     | string | Yes      | `MANAGED_CLUSTER` or `INDEPENDENT_CLUSTER`        |

### Response Parameters

| Field               | Type              | Description                         |
|---------------------|-------------------|-------------------------------------|
| ClusterVersion      | string            | Cluster version echoed back          |
| ClusterType         | string            | Cluster type echoed back             |
| AvailableExtraArgs  | AvailableExtraArgs| Available args per component         |

`AvailableExtraArgs` contains four optional lists of `Flag` objects:
- `KubeAPIServer`
- `KubeControllerManager`
- `KubeScheduler`
- `Kubelet`

Each `Flag`:
```
Name       *string   // argument name
Type       *string   // argument type
Usage      *string   // argument description
Default    *string   // default value
Constraint *string   // valid range / allowed values
```

## Schema Design

```hcl
data "tencentcloud_kubernetes_cluster_available_extra_args" "example" {
  cluster_version = "1.28.3"
  cluster_type    = "MANAGED_CLUSTER"
}
```

### Input (Required)
- `cluster_version` (string, Required) — cluster version
- `cluster_type`    (string, Required) — cluster type: `MANAGED_CLUSTER` or `INDEPENDENT_CLUSTER`

### Output (Computed)
- `cluster_version` — echoed cluster version
- `cluster_type`    — echoed cluster type
- `available_extra_args` (list, max 1) — nested block containing four component lists:
  - `kube_apiserver` (list of flag objects)
  - `kube_controller_manager` (list of flag objects)
  - `kube_scheduler` (list of flag objects)
  - `kubelet` (list of flag objects)

Each flag object:
- `name`       (string)
- `type`       (string)
- `usage`      (string)
- `default`    (string)
- `constraint` (string)

### Auxiliary
- `result_output_file` (string, Optional) — write results to file

## Resource ID

Use `helper.BuildToken()` to auto-generate a unique ID (no stable natural key needed for a read-only data source with volatile query params).

## Service Layer

Add method to `TkeService`:

```go
func (me *TkeService) DescribeClusterAvailableExtraArgs(ctx context.Context, clusterVersion, clusterType string) (resp *tke.DescribeClusterAvailableExtraArgsResponseParams, err error)
```

- Wraps the SDK call with `resource.Retry(tccommon.ReadRetryTimeout, ...)`.
- Returns nil, nil if response params are nil (not an error).

## File Structure

```
tencentcloud/services/tke/
├── data_source_tc_kubernetes_cluster_available_extra_args.go       # schema + read logic
├── data_source_tc_kubernetes_cluster_available_extra_args.md       # usage docs
├── data_source_tc_kubernetes_cluster_available_extra_args_test.go  # acceptance test
service_tencentcloud_tke.go                                          # new service method
tencentcloud/provider.go                                             # registration
```

## Code Style

Strictly follows `data_source_tc_igtm_instance_list.go` pattern:
- Package-level function returning `*schema.Resource`
- Private `read` function with `defer tccommon.LogElapsed(...)()` and `defer tccommon.InconsistentCheck(...)()`
- `paramMap` pattern for building request
- `resource.Retry` for all API calls
- Nil-guard every pointer dereference from API response
- `helper.BuildToken()` for SetId
- `result_output_file` support via `tccommon.WriteToFile`
