## Context

Terraform Provider for TencentCloud 需要新增一个 TEO（边缘安全加速平台）内容管理配额数据源。当前 Provider 中没有对应数据源可供用户查询缓存刷新和预热的配额信息。

云 API `DescribeContentQuota` 位于 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` 包中，接口入参仅 `ZoneId`，出参包含 `PurgeQuota` 和 `PrefetchQuota`，类型均为 `[]*Quota`。`Quota` 结构体包含 `Batch`（单次批量提交配额上限）、`Daily`（每日提交配额上限）、`DailyAvailable`（每日剩余配额）、`Type`（刷新预热缓存类型）四个字段。

该接口为同步接口，无需轮询。该接口无分页参数。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_teo_content_quota` 数据源，调用 `DescribeContentQuota` API 查询内容管理配额
- 数据源 schema 包含入参 `zone_id` 和出参 `purge_quota`、`prefetch_quota`
- `purge_quota` 和 `prefetch_quota` 为列表类型，每个元素包含 `batch`、`daily`、`daily_available`、`type` 四个字段
- 在 `provider.go` 和 `provider.md` 中注册新数据源
- 编写单元测试使用 gomonkey mock 云 API

**Non-Goals:**
- 不创建 resource（仅数据源）
- 不修改现有资源或数据源的 schema
- 不支持分页（该接口本身无分页）

## Decisions

1. **数据源 Read 函数中直接调用 API**：由于该接口入参简单（仅 ZoneId），无需通过 service 层封装，直接在 Read 函数中构造请求并调用 API，与现有 teo 数据源风格保持一致。

2. **Quota 字段使用 TypeList + schema.Resource 嵌套**：`PurgeQuota` 和 `PrefetchQuota` 为 `[]*Quota` 类型，在 Terraform schema 中使用 `TypeList` 搭配嵌套 `schema.Resource` 表示，每个元素包含 `batch`（TypeInt）、`daily`（TypeInt）、`daily_available`（TypeInt）、`type`（TypeString）。

3. **数据源 ID 使用 helper.BuildToken()**：与 igtm_instance_list 等数据源保持一致，使用 `helper.BuildToken()` 生成唯一 ID。

4. **API 调用添加 retry**：使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包装 API 调用，失败时使用 `tccommon.RetryError()` 包装错误。

5. **nil 检查**：在设置字段前检查 Response 中的字段是否为 nil，避免 panic。

6. **result_output_file**：保留标准的 `result_output_file` 可选参数，用于保存结果。

## Risks / Trade-offs

- [API 字段可能返回 null] → 在 Read 函数中对 PurgeQuota 和 PrefetchQuota 进行 nil 检查，仅当非 nil 时才设置到 state 中
- [Quota 结构体字段可能返回 null] → 对每个 Quota 元素的 Batch、Daily、DailyAvailable、Type 字段进行 nil 检查
