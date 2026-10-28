## Context

The `dlc` (Data Lake Compute) service already ships two "operation"-style
resources in this provider:

- `tencentcloud_dlc_attach_work_group_policy_operation` — calls
  `AttachWorkGroupPolicy` but its `Read`/`Delete` are no-ops.
- `tencentcloud_dlc_detach_work_group_policy_operation` — calls
  `DetachWorkGroupPolicy` but its `Read`/`Delete` are no-ops.

Neither tracks the binding lifecycle, so drift is invisible and the resources
cannot be refreshed or imported. The cloud APIs available in the vendored SDK
(`tencentcloud/dlc/v20210125`) are:

- `AttachWorkGroupPolicy(WorkGroupId *int64, PolicySet []*Policy)` →
  `Response.PolicySet []*Policy` (sync; returns the granted policies incl.
  `PolicyId`).
- `DescribeWorkGroupInfo(WorkGroupId, Type, ...)` →
  `Response.WorkGroupInfo *WorkGroupDetailInfo`. With `Type="DataAuth"` the
  `WorkGroupInfo.DataPolicyInfo.PolicySet` lists the bound data policies (each
  carrying a `PolicyId`).
- `DetachWorkGroupPolicy(WorkGroupId *int64, PolicySet []*Policy,
  PolicyIds []*string)` → `Response{RequestId}` (sync).

These three APIs are synchronous (no `TaskId`/async polling needed), so the
attachment is a straightforward Create/Read/Delete resource.

The existing sibling `tencentcloud_dlc_add_users_to_work_group_attachment`
already establishes the attachment pattern for this package (composite ID with
`tccommon.FILED_SP`, `ImportStatePassthrough`, no-op-less Read/Delete via the
`DlcService`). This change follows that pattern.

## Goals / Non-Goals

**Goals:**
- Provide a true `RESOURCE_KIND_ATTACHMENT` resource
  `tencentcloud_dlc_attach_work_group_policy_attachment` that binds one or more
  policies to a work group and unbinds them on destroy.
- Make `Read` actually verify membership via `DescribeWorkGroupInfo` so drift is
  detected (the binding disappears → `d.SetId("")`).
- Support import via a composite ID.
- Keep the schema consistent with the existing `policy_set` block used by the
  operation resources (same field names/types) so users can migrate trivially.

**Non-Goals:**
- Not modifying or deprecating the existing
  `tencentcloud_dlc_attach_work_group_policy_operation` /
  `tencentcloud_dlc_detach_work_group_policy_operation` resources (kept for
  backward compatibility).
- Not supporting in-place update (attachment resources are `ForceNew`-only; any
  change recreates).
- Not managing engine-policy or row-filter-policy bindings (only the
  `DataAuth`/data-policy dimension is covered, matching the primary use case of
  `AttachWorkGroupPolicy`).
- Not introducing async task polling — the three APIs are synchronous.

## Decisions

### Decision 1: Resource ID = `WorkGroupId#PolicyId` joined by `tccommon.FILED_SP`

Each attachment resource instance binds **one** policy to a work group. The ID
is `WorkGroupId + "#" + PolicyId` (e.g. `23184#policy-xxxx`), where `PolicyId`
is the deterministic `Policy.PolicyId` string returned by
`AttachWorkGroupPolicy`/`DescribeWorkGroupInfo`.

**Why one policy per resource (not the whole `PolicySet`):** The cloud API
accepts a list, but a Terraform attachment resource is most robust when it
represents a single bind relationship — this keeps Read/Delete deterministic
(locate exactly one policy by `PolicyId`) and matches the
`tencentcloud_dlc_add_users_to_work_group_attachment` granularity philosophy
(one logical binding per resource). The schema still accepts a `policy_set`
list (to stay compatible with the API shape and the existing operation
resource), but `MaxItems: 1` is enforced so one resource == one policy binding.
The `PolicyId` for the ID is taken from the first (only) element of the
response `PolicySet`.

**Alternative considered:** allow multiple policies per resource and join all
`PolicyId`s in the ID (like `add_users_to_work_group_attachment` joins
`userIds`). Rejected because partial unbind on a multi-policy resource would
require diffing the list on Read, complicating drift detection; one-policy-per-
resource is simpler and composes via `count`/`for_each`.

**Note on `PolicyId` availability:** `Policy.PolicyId` is the
"user-and-workgroup deterministic PolicyId" per the SDK comment and is returned
by both `AttachWorkGroupPolicy` (response) and `DescribeWorkGroupInfo`
(`DataPolicyInfo.PolicySet`). It is the stable identifier used for the
composite ID and for the delete request's `PolicyIds`.

### Decision 2: Schema mirrors the existing `policy_set` block

To keep user configs portable from the operation resource, the nested
`policy_set` block reuses the exact same field set and descriptions as
`resource_tc_dlc_attach_work_group_policy_operation.go` (`database`, `catalog`,
`table`, `operation`, `policy_type`, `function`, `view`, `column`,
`data_engine`, `re_auth`, `source`, `mode`, `operator`, `create_time`,
`source_id`, `source_name`, `id`), with `MaxItems: 1` added. `work_group_id`
is a top-level `ForceNew` `TypeInt` (matching the operation resource).

### Decision 3: Read uses `DescribeWorkGroupInfo` with `Type=DataAuth`

`DescribeWorkGroupInfo` is the only read API that returns the bound policy set.
`Type` is hard-coded to `"DataAuth"` (data permissions) and `Limit`/`Offset`
are set to the documented max (`100`) to page through all policies when
locating the `PolicyId`. If the `PolicyId` is not found in
`DataPolicyInfo.PolicySet`, the binding is gone →
`log.Printf("[CRUD] ..."); d.SetId("")`.

### Decision 4: Delete calls `DetachWorkGroupPolicy` with `PolicyIds=[policyId]`

`DetachWorkGroupPolicy` accepts both `PolicySet` and `PolicyIds`. Using
`PolicyIds=[policyId]` (the deterministic ID we stored) is the most reliable
unbind path and avoids having to reconstruct the full `Policy` struct on
delete. This aligns with the `policy_ids` schema field listed in the API
mapping.

### Decision 5: gomonkey-based unit tests (no terraform test suite)

Per the project rules for newly-added resources, the `_test.go` file uses
gomonkey to mock the DLC client methods (`AttachWorkGroupPolicy`,
`DescribeWorkGroupInfo`, `DetachWorkGroupPolicy`) and tests the business logic
of Create/Read/Delete directly, runnable via
`go test -gcflags=all=-l` (no real cloud credentials).

## Risks / Trade-offs

- **[Risk] `PolicyId` may be empty for some policy types** → The Create
  callback MUST guard that the response `PolicySet[0].PolicyId` is non-empty
  and return `NonRetryableError` if so (per project rule #9), preventing a
  broken empty-segment ID from being written to state.
- **[Risk] `DescribeWorkGroupInfo` pagination** → The API caps `Limit` at 100.
  The Read helper paginates with `Offset` increments until the target
  `PolicyId` is found or the returned page is short, so work groups with >100
  policies are handled.
- **[Risk] Drift: policy unbound out-of-band** → Read sets `d.SetId("")` when
  the `PolicyId` is absent, letting Terraform reconcile on next apply. This is
  the desired attachment semantics.
- **[Trade-off] `MaxItems:1` on `policy_set`** → One resource per binding. Users
  wanting multiple policies use multiple resources. Accepted for determinism.
- **[Trade-off] Reuses existing `DlcService` style** → A small
  `DescribeDlcWorkGroupPolicyAttachmentById` helper is added to the service
  layer (wrapping `DescribeWorkGroupInfo` + pagination + policy lookup), plus a
  `DeleteDlcAttachWorkGroupPolicyAttachmentById` helper (already exists as
  `DeleteDlcAttachWorkGroupPolicyAttachmentById` using `DetachWorkGroupPolicy`
  — but that one takes `[]*dlc.Policy`; we'll add/adjust to also support
  `PolicyIds`). Kept consistent with existing service-layer conventions.
