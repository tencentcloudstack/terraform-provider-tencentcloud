## 1. 服务层实现

- [x] 1.1 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中新增 `DescribeSecurityAPIResourceById` 方法，实现按 ZoneId 分页查询所有 API 资源（Limit=100），封装重试逻辑（tccommon.ReadRetryTimeout）

## 2. 资源 Schema 与 CRUD 实现

- [x] 2.1 创建 `tencentcloud/services/teo/resource_tc_teo_security_api_resource.go`，定义 `ResourceTencentCloudTeoSecurityApiResource()` 函数，包含完整 schema 定义：zone_id（required, ForceNew）、api_resources（TypeList 嵌套块，含 id/name/api_service_ids/path/methods/request_constraint）、api_resource_ids（computed），支持 Import 和 Timeouts
- [x] 2.2 实现 `resourceTencentCloudTeoSecurityApiResourceCreate` 函数：构建 CreateSecurityAPIResourceRequest，填充 ZoneId 和 APIResources 参数，调用 API（WriteRetryTimeout 重试），存储 api_resource_ids，设置 d.SetId(zoneId)，调用 Read 同步状态
- [x] 2.3 实现 `resourceTencentCloudTeoSecurityApiResourceRead` 函数：从 d.Id() 获取 zoneId，调用服务层 DescribeSecurityAPIResourceById，填充 api_resources 和 api_resource_ids，资源不存在时 d.SetId("")
- [x] 2.4 实现 `resourceTencentCloudTeoSecurityApiResourceUpdate` 函数：检测 api_resources 变更，构建 ModifySecurityAPIResourceRequest，填充 ZoneId 和带 Id 的 APIResources 列表，调用 API（WriteRetryTimeout 重试），调用 Read 同步状态
- [x] 2.5 实现 `resourceTencentCloudTeoSecurityApiResourceDelete` 函数：从 d.Get() 获取 zoneId 和 api_resource_ids，构建 DeleteSecurityAPIResourceRequest，调用 API（WriteRetryTimeout 重试）

## 3. Provider 注册

- [x] 3.1 在 `tencentcloud/provider.go` 中添加 teo 包的 import（如尚未引入）和资源注册：`"tencentcloud_teo_security_api_resource": teo.ResourceTencentCloudTeoSecurityApiResource()`
- [x] 3.2 在 `tencentcloud/provider.md` 的 TEO Resource 部分添加 `tencentcloud_teo_security_api_resource` 条目

## 4. 资源文档

- [x] 4.1 创建 `tencentcloud/services/teo/resource_tc_teo_security_api_resource.md`，包含一句话描述、Example Usage 和 Import 部分

## 5. 单元测试

- [x] 5.1 创建 `tencentcloud/services/teo/resource_tc_teo_security_api_resource_test.go`，使用 gomonkey mock 方式对 CRUD 操作编写单元测试，覆盖 Create/Read/Update/Delete 场景

## 6. 验证

- [x] 6.1 使用 `go test -gcflags=all=-l` 运行单元测试，确保所有测试通过
