## Why

The cloud API `OriginGroupReference` struct has been updated with new fields (`ZoneId`, `ZoneName`, `AliasZoneName`) that provide zone-level information for each reference entry. The current Terraform resource `tencentcloud_teo_origin_group` only exposes `instance_type`, `instance_id`, and `instance_name` in the `references` computed block, missing the zone-related information. Users need visibility into which zone each reference belongs to, especially for cross-zone reference scenarios.

## What Changes

- Add 3 new computed attributes to the `references` nested block of `tencentcloud_teo_origin_group` resource:
  - `zone_id`: The zone ID of the referenced instance (maps to `OriginGroupReference.ZoneId`)
  - `zone_name`: The zone name of the referenced instance (maps to `OriginGroupReference.ZoneName`)
  - `alias_zone_name`: The alias zone name of the referenced instance (maps to `OriginGroupReference.AliasZoneName`)

These are all read-only computed fields sourced from the `DescribeOriginGroup` API response, requiring no changes to Create/Update/Delete operations.

## Capabilities

### New Capabilities
- `teo-origin-group-reference-zone-fields`: Add zone_id, zone_name, alias_zone_name computed attributes to the references block of tencentcloud_teo_origin_group resource

### Modified Capabilities

## Impact

- `tencentcloud/services/teo/resource_tc_teo_origin_group.go`: Add new schema fields in references block and read logic
- `tencentcloud/services/teo/resource_tc_teo_origin_group_test.go`: Add unit tests for new computed fields
- `tencentcloud/services/teo/resource_tc_teo_origin_group.md`: Update example documentation
- No breaking changes - all new fields are computed and backward compatible
