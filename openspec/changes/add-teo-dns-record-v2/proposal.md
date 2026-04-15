## Why

用户需要通过 Terraform 管理 TEO（边缘安全加速平台）的 DNS 记录。当前 Terraform Provider 中缺少对 teo 产品 DNS 记录的支持，用户无法以基础设施即代码的方式管理 teo 的 DNS 记录资源。新增该资源可以满足用户对 teo DNS 记录的自动化管理需求，提升运维效率。

## What Changes

新增一个 Terraform 资源 `tencentcloud_teo_dns_record_v2`，用于管理 TEO 产品的 DNS 记录，提供完整的 CRUD（创建、查询、修改、删除）功能。

具体变更：
- 创建新的资源文件 `resource_tc_teo_dns_record_v2.go`
- 实现资源创建接口 `CreateDnsRecord`
- 实现资源查询接口 `DescribeDnsRecords`
- 实现资源修改接口 `ModifyDnsRecords`
- 实现资源删除接口 `DeleteDnsRecords`
- 添加资源的单元测试文件 `resource_tc_teo_dns_record_v2_test.go`
- 添加资源的 Terraform acceptance test（使用 mock 方式）

## Capabilities

### New Capabilities
- `teo-dns-record-v2-resource`: 管理 TEO 产品的 DNS 记录资源，支持创建、查询、修改和删除操作

### Modified Capabilities
- (无现有 capability 需要修改)

## Impact

受影响的代码和系统：
- 新增文件：`tencentcloud/services/teo/resource_tc_teo_dns_record_v2.go`
- 新增文件：`tencentcloud/services/teo/resource_tc_teo_dns_record_v2_test.go`
- 修改文件：`tencentcloud/services/teo/service_tencentcloud_teo.go`（添加新资源注册）
- 新增依赖：使用 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` 包的以下接口：
  - CreateDnsRecord
  - DescribeDnsRecords
  - ModifyDnsRecords
  - DeleteDnsRecords
- 依赖云 API：TEO 产品的 DNS 记录管理接口
