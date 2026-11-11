## 1. Schema and Read Logic

- [x] 1.1 Add `zone_name` Computed attribute to the schema in `tencentcloud/services/teo/resource_tc_teo_zone_setting.go`
- [x] 1.2 Add `d.Set("zone_name", ...)` with nil check in the Read function of `tencentcloud/services/teo/resource_tc_teo_zone_setting.go`

## 2. Unit Tests

- [x] 2.1 Add unit test in `tencentcloud/services/teo/resource_tc_teo_zone_setting_test.go` using gomonkey to mock the DescribeZoneSetting API and verify `zone_name` is correctly set

## 3. Documentation

- [x] 3.1 Update `tencentcloud/services/teo/resource_tc_teo_zone_setting.md` to include `zone_name` in the example usage or attribute description
