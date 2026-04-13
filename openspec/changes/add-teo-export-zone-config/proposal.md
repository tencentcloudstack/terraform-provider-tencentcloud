## Why

用户需要通过 Terraform 管理 TEO（TencentCloud Edge One）的站点配置导出功能，以便以基础设施即代码的方式自动化站点配置的导出和管理流程。目前 Terraform Provider 中缺少对导出站点配置能力的支持。

## What Changes

- 新增 `tencentcloud_teo_export_zone_config` 资源
- 实现该资源的 Create、Read、Update、Delete 操作函数
- 根据资源 UID `iacpres-ZHk6oZ2uSM` 对应的 CAPI 接口生成 Resource Schema 定义
- 添加对应的单元测试和验收测试代码

## Capabilities

### New Capabilities

- `teo-export-zone-config`: 提供导出 TEO 站点配置的能力，支持通过 Terraform 管理站点配置的导出操作

### Modified Capabilities

(无)

## Impact

- 新增文件：`tencentcloud/services/teo/resource_tencentcloud_teo_export_zone_config.go`
- 新增文件：`tencentcloud/services/teo/resource_tencentcloud_teo_export_zone_config_test.go`
- 新增文件：`tencentcloud/services/teo/resource_tencentcloud_teo_export_zone_config.md` (样例文档)
- 可能需要更新 `tencentcloud/services/teo/service_tencentcloud_teo.go` 添加相关辅助函数
