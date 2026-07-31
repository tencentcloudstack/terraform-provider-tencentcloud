# dbdc-db-custom-regions-datasource Specification

## Purpose
TBD - created by archiving change add-dbdc-db-custom-regions-datasource. Update Purpose after archive.
## Requirements
### Requirement: 数据源 schema 定义 region_set 展平字段
数据源 `tencentcloud_dbdc_db_custom_regions` SHALL 定义一个 `region_set` 计算字段，类型为 `TypeList`，其 `Elem` 为 `&schema.Resource{}`，包含展平后的计算字段：`region` (TypeString)、`region_state` (TypeString)。schema SHALL NOT 在 region 列表外再创建一层包装嵌套结构。

#### Scenario: schema 结构正确展平
- **WHEN** 检查数据源 schema
- **THEN** `region_set` 为 TypeList，其 Elem 包含一个 Resource，该 Resource 含 `region` 和 `region_state` 两个 Computed 字段，处于同一层级

### Requirement: 数据源无用户可见的过滤参数
数据源 SHALL 仅暴露 `result_output_file` 作为可选输入参数。不暴露 limit/offset 给用户，因为 `DescribeDBCustomRegions` API 无入参，无需分页。

#### Scenario: schema 中无分页参数
- **WHEN** 检查数据源 schema
- **THEN** schema 中没有 `offset` 或 `limit` 字段，只有 `result_output_file` 和 `region_set`

### Requirement: Read 函数通过 retry 调用 DescribeDBCustomRegions
Read 函数 SHALL 调用服务层的 `DescribeDBCustomRegions` 方法，并用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包装。来自 API 的错误 SHALL 使用 `tccommon.RetryError(e)` 包装。

#### Scenario: API 调用成功返回 region 列表
- **WHEN** 调用 Read 函数且 DescribeDBCustomRegions API 成功
- **THEN** `region_set` 字段被填充为所有可用 region，且调用 `d.SetId(helper.BuildToken())`

#### Scenario: API 瞬时失败
- **WHEN** DescribeDBCustomRegions API 返回可重试错误
- **THEN** 重试机制继续尝试直到成功或超时

### Requirement: 服务层方法处理 API 调用与重试
服务层方法 `DescribeDBCustomRegions` SHALL 在 `service_tencentcloud_dbdc.go` 中实现，使用 `dbdcv20201029.NewDescribeDBCustomRegionsRequest()` 构建请求（无入参），在 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 块中调用 `me.client.UseDbdcV20201029Client().DescribeDBCustomRegions(request)`。由于该 API 响应无分页字段（无 TotalCount/Limit/Offset），SHALL NOT 实现分页循环，单次调用返回完整 region 列表。

#### Scenario: 单次调用返回所有 region
- **WHEN** 调用 `DescribeDBCustomRegions` 服务方法
- **THEN** 服务方法发出一次 API 请求并返回完整的 RegionSet 列表

#### Scenario: 失败时记录日志
- **WHEN** API 调用失败
- **THEN** 在 defer 中通过 `log.Printf("[CRITAL]%s api[%s] fail ...")` 记录错误，日志中使用资源名称 `dbdc_db_custom_regions`

### Requirement: 空响应返回 NonRetryableError
在服务层中，如果 `DescribeDBCustomRegions` 返回 `nil` response、`nil` Response、或 `nil` 的 RegionSet，方法 SHALL 返回 `resource.NonRetryableError` 并附带描述性错误信息，而不是静默清空数据源 ID。SHALL 在重试失败路径保留 `log.Printf("[DATASOURCE] read empty, skip SetId")`。

#### Scenario: API 返回 nil 响应
- **WHEN** DescribeDBCustomRegions API 返回 nil 或空响应
- **THEN** 返回 NonRetryableError，且数据源 ID 不被清空

### Requirement: 设置每个字段前进行 nil 检查
在 Read 函数中，`RegionInfo` 的每个字段（`Region`、`RegionState`）在写入 Terraform state 前 SHALL 检查是否为 nil。值为 nil 的字段 SHALL NOT 被设置。

#### Scenario: region 含 nil 字段
- **WHEN** API 响应中某个 region 的 RegionState 为 nil
- **THEN** 该 region 条目的 `region_state` 字段不被设置，但其它非 nil 字段（如 `region`）仍被设置

### Requirement: 数据源在 provider.go 和 provider.md 中注册
数据源 SHALL 在 `tencentcloud/provider.go` 的 DataSourcesMap 中注册为 `"tencentcloud_dbdc_db_custom_regions": dbdc.DataSourceTencentCloudDbdcDbCustomRegions()`，并在 `tencentcloud/provider.md` 中列出。

#### Scenario: provider 注册
- **WHEN** provider 初始化
- **THEN** `tencentcloud_dbdc_db_custom_regions` 作为数据源可用

### Requirement: 单元测试使用 gomonkey mock 模式
测试文件 `data_source_tc_dbdc_db_custom_regions_test.go` SHALL 使用 gomonkey mock `DescribeDBCustomRegions` client 方法，复用已有的 `mockMetaDbdcDS`、`ptrStr`、`ptrInt64` 辅助函数。测试 SHALL 验证 schema 结构、Read 函数在 mock 响应下的行为、nil 字段处理、空响应错误处理。测试 SHALL 可通过 `go test -gcflags="all=-l"` 运行。

#### Scenario: 基本 read 测试含 mock 数据
- **WHEN** 调用 Read 函数，mock 的 DescribeDBCustomRegions 返回两个 region
- **THEN** `region_set` 包含两个条目且字段值正确，`d.Id()` 非空

#### Scenario: nil 字段测试
- **WHEN** 调用 Read 函数，mock 的 DescribeDBCustomRegions 返回含 nil RegionState 的 region
- **THEN** 该条目的 `region_state` 未被设置（Terraform SDK 默认空字符串），其它字段正常

#### Scenario: 空响应测试
- **WHEN** 调用 Read 函数，mock 的 DescribeDBCustomRegions 返回空 RegionSet
- **THEN** 返回错误（NonRetryableError 传播）

#### Scenario: schema 测试
- **WHEN** 检查数据源 schema
- **THEN** schema 含 `region_set`（TypeList, Computed）和 `result_output_file`，`region_set` 的 Elem 含 `region` 和 `region_state` 两个 Computed TypeString 字段

### Requirement: .md 文档文件遵循 gendoc 格式
`data_source_tc_dbdc_db_custom_regions.md` 文件 SHALL 包含一句话描述（提及 DBDC 云产品名称），格式为 "Use this data source to query ..."；包含 Example Usage 部分含 HCL 示例；不包含 Argument Reference 和 Attribute Reference 部分（这些由 make doc 自动生成）。该数据源为 RESOURCE_KIND_DATASOURCE 类型，不包含 Import 部分。

#### Scenario: 文档格式
- **WHEN** 生成 .md 文件
- **THEN** 以 "Use this data source to query ..." 开头并提及 DBDC，含 Example Usage，不含 Argument Reference 或 Attribute Reference 部分，不含 Import 部分

