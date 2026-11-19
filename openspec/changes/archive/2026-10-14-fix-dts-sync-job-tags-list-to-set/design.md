## Context

The `tencentcloud_dts_sync_job` resource defines its `tags` field as `schema.TypeList`. Tags are sent to the DTS cloud API via `CreateSyncJob` (as `[]*dts.TagItem`) and read back via `DescribeSyncJobs` (as `[]*dts.TagItem`). The cloud API does not guarantee tag ordering on read — the same set of tags may come back in a different order than what was sent. With `TypeList`, Terraform compares element-by-element by index, so any reordering produces a non-empty diff on every refresh. This is a well-known pattern in the codebase: other resources (e.g. MongoDB `availability_zone_list`, DTS sync_config sub-fields) already use `TypeSet` for unordered collections.

Constraints:
- The `tags` field is `ForceNew` + `Computed` + `Optional`, so it is set only at create time and refreshed on read. There is no update path that needs changing.
- The cloud API request/response types (`dts.TagItem`) are unchanged — only the Terraform-side collection type changes.
- HCL block-list syntax is identical for `TypeList` and `TypeSet`, so user `.tf` files require no modification.

## Goals / Non-Goals

**Goals:**
- Eliminate spurious diffs on the `tags` field caused by unordered cloud API responses.
- Keep the nested element schema (`tag_key`, `tag_value`) unchanged.
- Keep the HCL configuration syntax backwards-compatible (no change to user `.tf` files).

**Non-Goals:**
- Do not change the cloud API calls (`CreateSyncJob`, `DescribeSyncJobs`).
- Do not change the `dts.TagItem` SDK struct usage.
- Do not modify other fields of `tencentcloud_dts_sync_job`.
- Do not change the `tags` handling in the data source `tencentcloud_dts_sync_jobs` (it uses a different schema and is not affected by the reported issue).

## Decisions

### D1: TypeSet with default HashResource

**Decision**: Change `Type: schema.TypeList` to `Type: schema.TypeSet`, keeping `Elem: &schema.Resource{...}`. Use the default `Set` hash function (`schema.HashResource`), which hashes the full resource schema — i.e. the combination of `tag_key` + `tag_value`.

**Rationale**:
- `tag_key` + `tag_value` together form a unique identity for a tag. The default `HashResource` hashes both fields, so duplicate tags are correctly de-duplicated and order is irrelevant.
- An alternative is a custom `schema.HashResource` with a custom function, but that adds complexity for no benefit since the default already covers both fields.

**Alternatives considered**:
- Custom hash function keyed only on `tag_key`: rejected because two tags with the same key but different values would collide; the cloud API allows the same key with different values in some contexts, so hashing both fields is safer.

### D2: Create handler reads Set via `(*schema.Set).List()`

**Decision**: In `resourceTencentCloudDtsSyncJobCreate`, change `v.([]interface{})` to `v.(*schema.Set).List()` when iterating the `tags` field. The iteration body that builds `dts.TagItem` is unchanged.

**Rationale**: `d.GetOk("tags")` returns `*schema.Set` for a TypeSet field. Calling `.List()` converts it to `[]interface{}`, which is the same type the existing loop already expects.

### D3: Read handler stays as `d.Set("tags", tagsList)` with `[]interface{}`

**Decision**: The Read handler already builds `tagsList` as `[]interface{}` of `map[string]interface{}` and calls `d.Set("tags", tagsList)`. This works correctly with `TypeSet` — Terraform SDK converts the slice into the Set using the same hash function. No change needed in the read body besides keeping the existing logic.

**Rationale**: `schema.ResourceData.Set` accepts `[]interface{}` for both TypeList and TypeSet. For TypeSet it hashes each element with the configured hash function. The existing code already produces the right shape.

### D4: Test assertions updated to Set semantics

**Decision**: In `resource_tc_dts_sync_job_test.go`, replace the `resource.TestCheckResourceAttrSet(..., "tags.#")` check with `resource.TestCheckTypeSetElemNestedAttrs` to verify the specific tag key/value pair regardless of order. The HCL config stays in block syntax (valid for Set).

## Risks / Trade-offs

- **[R1] State format change** → On first `terraform plan`/`apply` after upgrade, the state representation of `tags` changes from an ordered list to a set. Users see a one-time refresh diff; no actual resource is modified. This is the standard SDKv2 List→Set migration, well-tested in this codebase (MongoDB precedent). Mitigation: document the one-time refresh in the changelog.
- **[R2] TypeSet hash collisions** → If two tags have identical `tag_key` + `tag_value`, they collapse into one Set element. The cloud API also treats identical tags as duplicates, so this matches API semantics. No risk.
- **[R3] Import compatibility** → `tags` is `ForceNew`+`Computed`, not a read-from-id field, so import (`schema.ImportStatePassthrough`) is unaffected.

## Migration Plan

1. Apply the schema + Create handler change in `resource_tc_dts_sync_job.go`.
2. Update test file assertions.
3. Run `gofmt` on modified `.go` files (finalize step).
4. Regenerate docs via `make doc` (finalize step).
5. Rollback: `git revert` the change; the state automatically reverts to List representation on next refresh with the old provider version.