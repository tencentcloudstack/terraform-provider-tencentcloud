## 1. Schema & CRUD 代码修改

- [x] 1.1 在 `tencentcloud/services/teo/resource_tc_teo_l4_proxy.go` 的 Schema 中添加 `proxy_id` 计算属性（Type: schema.TypeString, Computed: true）
- [x] 1.2 在 `resourceTencentCloudTeoL4ProxyCreate` 函数中，CreateL4Proxy API 调用成功后，将 `response.Response.ProxyId` 设置到 ResourceData 的 `proxy_id` 字段
- [x] 1.3 在 `resourceTencentCloudTeoL4ProxyRead` 函数中，DescribeL4Proxy 查询成功后，将 `respData.ProxyId` 设置到 ResourceData 的 `proxy_id` 字段

## 2. 单元测试

- [x] 2.1 在 `tencentcloud/services/teo/resource_tc_teo_l4_proxy_test.go` 中补充 `proxy_id` 字段的单元测试用例，使用 gomonkey mock 云API，验证 Create 和 Read 流程中 proxy_id 正确设置

## 3. 文档更新

- [x] 3.1 更新 `tencentcloud/services/teo/resource_tc_teo_l4_proxy.md`，在 Sub-markdown 中添加 `proxy_id` 属性说明
