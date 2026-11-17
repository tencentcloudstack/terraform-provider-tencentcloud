## ADDED Requirements

### Requirement: biz_type parameter in tencentcloud_cls_topic resource
The `tencentcloud_cls_topic` resource SHALL support a `biz_type` parameter of type integer that specifies the topic type. The parameter SHALL be Optional, Computed, and ForceNew. Valid values are 0 (log topic, default) and 1 (metric topic).

#### Scenario: Create topic with biz_type set to 1 (metric topic)
- **WHEN** a user creates a `tencentcloud_cls_topic` resource with `biz_type = 1`
- **THEN** the CreateTopic API SHALL be called with `BizType` set to 1, the resource ID SHALL be set to `"topic_id#1"`, and the resource SHALL be created as a metric topic

#### Scenario: Create topic without biz_type specified
- **WHEN** a user creates a `tencentcloud_cls_topic` resource without specifying `biz_type`
- **THEN** the CreateTopic API SHALL be called without the `BizType` parameter, the cloud API SHALL default to 0 (log topic), the resource ID SHALL be set to the plain `topic_id` (no suffix), and the `biz_type` field in state SHALL be populated with 0 from the Read response

#### Scenario: Resource ID format for log topic (biz_type=0)
- **WHEN** a resource is created with `biz_type = 0` or without specifying `biz_type`
- **THEN** the resource ID SHALL be the plain `topic_id` string (e.g., `"2f5764c1-c833-44c5-84c7-950979b2a278"`)

#### Scenario: Resource ID format for metric topic (biz_type=1)
- **WHEN** a resource is created with `biz_type = 1`
- **THEN** the resource ID SHALL be `"topic_id#1"` (e.g., `"2f5764c1-c833-44c5-84c7-950979b2a278#1"`)

### Requirement: ID parsing in Read/Update/Delete methods
The Read, Update, and Delete methods SHALL parse the resource ID using `parseClsTopicId` to extract the `topicId` and `bizType`.

#### Scenario: Read method parses ID
- **WHEN** the Read method is called for a resource with ID `"topic_id#1"`
- **THEN** `parseClsTopicId` SHALL return `topicId = "topic_id"` and `bizType = 1`, and `DescribeClsTopicById` SHALL be called with `topicId` and `bizType = 1`

#### Scenario: Read method parses plain ID
- **WHEN** the Read method is called for a resource with ID `"topic_id"` (no suffix)
- **THEN** `parseClsTopicId` SHALL return `topicId = "topic_id"` and `bizType = 0`, and `DescribeClsTopicById` SHALL be called with `topicId` and `bizType = nil`

#### Scenario: Update method parses ID
- **WHEN** the Update method is called for a resource
- **THEN** the method SHALL parse the ID to extract `topicId` for the `ModifyTopicRequest.TopicId`

#### Scenario: Delete method parses ID
- **WHEN** the Delete method is called for a resource
- **THEN** the method SHALL parse the ID to extract `topicId` for the `DeleteTopic` API call

### Requirement: Import support for metric topics
The resource SHALL support importing metric topics (biz_type=1) using the `"topic_id#1"` ID format.

#### Scenario: Import metric topic
- **WHEN** a user imports a metric topic with `terraform import tencentcloud_cls_topic.example "topic_id#1"`
- **THEN** the Read method SHALL parse the ID, extract `topicId` and `bizType = 1`, and correctly populate the state

#### Scenario: Import log topic
- **WHEN** a user imports a log topic with `terraform import tencentcloud_cls_topic.example "topic_id"`
- **THEN** the Read method SHALL parse the ID, extract `topicId` and `bizType = 0`, and correctly populate the state

#### Scenario: Read biz_type from DescribeTopics response
- **WHEN** the Read method is called for a `tencentcloud_cls_topic` resource
- **THEN** the `biz_type` field SHALL be set from `TopicInfo.BizType` in the DescribeTopics response, if the field is not nil

#### Scenario: Pass BizType to DescribeTopics when biz_type is set in state
- **WHEN** the Read method is called and the parsed ID has `bizType = 1`
- **THEN** the `DescribeClsTopicById` service method SHALL be called with `bizType` set, and the `DescribeTopicsRequest.BizType` SHALL be populated accordingly
- **AND** when the parsed ID has `bizType = 0`, `DescribeClsTopicById` SHALL be called with `nil` for `bizType`, and `DescribeTopicsRequest.BizType` SHALL NOT be set

#### Scenario: Attempt to update biz_type
- **WHEN** a user attempts to change the `biz_type` value on an existing `tencentcloud_cls_topic` resource
- **THEN** Terraform SHALL force resource recreation (ForceNew behavior), and the Update method SHALL reject the change via `immutableArgs` check

### Requirement: biz_type unit test coverage
The `tencentcloud_cls_topic` test file SHALL include unit tests that verify the `biz_type` parameter is correctly passed to the CreateTopic API request and correctly read from the DescribeTopics API response.

#### Scenario: Unit test for biz_type in Create
- **WHEN** the Create method is called with `biz_type` set to 1
- **THEN** the CreateTopic request SHALL contain `BizType` with value 1, and the resource ID SHALL be set to `"topic_id#1"`

#### Scenario: Unit test for biz_type in Read
- **WHEN** the Read method processes a TopicInfo with `BizType` set to 1
- **THEN** the `biz_type` field in the resource data SHALL be set to 1
