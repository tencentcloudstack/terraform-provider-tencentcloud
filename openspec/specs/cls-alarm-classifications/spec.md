# cls-alarm-classifications Specification

## Purpose
TBD - created by archiving change add-cls-alarm-classifications. Update Purpose after archive.
## Requirements
### Requirement: Schema 中添加 classifications 字段
`tencentcloud_cls_alarm` 资源 MUST 在 schema 中添加 `classifications` 字段,类型为字符串列表,属性为 Optional + Computed,用于支持告警分类管理。

#### Scenario: 定义 classifications 字段
- **WHEN** 定义资源 schema
- **THEN** 包含 `classifications` 字段,类型为 `schema.TypeList`,元素类型为 `schema.TypeString`
- **THEN** 字段属性设置为 `Optional: true` 和 `Computed: true`
- **THEN** 包含描述信息说明该字段用于告警分类

### Requirement: Create 操作支持设置 classifications
创建告警时,如果用户配置了 `classifications` 字段,系统 MUST 将该值传递给 `CreateAlarm` API 的 `Classifications` 参数。

#### Scenario: 用户配置了 classifications
- **WHEN** 用户在配置中设置 `classifications = ["category1", "category2"]`
- **THEN** 调用 `CreateAlarm` API 时,`request.Classifications` 参数包含对应的字符串数组

#### Scenario: 用户未配置 classifications
- **WHEN** 用户配置中未包含 `classifications` 字段
- **THEN** 调用 `CreateAlarm` API 时,不设置 `request.Classifications` 参数

### Requirement: Read 操作读取并同步 classifications
读取告警信息时,系统 MUST 从 `DescribeAlarms` API 响应中获取 `Classifications` 字段并同步到 Terraform state。

#### Scenario: API 返回 classifications 数据
- **WHEN** 调用 `DescribeAlarms` API 返回告警详情
- **THEN** 从响应的 `Classifications` 字段读取值
- **THEN** 使用 `d.Set("classifications", value)` 将值写入 state

#### Scenario: API 返回空 classifications
- **WHEN** API 返回的 `Classifications` 为 nil 或空数组
- **THEN** 不执行 `d.Set` 操作,保持 state 中该字段为空

### Requirement: Update 操作支持修改 classifications
更新告警时,如果 `classifications` 字段发生变化,系统 MUST 将新值传递给 `ModifyAlarm` API 的 `Classifications` 参数。

#### Scenario: 用户修改 classifications
- **WHEN** 用户修改配置中的 `classifications` 字段
- **THEN** `d.HasChange("classifications")` 检测到变更
- **THEN** 调用 `ModifyAlarm` API 时,`request.Classifications` 参数包含新的字符串数组

#### Scenario: classifications 未变更
- **WHEN** 用户更新告警但未修改 `classifications` 字段
- **THEN** `d.HasChange("classifications")` 返回 false
- **THEN** 不在 `ModifyAlarm` 请求中设置 `Classifications` 参数

### Requirement: 向后兼容性
新增的 `classifications` 字段 MUST 不破坏现有的 Terraform 配置和 state 文件。

#### Scenario: 升级 Provider 后现有配置仍可使用
- **WHEN** 用户升级到包含此变更的 Provider 版本
- **THEN** 现有未包含 `classifications` 字段的配置可正常运行
- **THEN** 执行 `terraform plan` 不显示 diff(因字段为 Optional + Computed)

#### Scenario: 新用户可选择性使用新字段
- **WHEN** 新用户创建告警资源
- **THEN** 可以选择配置或不配置 `classifications` 字段
- **THEN** 两种方式都能成功创建告警

### Requirement: 验收测试覆盖
MUST 添加验收测试用例验证 `classifications` 字段的完整 CRUD 流程。

#### Scenario: 测试创建带 classifications 的告警
- **WHEN** 验收测试创建告警并设置 `classifications`
- **THEN** 资源创建成功
- **THEN** state 中 `classifications` 值与配置一致

#### Scenario: 测试更新 classifications
- **WHEN** 验收测试修改告警的 `classifications` 字段
- **THEN** 资源更新成功
- **THEN** state 中反映新的 `classifications` 值

#### Scenario: 测试不配置 classifications
- **WHEN** 验收测试创建告警但不设置 `classifications`
- **THEN** 资源创建成功
- **THEN** state 中 `classifications` 为空或从云端读取的默认值

### Requirement: 文档更新
MUST 在资源文档中添加 `classifications` 字段的说明和使用示例。

#### Scenario: 文档包含字段说明
- **WHEN** 查看 `website/docs/r/cls_alarm.html.markdown`
- **THEN** Argument Reference 部分包含 `classifications` 字段描述
- **THEN** 说明字段类型、是否必填、作用

#### Scenario: 文档包含使用示例
- **WHEN** 查看资源文档的示例部分
- **THEN** 至少有一个示例展示如何配置 `classifications` 字段

