## ADDED Requirements

### Requirement: tencentcloud_sms_template 资源支持 status_code 字段
系统 SHALL 在 `tencentcloud_sms_template` 资源中增加 `status_code` 只读属性，用于展示短信模板的审核状态。

#### Scenario: Schema 定义中包含 status_code 字段
- **WHEN** 定义 `tencentcloud_sms_template` 资源 Schema
- **THEN** `status_code` 字段类型为 `schema.TypeInt`
- **AND** 字段为 `Computed: true`（只读输出字段，用户不可配置）
- **AND** `ForceNew: false`（字段变更不会触发资源重建）
- **AND** Description 说明为："Template status. 0: approved and effective, 1: pending review, 2: approved pending activation, -1: review failed or rejected"

#### Scenario: Read 操作从 API 响应中读取 status_code
- **WHEN** 执行 Read 操作刷新资源状态
- **THEN** 从 `DescribeTemplateListStatus.StatusCode` 字段读取模板审核状态码
- **AND** 使用 `d.Set("status_code", value)` 同步到 Terraform 状态
- **AND** 如果 API 响应中该字段为 nil，则不设置该字段（保持 Terraform 状态中为 null）

#### Scenario: 向后兼容性
- **WHEN** 用户升级到包含 `status_code` 字段的新版本 Provider
- **THEN** 现有的 Terraform 配置文件无需任何修改
- **AND** `terraform plan` 不会因新增字段而显示计划外的变更
- **AND** 已有的资源不受影响，Create/Update/Delete 操作行为不变