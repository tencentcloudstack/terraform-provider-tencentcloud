## Why

为了支持 VPC 端点资源的安全组绑定、标签管理和协议类型配置，需要为 `tencentcloud_vpc_end_point` 资源新增 3 个可选参数。这些功能在云 API 中已支持，但在 Terraform Provider 中尚未实现，导致用户无法通过 Terraform 管理这些重要属性。

## What Changes

在 `tencentcloud_vpc_end_point` 资源中新增 3 个字段：

- **SecurityGroupId** (string, 可选): 安全组 ID，用于绑定安全组到 VPC 端点
- **Tags** (list, 可选): 标签列表，用于为 VPC 端点打标签，包含 Key (必填) 和 Value (可选)
- **IpAddressType** (string, 可选): 协议类型，支持 "Ipv4" 或 "Ipv6"，默认为 "Ipv4"

更新 Create、Read、Update、Delete 函数以处理这些新字段：

- **Create**: 将新字段传递给 CreateVpcEndPoint API
- **Read**: 从 DescribeVpcEndPoints API 返回结果中读取新字段
- **Update**: 支持 UpdateVpcEndPointAttribute API 更新这些字段
- **Delete**: 无需特殊处理

更新相关测试代码以覆盖新增字段的功能。

## Capabilities

### New Capabilities

- `vpc-endpoint-security-group`: 为 VPC 端点资源添加安全组绑定功能
- `vpc-endpoint-tags`: 为 VPC 端点资源添加标签管理功能
- `vpc-endpoint-ip-address-type`: 为 VPC 端点资源添加协议类型配置功能

### Modified Capabilities

无（仅新增字段，不修改现有功能行为）

## Impact

**Affected Code**:
- `tencentcloud/services/vpc/resource_tc_vpc_end_point.go`: 修改 schema 和 CRUD 函数
- `tencentcloud/services/vpc/resource_tc_vpc_end_point_test.go`: 更新单元测试
- `tencentcloud/services/vpc/resource_tc_vpc_end_point.md`: 更新资源文档示例

**API Changes**:
- CreateVpcEndPoint API: 新增请求参数
- DescribeVpcEndPoints API: 新增返回字段读取
- UpdateVpcEndPointAttribute API: 新增更新参数

**Dependencies**:
- 依赖现有的 tencentcloud-sdk-go VPC API

**Systems**:
- Terraform Provider VPC 服务模块
- 用户现有的 Terraform 配置不受影响（新增字段均为可选）