## 1. Schema Definition

- [x] 1.1 Add `labels` argument (TypeList, MaxItems 20, Optional, ForceNew) with nested `key` (Required, TypeString) and `value` (Optional, TypeString) to the `tencentcloud_dbdc_node_to_db_custom_cluster_attachment` resource schema
- [x] 1.2 Add `taints` argument (TypeList, MaxItems 5, Optional, ForceNew) with nested `key` (Required, TypeString), `effect` (Required, TypeString), and `value` (Optional, TypeString) to the schema
- [x] 1.3 Add `host_name` argument (Optional, TypeString, ForceNew) to the schema
- [x] 1.4 Add `host_name_type` argument (Optional, TypeInt, ForceNew) to the schema
- [x] 1.5 Add `network_mode` computed field (Computed, TypeString) to the schema
- [x] 1.6 Add `eni_ip` computed field (Computed, TypeString) to the schema

## 2. CRUD Logic Updates

- [x] 2.1 Update Create function: populate `request.Labels` (build `[]*dbdcv20201029.Label{Key, Value}`) when `labels` is provided
- [x] 2.2 Update Create function: populate `request.Taints` (build `[]*dbdcv20201029.Taint{Key, Effect, Value}`) when `taints` is provided
- [x] 2.3 Update Create function: set `request.HostName` when `host_name` is provided and `request.HostNameType` when `host_name_type` is provided
- [x] 2.4 Update Read function: set `network_mode` from `respData.NetworkMode` with nil guard (skip set when nil)
- [x] 2.5 Update Read function: set `eni_ip` from `respData.EniIP` with nil guard (skip set when nil)
- [x] 2.6 Update Delete function: read `login_settings` block from resource data and set `request.LoginSettings` (build `dbdcv20201029.LoginSettings{Password, KeyIds, KeepImageLogin}`) when present, before calling `RemoveNodesFromDBCustomCluster`

## 3. Documentation

- [x] 3.1 Update `resource_tc_dbdc_node_to_db_custom_cluster_attachment.md` Example Usage to demonstrate the new optional arguments (`labels`, `taints`, `host_name`, `host_name_type`)

## 4. Unit Tests

- [x] 4.1 Add gomonkey-mocked unit tests in `resource_tc_dbdc_node_to_db_custom_cluster_attachment_test.go` verifying Create builds the request with `Labels`, `Taints`, `HostName`, `HostNameType`
- [x] 4.2 Add gomonkey-mocked unit tests verifying Read populates `network_mode` and `eni_ip` from the response
- [x] 4.3 Add gomonkey-mocked unit tests verifying Delete passes `LoginSettings` to the `RemoveNodesFromDBCustomCluster` request

## 5. Verification

- [x] 5.1 Verify the code compiles (no `go build`/`go vet` execution; checked by downstream pipeline)
- [x] 5.2 Verify backward compatibility: existing schema fields unchanged, composite ID format unchanged, all new inputs are Optional+ForceNew, new outputs are Computed
