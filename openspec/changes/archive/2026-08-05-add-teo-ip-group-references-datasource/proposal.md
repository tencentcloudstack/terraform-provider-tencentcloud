## Why

当前 Terraform Provider 中缺少查询 TEO（EdgeOne）IP 分组被引用信息的数据源。用户需要通过 `DescribeIPGroupReferences` 接口查询指定站点下某个 IP 组被哪些安全策略、DDoS 防护等实体引用，以便在 Terraform 配置中引用这些数据，用于依赖分析、影响面评估和配置审计。

## What Changes

- 新增数据源 `tencentcloud_teo_ip_group_references`，封装 `DescribeIPGroupReferences` API
- 支持按 `zone_id` 和 `group_id` 查询指定 IP 组的引用信息列表
- 内部自动分页，获取所有引用数据
- 在 `provider.go` 和 `provider.md` 中注册该数据源

## Capabilities

### New Capabilities
- `teo-ip-group-references-datasource`: 提供查询 TEO IP 分组引用信息的数据源能力，支持按站点 ID 和 IP 组 ID 查询该 IP 组被哪些安全策略、DDoS 防护等实体引用

### Modified Capabilities
<!-- 无需修改现有能力 -->

## Impact

- 新增文件: `tencentcloud/services/teo/data_source_tc_teo_ip_group_references.go`
- 新增文件: `tencentcloud/services/teo/data_source_tc_teo_ip_group_references_test.go`
- 新增文件: `tencentcloud/services/teo/data_source_tc_teo_ip_group_references.md`
- 修改文件: `tencentcloud/provider.go`（注册数据源）
- 修改文件: `tencentcloud/provider.md`（添加数据源文档条目）
- 依赖: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`
