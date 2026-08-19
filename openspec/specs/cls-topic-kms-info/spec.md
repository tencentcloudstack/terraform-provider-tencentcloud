## ADDED Requirements

### Requirement: kms_region parameter in tencentcloud_cls_topic resource
The `tencentcloud_cls_topic` resource SHALL support an optional `kms_region` parameter (TypeString, Optional, Computed) that specifies the KMS region for a custom KMS key. The parameter SHALL be passed to the `CreateTopic` and `ModifyTopic` APIs as `request.CustomKmsInfo.KmsRegion`, and is only effective when `encryption` is set to `1`. The parameter SHALL be read from `topic.CustomKmsInfo.KmsRegion` in the `DescribeTopics` response during Read, with a nil check on `CustomKmsInfo`.

#### Scenario: Create topic with encryption and kms_region
- **WHEN** a user creates a `tencentcloud_cls_topic` resource with `encryption = 1` and `kms_region = "ap-guangzhou"`
- **THEN** the CreateTopic API request SHALL contain `CustomKmsInfo` with `KmsRegion` set to `"ap-guangzhou"`

#### Scenario: Create topic without kms_region
- **WHEN** a user creates a `tencentcloud_cls_topic` resource without specifying `kms_region`
- **THEN** the CreateTopic API request SHALL NOT set `CustomKmsInfo` (or SHALL leave `KmsRegion` unset within it), and the CLS backend SHALL use the default key

#### Scenario: Read kms_region from DescribeTopics response
- **WHEN** the Read method processes a `TopicInfo` where `CustomKmsInfo` is non-nil and `CustomKmsInfo.KmsRegion` is non-nil
- **THEN** the `kms_region` field in resource data SHALL be set to the value of `CustomKmsInfo.KmsRegion`

#### Scenario: Read kms_region when CustomKmsInfo is nil
- **WHEN** the Read method processes a `TopicInfo` where `CustomKmsInfo` is nil
- **THEN** the Read method SHALL NOT call `d.Set("kms_region", ...)` and SHALL NOT panic

#### Scenario: Update kms_region
- **WHEN** a user changes `kms_region` on an existing `tencentcloud_cls_topic` resource with `encryption = 1`
- **THEN** the ModifyTopic API request SHALL contain `CustomKmsInfo` with the updated `KmsRegion` value

### Requirement: kms_key_id parameter in tencentcloud_cls_topic resource
The `tencentcloud_cls_topic` resource SHALL support an optional `kms_key_id` parameter (TypeString, Optional, Computed) that specifies the KMS key ID for a custom KMS key. The parameter SHALL be passed to the `CreateTopic` and `ModifyTopic` APIs as `request.CustomKmsInfo.KmsKeyId`, and is only effective when `encryption` is set to `1`. The parameter SHALL be read from `topic.CustomKmsInfo.KmsKeyId` in the `DescribeTopics` response during Read, with a nil check on `CustomKmsInfo`.

#### Scenario: Create topic with encryption and kms_key_id
- **WHEN** a user creates a `tencentcloud_cls_topic` resource with `encryption = 1` and `kms_key_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"`
- **THEN** the CreateTopic API request SHALL contain `CustomKmsInfo` with `KmsKeyId` set to the specified key ID

#### Scenario: Create topic without kms_key_id
- **WHEN** a user creates a `tencentcloud_cls_topic` resource without specifying `kms_key_id`
- **THEN** the CreateTopic API request SHALL NOT set `CustomKmsInfo` (or SHALL leave `KmsKeyId` unset within it), and the CLS backend SHALL use the default key

#### Scenario: Read kms_key_id from DescribeTopics response
- **WHEN** the Read method processes a `TopicInfo` where `CustomKmsInfo` is non-nil and `CustomKmsInfo.KmsKeyId` is non-nil
- **THEN** the `kms_key_id` field in resource data SHALL be set to the value of `CustomKmsInfo.KmsKeyId`

#### Scenario: Read kms_key_id when CustomKmsInfo is nil
- **WHEN** the Read method processes a `TopicInfo` where `CustomKmsInfo` is nil
- **THEN** the Read method SHALL NOT call `d.Set("kms_key_id", ...)` and SHALL NOT panic

#### Scenario: Update kms_key_id
- **WHEN** a user changes `kms_key_id` on an existing `tencentcloud_cls_topic` resource with `encryption = 1`
- **THEN** the ModifyTopic API request SHALL contain `CustomKmsInfo` with the updated `KmsKeyId` value

### Requirement: CustomKmsInfo construction gated on encryption
The provider SHALL only construct and assign the `CustomKmsInfo` struct on the CreateTopic / ModifyTopic request when `encryption` is set to `1` and at least one of `kms_region` or `kms_key_id` is provided. When `encryption` is not `1`, the provider SHALL NOT set `CustomKmsInfo` on the request.

#### Scenario: CustomKmsInfo set when encryption is 1 and kms fields provided
- **WHEN** `encryption = 1` and the user provides `kms_region` and/or `kms_key_id`
- **THEN** the request SHALL include a non-nil `CustomKmsInfo` struct populated with the provided fields

#### Scenario: CustomKmsInfo not set when encryption is not 1
- **WHEN** `encryption` is not set to `1` (or not set at all)
- **THEN** the request SHALL NOT include a `CustomKmsInfo` struct, even if `kms_region` or `kms_key_id` are set in the configuration

#### Scenario: CustomKmsInfo not set when encryption is 1 but no kms fields provided
- **WHEN** `encryption = 1` but neither `kms_region` nor `kms_key_id` is provided
- **THEN** the request SHALL NOT include a `CustomKmsInfo` struct, and the CLS backend SHALL use the default key

### Requirement: kms_region and kms_key_id unit test coverage
The `tencentcloud_cls_topic` test file SHALL include unit tests that verify the `kms_region` and `kms_key_id` parameters are correctly passed to the CreateTopic / ModifyTopic API requests and correctly read from the DescribeTopics API response, using gomonkey mocks.

#### Scenario: Unit test for kms fields in Create
- **WHEN** the Create method is called with `encryption = 1`, `kms_region`, and `kms_key_id` set
- **THEN** the CreateTopic request SHALL contain a non-nil `CustomKmsInfo` with the correct `KmsRegion` and `KmsKeyId` values

#### Scenario: Unit test for kms fields in Read
- **WHEN** the Read method processes a `TopicInfo` with a non-nil `CustomKmsInfo` containing `KmsRegion` and `KmsKeyId`
- **THEN** the `kms_region` and `kms_key_id` fields in resource data SHALL be set to the corresponding values
