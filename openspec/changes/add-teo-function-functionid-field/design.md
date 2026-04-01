## Context

当前 tencentcloud_teo_function 资源已支持基本的函数创建和管理功能，但缺少对 FunctionId 字段的支持。FunctionId 是 TEO 边缘函数的唯一标识符，由腾讯云服务端生成并在创建函数后返回。该字段在 API 层面存在于 CreateFunction 的响应参数中，但未在 Terraform Provider 的资源定义中暴露。

TencentCloud Terraform Provider 使用 Go 1.17+ 和 Terraform Plugin SDK v2 构建，遵循标准的资源 CRUD 模式。TEO 服务的代码位于 tencentcloud/services/teo/ 目录下，资源实现遵循 provider 项目的现有架构模式。

## Goals / Non-Goals

**Goals:**

- 在 tencentcloud_teo_function 资源的 Schema 中添加 FunctionId 字段（string 类型，Optional）
- 实现在 Create 函数中从 API 响应读取 FunctionId 并存储到 Terraform state 的逻辑
- 实现在 Read 函数中从 API 响应读取 FunctionId 并更新到 Terraform state 的逻辑
- 确保新增字段不会破坏现有配置和 state（向后兼容）
- 更新单元测试和验收测试以验证 FunctionId 字段的正确读写

**Non-Goals:**

- 修改 TEO Function 资源的其他字段或行为
- 修改 FunctionId 的生成逻辑（由腾讯云服务端负责）
- 实现 FunctionId 的可配置性（该字段为只读计算属性，由 API 返回）
- 修改 TEO Function 的 CRUD 流程或错误处理逻辑（仅添加字段读取）

## Decisions

### 1. FunctionId 字段作为 Computed Optional 字段

**决策**: 将 FunctionId 定义为 `Computed: true, Optional: false` 的字段。

**理由**: FunctionId 是由腾讯云服务端生成的唯一标识符，用户不能在 Terraform 配置中指定该值。根据 CreateFunction API 文档，FunctionId 是响应参数而非请求参数，因此应该定义为 Computed 字段。设置为 Optional: false 可以避免用户尝试在配置中手动指定该值。

**替代方案考虑**:
- 方案 A: `Computed: true, Optional: true` - 会允许用户配置该值，但实际会被忽略，可能导致混淆。
- 方案 B: `Computed: false, Optional: true` - 用户必须配置该值，但这不符合 API 设计，且会导致资源无法正确创建。

### 2. 字段映射方式

**决策**: 使用标准的 SDK 字段映射方式，直接将 API 响应中的 FunctionId 字段映射到 Terraform state 的 function_id 属性。

**理由**: 遵循 Provider 项目的现有模式，使用 `d.Set("function_id", response.FunctionId)` 进行字段设置。这种方式简单直接，与项目中其他资源的实现保持一致。

### 3. 在 Create 函数中的处理位置

**决策**: 在 Create 函数中调用 CreateFunction API 后，立即从响应中读取 FunctionId 并设置到 state 中。

**理由**: FunctionId 在创建成功后立即返回，需要在创建流程中尽早设置到 state 中，以确保后续操作（如 Update 或其他依赖操作）可以使用该 ID。

### 4. 在 Read 函数中的处理方式

**决策**: 在 Read 函数中调用 DescribeFunction 或相关查询 API 后，从响应中读取 FunctionId 并更新到 state 中。

**理由**: Read 函数的职责是同步资源状态到 Terraform state，需要包含所有可用的只读字段。这样即使 state 中的 FunctionId 丢失，也可以通过 Read 操作重新获取。

### 5. 测试覆盖策略

**决策**: 在单元测试中添加验证 FunctionId 字段设置的测试用例；在验收测试中验证完整的创建和读取流程。

**理由**:
- 单元测试可以快速验证字段映射逻辑的正确性
- 验收测试（TF_ACC=1）可以验证与真实 API 的集成
- 这两种测试的结合可以确保功能在各种场景下都能正常工作

## Risks / Trade-offs

### Risk 1: API 响应结构变更导致字段读取失败

**描述**: 如果腾讯云 API 响应中的 FunctionId 字段名称或位置发生变更，可能导致字段读取失败。

**缓解措施**:
- 遵循 TEO API 文档的当前定义
- 使用 SDK 提供的响应结构体，利用 Go 的类型系统进行编译时检查
- 在日志中记录详细的 API 响应，便于调试

### Risk 2: 现有 state 升级问题

**描述**: 用户现有的 state 中没有 FunctionId 字段，升级 Provider 版本后需要处理旧 state 的兼容性。

**缓解措施**:
- 新增字段设置为 Computed，确保旧 state 不会报错
- Read 函数会自动从 API 获取 FunctionId 并填充到 state 中
- 第一次 Read 操作会自动完成 state 的更新

### Trade-off: 字段命名

**描述**: API 中的字段名为 FunctionId（驼峰命名），Terraform 中使用 function_id（蛇形命名）。

**权衡**: 遵循 Terraform 的命名约定（蛇形命名）比保持 API 原始命名更重要。这种转换是 Terraform Provider 的常见做法。

**决定**: 使用 function_id 作为 Terraform 中的字段名。

## Migration Plan

**步骤**:

1. 代码变更
   - 修改 tencentcloud/services/teo/resource_tencentcloud_teo_function.go
   - 在 Schema 中添加 function_id 字段定义
   - 更新 Create 函数，添加 FunctionId 的读取和设置逻辑
   - 更新 Read 函数，添加 FunctionId 的读取和设置逻辑

2. 测试更新
   - 更新 resource_tencentcloud_teo_function_test.go
   - 添加验证 function_id 字段的单元测试
   - 添加验证 function_id 字段的验收测试

3. 验证
   - 运行单元测试验证代码变更
   - 运行验收测试（TF_ACC=1）验证与 API 的集成
   - 确认现有测试用例仍然通过（验证向后兼容性）

4. 文档更新（可选，如需要）
   - 更新 website/docs/r/teo_function.md（如果存在）
   - 添加 function_id 字段的说明文档

**回滚策略**:

如果发现严重问题，可以通过以下方式回滚：
- 移除 Schema 中的 function_id 字段定义
- 移除 Create 和 Read 函数中相关的字段读取逻辑
- 由于新增字段不影响现有配置，用户无需修改其 Terraform 代码

## Open Questions

无。当前需求清晰，技术实现路径明确，不存在未解决的技术或业务问题。
