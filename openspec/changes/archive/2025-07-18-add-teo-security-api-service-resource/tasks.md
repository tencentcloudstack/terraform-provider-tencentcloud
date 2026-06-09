## 1. Schema 定义与 CRUD 函数实现

- [x] 1.1 创建资源文件 `tencentcloud/services/teo/resource_tc_teo_security_api_service.go`，定义 `ResourceTencentCloudTeoSecurityApiService()` 函数，包含完整 schema（zone_id、api_services、api_service_ids、api_resources）和 CRUD + Import 支持
- [x] 1.2 实现 Create 函数 `resourceTencentCloudTeoSecurityApiServiceCreate`，调用 `CreateSecurityAPIService` API，设置复合 ID（zone_id + api_service_ids），处理重试逻辑
- [x] 1.3 实现 Read 函数 `resourceTencentCloudTeoSecurityApiServiceRead`，调用 `DescribeSecurityAPIService` API（Limit=100），根据 api_service_ids 过滤结果，处理资源不存在情况
- [x] 1.4 实现 Update 函数 `resourceTencentCloudTeoSecurityApiServiceUpdate`，检查 immutable 参数（zone_id），调用 `ModifySecurityAPIResource` API 更新 api_resources，处理重试逻辑
- [x] 1.5 实现 Delete 函数 `resourceTencentCloudTeoSecurityApiServiceDelete`，调用 `DeleteSecurityAPIService` API，从 d.Get() 获取参数，处理重试逻辑

## 2. Provider 注册与文档索引

- [x] 2.1 在 `tencentcloud/provider.go` 的 ResourcesMap 中添加 `"tencentcloud_teo_security_api_service": teo.ResourceTencentCloudTeoSecurityApiService()` 注册
- [x] 2.2 在 `tencentcloud/provider.md` 的 TEO Resource 部分添加 `tencentcloud_teo_security_api_service` 条目

## 3. 资源文档

- [x] 3.1 创建 `tencentcloud/services/teo/resource_tc_teo_security_api_service.md` 文件，包含一句话描述（提及 TEO 产品）、Example Usage（含 zone_id、api_services 配置）和 Import 部分

## 4. 单元测试

- [x] 4.1 创建测试文件 `tencentcloud/services/teo/resource_tc_teo_security_api_service_test.go`，使用 gomonkey mock 方式编写 Create、Read、Update、Delete 操作的单元测试
- [x] 4.2 使用 `go test -gcflags=all=-l` 运行单元测试，确保所有测试通过
