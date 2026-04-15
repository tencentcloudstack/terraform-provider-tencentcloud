## Why

当前 Terraform Provider for TencentCloud 缺少对 Teo 产品 DNS 记录的管理能力。用户无法通过 Terraform IaC 方式管理 Teo 产品的 DNS 记录，这影响了基础设施即代码的完整性和一致性。新增该资源可以让用户通过 Terraform 声明式地创建、查询、修改和删除 Teo DNS 记录，提升运维效率和自动化水平。

## What Changes

- 新增资源 `tencentcloud_teo_dns_record_10`，支持 Teo 产品的 DNS 记录的完整生命周期管理
- 实现资源的 CRUD 操作：
  - Create: 通过 `CreateDnsRecord` API 创建 DNS 记录
  - Read: 通过 `DescribeDnsRecords` API 查询 DNS 记录
  - Update: 通过 `ModifyDnsRecords` API 修改 DNS 记录
  - Delete: 通过 `DeleteDnsRecords` API 删除 DNS 记录
- 在 `tencentcloud/services/teo/` 目录下创建相关文件：
  - 资源实现文件：`resource_tc_teo_dns_record_10.go`
  - 资源测试文件：`resource_tc_teo_dns_record_10_test.go`
  - 资资源文档：`resource_tc_teo_dns_record_10.md`
  - 服务层文件（如需要）：`service_tencentcloud_teo.go`

## Capabilities

### New Capabilities

- `teo-dns-record-10`: 描述 Teo 产品的 DNS 记录资源的完整行为规范，包括资源标识、参数定义、CRUD 接口映射、错误处理和异步操作等

### Modified Capabilities

无。这是一个全新的资源，不涉及对现有规格的修改。

## Impact

**代码层面：**
- 新增文件：`tencentcloud/services/teo/resource_tc_teo_dns_record_10.go`（资源实现）
- 新增文件：`tencentcloud/services/teo/resource_tc_teo_dns_record_10_test.go`（单元测试）
- 新增文件：`website/docs/r/teo_dns_record_10.html.md`（资源文档）
- 可能修改文件：`tencentcloud/services/teo/service_tencentcloud_teo.go`（如需要新增服务层方法）
- 可能修改文件：`tencentcloud/provider.go`（注册新资源）

**API 依赖：**
- 依赖 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` 包的以下接口：
  - `CreateDnsRecord`
  - `DescribeDnsRecords`
  - `ModifyDnsRecords`
  - `DeleteDnsRecords`

**测试层面：**
- 需要编写单元测试（使用 mock 方式，避免实际调用云 API）
- 需要准备测试数据覆盖各种场景

**文档层面：**
- 需要在 `website/docs/` 目录下生成资源文档
- 需要提供 Terraform 配置示例

**兼容性：**
- 这是一个新增资源，不会破坏现有资源的向后兼容性
- 遵循 Terraform Provider 的命名和结构约定
