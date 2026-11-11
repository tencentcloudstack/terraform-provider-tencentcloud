## Context

The `tencentcloud_cls_logset` resource (`tencentcloud/services/cls/resource_tc_cls_logset.go`) currently defines a `tags` field of `schema.TypeMap` and manages tags indirectly through the separate TencentCloud Tag Service (`svctag.TagService.ModifyTags` / `DescribeResourceTags`). The CLS cloud API, however, natively supports a `Tags` array of `{Key, Value}` pairs on both `CreateLogset` and `ModifyLogset` requests, and returns `Tags` on the `DescribeLogsets` (`LogsetInfo.Tags`) response.

The SDK type `cls.Tag` (in `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016`) has exactly two fields: `Key *string` and `Value *string`. The requirement is to add parameters that map to `request.Tags.Key` and `request.Tags.Value`, i.e. bind tags natively through the CLS API rather than the Tag Service.

Other CLS resources in this codebase already follow the native-API pattern:
- `resource_tc_cls_topic.go` uses `TypeMap` `tags` but populates `request.Tags` directly.
- `resource_tc_cls_console.go` uses `TypeList` `tags` with nested `key`/`value` and populates `request.Tags` directly (ForceNew because ModifyConsole has no tags param).

For `cls_logset`, both `CreateLogset` and `ModifyLogset` accept `Tags`, so the tags can be fully managed (create + update) via the native CLS API.

## Goals / Non-Goals

**Goals:**
- Replace the existing `tags` (`TypeMap`, Tag-Service-managed) schema field on `tencentcloud_cls_logset` with a `tags` (`TypeList`) block containing nested `key` and `value` string fields, mapping directly to the CLS API `Tags []*Tag` parameter.
- Populate `CreateLogsetRequest.Tags` from the new schema in the Create function.
- Populate `ModifyLogsetRequest.Tags` from the new schema in the Update function (when `tags` changes).
- Read tags back from `DescribeLogsets` response (`LogsetInfo.Tags`) into the TypeList schema in the Read function.
- Remove the dependency on `svctag.TagService` for tag management on this resource.
- Update the resource documentation (`.md`) and test file accordingly.

**Non-Goals:**
- Changing the `logset_name`, `create_time`, `topic_count`, or `role_name` schema fields.
- Modifying the data source `tencentcloud_cls_logsets` (separate resource, out of scope).
- Adding support for the `LogsetId` custom-ID field of `CreateLogsetRequest` (not part of this requirement).
- Changes to the Tag Service infrastructure itself.

## Decisions

### Decision 1: Use `TypeList` with nested `key`/`value` (not `TypeMap`)

**Choice**: The `tags` field becomes `schema.TypeList` with `Elem: &schema.Resource{Schema: {"key": TypeString, "value": TypeString}}`.

**Rationale**: The requirement explicitly maps `request.Tags.Key` → SchemaName `Key` and `request.Tags.Value` → SchemaName `Value`, which describes a list of `{Key, Value}` objects — the natural Terraform representation is a `TypeList` of nested blocks. This also mirrors the cloud API `cls.Tag` struct exactly and is consistent with other resources in the codebase (`cls_console`, `sts_assume_role_operation`, `tat_command`) that use the same pattern for `cls.Tag` / similar `Key`/`Value` structs.

**Alternative considered**: Keep `TypeMap` and populate `request.Tags` from the map (as `cls_topic` does). Rejected because the requirement explicitly lists `Key` and `Value` as separate SchemaNames, which implies a structured block, not a flat map.

### Decision 2: Replace (not coexist with) the existing `tags` TypeMap field

**Choice**: The existing `tags` (`TypeMap`) field is replaced by the new `tags` (`TypeList`) field. There is no need for a separate deprecated field because the field name stays `tags` and the schema upgrade is handled by Terraform's schema migration.

**Rationale**: Terraform schema maps must have unique keys; two `tags` entries cannot coexist. Since the new field serves the same purpose (tag binding) but via a different mechanism (native CLS API vs. Tag Service), replacing is cleaner than introducing a differently-named field.

**Trade-off**: This is a **breaking change** for existing configurations that use the map-style `tags = { key = "value" }` syntax. Users must migrate to the list-style block syntax (`tags { key = "k" value = "v" }`). This is acceptable because it aligns the resource with the native CLS API contract.

### Decision 3: Manage tags through CLS API only (remove Tag Service usage)

**Choice**: Remove `svctag.TagService.ModifyTags` and `svctag.TagService.DescribeResourceTags` calls from the logset resource. Tags are created/updated via `CreateLogset`/`ModifyLogset` `Tags` parameter and read from `DescribeLogsets` `LogsetInfo.Tags`.

**Rationale**: The CLS API provides native tag CRUD support on logset create/modify/describe. Using the native API avoids the extra Tag Service round-trip and keeps tag state consistent within a single API call. The `ModifyLogset` API accepts `Tags`, so updates are supported without the Tag Service.

### Decision 4: Update path for tags

**Choice**: In the Update function, when `d.HasChange("tags")`, call `ModifyLogset` with the full new `Tags` list (not a diff). The `ModifyLogset` API replaces the entire tag set.

**Rationale**: `ModifyLogset` accepts `Tags []*Tag` as the desired final state, not incremental add/remove operations. Sending the complete new list is the correct semantic.

## Risks / Trade-offs

- **[Breaking change for existing users]** The `tags` field changes from `TypeMap` to `TypeList`. Existing HCL configurations using `tags = { "key" = "value" }` must be migrated to `tags { key = "k" value = "v" }` block syntax. → **Mitigation**: Document the migration clearly in the resource `.md` example and changelog. This is a necessary trade-off to align with the native CLS API contract.
- **[State migration]** Existing state files store `tags` as a map. After the schema change, Terraform will detect a type mismatch on refresh. → **Mitigation**: Users run `terraform apply -refresh-only` or `terraform plan` to let Terraform reconcile the state with the cloud (the CLS `DescribeLogsets` API returns `Tags` as a list, which maps cleanly to the new schema).
- **[Tag deduplication]** The CLS API allows max 10 tag key-value pairs with no duplicate keys. The Terraform schema does not enforce uniqueness at the `TypeList` level. → **Mitigation**: The CLS API will return an error on duplicate keys; the provider surfaces this error to the user. No additional client-side validation is needed.
- **[Removal of Tag Service dependency]** Resources previously tagged via the Tag Service will have their tags visible through the CLS API's `DescribeLogsets` response, so there is no data loss. → **Mitigation**: The Read function reads tags from `LogsetInfo.Tags`, which reflects all tags regardless of how they were originally set.
