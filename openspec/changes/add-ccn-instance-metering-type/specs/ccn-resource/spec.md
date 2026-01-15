# 云联网资源 InstanceMeteringType 参数规范

## MODIFIED Requirements

### Requirement: CCN 资源 Schema 定义
系统 SHALL 在 `tencentcloud_ccn` 资源 Schema 中支持 `instance_metering_type` 参数，用于配置云联网实例的计量模式。

#### Scenario: 添加计量类型参数到 Schema
- **WHEN** 定义 `tencentcloud_ccn` 资源 Schema
- **THEN** 包含 `instance_metering_type` 字段，类型为 `schema.TypeString`
- **AND** 字段为可选参数（`Optional: true`）
- **AND** 字段标记为不可修改（`ForceNew: true`），修改需要重建资源
- **AND** Description 说明："Instance metering type. This parameter cannot be modified after creation"

#### Scenario: 参数值传递到 CreateCcn API
- **WHEN** 用户在配置中指定 `instance_metering_type` 参数
- **THEN** 系统在调用 `CreateCcn` API 时将该值传递给 `request.InstanceMeteringType`
- **AND** 如果用户未指定该参数，则不设置 `request.InstanceMeteringType`（使用云平台默认值）

#### Scenario: 参数值从 DescribeCcns API 读取
- **WHEN** 系统执行 Read 操作刷新资源状态
- **THEN** 从 `CCN.InstanceMeteringType` 字段读取计量类型
- **AND** 使用 `d.Set("instance_metering_type", value)` 同步到 Terraform 状态
- **AND** 如果 API 响应中该字段为 nil（老实例可能不返回），则不设置状态或设置为空字符串

#### Scenario: 不可修改约束验证
- **WHEN** 用户尝试修改已创建资源的 `instance_metering_type` 值
- **THEN** Terraform 检测到 ForceNew 字段变更
- **AND** 显示计划提示将销毁并重新创建资源
- **AND** 执行 `terraform apply` 后，旧资源被删除，新资源使用新的计量类型创建

### Requirement: 服务层 CreateCcn 方法更新
系统 SHALL 更新 `VpcService.CreateCcn` 方法以支持 `instanceMeteringType` 参数。

#### Scenario: 方法签名扩展
- **WHEN** 更新 `CreateCcn` 方法定义
- **THEN** 方法签名为 `CreateCcn(ctx context.Context, name, description, qos, chargeType, bandWithLimitType, instanceMeteringType string)`
- **AND** `instanceMeteringType` 参数作为最后一个参数添加
- **AND** 保持与现有参数顺序一致的命名风格

#### Scenario: 参数传递到 SDK 请求
- **WHEN** 构造 `CreateCcnRequest` 请求对象
- **THEN** 如果 `instanceMeteringType` 不为空字符串，设置 `request.InstanceMeteringType = &instanceMeteringType`
- **AND** 如果 `instanceMeteringType` 为空字符串，则不设置该字段（保持 nil）
- **AND** 请求日志中包含 `InstanceMeteringType` 参数值

#### Scenario: 响应处理和日志记录
- **WHEN** API 调用成功返回
- **THEN** 记录 DEBUG 级别日志，包含请求和响应的完整 JSON
- **AND** 从响应中提取 CCN ID 并返回给调用者
- **AND** 错误情况下记录 CRITICAL 级别日志，包含请求参数和错误原因

### Requirement: 服务层 DescribeCcns 方法更新
系统 SHALL 更新 `VpcService.DescribeCcns` 方法以支持读取 `InstanceMeteringType` 字段。

#### Scenario: CcnBasicInfo 结构体扩展
- **WHEN** 定义 `CcnBasicInfo` 结构体
- **THEN** 添加 `instanceMeteringType string` 字段
- **AND** 字段命名遵循 camelCase 风格，与其他字段一致
- **AND** 字段位置放在 `bandWithLimitType` 之后，保持逻辑分组

#### Scenario: 解析 API 响应
- **WHEN** 遍历 `response.Response.CcnSet` 中的每个 CCN 实例
- **THEN** 检查 `item.InstanceMeteringType` 是否为 nil
- **AND** 如果不为 nil，设置 `basicInfo.instanceMeteringType = *item.InstanceMeteringType`
- **AND** 如果为 nil，设置 `basicInfo.instanceMeteringType = ""`（空字符串，表示未设置）
- **AND** 将解析后的 `basicInfo` 添加到返回列表

#### Scenario: 向后兼容性处理
- **WHEN** 查询老版本创建的 CCN 实例（可能不返回 InstanceMeteringType）
- **THEN** 系统不因该字段为 nil 而报错
- **AND** 状态中 `instance_metering_type` 字段显示为空或未设置
- **AND** 不影响其他字段的正常读取和显示

### Requirement: 资源 Create 操作实现
系统 SHALL 在 `resourceTencentCloudCcnCreate` 函数中正确处理 `instance_metering_type` 参数。

#### Scenario: 读取用户配置参数
- **WHEN** 执行资源创建操作
- **THEN** 使用 `d.GetOk("instance_metering_type")` 读取用户配置的计量类型
- **AND** 如果用户配置了该参数，将值存储到 `instanceMeteringType` 变量
- **AND** 如果用户未配置，`instanceMeteringType` 为空字符串

#### Scenario: 调用服务层方法
- **WHEN** 调用 `service.CreateCcn` 方法
- **THEN** 传递参数顺序为 `(ctx, name, description, qos, chargeType, bandwidthLimitType, instanceMeteringType)`
- **AND** 所有参数类型为 string
- **AND** 参数值从用户配置或默认值中获取

#### Scenario: 资源 ID 设置和状态刷新
- **WHEN** 创建成功并获得 CCN ID
- **THEN** 使用 `d.SetId(info.ccnId)` 设置资源 ID
- **AND** 调用 `resourceTencentCloudCcnRead` 刷新资源状态
- **AND** Read 操作将从云端读取包括 `instance_metering_type` 在内的所有字段

### Requirement: 资源 Read 操作实现
系统 SHALL 在 `resourceTencentCloudCcnRead` 函数中正确读取和同步 `instance_metering_type` 字段。

#### Scenario: 调用 DescribeCcn 获取资源详情
- **WHEN** 执行 Read 操作刷新状态
- **THEN** 调用 `service.DescribeCcn(ctx, d.Id())` 获取 CCN 详情
- **AND** 从返回的 `info` 结构体中读取 `instanceMeteringType` 字段
- **AND** 处理资源不存在的情况（has == 0 时清空 ID）

#### Scenario: 同步状态到 Terraform
- **WHEN** 成功获取 CCN 详情
- **THEN** 使用 `d.Set("instance_metering_type", info.instanceMeteringType)` 同步状态
- **AND** 同步操作与其他字段（name, description, qos, charge_type 等）一起执行
- **AND** 如果 Set 操作返回错误，记录错误但不中断其他字段的同步

#### Scenario: 处理空值情况
- **WHEN** `info.instanceMeteringType` 为空字符串（老实例未返回该字段）
- **THEN** 仍然执行 `d.Set("instance_metering_type", "")` 将状态设置为空
- **AND** Terraform 状态文件中该字段为 null 或空字符串
- **AND** 不影响资源的其他操作

### Requirement: 资源 Update 操作约束
系统 SHALL 确保 `instance_metering_type` 参数不支持原地更新。

#### Scenario: Update 函数不处理该字段
- **WHEN** 用户修改配置文件中的 `instance_metering_type` 值
- **THEN** `resourceTencentCloudCcnUpdate` 函数不检查该字段的变更
- **AND** Terraform 在 Plan 阶段检测到 ForceNew 字段变更
- **AND** 提示用户该操作将销毁并重新创建资源

#### Scenario: 变更计划显示
- **WHEN** 执行 `terraform plan` 且 `instance_metering_type` 发生变更
- **THEN** 计划输出显示资源将被强制替换（-/+ 符号）
- **AND** 显示旧值和新值
- **AND** 提示原因："forces replacement"

### Requirement: 测试覆盖
系统 SHALL 提供测试用例覆盖 `instance_metering_type` 参数的各种场景。

#### Scenario: 基础创建测试
- **WHEN** 运行 `TestAccTencentCloudCcn_basic` 测试
- **THEN** 测试配置中包含 `instance_metering_type` 参数
- **AND** 验证资源创建成功且 ID 非空
- **AND** 使用 `resource.TestCheckResourceAttr` 验证 `instance_metering_type` 值与配置一致
- **AND** 验证资源销毁后不存在

#### Scenario: 不指定参数测试
- **WHEN** 创建 CCN 资源时不指定 `instance_metering_type`
- **THEN** 资源创建成功
- **AND** 状态中 `instance_metering_type` 为空或云平台默认值
- **AND** 不影响其他字段的功能

#### Scenario: 资源导入测试
- **WHEN** 使用 `terraform import` 导入已存在的 CCN 资源
- **THEN** 导入后状态中包含 `instance_metering_type` 字段
- **AND** 字段值与云平台实际值一致
- **AND** 导入的资源可以正常进行后续操作

### Requirement: 错误处理和日志记录
系统 SHALL 正确处理与 `instance_metering_type` 相关的错误场景。

#### Scenario: 无效参数值处理
- **WHEN** 用户配置了无效的 `instance_metering_type` 值
- **THEN** API 调用返回 InvalidParameter 错误
- **AND** 错误信息清晰指出参数无效
- **AND** 记录 CRITICAL 级别日志，包含请求参数和错误原因

#### Scenario: API 调用失败重试
- **WHEN** CreateCcn API 调用因网络或临时错误失败
- **THEN** 使用 `resource.Retry` 机制进行重试
- **AND** 重试超时时间为 `tccommon.WriteRetryTimeout`
- **AND** 可重试错误使用 `tccommon.RetryError` 包装
- **AND** 不可重试错误直接返回

#### Scenario: 日志记录完整性
- **WHEN** 任何涉及 `instance_metering_type` 的操作执行
- **THEN** 操作开始时记录 `defer tccommon.LogElapsed` 日志
- **AND** API 请求和响应记录 DEBUG 级别日志，包含完整 JSON
- **AND** 错误情况记录 CRITICAL 级别日志，包含上下文信息

### Requirement: 向后兼容性保证
系统 SHALL 确保添加 `instance_metering_type` 参数不破坏现有功能。

#### Scenario: 现有配置不受影响
- **WHEN** 用户使用现有的 CCN 配置（不包含 `instance_metering_type`）
- **THEN** 资源创建和管理功能正常工作
- **AND** 不产生任何错误或警告
- **AND** 行为与未添加参数前完全一致

#### Scenario: 现有测试用例通过
- **WHEN** 运行所有现有的 CCN 相关测试用例
- **THEN** 所有测试用例保持通过状态
- **AND** 测试执行时间无明显增加
- **AND** 无新增的测试失败或不稳定情况

#### Scenario: 已创建资源状态刷新
- **WHEN** 刷新之前创建的 CCN 资源状态（没有 instance_metering_type）
- **THEN** Read 操作成功完成
- **AND** 状态中 `instance_metering_type` 为空或默认值
- **AND** 其他字段正常显示和更新
- **AND** 不触发资源重建或不必要的变更
