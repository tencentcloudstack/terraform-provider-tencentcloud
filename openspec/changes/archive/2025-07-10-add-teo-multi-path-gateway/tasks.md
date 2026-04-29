## 1. Schema 定义与 CRUD 函数实现

- [x] 1.1 创建 `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway.go` 文件，定义 ResourceTencentCloudTeoMultiPathGateway 资源，包含 Schema 定义（zone_id/ForceNew, gateway_type/ForceNew, gateway_name, gateway_port, region_id, gateway_ip, gateway_id/Computed, status/Computed, need_confirm/Computed）和 CRUD 函数框架
- [x] 1.2 实现 resourceTencentCloudTeoMultiPathGatewayCreate 函数：调用 CreateMultiPathGateway API，使用 helper.Retry 包裹，设置复合 ID（ZoneId:GatewayId），检查返回的 GatewayId 是否为空
- [x] 1.3 实现 resourceTencentCloudTeoMultiPathGatewayRead 函数：调用 DescribeMultiPathGateways API，使用 helper.Retry 包裹，从返回的 Gateways 列表中匹配 GatewayId，设置各字段值，处理网关不存在的情况
- [x] 1.4 实现 resourceTencentCloudTeoMultiPathGatewayUpdate 函数：调用 ModifyMultiPathGateway API，传入 ZoneId、GatewayId、GatewayName、GatewayIP、GatewayPort
- [x] 1.5 实现 resourceTencentCloudTeoMultiPathGatewayDelete 函数：调用 DeleteMultiPathGateway API，从复合 ID 中解析 ZoneId 和 GatewayId，使用 helper.Retry 包裹

## 2. Provider 注册

- [x] 2.1 在 `tencentcloud/provider.go` 中注册 tencentcloud_teo_multi_path_gateway 资源
- [x] 2.2 在 `tencentcloud/provider.md` 中添加 tencentcloud_teo_multi_path_gateway 资源条目

## 3. 资源文档

- [x] 3.1 创建 `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway.md` 文件，包含一句话描述（提及 TEO）、Example Usage（cloud 和 private 两种类型示例）、Import 部分（说明复合 ID 格式）

## 4. 单元测试

- [x] 4.1 创建 `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway_test.go` 文件，使用 gomonkey mock 云 API 调用，编写 Create、Read、Update、Delete 单元测试用例

## 5. 验证

- [x] 5.1 使用 go test -gcflags=all=-l 运行单元测试，确保所有测试通过
