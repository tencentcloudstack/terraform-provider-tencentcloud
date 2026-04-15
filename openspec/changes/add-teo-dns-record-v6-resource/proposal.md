## Why

用户需要通过 Terraform 管理边缘安全加速平台（TEO）的 DNS 记录。目前缺少对 TEO DNS 记录 v6 版本的管理能力，无法实现基础设施即代码的统一管理。

## What Changes

- 新增 Terraform 资源：`tencentcloud_teo_dns_record_v6`
- 实现资源的 CRUD 操作（Create、Read、Update、Delete）
- 根据 CAPI 接口定义生成对应的 Schema 参数
- 添加单元测试和验收测试代码

## Capabilities

### New Capabilities
- `teo-dns-record-v6`: 提供对 TEO DNS 记录的完整管理能力，包括创建、读取、更新和删除操作，支持 IPv6 地址类型的 DNS 记录管理

### Modified Capabilities
(无)

## Impact

- 新增资源文件：`tencentcloud/services/teo/resource_tencentcloud_teo_dns_record_v6.go`
- 新增测试文件：`tencentcloud/services/teo/resource_tencentcloud_teo_dns_record_v6_test.go`
- 新增文档文件：`website/docs/r/teo_dns_record_v6.md`
- 依赖 TEO 服务的 CAPI 接口：CreateDnsRecord、DescribeDnsRecords、ModifyDnsRecords、DeleteDnsRecords
