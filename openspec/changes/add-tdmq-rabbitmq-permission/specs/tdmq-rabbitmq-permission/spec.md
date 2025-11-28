# TDMQ RabbitMQ 权限管理规范

## ADDED Requirements

### Requirement: RabbitMQ 权限资源定义
系统 SHALL 提供 `tencentcloud_tdmq_rabbitmq_user_permission` 资源，用于管理 RabbitMQ 用户在特定 VirtualHost 下的访问权限。

#### Scenario: 创建权限配置
- **WHEN** 用户定义权限资源并执行 `terraform apply`
- **THEN** 系统调用 `ModifyRabbitMQPermission` API 创建权限配置
- **AND** 返回的资源 ID 格式为 `{instanceId}#{user}#{virtualHost}`
- **AND** 资源状态包含所有配置的权限信息

#### Scenario: 读取权限配置
- **WHEN** 系统需要刷新资源状态
- **THEN** 系统调用 `DescribeRabbitMQPermission` API 查询权限
- **AND** 使用 instance_id, user, virtual_host 作为查询条件
- **AND** 更新资源状态中的权限字段

#### Scenario: 更新权限配置
- **WHEN** 用户修改 config_regexp, write_regexp 或 read_regexp 字段
- **THEN** 系统调用 `ModifyRabbitMQPermission` API 更新权限
- **AND** 不允许修改 instance_id, user, virtual_host（强制重建）

#### Scenario: 删除权限配置
- **WHEN** 用户删除权限资源或执行 `terraform destroy`
- **THEN** 系统调用 `DeleteRabbitMQPermission` API 删除权限
- **AND** 验证权限已成功删除

#### Scenario: 导入已存在权限
- **WHEN** 用户使用 `terraform import` 导入权限
- **THEN** 系统接受格式为 `{instanceId}#{user}#{virtualHost}` 的 ID
- **AND** 系统解析 ID 并调用 Read 操作获取权限详情

### Requirement: 资源 Schema 定义
资源 Schema SHALL 包含以下字段，并遵循指定的约束。

#### Scenario: 必填字段验证
- **WHEN** 用户配置权限资源
- **THEN** instance_id, user, virtual_host, config_regexp, write_regexp, read_regexp 字段必须提供
- **AND** 缺少任何必填字段时返回验证错误

#### Scenario: 字段类型和描述
- **WHEN** 系统定义资源 Schema
- **THEN** 所有字段类型为 String
- **AND** 每个字段都有清晰的 Description 说明其用途
- **AND** instance_id 描述为 "Cluster instance ID"
- **AND** user 描述为 "Username"
- **AND** virtual_host 描述为 "VirtualHost name"
- **AND** config_regexp 描述为 "Configure permission regexp, controls which resources can be declared"
- **AND** write_regexp 描述为 "Write permission regexp, controls which resources can be written"
- **AND** read_regexp 描述为 "Read permission regexp, controls which resources can be read"

#### Scenario: 不可变字段保护
- **WHEN** 用户尝试修改 instance_id, user 或 virtual_host 字段
- **THEN** 系统返回错误："argument `{field_name}` cannot be changed"
- **AND** 建议用户删除并重新创建资源

### Requirement: API 集成实现
系统 SHALL 正确集成腾讯云 TDMQ RabbitMQ 权限管理 API。

#### Scenario: ModifyRabbitMQPermission API 调用
- **WHEN** 创建或更新权限时
- **THEN** 构造请求包含：InstanceId, User, VirtualHost, ConfigRegexp, WriteRegexp, ReadRegexp
- **AND** 使用重试机制处理临时失败（tccommon.WriteRetryTimeout）
- **AND** 记录请求和响应日志（包含 action、request body、response body）

#### Scenario: DescribeRabbitMQPermission API 调用
- **WHEN** 读取权限状态时
- **THEN** 构造请求包含：InstanceId, User (optional), VirtualHost (optional)
- **AND** 从返回的 RabbitMQPermissionList 中匹配对应的权限项
- **AND** 如果未找到权限，设置资源 ID 为空并记录警告日志

#### Scenario: DeleteRabbitMQPermission API 调用
- **WHEN** 删除权限时
- **THEN** 构造请求包含：InstanceId, User, VirtualHost
- **AND** 使用服务层方法 DeleteTdmqRabbitmqPermissionById
- **AND** 处理资源不存在的情况（视为成功）

### Requirement: 服务层方法实现
系统 SHALL 在 TdmqService 中提供权限管理的辅助方法。

#### Scenario: DescribeTdmqRabbitmqPermissionById 方法
- **WHEN** 调用 DescribeTdmqRabbitmqPermissionById(ctx, instanceId, user, virtualHost)
- **THEN** 返回 *tdmq.RabbitMQPermission 或 nil
- **AND** 使用 DescribeRabbitMQPermission API 查询
- **AND** 遍历返回列表匹配 user 和 virtualHost
- **AND** 包含错误处理和重试逻辑

#### Scenario: DeleteTdmqRabbitmqPermissionById 方法
- **WHEN** 调用 DeleteTdmqRabbitmqPermissionById(ctx, instanceId, user, virtualHost)
- **THEN** 调用 DeleteRabbitMQPermission API
- **AND** 返回错误信息或 nil
- **AND** 使用重试机制处理临时失败
- **AND** 记录操作日志

### Requirement: 资源生命周期管理
系统 SHALL 正确实现资源的完整生命周期管理。

#### Scenario: 资源 ID 生成和解析
- **WHEN** 创建资源后生成 ID
- **THEN** 使用 `strings.Join([]string{instanceId, user, virtualHost}, tccommon.FILED_SP)` 生成 ID
- **AND** 读取时使用 `strings.Split(d.Id(), tccommon.FILED_SP)` 解析 ID
- **AND** 验证解析后的分段数量为 3
- **AND** 解析失败时返回格式化错误："id is broken, %s"

#### Scenario: 错误日志记录
- **WHEN** 任何 CRUD 操作执行时
- **THEN** 使用 `defer tccommon.LogElapsed("resource.tencentcloud_tdmq_rabbitmq_user_permission.{operation}")()`
- **AND** 使用 `defer tccommon.InconsistentCheck(d, meta)()`
- **AND** 失败时记录 CRITICAL 级别日志，包含原因
- **AND** 成功时记录 DEBUG 级别日志，包含请求和响应

#### Scenario: 状态一致性检查
- **WHEN** 创建或更新成功后
- **THEN** 调用 Read 方法刷新状态
- **AND** 确保 Terraform 状态与云端实际状态一致

### Requirement: 测试覆盖
系统 SHALL 提供全面的验收测试覆盖。

#### Scenario: 基础创建和销毁测试
- **WHEN** 运行 TestAccTencentCloudTdmqRabbitmqUserPermission_basic
- **THEN** 创建完整的测试资源栈（instance, user, virtual_host, permission）
- **AND** 验证所有字段都正确设置
- **AND** 验证资源 ID 格式正确
- **AND** 测试资源销毁后不存在

#### Scenario: 更新测试
- **WHEN** 运行 TestAccTencentCloudTdmqRabbitmqUserPermission_update
- **THEN** 创建初始权限配置
- **AND** 修改 config_regexp, write_regexp, read_regexp
- **AND** 验证更新后的值正确
- **AND** 确认 instance_id, user, virtual_host 未变化

#### Scenario: 导入测试
- **WHEN** 运行导入测试
- **THEN** 使用格式 `{instanceId}#{user}#{virtualHost}` 的 ID
- **AND** 导入后验证所有字段匹配
- **AND** 确认状态文件正确生成

### Requirement: 文档完整性
系统 SHALL 提供完整、准确的资源文档。

#### Scenario: 资源文档内容
- **WHEN** 用户查看资源文档
- **THEN** 包含功能描述和使用场景
- **AND** 包含完整的参数列表和类型说明
- **AND** 包含示例配置（包括依赖资源）
- **AND** 包含导入示例
- **AND** 说明权限正则表达式的作用

#### Scenario: 示例代码可执行性
- **WHEN** 用户复制文档中的示例
- **THEN** 示例包含所有必需的依赖资源
- **AND** 变量值使用占位符或合理的示例值
- **AND** 示例遵循最佳实践（如使用 ".*" 表示所有资源）

### Requirement: 错误处理和边界情况
系统 SHALL 正确处理各种错误场景和边界情况。

#### Scenario: 资源不存在处理
- **WHEN** 读取不存在的权限配置
- **THEN** 设置资源 ID 为空字符串
- **AND** 记录警告日志
- **AND** 不返回错误（视为已删除）

#### Scenario: API 错误处理
- **WHEN** API 调用返回错误
- **THEN** 使用 tccommon.RetryError 包装可重试错误
- **AND** 返回明确的错误信息给用户
- **AND** 记录详细的错误日志用于排查

#### Scenario: 无效参数验证
- **WHEN** 用户提供无效的 instance_id, user 或 virtual_host
- **THEN** API 返回 InvalidParameter 错误
- **AND** 错误信息清晰指出哪个参数无效
- **AND** 建议用户检查资源是否存在

#### Scenario: ID 格式验证
- **WHEN** 解析资源 ID 失败（分段数量不为 3）
- **THEN** 返回格式化错误："id is broken, {id}"
- **AND** 阻止后续操作执行
