## 1. Schema 与 CRUD 函数实现

- [x] 1.1 创建资源文件 `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway_line.go`，定义 `ResourceTencentCloudTeoMultiPathGatewayLine()` 函数，包含完整 Schema 定义（zone_id Required+ForceNew, gateway_id Required+ForceNew, line_type Required, line_address Required, proxy_id Optional, rule_id Optional, line_id Computed）和 Import 支持
- [x] 1.2 实现 Create 函数 `resourceTencentCloudTeoMultiPathGatewayLineCreate`，调用 `CreateMultiPathGatewayLine` API，设置复合 ID `zone_id#gateway_id#line_id`，使用 WriteRetryTimeout 进行重试
- [x] 1.3 实现 Read 函数 `resourceTencentCloudTeoMultiPathGatewayLineRead`，从 d.Get() 获取 zone_id 和 gateway_id，从 d.Id() 解析 line_id，调用 `DescribeMultiPathGatewayLine` API，使用 ReadRetryTimeout 进行重试
- [x] 1.4 实现 Update 函数 `resourceTencentCloudTeoMultiPathGatewayLineUpdate`，检测 line_type/line_address/proxy_id/rule_id 变更，调用 `ModifyMultiPathGatewayLine` API，使用 WriteRetryTimeout 进行重试
- [x] 1.5 实现 Delete 函数 `resourceTencentCloudTeoMultiPathGatewayLineDelete`，调用 `DeleteMultiPathGatewayLine` API，使用 WriteRetryTimeout 进行重试

## 2. Service 层实现

- [x] 2.1 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中添加 `DescribeTeoMultiPathGatewayLine` 方法，封装 `DescribeMultiPathGatewayLine` API 调用逻辑

## 3. Provider 注册

- [x] 3.1 在 `tencentcloud/provider.go` 中添加资源注册 `"tencentcloud_teo_multi_path_gateway_line": teo.ResourceTencentCloudTeoMultiPathGatewayLine()`
- [x] 3.2 在 `tencentcloud/provider.md` 中添加资源文档条目

## 4. 文档生成

- [x] 4.1 创建 `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway_line.md` 资源示例文档，包含一句话描述、Example Usage 和 Import 部分

## 5. 单元测试

- [x] 5.1 创建 `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway_line_test.go`，使用 gomonkey mock 云 API 调用，编写 Create、Read、Update、Delete 操作的单元测试用例
- [x] 5.2 使用 `go test -gcflags=all=-l` 运行单元测试，确保所有测试通过
