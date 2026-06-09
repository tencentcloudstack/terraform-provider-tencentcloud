## 1. SDK Extension

- [x] 1.1 Add `ModifyClusterExtraArgs` request/response structs to `vendor/.../tke/v20180525/models.go`
- [x] 1.2 Add `NewModifyClusterExtraArgsRequest`, `ModifyClusterExtraArgs`, and `ModifyClusterExtraArgsWithContext` to `vendor/.../tke/v20180525/client.go`

## 2. Service Layer

- [x] 2.1 Append `DescribeKubernetesClusterExtraArgsConfig()` to `service_tencentcloud_tke.go`

## 3. Resource Implementation

- [x] 3.1 Create `resource_tc_kubernetes_cluster_extra_args_config.go` with schema
- [x] 3.2 Create handler: `d.SetId(clusterId)` then call Update
- [x] 3.3 Read handler: call `DescribeKubernetesClusterExtraArgsConfig`, populate all fields
- [x] 3.4 Update handler: call `ModifyClusterExtraArgs` with retry
- [x] 3.6 Update handler: poll `DescribeTasks` until `LifeState == "done"` after `ModifyClusterExtraArgs`
- [x] 3.5 Delete handler: no-op

## 4. Provider Registration

- [x] 4.1 Register `tencentcloud_kubernetes_cluster_extra_args_config` in `provider.go`

## 5. Documentation & Tests

- [x] 5.1 Create `resource_tc_kubernetes_cluster_extra_args_config.md`
- [x] 5.2 Create `resource_tc_kubernetes_cluster_extra_args_config_test.go`
