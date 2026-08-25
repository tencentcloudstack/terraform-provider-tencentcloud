## 1. Schema Changes (ForceNew Removal)

- [x] 1.1 Remove `ForceNew: true` from `description` field in schema
- [x] 1.2 Remove `ForceNew: true` from `rule` TypeList in schema
- [x] 1.3 Remove `ForceNew: true` from `rule.dest_namespace` field in schema
- [x] 1.4 Remove `ForceNew: true` from `rule.override` field in schema
- [x] 1.5 Remove `ForceNew: true` from `rule.filters` TypeList and its `type`/`value` children in schema
- [x] 1.6 Remove `ForceNew: true` from `rule.deletion` field in schema
- [x] 1.7 Verify `source_registry_id`, `destination_registry_id`, `rule.name`, `destination_region_id`, and `peer_replication_option` retain `ForceNew: true`

## 2. Core Implementation (Update Method)

- [x] 2.1 Implement `resourceTencentCloudTcrReplicationUpdate` function in `resource_tc_tcr_replication.go`
- [x] 2.2 Parse resource ID to extract `sourceRegistryId` and `ruleName` in Update
- [x] 2.3 Build `ModifyReplicationRequest` from schema data (d.GetOk/helper.InterfacesHeadMap) in Update
- [x] 2.4 Wrap `ModifyReplicationWithContext` API call with `resource.Retry(tccommon.WriteRetryTimeout, ...)` and `tccommon.RetryError`
- [x] 2.5 Call `resourceTencentCloudTcrReplicationRead` at the end of Update to refresh state
- [x] 2.6 Register `Update: resourceTencentCloudTcrReplicationUpdate` in `ResourceTencentCloudTcrReplication()` return value

## 3. Unit Tests

- [x] 3.1 Add unit test cases for the Update method in `resource_tc_tcr_replication_test.go` using gomonkey to mock `ModifyReplicationWithContext`
- [x] 3.2 Add test case for successful update of description
- [x] 3.3 Add test case for successful update of rule filters
- [x] 3.4 Add test case for update with API error (retry exhausted)

## 4. Documentation and Finalization

- [x] 4.1 Update `resource_tc_tcr_replication.md` to reflect updatable fields
- [ ] 4.2 Run `make doc` to regenerate `website/docs/` documentation (done in finalize phase)
- [ ] 4.3 Run `gofmt` on modified Go files (done in finalize phase)