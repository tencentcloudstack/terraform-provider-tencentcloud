## Why

TEO (TencentCloud EdgeOne) 服务需要提供一个导出站点配置的 Resource，使用户能够通过 Terraform 管理、导出站点配置。这是 EdgeOne 产品的重要组成部分，需要完整的 Terraform Provider 支持以实现基础设施即代码。

## What Changes

- 添加新的 Terraform Provider Resource: `tencentcloud_teo_export_zone_config`
- 实现完整的 CRUD 操作函数（Create、Read、Update、Delete）
- 根据 CAPI 接口定义生成 Resource Schema
- 添加单元测试和验收测试代码
- 添加 Resource 文档

## Capabilities

### New Capabilities
- `teo-export-zone-config`: TEO 站点配置导出功能，支持通过 Terraform 管理站点配置的导出操作

### Modified Capabilities

## Impact

- 新增文件：`tencentcloud/services/teo/resource_tc_teo_export_zone_config.go` (Resource 实现)
- 新增文件：`tencentcloud/services/teo/resource_tc_teo_export_zone_config_test.go` (单元测试)
- 新增文件：`tencentcloud/services/teo/resource_tc_teo_export_zone_config_acceptance_test.go` (验收测试)
- 新增文件：`website/docs/r/teo_export_zone_config.md` (资源文档)
- 新增文件：`examples/resources/tencentcloud_teo_export_zone_config/README.md` (使用样例)
