## 1. 服务层实现

- [x] 1.1 在 `tencentcloud/services/config/service_tencentcloud_config.go` 新增 `DescribeConfigListAggregatorsByFilter` 方法，参照 `DescribeConfigCompliancePacksByFilter` 分页骨架，调用 `ListAggregatorsWithContext`，内部 `Limit=100`/`Offset=0` 自动分页，返回 `items []*configv20220802.Aggregator` 与 `total uint64`
- [x] 1.2 服务层 retry 块内检查 `result == nil || result.Response == nil` 返回 `NonRetryableError`；retry 块内仅调用接口，分页累加与 break 判断放块外

## 2. 数据源实现

- [x] 2.1 创建 `tencentcloud/services/config/data_source_tc_config_list_aggregators.go`，定义 `DataSourceTencentCloudConfigListAggregators()`，schema 包含顶层 `total`（TypeInt，Computed）、`items`（TypeList，Computed，Elem 为 Resource，子字段：`name`/`description`/`owner_uin`/`create_time`/`account_count`/`type`/`account_group_id`/`aggregator_status`/`member_name` 均为 Computed）、`result_output_file`（Optional）
- [x] 2.2 实现 `dataSourceTencentCloudConfigListAggregatorsRead`，retry 块调用 `DescribeConfigListAggregatorsByFilter`；若返回空（items 为 nil 且 total 为 0）返回 `NonRetryableError` 并打印 `[DATASOURCE] read empty, skip SetId`；flatten 列表逐字段判 nil 后 set；`_ = d.Set("items", itemsList)`、`_ = d.Set("total", int(total))`；`d.SetId(helper.BuildToken())`；可选写 result_output_file

## 3. Provider 注册

- [x] 3.1 在 `tencentcloud/provider.go` 数据源 map（config 区块，约 line 1422-1428）新增 `"tencentcloud_config_list_aggregators": config.DataSourceTencentCloudConfigListAggregators()`

## 4. 文档

- [x] 4.1 创建 `tencentcloud/services/config/data_source_tc_config_list_aggregators.md`，包含一句话描述（带产品名 Config）、Example Usage（查询全部账号组示例）；不含 Argument/Attribute Reference（由 make doc 生成）

## 5. 单元测试

- [x] 5.1 创建 `tencentcloud/services/config/data_source_tc_config_list_aggregators_test.go`，使用 gomonkey mock `ListAggregatorsWithContext`，验证 Read 的 flatten 与 set 逻辑（items 列表与 total 正确填充）

## 6. 收尾（由 tfpacer-finalize 执行）

- [ ] 6.1 执行 `gofmt` 格式化新增/修改的 Go 文件
- [ ] 6.2 执行 `make doc` 生成 `website/docs/` 文档与 provider.md 数据源条目
- [ ] 6.3 在 `.changelog/` 下新增 changelog 文件