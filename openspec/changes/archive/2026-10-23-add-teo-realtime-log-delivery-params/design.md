## Context

`tencentcloud_teo_realtime_log_delivery` 是 TEO（边缘安全加速平台）实时日志投递任务的 Terraform 资源（RESOURCE_KIND_GENERAL），封装了 `CreateRealtimeLogDeliveryTask` / `DescribeRealtimeLogDeliveryTasks` / `ModifyRealtimeLogDeliveryTask` / `DeleteRealtimeLogDeliveryTask` 四个云 API，管理资源整个生命周期。复合 ID 格式为 `zoneId#taskId`（分隔符 `tccommon.FILED_SP`），支持 import。

当前状态（已实现的能力）：
- schema 已覆盖 `zone_id`、`task_name`、`task_type`、`entity_list`、`log_type`、`area`、`fields`、`custom_fields`、`delivery_conditions`、`sample`、`log_format`（含 format_type/batch_prefix/batch_suffix/record_prefix/record_suffix/record_delimiter/field_delimiter）、`cls`、`custom_endpoint`、`s3`、`delivery_status`、`task_id`。
- Create：构造 `CreateRealtimeLogDeliveryTaskRequest`，调用后取 `response.Response.TaskId`，`d.SetId(zoneId#taskId)`，然后通过 `checkRealtimeLogDeliveryStatus` 轮询等待 `enabled`。
- Read：调用 `service.DescribeTeoRealtimeLogDeliveryById(ctx, zoneId, taskId)`（内部封装 `DescribeRealtimeLogDeliveryTasks`，按 `task-id` 过滤），回填 state。
- Update：`immutableArgs = [task_type, log_type, area, cls]`；`mutableArgs = [task_name, delivery_status, entity_list, fields, custom_fields, delivery_conditions, sample, log_format, custom_endpoint, s3]`，任一变化即调用 `ModifyRealtimeLogDeliveryTask`。
- Delete：调用 `DeleteRealtimeLogDeliveryTask`，轮询等待 `deleted`。

本次新增需求（来源于云 API 已支持但 Terraform 未映射的参数）：
1. **aws_s3 投递目标块**：云 API 新增 `TaskType = aws_s3`（推送到 AWS 原生 S3，使用 IAM 角色授权免长期密钥），对应请求/响应结构 `AWSS3`（字段：Endpoint / BucketName / Path / RoleARN / CompressType）。
2. **custom_expression_fields 自定义表达式字段块**：云 API 新增 `CustomExpressionFields []*CustomExpressionField`（字段：Name / Expression / Enabled），用于通过自定义表达式实现个性化日志字段推送。
3. **log_format.record_template 日志模板子字段**：云 API 的 `LogFormat` 结构新增 `RecordTemplate` 字段，仅当 `FormatType = template` 时生效，用于自定义日志输出模板。

约束：
- 必须保持向后兼容：全部为新增 Optional / Optional+Computed 字段，不能破坏已有 TF 配置与 state。
- 复合 ID 不变（仍是 `zoneId#taskId`）。
- aws_s3 与 s3（S3 兼容存储）是不同的投递目标，互斥：当 `task_type = aws_s3` 时填 `aws_s3`；`task_type = s3` 时填 `s3`。
- `custom_expression_fields` 与 `custom_fields` 并列，两者都可配置自定义日志字段，但取值方式不同（表达式 vs HTTP 头/Cookie 提取）；若与 `Fields` / `CustomFields` 同名，以 `CustomExpressionFields` 为准。
- `log_format.record_template` 仅在 `format_type = template` 时生效，为已有 `log_format` 块追加子字段（不新增块）。

## Goals / Non-Goals

**Goals:**
- 将云 API 已支持的 `AWSS3`、`CustomExpressionFields`、`LogFormat.RecordTemplate` 映射到 Terraform schema，使资源支持 `aws_s3` 投递类型与自定义表达式字段、日志模板。
- 保持向后兼容：新增字段全部为 Optional / Optional+Computed，存量配置升级无 plan drift。
- Create / Update / Read 全链路适配新参数；Update 中将 `aws_s3` 与 `custom_expression_fields` 纳入 `mutableArgs`（可变更），`log_format` 已在 mutableArgs 内（新增子字段自动跟随）。
- 通过 gomonkey 单测覆盖新增参数的 Create/Read/Update 路径。

**Non-Goals:**
- 不改变现有 `cls` / `custom_endpoint` / `s3` 投递类型的逻辑。
- 不改变复合 ID 格式与 import 行为。
- 不新增独立资源（不拆分 `aws_s3` 为单独资源）。
- 不暴露 `DescribeRealtimeLogDeliveryTasks` 的 `Offset` / `Limit` / `Filters.Name` / `Filters.Fuzzy` 给用户 schema（这些是服务层内部查询参数，数据源分页按项目规范内部实现自动分页，资源 Read 不暴露）。
- 不修改 `task_type` 的取值校验（服务端返回错误即可，保持与现有字段一致）。

## Decisions

### Decision 1: aws_s3 作为新增顶层块（TypeList, MaxItems=1, Optional+Computed），而非平铺字段

**选择**：参照现有 `s3`、`cls`、`custom_endpoint` 块的既有模式，新增 `aws_s3` 顶层块，`MaxItems: 1`，`Optional + Computed`，内部 Elem 为 `schema.Resource`，子字段平铺（`endpoint` / `bucket_name` / `path` / `role_arn` / `compress_type`）。

**备选**：将 aws_s3 字段平铺到资源顶层（如 `aws_s3_endpoint` / `aws_s3_bucket_name`）。

**理由**：
- 与现有 `s3` / `cls` / `custom_endpoint` 块结构完全一致，代码风格统一，Read 用 `helper.InterfacesHeadMap` 解析，Create/Update 用 `teo.AWSS3{}` 构造。
- aws_s3 是一个逻辑聚合的配置对象，平铺会导致顶层字段过多且语义割裂。
- `MaxItems: 1` 保证只有一份配置。

### Decision 2: custom_expression_fields 作为新增顶层块（TypeList, Optional+Computed），元素结构与 custom_fields 一致风格

**选择**：新增 `custom_expression_fields` 顶层块，`TypeList`，`Optional + Computed`，`Elem` 为 `schema.Resource`，子字段：`name`（Required, TypeString）、`expression`（Required, TypeString）、`enabled`（Optional, TypeBool）。

**备选**：合并进 `custom_fields` 用一个类型标识区分。

**理由**：
- 云 API 中 `CustomFields`（`[]*CustomField`，字段 Name/Value/Enabled）与 `CustomExpressionFields`（`[]*CustomExpressionField`，字段 Name/Expression/Enabled）是两个独立的请求字段，语义不同（提取 vs 表达式），合并会增加耦合与歧义。
- 独立块与云 API 一一对应，Create/Update/Read 映射直观。

### Decision 3: log_format 新增 record_template 子字段（不新增块）

**选择**：在已有 `log_format` 块的 `Elem.Schema` 中追加 `record_template`（Optional, TypeString），对应 `LogFormat.RecordTemplate`。

**理由**：
- 云 API 的 `RecordTemplate` 是 `LogFormat` 结构的一个字段，与 `format_type` / `batch_prefix` 等同级，自然归入已有 `log_format` 块。
- Create/Read/Update 中 `log_format` 的解析逻辑已在 `helper.InterfacesHeadMap(d, "log_format")` 内逐字段填充，只需新增一行 `logFormat.RecordTemplate = ...`，改动最小。

### Decision 4: Read 回填新字段时遵循 nil 检查与展平规则

**选择**：
- `aws_s3`：参照 `s3` 的回填模式（`respData.AWSS3 != nil` 时构造 map，`[]interface{}{awsS3Map}` set 到 `d`）。
- `custom_expression_fields`：参照 `custom_fields` 的回填模式（遍历 `respData.CustomExpressionFields`，逐元素构造 map 加入 list）。
- `log_format.record_template`：在已有 `log_format` 回填分支内追加 `if respData.LogFormat.RecordTemplate != nil { logFormatMap["record_template"] = ... }`。
- 全部在调用 `d.Set` 前对 Response 字段做 nil 检查，符合代码生成要求。

### Decision 5: Update 中将 aws_s3、custom_expression_fields 纳入 mutableArgs

**选择**：在 `mutableArgs` 数组中追加 `"aws_s3"`、`"custom_expression_fields"`；`log_format` 已在其中（新增子字段自动触发）。`immutableArgs` 不变（`task_type` 仍是 ForceNew 等价语义——变更 task_type 会报错，因 aws_s3 与 s3 是不同 task_type，切换需重建资源）。

**理由**：
- 云 API `ModifyRealtimeLogDeliveryTask` 已支持 `AWSS3` 与 `CustomExpressionFields` 入参，可作为可变字段就地更新。
- `task_type` 在 immutableArgs 中，意味着从 `s3` 切到 `aws_s3` 会报错（需 destroy + recreate），符合云 API 语义（投递目标类型变更本质是重建任务）。

### Decision 6: 不暴露 DescribeRealtimeLogDeliveryTasks 的 Offset/Limit/Filters 给资源 schema

**选择**：本次需求虽然列出 `DescribeRealtimeLogDeliveryTasks` 的 `Offset` / `Limit` / `Filters.Name` / `Filters.Fuzzy` 为"新增入参"，但根据项目规范（数据源分页不暴露 limit/offset 给用户，资源 Read 内部实现）以及这些是查询过滤参数而非资源属性，**不在资源 schema 中新增** 这些字段。服务层 `DescribeTeoRealtimeLogDeliveryById` 继续按 `task-id` 精确过滤获取单个资源（已实现），无需改动分页逻辑。

**理由**：
- 资源 Read 只需按主键（zoneId + taskId）获取单个资源，不需要用户配置分页/过滤参数。
- 暴露分页参数会破坏资源 schema 的语义（资源不应有查询分页字段）。
- 这些参数属于数据源（datasource）场景，本次是 RESOURCE_KIND_GENERAL 资源，不在范围内。

## Risks / Trade-offs

- **Risk**：`aws_s3` 与 `s3` 互斥但 schema 层未做 `ConflictsWith` 强约束（与现有 `cls`/`custom_endpoint`/`s3` 风格一致，均未声明互斥）→ **Mitigation**：保持与现有风格一致，由云 API 校验（`task_type` 决定哪个块必填），服务端返回错误即可；不在 schema 层加 ConflictsWith 以避免与未来 API 扩展耦合。
- **Risk**：存量 state 中无 `aws_s3` / `custom_expression_fields` / `record_template`，首次 plan 可能因 Computed 字段未回填产生 diff → **Mitigation**：设为 `Optional + Computed`，未配置时由 Read 回填（云端返回 null 时跳过 set），符合 provider 标准模式，无 plan drift。
- **Risk**：`custom_expression_fields` 与 `custom_fields` 同名时以表达式为准，用户可能混淆 → **Mitigation**：在字段 Description 中注明优先级规则（与云 API 文档一致）。
- **Trade-off**：`log_format.record_template` 仅 `format_type = template` 生效，但 schema 层不做条件校验 → 可接受，服务端校验，与现有 `log_format` 子字段一致。

## Migration Plan

- 纯加法变更（新增 Optional / Optional+Computed 字段），无 state 迁移需求。
- 存量资源：`terraform plan` 对未在 HCL 配置新字段的资源不产生 diff（Optional+Computed 行为，Read 回填）。
- 文档更新：在 `resource_tc_teo_realtime_log_delivery.md` 中补充 `aws_s3`、`custom_expression_fields`、`log_format.record_template` 字段说明与示例（收尾阶段通过 `make doc` 生成 website docs）。
- 回滚：移除新增 schema 字段与对应 Create/Read/Update 分支即可，state 中残留的新字段值在回滚后会被忽略，无副作用。

## 参考信息（实施阶段所需，避免反复回读 vendor 大文件）

以下摘录自 `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901/models.go`（按行区间），供实施阶段直接引用。

### A. AWSS3 结构（models.go 第 68-83 行）

```go
type AWSS3 struct {
	// AWS S3 存储桶对应的服务端点域名，取值可参考 AWS 服务端点。例如：s3.us-west-2.amazonaws.com 。
	Endpoint *string `json:"Endpoint,omitnil,omitempty" name:"Endpoint"`

	// AWS S3 存储桶名称。可在 AWS S3 控制台 Buckets（存储桶）列表中查看对应桶名。
	// 长度须为 3–63 个字符，只能包含小写字母、数字、句点（.）和连字符（-），且必须以字母或数字开头和结尾。
	BucketName *string `json:"BucketName,omitnil,omitempty" name:"BucketName"`

	// 日志文件在 AWS S3 桶中的存储目录路径，以 / 分隔可创建多级目录，且不能包含 ./.. 等相对路径段，路径总长度不超过 1024 个字节。
	// 当目录路径存在时直接使用该路径，当目录路径不存在时会自动创建。
	// 示例：eo-logs/l7-access-logs 。
	Path *string `json:"Path,omitnil,omitempty" name:"Path"`

	// 授权 EdgeOne 访问 AWS S3 桶的 IAM 角色的 ARN。
	// 例如：arn:aws:iam::123456789012:role/eo-logs 。
	RoleARN *string `json:"RoleARN,omitnil,omitempty" name:"RoleARN"`

	// 数据压缩类型，取值有: gzip：gzip压缩。不填表示不启用压缩。
	CompressType *string `json:"CompressType,omitnil,omitempty" name:"CompressType"`
}
```

### B. CustomExpressionField 结构（models.go 第 6140-6149 行）

```go
type CustomExpressionField struct {
	// 自定义日志字段名称。可输入1-100个字符，允许的字符为字母、数字、_，仅能以字母开头，该名称不能重复。
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// 自定义日志字段的取值表达式，表达式长度上限 4KB。
	Expression *string `json:"Expression,omitnil,omitempty" name:"Expression"`

	// 是否投递该字段，不填表示不投递此字段。
	Enabled *bool `json:"Enabled,omitnil,omitempty" name:"Enabled"`
}
```

### C. LogFormat 结构（models.go 第 18977-19001 行）—— 重点新增 RecordTemplate 字段

```go
type LogFormat struct {
	// 日志输出格式，取值有：
	// json：使用预设日志输出格式 JSON Lines；
	// csv：使用预设日志输出格式 csv；
	// template：使用用户自定义输出模板，需配合 RecordTemplate 字段使用。
	FormatType *string `json:"FormatType,omitnil,omitempty" name:"FormatType"`

	BatchPrefix *string `json:"BatchPrefix,omitnil,omitempty" name:"BatchPrefix"`
	BatchSuffix *string `json:"BatchSuffix,omitnil,omitempty" name:"BatchSuffix"`
	RecordPrefix *string `json:"RecordPrefix,omitnil,omitempty" name:"RecordPrefix"`
	RecordSuffix *string `json:"RecordSuffix,omitnil,omitempty" name:"RecordSuffix"`
	RecordDelimiter *string `json:"RecordDelimiter,omitnil,omitempty" name:"RecordDelimiter"`

	// 日志模板，单条日志的输出模板，长度限制 4KB，仅当 FormatType = template 生效。
	// 支持对配置的推送字段按照模板进行自定义排版和拼接。
	RecordTemplate *string `json:"RecordTemplate,omitnil,omitempty" name:"RecordTemplate"`

	FieldDelimiter *string `json:"FieldDelimiter,omitnil,omitempty" name:"FieldDelimiter"`
}
```

### D. CreateRealtimeLogDeliveryTaskRequest 相关字段（models.go 第 5219-5269 行）

新增入参映射位置：
```go
// 第 5247 行：自定义表达式字段列表
CustomExpressionFields []*CustomExpressionField `json:"CustomExpressionFields,omitnil,omitempty" name:"CustomExpressionFields"`

// 第 5256 行：日志输出格式（LogFormat 结构内含 RecordTemplate）
LogFormat *LogFormat `json:"LogFormat,omitnil,omitempty" name:"LogFormat"`

// 第 5268 行：AWS S3 配置
AWSS3 *AWSS3 `json:"AWSS3,omitnil,omitempty" name:"AWSS3"`
```

### E. ModifyRealtimeLogDeliveryTaskRequest 相关字段（models.go 第 22264-22308 行）

新增入参映射位置：
```go
// 第 22289 行：自定义表达式字段列表
CustomExpressionFields []*CustomExpressionField `json:"CustomExpressionFields,omitnil,omitempty" name:"CustomExpressionFields"`

// 第 22298 行：日志输出格式
LogFormat *LogFormat `json:"LogFormat,omitnil,omitempty" name:"LogFormat"`

// 第 22307 行：AWS S3 配置
AWSS3 *AWSS3 `json:"AWSS3,omitnil,omitempty" name:"AWSS3"`
```

### F. DescribeRealtimeLogDeliveryTasksResponse / RealtimeLogDeliveryTask 出参（models.go 第 13646-13655、24838-24900 行）

Response 出参：
```go
type DescribeRealtimeLogDeliveryTasksResponseParams struct {
	// 符合查询条件的实时日志投递任务个数。
	TotalCount *uint64 `json:"TotalCount,omitnil,omitempty" name:"TotalCount"`
	// 符合查询条件的所有实时日志投递任务列表。
	RealtimeLogDeliveryTasks []*RealtimeLogDeliveryTask `json:"RealtimeLogDeliveryTasks,omitnil,omitempty" name:"RealtimeLogDeliveryTasks"`
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}
```

`RealtimeLogDeliveryTask` 出参（回填所需）：
```go
// 第 24867 行：自定义表达式字段列表（出参）
CustomExpressionFields []*CustomExpressionField `json:"CustomExpressionFields,omitnil,omitempty" name:"CustomExpressionFields"`

// 第 24877 行：日志输出格式（LogFormat 内含 RecordTemplate，注意：此字段可能返回 null）
LogFormat *LogFormat `json:"LogFormat,omitnil,omitempty" name:"LogFormat"`

// 第 24893 行：AWS S3 配置（注意：此字段可能返回 null）
AWSS3 *AWSS3 `json:"AWSS3,omitnil,omitempty" name:"AWSS3"`
```

### G. AdvancedFilter 结构（models.go 第 402-411 行）—— DescribeRealtimeLogDeliveryTasks 的 Filters 元素

```go
type AdvancedFilter struct {
	// 需要过滤的字段。
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`
	// 字段的过滤值。
	Values []*string `json:"Values,omitnil,omitempty" name:"Values"`
	// 是否启用模糊查询。
	Fuzzy *bool `json:"Fuzzy,omitnil,omitempty" name:"Fuzzy"`
}
```

注：`AdvancedFilter.Name` / `AdvancedFilter.Fuzzy` 属于查询过滤参数，按 Decision 6 不映射到资源 schema，仅服务层内部使用。

### H. 现有资源 schema 中与新增参数关联的结构说明

- `log_format` 块（resource_tc_teo_realtime_log_delivery.go 第 147-192 行）：TypeList, Optional+Computed, MaxItems=1，Elem.Schema 已含 `format_type`(Required) / `batch_prefix` / `batch_suffix` / `record_prefix` / `record_suffix` / `record_delimiter` / `field_delimiter`。本次在此 Elem.Schema 中追加 `record_template`（Optional, TypeString）。
- `custom_fields` 块（第 77-101 行）：TypeList, Optional+Computed，Elem.Schema 含 `name`/`value`/`enabled`。本次新增的 `custom_expression_fields` 块参照此结构，但子字段为 `name`/`expression`/`enabled`（无 value，有 expression）。
- `s3` 块（第 277-318 行）：TypeList, Optional+Computed, MaxItems=1，Elem.Schema 含 `endpoint`/`region`/`bucket`/`access_id`/`access_key`/`compress_type`。本次新增的 `aws_s3` 块参照此结构，子字段为 `endpoint`/`bucket_name`/`path`/`role_arn`/`compress_type`（注意：aws_s3 无 region/bucket/access_id/access_key，用 BucketName/Path/RoleARN 替代，使用 IAM 角色授权免长期密钥）。
- Create 中 `s3` 的填充模式（第 507-528 行）与 Read 中 `s3` 的回填模式（第 778-806 行）是 aws_s3 块的直接参照模板。
- Create 中 `custom_fields` 的填充模式（第 387-402 行）与 Read 中 `custom_fields` 的回填模式（第 622-643 行）是 custom_expression_fields 块的直接参照模板。
- Update 的 `mutableArgs`（第 833 行）需追加 `"aws_s3"`、`"custom_expression_fields"`；`immutableArgs`（第 819 行）不变。

## Open Questions

- 无