## 1. 服务层错误处理修改

- [x] 1.1 修改 `tencentcloud/services/vpc/service_tencentcloud_vpc.go` 中 `DescribeVpnGatewayRoutes` 方法 retry 块：在 `errRet != nil` 分支内，于 `tccommon.RetryError(errRet, tccommon.InternalError)` 之前，增加 `if tccommon.IsExpectError(errRet, []string{"ResourceNotFound"}) { return nil }` 判断，命中则跳出 retry 且不返回错误，使 `response` 保持 `nil` 并由方法后续 `response == nil` 分支返回空结果
- [x] 1.2 确认 `tccommon` 包已被该文件导入（现状已导入），无需新增 import；确认 `IsExpectError` 签名为 `func IsExpectError(err error, expectError []string) bool`

## 2. 验证一致性

- [x] 2.1 核对 vendor 云 API `DescribeVpnGatewayRoutes`（`vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312/client.go`）已声明 `RESOURCENOTFOUND = "ResourceNotFound"` 错误码，确认本次变更与云 API 行为一致
- [x] 2.2 核对资源 Read（`tencentcloud/services/vpn/resource_tc_vpn_gateway_route.go`）在 `route == nil` 时执行 `d.SetId("")` 并返回 `nil`，确认 `DescribeVpnGatewayRoutes` 返回空结果后能正确将资源从 state 清除

## 3. 单元测试

- [x] 3.1 在 `tencentcloud/services/vpn/resource_tc_vpn_gateway_route_test.go` 中补充单元测试用例：使用 mock（gomonkey）对 `VpcClient.DescribeVpnGatewayRoutes` 进行 mock，使其返回 `ResourceNotFound` 的 `TencentCloudSDKError`，断言 `DescribeVpnGatewayRoutes` 服务方法返回 `err == nil` 且结果为空，且资源 Read 能正常 `d.SetId("")` 返回 `nil`

## 4. 验证任务

- [x] 4.1 代码正确性检查：确认 retry 块内仅调用云 API 接口，未在 retry 块内执行设置 id 等成功操作；确认 `return nil` 后未触发 `d.SetId("")` 遗漏 id 日志的问题（资源 Read 在 `route == nil` 分支前已打印足够上下文）
- [x] 4.2 确认本次变更未修改资源/数据源 schema，未引入新增/修改的 `.md` 文档变更（无 schema 变更，无需 `make doc` 重新生成）