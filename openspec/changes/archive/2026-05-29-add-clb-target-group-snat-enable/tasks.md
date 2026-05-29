## 1. Schema and CRUD Implementation

- [x] 1.1 Add `snat_enable` schema field (Optional, bool) to `ResourceTencentCloudClbTargetGroup` in `tencentcloud/services/clb/resource_tc_clb_target_group.go`
- [x] 1.2 Add `snat_enable` extraction in `resourceTencentCloudClbTargetCreate` and pass `SnatEnable` to the `CreateTargetGroup` service method request
- [x] 1.3 Add `snat_enable` read logic in `resourceTencentCloudClbTargetRead` to set state from `TargetGroupInfo.SnatEnable`
- [x] 1.4 Add `snat_enable` update logic in `resourceTencentCloudClbTargetUpdate` to call `ModifyTargetGroupAttribute` with `SnatEnable`

## 2. Service Layer

- [x] 2.1 Update `ClbService.CreateTargetGroup` method to accept and pass `SnatEnable` parameter to `CreateTargetGroupRequest`

## 3. Unit Tests

- [x] 3.1 Add unit tests in `tencentcloud/services/clb/resource_tc_clb_target_group_snat_enable_unit_test.go` using gomonkey to mock cloud API calls, covering create/read/update scenarios for `snat_enable`

## 4. Documentation

- [x] 4.1 Update `tencentcloud/services/clb/resource_tc_clb_target_group.md` to include `snat_enable` parameter in example usage
