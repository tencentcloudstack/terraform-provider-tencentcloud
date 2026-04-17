## 1. Service Layer

- [x] 1.1 Append `DescribeKubernetesClusterSchedulerPolicy()` to `service_tencentcloud_tke.go`

## 2. Resource Implementation

- [x] 2.1 Create `resource_tc_kubernetes_cluster_scheduler_policy_config.go` with schema
- [x] 2.2 Create handler: `d.SetId(clusterId)` then call Update
- [x] 2.3 Read handler: call `DescribeClusterSchedulerPolicy`, populate all fields
- [x] 2.4 Update handler: call `ModifyClusterSchedulerPolicy`, then poll `DescribeTasks` until done
- [x] 2.5 Delete handler: no-op

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_kubernetes_cluster_scheduler_policy_config` in `provider.go`

## 4. Documentation & Tests

- [x] 4.1 Create `resource_tc_kubernetes_cluster_scheduler_policy_config.md`
- [x] 4.2 Create `resource_tc_kubernetes_cluster_scheduler_policy_config_test.go`

## 5. Refinements

- [x] 5.1 Change `plugin_set.enabled` and `plugin_set.disabled` from TypeList to TypeSet (hash by `name`), update `buildSchedulerPolicyConfigList` to use `.(*schema.Set).List()`
- [x] 5.2 `plugin_configs.args` transparent base64: user inputs raw JSON, TF encodes on write and decodes on read; update schema description, `buildSchedulerPolicyConfigList`, and `flattenSchedulerPolicyConfigList`

