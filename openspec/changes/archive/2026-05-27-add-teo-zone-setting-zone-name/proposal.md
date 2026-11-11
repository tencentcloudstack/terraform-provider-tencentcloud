## Why

The `tencentcloud_teo_zone_setting` resource currently does not expose the `zone_name` field, which is returned by the `DescribeZoneSetting` API in the `ZoneSetting.ZoneName` response field. Users need to be able to read the zone name associated with a zone setting without having to query a separate data source.

## What Changes

- Add a new `zone_name` computed attribute (string, read-only) to the `tencentcloud_teo_zone_setting` resource.
- The field is populated from the `DescribeZoneSetting` API response (`ZoneSetting.ZoneName`).
- The `ModifyZoneSetting` API does not accept `ZoneName`, so this field is strictly read-only (Computed only).

## Capabilities

### New Capabilities
- `zone-name-attribute`: Add a computed `zone_name` attribute to `tencentcloud_teo_zone_setting` that exposes the zone name from the DescribeZoneSetting API response.

### Modified Capabilities

(none)

## Impact

- `tencentcloud/services/teo/resource_tc_teo_zone_setting.go`: Add schema definition and read logic for `zone_name`.
- `tencentcloud/services/teo/resource_tc_teo_zone_setting_test.go`: Add unit test to verify `zone_name` is set correctly.
- `tencentcloud/services/teo/resource_tc_teo_zone_setting.md`: Update example usage documentation.
