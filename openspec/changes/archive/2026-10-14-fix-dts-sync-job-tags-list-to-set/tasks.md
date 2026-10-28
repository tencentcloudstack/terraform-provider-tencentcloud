## 1. Schema 修改

- [x] 1.1 在 `tencentcloud/services/dts/resource_tc_dts_sync_job.go` 中，将 `tags` 字段的 `Type` 从 `schema.TypeList` 改为 `schema.TypeSet`（约第 72 行），保持 `Optional`/`Computed`/`ForceNew` 及 `Elem: &schema.Resource{...}`（含 `tag_key`、`tag_value`）不变
- [x] 1.2 确认无需显式设置 `Set` 字段，使用默认 `schema.HashResource` 哈希 `tag_key` + `tag_value` 组合

## 2. Create handler 修改

- [x] 2.1 在 `resourceTencentCloudDtsSyncJobCreate` 函数中，将读取 `tags` 的类型断言从 `v.([]interface{})` 改为 `v.(*schema.Set).List()`（约第 169 行），循环体内构建 `dts.TagItem` 的逻辑不变

## 3. Read handler 确认

- [x] 3.1 确认 `resourceTencentCloudDtsSyncJobRead` 中的 `d.Set("tags", tagsList)` 逻辑（约第 283–297 行）对 TypeSet 兼容，`tagsList` 仍为 `[]interface{}` of `map[string]interface{}`，无需修改

## 4. 测试修改

- [x] 4.1 在 `tencentcloud/services/dts/resource_tc_dts_sync_job_test.go` 中，将 `resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_job.sync_job", "tags.#")` 替换为 `resource.TestCheckTypeSetElemNestedAttrs`，验证 tag 的 `tag_key`/`tag_value` 键值对（与顺序无关）
- [x] 4.2 确认 HCL 配置块 `testAccDtsSyncJob` 中的 `tags { ... }` 块语法无需修改（TypeList/TypeSet 共用块语法）

## 5. 文档更新

- [x] 5.1 确认 `tencentcloud/services/dts/resource_tc_dts_sync_job.md` 示例中的 `tags` 块语法保持不变（Set 与 List 语法相同，无需修改）

## 6. 收尾（tfpacer-finalize 阶段执行）

- [ ] 6.1 对修改的 `.go` 文件执行 `gofmt` 格式化
- [ ] 6.2 执行 `make doc` 生成 `website/docs/` 下的文档
- [ ] 6.3 创建 `.changelog/` 条目记录本次变更

## 7. 验证

- [x] 7.1 确认 `tags` 字段在 schema 中为 TypeSet
- [x] 7.2 确认 Create handler 使用 `(*schema.Set).List()` 读取 tags
- [x] 7.3 确认无其他文件引用 `tags` 字段的 List 类型断言