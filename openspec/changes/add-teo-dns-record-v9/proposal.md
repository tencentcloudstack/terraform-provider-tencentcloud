## Why

为 teo（EdgeOne）产品新增 DNS 记录管理能力，满足用户通过 Terraform 管理站点 DNS 解析记录的需求，提供基础设施即代码（IaC）的标准化管理方式。

## What Changes

- 新增 Terraform 资源 `tencentcloud_teo_dns_record_v9`，支持 teo 产品的 DNS 记录管理
- 实现资源的完整 CRUD 操作：创建、读取、更新、删除 DNS 记录
- 支持多种 DNS 记录类型：A、AAAA、MX、CNAME、TXT、NS、CAA、SRV
- 支持解析线路配置（适用于 A、AAAA、CNAME 记录类型）
- 支持记录权重配置（适用于 A、AAAA、CNAME 记录类型）
- 支持记录优先级配置（适用于 MX 记录类型）

## Capabilities

### New Capabilities

- `teo-dns-record-v9`: 管理腾讯云 EdgeOne（teo）产品的 DNS 记录，包括创建、更新、删除和查询操作，支持多种记录类型和高级配置选项

### Modified Capabilities

- 无

## Impact

- 新增资源文件：`tencentcloud/services/teo/resource_tc_teo_dns_record_v9.go`
- 新增资源测试文件：`tencentcloud/services/teo/resource_tc_teo_dns_record_v9_test.go`
- 新增资源样例文档：`tencentcloud/services/teo/resource_tc_teo_dns_record_v9.md`
- 修改服务层文件（如需要）：`tencentcloud/services/teo/service_tencentcloud_teo.go`
- 依赖 teo SDK 接口：CreateDnsRecord、DescribeDnsRecords、ModifyDnsRecords、DeleteDnsRecords
