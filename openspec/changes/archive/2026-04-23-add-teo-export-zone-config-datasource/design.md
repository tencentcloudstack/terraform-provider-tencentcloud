## Context

TEO (TencentCloud EdgeOne) 是腾讯云的边缘安全加速平台。当前 Terraform Provider 中已有多个 TEO 数据源（如 `tencentcloud_teo_origin_acl`、`tencentcloud_teo_zones` 等），但缺少站点配置导出的数据源。`ExportZoneConfig` 是 TEO SDK 中已有的同步接口，支持按配置类型导出站点配置，返回 JSON 格式的配置内容字符串。

本变更为纯新增数据源，不涉及任何已有资源或数据源的修改。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_teo_export_zone_config` 数据源，支持通过 `zone_id` 和 `types` 参数调用 `ExportZoneConfig` API 导出站点配置
- 返回 `content` 字段，包含 JSON 格式编码的站点配置内容
- 在 provider.go 和 provider.md 中正确注册该数据源
- 生成对应的 .md 文档示例文件
- 编写基于 gomonkey mock 的单元测试

**Non-Goals:**
- 不实现配置导入功能（`ImportZoneConfig` 属于另一资源）
- 不解析 `content` 字段的 JSON 内容为结构化 schema 字段（content 为原始 JSON 字符串）
- 不修改任何已有资源或数据源

## Decisions

### 1. 数据源 Schema 设计
- `zone_id`: Required, TypeString - 站点 ID，必填参数
- `types`: Optional, TypeList of TypeString - 导出配置类型列表，不填则导出所有类型配置。支持值：`L7AccelerationConfig`、`WebSecurity`
- `content`: Computed, TypeString - 导出的配置内容，JSON 格式字符串
- `result_output_file`: Optional, TypeString - 结果输出文件路径（标准模式）

**理由**: 遵循 provider 中已有数据源的模式（如 `tencentcloud_teo_origin_acl`），保持一致的用户体验。

### 2. 资源 ID 策略
- 使用 `helper.BuildToken()` 生成随机 token 作为数据源的 ID

**理由**: 数据源每次读取都应视为独立查询，使用随机 token 避免 ID 冲突，同时简化实现。

### 3. API 调用方式
- 在 service 层 `ExportZoneConfigByFilter` 中调用 `ExportZoneConfig` API，service 层已使用 `resource.Retry` 包装以 `tccommon.ReadRetryTimeout` 为超时时间
- data source 的 Read 函数直接调用 service 方法，不再额外包装 retry，避免重复 retry
- 不需要轮询，因为 `ExportZoneConfig` 是同步接口

**理由**: retry 逻辑统一放在 service 层，data source 层无需重复包装。

### 4. 测试策略
- 使用 gomonkey mock 方式进行单元测试，不使用 Terraform 测试套件

**理由**: 遵循项目规范，新增资源使用 mock 方式进行业务代码逻辑的单元测试。

## Risks / Trade-offs

- [Content 字段可能较大] → 导出所有类型配置时 content 可能包含大量 JSON 数据，建议用户通过 `types` 参数指定需要的配置类型以控制响应大小。此为 API 自身特性，无需特殊处理。
- [types 参数值可能扩展] → API 文档标注后续会支持更多导出类型，schema 中 `types` 字段不做值限制，由 API 端校验。
