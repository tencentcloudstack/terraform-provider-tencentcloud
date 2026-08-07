## Context

The `tencentcloud_tag_attachment` resource (in `tencentcloud/services/tag/resource_tc_tag_attachment.go`) binds a tag (`tag_key` + `tag_value`) to a cloud resource identified by a six-segment ARN (`resource`). Today the resource only implements Create/Read/Delete (no Update), and all three schema fields (`tag_key`, `tag_value`, `resource`) are `ForceNew: true`.

**Current state:**
- Resource file: `tencentcloud/services/tag/resource_tc_tag_attachment.go`
- Service layer: `tencentcloud/services/tag/service_tencentcloud_tag.go`
- SDK: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag/v20180813`
- Composite ID format: `tagKey + FILED_SP + tagValue + FILED_SP + resource` (three segments)

**Problem with current `ForceNew` on `tag_value`:**
When a user changes `tag_value` on an existing attachment, Terraform's default behavior is delete-then-create:
1. `DeleteResourceTag` removes the old `tag_key` binding from the resource.
2. `AddResourceTag` attempts to add `tag_key` + new `tag_value`.

Step 2 fails because the resource already has the tag key associated (the old binding was not fully cleared, or a race with eventual consistency), leaving the resource in a broken state.

**API behavior analysis:**

| API | TagKey | TagValue | Resource | Purpose |
|-----|--------|----------|----------|---------|
| `AddResourceTag` | Yes | Yes | Yes | Create: bind a tag to a resource |
| `DeleteResourceTag` | Yes | No | Yes | Delete: remove a tag key binding from a resource |
| `GetResources` | N/A | N/A | Yes (ResourceList) | Read: query tags bound to a resource |
| `UpdateResourceTagValue` | Yes | Yes (new value) | Yes | Update: modify the tag value of an already-associated tag (tag key unchanged) |

The `UpdateResourceTagValue` API is purpose-built for this case: it modifies the tag value of an already-associated tag in a single atomic call, with the tag key and resource unchanged.

## Goals / Non-Goals

**Goals:**
- Make `tag_value` updatable in place on `tencentcloud_tag_attachment` by removing `ForceNew: true` and adding an `Update` function.
- The Update function SHALL call `UpdateResourceTagValue` (request: `TagKey`, `TagValue`, `Resource`) when `tag_value` changes, performing a single-step modification.
- After a successful update, rebuild the composite ID with the new `tag_value` (since `tag_value` is part of the composite ID).
- Keep `tag_key` and `resource` as `ForceNew: true` (immutable), consistent with `UpdateResourceTagValue` semantics.
- Maintain full backward compatibility — existing configurations continue to work; only the behavior on `tag_value` change improves.

**Non-Goals:**
- Making `tag_key` or `resource` updatable (the API does not support changing the tag key or resource of an existing binding in place).
- Adding Update support for the `tencentcloud_tag` resource (out of scope).
- Changing the composite ID format or separator.

## Decisions

### Decision 1: Use `UpdateResourceTagValue` for tag_value changes

**Rationale:** `UpdateResourceTagValue` modifies the tag value of an already-associated tag in a single atomic call (tag key and resource unchanged). This directly solves the delete-then-create race that causes the second `AddResourceTag` request to fail. Among the available APIs (`ModifyResourcesTagValue`, `ModifyResourceTags`, `UpdateResourceTagValue`), `UpdateResourceTagValue` is the best fit because it accepts the full resource six-segment ARN directly (like `AddResourceTag` and `DeleteResourceTag` already used by this resource), whereas `ModifyResourcesTagValue` requires decomposing the ARN into `ServiceType`/`ResourceIds`/`ResourceRegion`/`ResourcePrefix`, and `ModifyResourceTags` operates on replace/delete tag sets (more complex and not needed for a single tag value change).

**Alternatives considered:**
- `ModifyResourcesTagValue`: Requires decomposing the six-segment ARN into separate fields (`ServiceType`, `ResourceIds`, `ResourceRegion`, `ResourcePrefix`). This is more complex and error-prone, and inconsistent with the resource's existing use of `AddResourceTag`/`DeleteResourceTag` which take the full ARN.
- `ModifyResourceTags`: Operates on `ReplaceTags`/`DeleteTags` arrays. Heavier than needed for a single tag value change, and its "replace" semantics could be confusing for a single-value update.

### Decision 2: Keep `tag_key` and `resource` as `ForceNew: true`

**Rationale:** `UpdateResourceTagValue` only modifies the tag value; the tag key and resource cannot change. If `tag_key` or `resource` changes, the resource must be recreated. Keeping `ForceNew: true` on these fields preserves existing behavior and matches API capabilities.

### Decision 3: Rebuild composite ID after update

**Rationale:** The composite ID is `tagKey + FILED_SP + tagValue + FILED_SP + resource`. Since `tag_value` is part of the ID and it changes during update, the Update function MUST rebuild the ID after a successful `UpdateResourceTagValue` call so subsequent Read operations use the correct ID. The Update function detects the change via `d.HasChange("tag_value")` and reads the new value via `d.GetOk("tag_value")`, then sets the new composite ID with `d.SetId(...)`. The `UpdateResourceTagValue` API only requires the new `TagValue` (the old value is not needed), so the old `tag_value` is not read.

### Decision 4: Add Update via the resource schema `Update` callback

**Rationale:** The resource currently has no `Update` callback. Adding `Update: resourceTencentCloudTagAttachmentUpdate` to the resource schema and removing `ForceNew: true` from `tag_value` makes Terraform call Update when only `tag_value` changes. When `tag_key` or `resource` changes, their `ForceNew: true` still triggers recreate as before.

## Risks / Trade-offs

- **[Risk] State drift if `UpdateResourceTagValue` fails mid-flight:** If the API call fails and retry is exhausted, the state may not match the cloud. 
  - **Mitigation:** Use `tccommon.WriteRetryTimeout` with `tccommon.RetryError` wrapping (consistent with Create/Delete), so transient failures are retried. On persistent failure, Terraform marks the update as failed and the user can re-run.

- **[Risk] Composite ID mismatch after update:** If the ID is not rebuilt, subsequent Read would use the old `tag_value` and fail to find the binding.
  - **Mitigation:** The Update function rebuilds the composite ID with the new `tag_value` after a successful API call, before returning. Read is then called automatically by Terraform to refresh state.

- **[Risk] Concurrent modifications to the same resource's tags:** If another process modifies the resource's tags between Terraform's plan and apply, the update may fail or produce unexpected state.
  - **Mitigation:** This is an inherent risk of tag management and exists with the current Create/Delete flow as well. No additional mitigation beyond retry is needed.

- **[Trade-off] `tag_value` change no longer recreates the resource:** This is the intended behavior change. Users relying on recreate-on-tag_value-change (unlikely) would see in-place update instead. This is a strictly better behavior since it avoids the delete-then-create failure mode.
