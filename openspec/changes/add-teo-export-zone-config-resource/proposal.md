## Why

用户需要通过 Terraform 管理腾讯云 EdgeOne (TEO) 站点配置的导出功能，目前缺少对应的 Terraform Provider Resource 支持。新增该资源将允许用户通过 IaC 方式导出站点配置，便于配置备份和迁移。

## What Changes

- 新增 `tencentcloud_teo_export_zone_config` Terraform Provider Resource
- 实现资源的 CRUD 操作函数：Create、Read、Update、Delete
- 生成对应的 Schema 定义，确保参数的 Required/Optional 属性与 CAPI 接口定义一致
- 生成单元测试代码（*_test.go）
- 生成验收测试代码
- 生成资源使用示例文档

## Capabilities

### New Capabilities
- `teo-export-zone-config`: 提供导出 TEO 站点配置的能力，支持通过 Terraform 管理配置导出操作

### Modified Capabilities
(无)

## Impact

- 新增文件：`tencentcloud/services/teo/resource_tencentcloud_teo_export_zone_config.go`
- 新增测试文件：`tencentcloud/services/teo/resource_tencentcloud_teo_export_zone_config_test.go`
- 新增文档：`website/docs/r/teo_export_zone_config.html.md`
- 更新服务层：可能需要在 `service_tencentcloud_teo.go` 中添加相关辅助函数
- 不影响现有资源和数据源，新增功能不涉及向后兼容问题
