## Why

当前 tencentcloud_vpc_end_point 资源缺少安全组、标签和 IP 地址类型这三个重要参数，这些参数在实际使用中对于安全配置、资源管理和网络配置是必要的。需要补充这些参数以满足用户的实际需求。

## What Changes

在 tencentcloud_vpc_end_point 资源中新增以下三个字段：

- **SecurityGroupId** (string, 可选): 安全组 ID，用于端点的安全配置
- **Tags** (list, 可选): 标签列表，用于资源管理和分类
  - Key (string, 必填): 标签键
  - Value (string, 可选): 标签值
- **IpAddressType** (string, 可选): IP 地址类型，支持 Ipv4 和 Ipv6，默认 Ipv4

## Capabilities

### New Capabilities

- `vpc-endpoint-params`: 新增 tencentcloud_vpc_end_point 资源的 SecurityGroupId、Tags 和 IpAddressType 字段支持，以及相应的 CRUD 操作逻辑

### Modified Capabilities

无（这是新增字段，不涉及现有 spec 需求变更）

## Impact

- 受影响的资源文件：`tencentcloud/services/vpc/resource_tencentcloud_vpc_end_point.go`
- 受影响的测试文件：`tencentcloud/services/vpc/resource_tencentcloud_vpc_end_point_test.go`
- 涉及的 API：CreateVpcEndPoint（Create 操作）、DescribeVpcEndPoints（Read 操作）
- 需要更新相关文档和使用示例
