## Why

The `tencentcloud_teo_l7_acc_setting` resource had a redundant `zone_setting` computed attribute that duplicated the data already available via `zone_name` and `zone_config`. This attribute added unnecessary schema complexity and maintenance burden. Additionally, the `ModifyL7AccSetting` API's `ZoneConfig` struct supports `NetworkErrorLogging` configuration, but the Terraform resource did not expose it.

## What Changes

- Remove the `zone_setting` computed attribute from the resource schema and its population logic in the Read function
- Remove the `zone_setting`-related unit tests (gomonkey mock tests)
- Add `network_error_logging` as an Optional sub-field under `zone_config` in the schema, Read function, and Update function
- Update documentation to include `network_error_logging` in the example

## Capabilities

### New Capabilities
- `teo-l7-acc-setting-network-error-logging`: Add `network_error_logging` configuration block (with `switch` field) under `zone_config` in the `tencentcloud_teo_l7_acc_setting` resource, mapping to `ZoneConfig.NetworkErrorLogging` in the `ModifyL7AccSetting`/`DescribeL7AccSetting` APIs.

### Modified Capabilities
- `teo-l7-acc-setting-zone-setting-removed`: Remove the redundant `zone_setting` computed attribute. Users should use `zone_name` and `zone_config` directly.

## Impact

- `tencentcloud/services/teo/resource_tc_teo_l7_acc_setting.go`: Remove `zone_setting` schema and Read logic; add `network_error_logging` to `zone_config` schema, Read, and Update functions
- `tencentcloud/services/teo/resource_tc_teo_l7_acc_setting_test.go`: Remove `zone_setting` unit tests
- `tencentcloud/services/teo/resource_tc_teo_l7_acc_setting.md`: Add `network_error_logging` to example
- API dependency: `ModifyL7AccSetting` and `DescribeL7AccSetting` (already used), `NetworkErrorLoggingParameters` struct in SDK
