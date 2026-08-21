## Context

The TencentCloud Terraform provider already manages tags via:
- `tencentcloud_tag` — manages a tag key/value pair (CRD: create/read/delete; no update since key+value are ForceNew).
- `tencentcloud_tag_attachment` — attaches a single tag key/value to a single resource, but `tag_value` is `ForceNew` (changing the value destroys and recreates the binding).

Neither gives a first-class, individually-addressable resource for the binding of a single tag key/value onto a single cloud resource (resource six-segment `qcs::...`) with a full CRUD lifecycle, including the ability to update only the tag value via `UpdateResourceTagValue`.

The cloud `tag/v20180813` SDK exposes four synchronous APIs that map cleanly to CRUD:
- `AddResourceTag(TagKey, TagValue, Resource)` — bind a tag to a resource six-segment.
- `GetResources(ResourceList, ...)` — query the tags bound to a list of resource six-segments; returns `ResourceTagMappingList []*ResourceTagMapping` where each mapping has `Resource` and `Tags []*Tag`, each `Tag` having `TagKey`, `TagValue`, `Category`, etc.
- `UpdateResourceTagValue(TagKey, TagValue, Resource)` — modify the value of a tag already bound to a resource.
- `DeleteResourceTag(TagKey, Resource)` — unbind a tag from a resource.

All four APIs are synchronous (no async polling needed).

## Goals / Non-Goals

**Goals:**
- Provide a `tencentcloud_tag_attachment_v2` RESOURCE_KIND_GENERAL resource that manages the full lifecycle (create/read/update/delete) of a single tag key/value bound to a single cloud resource six-segment.
- Support updating only the tag value (`tag_value`) while keeping `tag_key` and `resource` immutable.
- Use a stable composite ID (`tag_key` + `resource` joined by `tccommon.FILED_SP`) so each binding is uniquely addressable and importable.
- Follow the provider's established conventions: `tccommon` retry/error helpers, nil checks before `Set`, preserve id in read-before-clear logs, `NonRetryableError` on empty create responses, flat (non-nested) schema for list outputs.
- Add gomonkey-based unit tests that mock the cloud API client (no real cloud calls, no TF acceptance test suite).

**Non-Goals:**
- Batch binding/unbinding of multiple resources (that is `tencentcloud_tag_attachment`'s job).
- Managing the tag key/value pair existence itself (that is `tencentcloud_tag`'s job).
- Migrating any existing resource's schema or behavior.
- Exposing pagination as user-facing schema fields; the read path constructs the `GetResources` query from the composite ID and flattens the matched tag.

## Decisions

### Decision 1: Composite ID = `tag_key` + `resource` (joined by `tccommon.FILED_SP`)

**Rationale:** `AddResourceTag`/`UpdateResourceTagValue`/`DeleteResourceTag` key off `(TagKey, Resource)`. The natural unique identity of a binding is the pair `(tag_key, resource)`. The value `tag_value` is mutable (updated via `UpdateResourceTagValue`), so it must NOT be part of the ID (otherwise updating the value would recreate the resource).

**Alternative considered:** `tag_key` + `tag_value` + `resource` triple — rejected because `tag_value` is mutable and would force a destroy+recreate on value change, defeating the purpose of having an Update path.

### Decision 2: `tag_key` and `resource` are `ForceNew: true`; `tag_value` is mutable

**Rationale:** The cloud APIs do not allow changing the tag key or the resource of an existing binding in place — `UpdateResourceTagValue` only changes `TagValue`. Therefore changing `tag_key` or `resource` must recreate. `tag_value` changes call `UpdateResourceTagValue` in the Update function.

### Decision 3: Read uses `GetResources` and matches the tag whose `TagKey` equals the stored `tag_key`

**Rationale:** `GetResources` accepts the raw resource six-segment in `ResourceList` — no field decomposition needed. It returns `ResourceTagMappingList`, where each `ResourceTagMapping` has `Resource` and `Tags []*Tag`. The read path locates the mapping whose `Resource` matches the stored `resource` string, then selects the tag whose `TagKey` equals the stored `tag_key`. The matched tag's `TagValue` is set into state.

**Approach:** The `resource` field stored in state is the raw six-segment string (the same value passed to `AddResourceTag`/`UpdateResourceTagValue`/`DeleteResourceTag`). The service-layer method `DescribeTagAttachmentV2ById(ctx, tagKey, resource)` builds a `GetResources` request with `ResourceList = [resource]`, iterates `response.Response.ResourceTagMappingList` to find the mapping matching `resource`, and iterates its `Tags` to find the tag whose `TagKey == stored tag_key`. The matched `*tag.Tag` is returned (or nil if not found).

This avoids any decomposition of the six-segment, so it works uniformly for all resource types (including COS).

### Decision 4: Flat schema for the read result; no nested blocks in user schema

**Rationale:** Per provider rules, list outputs must be flattened, not nested under a `xxx_set`/`xxx_list` wrapper. Since `tencentcloud_tag_attachment_v2` represents a single binding (one tag key on one resource), the read result is a single matched tag. The `ResourceTagMappingList` / `Tags` fields from `GetResources` are used internally to locate the binding; they are NOT exposed as user-facing nested blocks. The user-facing schema is exactly the three arguments (`tag_key`, `tag_value`, `resource`) with no additional query-filter fields, because `GetResources` accepts the raw six-segment directly.

### Decision 5: Read empty handling

Per provider rules for RESOURCE_KIND_GENERAL (not datasource), if the read returns no matching row, we log `[CRUD] tag_attachment_v2 id=<id>` to preserve the scene, then `d.SetId("")`. We do NOT return `NonRetryableError` here (that rule is specific to RESOURCE_KIND_DATASOURCE). The Read path calls the service-layer `GetResources` lookup directly (with `ratelimit.Check` before the API call); transient API failures propagate to the caller.

### Decision 6: Tests use gomonkey mocks, not the TF acceptance suite

Per project rules for new resources, unit tests mock the cloud API client methods with gomonkey and exercise the business logic (Create/Read/Update/Delete functions) directly, without running `TF_ACC=1` against real cloud.

## Risks / Trade-offs

- **[Risk] Multiple tags on one resource** → `GetResources` returns all tags bound to the resource; we select the tag matching `tag_key`. If the binding was deleted out-of-band, no tag matches and we clear the ID. **Mitigation:** Log the id before clearing so the operator can trace the event.

- **[Trade-off] Mutability surface** → Only `tag_value` is mutable; changing `tag_key` or `resource` recreates. This matches the cloud API capabilities exactly and avoids pretending to support in-place changes the cloud cannot do.

- **[Trade-off] No async polling** → All four APIs are synchronous, so no `time.Sleep`/Read-loop polling is needed. If the cloud later makes any of these async, the CRUD functions will need to add polling; this is documented here for future maintainers.
