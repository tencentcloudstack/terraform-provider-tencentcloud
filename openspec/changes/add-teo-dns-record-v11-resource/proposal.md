## Why

为腾讯云 EdgeOne (TEO) 产品提供 DNS 记录管理能力。当前 Terraform Provider 缺少对 TEO DNS 记录的管理功能，用户无法通过 Terraform 代码管理和自动化 TEO DNS 记录的创建、修改、删除等操作。

## What Changes

- 新增 Terraform 资源 `tencentcloud_teo_dns_record_v11`，支持完整的 CRUD 操作
- 实现以下云 API 接口的集成：
  - CreateDnsRecord: 创建 DNS 记录
  - DescribeDnsRecords: 查询 DNS 记录列表
  - ModifyDnsRecords: 修改 DNS 记录
  - DeleteDnsRecords: 删除 DNS 记录
- 支持多种 DNS 记录类型：A、AAAA、MX、CNAME、TXT、NS、CAA、SRV
- 支持以下字段配置：
  - ZoneId: 站点 ID
  - Name: DNS 记录名
  - Type: DNS 记录类型
  - Content: DNS 记录内容
  - TTL: 缓存时间
  - Weight: DNS 记录权重
  - Priority: MX 记录优先级
  - Location: DNS 记录解析线路

## Capabilities

### New Capabilities

- `teo-dns-record-v11`: TEO DNS 记录 V11 资源管理能力，支持创建、读取、更新、删除 DNS 记录

### Modified Capabilities

无

## Impact

- 新增文件：`tencentcloud/services/teo/resource_tc_teo_dns_record_v11.go`
- 新增测试文件：`tencentcloud/services/teo/resource_tc_teo_dns_record_v11_test.go`
- 新增文档文件：`website/docs/r/teo_dns_record_v11.html.markdown`
- 新增示例文件：`tencentcloud/services/teo/resource_tc_teo_dns_record_v11.md`
- 修改文件：`tencentcloud/services/teo/service_tencentcloud_teo.go`（可能需要注册资源）
- 新增依赖：使用现有的 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` 包
