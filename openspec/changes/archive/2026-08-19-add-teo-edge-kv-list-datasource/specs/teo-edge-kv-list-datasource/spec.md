## ADDED Requirements

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
数据源 Read 方法 SHALL 通过 `TeoService` 的 `DescribeTeoEdgeKvListByFilter(ctx, paramMap)` service 层方法查询键名列表，不在 Read 函数内直接调用云 API。Read 方法 SHALL 将 schema 字段组装为 `paramMap`：
- `zone_id` → `paramMap["zone_id"]`
- `namespace` → `paramMap["namespace"]`
- `prefix` → `paramMap["prefix"]`（仅在用户配置非空时设置）
- `cursor` → `paramMap["cursor"]`（仅在用户配置非空时设置，首次查询不设置）

service 层方法 SHALL 在 for 循环内部固定 `request.Limit` 为 1000（云 API 注释标注的最大值），每次循环将上一次响应的 `Cursor` 填入 `request.Cursor`，直到响应 `Cursor` 为空字符串（或 nil）时跳出循环。service 层方法 SHALL 返回累计的 `Keys` 与最后一次响应的 `Cursor`。Read 方法 SHALL 一次性 `d.Set("keys", ...)`，并在最后一次响应 `Cursor` 非 nil 时 `d.Set("cursor", ...)`。由于 service 层方法内部已用 `resource.Retry` 包装 API 调用，Read 方法 SHALL NOT 再包一层 retry（避免 retry 嵌套）。

#### Scenario: Query all keys in a namespace
- **WHEN** user provides `zone_id` and `namespace` without `prefix` or `cursor`
- **THEN** the data source SHALL assemble `paramMap` with `zone_id` and `namespace` set, `prefix` and `cursor` unset, and call `DescribeTeoEdgeKvListByFilter` which sets `Limit` to 1000
- **AND** the service method SHALL loop using the returned `Cursor` until the response `Cursor` is an empty string, accumulating all `Keys` into the `keys` output attribute

#### Scenario: Query keys with prefix filter
- **WHEN** user provides `zone_id`, `namespace`, and `prefix`
- **THEN** the data source SHALL set `paramMap["prefix"]` to the configured value and the service method SHALL set `request.Prefix`, only returning keys starting with that prefix

#### Scenario: Continue traversal with cursor
- **WHEN** user provides `zone_id`, `namespace`, and a previously returned `cursor`
- **THEN** the data source SHALL set `paramMap["cursor"]` to that value and the service method SHALL set `request.Cursor`, beginning traversal from that position

#### Scenario: Pagination loop terminates on empty cursor
- **WHEN** the EdgeKVList response returns a `Cursor` that is an empty string (or nil)
- **THEN** the service method SHALL stop the pagination loop and return the accumulated key list, and the data source SHALL set `keys` to that list

### Requirement: Retry and empty response handling
service 层方法 `DescribeTeoEdgeKvListByFilter` SHALL 使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包装 `EdgeKVList` 调用，并按 `tccommon.RetryError()` 处理错误。在 retry 块内，若 `result == nil || result.Response == nil`，SHALL 返回 `resource.NonRetryableError(...)`，而不直接调用 `d.SetId("")`。在 retry 失败路径上 SHALL 保留 `log.Printf("[DATASOURCE] read empty, skip SetId")` 提示。Read 方法收到 service 返回的 error 后 SHALL 直接返回，不调用 `d.SetId("")`，从而保留 state id。

#### Scenario: API returns nil response
- **WHEN** the EdgeKVList API returns a nil `result` or nil `result.Response`
- **THEN** the service method SHALL return a `NonRetryableError` inside the retry block, and the Read function SHALL propagate the error without clearing the state id

#### Scenario: API returns retryable error
- **WHEN** the EdgeKVList API returns a retryable error
- **THEN** the service method SHALL retry using `tccommon.ReadRetryTimeout` via `resource.Retry` and wrap non-retryable errors with `tccommon.RetryError()`

#### Scenario: Empty keys list is a valid result
- **WHEN** the EdgeKVList API returns an empty `Keys` list but a non-empty `Cursor`
- **THEN** the service method SHALL continue the pagination loop and SHALL NOT treat the empty keys list as an error

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
