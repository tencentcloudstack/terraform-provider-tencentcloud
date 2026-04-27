## 1. 修改 Delete 函数实现

- [x] 1.1 修改 `resourceTencentCloudTeoApplicationProxyDelete` 函数，将 `zone_id` 和 `proxy_id` 的获取方式从 `d.Id()` 解析改为 `d.Get("zone_id")` 和 `d.Get("proxy_id")`
- [x] 1.2 在 delete 函数中，将 `service.DeleteTeoApplicationProxyById()` 调用替换为直接构造 `DeleteApplicationProxyRequest`，设置 `ZoneId` 和 `ProxyId` 字段，并使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 调用 `DeleteApplicationProxy` API
- [x] 1.3 保留现有的先 offline 再 delete 的两步删除流程，确保 `ModifyApplicationProxyStatus` 调用仍然使用 `d.Get()` 获取的 `zone_id` 和 `proxy_id`

## 2. 更新单元测试

- [x] 2.1 在 `resource_tc_teo_application_proxy_test.go` 中添加或更新 delete 函数的单元测试，验证 `zone_id` 和 `proxy_id` 从 `d.Get()` 获取并正确传递给 `DeleteApplicationProxyRequest`
- [x] 2.2 使用 go test（带 `-gcflags=all=-l` 参数）运行单元测试，验证测试通过

## 3. 更新资源示例文档

- [x] 3.1 更新 `resource_tc_teo_application_proxy.md` 示例文件（如有必要）
