## Why

The `tencentcloud_teo_l7_acc_setting` resource currently exposes the `DescribeL7AccSetting` API's `ZoneSetting` response as two separate attributes (`zone_name` and `zone_config`), but does not provide a computed `zone_setting` attribute that represents the full `ZoneConfigParameters` response. Adding the `zone_setting` computed attribute allows users to access the complete site acceleration configuration as a single structured block, improving usability for programmatic consumption.

## What Changes

- Add a new computed attribute `zone_setting` to the `tencentcloud_teo_l7_acc_setting` resource schema, mapping to `DescribeL7AccSetting` response's `ZoneSetting` field (type `ZoneConfigParameters`).
- The `zone_setting` attribute will be of type `TypeList` with `MaxItems: 1` and `Computed: true`, containing `zone_name` (string, computed) and `zone_config` (list, computed) sub-fields that mirror the existing `ZoneConfigParameters` API structure.
- Update the `resourceTencentCloudTeoL7AccSettingRead` function to populate the `zone_setting` attribute from the `DescribeL7AccSetting` API response.

## Capabilities

### New Capabilities
- `teo-l7-acc-setting-zone-setting`: Add computed `zone_setting` attribute to the `tencentcloud_teo_l7_acc_setting` resource

### Modified Capabilities

## Impact

- Affected file: `tencentcloud/services/teo/resource_tc_teo_l7_acc_setting.go` (schema definition and read function)
- Affected file: `tencentcloud/services/teo/resource_tc_teo_l7_acc_setting_test.go` (unit tests)
- Affected file: `tencentcloud/services/teo/resource_tc_teo_l7_acc_setting.md` (documentation)
- API dependency: `DescribeL7AccSetting` (already used by this resource)
- No breaking changes: the new attribute is computed-only and additive
