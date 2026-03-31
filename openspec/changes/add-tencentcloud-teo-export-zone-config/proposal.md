# Add tencentcloud_teo_export_zone_config resource

## Why

用户需要导出站点配置的功能，以便将站点配置导出为 JSON 格式，便于配置管理和备份。目前 Terraform Provider 中缺少导出站点配置的能力，需要添加新的 resource 来满足这一需求。

## What Changes

- 添加新的 resource: `tencentcloud_teo_export_zone_config`
- 支持通过 ZoneId 和 Types 参数导出站点配置
- 返回导出的配置内容（JSON 格式）
- 实现完整的 CRUD 操作函数
- 添加单元测试和验收测试
- 添加资源文档和使用样例

## Capabilities

### New Capabilities
- `teo-export-zone-config`: TEO 站点配置导出能力，支持通过 ZoneId 和配置类型列表导出站点配置

### Modified Capabilities

无

## Impact

- 新增文件：
  - `tencentcloud/services/teo/resource_tencentcloud_teo_export_zone_config.go`
  - `tencentcloud/services/teo/resource_tencentcloud_teo_export_zone_config_test.go`
  - `website/docs/r/teo_export_zone_config.html.markdown`
  - `examples/resources/teo_export_zone_config/README.md`
- 调用 TencentCloud Teo API: ExportZoneConfig
- 无破坏性变更
- 向后兼容
