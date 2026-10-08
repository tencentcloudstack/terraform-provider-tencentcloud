## ADDED Requirements

### Requirement: DescribeVpnGatewayRoutes 应容忍 ResourceNotFound 错误

`VpcService.DescribeVpnGatewayRoutes` 在调用云 API `DescribeVpnGatewayRoutes` 失败时，SHALL 识别 `ResourceNotFound` 错误码；当命中该错误码时 MUST 返回 `nil`（不报错且不重试），其余错误 MUST 维持原有重试行为（`tccommon.RetryError(errRet, tccommon.InternalError)`）。

#### Scenario: 资源已被远端删除时 Read 不报错
- **WHEN** `DescribeVpnGatewayRoutes` 云 API 返回 `ResourceNotFound` 错误码
- **THEN** `VpcService.DescribeVpnGatewayRoutes` 返回 `errRet == nil` 且 `result` 为空，使上层资源 Read 进入"路由不存在"分支并执行 `d.SetId("")`

#### Scenario: 其他错误仍走重试
- **WHEN** `DescribeVpnGatewayRoutes` 云 API 返回非 `ResourceNotFound` 的错误（如 `InternalError`）
- **THEN** `VpcService.DescribeVpnGatewayRoutes` 通过 `tccommon.RetryError(errRet, tccommon.InternalError)` 触发重试，行为与变更前一致

#### Scenario: 资源正常存在时读行为不变
- **WHEN** `DescribeVpnGatewayRoutes` 云 API 正常返回路由列表
- **THEN** `VpcService.DescribeVpnGatewayRoutes` 返回对应路由列表，行为与变更前一致