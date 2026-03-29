## Why

用户需要导出 TEO (TencentCloud EdgeOne) 站点的配置信息，以便进行配置备份、迁移或审计分析。当前 Terraform Provider for TencentCloud 缺少导出站点配置的数据源能力，无法满足用户的这一需求。

## What Changes

- 新增数据源资源 `tencentcloud_teo_export_zone_config`，用于导出 TEO 站点的完整配置
- 添加数据源定义文件 `data_source_tc_teo_export_zone_config.go`
- 添加数据源测试文件 `data_source_tc_teo_export_zone_config_test.go`
- 添加数据源文档文件 `data_source_tc_teo_export_zone_config.md`
- 集成到 TEO 服务层，调用对应的 TencentCloud API

## Capabilities

### New Capabilities
- `teo-export-zone-config`: 导出 TEO 站点配置的数据源能力，支持查询指定站点 ID 的完整配置信息

### Modified Capabilities
无

## Impact

- **代码变更**: 新增 TEO 服务相关代码文件
- **API 调用**: 需要调用 TEO 导出站点配置的相关 API
- **文档**: 需要编写数据源使用文档和示例
- **测试**: 需要编写验收测试用例，验证数据源查询功能
- **依赖**: 可能需要更新 TEO 服务的 API SDK 版本
