## Context

当前 tencentcloud_teo_function 资源使用 CreateFunction API 创建函数，其中 function_id 参数仅在响应中返回，不支持在创建时指定。这限制了用户在某些场景下对 FunctionId 的控制，例如需要复用已存在的 FunctionId 或进行资源关联。

## Goals / Non-Goals

**Goals:**
- 允许用户在创建 tencentcloud_teo_function 资源时指定 function_id 参数
- 保持向后兼容性，不破坏现有的 Terraform 配置和状态
- 确保在未指定 function_id 时仍由服务端自动生成

**Non-Goals:**
- 不修改 Update 或 Delete 操作中的 function_id 行为
- 不支持在已存在的资源上修改 function_id
- 不改变资源 ID 的组成格式

## Decisions

### 1. Schema 修改策略
将 function_id 参数从 Computed 改为 Computed+Optional。这允许用户在创建时指定值，同时仍由服务端填充该字段。

**理由：** 这种方式符合 Terraform 最佳实践，对于 ID 类字段，通常采用 Computed+Optional 模式。使用户既可以选择指定，也可以让系统自动生成。

### 2. Create 函数修改
在 `resourceTencentCloudTeoFunctionCreate` 函数中，检查用户是否指定了 function_id。如果指定，则将其包含在 CreateFunctionRequest 中；否则，不传递该参数。

**理由：** 保持与 API 的兼容性，如果 API 不支持 FunctionId 参数，可以优雅地处理。

### 3. 向后兼容性保证
由于 function_id 之前仅为 Computed，现有配置中不会包含该参数的显式设置。因此将其改为 Optional 不会破坏现有配置。

**理由：** 在 Terraform schema 中，将 Computed 改为 Computed+Optional 是向后兼容的变更，不会强制用户修改其配置。

### 4. API 兼容性处理
需要先验证 CreateFunction API 是否支持 FunctionId 参数。如果 API 不支持该参数，则需要在代码中进行条件判断，仅在 API 支持时才传递该参数。

**理由：** 避免向不支持的 API 传递参数导致错误。

## Risks / Trade-offs

### Risk 1: CreateFunction API 可能不支持 FunctionId 参数
**风险：** 如果 API 不支持 FunctionId 参数，直接传递可能导致错误。
**缓解措施：** 在代码中添加条件判断，检查用户是否指定了 function_id，并仅在 API 支持时才传递该参数。可以通过查阅 API 文档或测试验证 API 能力。

### Risk 2: 用户指定了无效的 FunctionId
**风险：** 用户可能指定了已存在或格式无效的 FunctionId，导致创建失败。
**缓解措施：** 依赖 API 的错误返回机制，将 API 错误信息返回给用户。在文档中说明 FunctionId 的约束条件。

### Risk 3: State 不一致
**风险：** 如果用户指定了 FunctionId 但 API 返回了不同的 FunctionId，可能导致状态不一致。
**缓解措施：** 在 Read 函数中从 API 响应中读取实际的 FunctionId 并更新状态，确保状态与 API 返回值一致。

### Trade-off 1: API 调用复杂性增加
**权衡：** 需要在 Create 函数中添加条件逻辑来判断是否传递 FunctionId。
**收益：** 提供了更灵活的资源创建方式。

### Trade-off 2: 测试覆盖增加
**权衡：** 需要新增测试用例覆盖用户指定和不指定 FunctionId 的场景。
**收益：** 确保新功能的正确性和向后兼容性。

## Migration Plan

### 部署步骤
1. 修改 `resource_tc_teo_function.go` 中的 schema，将 function_id 改为 Computed+Optional
2. 修改 `resourceTencentCloudTeoFunctionCreate` 函数，添加 FunctionId 参数的条件传递逻辑
3. 更新测试文件，添加新的测试用例
4. 更新文档文件，说明 function_id 参数的新用途
5. 运行完整的测试套件，确保所有测试通过

### 回滚策略
如果出现问题，可以回滚到之前的版本。由于 function_id 在配置中的存在性是可选的，回滚后现有配置仍能正常工作。

## Open Questions

1. CreateFunction API 是否支持 FunctionId 参数？需要查阅 API 文档或通过测试验证。
2. FunctionId 的格式和约束条件是什么？需要在 API 文档中确认。
3. 如果 API 不支持 FunctionId 参数，是否需要提供其他方式来实现类似功能？
