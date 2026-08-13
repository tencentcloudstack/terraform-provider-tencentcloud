## Context

`tencentcloud_vpc_private_nat_gateway` 资源（`resource_tc_vpc_private_nat_gateway.go`）当前管理私网 NAT 网关的创建、读取、更新和删除，但未暴露 tags 参数。云 API `CreatePrivateNatGateway` 已支持 `Tags []*Tag` 入参（`Tag` 结构包含 `Key` 和 `Value`），`DescribePrivateNatGateways` 返回的 `PrivateNatGateway.TagSet` 也已包含 `[]*Tag`。因此具备在 Terraform 中支持 tags 的条件。

现有资源 schema 使用扁平结构，已有字段：`nat_gateway_name`、`vpc_id`、`cross_domain`、`vpc_type`、`ccn_id`。Update 方法使用 `ModifyPrivateNatGatewayAttribute` 接口，该接口不支持修改标签，因此 tags 变更需 ForceNew。

## Goals / Non-Goals

**Goals:**
- 在 `tencentcloud_vpc_private_nat_gateway` 资源 schema 中新增 `tags` 参数（TypeMap, Optional），支持标签键值对配置。
- Create 方法将 tags map 转换为 `[]*vpc.Tag` 并设置到 `CreatePrivateNatGatewayRequest.Tags`。
- Read 方法将 `PrivateNatGateway.TagSet` 转换为 map 并 set 到 state。
- 保持向后兼容，新增字段为 Optional，不影响已有配置。

**Non-Goals:**
- 不支持通过 Update 方法动态修改标签（标签变更需重建资源）。
- 不修改 `ModifyPrivateNatGatewayAttribute` 接口调用逻辑。
- 不修改 `DescribeVpcPrivateNatGatewayById` 服务层方法（该方法已返回完整的 `PrivateNatGateway` 对象，包含 `TagSet`）。

## Decisions

### 1. tags schema 类型选择 TypeMap

采用 `schema.TypeMap`（key 为 string, value 为 string），与同仓库中 `resource_tc_nat_gateway.go` 的 tags 实现保持一致。

**理由**: TypeMap 使用简单，符合腾讯云标签键值对的使用习惯，与 vpc 服务的其他资源（如 `tencentcloud_nat_gateway`）保持风格统一。

### 2. tags 转换为 `[]*vpc.Tag` 的实现

Create 方法中遍历 `d.GetOk("tags")` 返回的 map，构建 `[]*vpc.Tag` 列表，每个元素设置 `Key` 和 `Value`。

**理由**: 云 API `CreatePrivateNatGatewayRequest.Tags` 的类型是 `[]*Tag`，需要从 Terraform 的 map 类型转换。

### 3. TagSet 转换为 map 的实现

Read 方法中遍历 `privateNatGateway.TagSet`，将每个 `Tag` 的 `Key` 和 `Value` 写入 map，然后 `d.Set("tags", tagsMap)`。

**理由**: 将云 API 的 `[]*Tag` 列表转换为 Terraform 的 map 格式存储到 state。

### 4. tags 加入 immutableArgs

由于 `ModifyPrivateNatGatewayAttribute` 接口不支持修改标签，将 `tags` 加入 Update 方法的 `immutableArgs` 数组。当 tags 发生变更时，返回错误提示用户需重建资源。

**理由**: 与现有 `vpc_id`、`cross_domain`、`vpc_type`、`ccn_id` 的处理方式一致，这些字段同样不支持通过 Update 修改。

## Risks / Trade-offs

- **[tags 变更需重建资源]** → 与 vpc 服务其他资源一致，用户需通过 `terraform taint` 或删除重建来变更标签。在 immutableArgs 中明确报错提示。
- **[空 TagSet 处理]** → Read 方法中当 `TagSet` 为 nil 或空时，不执行 `d.Set("tags", ...)`，避免设置空 map 覆盖用户配置。参考现有代码中其他字段的 nil 检查模式。
