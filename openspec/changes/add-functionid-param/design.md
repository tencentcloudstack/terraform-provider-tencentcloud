## Context

当前 tencentcloud_teo_function 资源的 function_id 字段定义为 Computed，在创建资源时由腾讯云 CreateFunction API 自动生成 FunctionId 并返回。用户反馈需要在创建时指定自定义的 FunctionId，以满足特定的业务需求和管理规范（例如，与内部系统集成、使用固定标识符等）。

Terraform Provider for TencentCloud 基于 Go 1.17+ 和 Terraform Plugin SDK v2 开发，资源文件位于 `tencentcloud/services/teo/` 目录下。资源使用腾讯云 SDK (teov20220901) 调用云 API。

## Goals / Non-Goals

**Goals:**
- 将 tencentcloud_teo_function 资源的 function_id 字段从 Computed 改为 Optional + Computed
- 修改 resourceTencentCloudTeoFunctionCreate 函数，支持传入用户指定的 FunctionId 到 CreateFunction API
- 确保向后兼容性：未指定 FunctionId 时保持原有行为（由 API 自动生成）
- 更新相关测试用例和文档

**Non-Goals:**
- 不修改 tencentcloud_teo_function 资源的其他字段
- 不修改资源的 ID 格式（zoneId#functionId）
- 不修改其他 TEO 相关资源或数据源
- 不改变 Read、Update、Delete 函数的逻辑

## Decisions

### Schema 修改策略

将 function_id 字段定义为 Optional + Computed：
- **Optional**: 允许用户在配置中指定 FunctionId
- **Computed**: 如果用户未指定，则接收 API 自动生成的 FunctionId

这是 Terraform Provider 中处理"用户可输入但也可由系统生成"字段的常见模式。

### API 调用逻辑

在 resourceTencentCloudTeoFunctionCreate 函数中：
1. 检查用户是否提供了 function_id 参数
2. 如果提供了，则设置 `request.FunctionId` 参数
3. 如果未提供，则不设置该参数（让 API 自动生成）

这种设计确保了向后兼容性：现有配置不包含 function_id 时，行为与之前完全一致。

### 错误处理

如果用户传入的 FunctionId 与现有资源冲突，CreateFunction API 会返回错误。这个错误会通过 helper.Retry 机制传递给 Terraform，用户会看到明确的错误信息。

## Risks / Trade-offs

### 风险 1: CreateFunction API 可能不支持传入 FunctionId 参数

**影响**: 如果 API 不支持此参数，需要修改 API 文档或与云服务团队沟通。

**缓解措施**: 在实现前先查阅腾讯云 TEO CreateFunction API 文档，确认 FunctionId 参数是否为可选输入参数。如果不支持，则需要通过其他方式实现（如创建后使用 ModifyFunction API 修改，但这需要 API 支持）。

### 风险 2: 用户传入的 FunctionId 格式可能不符合要求

**影响**: API 验证失败，创建失败。

**缓解措施**: 依赖 API 的验证逻辑，Terraform 会将 API 错误返回给用户。在文档中明确说明 FunctionId 的格式要求（参考 API 文档）。

### 风险 3: 向后兼容性问题

**影响**: 现有用户的状态管理可能受影响。

**缓解措施**:
- 使用 Optional + Computed 模式确保未指定时的行为与之前一致
- 确保 Read 函数正确读取并设置 function_id
- 不修改资源 ID 格式（zoneId#functionId）
- 在变更前运行现有的验收测试，确保无破坏性变更

### 权衡: 灵活性 vs 简单性

**决策**: 选择灵活性，允许用户指定 FunctionId，因为这满足了明确的业务需求，且通过 Optional + Computed 模式保持了简单性。

## Migration Plan

1. **开发阶段**:
   - 修改 resource_tc_teo_function.go 的 schema
   - 修改 resourceTencentCloudTeoFunctionCreate 函数逻辑
   - 添加新的测试用例（指定 FunctionId 的情况）
   - 更新文档示例

2. **测试阶段**:
   - 运行现有的验收测试（TF_ACC=1），确保向后兼容
   - 运行新的测试用例，验证 FunctionId 输入功能
   - 手动测试：指定和不指定 FunctionId 两种场景

3. **部署阶段**:
   - 合并代码到主分支
   - 发布新版本的 Terraform Provider
   - 更新版本说明和文档

4. **回滚策略**:
   - 如果发现严重问题，回退到上一个版本
   - 由于保持了向后兼容性，现有用户不会受到影响

## Open Questions

1. CreateFunction API 是否支持 FunctionId 作为可选输入参数？
   - 需要查阅 TEO API 文档确认
   - 如果不支持，需要寻找替代方案

2. FunctionId 的格式限制是什么？
   - 需要查阅 API 文档
   - 如果有格式要求，需要在文档中明确说明

3. 是否需要为 FunctionId 添加验证逻辑？
   - 通常依赖 API 验证即可
   - 如果格式复杂，可以在 Provider 层添加验证（但会增加维护成本）
