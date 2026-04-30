## Why

TencentCloud EdgeOne (TEO) 提供了预热回源限速（Prefetch Origin Limit）功能，允许用户对预热回源带宽进行限速控制，避免预热任务占用过多源站带宽。目前 Terraform Provider 中缺少对该配置的管理能力，用户无法通过 Terraform 对预热回源限速进行配置和管理。

## What Changes

- 新增 Terraform CONFIG 类型资源 `tencentcloud_teo_prefetch_origin_limit`，支持预热回源限速配置的管理
- Create: 调用 `ModifyPrefetchOriginLimit` 接口创建/启用限速配置（Enabled=on）
- Read: 调用 `DescribePrefetchOriginLimit` 接口读取限速配置
- Update: 调用 `ModifyPrefetchOriginLimit` 接口更新限速配置
- Delete: 调用 `ModifyPrefetchOriginLimit` 接口设置 Enabled=off 删除限速配置
- 使用 `zone_id + domain_name + area` 作为联合 ID（FILED_SP 分隔）
- 在 `provider.go` 和 `provider.md` 中注册该资源
- 生成对应的 `.md` 文档

## Capabilities

### New Capabilities
- `teo-prefetch-origin-limit-config`: 管理 TEO 预热回源限速配置，包括设置限速带宽、加速区域及启停控制

### Modified Capabilities
（无）

## Impact

- 新增文件: `tencentcloud/services/teo/resource_tc_teo_prefetch_origin_limit_config.go`
- 新增文件: `tencentcloud/services/teo/resource_tc_teo_prefetch_origin_limit_config_test.go`
- 新增文件: `tencentcloud/services/teo/resource_tc_teo_prefetch_origin_limit_config.md`
- 修改文件: `tencentcloud/provider.go`（注册新资源）
- 修改文件: `tencentcloud/provider.md`（添加资源文档链接）
- 修改文件: `tencentcloud/services/teo/service_tencentcloud_teo.go`（新增 service 层方法）
- 依赖云 API: `ModifyPrefetchOriginLimit`、`DescribePrefetchOriginLimit`（teo v20220901）
