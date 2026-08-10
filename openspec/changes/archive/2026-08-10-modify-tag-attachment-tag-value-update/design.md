## Context

The `tencentcloud_tag_attachment` resource binds a tag key-value pair to a cloud resource (identified by a six-segment resource description). The composite id is `tagKey + FILED_SP + tagValue + FILED_SP + resource`.

**Current state:**
- Resource file: `tencentcloud/services/tag/resource_tc_tag_attachment.go`
- Service layer: `tencentcloud/services/tag/service_tencentcloud_tag.go` (contains `DescribeTagTagAttachmentById` and `DeleteTagTagAttachmentById`)
- All three schema fields (`tag_key`, `tag_value`, `resource`) are `ForceNew: true`, so any change triggers destroy-and-recreate.

**Problem:** When `tag_value` changes, Terraform first deletes the old binding (`运营产品:A`) and then adds the new one (`运营产品:B`). For resources that use tag key-value pairs for access control, deleting the binding removes access, so the subsequent add request fails.

**API analysis (vendor SDK `tag/v20180813`):**

| API | TagKey | TagValue | Resource | Notes |
|-----|--------|----------|----------|-------|
| `AddResourceTag` | Yes (existing key) | Yes (existing value) | Yes | Create binding |
| `DeleteResourceTag` | Yes | — (not required by this API) | Yes | Delete binding (used by Delete) |
| `GetResources` | N/A (query by ResourceList) | N/A | N/A | Read: list tags for a resource |
| `UpdateResourceTagValue` | Yes (key of the binding) | Yes (the **new** value) | Yes | Update the tag value in place |

`UpdateResourceTagValue` is the right API: it accepts `TagKey`, `TagValue` (the new value), and `Resource`, and updates the tag value without removing the binding. The SDK already contains this API and its request/response models — no vendor upgrade is required.

## Goals / Non-Goals

**Goals:**
- Remove `ForceNew: true` from the `tag_value` field so modifying it triggers an in-place update.
- Implement an `Update` function that calls `UpdateResourceTagValue` when `tag_value` changes.
- Add a `ModifyTagTagAttachmentValue` service-layer method wrapping the `UpdateResourceTagValue` API.
- Update the composite id after a successful `tag_value` update (since `tagValue` is embedded in the id).
- Update the resource documentation to reflect that `tag_value` is updatable.
- Maintain full backward compatibility — existing configurations continue to work unchanged.

**Non-Goals:**
- Making `tag_key` or `resource` updatable (the `UpdateResourceTagValue` API only updates `tag_value`; `tag_key` and `resource` remain `ForceNew`).
- Adding a `Timeouts` block — `UpdateResourceTagValue` is a synchronous API call, no async waiting is needed.

## Decisions

### Decision 1: Use `UpdateResourceTagValue` API for in-place tag_value update

**Rationale:** `UpdateResourceTagValue` updates the tag value on a bound resource without removing the binding. This directly solves the access-control problem described in the proposal: the resource never loses its tag binding, so access is preserved across the update. The alternative (delete-then-add) is exactly the broken behavior we are replacing.

**Alternatives considered:**
- `ModifyResourcesTagValue` — this API takes `ServiceType`, `ResourceIds`, `TagKey`, `TagValue`, `ResourceRegion`, `ResourcePrefix`, i.e. it requires decomposing the six-segment resource description into service/region/prefix/id. The resource already stores the full six-segment `Resource` string, so `UpdateResourceTagValue` (which accepts `Resource` directly) is simpler and less error-prone.

### Decision 2: Keep `tag_key` and `resource` as ForceNew

**Rationale:** The `UpdateResourceTagValue` API only updates `tag_value`. There is no API to change `tag_key` or `resource` in place — those changes legitimately require destroy-and-recreate. Keeping them `ForceNew` preserves correct behavior.

### Decision 3: Update the composite id after a successful update

**Rationale:** The id is `tagKey + FILED_SP + tagValue + FILED_SP + resource`. Since `tagValue` changes, the id must be rebuilt with the new value so the subsequent `Read` looks up the correct binding. This mirrors how the `Create` function builds the id.

### Decision 4: Service-layer method wraps the API call

**Rationale:** Consistent with the existing pattern (`DescribeTagTagAttachmentById`, `DeleteTagTagAttachmentById`), a dedicated `ModifyTagTagAttachmentValue(ctx, tagKey, tagValue, resource)` method keeps the API-call logic in the service layer and the CRUD flow in the resource file. The resource `Update` function wraps the call with `resource.Retry(tccommon.WriteRetryTimeout, ...)` and `tccommon.RetryError`, consistent with the `Create` function.

## Risks / Trade-offs

- **[Risk] Behavior change for existing users who modify `tag_value`**: Previously this triggered destroy-and-recreate; now it triggers an in-place update.
  - **Mitigation:** This is the desired behavior and is a positive change. Existing users who relied on the recreate behavior are rare, and the in-place update is what they would expect. Document the behavior change in the changelog and `.md` file.

- **[Risk] `UpdateResourceTagValue` fails if the new tag value does not exist as a tag value under the tag key**: The API requires that the new `TagValue` already exists under the given `TagKey` in the tag system.
  - **Mitigation:** This is an API-level constraint. The provider surfaces the API error to the user. No provider-side pre-validation is added (consistent with the existing `Create` flow, which does not pre-validate tag existence either). The error message from the API is descriptive enough for the user to take action.
