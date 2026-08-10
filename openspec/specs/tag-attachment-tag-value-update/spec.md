# tag-attachment-tag-value-update Specification

## Purpose
TBD - created by archiving change modify-tag-attachment-tag-value-force-new. Update Purpose after archive.
## Requirements
### Requirement: tag_value 字段支持原地更新
`tencentcloud_tag_attachment` 资源的 `tag_value` 字段 SHALL 支持原地更新，不再触发资源重建。当 `tag_value` 发生变化时，Provider SHALL 调用 `UpdateResourceTagValue` API 修改资源已关联的标签值（标签键不变）。

#### Scenario: 用户修改 tag_value 触发原地更新
- **GIVEN** 一个已创建的 `tencentcloud_tag_attachment` 资源，tag_key 为 "env"，tag_value 为 "dev"，resource 为某资源六段式
- **WHEN** 用户修改 Terraform 配置中的 `tag_value` 为 "prod"
- **THEN** Provider 应调用 `UpdateResourceTagValue` API，传入 TagKey="env"、TagValue="prod"、Resource=原资源六段式
- **AND** 不应调用 DeleteResourceTag 或 AddResourceTag
- **AND** 更新成功后应重新设置资源 ID 为 `tagKey + FILED_SP + newTagValue + FILED_SP + resourceId`
- **AND** 调用 Read 操作同步最新状态

#### Scenario: tag_value 未变化时不调用更新 API
- **GIVEN** 一个已创建的 `tencentcloud_tag_attachment` 资源
- **WHEN** 用户执行 `terraform apply` 但未修改 `tag_value`
- **THEN** Provider 不应调用 `UpdateResourceTagValue` API

#### Scenario: tag_value 字段不包含 ForceNew
- **GIVEN** `tencentcloud_tag_attachment` 资源的 Schema 定义
- **WHEN** 检查 `tag_value` 字段属性
- **THEN** 该字段应为 `Required`、`TypeString`
- **AND** 不应包含 `ForceNew: true`

### Requirement: tag_key 和 resource 字段保持 ForceNew
`tencentcloud_tag_attachment` 资源的 `tag_key` 和 `resource` 字段 SHALL 保持 `ForceNew: true`，修改这两个字段仍触发资源重建。

#### Scenario: 修改 tag_key 触发重建
- **GIVEN** 一个已创建的 `tencentcloud_tag_attachment` 资源
- **WHEN** 用户修改 `tag_key` 字段
- **THEN** Provider 应触发资源重建（先删除后创建）

#### Scenario: 修改 resource 触发重建
- **GIVEN** 一个已创建的 `tencentcloud_tag_attachment` 资源
- **WHEN** 用户修改 `resource` 字段
- **THEN** Provider 应触发资源重建（先删除后创建）

### Requirement: 服务层 UpdateTagTagAttachment 方法
TagService SHALL 提供 `UpdateTagTagAttachment` 方法，封装 `UpdateResourceTagValue` API 调用。

#### Scenario: 调用 UpdateResourceTagValue API
- **GIVEN** tag_key、new_tag_value、resource 参数
- **WHEN** 调用 `service.UpdateTagTagAttachment(ctx, tagKey, newTagValue, resource)`
- **THEN** 应构造 `UpdateResourceTagValueRequest`，设置 TagKey、TagValue、Resource 字段
- **AND** 调用 `UpdateResourceTagValue` API
- **AND** 记录请求和响应日志

#### Scenario: API 调用失败处理
- **GIVEN** `UpdateResourceTagValue` API 返回错误
- **WHEN** 错误为可重试类型
- **THEN** 使用 `tccommon.RetryError()` 包装错误
- **AND** 由外层 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 进行重试

### Requirement: Update 操作的重试机制
`resourceTencentCloudTagAttachmentUpdate` 函数 SHALL 使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包装 `UpdateResourceTagValue` API 调用，retry 块内仅执行 API 调用，不执行 SetId 等成功操作。

#### Scenario: 重试块内仅执行 API 调用
- **GIVEN** Update 函数中的 retry 块
- **WHEN** 执行 retry 闭包
- **THEN** retry 块内应仅调用 `UpdateResourceTagValue` API
- **AND** 不应在 retry 块内执行 `d.SetId()` 等成功后操作
- **AND** SetId 操作应在 retry 块外、错误处理后执行
