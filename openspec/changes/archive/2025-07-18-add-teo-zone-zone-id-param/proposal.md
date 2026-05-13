## Why

The `tencentcloud_teo_zone` resource currently derives the `ZoneId` parameter for the `DeleteZone` API call from `d.Id()`. However, the `DeleteZone` API requires `ZoneId` as an explicit input parameter. Adding `zone_id` as a schema parameter makes the resource's delete operation more explicit and consistent, allowing users to reference the zone ID directly in their Terraform configuration.

## What Changes

- Add a new `zone_id` parameter (Computed, Optional) to the `tencentcloud_teo_zone` resource schema
- Update the `resourceTencentCloudTeoZoneDelete` function to read `zone_id` from `d.Get("zone_id")` instead of `d.Id()` for the `DeleteZone` API request's `ZoneId` field
- Update the `resourceTencentCloudTeoZoneRead` function to set `zone_id` from the API response
- Update the resource's `.md` documentation file accordingly

## Capabilities

### New Capabilities
- `teo-zone-zone-id-param`: Add `zone_id` parameter to the `tencentcloud_teo_zone` resource, mapping to `request.ZoneId` in the `DeleteZone` API

### Modified Capabilities

## Impact

- `tencentcloud/services/teo/resource_tc_teo_zone.go`: Schema definition, read and delete functions
- `tencentcloud/services/teo/resource_tc_teo_zone.md`: Documentation update
- `tencentcloud/services/teo/resource_tc_teo_zone_test.go`: Unit test update
