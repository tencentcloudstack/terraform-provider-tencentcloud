## 1. Schema Definition

- [x] 1.1 Add `input_type` field (Optional, TypeString) to the `tencentcloud_cls_config` resource schema in `resource_tc_cls_config.go`

## 2. CRUD Implementation

- [x] 2.1 Wire `input_type` into the Create function (`resourceTencentCloudClsConfigCreate`): read from `d.GetOk("input_type")` and set `request.InputType`
- [x] 2.2 Wire `InputType` into the Read function (`resourceTencentCloudClsConfigRead`): check `config.InputType != nil` and call `d.Set("input_type", ...)`
- [x] 2.3 Wire `input_type` into the Update function (`resourceTencentCloudClsConfigUpdate`): use `d.HasChange("input_type")` guard and set `request.InputType`

## 3. Documentation

- [x] 3.1 Update `resource_tc_cls_config.md` to include the `input_type` parameter in the Example Usage section

## 4. Unit Tests

- [x] 4.1 Add unit test cases for the new `input_type` parameter in `resource_tc_cls_config_test.go`