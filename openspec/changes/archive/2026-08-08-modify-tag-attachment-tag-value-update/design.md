## Context

The `tencentcloud_tag_attachment` resource binds a tag key-value pair to a cloud resource (described by a six-segment QCS string). The resource currently has three schema fields — `tag_key`, `tag_value`, and `resource` — all marked `ForceNew: true`. There is no `Update` function; any field change destroys and recreates the attachment.

**Current state:**
- Resource file: `tencentcloud/services/tag/resource_tc_tag_attachment.go`
- Service layer: `tencentcloud/services/tag/service_tencentcloud_tag.go` (`DescribeTagTagAttachmentById`, `DeleteTagTagAttachmentById`)
- Composite ID format: `tagKey + FILED_SP + tagValue + FILED_SP + resource`
- CRUD: Create (`AddResourceTag`), Read (`GetResources` via `DescribeTagTagAttachmentById`), Delete (`DeleteResourceTag`); no Update.

**Problem:** When `tag_value` changes, ForceNew triggers a delete of the old `tag_key:tag_value` binding followed by a create of `tag_key:new_tag_value`. The delete step removes the tag value from the resource, and the subsequent create can fail (e.g., permission policies that require the tag to remain attached, or tag-state races), leaving the resource without the intended tag.

**API behavior analysis:**

| API | TagKey | TagValue | Resource (six-segment) | Purpose |
|-----|--------|----------|------------------------|---------|
| `AddResourceTag` | Yes | Yes | Yes | Create: attach tag key-value to resource |
| `DeleteResourceTag` | Yes | No | Yes | Delete: detach tag key from resource |
| `GetResources` | N/A | N/A | Yes (ResourceList) | Read: query tags attached to resource |
| `UpdateResourceTagValue` | Yes | Yes (new value) | Yes | Update: change tag value for a given key on a single resource, atomically |
| `ModifyResourcesTagValue` | Yes | Yes | Requires parsing six-segment into ServiceType/ResourceIds/Region/Prefix | Update: batch modify tag value across multiple resources |
| `ModifyResourceTags` | Via ReplaceTags array | Via ReplaceTags array | Yes | Generic tag add/delete on a resource |

**Key constraint:** `UpdateResourceTagValue` directly accepts the same three fields as the existing schema (`TagKey`, `TagValue`, `Resource` six-segment), so no six-segment parsing is needed. It changes only the tag value while keeping the tag key attached to the resource — exactly the desired in-place update.

## Goals / Non-Goals

**Goals:**
- Remove `ForceNew: true` from the `tag_value` schema field so Terraform routes tag-value changes through Update.
- Add an `Update` function that calls `UpdateResourceTagValue` to atomically change the tag value on the resource in a single API request.
- Add a service-layer wrapper `UpdateTagAttachmentTagValue` that calls `UpdateResourceTagValue` with retry handling (`tccommon.WriteRetryTimeout`), consistent with existing service functions.
- After a successful update, re-set the composite ID with the new `tag_value` so subsequent reads locate the correct attachment.
- Keep `tag_key` and `resource` as `ForceNew: true` (immutable); only `tag_value` becomes updatable.
- Update the resource documentation and unit tests.

**Non-Goals:**
- Making `tag_key` or `resource` updatable (changing the key or target resource is a different attachment; ForceNew is correct).
- Using `ModifyResourcesTagValue` (requires six-segment parsing) or `ModifyResourceTags` (generic, more complex) — `UpdateResourceTagValue` is the simplest, directly-matching API.
- Changing the composite ID format or the import behavior.

## Decisions

### Decision 1: Use `UpdateResourceTagValue` for tag-value updates

**Rationale:** Three candidate APIs were considered:
1. `UpdateResourceTagValue` — accepts `TagKey`, `TagValue`, `Resource` (six-segment). Directly matches the existing schema fields with no parsing. Changes the tag value in one atomic call while keeping the key attached.
2. `ModifyResourcesTagValue` — accepts `ServiceType`, `ResourceIds` (array), `TagKey`, `TagValue`, `ResourceRegion`, `ResourcePrefix`. Requires parsing the six-segment `resource` string into multiple parts, adding complexity and fragility.
3. `ModifyResourceTags` — generic tag add/delete via `ReplaceTags`/`DeleteTags` arrays. Over-engineered for a single tag-value change.

`UpdateResourceTagValue` is the best fit: minimal, direct, and already available in the vendored SDK.

### Decision 2: Remove `ForceNew` only from `tag_value`, keep it on `tag_key` and `resource`

**Rationale:** Changing `tag_key` or `resource` describes a fundamentally different attachment (different key or different target). These remain `ForceNew`. Only `tag_value` has a valid in-place update path via `UpdateResourceTagValue`.

### Decision 3: Re-set the composite ID after update

**Rationale:** The composite ID is `tagKey + FILED_SP + tagValue + FILED_SP + resource`. When `tag_value` changes, the old ID no longer matches the new state. After a successful `UpdateResourceTagValue` call, the Update function re-sets `d.SetId(newTagKey + FILED_SP + newTagValue + FILED_SP + resource)` so the ID reflects the new tag value. The Read that follows then locates the correct attachment by the new value.

### Decision 4: Service-layer wrapper with retry

**Rationale:** Consistent with the existing service functions (`DescribeTagTagAttachmentById`, `DeleteTagTagAttachmentById`), the new `UpdateTagAttachmentTagValue` wraps the API call with `resource.Retry(tccommon.WriteRetryTimeout, ...)` and `tccommon.RetryError(err)` for error handling, plus `ratelimit.Check`.

### Decision 5: Update reads the new value from state, not re-queried

**Rationale:** After `UpdateResourceTagValue` succeeds and the ID is re-set, the Update function calls `resourceTencentCloudTagAttachmentRead(d, meta)` to refresh state from the cloud. The Read function already parses the ID and queries `GetResources` to verify the attachment, so the state is consistent with the cloud after update.

## Risks / Trade-offs

- **[Risk] ID changes after update break references**: The composite ID changes when `tag_value` changes, but this is expected and consistent with the existing ID design. Terraform handles ID changes within Update correctly.
  - **Mitigation:** Re-set the ID before calling Read so subsequent operations use the new ID.

- **[Risk] `UpdateResourceTagValue` fails if tag key not attached**: If the tag key was manually detached, the update returns `ResourceNotFound.AttachedTagKeyNotFound`.
  - **Mitigation:** The retry layer surfaces the error to the user; no silent state corruption because the ID is only re-set on success.

- **[Risk] Backward compatibility for existing state**: Existing state files have `tag_value` with `ForceNew`. Removing `ForceNew` does not break existing state — Terraform simply routes future changes through Update instead of destroy-recreate.
  - **Mitigation:** No state migration needed; the schema field type and optionality are unchanged.
