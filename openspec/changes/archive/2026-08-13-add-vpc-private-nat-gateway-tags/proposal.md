## Why

`tencentcloud_vpc_private_nat_gateway` 资源当前不支持标签（Tags）参数。云 API 的 `CreatePrivateNatGateway` 接口已支持通过 `Tags` 字段传入标签键值对，`DescribePrivateNatGateways` 接口的 `PrivateNatGatewaySet.TagSet` 也已返回标签数据。用户无法在 Terraform 中为私网 NAT 网关配置和管理标签，导致无法使用标签进行资源分组、成本归集和权限管理。

## What Changes

- 在 `tencentcloud_vpc_private_nat_gateway` 资源的 schema 中新增 `tags` 参数（TypeMap，Optional），用于配置私网 NAT 网关的标签键值对。
- 在资源 Create 方法中，将 schema 中的 `tags` 映射转换为云 API 的 `[]*Tag` 列表，设置到 `CreatePrivateNatGatewayRequest.Tags`。
- 在资源 Read 方法中，将云 API 返回的 `PrivateNatGateway.TagSet`（`[]*Tag`）转换为 map 并 set 到 schema 的 `tags` 字段。
- 由于 `tags` 参数在 Create 时传入，Update 方法中将其加入 `immutableArgs`，标签变更需重建资源。

## Capabilities

### New Capabilities
- `vpc-private-nat-gateway-tags`: 为 `tencentcloud_vpc_private_nat_gateway` 资源新增 tags 参数，支持在创建私网 NAT 网关时配置标签，并支持在 Read 时回显标签。

### Modified Capabilities
<!-- 无 -->

## Impact

- **代码文件**: `tencentcloud/services/vpc/resource_tc_vpc_private_nat_gateway.go` — 新增 schema 字段、修改 Create/Read/Update 方法。
- **测试文件**: `tencentcloud/services/vpc/resource_tc_vpc_private_nat_gateway_test.go` — 补充 tags 参数的单元测试用例。
- **文档文件**: `tencentcloud/services/vpc/resource_tc_vpc_private_nat_gateway.md` — 更新示例和文档。
- **云 API 依赖**: 依赖 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312` 中已有的 `Tag` 结构和 `CreatePrivateNatGatewayRequest.Tags` / `PrivateNatGateway.TagSet` 字段，无需更新 vendor。
- **向后兼容**: 新增 Optional 字段，不影响已有配置和 state。
