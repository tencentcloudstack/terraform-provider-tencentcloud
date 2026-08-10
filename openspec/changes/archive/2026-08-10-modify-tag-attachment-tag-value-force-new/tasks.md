## 1. Schema 修改

- [x] 1.1 在 `tencentcloud/services/tag/resource_tc_tag_attachment.go` 中将 `tag_value` 字段的 `ForceNew: true` 改为 `ForceNew: false`
- [x] 1.2 在 `ResourceTencentCloudTagAttachment()` 的 `schema.Resource` 中注册 `Update: resourceTencentCloudTagAttachmentUpdate`
- [x] 1.3 保持 `tag_key` 和 `resource` 字段的 `ForceNew: true` 不变

## 2. 服务层实现

- [x] 2.1 在 `tencentcloud/services/tag/service_tencentcloud_tag.go` 中新增 `UpdateTagTagAttachment(ctx context.Context, tagKey string, newTagValue string, resource string) (errRet error)` 方法
- [x] 2.2 构造 `UpdateResourceTagValueRequest`，设置 TagKey、TagValue、Resource 字段
- [x] 2.3 调用 `me.client.UseTagClient().UpdateResourceTagValue(request)`，记录请求和响应日志，处理错误

## 3. Update 函数实现

- [x] 3.1 在 `tencentcloud/services/tag/resource_tc_tag_attachment.go` 中新增 `resourceTencentCloudTagAttachmentUpdate` 函数，包含 `defer tccommon.LogElapsed()` 和 `defer tccommon.InconsistentCheck()`
- [x] 3.2 从 `d.Id()` 解析出 tagKey、旧 tagValue、resourceId（使用 `tccommon.FILED_SP` 分隔）
- [x] 3.3 检测 `d.HasChange("tag_value")`，获取旧值和新值
- [x] 3.4 使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包装 `service.UpdateTagTagAttachment` 调用，retry 块内仅执行 API 调用，失败时使用 `tccommon.RetryError()` 包装错误
- [x] 3.5 retry 块外、错误处理后，使用新 tag_value 重新设置 ID：`d.SetId(tagKey + tccommon.FILED_SP + newTagValue + tccommon.FILED_SP + resourceId)`
- [x] 3.6 调用 `resourceTencentCloudTagAttachmentRead(d, meta)` 同步状态并返回

## 4. 单元测试

- [x] 4.1 在 `tencentcloud/services/tag/resource_tc_tag_attachment_test.go` 中新增 Update 相关单元测试用例，使用 gomonkey mock 云 API（mock `UpdateResourceTagValue`、`GetResources` 等）
- [x] 4.2 测试场景：tag_value 变更时正确调用 UpdateResourceTagValue 并更新 ID
- [x] 4.3 测试场景：tag_value 未变更时不调用 UpdateResourceTagValue
- [x] 4.4 测试场景：API 调用失败时返回错误且不更新 ID

## 5. 文档更新

- [x] 5.1 更新 `tencentcloud/services/tag/resource_tc_tag_attachment.md`，补充 tag_value 可修改的说明

## 6. 验证

- [x] 6.1 确认代码编译通过（由后续流程的 go build 验证）
- [x] 6.2 确认单元测试用例可在当前环境下正确构建（由后续流程验证）
