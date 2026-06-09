## 1. Schema & CRUD 函数修改

- [x] 1.1 在 `tencentcloud/services/teo/resource_tc_teo_l4_proxy.go` 的 Schema 中添加 `proxy_id` computed 属性（TypeString, Computed: true, Description: "L4 proxy instance ID."）
- [x] 1.2 在 `resourceTencentCloudTeoL4ProxyCreate` 函数中，CreateL4Proxy API 调用成功后，通过 `d.Set("proxy_id", proxyId)` 设置 `proxy_id`
- [x] 1.3 在 `resourceTencentCloudTeoL4ProxyRead` 函数中，当 `respData.ProxyId` 不为 nil 时，设置 `d.Set("proxy_id", respData.ProxyId)`
- [x] 1.4 将 `ddos_protection_config` 及其子字段标记为 deprecated（添加 `Deprecated` 和 `Computed: true` 属性）

## 2. 测试

- [x] 2.1 在 `tencentcloud/services/teo/resource_tc_teo_l4_proxy_test.go` 中添加 gomonkey mock 单元测试（TestTeoL4Proxy_Create、TestTeoL4Proxy_Read、TestTeoL4Proxy_Read_NotFound、TestTeoL4Proxy_Schema），验证 `proxy_id` 在 Create 和 Read 操作中被正确设置

## 3. 文档

- [x] 3.1 更新 `tencentcloud/services/teo/resource_tc_teo_l4_proxy.md` 示例文件，添加 `proxy_id` 属性说明和 output 示例
