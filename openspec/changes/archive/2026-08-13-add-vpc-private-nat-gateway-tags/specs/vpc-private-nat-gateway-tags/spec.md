## ADDED Requirements

### Requirement: Tags 参数定义

`tencentcloud_vpc_private_nat_gateway` 资源 SHALL 新增 `tags` 参数，类型为 `schema.TypeMap`，`Optional: true`，用于配置私网 NAT 网关的标签键值对。Key 和 Value 均为字符串类型。

#### Scenario: 用户在 Terraform 配置中指定 tags

- **WHEN** 用户在 `tencentcloud_vpc_private_nat_gateway` 资源配置中设置 `tags = { "key1" = "value1", "key2" = "value2" }`
- **THEN** 资源 schema SHALL 接受该配置，并在 plan 阶段显示 tags 参数

#### Scenario: 用户未配置 tags

- **WHEN** 用户在 `tencentcloud_vpc_private_nat_gateway` 资源配置中未设置 tags 参数
- **THEN** 资源 SHALL 正常创建，tags 参数为空

### Requirement: Create 方法处理 tags

资源 Create 方法 SHALL 将 schema 中的 `tags` map 转换为 `[]*vpc.Tag` 列表，并设置到 `CreatePrivateNatGatewayRequest.Tags` 字段。每个 Tag 元素的 `Key` 和 `Value` 从 map 的键值对获取。

#### Scenario: 创建资源时传入标签

- **WHEN** 用户创建 `tencentcloud_vpc_private_nat_gateway` 资源并配置了 tags
- **THEN** Create 方法 SHALL 将 tags map 转换为 `[]*vpc.Tag` 列表，并设置到 `CreatePrivateNatGatewayRequest.Tags`

#### Scenario: 创建资源时未传入标签

- **WHEN** 用户创建 `tencentcloud_vpc_private_nat_gateway` 资源且未配置 tags
- **THEN** Create 方法 SHALL 不设置 `CreatePrivateNatGatewayRequest.Tags` 字段

### Requirement: Read 方法回显 tags

资源 Read 方法 SHALL 将云 API 返回的 `PrivateNatGateway.TagSet`（`[]*vpc.Tag`）转换为 map，并 set 到 schema 的 `tags` 字段。

#### Scenario: 读取资源时云 API 返回标签

- **WHEN** Read 方法调用 `DescribePrivateNatGateways` 返回的 `PrivateNatGateway.TagSet` 包含标签数据
- **THEN** Read 方法 SHALL 将 TagSet 转换为 map 并 `d.Set("tags", tagsMap)`

#### Scenario: 读取资源时云 API 未返回标签

- **WHEN** Read 方法调用 `DescribePrivateNatGateways` 返回的 `PrivateNatGateway.TagSet` 为 nil 或空
- **THEN** Read 方法 SHALL 不执行 `d.Set("tags", ...)`，保留 state 中已有的值

### Requirement: Update 方法禁止修改 tags

资源 Update 方法 SHALL 将 `tags` 加入 `immutableArgs` 数组。当 tags 发生变更时，返回错误提示参数不可变更。

#### Scenario: 用户尝试修改 tags

- **WHEN** 用户修改 `tencentcloud_vpc_private_nat_gateway` 资源的 tags 参数并执行 `terraform apply`
- **THEN** Update 方法 SHALL 返回错误 `argument 'tags' cannot be changed`
