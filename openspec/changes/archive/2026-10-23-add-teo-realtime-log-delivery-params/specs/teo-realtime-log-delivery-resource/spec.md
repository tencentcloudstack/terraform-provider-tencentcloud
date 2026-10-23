## ADDED Requirements

### Requirement: Resource Schema Definition
The system SHALL define a Terraform resource `tencentcloud_teo_realtime_log_delivery` that manages a TEO realtime log delivery task, with a composite ID of `zoneId#taskId` (using `tccommon.FILED_SP` as separator) and import support. The schema SHALL include the following fields:
- `zone_id` (Required, ForceNew, TypeString): 站点 ID
- `task_name` (Required, TypeString): 实时日志投递任务名称
- `task_type` (Required, TypeString): 任务类型，取值含 `cls`/`custom_endpoint`/`s3`/`aws_s3`/`log_analysis`
- `entity_list` (Required, TypeList of TypeString): 实体列表
- `log_type` (Required, TypeString): 数据投递类型
- `area` (Required, TypeString): 数据投递区域
- `fields` (Required, TypeList of TypeString): 预设字段列表
- `custom_fields` (Optional, Computed, TypeList): 自定义字段列表（HTTP 头/Cookie 提取），每元素含 `name`/`value`/`enabled`
- `custom_expression_fields` (Optional, Computed, TypeList): 自定义表达式字段列表，每元素含 `name`(Required, TypeString)、`expression`(Required, TypeString)、`enabled`(Optional, TypeBool)。对应云 API `CustomExpressionFields`。注意：若与 `Fields`/`CustomFields` 同名，以 `CustomExpressionFields` 取值为准。
- `delivery_conditions` (Optional, Computed, TypeList): 投递过滤条件
- `sample` (Required, TypeInt): 采样比例
- `log_format` (Optional, Computed, TypeList, MaxItems=1): 日志输出格式块，含 `format_type`/`batch_prefix`/`batch_suffix`/`record_prefix`/`record_suffix`/`record_delimiter`/`field_delimiter`，以及新增的 `record_template`(Optional, TypeString，仅当 `format_type = template` 时生效)
- `cls` (Optional, Computed, TypeList, MaxItems=1): CLS 配置块
- `custom_endpoint` (Optional, Computed, TypeList, MaxItems=1): 自定义 HTTP 配置块
- `s3` (Optional, Computed, TypeList, MaxItems=1): S3 兼容存储配置块
- `aws_s3` (Optional, Computed, TypeList, MaxItems=1): AWS 原生 S3 配置块（当 `task_type = aws_s3` 时必填），含 `endpoint`/`bucket_name`/`path`/`role_arn`/`compress_type`
- `delivery_status` (Optional, Computed, TypeString): 任务状态
- `task_id` (Computed, TypeString): 任务 ID

The schema SHALL NOT expose `DescribeRealtimeLogDeliveryTasks` query parameters (`Offset`/`Limit`/`Filters.Name`/`Filters.Fuzzy`) as resource schema fields; these are internal service-layer query parameters.

#### Scenario: Schema defines all CRUD fields including new aws_s3, custom_expression_fields, record_template
- **WHEN** the resource schema is defined
- **THEN** it SHALL include `aws_s3` block with `endpoint`/`bucket_name`/`path`/`role_arn`/`compress_type` sub-fields, `custom_expression_fields` block with `name`/`expression`/`enabled` sub-fields, and `record_template` sub-field inside the `log_format` block

#### Scenario: New fields are optional and computed for backward compatibility
- **WHEN** a user does not set `aws_s3`, `custom_expression_fields`, or `log_format.record_template` in the Terraform configuration
- **THEN** the system SHALL populate these fields from the API response (Computed behavior) without triggering a plan diff

#### Scenario: ForceNew fields prevent in-place update
- **WHEN** `zone_id` is changed in the Terraform configuration
- **THEN** the resource SHALL be destroyed and recreated

#### Scenario: Composite ID format
- **WHEN** a realtime log delivery task is created
- **THEN** the resource ID SHALL be `zoneId#taskId` using `tccommon.FILED_SP` as separator

### Requirement: Resource Create Operation
The system SHALL implement Create by calling `CreateRealtimeLogDeliveryTask` API with ZoneId and the configured task parameters. When `custom_expression_fields` is set, the system SHALL populate `request.CustomExpressionFields` (a `[]*CustomExpressionField` with Name/Expression/Enabled). When `aws_s3` is set, the system SHALL populate `request.AWSS3` (an `*AWSS3` with Endpoint/BucketName/Path/RoleARN/CompressType). When `log_format` is set with `record_template`, the system SHALL set `request.LogFormat.RecordTemplate`. After obtaining the TaskId, the system SHALL set the composite ID and then poll `DescribeRealtimeLogDeliveryTasks` until the task status reaches `enabled` (via `checkRealtimeLogDeliveryStatus`).

#### Scenario: Create task with custom_expression_fields
- **WHEN** the Terraform configuration includes `custom_expression_fields` blocks
- **THEN** the system SHALL call `CreateRealtimeLogDeliveryTask` with `CustomExpressionFields` populated, mapping each block's `name`/`expression`/`enabled` to the corresponding `CustomExpressionField` struct fields

#### Scenario: Create task with aws_s3 delivery target
- **WHEN** `task_type = "aws_s3"` and the `aws_s3` block is configured
- **THEN** the system SHALL call `CreateRealtimeLogDeliveryTask` with `AWSS3` populated, mapping `endpoint`/`bucket_name`/`path`/`role_arn`/`compress_type` to the `AWSS3` struct fields

#### Scenario: Create task with log_format.record_template
- **WHEN** the `log_format` block is configured and `format_type = "template"` with `record_template` set
- **THEN** the system SHALL set `request.LogFormat.RecordTemplate` to the configured value

#### Scenario: Create validates non-empty TaskId
- **WHEN** `CreateRealtimeLogDeliveryTask` returns successfully
- **THEN** the system SHALL verify `response.Response.TaskId` is non-empty before setting the ID; if empty, SHALL return a NonRetryableError

### Requirement: Resource Read Operation
The system SHALL implement Read by calling `DescribeRealtimeLogDeliveryTasks` (via `service.DescribeTeoRealtimeLogDeliveryById`) filtered by `task-id`. The system SHALL backfill `custom_expression_fields` (from `respData.CustomExpressionFields`), `aws_s3` (from `respData.AWSS3`, note: may be null), and `log_format.record_template` (from `respData.LogFormat.RecordTemplate`) into state. All Response field accesses SHALL perform nil checks before calling `d.Set`. If the response is empty (no matching task), the system SHALL log `[CRUD]` with the id and then `d.SetId("")`.

#### Scenario: Read backfills custom_expression_fields
- **WHEN** the API response contains `CustomExpressionFields`
- **THEN** the system SHALL populate the state's `custom_expression_fields` list, mapping Name/Expression/Enabled for each element

#### Scenario: Read backfills aws_s3
- **WHEN** the API response contains a non-null `AWSS3`
- **THEN** the system SHALL populate the state's `aws_s3` block with endpoint/bucket_name/path/role_arn/compress_type mapped from the AWSS3 struct

#### Scenario: Read backfills log_format.record_template
- **WHEN** the API response contains `LogFormat` with a non-null `RecordTemplate`
- **THEN** the system SHALL populate the state's `log_format.0.record_template`

#### Scenario: Read handles null aws_s3 without error
- **WHEN** the API response has `AWSS3 = null` (e.g., task_type is not aws_s3)
- **THEN** the system SHALL skip setting the `aws_s3` block without error

### Requirement: Resource Update Operation
The system SHALL implement Update by calling `ModifyRealtimeLogDeliveryTask` API. The `immutableArgs` SHALL be `["task_type", "log_type", "area", "cls"]` (changing these returns an error). The `mutableArgs` SHALL include `task_name`, `delivery_status`, `entity_list`, `fields`, `custom_fields`, `custom_expression_fields`, `delivery_conditions`, `sample`, `log_format`, `custom_endpoint`, `s3`, and `aws_s3`; when any changes, the system SHALL call `ModifyRealtimeLogDeliveryTask` with `CustomExpressionFields`/`AWSS3`/`LogFormat.RecordTemplate` populated the same way as Create. After the call, the system SHALL poll `DescribeRealtimeLogDeliveryTasks` (via `checkRealtimeLogDeliveryStatus`) until status stabilizes.

#### Scenario: Update custom_expression_fields
- **WHEN** `custom_expression_fields` changes in the Terraform configuration
- **THEN** the system SHALL call `ModifyRealtimeLogDeliveryTask` with `CustomExpressionFields` populated

#### Scenario: Update aws_s3 configuration
- **WHEN** `aws_s3` block changes in the Terraform configuration
- **THEN** the system SHALL call `ModifyRealtimeLogDeliveryTask` with `AWSS3` populated

#### Scenario: Update log_format.record_template
- **WHEN** `log_format.record_template` changes in the Terraform configuration
- **THEN** the system SHALL call `ModifyRealtimeLogDeliveryTask` with `LogFormat.RecordTemplate` populated

#### Scenario: Immutable task_type change rejected
- **WHEN** `task_type` changes (e.g., from `s3` to `aws_s3`)
- **THEN** the system SHALL return an error indicating the argument cannot be changed

### Requirement: Resource Delete Operation
The system SHALL implement Delete by calling `DeleteRealtimeLogDeliveryTask` API with ZoneId and TaskId (parsed from the composite ID). The call SHALL be wrapped with `resource.Retry(tccommon.WriteRetryTimeout, ...)` and `tccommon.RetryError`. After success, the system SHALL poll `DescribeRealtimeLogDeliveryTasks` (via `checkRealtimeLogDeliveryStatus`) until the task is no longer found (deleted).

#### Scenario: Delete existing task
- **WHEN** the resource is destroyed
- **THEN** the system SHALL call `DeleteRealtimeLogDeliveryTask` with ZoneId and TaskId from the composite ID

### Requirement: Unit Tests
The system SHALL provide unit tests in `resource_tc_teo_realtime_log_delivery_test.go` using gomonkey to mock cloud API calls (CreateRealtimeLogDeliveryTask, DescribeRealtimeLogDeliveryTasks, ModifyRealtimeLogDeliveryTask, DeleteRealtimeLogDeliveryTask), testing Create, Read, Update, and Delete operations including the new `custom_expression_fields`, `aws_s3`, and `log_format.record_template` parameters. Tests SHALL NOT use terraform test suite for new scenarios (use gomonkey mock).

#### Scenario: Unit tests pass with new fields
- **WHEN** `go test` is run on the test file
- **THEN** all test cases covering Create/Read/Update/Delete with `custom_expression_fields`, `aws_s3`, and `record_template` SHALL pass

### Requirement: Resource Documentation
The system SHALL provide a markdown documentation file `resource_tc_teo_realtime_log_delivery.md` with description, example usage (including examples of `aws_s3`, `custom_expression_fields`, and `log_format.record_template`), and import section. The `website/docs/` documentation SHALL be generated via `make doc` in the finalize phase (no manual edits to `website/` directory).

#### Scenario: Documentation includes new fields and examples
- **WHEN** the resource is documented
- **THEN** the .md file SHALL include example usage demonstrating `aws_s3` block, `custom_expression_fields` block, and `record_template` within `log_format`, plus import section showing the composite ID format