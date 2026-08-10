## 1. Schema 修改

- [x] 1.1 在 `tencentcloud/services/tag/resource_tc_tag_attachment.go` 的 `ResourceTencentCloudTagAttachment()` Schema 中，将 `tag_value` 字段的 `ForceNew` 从 `true` 改为 `false`（删除 `ForceNew: true` 行或显式设为 false）
- [x] 1.2 在 `schema.Resource` 中注册 `Update: resourceTencentCloudTagAttachmentUpdate` 回调，同时保留原有 `Create`、`Read`、`Delete`、`Importer` 定义

## 2. service 层实现

**文件**: `tencentcloud/services/tag/service_tencentcloud_tag.go`

- [x] 2.1 新增 `UpdateTagAttachmentTagValue(ctx context.Context, tagKey string, tagValue string, resource string) (errRet error)` 方法
- [x] 2.2 构造 `tag.NewUpdateResourceTagValueRequest()`，设置 `TagKey`、`TagValue`（新值）、`Resource` 参数
- [x] 2.3 使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包裹 API 调用，失败时用 `tccommon.RetryError(e)` 包装错误
- [x] 2.4 添加 request/response body 日志记录（`log.Printf("[DEBUG]..."）和错误处理 defer 日志

## 3. Update 函数实现

**文件**: `tencentcloud/services/tag/resource_tc_tag_attachment.go`

- [x] 3.1 新增 `resourceTencentCloudTagAttachmentUpdate(d *schema.ResourceData, meta interface{}) error` 函数，包含 `defer tccommon.LogElapsed()` 和 `defer tccommon.InconsistentCheck()`
- [x] 3.2 从 `d.Id()` 按 `tccommon.FILED_SP` 拆分出三段（tagKey、旧 tagValue、resource），校验 len==3
- [x] 3.3 使用 `d.HasChange("tag_value")` 检测 tag_value 变更
- [x] 3.4 变更时取 `d.Get("tag_value")` 新值，调用 `service.UpdateTagAttachmentTagValue(ctx, tagKey, newTagValue, resource)`
- [x] 3.5 API 调用成功后（retry 块外）使用新 tag_value 重建 ID：`d.SetId(tagKey + tccommon.FILED_SP + newTagValue + tccommon.FILED_SP + resource)`
- [x] 3.6 函数末尾调用 `resourceTencentCloudTagAttachmentRead(d, meta)` 同步状态

## 4. 单元测试

**文件**: `tencentcloud/services/tag/resource_tc_tag_attachment_test.go`

- [x] 4.1 新增 `TestAccTencentCloudTagAttachmentResource_tagValueUpdate` 测试用例（使用 terraform 测试套件，因为是修改现有资源）
- [x] 4.2 测试步骤覆盖：创建资源 → 修改 tag_value → 验证原地更新（非重建）→ 验证新 tag_value 已生效
- [x] 4.3 添加 ImportState 验证步骤，确认更新后 ID 可正确导入

## 5. 文档更新

- [x] 5.1 更新 `tencentcloud/services/tag/resource_tc_tag_attachment.md`，补充修改 tag_value 的示例场景（创建后修改 tag_value 的用法）

## 6. 代码质量检查与文档生成（收尾阶段）

- [ ] 6.1 执行 `gofmt` 格式化变更的 Go 文件
- [ ] 6.2 执行 `make doc` 生成 `website/docs/` 下的文档
- [ ] 6.3 创建 `.changelog/<PR_NUMBER>.txt` 变更日志条目（enhancement: resource/tencentcloud_tag_attachment: support tag_value in-place update）

## 依赖关系

- 任务 2（service 层）和任务 1（Schema）可并行执行
- 任务 3（Update 函数）依赖任务 1 和任务 2 完成
- 任务 4（测试）依赖任务 3 完成
- 任务 5（文档）可与任务 3、4 并行
- 任务 6（收尾）依赖所有开发任务完成，在 tfpacer-finalize 阶段统一执行
