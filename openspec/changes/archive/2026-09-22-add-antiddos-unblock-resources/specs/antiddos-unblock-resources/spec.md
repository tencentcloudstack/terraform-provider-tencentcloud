## ADDED Requirements

### Requirement: 解封被封堵的 AntiDDoS 资源

framework action 资源 `tencentcloud_antiddos_unblock_resources` SHALL 通过 `UnblockResources` API 申请解封被封堵的 AntiDDoS 资源（公网 IP 列表）。该 action 为一次性操作，操作完成后不持久化任何云侧状态。

#### Scenario: 成功解封资源

- **WHEN** 用户在 Terraform 配置中声明该 action 并提供非空的 `resources`（公网 IP 列表）参数
- **THEN** action 在 `Invoke` 中将 `resources` 转换为 `[]*string` 并调用 `UseAntiddosV20250903Client().UnblockResourcesWithContext(ctx, request)`
- **AND** 调用包裹在 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 重试机制中
- **AND** 调用前执行 `ratelimit.Check(request.GetAction())` 限流检查
- **AND** API 错误通过 `tccommon.RetryError(e)` 包装后交由 retry 机制处理

### Requirement: 必填入参校验

action 的 `Invoke` 方法 SHALL 在调用云 API 之前校验 `resources` 参数非空。

#### Scenario: 缺少 resources 参数时拒绝调用

- **WHEN** `resources` 为 null、unknown 或空列表
- **THEN** action 通过 `resp.Diagnostics.AddError` 返回 "Missing resources" 错误
- **AND** 不发起任何云 API 调用

### Requirement: Provider 未配置时的保护

action 的 `Invoke` 方法 SHALL 在调用服务前校验 client 是否可用。

#### Scenario: Provider client 为 nil

- **WHEN** action 的 `Client()` 返回 nil
- **THEN** action 通过 `resp.Diagnostics.AddError` 返回 "Provider not configured" 错误
- **AND** 不发起任何云 API 调用

### Requirement: action schema 定义

action 的 schema SHALL 定义一个必填的 `resources` 属性（`List` of `String`），对应 `UnblockResources` API 的 `request.Resources`。

#### Scenario: schema 校验

- **WHEN** 查询 action 的 schema
- **THEN** schema 包含 `resources` 属性，其类型为 List、元素类型为 String
- **AND** `resources` 属性标记为 Required
- **AND** action 的 TypeName 为 `tencentcloud_antiddos_unblock_resources`

### Requirement: framework action 注册

该 action SHALL 通过 `tencentcloud/framework/registry.go` 的 `actionFactories` 注册，工厂函数为 `antiddos.NewAntiddosUnblockResources`。

#### Scenario: action 可被 provider 识别

- **WHEN** framework provider 收集 action 工厂列表
- **THEN** `antiddos.NewAntiddosUnblockResources` 出现在 `actionFactories` 切片中
- **AND** 该 action 可被 Terraform 识别为 `tencentcloud_antiddos_unblock_resources`