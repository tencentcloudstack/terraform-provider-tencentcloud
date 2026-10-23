## Why

`tencentcloud_teo_realtime_log_delivery` 资源目前仅覆盖 CLS、自定义 HTTP(S)（custom_endpoint）、S3 兼容（s3）三种投递目标，缺少对云 API 已支持的 `aws_s3`（AWS 原生 S3，IAM 角色授权免长期密钥）投递类型、自定义表达式字段（CustomExpressionFields）以及日志自定义输出模板（LogFormat.RecordTemplate）的支持。腾讯云 TEO SDK（`teov20220901`）的 `CreateRealtimeLogDeliveryTask` / `ModifyRealtimeLogDeliveryTask` / `DescribeRealtimeLogDeliveryTasks` 三个接口已提供这些参数，需要将其映射到 Terraform 资源 schema，使 gusty 用户能够通过 Terraform 声明式管理这些能力。

## What Changes

- 在 `tencentcloud_teo_realtime_log_delivery` 资源 schema 中新增顶层块 `aws_s3`（TypeList, MaxItems=1, Optional+Computed），用于配置 AWS 原生 S3 投递目标（当 `task_type = aws_s3` 时必填），包含字段：`endpoint`、`bucket_name`、`path`、`role_arn`、`compress_type`。
- 在资源 schema 中新增顶层块 `custom_expression_fields`（TypeList, Optional+Computed），用于自定义表达式日志字段，每个元素包含字段：`name`、`expression`、`enabled`。
- 在已有 `log_format` 块中新增子字段 `record_template`（Optional, TypeString），对应云 API `LogFormat.RecordTemplate`，仅当 `format_type = template` 时生效。
- Create/Update 逻辑适配：构造请求时填充 `CustomExpressionFields`、`AWSS3`，以及在 `LogFormat` 中填充 `RecordTemplate`。
- Read 逻辑适配：从 `DescribeRealtimeLogDeliveryTasks` 响应中回填上述新增字段到 state。
- 单元测试补充（使用 gomonkey mock 云 API）覆盖新增 `custom_expression_fields`、`aws_s3`、`log_format.record_template` 的 Create/Read/Update 场景。
- 通过 `make doc` 在收尾阶段同步更新 `website/docs/` 文档（禁止手改 website 目录）。

非破坏性：全部为新增 Optional（或 Optional+Computed）字段，不影响已有 TF 配置与 state。

## Capabilities

### New Capabilities
- `teo-realtime-log-delivery-resource`: 管理 TEO 实时日志投递任务（`tencentcloud_teo_realtime_log_delivery`）的 CRUD 能力，本次新增 `aws_s3` 投递目标块、`custom_expression_fields` 自定义表达式字段块、以及 `log_format.record_template` 日志模板字段。

### Modified Capabilities
<!-- 无已有 capability 需要修改，该资源此前未在 openspec/specs 中登记 -->

## Impact

- 代码：
  - `tencentcloud/services/teo/resource_tc_teo_realtime_log_delivery.go`（schema、Create、Read、Update 逻辑）
  - `tencentcloud/services/teo/resource_tc_teo_realtime_log_delivery_test.go`（新增单测用例）
  - `tencentcloud/services/teo/resource_tc_teo_realtime_log_delivery.md`（示例补充与字段说明）
- 依赖：使用已 vendored 的 `tencentcloud-sdk-go` 中 `teov20220901` 包的 `AWSS3`、`CustomExpressionField`、`LogFormat.RecordTemplate` 字段，无需变更 vendor。
- 向后兼容：新增 Optional / Optional+Computed 字段，不影响已有 state 与 TF 配置。
- 文档：需要同步更新 website docs（由收尾阶段 `make doc` 自动生成）。