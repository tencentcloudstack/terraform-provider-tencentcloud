## Why

用户需要通过 Terraform 管理腾讯云边缘安全加速平台 (TEO) 的 DNS 记录资源。当前 Terraform Provider 中缺乏对 TEO DNS 记录的支持，用户无法通过基础设施即代码的方式自动化管理 DNS 记录的创建、修改和删除。

## What Changes

为云产品 `teo` 新增以下 Terraform 资源：
- **新增资源**: `tencentcloud_teo_dns_record_10` - 管理 TEO DNS 记录的通用资源
  - 支持 DNS 记录的创建
  - 支持 DNS 记录的读取
  - 支持 DNS 记录的更新
  - 支持 DNS 记录的删除

资源类型为 **RESOURCE_KIND_GENERAL**，实现完整的 CRUD 生命周期管理。

## Capabilities

### New Capabilities
- `teo-dns-record-10`: 管理 TEO DNS 记录资源，支持 A、AAAA、MX、CNAME、TXT、NS、CAA、SRV 等多种记录类型的创建、查询、修改和删除

### Modified Capabilities

## Impact

**新增文件**:
- `tencentcloud/services/teo/resource_tc_teo_dns_record_10.go` - 资源实现文件
- `tencentcloud/services/teo/resource_tc_teo_dns_record_10_test.go` - 资源单元测试
- `website/docs/r/teo_dns_record_10.md` - 资源文档

**修改文件**:
- `tencentcloud/services/teo/service_tencentcloud_teo.go` - 注册新资源

**依赖**:
- 使用 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` 包中的以下接口:
  - `CreateDnsRecord` - 创建 DNS 记录
  - `DescribeDnsRecords` - 查询 DNS 记录列表
  - `ModifyDnsRecords` - 批量修改 DNS 记录
  - `DeleteDnsRecords` - 批量删除 DNS 记录

**API 调用**:
- CreateDnsRecord: ZoneId, Name, Type, Content, Location, TTL, Weight, Priority
- DescribeDnsRecords: ZoneId, Filters (id, name, content, type, ttl)
- ModifyDnsRecords: ZoneId, DnsRecords (包含 RecordId, Name, Type, Location, Content, TTL, Weight, Priority)
- DeleteDnsRecords: ZoneId, RecordIds

**兼容性**:
- 新增资源，不影响现有资源
- 保持向后兼容性
