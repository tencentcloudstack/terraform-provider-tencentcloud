## 1. Schema Definition

- [x] 1.1 Modify the `tags` field in `ResourceTencentCloudClsLogset()` schema from `TypeMap` to `TypeList` with nested `key`/`value` string fields in `tencentcloud/services/cls/resource_tc_cls_logset.go`
- [x] 1.2 Remove the `svctag` import from `tencentcloud/services/cls/resource_tc_cls_logset.go` (no longer needed for tag management)

## 2. Create Function

- [x] 2.1 In `resourceTencentCloudClsLogsetCreate`, populate `CreateLogsetRequest.Tags` from the new `tags` TypeList schema (iterate list, build `[]*cls.Tag{Key, Value}`)
- [x] 2.2 Remove the `svctag.TagService.ModifyTags` call from the Create function (tags are now set natively via `CreateLogset`)

## 3. Read Function

- [x] 3.1 In `resourceTencentCloudClsLogsetRead`, read tags from `DescribeLogsets` response (`LogsetInfo.Tags`), convert `[]*cls.Tag` to `[]map[string]interface{}` with `key`/`value`, and `d.Set("tags", ...)`
- [x] 3.2 Remove the `svctag.TagService.DescribeResourceTags` call from the Read function

## 4. Update Function

- [x] 4.1 In `resourceTencentCloudClsLogsetUpdate`, when `d.HasChange("tags")`, populate `ModifyLogsetRequest.Tags` with the full new `[]*cls.Tag` list and call `ModifyLogset`
- [x] 4.2 Remove the `svctag.TagService.ModifyTags` / `svctag.DiffTags` call from the Update function

## 5. Tests

- [x] 5.1 Update `testAccClsLogset_basic` in `tencentcloud/services/cls/resource_tc_cls_logset_test.go` to use the new TypeList `tags` block syntax
- [x] 5.2 Add a test step that verifies tag update (change tags and re-apply) using the Terraform test suite (existing resource modification uses TF test suite)

## 6. Documentation

- [x] 6.1 Update the example in `tencentcloud/services/cls/resource_tc_cls_logset.md` to use the new `tags` TypeList block syntax
