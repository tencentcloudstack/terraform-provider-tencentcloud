## 1. Schema Definition

- [x] 1.1 Add `kms_region` parameter to the `tencentcloud_cls_topic` resource schema in `tencentcloud/services/cls/resource_tc_cls_topic.go` (TypeString, Optional, Computed, description: "KMS region for the custom KMS key. Only effective when `encryption` is set to 1. If not set, the CLS default key (alias KMS-CLS) is used.")
- [x] 1.2 Add `kms_key_id` parameter to the `tencentcloud_cls_topic` resource schema in `tencentcloud/services/cls/resource_tc_cls_topic.go` (TypeString, Optional, Computed, description: "KMS key ID for the custom KMS key. Only effective when `encryption` is set to 1. If not set, the CLS default key (alias KMS-CLS) is used.")

## 2. Create Method Changes

- [x] 2.1 In `resourceTencentCloudClsTopicCreate`, after the existing `encryption` handling block, when `encryption = 1` and at least one of `kms_region` / `kms_key_id` is set, construct a `cls.CustomKmsInfo{}` struct populated with `KmsRegion` and/or `KmsKeyId` (using `helper.String(...)`) and assign it to `request.CustomKmsInfo`
- [x] 2.2 Ensure `CustomKmsInfo` is NOT set on the request when `encryption` is not `1` or neither kms field is provided

## 3. Read Method Changes

- [x] 3.1 In `resourceTencentCloudClsTopicRead`, after the existing `encryption` read-back block, add a nil check on `topic.CustomKmsInfo`; when non-nil, set `kms_region` from `topic.CustomKmsInfo.KmsRegion` (with nil check on the sub-field) and `kms_key_id` from `topic.CustomKmsInfo.KmsKeyId` (with nil check on the sub-field)

## 4. Update Method Changes

- [x] 4.1 In `resourceTencentCloudClsTopicUpdate`, within the existing `if d.HasChange("encryption")` block (and also when `kms_region` / `kms_key_id` change), construct a `cls.CustomKmsInfo{}` struct populated with the current `kms_region` and `kms_key_id` values when `encryption = 1` and at least one kms field is set, and assign it to `request.CustomKmsInfo`
- [x] 4.2 Add `d.HasChange("kms_region")` and `d.HasChange("kms_key_id")` to the conditions that set `hasChange = true` so the `ModifyTopic` call is triggered

## 5. Unit Tests

- [x] 5.1 Add a unit test in `tencentcloud/services/cls/resource_tc_cls_topic_test.go` to verify that when `encryption = 1` and `kms_region` / `kms_key_id` are set, the `CreateTopic` request contains a non-nil `CustomKmsInfo` with correct `KmsRegion` and `KmsKeyId` values (using gomonkey mocks)
- [x] 5.2 Add a unit test to verify that when `encryption` is not `1`, the `CreateTopic` request does NOT contain a `CustomKmsInfo` struct
- [x] 5.3 Add a unit test to verify that `kms_region` and `kms_key_id` are correctly read from a `TopicInfo` response with a non-nil `CustomKmsInfo`
- [x] 5.4 Add a unit test to verify the Read method does not panic when `CustomKmsInfo` is nil

## 6. Documentation

- [x] 6.1 Update `tencentcloud/services/cls/resource_tc_cls_topic.md` to add an example showing `encryption = 1` with `kms_region` and `kms_key_id` parameters

## 7. Validation

- [x] 7.1 Verify the code compiles successfully
- [x] 7.2 Verify no lint errors
