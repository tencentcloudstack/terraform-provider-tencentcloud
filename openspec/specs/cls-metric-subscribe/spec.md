# cls-metric-subscribe Specification

## Purpose
TBD - created by archiving change add-cls-metric-subscribe-resource. Update Purpose after archive.
## Requirements
### Requirement: Resource schema definition

The system SHALL provide a Terraform resource `tencentcloud_cls_metric_subscribe` with the following top-level schema fields:
- `name` (Required, string): 订阅任务名称，长度不超过64字符
- `topic_id` (Required, string, ForceNew): 日志主题 ID
- `namespace` (Required, string): 云产品命名空间
- `metrics` (Required, list): 指标配置信息列表
- `instance_info` (Required, list, MaxItems=1): 实例配置信息
- `enable` (Optional, int): 任务开关，1 暂停 / 2 启用
- `task_id` (Computed, string): 订阅任务 ID
- `status` (Computed, int): 运行状态，0 创建中 / 1 暂停 / 2 运行中 / 3 异常
- `create_time` (Computed, int): 创建时间（秒级时间戳）
- `update_time` (Computed, int): 更新时间（秒级时间戳）

#### Scenario: Schema contains all required fields
- **WHEN** the resource schema is defined
- **THEN** the schema SHALL include `name`, `topic_id`, `namespace`, `metrics`, `instance_info` as Required fields and `task_id`, `status`, `create_time`, `update_time` as Computed fields

#### Scenario: topic_id is ForceNew
- **WHEN** the user changes `topic_id` on an existing resource
- **THEN** Terraform SHALL destroy and recreate the resource rather than update in place

### Requirement: Nested metrics schema

The `metrics` field SHALL be a TypeList of `MetricConfig` objects with the following nested schema:
- `metric_name` (Required, string): 指标名称
- `periods` (Optional, list of int): 统计周期（秒）
- `metric_labels` (Optional, list): 自定义指标标签

The `metric_labels` nested list SHALL contain:
- `key` (Required, string): 指标标签名称
- `value` (Required, string): 指标标签内容

#### Scenario: Metrics list with metric labels
- **WHEN** the user provides a metrics block with metric_name, periods, and metric_labels
- **THEN** the system SHALL correctly map these to the cloud API `MetricConfig` structure

### Requirement: Nested instance_info schema

The `instance_info` field SHALL be a TypeList with MaxItems=1 of `InstanceConfig` objects with the following nested schema:
- `instance_dimension` (Optional, list of string): 实例维度
- `instances` (Optional, list): 实例值列表

The `instances` nested list SHALL contain:
- `values` (Optional, list of string): 实例信息值列表

#### Scenario: Instance info with dimensions and instances
- **WHEN** the user provides an instance_info block with instance_dimension and instances
- **THEN** the system SHALL correctly map these to the cloud API `InstanceConfig` structure

### Requirement: Create operation

The system SHALL implement a Create function that calls the `CreateMetricSubscribe` API with `Name`, `TopicId`, `Namespace`, `Metrics`, `InstanceInfo` parameters and sets the resource ID to the composite format `topicId#taskId` (using `tccommon.FILED_SP` as separator).

#### Scenario: Successful creation
- **WHEN** the user creates a `tencentcloud_cls_metric_subscribe` resource with all required fields
- **THEN** the system SHALL call `CreateMetricSubscribe`, retrieve the returned `TaskId`, and set the resource ID to `topicId#taskId`

#### Scenario: Create returns empty TaskId
- **WHEN** the `CreateMetricSubscribe` API returns an empty `TaskId`
- **THEN** the system SHALL return a `NonRetryableError` to avoid writing an empty ID into state

### Requirement: Read operation

The system SHALL implement a Read function that calls the `DescribeMetricSubscribes` API with `TopicId` and a `Filter` (Key=taskId, Values=[taskId]) to query the resource, then populates the state from the first matching `MetricSubscribeInfo` record.

#### Scenario: Resource exists
- **WHEN** the Read function queries the resource by topicId and taskId
- **THEN** the system SHALL populate all schema fields from the returned `MetricSubscribeInfo` including computed fields `task_id`, `status`, `create_time`, `update_time`

#### Scenario: Resource not found
- **WHEN** the `DescribeMetricSubscribes` API returns an empty list
- **THEN** the system SHALL log the resource ID and call `d.SetId("")` to remove the resource from state

### Requirement: Update operation

The system SHALL implement an Update function that calls the `ModifyMetricSubscribe` API when any of `name`, `namespace`, `metrics`, `instance_info`, `enable` fields change, passing `TopicId` and `TaskId` as required identifiers along with the changed fields.

#### Scenario: Update name
- **WHEN** the user changes the `name` field
- **THEN** the system SHALL call `ModifyMetricSubscribe` with the new `Name` value, along with `TopicId` and `TaskId`

#### Scenario: Update enable status
- **WHEN** the user changes the `enable` field
- **THEN** the system SHALL call `ModifyMetricSubscribe` with the new `Enable` value

### Requirement: Delete operation

The system SHALL implement a Delete function that calls the `DeleteMetricSubscribe` API with `TaskId` and `TopicId` parameters.

#### Scenario: Successful deletion
- **WHEN** the user destroys the resource
- **THEN** the system SHALL call `DeleteMetricSubscribe` with `TaskId` and `TopicId` extracted from the composite resource ID

### Requirement: Import support

The resource SHALL support Terraform import using the composite ID format `topicId#taskId`.

#### Scenario: Import by composite ID
- **WHEN** the user runs `terraform import tencentcloud_cls_metric_subscribe.example topicId#taskId`
- **THEN** the system SHALL parse the composite ID, call Read, and populate the resource state

### Requirement: Provider registration

The system SHALL register the `tencentcloud_cls_metric_subscribe` resource in `tencentcloud/provider.go` and add the corresponding entry in `tencentcloud/provider.md`.

#### Scenario: Resource registered in provider
- **WHEN** the provider is initialized
- **THEN** the resource `tencentcloud_cls_metric_subscribe` SHALL be available for use in Terraform configurations

### Requirement: Error handling and retry

All CRUD operations SHALL use appropriate retry timeouts (`tccommon.WriteRetryTimeout` for Create/Update/Delete, `tccommon.ReadRetryTimeout` for Read) and wrap API errors with `tccommon.RetryError()`.

#### Scenario: API transient failure during create
- **WHEN** the `CreateMetricSubscribe` API returns a transient error
- **THEN** the system SHALL retry within the `WriteRetryTimeout` period before failing

#### Scenario: API transient failure during read
- **WHEN** the `DescribeMetricSubscribes` API returns a transient error
- **THEN** the system SHALL retry within the `ReadRetryTimeout` period before failing
