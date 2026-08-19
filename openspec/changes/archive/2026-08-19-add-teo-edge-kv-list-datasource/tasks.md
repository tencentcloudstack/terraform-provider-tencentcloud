## 1. 数据源实现

- [x] 1.1 创建 `tencentcloud/services/teo/data_source_tc_teo_edge_kv_list.go` 文件
- [x] 1.2 定义数据源 Schema：`zone_id`(Required, String)、`namespace`(Required, String)、`prefix`(Optional, String)、`cursor`(Optional, Computed, String)、`keys`(Computed, TypeList of String)、`result_output_file`(Optional, String)
- [x] 1.3 实现 `DataSourceTencentCloudTeoEdgeKvList()` 函数返回 `*schema.Resource`，Read 指向 `dataSourceTencentCloudTeoEdgeKvListRead`
- [x] 1.4 实现 `dataSourceTencentCloudTeoEdgeKvListRead()` 函数：
  - `defer tccommon.LogElapsed("data_source.tencentcloud_teo_edge_kv_list.read")()` 与 `defer tccommon.InconsistentCheck(d, meta)()`
  - 通过 `meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client()` 获取客户端（与 edge_kv 系列资源保持一致）
  - 构造 `teo.NewEdgeKVListRequest()`，从 schema 读取 `zone_id`/`namespace`（必填）、`prefix`/`cursor`（可选，非空时设置）填入 request
  - for 循环内固定 `request.Limit = helper.IntInt64(1000)`，`resource.Retry(tccommon.ReadRetryTimeout, ...)` 包装 `client.EdgeKVListWithContext(ctx, request)` 调用
  - retry 块内检查 `response == nil || response.Response == nil`，返回 `resource.NonRetryableError(errors.New(...))`；失败路径保留 `log.Printf("[DATASOURCE] read empty, skip SetId, ...")`
  - 错误用 `tccommon.RetryError(e)` 包装
  - 累计 `response.Response.Keys` 到 `keys` 切片；将 `response.Response.Cursor` 填入下次 `request.Cursor`
  - 当响应 `Cursor` 为空字符串（或 nil）时跳出循环
  - 循环结束后：非 nil 时 `_ = d.Set("cursor", *lastCursor)`，`_ = d.Set("keys", keysList)`
  - `d.SetId(helper.BuildToken())`
  - 处理 `result_output_file`：非空时调用 `tccommon.WriteToFile(output.(string), d)`
- [x] 1.5 在 `set` 字段前判断 response 对应字段是否为 nil，keys/cursor 仅在非 nil 时 set

## 2. Provider 注册

- [x] 2.1 在 `tencentcloud/provider.go` 的 DataSourcesMap 中按字母序注册 `"tencentcloud_teo_edge_kv_list": teo.DataSourceTencentCloudTeoEdgeKvList()`
- [x] 2.2 在 `tencentcloud/provider.md` 的数据源列表中按字母序添加 `tencentcloud_teo_edge_kv_list` 条目

## 3. 文档模板

- [x] 3.1 创建 `tencentcloud/services/teo/data_source_tc_teo_edge_kv_list.md` 文件
- [x] 3.2 添加一句话描述（Use this data source to query ...，提及 TEO EdgeKV）
- [x] 3.3 添加 Example Usage 部分，演示 `zone_id`、`namespace` 必填用法，以及 `prefix`、`cursor` 可选用法
- [x] 3.4 确认不手动添加 `Argument Reference` 和 `Attribute Reference` 部分（由 `make doc` 自动生成）

## 4. 单元测试（gomonkey mock）

- [x] 4.1 创建 `tencentcloud/services/teo/data_source_tc_teo_edge_kv_list_test.go` 文件
- [x] 4.2 实现 `TestTeoEdgeKvListDataSource_ReadSuccess`：mock `EdgeKVListWithContext` 返回单页非空 keys + 空 cursor，验证 `keys` 正确填充且 id 已设置
- [x] 4.3 实现 `TestTeoEdgeKvListDataSource_Paginated`：mock `EdgeKVListWithContext` 首次返回 keys + 非空 cursor，第二次返回 keys + 空 cursor，验证多页 keys 被累计
- [x] 4.4 实现 `TestTeoEdgeKvListDataSource_PrefixFilter`：mock `EdgeKVListWithContext`，验证 `request.Prefix` 被正确设置
- [x] 4.5 实现 `TestTeoEdgeKvListDataSource_NilResponse`：mock `EdgeKVListWithContext` 返回 nil response，验证 Read 返回 NonRetryableError 且不清空 state id
- [x] 4.6 确保测试使用 gomonkey mock 云 API，不依赖真实云环境；检查所有 error 返回值（必定不出错的用 `_ =` 赋值）

## 5. 代码正确性检查

- [x] 5.1 对照 `vendor/.../teo/v20220901/models.go` 中 `EdgeKVListRequest`/`EdgeKVListResponseParams` 字段，确认入参（ZoneId/Namespace/Prefix/Cursor/Limit）与出参（Keys/Cursor）映射正确
- [x] 5.2 确认 `EdgeKVList` 为同步接口，无需轮询 Read 接口
- [x] 5.3 确认 retry 块内仅调用 API 与空响应检查，未在 retry 块内执行 set id 等成功操作
- [x] 5.4 确认未引入 `_extension.go` 文件、未在 go 文件开头添加注释、日志中统一使用 `teo_edge_kv_list` 资源名

## 6. 验证（收尾阶段执行）

- [ ] 6.1 运行 `gofmt` 格式化 `data_source_tc_teo_edge_kv_list.go` 与 `data_source_tc_teo_edge_kv_list_test.go`
- [ ] 6.2 执行 `make doc` 生成 `website/docs/d/teo_edge_kv_list.html.markdown` 文档
- [ ] 6.3 检查生成的 website 文档格式正确，包含参数说明与示例
