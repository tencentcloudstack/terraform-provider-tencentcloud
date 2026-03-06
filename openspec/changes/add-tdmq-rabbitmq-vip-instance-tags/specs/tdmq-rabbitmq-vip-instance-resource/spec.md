# Spec: TDMQ RabbitMQ VIP Instance Resource

## ADDED Requirements

### Requirement: 标签参数支持
`tencentcloud_tdmq_rabbitmq_vip_instance` 资源 SHALL 支持 `tags` 参数,用于管理实例的资源标签。

#### Scenario: 创建实例时绑定标签
- **GIVEN** 用户在配置中指定了 `tags` 参数
- **WHEN** 执行 `terraform apply` 创建实例
- **THEN** Provider 应将 tags 转换为 `[]*tdmq.Tag` 格式
- **AND** 通过 `CreateRabbitMQVipInstance` API 的 `ResourceTags` 参数传递标签
- **AND** 实例创建成功后标签应绑定到实例

#### Scenario: 读取实例标签
- **GIVEN** 实例已经存在且有标签
- **WHEN** 执行 `terraform refresh` 或 `terraform plan`
- **THEN** Provider 应通过 `DescribeRabbitMQVipInstance` API 获取实例信息
- **AND** 从响应中提取 `Tags` 字段
- **AND** 将 `[]*tdmq.Tag` 转换为 `map[string]string` 格式
- **AND** 同步标签到 Terraform 状态

#### Scenario: 更新实例标签
- **GIVEN** 实例已经存在
- **WHEN** 用户修改配置中的 `tags` 参数并执行 `terraform apply`
- **THEN** Provider 应检测到 `tags` 字段变更
- **AND** 获取最新的 tags 值
- **AND** 将 tags 转换为 `[]*tdmq.Tag` 格式
- **AND** 调用 `ModifyRabbitMQVipInstance` API 并设置 `Tags` 参数(全量替换)
- **AND** 标签更新成功后同步到 Terraform 状态
- **NOTE**: API 参考 https://cloud.tencent.com/document/api/1179/88450,Tags 参数为全量标签更新(非增量)

#### Scenario: 删除实例标签
- **GIVEN** 实例已经有标签
- **WHEN** 用户从配置中移除某些标签键并执行 `terraform apply`
- **THEN** Provider 应识别标签变更
- **AND** 将剩余的标签列表转换为 `[]*tdmq.Tag` 格式
- **AND** 调用 `ModifyRabbitMQVipInstance` API 进行全量替换
- **AND** Terraform 状态应反映最新的标签列表

#### Scenario: 创建实例时不指定标签
- **GIVEN** 用户配置中未指定 `tags` 参数
- **WHEN** 执行 `terraform apply` 创建实例
- **THEN** Provider 应正常创建实例
- **AND** 不传递 `ResourceTags` 参数
- **AND** 实例创建成功且没有标签

#### Scenario: 导入已有实例时同步标签
- **GIVEN** 云端已存在一个带标签的 RabbitMQ VIP 实例
- **WHEN** 用户执行 `terraform import` 导入实例
- **THEN** Provider 应通过 Read 函数获取实例信息
- **AND** 标签应自动同步到 Terraform 状态
- **AND** 后续 `terraform plan` 应显示标签已同步

### Requirement: Schema 定义规范
`tags` 字段 SHALL 遵循 Terraform Provider 的标准 Schema 定义。

#### Scenario: Schema 字段属性
- **GIVEN** 定义 `tags` Schema
- **THEN** 类型应为 `schema.TypeMap`
- **AND** `Optional` 应为 `true`(允许用户不配置)
- **AND** `Computed` 应为 `true`(支持从云端读取)
- **AND** `Description` 应包含标签使用说明和限制链接

### Requirement: 标签格式转换
Provider SHALL 正确处理 Terraform 和腾讯云 API 之间的标签格式转换。

#### Scenario: Terraform 到 API 格式转换
- **GIVEN** Terraform 配置中的 tags 格式为 `map[string]interface{}`
- **WHEN** 调用 Create API
- **THEN** Provider 应使用 `helper.GetTags` 获取标签
- **AND** 转换为 `[]*tdmq.Tag` 格式
- **AND** 每个 Tag 对象包含 `TagKey` 和 `TagValue` 字段

#### Scenario: API 到 Terraform 格式转换
- **GIVEN** API 响应中的 Tags 格式为 `[]*tdmq.Tag`
- **WHEN** 读取实例信息
- **THEN** Provider 应将每个 Tag 对象提取为键值对
- **AND** 转换为 `map[string]string` 格式
- **AND** 设置到 Terraform 状态

### Requirement: 错误处理
Provider SHALL 妥善处理标签相关的错误情况。

#### Scenario: 标签更新调用失败
- **GIVEN** 标签更新操作
- **WHEN** `ModifyRabbitMQVipInstance` API 调用返回错误
- **THEN** Provider 应返回错误信息
- **AND** 错误信息应包含失败原因
- **AND** Terraform 操作应终止,不更新状态
- **NOTE**: 如果当前 SDK 版本不支持 `Tags` 参数,可使用 `svctag.TagService.ModifyTags` 作为备选方案

#### Scenario: 标签格式错误
- **GIVEN** 用户提供的标签不符合腾讯云限制
- **WHEN** 执行 Create 或 Update 操作
- **THEN** API 调用应返回验证错误
- **AND** Provider 应将错误传递给用户
- **AND** 用户应看到清晰的错误提示
