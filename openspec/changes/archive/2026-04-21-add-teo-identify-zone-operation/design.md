## Context

EdgeOne (TEO) 是腾讯云的边缘加速服务，在使用前需要对站点进行所有权验证。当前 Terraform Provider for TencentCloud 已经支持 TEO 的其他资源（如站点、规则引擎等），但缺少站点认证操作的 Terraform 资源支持。用户需要手动在控制台完成站点验证，这破坏了基础设施即代码（IaC）的完整流程。

TencentCloud SDK for Go 已经提供了 `IdentifyZone` 接口（在 teo/v20220901 包中），支持获取 DNS 校验和文件校验的配置信息。该接口返回认证所需的详细配置，用户可以根据这些信息在 DNS 服务或文件系统完成站点验证。

## Goals / Non-Goals

**Goals:**
- 创建 `tencentcloud_teo_identify_zone` Terraform 资源，支持执行站点认证操作
- 实现创建接口调用 `IdentifyZone` 云 API
- 支持传入站点名称和可选的子域名参数
- 返回 DNS 校验信息和文件校验信息
- 提供完整的单元测试，使用 mock 方式对云 API 进行模拟

**Non-Goals:**
- 不提供读取、更新、删除接口（一次性操作）
- 不自动完成 DNS 配置或文件创建，仅返回配置信息
- 不处理异步操作（IdentifyZone 接口是同步返回的）

## Decisions

### 资源类型选择
决定使用 `RESOURCE_KIND_OPERATION` 资源类型。这是因为站点认证是一个一次性操作，操作完成后不需要维护任何状态。用户调用资源获取认证配置后，可以自行完成 DNS 或文件验证。

**替代方案考虑:**
- 使用标准资源类型（带状态管理）：不适用，因为认证配置不需要持久化在 Terraform state 中
- 使用数据源：不适用，因为数据源用于查询已有资源，而认证操作是一个主动的动作

### Schema 设计
Schema 中包含必需参数和可选参数：
- `zone_name` (Required): 站点名称
- `domain` (Optional): 子域名，仅验证站点下的子域名时需要
- `ascription` (Computed): DNS 校验信息，包含 Subdomain, RecordType, RecordValue
- `file_ascription` (Computed): 文件校验信息，包含 IdentifyPath, IdentifyContent

`ascription` 和 `file_ascription` 设计为嵌套对象（SchemaTypeList with MaxItems=1），而不是扁平字段，这样可以保持与云 API 结构的一致性，也便于后续扩展。

**替代方案考虑:**
- 将 ascription 的字段扁平化到根级别：会失去结构化信息的语义，不便于理解和管理
- 使用 SchemaTypeMap: 不适合，因为字段是固定的且有明确的类型和语义

### 测试策略
决定使用 mock（gomonkey）方式对云 API 进行单元测试。这是因为：
1. 认证操作需要真实的站点和 DNS 环境，集成测试成本较高
2. 云 API 接口行为确定，可以通过 mock 验证业务逻辑正确性
3. 符合 TEO 服务现有资源的测试模式

**替代方案考虑:**
- 使用 Terraform 测试套件（TF_ACC）：需要真实云环境和站点配置，不适合纯业务逻辑测试

### 文件命名规范
按照 TEO 服务现有资源的命名规范：
- 主文件：`resource_tc_teo_identify_zone_operation.go`（明确标识为 operation 类型）
- 测试文件：`resource_tc_teo_identify_zone_operation_test.go`

**替代方案考虑:**
- 使用 `resource_tc_teo_identify_zone.go`: 不够明确，无法区分为一次性操作

### 错误处理
- 使用 Terraform Plugin SDK v2 的标准错误处理机制
- 使用 `defer tccommon.LogElapsed()` 记录耗时
- 使用 `defer tccommon.InconsistentCheck()` 处理最终一致性错误
- 对于云 API 返回的错误，通过 `diag.Errorf()` 返回给用户

## Risks / Trade-offs

### Risk 1: 云 API 参数变化
**风险**: 云 API 可能在未来版本中新增或修改参数，导致 Terraform 资源与实际 API 不匹配
**缓解措施**: 依赖 vendor 目录中的云 SDK 版本，定期更新 vendor；在代码中添加注释说明 API 版本

### Risk 2: 测试覆盖不足
**风险**: 单元测试可能无法完全覆盖所有边界场景，导致生产环境问题
**缓解措施**: 设计多个测试用例覆盖不同场景（成功、失败、缺少必需参数等）；定期审查测试覆盖率

### Trade-off 1: Schema 复杂度 vs 可用性
**权衡**: 使用嵌套对象结构增加了 Schema 的复杂度，但提高了可维护性和扩展性
**决策**: 选择嵌套对象，因为代码可维护性优先于轻微的复杂度增加

### Trade-off 2: 测试方式 vs 真实性
**权衡**: 单元测试（mock）无法验证真实云 API 的集成，但开发成本低、执行速度快
**决策**: 使用单元测试，因为认证操作的核心是业务逻辑，不是 API 集成

## Migration Plan

由于这是一个全新资源，不需要迁移计划。用户可以直接使用新资源 `tencentcloud_teo_identify_zone`。

## Open Questions

（无）
