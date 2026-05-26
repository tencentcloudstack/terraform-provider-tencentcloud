## 1. Schema Definition

- [x] 1.1 Add `snat_enable` field (TypeBool, Optional) to the resource schema in `tencentcloud/services/clb/resource_tc_clb_target_group.go`

## 2. Create Logic

- [x] 2.1 Extract `snat_enable` from resource data in `resourceTencentCloudClbTargetCreate` and pass it to `ClbService.CreateTargetGroup`
- [x] 2.2 Extend `ClbService.CreateTargetGroup` method signature in `tencentcloud/services/clb/service_tencentcloud_clb.go` to accept `snatEnable *bool` parameter and set `request.SnatEnable` when non-nil

## 3. Update Logic

- [x] 3.1 Add `snat_enable` change detection in `resourceTencentCloudClbTargetUpdate` and pass `SnatEnable` to the `ModifyTargetGroupAttributeRequest`

## 4. Documentation

- [x] 4.1 Update `tencentcloud/services/clb/resource_tc_clb_target_group.md` to add example usage with `snat_enable` parameter

## 5. Testing

- [x] 5.1 Add unit tests in `tencentcloud/services/clb/resource_tc_clb_target_group_test.go` using gomonkey mock to test Create and Update with `snat_enable` parameter
