## Context

The `tencentcloud_tag_attachment` resource binds a tag key/value pair to a cloud resource (identified by a six-segment resource string). The resource currently has only Create/Read/Delete operations — no Update. Both `tag_key`, `tag_value`, and `resource` are all `ForceNew: true`, meaning any change triggers destroy-and-recreate.

The composite id is `tagKey + FILED_SP + tagValue + FILED_SP + resource`.

**Problem:** Some cloud resources use the tag KV pair to control access permissions. When `tag_value` changes, the current ForceNew behavior makes Terraform first delete the attachment (removing the tag key entirely) then re-add it with the new value. After the delete step, the caller loses the authorizing tag and the subsequent add request fails.

**Cloud API support:** The Tag service (`tag/v20180813`) provides `UpdateResourceTagValue`, documented as "本接口用于修改资源已关联的标签值（标签键不变）" (modify the tag value associated with a resource, keeping the tag key unchanged). Its request takes `TagKey`, `TagValue` (the new value), and `Resource` — exactly the fields we already have. The method is available in the upgraded Tag SDK (`v1.3.110`), so the SDK module was upgraded from `v1.0.860` to `v1.3.110` (and re-vendored) to obtain `Client.UpdateResourceTagValue`.

**Current state:**
- Resource file: `tencentcloud/services/tag/resource_tc_tag_attachment.go`
- Service layer: `tencentcloud/services/tag/service_tencentcloud_tag.go` (`DescribeTagTagAttachmentById`, `DeleteTagTagAttachmentById`)
- SDK: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag/v20180813` — upgraded to `v1.3.110` to obtain `UpdateResourceTagValue`

**API behavior analysis:**

| API | TagKey in Request | TagValue in Request | Resource in Request | Purpose |
|-----|-------------------|---------------------|---------------------|---------|
| `AddResourceTag` | Yes | Yes (value to bind) | Yes | Create attachment |
| `GetResources` | N/A | N/A | Yes (ResourceList) | Read/verify attachment |
| `UpdateResourceTagValue` | Yes | Yes (new value) | Yes | Update tag value in place (key unchanged) |
| `DeleteResourceTag` | Yes | No | Yes | Delete attachment |

## Goals / Non-Goals

**Goals:**
- Remove `ForceNew: true` from the `tag_value` schema field so changing it does not destroy the resource
- Add an `Update` function to `tencentcloud_tag_attachment` that calls `UpdateResourceTagValue` when `tag_value` changes
- After a successful update, refresh the composite id to reflect the new `tag_value`
- Keep `tag_key` and `resource` immutable (`ForceNew: true`) since `UpdateResourceTagValue` identifies the attachment by tag key + resource and does not support changing them
- Maintain full backward compatibility — existing configurations and state continue to work unchanged

**Non-Goals:**
- Making `tag_key` or `resource` updatable (the cloud API does not support changing the tag key of an existing attachment)
- Adding new schema fields
- Changing the create, read, or delete behavior

## Decisions

### Decision 1: Update `tag_value` in place via `UpdateResourceTagValue`

**Rationale:** The cloud API `UpdateResourceTagValue` directly modifies the tag value of an existing resource-tag association without removing the tag key. This avoids the delete-then-create cycle that drops the authorizing tag and causes the second request to fail. The request parameters (`TagKey`, `TagValue`, `Resource`) are all already available in the resource schema.

### Decision 2: Keep `tag_key` and `resource` as `ForceNew`

**Rationale:** `UpdateResourceTagValue` only changes the tag value; it does not support changing the tag key or the resource. Changing the tag key or resource of an attachment is semantically a different attachment, so destroy-and-recreate is the correct behavior. Keeping these as `ForceNew` is consistent with the API.

### Decision 3: Refresh composite id after update

**Rationale:** The composite id embeds `tag_value` (`tagKey#tagValue#resource`). After a successful `UpdateResourceTagValue` call, the stored id must be updated to `tagKey#newTagValue#resource` so subsequent Read/Delete operations target the correct attachment. The Update function calls `d.SetId()` with the new composite id before delegating to Read.

### Decision 4: Wrap the update API call in `resource.Retry` with `WriteRetryTimeout`

**Rationale:** Consistent with the Create function and the project's conventions for write operations. The retry block only calls the API; setting the id and other success-path operations happen outside the retry block, after the error-handling path.

### Decision 5: Use `d.HasChange("tag_value")` to gate the update call

**Rationale:** The Update function should only call `UpdateResourceTagValue` when `tag_value` actually changed. Since `tag_key` and `resource` are ForceNew, they cannot trigger Update. This mirrors the `mutableArgs`/`d.HasChange` pattern used by other resources (e.g. `tencentcloud_igtm_strategy`).

## Risks / Trade-offs

- **[Risk] State drift if `UpdateResourceTagValue` partially fails**: The API call either succeeds or fails atomically (no partial state). On failure, the retry mechanism retries; on exhaustion the error propagates and the state is not updated, so the next plan will retry.
  - **Mitigation:** The retry block wraps only the API call; `d.SetId()` is called only after the retry succeeds.

- **[Risk] Existing state with old id format**: Existing state already uses the `tagKey#tagValue#resource` format, which is unchanged by this work (only the `tagValue` segment is refreshed after an update). No migration of existing state is required.
  - **Mitigation:** No state migration needed; the id format is identical.

- **[Risk] Update function called with no actual change to `tag_value`**: If Terraform invokes Update for a no-op change, the function would still call the API unnecessarily.
  - **Mitigation:** Gate the API call behind `d.HasChange("tag_value")`; if no change, skip directly to Read.
