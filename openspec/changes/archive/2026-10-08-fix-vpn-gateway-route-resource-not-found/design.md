## Context

资源 `tencentcloud_vpn_gateway_route`（定义于 `tencentcloud/services/vpn/resource_tc_vpn_gateway_route.go`）的 Read/Update/Delete 逻辑均依赖服务方法 `VpcService.DescribeVpnGatewayRoutes`（定义于 `tencentcloud/services/vpc/service_tencentcloud_vpc.go`）。该方法在 `resource.Retry` 块中调用云 API `DescribeVpnGatewayRoutes`，当前实现未区分错误码，任何错误都会触发 `tccommon.RetryError(errRet, tccommon.InternalError)` 走重试逻辑。

当 VPN 网关路由已被远端删除时，云 API 返回 `ResourceNotFound` 错误码（vendor SDK 中 `DescribeVpnGatewayRoutes` 已声明该错误码）。当前实现会把该错误持续重试直至超时，导致资源 Read 报错而非将资源从 state 中清空，阻塞 `terraform destroy`/`terraform refresh` 等操作。

## Goals / Non-Goals

**Goals:**
- 在 `DescribeVpnGatewayRoutes` 服务方法的 retry 块内识别 `ResourceNotFound` 错误码，命中时 `return nil`（不报错），使上层资源 Read 进入"路由列表为空"的分支，进而 `d.SetId("")` 正常清空 id。
- 复用现有 `tccommon.IsExpectError` 工具函数，与代码库中已有处理模式（如 `tencentcloud/services/tcm/resource_tc_tcm_mesh.go`）保持一致。
- 保持资源 `tencentcloud_vpn_gateway_route` 与数据源 `tencentcloud_vpn_gateway_routes` 的正常读行为不受影响（仅远端删除场景行为改变）。

**Non-Goals:**
- 不修改资源/数据源的 schema。
- 不修改 Create/Update/Delete 中除 `DescribeVpnGatewayRoutes` 调用结果处理以外的逻辑。
- 不调整分页逻辑或 `VPN_DESCRIBE_LIMIT` 常量。

## Decisions

### 决策 1：在 retry 块内判断 `ResourceNotFound` 并 `return nil`

**方案**：将 retry 块内的错误处理由：

```go
if errRet != nil {
    return tccommon.RetryError(errRet, tccommon.InternalError)
}
```

修改为：

```go
if errRet != nil {
    if tccommon.IsExpectError(errRet, []string{"ResourceNotFound"}) {
        return nil
    }
    return tccommon.RetryError(errRet, tccommon.InternalError)
}
```

**理由**：`ResourceNotFound` 表示路由已不存在，属于预期内的"资源已删除"语义，应跳出 retry 且不返回错误，使 `response` 保持为 `nil`，进而命中方法后续的 `response == nil || response.Response == nil` 分支并返回空结果（`errRet == nil, result == nil`）。上层资源 Read 在 `route == nil` 时执行 `d.SetId("")`，符合 Terraform 对"远端资源已删除"的标准处理方式。

**备选方案**：
- 在 retry 块外（`if errRet != nil { return errRet, nil }` 之前）判断：不可行，因为 retry 块内 `return nil` 后 `errRet` 已被清空，块外无法再捕获到该错误码。
- 在资源 Read 层（`resource_tc_vpn_gateway_route.go`）判断错误码：不符合分层设计，且数据源 `tencentcloud_vpn_gateway_routes` 也会重复处理；统一在服务方法层处理更合理。

### 决策 2：使用 `tccommon.IsExpectError` 而非字符串匹配

**理由**：`tccommon.IsExpectError` 已封装了对 `TencentCloudSDKError`（含国际版英文错误）的匹配，并支持长/短错误码匹配，是代码库中处理"期望错误码"的标准做法，避免手写类型断言与字符串比较。

## Risks / Trade-offs

- [风险] `ResourceNotFound` 在某些非"资源已删除"场景下也可能返回（如传入错误的 vpnGatewayId），命中 `return nil` 后会静默返回空结果，可能掩盖配置错误。→ **缓解**：该行为符合 Terraform 对"资源不存在即从 state 移除"的约定；且数据源场景下返回空列表为合理表现，用户可通过 `terraform plan` 发现配置问题。
- [风险] `return nil` 跳过 retry 后 `response` 仍为 `nil`，方法后续 `response == nil` 分支返回 `(nil, nil)`，需确认上层资源 Read 能正确处理空结果。→ **已验证**：资源 Read 在 `routeList` 为空时 `route == nil`，执行 `d.SetId("")` 并返回 `nil`，行为正确。