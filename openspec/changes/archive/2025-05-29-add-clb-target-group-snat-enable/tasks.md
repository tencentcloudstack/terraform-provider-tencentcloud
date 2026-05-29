## 1. Schema Definition

- [x] 1.1 Add `snat_enable` field (TypeBool, Optional, Computed) to the resource schema in `tencentcloud/services/clb/resource_tc_clb_target_group.go`

## 2. Create Logic

- [x] 2.1 Add `snatEnable *bool` parameter to `ClbService.CreateTargetGroup` method signature in `tencentcloud/services/clb/service_tencentcloud_clb.go`
- [x] 2.2 Set `request.SnatEnable` in `CreateTargetGroup` method when `snatEnable` is not nil
- [x] 2.3 Extract `snat_enable` from resource data in `resourceTencentCloudClbTargetCreate` and pass it to `CreateTargetGroup`

## 3. Read Logic

- [x] 3.1 Read `SnatEnable` from `TargetGroupInfo` response and set it in state in `resourceTencentCloudClbTargetRead`

## 4. Update Logic

- [x] 4.1 Add `snat_enable` change detection in `resourceTencentCloudClbTargetUpdate` and set `request.SnatEnable` on `ModifyTargetGroupAttributeRequest`

## 5. Documentation

- [x] 5.1 Update `tencentcloud/services/clb/resource_tc_clb_target_group.md` to include `snat_enable` parameter in example usage

## 6. Testing

- [x] 6.1 Add unit test for `snat_enable` parameter in `tencentcloud/services/clb/resource_tc_clb_target_group_test.go` using gomonkey mock
