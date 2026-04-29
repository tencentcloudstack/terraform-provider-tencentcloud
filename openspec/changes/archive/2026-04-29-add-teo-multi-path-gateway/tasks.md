## 1. Schema 与 CRUD 函数实现

- [x] 1.1 创建资源文件 `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway.go`，定义 ResourceTencentCloudTeoMultiPathGateway 函数，包含完整 Schema 定义（zone_id/ForceNew、gateway_type/ForceNew、gateway_name、gateway_port、region_id/ForceNew、gateway_ip、gateway_id/Computed、status/Computed、need_confirm/Computed）和 Import 支持
- [x] 1.2 实现 resourceTencentCloudTeoMultiPathGatewayCreate 函数，调用 CreateMultiPathGateway API，设置联合 ID（zone_id + FILED_SP + gateway_id），检查返回值 GatewayId 是否为空
- [x] 1.3 实现 resourceTencentCloudTeoMultiPathGatewayRead 函数，调用 DescribeMultiPathGateway API，查询网关详情，回写所有字段，处理网关不存在的情况
- [x] 1.4 实现 resourceTencentCloudTeoMultiPathGatewayUpdate 函数，调用 ModifyMultiPathGateway API，传入 zone_id、gateway_id 及可修改字段（gateway_name、gateway_ip、gateway_port）
- [x] 1.5 实现 resourceTencentCloudTeoMultiPathGatewayDelete 函数，调用 DeleteMultiPathGateway API，传入 zone_id 和 gateway_id

## 2. Provider 注册

- [x] 2.1 在 `tencentcloud/provider.go` 中注册 `tencentcloud_teo_multi_path_gateway` 资源
- [x] 2.2 在 `tencentcloud/provider.md` 中添加资源文档链接

## 3. 资源文档

- [x] 3.1 创建 `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway.md` 文件，包含一句话描述（提及 TEO）、Example Usage（含 cloud 和 private 类型示例）和 Import 部分（说明联合 ID 格式）

## 4. 单元测试

- [x] 4.1 创建 `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway_test.go`，使用 gomonkey mock 云 API 调用，测试 Create/Read/Update/Delete 业务逻辑
- [x] 4.2 使用 `go test -gcflags=all=-l` 运行单元测试并确认通过
