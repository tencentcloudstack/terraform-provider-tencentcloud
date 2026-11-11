## Context

The CLS (Cloud Log Service) product is already integrated into the Tencent Cloud Terraform provider under `tencentcloud/services/cls/`. A managed resource `tencentcloud_cls_machine_group` exists for create/update/delete, and a data source `tencentcloud_cls_machine_group_configs` exists to query scrape configs of a single machine group. However, there is no data source to list/look up machine groups themselves (e.g. by name, group ID, OS type, or tag). The CLS SDK (vendored at `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016`) exposes `DescribeMachineGroups` which returns a list of `MachineGroupInfo` with pagination.

This is an additive, single-service change that follows the established data-source pattern already used by `data_source_tc_cls_alarm_notices.go` and `data_source_tc_cls_machine_group_configs.go`.

## Goals / Non-Goals

**Goals:**
- Provide a read-only data source `tencentcloud_cls_machine_groups` to query CLS machine groups.
- Support the API's `Filters` input (machineGroupName, machineGroupId, osType, tagKey, tag:tagKey).
- Expose the full `MachineGroupInfo` fields as a computed `machine_groups` list (group_id, group_name, machine_group_type, create_time, tags, auto_update, update_start_time, update_end_time, service_logging, delay_cleanup_time, meta_tags, os_type).
- Implement internal pagination (not exposing `limit`/`offset` to users) to retrieve all matching groups.
- Follow the project's data-source conventions: retry with `tccommon.ReadRetryTimeout`, `tccommon.RetryError`, `defer tccommon.LogElapsed()` / `InconsistentCheck()`, `result_output_file` support, and a deterministic `d.SetId()` derived from the returned group IDs.
- Register the data source in the provider and generate docs via `make doc`.

**Non-Goals:**
- Creating, updating, or deleting machine groups (handled by `tencentcloud_cls_machine_group` resource).
- Exposing `limit`/`offset` as user-facing schema fields (pagination is internal per project conventions).
- Modifying any existing resource or data source schema.
- Terraform acceptance test suite for the new data source (unit tests with gomonkey mocks are used instead, per project rules for new data sources).

## Decisions

### Decision 1: Single Read function, no service-layer pagination helper duplication
The data source's Read function will call `DescribeMachineGroups` directly through `me.client.UseClsClient().DescribeMachineGroups(request)` inside a `resource.Retry(tccommon.ReadRetryTimeout, ...)` block, and paginate with a `for` loop in the Read function itself (incrementing `offset` by `limit` until a page returns fewer than `limit` items).

**Rationale**: This mirrors the established pattern in `data_source_tc_cls_alarm_notices.go` (service layer) and keeps the logic close to the schema mapping. The project rule "data source pagination: do not expose limit/offset to users, internally implement automatic pagination to fetch all data" is satisfied by using a fixed internal `limit` (the API maximum 100) and looping.

**Alternative considered**: Putting pagination in a `service_tencentcloud_cls.go` method returning `[]*cls.MachineGroupInfo`. Rejected to keep the change minimal and consistent with the machine_group_configs data source which calls a service method. We will still add a thin service-layer wrapper `DescribeClsMachineGroupsByFilter` to keep retry/pagination encapsulated and testable, consistent with `DescribeClsMachineGroupConfigsByFilter`.

### Decision 2: Schema field mapping
Input schema:
- `filters` (Optional, TypeList of `name`+`values`) — maps to `request.Filters []*cls.Filter`.
- `result_output_file` (Optional, TypeString) — save results to file.

Output schema (`machine_groups` computed list, fields flattened — no extra `machine_groups` wrapper nesting beyond the list element):
- `group_id` (string)
- `group_name` (string)
- `machine_group_type` (list of one: `type` string, `values` set of string) — mirrors `MachineGroupTypeInfo`.
- `create_time` (string)
- `tags` (list of `key`/`value`)
- `auto_update` (string)
- `update_start_time` (string)
- `update_end_time` (string)
- `service_logging` (bool)
- `delay_cleanup_time` (int)
- `meta_tags` (list of `key`/`value`)
- `os_type` (int)
- `total_count` (int, top-level computed)

**Rationale**: Directly reflects `MachineGroupInfo` struct fields. `machine_group_type` is a single object so it is modeled as a TypeList with one element (consistent with other CLS schemas). Per project rule #13, list data is flattened — the `machine_groups` list elements contain the fields directly (no nested `xxx_set` wrapper).

### Decision 3: Empty-response handling for data source
Per project rule #14 for `RESOURCE_KIND_DATASOURCE`, inside the retry block we MUST check for empty API responses (`response == nil || response.Response == nil`) and return `NonRetryableError` instead of clearing the ID. This prevents transient API blips from wiping local state. On the outer retry-exhausted path, log `[DATASOURCE] read empty, skip SetId` and return the error without calling `d.SetId("")`.

### Decision 4: ID computation
`d.SetId(helper.DataResourceIdsHash(ids))` where `ids` are the returned `GroupId` values, consistent with `data_source_tc_cls_machine_group_configs.go`. This gives a stable, deterministic ID for the data source instance.

### Decision 5: Tests use gomonkey mocks
Per project rules, new data sources use gomonkey to mock the CLS API client (no terraform test suite). Tests will mock `DescribeMachineGroups` to return canned `MachineGroupInfo` slices and assert the Read function maps fields correctly, including nil-safe handling and pagination across two pages.

## Risks / Trade-offs

- [Risk: API response nil/empty during transient blip] → Mitigation: return `NonRetryableError` inside retry so the outer retry continues; do not clear state ID. Log for traceability.
- [Risk: Pagination edge cases (exactly N*100 items)] → Mitigation: loop breaks when a page returns fewer items than `limit`; also bounded by `TotalCount`.
- [Risk: Field name/type mismatches with SDK] → Mitigation: verified against vendored `MachineGroupInfo` struct (`GroupId`, `GroupName`, `MachineGroupType`, `CreateTime`, `Tags`, `AutoUpdate`, `UpdateStartTime`, `UpdateEndTime`, `ServiceLogging`, `DelayCleanupTime`, `MetaTags`, `OSType`).
- [Risk: Documentation drift] → Mitigation: docs generated via `make doc` from the `.md` source in the finalize phase, per project rules.
- [Trade-off: internal pagination hides limit/offset from users] → Accepted per project convention; users who need raw paging can use the API directly.
