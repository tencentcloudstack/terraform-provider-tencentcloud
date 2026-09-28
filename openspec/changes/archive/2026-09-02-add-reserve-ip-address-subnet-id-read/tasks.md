## 1. Schema 修改

- [x] 1.1 在 `tencentcloud/services/vpc/resource_tc_reserve_ip_address.go` 中，将 `subnet_id` schema 字段从 `Optional: true` 修改为 `Optional: true, Computed: true`
- [x] 1.2 确认 `subnet_id` 仍在 Update 方法的 `immutableArgs` 数组中（无需修改，仅验证）

## 2. Read 方法修改

- [x] 2.1 在 `resourceTencentCloudReserveIpAddressRead` 方法中，在现有 `d.Set` 调用之后、读取 tags 之前，添加 `subnet_id` 的回显逻辑
  - 先判断 `reserveIpAddress.SubnetId` 是否为 nil
  - 若不为 nil，执行 `_ = d.Set("subnet_id", reserveIpAddress.SubnetId)`

## 3. 测试修改

- [x] 3.1 在 `tencentcloud/services/vpc/resource_tc_reserve_ip_address_test.go` 中，从 `TestAccTencentCloudReserveIpAddressesResource_SetIpAddress` 测试用例的 `ImportStateVerifyIgnore` 中移除 `"subnet_id"`
- [x] 3.2 在 `tencentcloud/services/vpc/resource_tc_reserve_ip_address_test.go` 中，从 `TestAccTencentCloudReserveIpAddressesResource_NotSetIpAddress` 测试用例的 `ImportStateVerifyIgnore` 中移除 `"subnet_id"`

## 4. 文档

- [x] 4.1 检查 `tencentcloud/services/vpc/resource_tc_reserve_ip_address.md` 是否需要更新（`subnet_id` 已在示例中，预期无需修改）

## 5. 验证

- [x] 5.1 代码正确性检查：确认 Read 方法中 `SubnetId` 字段路径与 vendor SDK 中 `ReserveIpAddressInfo.SubnetId` 一致
- [x] 5.2 代码正确性检查：确认 schema 中 `subnet_id` 的 `Computed` 标记不会影响现有 Create 逻辑（Create 中仍从 `d.GetOk("subnet_id")` 读取）
