# Proposal: Add tencentcloud_kubernetes_cluster_available_extra_args Data Source

## What

Add a new Terraform data source `tencentcloud_kubernetes_cluster_available_extra_args` that queries the available custom extra arguments for TKE (Tencent Kubernetes Engine) cluster components, by calling the `DescribeClusterAvailableExtraArgs` API.

## Why

TKE supports customizing startup arguments for core Kubernetes components (kube-apiserver, kube-controller-manager, kube-scheduler, kubelet). Before applying custom extra args to a cluster, operators need to know which arguments are supported for a given cluster version and type. Currently there is no Terraform data source exposing this information, forcing users to consult the console or call the API manually.

This data source fills that gap, enabling:
- Infrastructure-as-code workflows that validate available extra args before applying them.
- Dynamic configuration: use data source output as input to `tencentcloud_kubernetes_cluster_extra_args_config` or cluster resources.

## Scope

- **New file**: `tencentcloud/services/tke/data_source_tc_kubernetes_cluster_available_extra_args.go`
- **New file**: `tencentcloud/services/tke/data_source_tc_kubernetes_cluster_available_extra_args.md`
- **New file**: `tencentcloud/services/tke/data_source_tc_kubernetes_cluster_available_extra_args_test.go`
- **Modified file**: `tencentcloud/provider.go` — register the new data source
- **Modified file**: `tencentcloud/services/tke/service_tencentcloud_tke.go` — add service method wrapping the API call

## Out of Scope

- No create / update / delete operations (read-only data source).
- No SDK changes required (SDK already contains `DescribeClusterAvailableExtraArgs`).
