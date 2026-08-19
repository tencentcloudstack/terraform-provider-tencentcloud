# teo-edge-kv-list-datasource Specification

## Purpose
TBD - created by archiving change add-teo-edge-kv-list-datasource. Update Purpose after archive.
## Requirements
### Requirement: Data source schema definition
数据源 `tencentcloud_teo_edge_kv_list` SHALL 定义以下 schema 字段：

- `zone_id` (Required, String): 站点 ID。
- `namespace` (Required, String): 命名空间名称。
- `prefix` (Optional, String): 键名前缀过滤，只返回以指定前缀开头的键名。不填写表示返回所有键名。
- `cursor` (Optional, Computed, String): 游标位置，用于续续遍历。首次查询不填写，从头开始遍历；后续查询填写上一次返回的 cursor 值。Read 末尾 SHALL 用最后一次响应的 Cursor 回填该字段。
- `keys` (Computed, TypeList of String): 键名列表。元素为标量字符串，直接平铺在资源 schema 顶层，不嵌套对象。
- `result_output_file` (Optional, String): 用于将结果保存到文件。

#### Scenario: Schema fields are correctly defined
- **WHEN** user defines a `data "tencentcloud_teo_edge_kv_list"` block in HCL
- **THEN** Terraform SHALL validate that `zone_id` and `namespace` are provided as required fields, and SHALL accept optional `prefix`, `cursor`, and `result_output_file` fields

#### Scenario: cursor is optional and computed
- **WHEN** a user does not set `cursor` in the Terraform configuration
- **THEN** the schema SHALL populate `cursor` from the last API response (Computed behavior) without triggering a plan diff

#### Scenario: cursor accepts user input for pagination continuation
- **WHEN** a user sets `cursor` to a previously returned value
- **THEN** the schema SHALL accept the value and pass it to the EdgeKVList API as the starting position

### Requirement: Data source Read operation queries key list
数据源 Read 方法 SHALL 调用 TEO `EdgeKVList` API。入参映射如下：
- `zone_id` → `request.ZoneId`
- `namespace` → `request.Namespace`
- `prefix` → `request.Prefix`（仅在用户配置非空时设置）
- `cursor` → `request.Cursor`（仅在用户配置非空时设置，首次查询不设置）

Read 方法 SHALL 在 for 循环内部固定 `request.Limit` 为 1000（云 API 注释标注的最大值），每次循环将上一次响应的 `Cursor` 填入 `request.Cursor`，直到响应 `Cursor` 为空字符串（或 nil）时跳出循环。累计的 `Keys` SHALL 一次性 `d.Set("keys", ...)`。循环结束后，Read 方法 SHALL 用最后一次响应的 `Cursor`（若非 nil）回填 `d.Set("cursor", ...)`。

#### Scenario: Query all keys in a namespace
- **WHEN** user provides `zone_id` and `namespace` without `prefix` or `cursor`
- **THEN** the data source SHALL call EdgeKVList with `ZoneId` and `Namespace` set, `Prefix` and `Cursor` unset, and `Limit` set to 1000
- **AND** the data source SHALL loop using the returned `Cursor` until the response `Cursor` is an empty string, accumulating all `Keys` into the `keys` output attribute

#### Scenario: Query keys with prefix filter
- **WHEN** user provides `zone_id`, `namespace`, and `prefix`
- **THEN** the data source SHALL set `request.Prefix` to the configured value and only return keys starting with that prefix

#### Scenario: Continue traversal with cursor
- **WHEN** user provides `zone_id`, `namespace`, and a previously returned `cursor`
- **THEN** the data source SHALL set `request.Cursor` to that value and begin traversal from that position

#### Scenario: Pagination loop terminates on empty cursor
- **WHEN** the EdgeKVList response returns a `Cursor` that is an empty string (or nil)
- **THEN** the data source SHALL stop the pagination loop and set `keys` to the accumulated key list

### Requirement: Retry and empty response handling
数据源 Read 方法 SHALL 使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包装 `EdgeKVList` 调用，并按 `tccommon.RetryError()` 处理错误。在 retry 块内，若 `response == nil || response.Response == nil`，SHALL 返回 `resource.NonRetryableError(...)`，而不直接调用 `d.SetId("")`。在 retry 失败路径上 SHALL 保留 `log.Printf("[DATASOURCE] read empty, skip SetId")` 提示。

#### Scenario: API returns nil response
- **WHEN** the EdgeKVList API returns a nil `response` or nil `response.Response`
- **THEN** the data source SHALL return a `NonRetryableError` inside the retry block instead of clearing the state id

#### Scenario: API returns retryable error
- **WHEN** the EdgeKVList API returns a retryable error
- **THEN** the data source SHALL retry using `tccommon.ReadRetryTimeout` via `resource.Retry` and wrap non-retryable errors with `tccommon.RetryError()`

#### Scenario: Empty keys list is a valid result
- **WHEN** the EdgeKVList API returns an empty `Keys` list but a non-empty `Cursor`
- **THEN** the data source SHALL continue the pagination loop and SHALL NOT treat the empty keys list as an error

### Requirement: Data source ID generation
数据源 Read 方法 SHALL 在成功完成查询后调用 `d.SetId(helper.BuildToken())` 设置 state 占位 ID。

#### Scenario: Set token id on success
- **WHEN** the Read method completes successfully
- **THEN** the data source SHALL set its id to a generated token via `helper.BuildToken()`

### Requirement: Result output file support
数据源 SHALL 支持可选的 `result_output_file` 参数，当用户设置该参数为非空字符串时，SHALL 调用 `tccommon.WriteToFile()` 将结果保存到指定文件。

#### Scenario: Save results to file
- **WHEN** user sets `result_output_file = "./keys.json"`
- **THEN** the data source SHALL write the query results to the specified file in JSON format

#### Scenario: No file written when parameter absent
- **WHEN** user does not set `result_output_file`
- **THEN** the data source SHALL NOT write any file and SHALL only return data to Terraform state

### Requirement: Provider registration
数据源 SHALL 在 `tencentcloud/provider.go` 的 DataSourcesMap 中注册为 `tencentcloud_teo_edge_kv_list`，并在 `tencentcloud/provider.md` 中添加对应条目。

#### Scenario: Data source is accessible
- **WHEN** user references `data "tencentcloud_teo_edge_kv_list"` in their Terraform configuration
- **THEN** Terraform SHALL recognize it as a valid data source type after the provider is initialized

### Requirement: Unit tests with gomonkey mock
数据源 SHALL 提供单元测试文件 `data_source_tc_teo_edge_kv_list_test.go`，使用 gomonkey mock 云 API 调用，覆盖 Read 主流程与空响应路径。

#### Scenario: Unit test for successful read
- **WHEN** unit test for Read is executed with mocked EdgeKVList returning a non-empty keys list and an empty cursor
- **THEN** the test SHALL verify that `keys` is populated correctly and the data source id is set

#### Scenario: Unit test for paginated read
- **WHEN** unit test for Read is executed with mocked EdgeKVList returning keys across multiple pages (first response with non-empty cursor, second response with empty cursor)
- **THEN** the test SHALL verify that keys from all pages are accumulated into the `keys` attribute

#### Scenario: Unit test for prefix filter
- **WHEN** unit test for Read is executed with `prefix` set and mocked EdgeKVList
- **THEN** the test SHALL verify that `request.Prefix` is set to the configured value

#### Scenario: Unit test for nil response returns NonRetryableError
- **WHEN** unit test for Read is executed with mocked EdgeKVList returning a nil response
- **THEN** the test SHALL verify that the Read function returns a NonRetryableError and does not clear the state id

### Requirement: Data source documentation
数据源 SHALL 提供 `data_source_tc_teo_edge_kv_list.md` 文档模板文件，包含一句话描述（提及 TEO）、Example Usage 部分。文档 SHALL NOT 手动添加 `Argument Reference` 和 `Attribute Reference` 部分（由工具自动生成）。

#### Scenario: Documentation file exists with required sections
- **WHEN** the data source is created
- **THEN** a `.md` file SHALL exist with a one-line description mentioning TEO, an Example Usage section demonstrating `zone_id`, `namespace`, and optional `prefix`/`cursor` usage

