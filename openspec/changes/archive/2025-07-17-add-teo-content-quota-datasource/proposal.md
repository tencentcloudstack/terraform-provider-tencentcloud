## Why

目前 Terraform Provider for TencentCloud 中没有数据源可以查询 TEO（边缘安全加速平台）的内容管理接口配额信息（包括缓存刷新配额和预热配额）。用户在编写 Terraform 配置时，无法获取当前账号的刷新和预热配额使用情况，不利于资源规划和自动化管理。

## What Changes

- 新增数据源 `tencentcloud_teo_content_quota`，用于查询 TEO 内容管理接口配额
- 数据源调用 `DescribeContentQuota` API，入参为 `zone_id`，出参包括 `purge_quota`（刷新配额列表）和 `prefetch_quota`（预热配额列表）
- 在 `provider.go` 和 `provider.md` 中注册新数据源

## Capabilities

### New Capabilities
- `teo-content-quota-datasource`: 新增 TEO 内容管理配额数据源，支持按站点 ID 查询缓存刷新和预热的配额信息

### Modified Capabilities

## Impact

- 新增文件：`tencentcloud/services/teo/data_source_tc_teo_content_quota.go`
- 新增文件：`tencentcloud/services/teo/data_source_tc_teo_content_quota_test.go`
- 新增文件：`tencentcloud/services/teo/data_source_tc_teo_content_quota.md`
- 修改文件：`tencentcloud/provider.go`（注册数据源）
- 修改文件：`tencentcloud/provider.md`（添加数据源文档）
- 依赖：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` 中的 `DescribeContentQuota` 接口
