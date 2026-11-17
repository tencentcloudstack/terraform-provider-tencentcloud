## 1. Schema 变更：移除 ForceNew

- [x] 1.1 在 `tencentcloud/services/bh/resource_tc_bh_resource.go` 的 `web_access` schema 字段中移除 `ForceNew: true`（约第 124–130 行），其余属性（`Type: schema.TypeInt`、`Optional: true`、`Computed: true`、`Description`）保持不变
- [x] 1.2 在 `tencentcloud/services/bh/resource_tc_bh_resource.go` 的 `client_access` schema 字段中移除 `ForceNew: true`（约第 132–138 行），其余属性保持不变

## 2. Update 逻辑：新增 web_access 变更分支

- [x] 2.1 在 `resourceTencentCloudBhResourceUpdate` 中、现有 `external_access` 处理块之后新增 `if d.HasChange("web_access")` 分支
- [x] 2.2 当新值为 `1` 时，构造 `EnableWebAccessRequest`（仅设置 `ResourceId = helper.String(resourceId)`），用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包裹 `EnableWebAccessWithContext` 调用，失败用 `tccommon.RetryError(e)` 包装
- [x] 2.3 Enable 调用成功后，构造 `DescribeResourcesRequest`（`ResourceIds = []*string{&resourceId}`），用 `resource.Retry(tccommon.ReadRetryTimeout*7, ...)` 轮询，等待 `result.Response.ResourceSet[0].WebAccess == 1` 返回 nil，否则 `RetryableError`
- [x] 2.4 当新值为 `0` 时，调用 `DisableWebAccessWithContext`（入参同上），成功后轮询等待 `WebAccess == 0`
- [x] 2.5 wait 块中保留空指针保护：`result == nil || result.Response == nil || result.Response.ResourceSet == nil || len(...) == 0` → `NonRetryableError`；`WebAccess == nil` → `NonRetryableError`

## 3. Update 逻辑：新增 client_access 变更分支

- [x] 3.1 在 `web_access` 分支之后新增 `if d.HasChange("client_access")` 分支
- [x] 3.2 当新值为 `1` 时，构造 `EnableClientTcpAccessRequest`（仅设置 `ResourceId`），用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包裹 `EnableClientTcpAccessWithContext` 调用
- [x] 3.3 Enable 调用成功后，用 `DescribeResources` 轮询等待 `result.Response.ResourceSet[0].ClientAccess == 1`
- [x] 3.4 当新值为 `0` 时，调用 `DisableClientTcpAccessWithContext`，成功后轮询等待 `ClientAccess == 0`
- [x] 3.5 wait 块保留与 web_access 一致的空指针保护逻辑

## 4. 确认不误调 ModifyResource

- [x] 4.1 确认 `mutableArgs` 数组仍为 `{"resource_edition", "resource_node", "auto_renew_flag"}`，未将 `web_access`/`client_access` 加入（这两个字段不在 `ModifyResource` 接口入参中）

## 5. 测试：补充单元测试用例

- [x] 5.1 在 `tencentcloud/services/bh/resource_tc_bh_resource_test.go` 中使用 gomonkey mock `EnableWebAccessWithContext` / `DisableWebAccessWithContext` / `DescribeResourcesWithContext`，补充 `web_access` 由 0→1 与 1→0 的 update 单测用例
- [x] 5.2 使用 gomonkey mock `EnableClientTcpAccessWithContext` / `DisableClientTcpAccessWithContext` / `DescribeResourcesWithContext`，补充 `client_access` 由 0→1 与 1→0 的 update 单测用例
- [x] 5.3 确保测试中对 mock 函数返回的 error 均被检查，无未使用变量；测试代码可在当前环境编译通过

## 6. 文档更新

- [x] 6.1 更新 `tencentcloud/services/bh/resource_tc_bh_resource.md`，在示例中体现 `web_access`、`client_access` 支持在 update 阶段修改（移除重建的误导），保持与 `.go` 字段含义一致

## 7. 验证（收尾阶段执行，非本提案阶段）

- [x] 7.1 由收尾阶段通过 `make doc` 生成 `website/docs/` 文档，禁止手动编写 website 文件
- [x] 7.2 由收尾阶段执行 `gofmt` 格式化变更的 Go 文件
- [x] 7.3 由收尾阶段创建 `.changelog/` 下的 changelog 文件