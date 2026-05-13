## 1. 服务层实现

- [x] 1.1 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中新增 `DescribeTeoSecurityAPIResourceById` 方法，实现按 ZoneId 分页查询所有 API 资源（Limit=100），遍历匹配 apiResourceId，返回单个 `*APIResource`

## 2. 资源 Schema 与 CRUD 实现

- [x] 2.1 创建 `tencentcloud/services/teo/resource_tc_teo_security_api_resource.go`，定义 `ResourceTencentCloudTeoSecurityAPIResource()` 函数，包含完整 schema 定义：zone_id（required, ForceNew）、api_resources（TypeList 嵌套块，MaxItems:1，含 name/Required、path/Required、api_service_ids/Optional、methods/Optional、request_constraint/Optional、id/Computed），支持 Import
- [x] 2.2 实现 `resourceTencentCloudTeoSecurityAPIResourceCreate` 函数：构建 CreateSecurityAPIResourceRequest，填充 ZoneId 和单个 APIResource 参数（使用 buildSecurityAPIResourceFromMap 辅助函数，id 传空），调用 API（WriteRetryTimeout 重试），检查返回的 APIResourceIds，使用复合 ID zoneId#apiResourceId 设置 d.SetId()，调用 Read 同步状态
- [x] 2.3 实现 `resourceTencentCloudTeoSecurityAPIResourceRead` 函数：从 d.Id() 解析复合 ID 获取 zoneId 和 apiResourceId，调用服务层 DescribeTeoSecurityAPIResourceById，填充 api_resources 和 zone_id，资源不存在时 d.SetId("")
- [x] 2.4 实现 `resourceTencentCloudTeoSecurityAPIResourceUpdate` 函数：从 d.Id() 解析复合 ID，构建 ModifySecurityAPIResourceRequest，填充 ZoneId 和带 Id 的 APIResource（使用 buildSecurityAPIResourceFromMap 辅助函数），调用 API（WriteRetryTimeout 重试），调用 Read 同步状态
- [x] 2.5 实现 `resourceTencentCloudTeoSecurityAPIResourceDelete` 函数：从 d.Id() 解析复合 ID 获取 zoneId 和 apiResourceId，构建 DeleteSecurityAPIResourceRequest，填充 ZoneId 和 APIResourceIds，调用 API（WriteRetryTimeout 重试）
- [x] 2.6 实现 `buildSecurityAPIResourceFromMap` 辅助函数：将 schema map 转换为 *teo.APIResource，id 参数在创建时传空、更新时传实际值

## 3. Provider 注册

- [x] 3.1 在 `tencentcloud/provider.go` 中添加资源注册：`"tencentcloud_teo_security_api_resource": teo.ResourceTencentCloudTeoSecurityAPIResource()`
- [x] 3.2 在 `tencentcloud/provider.md` 的 TEO Resource 部分添加 `tencentcloud_teo_security_api_resource` 条目

## 4. 资源文档

- [x] 4.1 创建 `tencentcloud/services/teo/resource_tc_teo_security_api_resource.md`，包含一句话描述、Example Usage 和 Import 部分（Import 使用 zoneId#apiResourceId 格式）

## 5. 单元测试

- [x] 5.1 创建 `tencentcloud/services/teo/resource_tc_teo_security_api_resource_test.go`，使用 gomonkey mock 方式对 CRUD 操作编写单元测试，覆盖 Create/Read/Update/Delete/Schema 场景

## 6. 验证

- [x] 6.1 使用 `go test -gcflags=all=-l` 运行单元测试，确保所有测试通过
