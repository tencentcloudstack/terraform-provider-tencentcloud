## ADDED Requirements

### Requirement: storage_pay_mode 支持从预付费转为后付费
`tencentcloud_cynosdb_cluster` 资源 SHALL 支持通过修改 `storage_pay_mode` 参数将集群的预付费存储（`1`）转为后付费存储（`0`），调用 `TransferStoragePrepayToPostpay` 接口完成转换。

#### Scenario: 预付费存储转为后付费存储
- **WHEN** 用户将已创建集群的 `storage_pay_mode` 参数从 `1`（预付费）修改为 `0`（后付费）
- **THEN** Provider SHALL 从 Update 函数的 `immutableArgs` 校验中放行该字段变更
- **AND** 调用 `TransferStoragePrepayToPostpay` 接口，传入集群 ID（`ClusterId`）
- **AND** 通过返回值中的 `DealNames` 调用 `DescribeResourcesByDealName` 接口轮询等待转换完成
- **AND** 转换完成后通过 Read 操作将新的 `storage_pay_mode` 值同步到 Terraform 状态

#### Scenario: 拒绝后付费转预付费
- **WHEN** 用户将 `storage_pay_mode` 参数从 `0`（后付费）修改为 `1`（预付费）
- **THEN** Provider SHALL 返回错误，明确说明 `storage_pay_mode` 仅支持从预付费(1)转为后付费(0)

#### Scenario: 字段未变更时不触发转换
- **WHEN** 用户未修改 `storage_pay_mode` 参数
- **THEN** Provider SHALL 不调用 `TransferStoragePrepayToPostpay` 接口
- **AND** 其他字段的更新行为不受影响

### Requirement: 异步等待存储付费模式转换完成
Provider SHALL 在调用 `TransferStoragePrepayToPostpay` 接口成功后，通过返回值中的 `DealNames` 调用 `DescribeResourcesByDealName` 接口轮询，等待存储付费模式转换完成，参考 create 模块中的异步等待实现方式。

#### Scenario: 转换进行中时持续轮询
- **WHEN** `TransferStoragePrepayToPostpay` 接口调用成功并返回 `DealNames`
- **AND** 调用 `DescribeResourcesByDealName` 接口返回错误码 `InvalidParameterValue.DealNameNotFound`（计费未回调发货）
- **THEN** Provider SHALL 将该错误作为可重试错误继续轮询
- **AND** 使用 `tccommon.ReadRetryTimeout` 作为最大等待时间

#### Scenario: 转换完成时停止轮询
- **WHEN** 调用 `DescribeResourcesByDealName` 接口成功返回 `BillingResourceInfos`
- **THEN** Provider SHALL 停止轮询
- **AND** 继续执行后续 Update 逻辑

#### Scenario: DealNames 为空时返回错误
- **WHEN** `TransferStoragePrepayToPostpay` 接口调用成功但返回的 `DealNames` 为空
- **THEN** Provider SHALL 返回错误
- **AND** 错误信息中包含集群 ID 便于排障

#### Scenario: 轮询超时
- **WHEN** 在 `tccommon.ReadRetryTimeout` 时间内 `DescribeResourcesByDealName` 未成功返回资源信息
- **THEN** Provider SHALL 返回超时错误
- **AND** 错误信息中包含 clusterId 与 DealName

### Requirement: 服务层 TransferStoragePrepayToPostpay 方法
CynosdbService SHALL 提供 `TransferStoragePrepayToPostpay` 方法用于将预付费存储转为后付费存储。

#### Scenario: 调用 TransferStoragePrepayToPostpay API
- **WHEN** 调用 `cynosdbService.TransferStoragePrepayToPostpay(ctx, clusterId)`
- **THEN** 应构造 `TransferStoragePrepayToPostpayRequest` 并设置 `ClusterId` 参数
- **AND** 调用 API 并返回 `DealNames` 列表

#### Scenario: API 调用失败处理
- **WHEN** `TransferStoragePrepayToPostpay` API 返回错误
- **THEN** 使用 `tccommon.RetryError` 包装错误
- **AND** 使用 `resource.Retry` 与 `tccommon.WriteRetryTimeout` 进行重试

### Requirement: immutableArgs 列表更新
`storage_pay_mode` SHALL 从 Update 函数的 `immutableArgs` 列表中移除，以支持更新操作。

#### Scenario: storage_pay_mode 不在 immutableArgs 中
- **WHEN** 检查 Update 函数的 `immutableArgs` 列表
- **THEN** 该列表 SHALL 不包含 `storage_pay_mode`
- **AND** Update 函数应能通过 `d.HasChange("storage_pay_mode")` 检测到字段变更

### Requirement: 文档更新
文档 SHALL 明确说明 `storage_pay_mode` 参数支持从预付费转为后付费，且为单向操作。

#### Scenario: 参数说明更新
- **WHEN** 查看 `resource_tc_cynosdb_cluster.md` 文档中 `storage_pay_mode` 参数描述
- **THEN** 应说明该参数支持修改
- **AND** 应说明仅支持从预付费(`1`)转为后付费(`0`)，不支持反向操作
