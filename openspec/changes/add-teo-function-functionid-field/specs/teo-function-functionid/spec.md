## ADDED Requirements

### Requirement: FunctionId 字段在资源创建后被正确设置

在成功创建 TEO Function 后，Terraform Provider MUST 从 CreateFunction API 响应中读取 FunctionId 字段，并将其设置到 Terraform state 的 function_id 属性中。

#### Scenario: 成功创建函数后 FunctionId 被设置到 state

- **WHEN** 用户通过 Terraform 配置创建一个新的 tencentcloud_teo_function 资源
- **AND** CreateFunction API 调用成功并返回 FunctionId（例如 "ef-jjhvk7ec"）
- **THEN** Provider MUST 将 FunctionId 从 API 响应中提取
- **AND** Provider MUST 将 FunctionId 设置到 Terraform state 的 function_id 属性中
- **AND** 在 `terraform show` 或 `terraform output` 中应该能够看到 function_id 的值

#### Scenario: FunctionId 字段在 Schema 中正确定义

- **WHEN** tencentcloud_teo_function 资源的 Schema 被加载
- **THEN** Schema MUST 包含 function_id 字段
- **AND** function_id 字段必须设置为 Computed（Computed: true）
- **AND** function_id 字段必须设置为非 Optional（Optional: false）
- **AND** function_id 字段的类型必须为 string

### Requirement: FunctionId 字段在资源读取操作中被正确更新

当读取 TEO Function 资源时，Terraform Provider MUST 从 DescribeFunction 或相关 API 响应中读取 FunctionId 字段，并更新到 Terraform state 中。

#### Scenario: 读取已存在的函数时 FunctionId 被更新到 state

- **WHEN** 用户执行 `terraform refresh` 或 `terraform apply` 来读取现有的 tencentcloud_teo_function 资源
- **AND** DescribeFunction 或相关 API 返回 FunctionId
- **THEN** Provider MUST 从 API 响应中读取 FunctionId
- **AND** Provider MUST 将 FunctionId 更新到 Terraform state 的 function_id 属性中
- **AND** 更新后的 state 应该与 API 返回的 FunctionId 值一致

#### Scenario: Read 函数处理 FunctionId 字段

- **WHEN** tencentcloud_teo_function 资源的 Read 函数被调用
- **THEN** Read 函数 MUST 调用相应的查询 API（如 DescribeFunction）
- **AND** Read 函数 MUST 从 API 响应中提取 FunctionId 字段
- **AND** Read 函数 MUST 使用 `d.Set("function_id", response.FunctionId)` 设置字段值

### Requirement: FunctionId 字段必须保持向后兼容性

新增的 FunctionId 字段 MUST 不能影响现有的 Terraform 配置和 state，确保使用旧版本 Provider 创建的资源仍然可以正常工作。

#### Scenario: 旧版本 state 升级到新版本时正常工作

- **WHEN** 用户使用旧版本 Provider 创建的 tencentcloud_teo_function 资源
- **AND** 该资源的 state 中没有 function_id 字段
- **AND** 用户升级到新版本 Provider
- **AND** 用户执行 `terraform refresh` 或 `terraform apply`
- **THEN** Provider MUST 能够读取该资源的状态
- **AND** 不会因为缺少 function_id 字段而报错
- **AND** Read 函数会自动从 API 获取 FunctionId 并填充到 state 中

#### Scenario: 现有 Terraform 配置无需修改

- **WHEN** 用户的现有 Terraform 配置中不包含 function_id 字段
- **AND** 用户升级到新版本 Provider
- **THEN** 现有配置 MUST 能够正常工作，无需修改
- **AND** `terraform plan` 不会显示任何需要添加 function_id 字段的变更
- **AND** `terraform apply` 不会因为缺少 function_id 字段而失败

### Requirement: FunctionId 字段不能通过用户配置设置

FunctionId 字段 MUST 作为只读计算属性，用户不能在 Terraform 配置中手动指定该字段的值。

#### Scenario: 用户尝试配置 function_id 字段

- **WHEN** 用户尝试在 Terraform 配置中手动设置 function_id 字段
- **THEN** 该配置 MUST 被忽略或报错（取决于字段设置）
- **AND** 实际使用的 FunctionId 值必须来自 API 响应
- **AND** 用户配置的值不应该影响资源的实际状态

### Requirement: 单元测试验证 FunctionId 字段设置

单元测试 MUST 验证 Create 和 Read 函数中 FunctionId 字段的设置逻辑是否正确。

#### Scenario: 单元测试验证 Create 函数的 FunctionId 设置

- **WHEN** 单元测试模拟 CreateFunction API 返回带有 FunctionId 的响应
- **THEN** 测试 MUST 验证 Create 函数正确地调用了 `d.Set("function_id", expectedValue)`
- **AND** 测试必须覆盖不同的 FunctionId 值（包括空值、有效值等）

#### Scenario: 单元测试验证 Read 函数的 FunctionId 设置

- **WHEN** 单元测试模拟 DescribeFunction API 返回带有 FunctionId 的响应
- **THEN** 测试 MUST 验证 Read 函数正确地调用了 `d.Set("function_id", expectedValue)`
- **AND** 测试必须覆盖 API 返回 FunctionId 和不返回 FunctionId 的情况

### Requirement: 验收测试验证 FunctionId 字段集成

验收测试 MUST 验证与真实腾讯云 API 的集成，确保 FunctionId 字段的完整工作流程。

#### Scenario: 验收测试验证创建和读取流程

- **WHEN** 验收测试创建一个真实的 tencentcloud_teo_function 资源
- **THEN** 测试 MUST 验证创建后 state 中包含正确的 function_id
- **AND** 测试 MUST 验证 function_id 的值与腾讯云服务端返回的一致
- **AND** 测试 MUST 验证刷新或再次读取操作时 function_id 保持一致

#### Scenario: 验收测试验证向后兼容性

- **WHEN** 验收测试使用标准的 Terraform 配置创建资源（不包含 function_id）
- **THEN** 测试 MUST 验证资源创建成功
- **AND** 测试 MUST 验证 state 中正确填充了 function_id
- **AND** 测试 MUST 验证后续的 refresh 和 apply 操作正常工作
