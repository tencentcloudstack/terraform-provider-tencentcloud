## 1. 服务层实现

- [x] 1.1 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中新增 `DescribeTeoMultiPathGatewaysByFilter` 方法，接收 `context.Context` 和 `map[string]interface{}` 参数，调用 `DescribeMultiPathGateways` API，使用 Limit=1000 进行分页，返回 `[]*teov20220901.MultiPathGateway`

## 2. 数据源 Schema 与 Read 函数实现

- [x] 2.1 创建 `tencentcloud/services/teo/data_source_tc_teo_multi_path_gateway.go`，定义 `DataSourceTencentCloudTeoMultiPathGateway()` 函数，包含 schema 定义（zone_id Required、filters Optional、gateways Computed、result_output_file Optional）和 Read 函数
- [x] 2.2 在 Read 函数中实现：构建 paramMap → 调用服务层方法（含 Retry）→ 扁平化响应到 gateways 列表 → 设置 d.Id() → 写入 result_output_file

## 3. Provider 注册

- [x] 3.1 在 `tencentcloud/provider.go` 的 DataSourcesMap 中添加 `"tencentcloud_teo_multi_path_gateway": teo.DataSourceTencentCloudTeoMultiPathGateway()`
- [x] 3.2 在 `tencentcloud/provider.md` 中添加对应的数据源注册信息

## 4. 文档

- [x] 4.1 创建 `tencentcloud/services/teo/data_source_tc_teo_multi_path_gateway.md`，包含描述和示例用法（按 zone_id 查询、按 zone_id + filters 查询）

## 5. 单元测试

- [x] 5.1 创建 `tencentcloud/services/teo/data_source_tc_teo_multi_path_gateway_test.go`，使用 gomonkey mock 云API，编写 Read 函数的单元测试

## 6. 验证

- [x] 6.1 使用 `go test -gcflags=all=-l` 运行单元测试确保通过
