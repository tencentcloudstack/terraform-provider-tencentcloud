## ADDED Requirements

### Requirement: Resource Schema 定义
资源 SHALL 定义包含以下字段的 Schema：
- `zone_id` (string, Required): 站点 ID
- `namespace` (string, Required): 命名空间名称
- `keys` (list of string, Required): 键名列表，数组长度上限为 20
- `data` (list of object, Computed): 查询结果，包含键值对数据列表
  - `key` (string, Computed): 键名
  - `value` (string, Computed): 键值
  - `expiration` (string, Computed): 过期时间

#### Scenario: Schema 定义验证
- **WHEN** 用户定义资源时提供 zone_id、namespace 和 keys 参数
- **THEN** 资源 SHALL 成功创建并验证参数类型和约束
- **AND** data 字段 SHALL 标记为 Computed

#### Scenario: 键名列表约束
- **WHEN** 用户提供的 keys 列表长度超过 20
- **THEN** 资源 SHALL 返回验证错误

#### Scenario: 键名格式验证
- **WHEN** 用户提供的键名包含非法字符（非字母、数字、中划线、下划线）
- **THEN** 资源 SHALL 返回验证错误

### Requirement: Create 操作
Create 操作 SHALL 调用 EdgeKVGet API 查询指定的键值对数据，并将结果存储到 Terraform State 中。

#### Scenario: 成功创建并查询数据
- **WHEN** 用户执行 terraform apply 创建资源
- **AND** 提供 zone_id、namespace 和 keys 参数
- **THEN** 资源 SHALL 调用 EdgeKVGet API 查询数据
- **AND** 将查询结果存储到 State 中
- **AND** 返回包含 key、value、expiration 的 data 列表

#### Scenario: 查询不存在的键
- **WHEN** 查询的键不存在
- **THEN** data 列表中对应项的 value 字段 SHALL 返回空字符串

#### Scenario: API 调用失败
- **WHEN** EdgeKVGet API 调用失败
- **THEN** 资源 SHALL 返回错误信息给用户

### Requirement: Read 操作
Read 操作 SHALL 调用 EdgeKVGet API 重新查询数据，并更新 Terraform State 中的数据。

#### Scenario: 成功读取并更新数据
- **WHEN** 用户执行 terraform refresh 或 plan
- **THEN** 资源 SHALL 调用 EdgeKVGet API 查询最新数据
- **AND** 更新 State 中的 data 字段

#### Scenario: 数据变化
- **WHEN** 键值在云端被修改
- **THEN** Read 操作 SHALL 检测到变化并更新 State

### Requirement: Update 操作
Update 操作 SHALL 支持修改 zone_id、namespace 或 keys 参数。

#### Scenario: 修改键名列表
- **WHEN** 用户修改 keys 列表
- **THEN** 资源 SHALL 重新调用 EdgeKVGet API 查询新的键值对
- **AND** 更新 State 中的 data 字段

#### Scenario: 修改命名空间
- **WHEN** 用户修改 namespace 参数
- **THEN** 资源 SHALL 使用新的命名空间查询数据

### Requirement: Delete 操作
Delete 操作 SHALL 从 Terraform State 中移除资源，但不会删除云端数据。

#### Scenario: 成功删除资源
- **WHEN** 用户执行 terraform destroy
- **THEN** 资源 SHALL 从 State 中移除
- **AND** 不影响云端的数据

### Requirement: 错误处理
资源 SHALL 实现统一的错误处理机制，包括 API 调用失败、参数验证失败等场景。

#### Scenario: API 错误处理
- **WHEN** EdgeKVGet API 返回错误
- **THEN** 资源 SHALL 返回清晰的错误信息
- **AND** 包含 RequestId 用于问题定位

#### Scenario: 网络超时重试
- **WHEN** API 调用因网络问题超时
- **THEN** 资源 SHALL 进行重试
- **AND** 最终一致性检查 SHALL 确保数据一致性

### Requirement: 资源 ID 格式
资源 SHALL 使用复合 ID 格式：`zoneId#namespace#keysHash`，其中 keysHash 是键名列表的哈希值。

#### Scenario: ID 生成
- **WHEN** 创建资源
- **THEN** 资源 SHALL 生成格式为 `zoneId#namespace#keysHash` 的 ID
- **AND** keysHash SHALL 基于键名列表计算

#### Scenario: ID 解析
- **WHEN** 读取资源
- **THEN** 资源 SHALL 能够从 ID 中解析出 zoneId 和 namespace

### Requirement: 文档和测试
资源 SHALL 包含完整的文档和测试代码。

#### Scenario: 单元测试
- **WHEN** 运行单元测试
- **THEN** 所有测试用例 SHALL 通过
- **AND** 测试 SHALL 覆盖所有 CRUD 操作

#### Scenario: 验收测试
- **WHEN** 运行验收测试（TF_ACC=1）
- **THEN** 资源 SHALL 能够与真实的 TEO API 交互
- **AND** 测试 SHALL 覆盖正常和异常场景

#### Scenario: 文档完整性
- **WHEN** 查看资源文档
- **THEN** 文档 SHALL 包含参数说明、示例和注意事项
