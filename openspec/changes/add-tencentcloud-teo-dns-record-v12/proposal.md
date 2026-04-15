## Why

当前 Terraform TencentCloud Provider 缺少对 TEO (TencentCloud EdgeOne) 产品 DNS 记录资源的支持。用户无法通过 Terraform 管理边缘加速服务（TEO）的 DNS 记录配置，导致无法实现边缘 DNS 的自动化部署和管理。为满足用户对边缘 DNS 记录进行基础设施即代码（IaC）管理的需求，需要新增 `tencentcloud_teo_dns_record_v12` 资源。

## What Changes

- **新增 Terraform 资源**: `tencentcloud_teo_dns_record_v12`，用于管理 TEO 产品的 DNS 记录
- **资源类型**: RESOURCE_KIND_GENERAL（通用资源），支持完整的 CRUD 操作
- **代码文件**: 创建 `resource_tc_teo_dns_record_v12.go` 资源实现文件
- **API 映射**:
  - 创建资源：调用 `CreateDnsRecord` 接口
  - 读取资源：调用 `DescribeDnsRecords` 接口
  - 更新资源：调用 `ModifyDnsRecords` 接口
  - 删除资源：调用 `DeleteDnsRecords` 接口
- **资源参数**:
  - `zone_id`: 站点 ID（必需）
  - `name`: DNS 记录名（必需）
  - `type`: DNS 记录类型（必需）
  - `content`: DNS 记录内容（必需）
  - `location`: DNS 记录解析线路（可选）
  - `ttl`: 缓存时间，范围 60~86400 秒（可选，默认 300）
  - `weight`: DNS 记录权重，范围 -1~100（可选，默认 -1）
  - `priority`: MX 记录优先级，范围 0~50（可选，默认 0）
  - `record_id`: DNS 记录 ID（计算属性，由系统返回）
  - `status`: DNS 记录解析状态（只读，enable/disable）
  - `created_on`: 创建时间（只读）

## Capabilities

### New Capabilities
- `teo-dns-record-v12`: TEO DNS 记录资源管理能力，支持通过 Terraform 创建、读取、更新和删除 TEO 产品的 DNS 记录，包括 A、AAAA、MX、CNAME、TXT、NS、CAA、SRV 等类型的记录管理，支持配置解析线路、缓存时间、权重和优先级等属性。

### Modified Capabilities
- 无现有能力的需求变更

## Impact

- **新增代码文件**:
  - `tencentcloud/services/teo/resource_tc_teo_dns_record_v12.go`: 资源实现代码
  - `tencentcloud/services/teo/resource_tc_teo_dns_record_v12_test.go`: 资源单元测试
  - `website/docs/r/teo_dns_record_v12.md`: 资源文档
  - `examples/resources/tencentcloud_teo_dns_record_v12/resource.tf`: 资源使用示例

- **依赖变更**:
  - 使用现有的 tencentcloud-sdk-go TEO 服务包（v20220901）
  - 新增对 `CreateDnsRecord`、`DescribeDnsRecords`、`ModifyDnsRecords`、`DeleteDnsRecords` 四个 API 的依赖

- **受影响系统**:
  - Terraform Provider 的资源注册表
  - 文档生成系统（website/docs/）
  - 测试套件

- **向后兼容性**:
  - 无破坏性变更，仅新增资源
  - 不影响现有资源和数据源
