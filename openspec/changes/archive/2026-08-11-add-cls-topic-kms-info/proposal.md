## Why

The `tencentcloud_cls_topic` resource already supports an `encryption` parameter that enables KMS-CLS cloud product key encryption, but users currently have no way to specify a **custom** KMS key (region + key ID) through Terraform. The CLS `ModifyTopic` and `CreateTopic` APIs accept a `CustomKmsInfo` object containing `KmsRegion` and `KmsKeyId`, and `DescribeTopics` returns the same object in each `TopicInfo`. Exposing these two parameters lets users bring their own KMS keys for CLS topic encryption instead of relying on the system default key (alias `KMS-CLS`).

## What Changes

- Add a new `kms_region` parameter (TypeString, Optional, Computed) to the `tencentcloud_cls_topic` resource schema. It maps to `request.CustomKmsInfo.KmsRegion` in the `ModifyTopic` / `CreateTopic` APIs and is read from `response.Topics[].CustomKmsInfo.KmsRegion` in the `DescribeTopics` API.
- Add a new `kms_key_id` parameter (TypeString, Optional, Computed) to the `tencentcloud_cls_topic` resource schema. It maps to `request.CustomKmsInfo.KmsKeyId` in the `ModifyTopic` / `CreateTopic` APIs and is read from `response.Topics[].CustomKmsInfo.KmsKeyId` in the `DescribeTopics` API.
- Both parameters are only effective when `encryption = 1` (KMS-CLS cloud product key encryption). When `CustomKmsInfo` is not provided, the CLS backend uses the default key.
- In the **Create** method: when `encryption = 1` and either `kms_region` or `kms_key_id` is set, construct a `cls.CustomKmsInfo{}` struct and assign it to `request.CustomKmsInfo`.
- In the **Update** method: when `encryption` changes (existing behavior) or `kms_region` / `kms_key_id` change, construct the `CustomKmsInfo` struct and assign it to `request.CustomKmsInfo` so the `ModifyTopic` API call carries the user's custom key info.
- In the **Read** method: read `KmsRegion` and `KmsKeyId` from `topic.CustomKmsInfo` (with nil checks on `CustomKmsInfo` and each sub-field) and set them in state.

## Capabilities

### New Capabilities
- `cls-topic-kms-info`: Adds custom KMS key (`kms_region`, `kms_key_id`) parameter support to the `tencentcloud_cls_topic` resource, allowing users to specify their own KMS key region and key ID for CLS topic encryption instead of using the system default key.

### Modified Capabilities
<!-- No existing capability requirements are changing -->

## Impact

- **Resource file**: `tencentcloud/services/cls/resource_tc_cls_topic.go` — add `kms_region` and `kms_key_id` schema fields; wire through Create, Update, and Read methods.
- **Test file**: `tencentcloud/services/cls/resource_tc_cls_topic_test.go` — add unit tests for the new parameters using gomonkey mocks.
- **Documentation**: `tencentcloud/services/cls/resource_tc_cls_topic.md` — update example usage to show `kms_region` and `kms_key_id` alongside `encryption = 1`.
- **Cloud API**: `CreateTopic` (input: `CustomKmsInfo.KmsRegion`, `CustomKmsInfo.KmsKeyId`), `ModifyTopic` (input: `CustomKmsInfo.KmsRegion`, `CustomKmsInfo.KmsKeyId`), `DescribeTopics` (output: `TopicInfo.CustomKmsInfo.KmsRegion`, `TopicInfo.CustomKmsInfo.KmsKeyId`). The vendor SDK already includes the `CustomKmsInfo` struct with both fields — no SDK upgrade required.
- **Backward compatibility**: fully backward compatible — both new parameters are Optional and Computed. Existing configurations without `kms_region` / `kms_key_id` continue to work unchanged (the CLS backend uses the default key).
