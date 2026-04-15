## Why

用户需要通过 Terraform 管理 EdgeOne（TEO）服务的 DNS 记录（V10 版本），以实现基础设施即代码的管理方式，提升 DNS 配置的可维护性和一致性。

## What Changes

- 新增 Terraform 资源 `tencentcloud_teo_dns_record_v10`，用于管理 TEO 服务的 DNS 记录
- 支持完整的 CRUD 操作：创建、读取、更新、删除 DNS 记录
- 支持的 DNS 记录类型包括：A、AAAA、MX、CNAME、TXT、NS、CAA、SRV
- 支持配置解析线路、缓存时间、权重、MX 优先级等参数
- 实现资源状态管理和异步操作支持

## Capabilities

### New Capabilities
- `teo-dns-record-v10`: 管理 TEO DNS 记录资源，支持创建、查询、修改和删除 DNS 记录，支持多种记录类型和高级配置参数

### Modified Capabilities
无

## Impact

- 新增文件：`tencentcloud/services/teo/resource_tc_teo_dns_record_v10.go`（资源实现）
- 新增文件：`tencentcloud/services/teo/resource_tc_teo_dns_record_v10_test.go`（单元测试）
- 新增文件：`website/docs/r/teo_dns_record_v10.md`（资源文档）
- 修改文件：`tencentcloud/services/teo/service_tencentcloud_teo.go`（添加服务方法）
- 修改文件：`tencentcloud/services/teo/tencentcloud_teo_suite_test.go`（添加测试配置）
- 新增文件：`website/docs/d/teo_*_test.go`（示例文档，如需要）
