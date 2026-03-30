## Context

当前 tencentcloud_teo_function 资源通过 CreateFunction API 创建腾讯云边缘函数（TEO）服务中的函数。用户需要在 Terraform 配置中定义函数内容，资源会调用 API 创建新函数。

新的需求是允许用户在创建函数时指定一个已存在的 FunctionId，而不是总是创建新函数。这可以支持以下场景：
- 用户在腾讯云控制台预创建了函数，希望在 Terraform 中管理
- 用户需要复用跨项目/跨环境的函数实例
- 用户希望通过 FunctionId 引用已有的函数逻辑

当前资源实现位于 `tencentcloud/services/teo/resource_tencentcloud_teo_function.go`。

## Goals / Non-Goals

**Goals:**
- 为 tencentcloud_teo_function 资源添加可选的 FunctionId 参数
- 在 create 操作中将 FunctionId 参数传递给 CreateFunction API
- 保持完全向后兼容，现有配置无需修改
- 确保资源状态正确读取和存储 FunctionId

**Non-Goals:**
- 不修改 update、delete 或 read 操作的行为（除非因 FunctionId 而必要）
- 不改变现有的资源 schema 结构（仅新增字段）
- 不引入新的依赖或外部服务
- 不修改其他 TEO 相关资源或数据源

## Decisions

### Schema Design: Optional Computed Field
**Decision:** 将 FunctionId 定义为可选的 Schema 属性，设置为 `Optional: true`。

**Rationale:**
- 可选字段确保向后兼容，现有配置无需修改
- 用户可以选择是否提供 FunctionId
- 符合 Terraform 最佳实践（可选的引用参数）

**Alternative Considered:** 将 FunctionId 设为必填。
- **Rejected:** 会破坏向后兼容性，强制所有用户提供 FunctionId。

### Parameter Mapping to API
**Decision:** 在 create 函数中检查 FunctionId 是否存在，如果存在则传递给 CreateFunction API 请求。

**Rationale:**
- 最小化代码变更，仅影响创建逻辑
- 保持现有流程不变，仅在必要时添加参数
- 利用现有的 SDK 调用机制

**Implementation Approach:**
```go
func resourceTencentcloudTeoFunctionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    // 现有代码...
    request := teo.NewCreateFunctionRequest()
    // 现有参数设置...

    if v, ok := d.GetOk("function_id"); ok {
        request.FunctionId = v.(string)
    }

    // 其余代码...
}
```

### State Management: Read Back from API
**Decision:** 在 read 函数中从 API 响应读取 FunctionId 并更新资源状态。

**Rationale:**
- 确保状态一致性
- 允许后续的 refresh 操作正确读取 FunctionId
- 符合 Terraform Provider 状态管理最佳实践

**Implementation Approach:**
```go
func resourceTencentcloudTeoFunctionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    // 现有代码获取响应...

    if response.FunctionId != "" {
        _ = d.Set("function_id", response.FunctionId)
    }

    // 其余代码...
}
```

### Testing Strategy
**Decision:** 添加新的测试场景来验证 FunctionId 参数，同时保留现有测试用例。

**Rationale:**
- 确保新功能按预期工作
- 验证向后兼容性
- 遵循项目测试惯例

**Test Scenarios:**
1. 测试使用 FunctionId 创建资源
2. 测试不使用 FunctionId 创建资源（确保向后兼容）
3. 测试读取 FunctionId 从状态中
4. 测试更新资源时 FunctionId 的行为

### Documentation Updates
**Decision:** 更新 `website/docs/r/teo_function.md` 文档，添加 FunctionId 参数说明和示例。

**Rationale:**
- 用户需要了解新参数的用途和使用方法
- 文档是 Terraform Provider 的重要组成部分
- 提供清晰的使用示例

## Risks / Trade-offs

**Risk 1:** API 可能不支持在创建时传入 FunctionId
- **Mitigation:** 在实现前先验证 CreateFunction API 的 SDK 是否支持此参数，查阅腾讯云 TEO API 文档

**Risk 2:** 传入无效的 FunctionId 可能导致 API 调用失败
- **Mitigation:** 不添加客户端验证（避免过度检查），让 API 返回错误信息给用户，符合 Terraform Provider 惯例

**Risk 3:** 现有用户的状态可能与新 schema 不兼容（如果状态中已包含 FunctionId）
- **Mitigation:** 由于是新参数，现有状态中不应包含此字段，无需迁移。新资源创建后会自动写入状态

**Trade-off:** 不添加 FunctionId 的格式验证（如 UUID 格式检查）
- **Rationale:** 避免过度限制用户，让后端 API 验证参数格式
- **Impact:** 错误信息可能稍晚返回（API 调用时而非配置解析时）

## Migration Plan

此变更无需数据迁移，因为：
- FunctionId 是新增的可选字段
- 现有资源和状态不受影响
- 新资源创建后会自动包含 FunctionId（如果提供）

**Deployment Steps:**
1. 修改 `tencentcloud/services/teo/resource_tencentcloud_teo_function.go` 添加 FunctionId 到 schema
2. 更新 create 函数支持传递 FunctionId
3. 更新 read 函数支持读取 FunctionId
4. 添加/更新测试用例
5. 更新文档
6. 运行 `TF_ACC=1 go test` 验证所有测试通过
7. 提交代码变更

**Rollback Strategy:**
- 如果发现问题，可以回滚代码变更
- 由于是可选字段，已创建的资源不受影响
- 回滚后用户仍可使用原有方式创建资源（不提供 FunctionId）

## Open Questions

无。这是一个 straightforward 的参数添加，没有需要进一步澄清的技术决策。
