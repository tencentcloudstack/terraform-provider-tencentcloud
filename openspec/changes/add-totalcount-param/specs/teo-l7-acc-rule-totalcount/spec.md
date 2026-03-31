## ADDED Requirements

### Requirement: TotalCount 字段应在资源 Schema 中定义
tencentcloud_teo_l7_acc_rule 资源 SHALL 在 Schema 中定义 TotalCount 字段，该字段 SHALL 为整型（TypeInt），SHALL 设置为 Optional（可选），SHALL 设置为 Computed（计算字段），SHALL 表示规则总数。

#### Scenario: Schema 定义验证
- **WHEN** Terraform Provider 加载 tencentcloud_teo_l7_acc_rule 资源
- **THEN** 资源 Schema SHALL 包含 TotalCount 字段
- **THEN** TotalCount 字段类型 SHALL 为 TypeInt
- **THEN** TotalCount 字段 SHALL 设置 Optional 为 true
- **THEN** TotalCount 字段 SHALL 设置 Computed 为 true

### Requirement: TotalCount 字段应从 API 响应中读取
在 Read 函数中，Terraform Provider SHALL 从 DescribeL7AccRules API 的响应中读取 TotalCount 字段，SHALL 将该值设置到资源状态中。

#### Scenario: 成功读取 TotalCount
- **WHEN** 用户读取 tencentcloud_teo_l7_acc_rule 资源
- **AND** DescribeL7AccRules API 响应中包含 TotalCount 字段
- **THEN** Provider SHALL 从 API 响应中读取 TotalCount 值
- **THEN** Provider SHALL 将 TotalCount 值设置到资源状态中
- **THEN** TotalCount 值 SHALL 与 API 响应中的值一致

#### Scenario: TotalCount 字段为 null
- **WHEN** 用户读取 tencentcloud_teo_l7_acc_rule 资源
- **AND** DescribeL7AccRules API 响应中 TotalCount 字段为 null
- **THEN** Provider SHALL 设置 TotalCount 值为 0

### Requirement: TotalCount 字段应为只读字段
TotalCount 字段 SHALL 为只读字段，SHALL 仅从 API 响应中获取，SHALL 不需要用户在配置中指定，SHALL 不在 Create/Update/Delete 操作中处理。

#### Scenario: 用户不需要在配置中指定 TotalCount
- **WHEN** 用户创建或更新 tencentcloud_teo_l7_acc_rule 资源
- **THEN** 用户 SHALL 可以不指定 TotalCount 字段
- **THEN** Provider SHALL 不要求用户在配置中提供 TotalCount 值

#### Scenario: TotalCount 字段不参与 Create/Update/Delete
- **WHEN** 用户创建或更新 tencentcloud_teo_l7_acc_rule 资源
- **THEN** Provider SHALL 不在 Create 函数中处理 TotalCount 字段
- **THEN** Provider SHALL 不在 Update 函数中处理 TotalCount 字段
- **THEN** Provider SHALL 不在 Delete 函数中处理 TotalCount 字段

### Requirement: TotalCount 字段应在单元测试中验证
单元测试 SHALL 包含对 TotalCount 字段的验证，SHALL 测试字段的读取逻辑，SHALL 验证字段类型和值。

#### Scenario: 单元测试验证 TotalCount 字段
- **WHEN** 运行 tencentcloud_teo_l7_acc_rule 资源的单元测试
- **THEN** 测试 SHALL 包含 TotalCount 字段的测试用例
- **THEN** 测试 SHALL 验证 TotalCount 字段的类型为整型
- **THEN** 测试 SHALL 验证 TotalCount 字段的值正确
