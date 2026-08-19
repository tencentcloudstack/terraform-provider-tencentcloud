## 1. Schema 定义

- [x] 1.1 在 `resource_tc_dlc_data_engine.go` 的 `Schema` map 中新增 `tags` 字段：类型 `schema.TypeList`，`Optional: true`，`ForceNew: true`，`Description` 说明该参数变更会触发重建
- [x] 1.2 为 `tags` 的 `Elem` (`schema.Resource`) 定义两个子字段：`tag_key`（`schema.TypeString`, `Required: true`）与 `tag_value`（`schema.TypeString`, `Optional: true`），各自补充 `Description`

## 2. 资源 CRUD 实现

- [x] 2.1 在 `resourceTencentCloudDlcDataEngineCreate` 中读取 `tags` 块，遍历构造 `[]*dlc.TagInfo`（`TagKey` 来自 `tag_key`、`TagValue` 来自 `tag_value`），赋值给 `request.Tags`；用户未配置时保持 `request.Tags` 为 nil
- [x] 2.2 在 `resourceTencentCloudDlcDataEngineRead` 中，当 `dataEngine.TagList` 非 nil 时，将其展平为 `[]map[string]interface{}`（仅当 `TagKey`/`TagValue` 非 nil 时写入对应 `tag_key`/`tag_value`），通过 `d.Set("tags", ...)` 同步状态；`TagList` 为 nil/空时不报错
- [x] 2.3 在 `resourceTencentCloudDlcDataEngineUpdate` 的 `immutableArgs` 列表中追加 `"tags"`，与现有 create-only 参数保持一致的不可变约束
- [x] 2.4 确认 `UpdateDataEngineRequest` 无 `Tags` 字段，`Update` 路径不向其传递 tags（由 `ForceNew` 在 plan 阶段处理重建）

## 3. 测试实现

- [x] 3.1 在 `resource_tc_dlc_data_engine_test.go` 中使用 gomonkey mock 云 API，补充 create 路径单测：验证配置 `tags` 后 `request.Tags` 包含期望的 `TagInfo`（正确的 `TagKey`/`TagValue`）
- [x] 3.2 补充 read 路径单测：mock `DescribeDataEngine` 返回带 `TagList` 的 `DataEngineInfo`，验证 `d.Get("tags")` 展平结果正确；并覆盖 `TagList` 为 nil/空时不报错的场景

## 4. 文档更新

- [x] 4.1 更新 `tencentcloud/services/dlc/resource_tc_dlc_data_engine.md`：在 Example Usage 中补充 `tags` 块示例（含 `tag_key`/`tag_value`），不手写 `Argument Reference`/`Attribute Reference`（由 `make doc` 自动生成）
- [x] 4.2 收尾阶段通过 `make doc` 生成 `website/docs/` 下文档（不在本步手动编辑 website/）

## 5. 代码正确性检查

- [x] 5.1 校验新增参数 `tags` 仅出现在 `CreateDataEngine` 入参和 `DescribeDataEngine`/`DescribeDataEngines` 出参中，且 `UpdateDataEngine` 不支持 `Tags`，CRUD 各路径的参数与对应云 API 接口一致
- [x] 5.2 确认 `Read` 中调用 `setXX` 前先判断 `TagList` 及子字段非 nil，符合 nil 检查规则
- [x] 5.3 确认 `Create` 检查云 API 返回值非空逻辑未受影响，未引入新的未检查 error
