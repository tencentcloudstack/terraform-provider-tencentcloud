## ADDED Requirements

### Requirement: biz_type parameter in tencentcloud_cls_topic resource
The `tencentcloud_cls_topic` resource SHALL support a `biz_type` parameter of type integer that specifies the topic type. The parameter SHALL be Optional, Computed, and ForceNew. Valid values are 0 (log topic, default) and 1 (metric topic).

#### Scenario: Create topic with biz_type set to 1 (metric topic)
- **WHEN** a user creates a `tencentcloud_cls_topic` resource with `biz_type = 1`
- **THEN** the CreateTopic API SHALL be called with `BizType` set to 1, and the resource SHALL be created as a metric topic

#### Scenario: Create topic without biz_type specified
- **WHEN** a user creates a `tencentcloud_cls_topic` resource without specifying `biz_type`
- **THEN** the CreateTopic API SHALL be called without the `BizType` parameter, the cloud API SHALL default to 0 (log topic), and the `biz_type` field in state SHALL be populated with 0 from the Read response

#### Scenario: Read biz_type from DescribeTopics response
- **WHEN** the Read method is called for a `tencentcloud_cls_topic` resource
- **THEN** the `biz_type` field SHALL be set from `TopicInfo.BizType` in the DescribeTopics response, if the field is not nil

#### Scenario: Attempt to update biz_type
- **WHEN** a user attempts to change the `biz_type` value on an existing `tencentcloud_cls_topic` resource
- **THEN** Terraform SHALL force resource recreation (ForceNew behavior), and the Update method SHALL reject the change via `immutableArgs` check

### Requirement: biz_type unit test coverage
The `tencentcloud_cls_topic` test file SHALL include unit tests that verify the `biz_type` parameter is correctly passed to the CreateTopic API request and correctly read from the DescribeTopics API response.

#### Scenario: Unit test for biz_type in Create
- **WHEN** the Create method is called with `biz_type` set to 1
- **THEN** the CreateTopic request SHALL contain `BizType` with value 1

#### Scenario: Unit test for biz_type in Read
- **WHEN** the Read method processes a TopicInfo with `BizType` set to 1
- **THEN** the `biz_type` field in the resource data SHALL be set to 1
