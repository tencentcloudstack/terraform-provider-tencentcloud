## 1. Schema 扩展

- [x] 1.1 在 `tencentcloud/services/dts/resource_tc_dts_sync_config.go` 的 `objects.databases.tables[]` Elem schema 中新增 `column_mode` 字段：`schema.TypeString`, `Optional`，对应云API `Table.ColumnMode`
- [x] 1.2 在同一 `tables[]` Elem schema 中新增 `columns` 字段：`schema.TypeList`, `Optional`，Elem 为 `schema.Resource`，含 `column_name`(`schema.TypeString`, `Optional` → `Column.ColumnName`) 和 `new_column_name`(`schema.TypeString`, `Optional` → `Column.NewColumnName`)
- [x] 1.3 在同一 `tables[]` Elem schema 中新增 `tmp_tables` 字段：`schema.TypeSet`, `Elem: schema.TypeString`, `Optional`，对应云API `Table.TmpTables`
- [x] 1.4 在同一 `tables[]` Elem schema 中新增 `table_edit_mode` 字段：`schema.TypeString`, `Optional`，对应云API `Table.TableEditMode`

## 2. Read 路径实现

- [x] 2.1 在 `resourceTencentCloudDtsSyncConfigRead` 中，`databases.Tables` 遍历构建 `tablesMap` 时，当 `tables.ColumnMode` 非 nil 时设置 `tablesMap["column_mode"]`
- [x] 2.2 在同一遍历中，当 `tables.Columns` 非 nil 且非空时，遍历构建 `columnsList`，对每个元素仅当 `ColumnName`/`NewColumnName` 非 nil 时写入 `column_name`/`new_column_name`，并赋值 `tablesMap["columns"]`
- [x] 2.3 在同一遍历中，当 `tables.TmpTables` 非 nil 且非空时，使用 `helper.StringsInterfaces(tables.TmpTables)` 设置 `tablesMap["tmp_tables"]`
- [x] 2.4 在同一遍历中，当 `tables.TableEditMode` 非 nil 时设置 `tablesMap["table_edit_mode"]`

## 3. Write 路径实现

- [x] 3.1 在 `resourceTencentCloudDtsSyncConfigUpdate` 中，构建 `dts.Table` 的循环里，读取 `tablesMap["column_mode"]` 并设置 `table.ColumnMode`
- [x] 3.2 在同一循环里，读取 `tablesMap["columns"]` 列表，遍历构建 `dts.Column` 元素（设置 `ColumnName`/`NewColumnName`），append 到 `table.Columns`
- [x] 3.3 在同一循环里，读取 `tablesMap["tmp_tables"]` 的 `*schema.Set`，转换为 `[]*string` 并赋值给 `table.TmpTables`
- [x] 3.4 在同一循环里，读取 `tablesMap["table_edit_mode"]` 并设置 `table.TableEditMode`

## 4. 文档更新

- [x] 4.1 更新 `tencentcloud/services/dts/resource_tc_dts_sync_config.md` 示例，展示新增的 `column_mode`、`columns`、`tmp_tables`、`table_edit_mode` 字段用法

## 5. 测试

- [x] 5.1 在 `tencentcloud/services/dts/resource_tc_dts_sync_config_test.go` 中补充测试用例，覆盖新增字段的 schema 与读写逻辑（沿用该资源现有测试风格扩展）

## 6. 验证与收尾

- [x] 6.1 代码正确性检查：确认新增字段在 `ConfigureSyncJob` 请求结构 `Table` 和 `DescribeSyncJobs` 响应结构 `Table` 中均存在且路径一致
- [ ] 6.2 收尾阶段通过 `make doc` 重新生成 `website/docs/r/dts_sync_config.html.markdown`，并执行 `gofmt` 格式化（留待收尾阶段执行）
