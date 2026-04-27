## 1. Schema Definition

- [x] 1.1 Add `zone_setting` computed attribute to the `tencentcloud_teo_l7_acc_setting` resource schema in `resource_tc_teo_l7_acc_setting.go`. The attribute SHALL be of type `TypeList` with `MaxItems: 1`, `Computed: true`, and contain `zone_name` (TypeString, Computed) and `zone_config` (TypeList, MaxItems:1, Computed) sub-attributes. The `zone_config` sub-attribute schema SHALL match the existing top-level `zone_config` attribute structure.

## 2. Read Function Update

- [x] 2.1 Update `resourceTencentCloudTeoL7AccSettingRead` in `resource_tc_teo_l7_acc_setting.go` to populate the `zone_setting` computed attribute from the `DescribeL7AccSetting` API response's `ZoneSetting` field. Map `ZoneSetting.ZoneName` to `zone_setting.zone_name` and `ZoneSetting.ZoneConfig` to `zone_setting.zone_config` following the same pattern used for the existing top-level attributes.

## 3. Unit Tests

- [x] 3.1 Add unit test cases in `resource_tc_teo_l7_acc_setting_test.go` using gomonkey mocks to verify that the `zone_setting` computed attribute is correctly populated in the read function when the `DescribeL7AccSetting` API returns data.

## 4. Documentation

- [x] 4.1 Update `resource_tc_teo_l7_acc_setting.md` to document the new `zone_setting` computed attribute with its sub-attributes, following the existing documentation format in this resource.
