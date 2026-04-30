## Context

需要实现一个新的 CVM 数据源 `tencentcloud_cvm_account_quota`，用于查询账户配额详情。该数据源调用腾讯云 CVM DescribeAccountQuota API，返回用户在不同地域和可用区的各类配额信息（后付费、预付费、竞价实例、镜像、置放群组等）。

当前 Provider 已有完善的 CVM 服务支持，本次新增为独立的只读数据源，无需修改现有代码，风险较低。

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_cvm_account_quota` 数据源，支持查询账户配额详情
- 支持通过 Filters 按可用区 (zone) 和配额类型 (quota-type) 过滤
- 返回完整的配额数据结构，包括所有配额类型
- 提供完整的文档和测试用例

**Non-Goals:**
- 不支持修改配额（API 本身为只读查询）
- 不实现配额预警或监控功能
- 不实现分页（API 返回完整数据，无需分页）

## Decisions

### 1. Schema 设计

**Filter 参数:**
- `zone` (Optional, Set of String): 可用区列表，如 `["ap-guangzhou-3", "ap-guangzhou-4"]`
- `quota_type` (Optional, String): 配额类型，支持: `PostPaidQuotaSet`, `PrePaidQuotaSet`, `SpotPaidQuotaSet`, `ImageQuotaSet`, `DisasterRecoverGroupQuotaSet`

**输出属性:**
- `app_id` (Int): 用户 AppId
- `account_quota_overview` (List): 配额数据概览
  - `region` (String): 地域
  - `account_quota` (List): 配额详情
    - `post_paid_quota_set` (List): 后付费配额列表
    - `pre_paid_quota_set` (List): 预付费配额列表
    - `spot_paid_quota_set` (List): 竞价实例配额列表
    - `image_quota_set` (List): 镜像配额列表
    - `disaster_recover_group_quota_set` (List): 置放群组配额列表
- `result_output_file` (Optional, String): 结果输出文件

**理由:** 遵循 Provider 现有的 Schema 命名规范（snake_case），Filter 参数使用 Set 类型支持多值过滤。

### 2. API 调用方式

使用 `tencentcloud-sdk-go/tencentcloud/cvm/v20170312` 包的 `DescribeAccountQuota` 方法。

```go
request := cvm.NewDescribeAccountQuotaRequest()
// 设置 Filters
response, err := cvmService.client.UseCvmClient().DescribeAccountQuota(request)
```

**理由:** 遵循现有 CVM 数据源的实现模式，复用 cvmService 的客户端连接。

### 3. 数据转换

API 返回的 `AccountQuotaOverview` 是复杂的嵌套结构，需要递归转换为 Terraform Schema 支持的 map 格式。

**转换函数:**
- `flattenAccountQuotaOverview()`: 转换 AccountQuotaOverview
- `flattenAccountQuota()`: 转换 AccountQuota
- `flattenPostPaidQuotaSet()`, `flattenPrePaidQuotaSet()` 等: 转换各类配额集合

**理由:** 保持与其他 CVM 数据源一致的 flatten 模式，便于维护。

### 4. 错误处理

使用标准的错误处理模式:
```go
defer tccommon.LogElapsed("data_source.tencentcloud_cvm_account_quota.read")()
if err != nil {
    return diag.FromErr(err)
}
```

**理由:** 遵循 Provider 的错误处理规范。

### 5. 测试策略

**验收测试:**
- 基本查询测试（不带 filter）
- 按可用区过滤测试
- 按配额类型过滤测试

**测试数据:** 使用真实的云 API 调用，需要有效的 TENCENTCLOUD_SECRET_ID 和 TENCENTCLOUD_SECRET_KEY。

**理由:** 数据源为只读操作，测试不会产生资源消耗，可安全地对真实 API 进行测试。

## Risks / Trade-offs

**[风险] API 返回数据结构可能随版本变化** → 使用 tencentcloud-sdk-go 的类型定义，随 SDK 更新自动适配

**[风险] 不同地域的配额数据结构可能不一致** → 使用 Computed + Optional 的 Schema 设计，兼容字段缺失情况

**[风险] API 频率限制 (20次/秒)** → 文档中说明查询频率限制，建议用户合理使用

**[取舍] Filter 参数设计为独立字段 vs 使用 Filter 结构** → 选择独立字段（zone, quota_type）以提供更清晰的用户体验，牺牲了与 API 参数的一致性

**[取舍] 不实现配额监控功能** → 专注于查询功能，监控可通过 Terraform 外部工具实现
