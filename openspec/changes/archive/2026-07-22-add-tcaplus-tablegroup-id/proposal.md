## Why

The TcaplusDB `CreateTableGroup` API supports an optional `TableGroupId` parameter that allows users to specify a custom table group ID at creation time (when not specified, the API auto-increments the ID). However, the Terraform resource `tencentcloud_tcaplus_tablegroup` does not expose this parameter, so users cannot set a fixed table group ID through Terraform. Users who need a deterministic, pre-known table group ID (e.g., for cross-environment consistency or downstream automation referencing the ID) are forced to use the console or API directly.

## What Changes

- Add `table_group_id` (Optional, immutable after creation) parameter to the `tencentcloud_tcaplus_tablegroup` resource. When set, it is passed to the `CreateTableGroup` API request as `TableGroupId`. When not set, the API defaults to auto-increment mode.
- The `CreateTableGroup` API response returns the created `TableGroupId` (both for user-specified and auto-incremented cases), which is read back into state.
- Treat `table_group_id` as immutable after creation: the `ModifyTableGroupName` API only accepts `TableGroupName` (not `TableGroupId`), so changes to `table_group_id` after creation SHALL be rejected via the `immutableArgs` check in the Update function.
- The existing resource ID format (`clusterId:tableGroupId`) is preserved; the `table_group_id` schema field is the human-facing mirror of the group id portion already embedded in `d.Id()`.

## Capabilities

### New Capabilities
- `tcaplus-tablegroup-id`: Enable the optional `table_group_id` parameter on the `tencentcloud_tcaplus_tablegroup` resource to allow users to specify a custom table group ID when creating a TcaplusDB table group.

### Modified Capabilities
<!-- No existing specs require modification -->

## Impact

- **Affected files:**
  - `tencentcloud/services/tcaplusdb/resource_tc_tcaplus_tablegroup.go` — add `table_group_id` schema field (Optional, immutable), wire through Create flow, add immutable args check in Update, set the field in Read
  - `tencentcloud/services/tcaplusdb/service_tencentcloud_tcaplus.go` — update `CreateGroup` service function to accept and pass the optional `TableGroupId` to the `CreateTableGroup` API request
  - `tencentcloud/services/tcaplusdb/resource_tc_tcaplus_tablegroup.md` — update documentation with `table_group_id` usage example
- **SDK dependency:** `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcaplusdb/v20190823` — `CreateTableGroupRequest` already includes the optional `TableGroupId *string` field, and `CreateTableGroupResponse` returns `TableGroupId *string`. No SDK update required.
- **Backward compatibility:** fully backward compatible — the new parameter is Optional; existing configurations that do not set `table_group_id` continue to use API auto-increment behavior unchanged.
- **API constraints:** `TableGroupId` is only accepted by `CreateTableGroup` (the `ModifyTableGroupName` and `DeleteTableGroup` APIs do not accept it as a writable attribute), so this parameter is immutable after creation.
