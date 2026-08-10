# tag-attachment-tag-value-update Specification

## Purpose
TBD - created by archiving change modify-tag-attachment-tag-value-update. Update Purpose after archive.
## Requirements
### Requirement: tag_value 字段支持原地更新
`tencentcloud_tag_attachment` 资源的 `tag_value` 字段 SHALL 支持原地更新（`ForceNew: false`），修改标签值时不触发资源删除重建，而是调用云 API `UpdateResourceTagValue` 在标签键不变的情况下修改标签值。

#### Scenario: 修改 tag_value 触发原地更新
- **GIVEN** 一个已创建的 `tencentcloud_tag_attachment` 资源，tag_key 为 "env"，tag_value 为 "dev"，resource 为某云资源六段式
- **WHEN** 用户修改 Terraform 配置中的 `tag_value` 为 "prod"
- **THEN** Provider SHALL 调用 `UpdateResourceTagValue` API，传入 TagKey="env"、TagValue="prod"、Resource=原资源六段式
- **AND** 不 SHALL 调用 `DeleteResourceTag` 或 `AddResourceTag`
- **AND** 更新成功后 SHALL 用新 tag_value 重建资源 ID（tagKey#新tagValue#resource）

#### Scenario: tag_value 未变更时不调用更新 API
- **GIVEN** 一个已创建的 `tencentcloud_tag_attachment` 资源
- **WHEN** 用户执行 terraform apply 但 `tag_value` 未发生变更
- **THEN** Provider SHALL 不调用 `UpdateResourceTagValue` API

#### Scenario: tag_value 字段 ForceNew 为 false
- **GIVEN** `tencentcloud_tag_attachment` 资源的 Schema 定义
- **WHEN** 检查 `tag_value` 字段属性
- **THEN** 该字段 SHALL 为 `Required`、`TypeString`
- **AND** SHALL 不包含 `ForceNew: true`（即 ForceNew 为 false）

### Requirement: tag_key 与 resource 字段保持 ForceNew
`tencentcloud_tag_attachment` 资源的 `tag_key` 和 `resource` 字段 SHALL 保持 `ForceNew: true`，修改这两个字段仍触发资源删除重建。

#### Scenario: 修改 tag_key 触发重建
- **GIVEN** 一个已创建的 `tencentcloud_tag_attachment` 资源
- **WHEN** 用户修改 `tag_key` 字段
- **THEN** Provider SHALL 触发资源删除并重新创建（ForceNew 行为）

#### Scenario: 修改 resource 触发重建
- **GIVEN** 一个已创建的 `tencentcloud_tag_attachment` 资源
- **WHEN** 用户修改 `resource` 字段
- **THEN** Provider SHALL 触发资源删除并重新创建（ForceNew 行为）

### Requirement: Update 函数注册
`tencentcloud_tag_attachment` 资源 SHALL 注册 `Update` 回调函数 `resourceTencentCloudTagAttachmentUpdate`。

#### Scenario: 资源定义包含 Update 回调
- **GIVEN** `ResourceTencentCloudTagAttachment()` 返回的 `schema.Resource`
- **WHEN** 检查资源定义
- **THEN** SHALL 包含 `Update: resourceTencentCloudTagAttachmentUpdate`
- **AND** SHALL 同时保留原有的 `Create`、`Read`、`Delete`、`Importer` 定义

### Requirement: 更新成功后同步资源 ID
`resourceTencentCloudTagAttachmentUpdate` SHALL 在 `UpdateResourceTagValue` API 调用成功后，使用新的 tag_value 重建复合资源 ID。

#### Scenario: 更新后 ID 包含新 tag_value
- **GIVEN** 资源当前 ID 为 "env#dev#qcs::cvm:ap-guangzhou:uin/123:instance/ins-xxx"
- **WHEN** 用户将 tag_value 从 "dev" 修改为 "prod" 且 `UpdateResourceTagValue` API 调用成功
- **THEN** Provider SHALL 将资源 ID 更新为 "env#prod#qcs::cvm:ap-guangzhou:uin/123:instance/ins-xxx"
- **AND** SHALL 调用 Read 函数用新 ID 同步资源状态

### Requirement: service 层封装 UpdateResourceTagValue
`TagService` SHALL 提供 `UpdateTagAttachmentTagValue` 方法，封装云 API `UpdateResourceTagValue` 调用。

#### Scenario: 调用 UpdateResourceTagValue API
- **GIVEN** tag_key、新 tag_value、resource 六段式
- **WHEN** 调用 `TagService.UpdateTagAttachmentTagValue(ctx, tagKey, tagValue, resource)`
- **THEN** SHALL 构造 `UpdateResourceTagValueRequest`，设置 TagKey、TagValue、Resource 参数
- **AND** SHALL 调用 `UpdateResourceTagValue` API
- **AND** SHALL 使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 进行重试
- **AND** API 失败时 SHALL 使用 `tccommon.RetryError(e)` 包装错误

#### Scenario: API 调用成功
- **GIVEN** 一个有效的标签关联关系
- **WHEN** 调用 `UpdateResourceTagValue` 且 API 返回成功
- **THEN** SHALL 返回 nil（无错误）
- **AND** SHALL 记录 request/response body 日志

#### Scenario: API 调用失败时重试
- **GIVEN** `UpdateResourceTagValue` API 返回可重试错误（如网络错误、限流）
- **WHEN** 调用 `UpdateTagAttachmentTagValue`
- **THEN** SHALL 通过 `resource.Retry` 机制在 `tccommon.WriteRetryTimeout` 超时时间内重试
- **AND** 重试耗尽后 SHALL 返回错误

### Requirement: 文档更新
`tencentcloud_tag_attachment` 资源的 `.md` 示例文件 SHALL 更新，补充 tag_value 更新场景的示例。

#### Scenario: 文档包含更新示例
- **GIVEN** `tencentcloud/services/tag/resource_tc_tag_attachment.md` 文件
- **WHEN** 查看文档内容
- **THEN** SHALL 包含修改 tag_value 的使用示例

