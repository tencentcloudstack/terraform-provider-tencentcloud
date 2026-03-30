## Context

当前 `tencentcloud_teo_function` 资源的 `function_id` 参数被定义为 Computed 字段，仅在资源创建后由系统自动生成并填充。这种设计限制了用户在特定场景下的灵活性，例如：
- 需要与现有系统保持一致的 function_id 命名规范
- 需要预设 function_id 以满足业务流程要求
- 需要在迁移场景中保持原有 function_id

项目使用 Terraform Plugin SDK v2 和 tencentcloud-sdk-go，资源文件位于 `tencentcloud/services/teo/resource_tc_teo_function.go`。资源通过 CreateFunction API 创建，通过 ModifyFunction API 更新，通过 DeleteFunction API 删除。

## API Limitation Discovery

**重要发现：** 在实施过程中发现，腾讯云 TEO API 的 `CreateFunction` 接口**不支持**用户自定义 `function_id` 参数。经过检查 API 定义：

- `CreateFunctionRequest` 结构体**没有** `FunctionId` 字段
- `function_id` 只能由系统自动生成
- `ModifyFunctionRequest` 需要 `FunctionId` 作为标识符，但**不能修改** `FunctionId` 本身

**结论：** 由于 API 层面的限制，当前无法实现在创建时指定自定义 `function_id` 的需求。

## Goals / Non-Goals

**Goals:**
- ~~允许用户在创建 `tencentcloud_teo_function` 资源时指定 `function_id` 参数~~（因 API 限制无法实现）
- ~~保持向后兼容，未提供 `function_id` 时仍由系统自动生成~~（保持现状）
- ~~确保 Schema 变更符合 Terraform 最佳实践（Optional + Computed）~~
- 记录 API 限制，为未来可能的 API 改进提供参考

**Non-Goals:**
- 不修改 function_id 的唯一性约束（仍由 API 侧保证）
- 不修改资源的其他行为（remark、content、name 等参数）
- 不影响资源的更新和删除逻辑
- **新增：** 不通过 Terraform Provider 实现绕过 API 限制的变通方案

## Decisions

### 1. Schema 修改方式

**决策：** 将 `function_id` 参数从 `Computed` 改为 `Optional + Computed`

**理由：**
- `Optional` 允许用户在配置中显式指定 function_id
- `Computed` 保留系统自动生成的能力，并在读取时返回 API 返回的值
- 这是 Terraform 中处理"可配置但也可自动生成"字段的常用模式
- 完全向后兼容，现有配置无需修改

### 2. API 调用逻辑

**决策：** 在 `resourceTencentCloudTeoFunctionCreate` 函数中，如果用户提供了 `function_id`，则将其传递给 `CreateFunction` API；否则不传递

**理由：**
- 遵循最小修改原则，仅添加必要的逻辑
- 让后端 API 决定未提供 FunctionId 时的行为（自动生成）
- 避免在 Terraform Provider 端处理 ID 生成逻辑

### 3. 错误处理

**决策：** 如果用户提供的 function_id 在 API 调用时产生冲突（例如 ID 已存在），直接返回 API 错误信息

**理由：**
- API 会处理 ID 冲突并返回明确的错误信息
- Terraform Provider 不需要额外的重复验证逻辑
- 保持错误信息的一致性和准确性

### 4. 状态刷新逻辑

**决策：** 状态刷新逻辑保持不变，仍使用 zone_id 和 function_id 的组合作为资源 ID

**理由：**
- 当前实现已经正确处理了 function_id 为空和有值的情况
- 无需修改资源 ID 的组合方式
- 状态刷新的等待逻辑保持稳定

## Risks / Trade-offs

### 风险 1：function_id 冲突

**风险：** 用户指定的 function_id 可能与已存在的函数冲突，导致创建失败

**缓解措施：**
- API 会返回明确的错误信息，用户可以根据错误调整配置
- 在文档中说明 function_id 的唯一性约束
- 在错误信息中提供清晰的指引

### 风险 2：测试覆盖不足

**风险：** 新增功能可能缺少充分的测试用例，导致边界场景未被覆盖

**缓解措施：**
- 新增测试用例覆盖：
  - 提供有效 function_id 的创建场景
  - 未提供 function_id 的创建场景（验证向后兼容）
  - function_id 冲突的错误场景
- 使用 TF_ACC=1 运行完整的验收测试

### 权衡：复杂性 vs 灵活性

**权衡：** 增加 Schema 的复杂性（从 Computed 到 Optional+Computed）以换取用户的灵活性

**分析：**
- 复杂性增加很小，仅涉及 Schema 定义和条件逻辑
- 灵活性提升显著，满足用户在特定场景下的需求
- 向后兼容性得到完全保留
- **结论：** 权衡合理，收益大于成本
