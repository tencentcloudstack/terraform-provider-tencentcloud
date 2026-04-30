## Why

Terraform Provider for TencentCloud 目前缺少对 TEO（EdgeOne）内容管理配额的查询支持。用户在使用 Terraform 管理 TEO 站点时，无法通过数据源获取缓存刷新（Purge）和预热（Prefetch）的配额信息，包括每日配额上限、每日剩余配额和单次批量提交上限。新增 `tencentcloud_teo_content_quota` 数据源可填补这一空白，使用户能够在 Terraform 配置中查询并引用这些配额信息。

## What Changes

- 新增数据源 `tencentcloud_teo_content_quota`，用于查询 TEO 站点的内容管理接口配额
- 支持通过 `zone_id` 入参查询指定站点的刷新配额（PurgeQuota）和预热配额（PrefetchQuota）
- 每种配额包含以下信息：类型（Type）、单次批量提交上限（Batch）、每日提交上限（Daily）、每日剩余可用配额（DailyAvailable）
- 在 `provider.go` 和 `provider.md` 中注册新数据源
- 生成对应的 `.md` 文档

## Capabilities

### New Capabilities
- `teo-content-quota-datasource`: 新增 TEO 内容管理配额数据源，支持查询指定站点的缓存刷新和预热配额信息

### Modified Capabilities

## Impact

- 新增文件：`tencentcloud/services/teo/data_source_tc_teo_content_quota.go`
- 新增文件：`tencentcloud/services/teo/data_source_tc_teo_content_quota_test.go`
- 修改文件：`tencentcloud/services/teo/service_tencentcloud_teo.go`（新增 DescribeContentQuotaByFilter 方法）
- 修改文件：`tencentcloud/provider.go`（注册新数据源）
- 修改文件：`tencentcloud/provider.md`（添加新数据源文档条目）
- 新增文件：`tencentcloud/services/teo/data_source_tc_teo_content_quota.md`
- 依赖云 API：`teo/v20220901.DescribeContentQuota`
