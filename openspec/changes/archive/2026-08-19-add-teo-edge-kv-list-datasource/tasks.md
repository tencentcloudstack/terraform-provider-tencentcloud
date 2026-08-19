## 1. service 层方法实现

- [x] 1.1 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中新增 `DescribeTeoEdgeKvListByFilter(ctx context.Context, param map[string]interface{}) (keys []*string, cursor *string, errRet error)` 方法
- [x] 1.2 从 `param` 读取 `zone_id`/`namespace`/`prefix`/`cursor`（非空时）填入 `request`
- [x] 1.3 for 循环内固定 `request.Limit = helper.IntInt64(1000)`，`resource.Retry(tccommon.ReadRetryTimeout, ...)` 包装 `me.client.UseTeoClient().EdgeKVListWithContext(ctx, request)` 调用
- [x] 1.4 retry 块内检查 `result == nil || result.Response == nil`，返回 `resource.NonRetryableError(...)`；失败路径保留 `log.Printf("[DATASOURCE] read empty, skip SetId, ...")`
- [x] 1.5 错误用 `tccommon.RetryError(e)` 包装
- [x] 1.6 累计 `response.Response.Keys` 到 `keys`；将 `response.Response.Cursor` 填入下次 `request.Cursor`；当 `Cursor` 为空字符串（或 nil）时跳出循环
- [x] 1.7 返回累计的 `keys` 与最后一次的 `cursor`
- [x] 1.8 defer 中 errRet 非空时打印 `[CRITAL]` 日志

## 2. 数据源实现

- [x] 2.1 创建 `tencentcloud/services/teo/data_source_tc_teo_edge_kv_list.go` 文件
- [x] 2.2 定义数据源 Schema：`zone_id`(Required, String)、`namespace`(Required, String)、`prefix`(Optional, String)、`cursor`(Optional, Computed, String)、`keys`(Computed, TypeList of String)、`result_output_file`(Optional, String)
- [x] 2.3 实现 `DataSourceTencentCloudTeoEdgeKvList()` 函数返回 `*schema.Resource`，Read 指向 `dataSourceTencentCloudTeoEdgeKvListRead`
- [x] 2.4 实现 `dataSourceTencentCloudTeoEdgeKvListRead()` 函数：
  - `defer tccommon.LogElapsed("data_source.tencentcloud_teo_edge_kv_list.read")()` 与 `defer tccommon.InconsistentCheck(d, meta)()`
  - 通过 `tccommon.NewResourceLifeCycleHandleFuncContext(...)` 构造 ctx，实例化 `TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}`
  - 组装 `paramMap`（zone_id/namespace 必填，prefix/cursor 可选非空时设置）
  - 调用 `service.DescribeTeoEdgeKvListByFilter(ctx, paramMap)`，**不**再包一层 retry（service 内部已有 retry，避免嵌套）
  - `cursor` 非 nil 时 `_ = d.Set("cursor", *lastCursor)`，`keys` 非 nil 时 `_ = d.Set("keys", keysList)`
  - `d.SetId(helper.BuildToken())`
  - 处理 `result_output_file`：非空时调用 `tccommon.WriteToFile(output.(string), d)`
- [x] 2.5 在 set 字段前判断 response 对应字段是否为 nil，keys/cursor 仅在非 nil 时 set

## 3. Provider 注册

- [x] 3.1 在 `tencentcloud/provider.go` 的 DataSourcesMap 中按字母序注册 `"tencentcloud_teo_edge_kv_list": teo.DataSourceTencentCloudTeoEdgeKvList()`
- [x] 3.2 在 `tencentcloud/provider.md` 的数据源列表中按字母序添加 `tencentcloud_teo_edge_kv_list` 条目

## 4. 文档模板

- [x] 4.1 创建 `tencentcloud/services/teo/data_source_tc_teo_edge_kv_list.md` 文件
- [x] 4.2 添加一句话描述（Use this data source to query ...，提及 TEO EdgeKV）
- [x] 4.3 添加 Example Usage 部分，演示 `zone_id`、`namespace` 必填用法，以及 `prefix`、`cursor` 可选用法
- [x] 4.4 确认不手动添加 `Argument Reference` 和 `Attribute Reference` 部分（由 `make doc` 自动生成）

## 5. 单元测试（gomonkey mock）

- [x] 5.1 创建 `tencentcloud/services/teo/data_source_tc_teo_edge_kv_list_test.go` 文件
- [x] 5.2 mock `UseTeoClient` 返回 `teov20220901.Client{}`，mock `EdgeKVListWithContext` 验证业务逻辑
- [x] 5.3 实现 `TestTeoEdgeKvListDataSource_ReadSuccess`：mock 返回单页非空 keys + 空 cursor，验证 `keys` 正确填充且 id 已设置
- [x] 5.4 实现 `TestTeoEdgeKvListDataSource_Paginated`：mock 首次返回 keys + 非空 cursor，第二次返回 keys + 空 cursor，验证多页 keys 被累计
- [x] 5.5 实现 `TestTeoEdgeKvListDataSource_PrefixFilter`：mock，验证 `request.Prefix` 被正确设置
- [x] 5.6 实现 `TestTeoEdgeKvListDataSource_NilResponse`：mock 返回 nil response，验证 Read 返回错误且不清空 state id
- [x] 5.7 实现 `TestTeoEdgeKvListDataSource_APIError`：mock 返回 SDK error，验证错误被传播
- [x] 5.8 实现 `TestTeoEdgeKvListDataSource_Schema`：验证 schema 字段定义
- [x] 5.9 确保测试使用 gomonkey mock 云 API，不依赖真实云环境；检查所有 error 返回值（必定不出错的用 `_ =` 赋值）

## 6. 代码正确性检查

- [x] 6.1 对照 `vendor/.../teo/v20220901/models.go` 中 `EdgeKVListRequest`/`EdgeKVListResponseParams` 字段，确认入参（ZoneId/Namespace/Prefix/Cursor/Limit）与出参（Keys/Cursor）映射正确
- [x] 6.2 确认 `EdgeKVList` 为同步接口，无需轮询 Read 接口
- [x] 6.3 确认 service 层 retry 块内仅调用 API 与空响应检查，未在 retry 块内执行 set id 等成功操作
- [x] 6.4 确认 Read 函数未在 service 已有 retry 之外再包一层 retry（避免 retry 嵌套）
- [x] 6.5 确认未引入 `_extension.go` 文件、未在 go 文件开头添加注释、日志中统一使用 `teo_edge_kv_list` 资源名

## 7. 验证（收尾阶段执行）

- [ ] 7.1 运行 `gofmt` 格式化 `data_source_tc_teo_edge_kv_list.go`、`data_source_tc_teo_edge_kv_list_test.go` 与 `service_tencentcloud_teo.go`
- [ ] 7.2 执行 `make doc` 生成 `website/docs/d/teo_edge_kv_list.html.markdown` 文档
- [ ] 7.3 检查生成的 website 文档格式正确，包含参数说明与示例
