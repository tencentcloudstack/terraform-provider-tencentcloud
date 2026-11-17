### Requirement: subnet_id 字段标记为 Computed

`tencentcloud_reserve_ip_address` 资源的 `subnet_id` schema 字段 SHALL 标记为 `Optional: true, Computed: true`，以允许 Terraform 从云 API `DescribeReserveIpAddresses` 响应中刷新该值，即使用户在配置中未显式设置。

**Rationale**: 云 API `DescribeReserveIpAddresses` 的 `ReserveIpAddressSet` 中每个 `ReserveIpAddressInfo` 包含 `SubnetId` 字段。标记为 `Computed` 可避免当用户未设置 `subnet_id` 时产生 plan diff。

#### Scenario: 用户在配置中指定 subnet_id

- **WHEN** 用户在 `tencentcloud_reserve_ip_address` 资源配置中设置 `subnet_id = "subnet-xxxxxx"`
- **THEN** 资源 schema SHALL 接受该配置，Create 方法将值传入 `CreateReserveIpAddressesRequest.SubnetId`，Read 方法从 `DescribeReserveIpAddresses` 响应中回显该值

#### Scenario: 用户未在配置中指定 subnet_id

- **WHEN** 用户在 `tencentcloud_reserve_ip_address` 资源配置中未设置 `subnet_id`
- **THEN** 资源 SHALL 正常创建，Read 方法从 `DescribeReserveIpAddresses` 响应的 `ReserveIpAddressSet[0].SubnetId` 读取值并 `d.Set("subnet_id", ...)`，不会产生 plan diff

### Requirement: Read 方法回显 subnet_id

`tencentcloud_reserve_ip_address` 资源的 Read 方法 SHALL 从 `DescribeReserveIpAddresses` API 响应的 `ReserveIpAddressSet[0].SubnetId` 字段读取 `SubnetId` 值，并通过 `d.Set("subnet_id", ...)` 写入 Terraform state。

#### Scenario: 读取资源时云 API 返回 SubnetId

- **WHEN** Read 方法调用 `DescribeReserveIpAddresses` 返回的 `ReserveIpAddressSet[0].SubnetId` 不为 nil
- **THEN** Read 方法 SHALL 执行 `d.Set("subnet_id", reserveIpAddress.SubnetId)` 将值写入 state

#### Scenario: 读取资源时云 API 未返回 SubnetId

- **WHEN** Read 方法调用 `DescribeReserveIpAddresses` 返回的 `ReserveIpAddressSet[0].SubnetId` 为 nil
- **THEN** Read 方法 SHALL 不执行 `d.Set("subnet_id", ...)`，保留 state 中已有的值

### Requirement: Import 验证覆盖 subnet_id

`tencentcloud_reserve_ip_address` 资源的验收测试 SHALL 从 `ImportStateVerifyIgnore` 列表中移除 `subnet_id`，使 import 验证覆盖该字段。

#### Scenario: 资源导入后验证 subnet_id

- **WHEN** 使用 `terraform import` 导入 `tencentcloud_reserve_ip_address` 资源并执行 `ImportStateVerify`
- **THEN** 测试 SHALL 验证 `subnet_id` 字段从 API 响应正确读取，不再被忽略
