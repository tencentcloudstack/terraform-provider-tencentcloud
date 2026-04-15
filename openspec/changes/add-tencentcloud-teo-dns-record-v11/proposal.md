## Why

用户需要通过 Terraform 管理 TEO（TencentCloud EdgeOne）产品的 DNS 记录资源。当前 Terraform Provider 中缺少对 TEO DNS 记录的管理支持，无法实现基础设施即代码的完整管理。新增 `tencentcloud_teo_dns_record_v11` 资源可以让用户声明式地管理 DNS 记录，提高运维效率和配置一致性。

## What Changes

- **新增资源**: 添加 `tencentcloud_teo_dns_record_v11` Terraform 资源，支持完整的 CRUD 操作
- **资源类型**: RESOURCE_KIND_GENERAL（通用资源）
- **支持的云 API 接口**:
  - CreateDnsRecord: 创建 DNS 记录
  - DescribeDnsRecords: 查询 DNS 记录列表
  - ModifyDnsRecords: 批量修改 DNS 记录
  - DeleteDnsRecords: 批量删除 DNS 记录

## Capabilities

### New Capabilities
- `teo-dns-record-v11`: 管理 TEO DNS 记录的完整生命周期，包括创建、读取、更新和删除操作

### Modified Capabilities
- (无)

## Impact

- **新增代码文件**:
  - `tencentcloud/services/teo/resource_tc_teo_dns_record_v11.go`: 资源实现
  - `tencentcloud/services/teo/resource_tc_teo_dns_record_v11_test.go`: 单元测试
  - `website/docs/r/teo_dns_record_v11.html.markdown`: 资源文档

- **修改代码文件**:
  - `tencentcloud/services/teo/service_tencentcloud_teo.go`: 注册新资源

- **依赖**:
  - 使用 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` 包中的 API

- **兼容性**: 向后兼容，不影响现有资源
