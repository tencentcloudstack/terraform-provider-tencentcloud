## Context

The `tencentcloud_tcaplus_tablegroup` resource currently supports creating TcaplusDB table groups with `cluster_id` and `tablegroup_name` parameters. The TcaplusDB `CreateTableGroup` API also accepts an optional `TableGroupId` parameter that lets users specify a custom table group ID at creation time (when not specified, the API auto-increments the ID), and the response returns the created `TableGroupId`.

**Current state:**
- Resource file: `tencentcloud/services/tcaplusdb/resource_tc_tcaplus_tablegroup.go`
- Service layer: `tencentcloud/services/tcaplusdb/service_tencentcloud_tcaplus.go` (`CreateGroup` function at line 277)
- SDK: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcaplusdb/v20190823`

**API behavior analysis:**

| API | TableGroupId in Request | TableGroupId in Response | Notes |
|-----|-------------------------|--------------------------|-------|
| `CreateTableGroup` | Yes (`TableGroupId *string`, optional — user-specified or auto-increment) | Yes (`TableGroupId *string` — created id) | Writable at creation only |
| `DescribeTableGroups` | N/A | Yes (`TableGroupInfo.TableGroupId`) | Readable for state refresh |
| `ModifyTableGroupName` | Yes (as locator — "待修改名称的表格组ID") | N/A | `TableGroupId` identifies the target, not a writable attribute |
| `DeleteTableGroup` | Yes (as locator) | N/A | `TableGroupId` identifies the target |

**Key constraint:** `TableGroupId` is only writable through `CreateTableGroup`. The `ModifyTableGroupName` API uses `TableGroupId` solely to locate the table group whose name is being changed — it cannot change the id itself. Therefore `table_group_id` is immutable after creation.

## Goals / Non-Goals

**Goals:**
- Add `table_group_id` (Optional, immutable after creation) parameter to `tencentcloud_tcaplus_tablegroup`
- Pass `TableGroupId` to the `CreateTableGroup` API request when the user specifies `table_group_id`
- Read `TableGroupId` from the `DescribeTableGroups` API response (`TableGroupInfo.TableGroupId`) to support state refresh and import
- Implement immutable args check in the Update function so changes to `table_group_id` return a clear error (not ForceNew, to avoid silent resource destruction)
- Maintain full backward compatibility — existing configurations that omit `table_group_id` continue to use API auto-increment behavior unchanged

**Non-Goals:**
- Making `table_group_id` updatable (the `ModifyTableGroupName` API does not support changing the id)
- Adding `table_group_id` to the `tencentcloud_tcaplus_tablegroups` data source (out of scope for this single-parameter change)

## Decisions

### Decision 1: `table_group_id` is immutable (not ForceNew)

**Rationale:** The `ModifyTableGroupName` API only accepts `TableGroupName` as a writable attribute; `TableGroupId` is used only to locate the target table group. Instead of using `ForceNew: true` (which would silently destroy and recreate the resource, risking data loss on the table group), we use an `immutableArgs` array pattern in the Update function that returns a clear error when `table_group_id` changes. This gives users a better error message than silent destruction.

### Decision 2: Optional (not Computed) for `table_group_id`

**Rationale:** The API auto-increments the id when `TableGroupId` is not supplied in the create request. Since the value is always echoed back in the `CreateTableGroup` response and in `DescribeTableGroups`, the Read function sets it into state. Marking the field `Optional` (without `Computed`) is acceptable because the field is always populated by the Read path after creation. This keeps the schema simple and consistent with the existing resource style.

### Decision 3: Read `TableGroupId` from Describe API response

**Rationale:** The `DescribeTableGroups` response's `TableGroupInfo` struct includes `TableGroupId`. Reading this value enables proper state refresh and supports imported resources where the user did not specify the id at create time.

### Decision 4: Preserve existing resource ID format

**Rationale:** The existing resource ID format is `clusterId:tableGroupId` (set via `d.SetId(fmt.Sprintf("%s:%s", clusterId, groupId))`). The `table_group_id` schema field is the human-facing mirror of the group id portion already embedded in `d.Id()`. We keep `d.Id()` unchanged to preserve backward compatibility with existing state files.

### Decision 5: Update `CreateGroup` service function signature

**Rationale:** The `CreateGroup` service function (line 277) currently accepts `(ctx, id, groupName)` and returns `(groupId string, errRet error)`. We extend its signature to accept an optional `tableGroupId string` parameter, passing it to `request.TableGroupId` when non-empty. This is the minimal change needed and keeps the function returning the created `groupId` (which already comes from `response.Response.TableGroupId`).

## Risks / Trade-offs

- **[Risk] Changing `table_group_id` after creation**: Using the immutable args pattern (not ForceNew) means users receive a clear error rather than silent destruction.
  - **Mitigation:** The Update function returns a formatted error explaining the field cannot be changed.

- **[Risk] User-specified id collision**: If a user specifies a `table_group_id` that already exists within the cluster, the `CreateTableGroup` API will return an error.
  - **Mitigation:** No provider-side workaround; the API error is surfaced to the user directly, which is the expected behavior.
