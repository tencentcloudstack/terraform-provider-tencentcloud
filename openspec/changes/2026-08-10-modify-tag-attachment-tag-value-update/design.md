## Context

The `tencentcloud_tag_attachment` resource manages the association between a cloud resource, a tag key, and a tag value. Its composite id is `tagKey + FILED_SP + tagValue + FILED_SP + resource`.

**Current state:**
- Resource file: `tencentcloud/services/tag/resource_tc_tag_attachment.go`
- Service layer: `tencentcloud/services/tag/service_tencentcloud_tag.go`
- The resource currently only implements `Create`, `Read`, `Delete` (no `Update`).
- All three schema fields (`tag_key`, `tag_value`, `resource`) are `ForceNew: true`.
- Create uses `AddResourceTag` API; Delete uses `DeleteResourceTag` API (deletes by `tag_key` + `resource`, without `tag_value`).

**Problem with the current `ForceNew` on `tag_value`:**
When a user changes `tag_value`, Terraform treats it as a replacement: it runs `Delete` (which calls `DeleteResourceTag`, removing the `tag_key` association entirely) and then `Create` (which calls `AddResourceTag` with the new value). These are **two separate API calls**. If the second call fails, the resource is left with **no tag attached**. For self-developed cloud onboarding systems that gate resource access on the presence of a specific tag KV (e.g. `运营产品:A`), losing the tag means losing access to the resource entirely.

**Cloud API behavior analysis (tag v20180813, vendored SDK):**

| API | Purpose | Key request fields |
|-----|---------|--------------------|
| `AddResourceTag` | Attach a `{tag_key, tag_value}` to a `resource` (single KV) | `TagKey`, `TagValue`, `Resource` |
| `DeleteResourceTag` | Detach a `tag_key` association from a `resource` | `TagKey`, `Resource` |
| `ModifyResourceTags` | Atomically replace and/or delete tags on a `resource` | `Resource`, `ReplaceTags []*Tag{TagKey,TagValue}`, `DeleteTags []*TagKeyObject{TagKey}` |

The `ModifyResourceTags` API is the key enabler: it accepts `ReplaceTags` and `DeleteTags` in a **single request**. Per the API contract, when a resource already has a `tag_key` associated, placing `{tag_key, new_value}` in `ReplaceTags` changes the existing value to the new value in place (the API states: "若已关联，则将该资源关联的键对应的标签值修改为输入值"). Additionally, `ReplaceTags` and `DeleteTags` must not contain the same tag key, so for a pure tag-value change the update uses only `ReplaceTags = [{tag_key, new_value}]` and leaves `DeleteTags` empty. This avoids the intermediate "no tag" state that `ForceNew` (delete `DeleteResourceTag` + add `AddResourceTag`) produces, which breaks tag-KV-gated access control when the second call fails.

The service layer already wraps `ModifyResourceTags` in `TagService.ModifyTags(ctx, resourceName, replaceTags map[string]string, deleteKeys []string)` (the `deleteKeys` slice is only turned into `DeleteTags` when non-empty, so passing `nil` is safe).

## Goals / Non-Goals

**Goals:**
- Remove `ForceNew: true` from `tag_value` so changing it performs an in-place update.
- Add an `Update` function (`resourceTencentCloudTagAttachmentUpdate`) that updates `tag_value` in a single atomic `ModifyResourceTags` request.
- Reuse the existing `TagService.ModifyTags` service method — no new SDK calls or service methods.
- Refresh the composite id (`tagKey + FILED_SP + newTagValue + FILED_SP + resource`) after a successful update so Read resolves the correct attachment.
- Maintain full backward compatibility: field type/name unchanged, existing configs and state keep working.

**Non-Goals:**
- Making `tag_key` or `resource` updatable — these identify the attachment target and remain `ForceNew`.
- Changing the Create (`AddResourceTag`) or Delete (`DeleteResourceTag`) flows.
- Bumping the SDK version — `ModifyResourceTags` is already vendored and wrapped.

## Decisions

### Decision 1: Update via `ModifyResourceTags` (atomic), not delete-then-add

**Rationale:** `ModifyResourceTags` accepts `ReplaceTags` and `DeleteTags` in a single API request. Per the API contract, when a resource already has a `tag_key` associated, placing `{tag_key, new_value}` in `ReplaceTags` changes the existing value to the new value in place. Because `ReplaceTags` and `DeleteTags` must not contain the same tag key, for a pure tag-value change the update uses only `ReplaceTags = [{tag_key, new_value}]` and leaves `DeleteTags` empty. This updates the tag value atomically in one request, avoiding the intermediate "no tag" state that `ForceNew` (delete `DeleteResourceTag` + add `AddResourceTag`) produces, which breaks tag-KV-gated access control when the second call fails.

**Alternatives considered:**
- *Put old value's key in `DeleteTags` and new value in `ReplaceTags` of the same request:* Rejected — the API forbids `ReplaceTags` and `DeleteTags` from sharing the same tag key, so this is not allowed for the same `tag_key`.
- *Keep ForceNew but make delete+add transactional:* Not feasible — `DeleteResourceTag` and `AddResourceTag` are separate APIs with no transactional guarantee.
- *Add `UpdateResourceTag` API:* No such API exists in the tag v20180813 SDK; `ModifyResourceTags` is the supported way to change a resource's tag value.

### Decision 2: Reuse `TagService.ModifyTags`, no new service method

**Rationale:** `TagService.ModifyTags(ctx, resourceName, replaceTags, deleteKeys)` already builds the `ModifyResourceTagsRequest` (`ReplaceTags` from the map, `DeleteTags` from the keys slice) and wraps it with `WriteRetryTimeout` retry via `ratelimit.Check` and `tccommon.RetryError`. The Update function calls it with `replaceTags = {tag_key: new_value}` and `deleteKeys = nil` (leaving `DeleteTags` empty, because `ReplaceTags` and `DeleteTags` must not contain the same tag key per the API contract). Adding a new service method would duplicate this logic.

### Decision 3: `tag_key` and `resource` stay `ForceNew`

**Rationale:** `tag_key` and `resource` identify *which* resource and *which* tag key the attachment refers to — they are part of the composite id and the attachment identity. Changing either means a different attachment, so it must force replacement. Only `tag_value` is the mutable attribute of an existing attachment.

### Decision 4: Refresh the composite id after update

**Rationale:** The id encodes `tagValue`. After updating `tag_value` from `old_value` to `new_value`, the id must be refreshed to `tagKey + FILED_SP + new_value + FILED_SP + resource` so the subsequent Read (and future operations) resolve the correct attachment by the new value. The Update function sets the new id, then calls Read to refresh state.

### Decision 5: Read new value from `d.GetChange("tag_value")`

**Rationale:** To build the `ReplaceTags` correctly the Update needs the new value. Terraform Plugin SDK v2 provides `oldValue, newValue := d.GetChange("tag_value")`; only `newValue` is used (placed in `ReplaceTags`), and `deleteKeys` is `nil` since `ModifyResourceTags` updates the existing value in place. Since `tag_key` and `resource` are ForceNew they are unchanged, so the old attachment is identified by the current id split (`idSplit[0]` = tag_key, `idSplit[2]` = resource). Only `tag_value` changed, so only that field is diffed.

## Risks / Trade-offs

- **[Risk] Behavior change for users relying on ForceNew semantics:** Previously changing `tag_value` recreated the resource; now it updates in-place. This is strictly safer (no data loss), but the resource is no longer replaced.
  - **Mitigation:** This is the intended fix and is backward compatible — the field type/name and the resulting attachment state are the same; only the lifecycle (in-place vs. replace) changes for the better.

- **[Risk] `ModifyResourceTags` requires the tag key to already exist on the resource for the delete path:** The update path assumes the attachment exists (it was created by `AddResourceTag`). If the tag was removed out-of-band, the delete portion is a no-op and the replace adds the new value — which is the desired end state anyway.
  - **Mitigation:** The subsequent Read refreshes state from the API, so drift is reconciled.

- **[Risk] Id refresh must happen after the API call succeeds:** Setting the id before the API call would leave stale state on failure.
  - **Mitigation:** The Update function calls `ModifyTags` first (with retry), and only refreshes `d.SetId(...)` after it returns without error, then calls Read.
