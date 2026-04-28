## 1. Service Layer

- [x] 1.1 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中添加 `DescribeTeoMultiPathGatewayById` 方法，调用 DescribeMultiPathGateway API，返回单个 MultiPathGateway 对象
- [x] 1.2 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中添加 `DeleteTeoMultiPathGateway` 方法，调用 DeleteMultiPathGateway API

## 2. Resource Schema & CRUD

- [x] 2.1 创建 `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway.go`，定义 ResourceTencentCloudTeoMultiPathGateway 函数，包含完整 schema 定义（zone_id、gateway_type、gateway_name、gateway_port、region_id、gateway_ip、gateway_id、status、need_confirm、lines）
- [x] 2.2 实现 resourceTencentCloudTeoMultiPathGatewayCreate 函数，调用 CreateMultiPathGateway API，设置复合 ID（zone_id#gateway_id），添加 WriteRetryTimeout 重试逻辑
- [x] 2.3 实现 resourceTencentCloudTeoMultiPathGatewayRead 函数，解析复合 ID，调用服务层 Describe 方法，设置所有 schema 字段，处理 nil 字段和资源不存在的情况
- [x] 2.4 实现 resourceTencentCloudTeoMultiPathGatewayUpdate 函数，使用 d.HasChange() 检查 gateway_name、gateway_ip、gateway_port 变更，调用 ModifyMultiPathGateway API，添加 WriteRetryTimeout 重试逻辑
- [x] 2.5 实现 resourceTencentCloudTeoMultiPathGatewayDelete 函数，解析复合 ID，调用服务层 Delete 方法，添加 WriteRetryTimeout 重试逻辑
- [x] 2.6 添加 Importer 支持（schema.ImportStatePassthrough）

## 3. Provider Registration

- [x] 3.1 在 `tencentcloud/provider.go` 的 ResourcesMap 中添加 `tencentcloud_teo_multi_path_gateway` 资源注册
- [x] 3.2 在 `tencentcloud/provider.md` 中添加资源文档入口

## 4. Documentation

- [x] 4.1 创建 `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway.md` 文档文件，包含一句话描述、Example Usage（cloud 和 private 两种类型示例）、Import 部分

## 5. Unit Tests

- [x] 5.1 创建 `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway_test.go`，使用 gomonkey mock 方式编写 Create、Read、Update、Delete 的单元测试
- [x] 5.2 运行 `go test -gcflags=all=-l` 验证单元测试通过

## 6. Verification

- [x] 6.1 检查所有云 API 接口参数映射的正确性，确保 Create 参数在 Create API 中存在、Modify 参数在 Modify API 中存在、Delete 参数在 Delete API 中存在
- [x] 6.2 检查代码可编译性，确保没有语法错误和未使用的 import
