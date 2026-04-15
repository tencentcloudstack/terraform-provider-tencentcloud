## Why

TEO (TencentCloud EdgeOne) 是腾讯云的边缘加速产品，DNS 记录管理是 TEO 的核心功能之一。当前 Terraform Provider 缺少对 TEO DNS 记录的资源支持，用户无法通过 Terraform 以基础设施即代码的方式管理 TEO 的 DNS 记录。新增此资源将填补这一空白，为用户提供完整的 IaC 管理能力。

## What Changes

- 新增 Terraform 资源 `tencentcloud_teo_dns_record_v13`，支持 TEO DNS 记录的完整生命周期管理
- 实现 Create 操作：通过 `CreateDnsRecord` API 创建 DNS 记录
- 实现 Read 操作：通过 `DescribeDnsRecords` API 查询 DNS 记录
- 实现 Update 操作：通过 `ModifyDnsRecords` API 批量修改 DNS 记录
- 实现 Delete 操作：通过 `DeleteDnsRecords` API 批量删除 DNS 记录
- 支持所有 DNS 记录类型：A、AAAA、CNAME、TXT、MX、NS、CAA、SRV
- 支持解析线路配置（Location）和记录权重（Weight）
- 支持记录的启用/停用状态管理
- 生成对应的单元测试文件和文档

## Capabilities

### New Capabilities

- `teo-dns-record`: TEO DNS 记录管理能力，包括创建、查询、修改和删除 DNS 记录，支持所有记录类型、解析线路、权重、TTL、优先级等配置项

### Modified Capabilities

（无）

## Impact

**新增文件**:
- `tencentcloud/services/teo/resource_tc_teo_dns_record_v13.go` - 资源实现
- `tencentcloud/services/teo/resource_tc_teo_dns_record_v13_test.go` - 单元测试
- `website/docs/r/teo_dns_record_v13.html.markdown` - 资源文档（通过 make doc 生成）
- `tencentcloud/services/teo/tea_dns_record_v13_helper.go` - 辅助函数（如需要）

**修改文件**:
- `tencentcloud/services/teo/service_tencentcloud_teo.go` - 注册新资源

**依赖**:
- 使用 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` 包中的云 API
- 依赖 Terraform Plugin SDK v2 和标准 helper 函数

**API 接口**:
- `CreateDnsRecord` - 创建 DNS 记录
- `DescribeDnsRecords` - 查询 DNS 记录列表
- `ModifyDnsRecords` - 批量修改 DNS 记录
- `DeleteDnsRecords` - 批量删除 DNS 记录
