## Why

需要支持用户通过 Terraform 导出 TEO (EdgeOne) 站点配置。当前用户无法通过 Terraform Provider 导出站点配置，导致无法对站点配置进行版本控制和备份。

## What Changes

- 新增 TEO 导出站点配置资源：`tencentcloud_teo_export_zone_config`
- 实现该资源的 Create、Read、Update、Delete 四个操作函数
- 生成对应的单元测试和验收测试代码
- 根据 CAPI 接口定义生成完整的 Resource Schema

## Capabilities

### New Capabilities
- `teo-export-zone-config`: 支持导出 TEO 站点配置的 Terraform 资源，提供站点配置的导出功能

### Modified Capabilities
- (无现有 capability 需要修改)

## Impact

- **代码变更**:
  - 新增文件：`tencentcloud/services/teo/resource_tencentcloud_teo_export_zone_config.go`
  - 新增文件：`tencentcloud/services/teo/resource_tencentcloud_teo_export_zone_config_test.go`
  - 可能需要修改：`tencentcloud/services/teo/service_tencentcloud_teo.go` (添加注册)
- **API 依赖**: 使用 TEO 相关的 CAPI 接口
- **文档**: 需要添加对应的文档和使用示例
