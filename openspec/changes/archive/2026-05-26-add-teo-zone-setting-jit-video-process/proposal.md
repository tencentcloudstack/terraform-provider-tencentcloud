## Why

The `tencentcloud_teo_zone_setting` resource is missing the `jit_video_process` (视频即时处理) configuration parameter that is already supported by the TencentCloud TEO API (`DescribeZoneSetting` and `ModifyZoneSetting`). Users need the ability to enable/disable JIT video processing for their TEO zones through Terraform.

## What Changes

- Add a new `jit_video_process` parameter to the `tencentcloud_teo_zone_setting` resource schema, containing a `switch` field (on/off) to control the JIT video processing feature.
- Update the resource's Read method to read `JITVideoProcess` from the `DescribeZoneSetting` API response and set it into state.
- Update the resource's Update method to include `JITVideoProcess` in the `ModifyZoneSetting` API request when the parameter is changed.
- Add corresponding unit tests and documentation.

## Capabilities

### New Capabilities
- `jit-video-process-setting`: Add JIT video processing configuration support to the TEO zone setting resource, allowing users to enable or disable the feature via a simple on/off switch.

### Modified Capabilities

## Impact

- `tencentcloud/services/teo/resource_tc_teo_zone_setting.go`: Add schema definition, read logic, and update logic for `jit_video_process`.
- `tencentcloud/services/teo/resource_tc_teo_zone_setting_test.go`: Add unit tests for the new parameter.
- `tencentcloud/services/teo/resource_tc_teo_zone_setting.md`: Update example usage documentation.
- `tencentcloud/provider.go`: No change needed (resource already registered).
- Dependency: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo v1.3.102` (already vendored).
