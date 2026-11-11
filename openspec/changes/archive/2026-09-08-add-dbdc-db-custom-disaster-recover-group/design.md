## Context

The `dbdc` (DB Custom) service manages dedicated database clusters and nodes. A
**placement group (置放群组)** controls how physical hosts are distributed to
reduce correlated failures. The cloud API exposes four placement-group
operations:

| Operation | API | Sync/Async |
|---|---|---|
| Create | `CreateDBCustomDisasterRecoverGroup` | Sync response with `Status` field; group becomes `Available` after creation completes |
| Read | `DescribeDBCustomDisasterRecoverGroups` | Sync |
| Update | `ModifyDBCustomDisasterRecoverGroupAttribute` | Sync (only `name` + `affinity` mutable) |
| Delete | `DeleteDBCustomDisasterRecoverGroups` | Async — returns `TaskId` |

The Terraform provider already has `dbdc` resources (`dbdc_db_custom_cluster`,
`dbdc_db_custom_node`, `dbdc_node_to_db_custom_cluster_attachment`) and a shared
async helper `waitDBCustomTaskSucceeded` in the `dbdc` package, plus a
`DescribeDBCustomTaskStatusById` service method. This change adds the missing
placement-group resource following the established `RESOURCE_KIND_GENERAL`
conventions (reference: `tencentcloud_igtm_strategy`).

## Goals / Non-Goals

**Goals:**
- Provide full CRUD lifecycle management of a DB Custom placement group via
  Terraform.
- Handle the async nature of delete (poll `TaskId` via
  `DescribeDBCustomTaskStatus`) and the create finalization (poll
  `DescribeDBCustomDisasterRecoverGroups` until `Status == "Available"`).
- Support import via `DisasterRecoverGroupId`.
- Honor provider retry/ratelimit and nil-safety conventions.

**Non-Goals:**
- Managing nodes inside a placement group (that is a separate concern; existing
  node APIs cover membership).
- Bulk placement-group management (one resource = one group).
- Exposing the Describe `Filters`/`Tags` query parameters as user-facing
  schema arguments — they are internal read lookup helpers.

## Decisions

### 1. Resource ID — single `DisasterRecoverGroupId` (not composite)

The create response returns a single `DisasterRecoverGroupId`. Unlike
`igtm_strategy` (which needs `instanceId#strategyId`), this resource has a
single natural key, so `d.SetId(disasterRecoverGroupId)` is used directly.
Import uses `ImportStatePassthrough`.

**Alternative considered**: composite ID — rejected because there is no second
key to encode.

### 2. Create finalization — poll Describe until `Available`

`CreateDBCustomDisasterRecoverGroup` returns the group immediately with
`Status == "Creating"`. To avoid races where `Read` runs before the group is
usable, after a successful create we poll
`DescribeDBCustomDisasterRecoverGroups` (by id) inside `resource.Retry` until
`Status == "Available"`. On `CreateFailed` we return a non-retryable error.
This mirrors the async-polling discipline used elsewhere in the provider and the
project rule for异步接口.

**Alternative considered**: rely solely on retry in Read — rejected because it
delays failure surfacing and can produce confusing empty-state logs.

### 3. Delete finalization — reuse `waitDBCustomTaskSucceeded`

`DeleteDBCustomDisasterRecoverGroups` returns a `TaskId`. We reuse the existing
`waitDBCustomTaskSucceeded(ctx, &service, taskId, d.Timeout(schema.TimeoutDelete))`
helper (already in `resource_tc_dbdc_db_custom_cluster.go`) which polls
`DescribeDBCustomTaskStatus` until `Succeeded`/`Failed`. No new helper needed.

### 4. Mutable vs immutable fields

`ModifyDBCustomDisasterRecoverGroupAttribute` only accepts `name` and
`affinity`. Therefore:
- `name` → Optional, mutable (not ForceNew).
- `affinity` → Optional, mutable (not ForceNew).
- `type`, `strategy` → Optional, `ForceNew` (no modify support).
- `tags` → Optional, `ForceNew` (Create accepts `Tags`; Modify has no tag API;
  tags are not updatable in-place). Per convention, modeled as `TypeList` of
  `{key,value}` blocks matching the SDK `Tag` struct (the Create API takes a
  list of `Tag` objects, not a map).
- `client_token` → Optional, `ForceNew` (idempotency token, create-only).

For the `immutableArgs` discipline: in `Update`, we build `immutableArgs` from
the top-level fields that are NOT supported by Modify (`type`, `strategy`,
`tags`, `client_token`) and reject changes to them with an error. Only `name`
and `affinity` trigger the Modify call.

### 5. `tags` schema shape — `TypeList` of `{key, value}` blocks

The cloud API `Tags` field is `[]*Tag` where `Tag{Key, Value}`. The user query
maps `request.Tags.Key → key` and `request.Tags.Value → value` (both required).
Following the igtm_strategy nested-block convention, we model `tags` as
`TypeList` with `Elem: &schema.Resource{Schema: {key, value}}` rather than a
`TypeMap`, because the SDK structure is an explicit object list.

### 6. Read lookup — new `DescribeDBCustomDisasterRecoverGroupById` service helper

Add a service-layer helper that calls
`DescribeDBCustomDisasterRecoverGroups` with `DisasterRecoverGroupIds=[id]`,
`Limit=100` (documented max), wrapped in `resource.Retry(ReadRetryTimeout)` +
`ratelimit.Check`, and returns the single matching
`*dbdcv20201029.DisasterRecoverGroup` (nil if not found). Pagination is
unnecessary for a single-id lookup but the loop is kept for consistency.

### 7. Computed fields

Read populates computed fields from `DisasterRecoverGroup`: `status`,
`node_quota_total`, `current_num`, `created_time`, `node_ids`, plus echoes back
`name`, `type`, `affinity`, `strategy`, `tags`, and the id
`disaster_recover_group_id`.

## Risks / Trade-offs

- **[Risk] Create polling may loop on non-terminal statuses** → We treat
  `CreateFailed` as a terminal non-retryable failure and `Available` as
  success; any other status is retryable. Timeout is bounded by
  `d.Timeout(schema.TimeoutCreate)`.
- **[Risk] Tags not updatable in-place** → Modeled as `ForceNew`; changing tags
  recreates the group. This is consistent with the absence of a Modify-tags API
  and is clearly documented.
- **[Risk] `Affinity` type mismatch between Create (int64) and Read
  (uint64)** → Read sets the value via the SDK-provided pointer; we use
  `helper.Int64`/appropriate conversion at the schema boundary so the int
  schema round-trips correctly.
- **[Risk] Delete of a non-empty group fails server-side** → The API requires
  nodes to be removed first; this is a server-side validation surfaced as an
  error to the user, not handled in provider code (consistent with API
  semantics).
