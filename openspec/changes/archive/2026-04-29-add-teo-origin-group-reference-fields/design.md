## Context

The `tencentcloud_teo_origin_group` resource manages TEO origin groups. The current `references` computed block in the resource schema only exposes `instance_type`, `instance_id`, and `instance_name` from the `OriginGroupReference` struct. The cloud API SDK (`OriginGroupReference`) has been updated with 3 additional fields: `ZoneId`, `ZoneName`, and `AliasZoneName` that provide zone-level context for each reference entry. These fields are important for cross-zone reference scenarios where origin groups are referenced by resources in different zones.

Current state of the `references` block schema:
- `instance_type` (Computed, string)
- `instance_id` (Computed, string)
- `instance_name` (Computed, string)

Cloud API `OriginGroupReference` struct fields:
- `InstanceType` - already mapped
- `InstanceId` - already mapped
- `InstanceName` - already mapped
- `ZoneId` - NOT yet mapped (new)
- `ZoneName` - NOT yet mapped (new)
- `AliasZoneName` - NOT yet mapped (new)

## Goals / Non-Goals

**Goals:**
- Add `zone_id`, `zone_name`, and `alias_zone_name` computed attributes to the `references` nested block
- Ensure the Read method populates these new fields from the `DescribeOriginGroup` API response
- Maintain full backward compatibility - these are additive computed fields

**Non-Goals:**
- Modifying Create, Update, or Delete operations (these fields are read-only from the API)
- Changing the existing `references` sub-attributes or their behavior
- Adding any new top-level resource parameters

## Decisions

1. **Schema type for new fields**: Use `schema.TypeString` with `Computed: true` for all three new fields. This matches the existing pattern in the `references` block and the cloud API field types (all `*string`).

2. **Read logic placement**: Add the new field mappings inside the existing `references` loop in `resourceTencentCloudTeoOriginGroupRead`, following the same nil-check pattern as existing fields (`if respData.ZoneId != nil { ... }`).

3. **No changes to mutableArgs**: Since these are computed-only fields in the references block (not top-level mutable fields), no changes needed to the `mutableArgs` list in the Update method.

4. **Test approach**: Use gomonkey-based mock unit tests (not Terraform acceptance tests) since this is a modification to an existing resource. Mock the `DescribeTeoOriginGroupById` service method to return test data with the new fields populated.

## Risks / Trade-offs

- [Risk] Cloud API may return nil for these new fields → Mitigation: Follow existing nil-check pattern before setting field values
- [Risk] State migration not needed → Mitigation: Adding computed fields is backward compatible; existing state files will simply have these fields empty until next refresh
