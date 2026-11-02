## 1. Schema 定义

- [x] 1.1 在 `tencentcloud/services/dts/data_source_tc_dts_sync_jobs.go` 的 `list.objects.databases` 嵌套 schema 中，新增 `schema_mode` 字段（`schema.TypeString`, `Computed: true`），描述 schema 选择模式（All/Partial，pg 和 sqlserver 使用），放置位置紧邻 `new_schema_name` / `table_mode` 等同级字段。

## 2. Read 函数

- [x] 2.1 在 `data_source_tc_dts_sync_jobs.go` 的 Read 处理函数（`dataSourceTencentCloudDtsSyncJobsRead`）中，`databases[]` 扁平化循环内，紧邻 `databases.ProcedureMode`/`databases.Procedures` 设置逻辑附近，新增 nil 守卫：`if databases.SchemaMode != nil { databasesMap["schema_mode"] = databases.SchemaMode }`，与同级字段保持一致的 nil 检查模式。

## 3. 单元测试

- [x] 3.1 在 `tencentcloud/services/dts/data_source_tc_dts_sync_jobs_test.go` 中，使用 gomonkey mock 云 API（`DescribeDtsSyncJobsByFilter`）的方式新增单元测试用例，构造包含 `Database.SchemaMode` 非空值的 mock 响应，验证 Read 后 `list.objects.databases.schema_mode` 被正确设置；并补充 `SchemaMode` 为 nil 时不设置该字段的用例。

## 4. 文档

- [x] 4.1 通过收尾阶段的 `make doc` 命令自动生成 `website/docs/` 下对应数据源文档（含新增 `schema_mode` 字段），禁止手动编写 website/ 下的文件。

## 5. 验证

- [x] 5.1 执行 `gofmt` 格式化变更的 Go 代码（收尾阶段 tfpacer-finalize skill 执行）。
- [x] 5.2 确认变更仅新增了一个 `schema_mode` 参数，未涉及其他参数改动。