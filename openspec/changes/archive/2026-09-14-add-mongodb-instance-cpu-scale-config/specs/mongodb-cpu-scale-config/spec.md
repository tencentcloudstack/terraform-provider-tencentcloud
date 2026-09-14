## ADDED Requirements

### Requirement: CPU 弹性扩容配置管理
系统 SHALL 提供 `tencentcloud_mongodb_instance_cpu_scale` 配置型资源，用于管理已存在的 MongoDB 实例的 CPU 弹性扩容配置。

#### Scenario: 创建 CPU 弹性扩容配置
- **WHEN** 用户声明 `tencentcloud_mongodb_instance_cpu_scale` 资源，指定 `instance_id` 和 `extra_cpu`
- **THEN** 系统调用 `ScaleUpDBInstanceCpu` API，传入 `instance_id` 和 `extra_cpu`，开启 CPU 弹性扩容
- **THEN** 资源 ID 设置为 `instance_id` 的值

#### Scenario: 读取 CPU 弹性扩容配置
- **WHEN** terraform 执行 Read 操作
- **THEN** 系统调用 `DescribeDBInstances` API 通过 `instance_id` 验证实例存在
- **THEN** 若实例存在，保留当前 state 中的配置参数
- **THEN** 若实例不存在或 API 返回空，将 `d.SetId("")` 表示资源已不存在

#### Scenario: 更新额外 CPU 核数
- **WHEN** 用户修改 `extra_cpu` 的值
- **THEN** 系统调用 `ScaleUpDBInstanceCpu` API，传入新的 `instance_id` 和 `extra_cpu`
- **THEN** 系统通过轮询 Read 等待异步操作生效

#### Scenario: 关闭弹性扩容（移除 extra_cpu）
- **WHEN** 用户从配置中移除 `extra_cpu` 字段
- **THEN** 系统调用 `ScaleDownDBInstanceCpu` API，传入 `instance_id`
- **THEN** 系统通过轮询 Read 等待异步操作生效

#### Scenario: 删除资源
- **WHEN** 用户执行 `terraform destroy` 删除 `tencentcloud_mongodb_instance_cpu_scale` 资源
- **THEN** 系统调用 `ScaleDownDBInstanceCpu` API 关闭 CPU 弹性扩容
- **THEN** 资源从 terraform state 中移除

### Requirement: 资源 Import 支持
系统 SHALL 支持通过 `terraform import` 导入已存在的 MongoDB 实例 CPU 弹性扩容配置。

#### Scenario: 导入已有配置
- **WHEN** 用户执行 `terraform import tencentcloud_mongodb_instance_cpu_scale.example cmgo-xxxxxxxx`
- **THEN** 系统接受 `instance_id` 作为导入 ID
- **THEN** 系统调用 `DescribeDBInstances` API 验证实例存在后完成导入

### Requirement: 异步操作处理
系统 SHALL 在调用异步 API（ScaleUpDBInstanceCpu / ScaleDownDBInstanceCpu）后轮询等待操作生效。

#### Scenario: 异步操作轮询
- **WHEN** 系统调用 ScaleUpDBInstanceCpu 或 ScaleDownDBInstanceCpu 并成功返回 FlowId
- **THEN** 系统使用 `resource.Retry` 和 `ReadRetryTimeout` 配置重试等待
- **THEN** 系统在重试中调用 Read 验证操作结果
- **THEN** 若超时仍不生效，返回错误提示