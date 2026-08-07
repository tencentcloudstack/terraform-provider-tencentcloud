## Context

The `tencentcloud_tag_attachment` resource binds a single tag key/value pair to a cloud resource described by its six-segment QCS string. Today all three schema fields (`tag_key`, `tag_value`, `resource`) are `ForceNew: true`, so any change destroys and recreates the attachment.

The destroy+recreate sequence is problematic for tag value changes. When a resource already binds `运营产品:A` and the user changes the value to `运营产品:B`, Terraform first deletes the `运营产品:A` binding and then attempts to add `运营产品:B`. The second request can fail (e.g. when the resource still carries a value for that key during propagation), and the two-step sequence is non-atomic.

**Current state:**
- Resource file: `tencentcloud/services/tag/resource_tc_tag_attachment.go`
- Service layer: `tencentcloud/services/tag/service_tencentcloud_tag.go`
- Composite id: `tagKey + FILED_SP + tagValue + FILED_SP + resource`
- SDK: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag/v20180813`

**Available tag service APIs for modifying a bound tag value (verified in vendored SDK):**

| API | Request fields | Notes |
|-----|-----------------|-------|
| `UpdateResourceTagValue` | `TagKey`, `TagValue` (new), `Resource` (six-segment) | Matches the existing schema fields exactly. |
| `ModifyResourcesTagValue` | `ServiceType`, `ResourceIds[]`, `TagKey`, `TagValue`, `ResourceRegion`, `ResourcePrefix` | Requires decomposing the six-segment resource string into separate parts. |
| `ModifyResourceTags` | `Resource`, `ReplaceTags[]`, `DeleteTags[]` | Generic add/replace/delete; accepts the full six-segment `Resource`. |

## Goals / Non-Goals

**Goals:**
- Make `tag_value` updatable in place (remove `ForceNew`) on `tencentcloud_tag_attachment`.
- Implement an `Update` function that modifies the bound tag value in a single API call when only `tag_value` changes.
- Keep `tag_key` and `resource` immutable (`ForceNew: true`).
- Refresh the composite id after a successful tag value update so `d.Id()` reflects the new value segment.
- Add a service-layer helper that wraps the chosen API with `tccommon.WriteRetryTimeout` retry and `ratelimit.Check`.
- Add an acceptance test that covers the in-place update of `tag_value`.
- Update the `.md` example to document the in-place update behavior.

**Non-Goals:**
- Updating `tag_key` or `resource` in place (they remain `ForceNew`).
- Changing the composite id format.
- Adding support for batch/multi-resource tag operations.
- Migrating existing state (the id format is unchanged; only the value segment changes).

## Decisions

### Decision 1: Use `UpdateResourceTagValue` as the update API

**Rationale:** The `UpdateResourceTagValue` request takes exactly `TagKey`, `TagValue` (the new value), and `Resource` (the six-segment QCS string). These map 1:1 to the existing schema fields (`tag_key`, `tag_value`, `resource`), so no decomposition of the resource six-segment is needed. `ModifyResourcesTagValue` would require splitting the QCS into `ServiceType`/`ResourcePrefix`/`ResourceRegion`/`ResourceIds`, adding fragile parsing logic. `ModifyResourceTags` is more generic (replace/delete tag sets) and is already used elsewhere for multi-key rewrites; using it for a single value change is overkill and adds unrelated complexity. `UpdateResourceTagValue` is the simplest, most direct fit.

**Alternatives considered:**
- `ModifyResourcesTagValue`: rejected because it requires decomposing the six-segment resource string into `ServiceType`, `ResourceRegion`, `ResourcePrefix`, and `ResourceIds`, which is fragile and duplicates parsing the QCS.
- `ModifyResourceTags` with a single `ReplaceTags` entry: rejected because it is a generic add/replace/delete interface already used for multi-key rewrites; using it for a single tag value change adds unnecessary generality and would not clearly express "update this one value".

### Decision 2: Remove `ForceNew` only from `tag_value`; keep `tag_key` and `resource` immutable

**Rationale:** The requirement is specifically about updating the tag value. `tag_key` and `resource` identify which binding is being modified; changing them is a different attachment and must remain `ForceNew`. This keeps the change minimal and backward compatible.

### Decision 3: Refresh the composite id after update

**Rationale:** The composite id is `tagKey#tagValue#resource`. After a successful `UpdateResourceTagValue`, the `tagValue` segment changes, so the Update function recomputes and sets `d.SetId(tagKey + FILED_SP + newTagValue + FILED_SP + resource)` before calling Read. This keeps state consistent for subsequent plan/apply cycles.

### Decision 4: Add a dedicated service-layer helper `UpdateTagTagValue`

**Rationale:** Consistent with the existing service layer (`DescribeTagTagAttachmentById`, `DeleteTagTagAttachmentById`). The helper builds the `UpdateResourceTagValueRequest`, applies `ratelimit.Check`, and wraps the call in `resource.Retry(tccommon.WriteRetryTimeout, ...)` with `tccommon.RetryError(err)` on failure — matching the retry pattern used across the tag service. The resource six-segment is passed as the `resourceName` parameter (matching the `ModifyTags` naming convention) to avoid shadowing the imported `helper/resource` package within the retry closure.

### Decision 5: Read tag_key and resource from the existing state id (not a separate Describe)

**Rationale:** The Update function parses the current id (`tagKey#oldTagValue#resource`) to validate the id and obtain the `tag_key` (first segment) and `resource` (third segment). It reads the new tag value from `d.Get("tag_value")`. The old tag value segment is not needed because `UpdateResourceTagValue` only takes the key, the new value, and the resource. This avoids an extra read API call before updating and matches how the existing Read/Delete functions already split the id.

## Risks / Trade-offs

- **[Risk] Idempotency when the new tag value equals the old value**: Terraform only calls Update when the diff is non-empty, so a no-op change does not trigger an API call. No special handling needed.
  - **Mitigation:** None required; the SDK/plan layer guards against no-op updates.

- **[Risk] The `UpdateResourceTagValue` API could fail if the resource does not currently bind `tagKey`**: This is expected API behavior and surfaces as an error to the user.
  - **Mitigation:** The retry wrapper handles transient errors; persistent failures are returned to the user with a clear log line.

- **[Risk] State drift between Delete and the old ForceNew path**: Removing `ForceNew` from `tag_value` changes the plan behavior. Existing state is unaffected because the id format and schema types are unchanged.
  - **Mitigation:** Composite id format is preserved; only the value segment is refreshed after update.
