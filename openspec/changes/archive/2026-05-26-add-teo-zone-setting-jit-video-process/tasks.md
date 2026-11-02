## 1. Schema Definition

- [x] 1.1 Add `jit_video_process` parameter to the resource schema in `tencentcloud/services/teo/resource_tc_teo_zone_setting.go` as TypeList, Optional, Computed, MaxItems=1, with inner `switch` field (TypeString, Required)

## 2. Read Method

- [x] 2.1 Add read logic for `JITVideoProcess` in `resourceTencentCloudTeoZoneSettingRead` function: check nil, extract Switch field, set to state as `jit_video_process`

## 3. Update Method

- [x] 3.1 Add `jit_video_process` to the `mutableArgs` slice in `resourceTencentCloudTeoZoneSettingUpdate`
- [x] 3.2 Add update logic to construct `teo.JITVideoProcess{}` from state and assign to `request.JITVideoProcess`

## 4. Unit Tests

- [x] 4.1 Add unit tests in `tencentcloud/services/teo/resource_tc_teo_zone_setting_test.go` using gomonkey to mock the TEO API calls, verifying read and update behavior for `jit_video_process`

## 5. Documentation

- [x] 5.1 Update `tencentcloud/services/teo/resource_tc_teo_zone_setting.md` to include `jit_video_process` in the example usage
