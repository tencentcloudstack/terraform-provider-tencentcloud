## Why

The `tencentcloud_tcr_replication` resource currently has all schema fields marked as `ForceNew: true`, which means any modification to the resource (even changing the rule description or filter configuration) triggers a destroy-and-recreate cycle. This is disruptive in production environments where users need to update replication rules without losing existing synchronization state. The TCR API already provides a `ModifyReplication` endpoint for in-place updates, and the Terraform resource should leverage it to support smooth updates.

## What Changes

- Add an `Update` method (`resourceTencentCloudTcrReplicationUpdate`) to the resource, using the `ModifyReplication` API
- Remove `ForceNew: true` from fields that are supported by the `ModifyReplication` API: `description`, `rule.dest_namespace`, `rule.override`, `rule.filters`, `rule.deletion`
- Keep `ForceNew: true` for fields that cannot be modified via the API: `source_registry_id`, `destination_registry_id`, `rule.name` (used as the rule identifier), `destination_region_id`, and `peer_replication_option` (cross-account settings)
- Register the `Update` function in the resource schema
- Add unit test coverage for the new Update logic

## Capabilities

### New Capabilities
- `tcr-replication-update`: Enable in-place update of TCR replication rules via the `ModifyReplication` API, allowing users to modify description, destination namespace, override behavior, filters, and deletion settings without destroying and recreating the resource.

### Modified Capabilities
<!-- No existing capabilities are being modified at the spec level -->

## Impact

- **Affected code**: `tencentcloud/services/tcr/resource_tc_tcr_replication.go` (add Update method, modify ForceNew flags), `tencentcloud/services/tcr/resource_tc_tcr_replication_test.go` (add unit tests)
- **Affected documentation**: `tencentcloud/services/tcr/resource_tc_tcr_replication.md` (update to reflect updatable fields)
- **API dependency**: `ModifyReplication` from `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924` (already present in vendor)
- **Backward compatibility**: Fully backward compatible — existing configurations will continue to work; the only change is that previously ForceNew fields now support in-place updates