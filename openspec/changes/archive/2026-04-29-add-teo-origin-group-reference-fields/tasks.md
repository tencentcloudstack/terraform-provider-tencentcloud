## 1. Schema and Read Logic

- [x] 1.1 Add `zone_id`, `zone_name`, `alias_zone_name` computed string fields to the `references` nested block schema in `ResourceTencentCloudTeoOriginGroup()` in `tencentcloud/services/teo/resource_tc_teo_origin_group.go`
- [x] 1.2 Add nil-check and field-setting logic for `zone_id`, `zone_name`, `alias_zone_name` in the `references` loop within `resourceTencentCloudTeoOriginGroupRead()` in `tencentcloud/services/teo/resource_tc_teo_origin_group.go`

## 2. Unit Tests

- [x] 2.1 Add gomonkey-based mock unit tests for the new `zone_id`, `zone_name`, `alias_zone_name` computed fields in `tencentcloud/services/teo/resource_tc_teo_origin_group_test.go`
- [x] 2.2 Run unit tests with `go test -gcflags=all=-l` to verify the new fields are correctly populated

## 3. Documentation

- [x] 3.1 Update `tencentcloud/services/teo/resource_tc_teo_origin_group.md` to reflect the new computed fields in the references block
