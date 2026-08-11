## 1. Schema 定义

- [x] 1.1 在 `tencentcloud/services/teo/resource_tc_teo_config_group_version.go` 中，确认 `version_id` 字段在 schema 中声明为 `Computed: true, TypeString`，Description 说明该字段来源于 `DescribeConfigGroupVersionDetail` 接口的 `Response.ConfigGroupVersionInfo.VersionId`
- [x] 1.2 确认 `version_id` 未设置 `Optional`/`Required`/`ForceNew`（保持只读 Computed 出参语义，不改变资源生命周期）

## 2. Read 函数

- [x] 2.1 在 `resourceTencentCloudTeoConfigGroupVersionRead` 中，确认从 `respData.ConfigGroupVersionInfo.VersionId` 读取并 `_ = d.Set("version_id", respData.ConfigGroupVersionInfo.VersionId)`，且仅在 `respData.ConfigGroupVersionInfo.VersionId != nil` 时调用 set
- [x] 2.2 确认 Read 方法在 `respData == nil` 时，先 `log.Printf("[CRUD] teo_config_group_version id=%s", d.Id())` 保留现场，再 `d.SetId("")`，避免日志中无法定位被清空的 id
- [x] 2.3 确认 Read 方法复用 `service.DescribeTeoConfigGroupVersionById(ctx, zoneId, groupId, versionId)`，复合 ID 按 `tccommon.FILED_SP` 拆分为三段，段数不为 3 时返回错误

## 3. 单元测试

- [x] 3.1 在 `tencentcloud/services/teo/resource_tc_teo_config_group_version_test.go` 中，使用 gomonkey mock `DescribeConfigGroupVersionDetail` 云 API，新增测试用例覆盖 Read 成功时 `version_id` 被正确回填（如 `ver-2kplomhisdcb`）
- [x] 3.2 新增测试用例覆盖 `ConfigGroupVersionInfo` 为 nil 时 Read 不 panic 且正确清空 id（保留 `[CRUD]` 日志）
- [x] 3.3 确保新增单测代码可在当前环境下正确构建（禁止执行 `go test`，保证代码可编译）

## 4. 文档同步

- [x] 4.1 在 `tencentcloud/services/teo/resource_tc_teo_config_group_version.md` 中补充 `version_id` 出参说明（一句话描述中带上所属云产品名称 TEO）
- [x] 4.2 确认 `.md` 文件未手动添加 `Argument Reference`/`Attribute Reference` 部分（由工具自动生成），仅保留一句话描述、Example Usage、Import 部分

## 5. 验证

- [x] 5.1 检查所有函数返回的 error，确保非 `NonRetryableError` 的路径已正确处理；不会出现未使用变量错误
- [x] 5.2 确认 `website/docs/` 下文档不在本阶段手动修改，统一在收尾阶段通过 `make doc` 生成
