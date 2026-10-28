## 1. 数据源代码实现

- [x] 1.1 创建 `tencentcloud/services/teo/data_source_tc_teo_billing_data.go`，定义 `tencentcloud_teo_billing_data` 数据源 schema：入参 `start_time`(string,Required)、`end_time`(string,Required)、`zone_ids`(*[]string,Required)、`metric_name`(string,Required)、`interval`(string,Optional)、`filters`(*schema.List,Optional,elem 含 `type`/`value`)、`group_by`(*[]string,Optional)；出参 `id`(string,Computed)、`data`(*schema.List,Computed,elem 含 `time`/`value`/`zone_id`/`host`/`proxy_id`/`region_id`)
- [x] 1.2 实现数据源 Read 函数：在 `helper.Retry()` 块内（使用 `tccommon.ReadRetryTimeout`）调用 `DescribeBillingData`，组装请求参数（StartTime/EndTime/ZoneIds/MetricName/Interval/Filters/GroupBy）；retry 块内仅执行 API 调用
- [x] 1.3 在 Read 的 retry 块内检查空响应（`response == nil` / `response.Response == nil`），若为空返回 `NonRetryableError` 并在失败路径打印 `log.Printf("[DATASOURCE] read empty, skip SetId")`，不直接 `d.SetId("")`
- [x] 1.4 在 retry 块外解析 `response.Response.Data`，将每个 `BillingData` 的 `Time`/`Value`/`ZoneId`/`Host`/`ProxyId`/`RegionId` 映射到 `data` 列表字段 `time`/`value`/`zone_id`/`host`/`proxy_id`/`region_id`（`Data` 为空时 set 空列表）；设置 `d.SetId()` 为 `metric_name#start_time#end_time`（使用 `tccommon.FIELD_SP`）；调用失败用 `tccommon.RetryError()` 包装错误；defer `tccommon.LogElapsed()` 与 `tccommon.InconsistentCheck()`

## 2. Provider 注册

- [x] 2.1 在 `tencentcloud/provider.go` 中注册 `tencentcloud_teo_billing_data` 数据源（参考 `tencentcloud_teo_content_quota` 等现有 TEO 数据源注册方式）
- [x] 2.2 在 `tencentcloud/provider.md` 中添加 `tencentcloud_teo_billing_data` 数据源文档条目

## 3. 文档样例

- [x] 3.1 创建 `tencentcloud/services/teo/data_source_tc_teo_billing_data.md`：一句话描述（"Use this data source to query ..." 带 TEO 产品名）、Example Usage（含完整参数示例，filters 使用合理值）、不包含 Argument/Attribute Reference 部分

## 4. 单元测试

- [x] 4.1 创建 `tencentcloud/services/teo/data_source_tc_teo_billing_data_test.go`，使用 gomonkey mock `DescribeBillingData` 云 API，编写测试用例：mock 返回含多个 `BillingData` 的响应，验证 Read 正确组装请求参数并设置 `data` 输出字段
- [x] 4.2 补充 mock 返回空 `Data` 的测试用例，验证 `data` 设置为空列表且不报错
- [x] 4.3 使用 `go test -gcflags=all=-l` 运行该测试文件，确保测试通过

## 5. 验证

- [x] 5.1 检查所有函数返回的 error 均被处理（必不出错的用 `_ =` 赋值）
- [x] 5.2 确认未修改 `website/` 目录下任何文件（文档由 `make doc` 在收尾阶段生成）
- [x] 5.3 确认未在 `.changelog/` 下新增文件（由收尾阶段处理）
