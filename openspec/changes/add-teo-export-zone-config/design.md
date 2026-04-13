## Context

Terraform Provider for TencentCloud 已经支持 TEO（TencentCloud Edge One）服务的部分功能，但目前缺少导出站点配置（ExportZoneConfig）的资源支持。TEO 是腾讯云边缘安全加速平台，提供全球边缘节点网络加速和安全防护能力。导出站点配置功能允许用户获取站点当前配置的详细信息，用于备份、审计或迁移场景。

当前 TEO 服务的其他资源已经建立了标准的实现模式，包括通过 tencentcloud-sdk-go 调用云 API、使用 helper.Retry() 处理最终一致性、以及标准的错误处理模式。

## Goals / Non-Goals

**Goals:**
- 实现完整的 `tencentcloud_teo_export_zone_config` 资源，包括 Create、Read、Update、Delete 四个操作
- 根据 CAPI 接口生成准确的 Resource Schema 定义，确保 Required/Optional 属性与 API 定义一致
- 提供完整的单元测试和验收测试覆盖
- 遵循现有 TEO 服务的代码规范和模式
- 添加完整的文档和示例

**Non-Goals:**
- 不修改现有 TEO 资源的行为或 schema
- 不引入新的外部依赖（复用现有的 tencentcloud-sdk-go）
- 不涉及复杂的数据迁移或状态迁移

## Decisions

**1. 资源 ID 策略**
- 决策：使用复合 ID，格式为 `zone_id`，因为导出站点配置是基于特定站点的操作
- 理由：与其他 TEO 资源保持一致，简化状态管理

**2. Schema 定义来源**
- 决策：根据资源 UID `iacpres-ZHk6oZ2uSM` 对应的 CAPI 接口定义生成 Schema
- 理由：确保与云 API 的接口定义完全一致，减少人为错误

**3. 异步操作处理**
- 决策：在 Schema 中声明 Timeouts 块，在 CRUD 函数中使用 ctx 支持异步操作
- 理由：遵循 Provider 的硬约束要求，提供更好的用户体验

**4. 错误处理策略**
- 决策：使用标准错误处理模式：`defer tccommon.LogElapsed()` 和 `defer tccommon.InconsistentCheck()`
- 理由：与现有代码保持一致，提供统一的错误日志和一致性检查

**5. 测试策略**
- 决策：同时生成单元测试（*_test.go）和验收测试（TF_ACC=1）
- 理由：确保代码质量和功能正确性，验收测试验证与实际云 API 的集成

## Risks / Trade-offs

**[风险 1] CAPI 接口参数变更** → 如果 CAPI 接口在实现过程中发生变更，可能需要调整 Schema 和实现
- 缓解措施：在实现前确认 CAPI 接口的最新版本，并在代码中添加注释说明接口版本

**[风险 2] 缺少实际测试环境** → 验收测试需要真实的 TEO 资源和密钥
- 缓解措施：确保测试文档中明确说明所需的环境变量（TENCENTCLOUD_SECRET_ID/KEY）

**[权衡 1] 代码生成 vs 手动实现** → 自动生成代码效率高但可能不够灵活
- 决策：采用自动生成基础代码，然后根据实际情况进行微调的策略

**[权衡 2] 导出操作的幂等性** → ExportZoneConfig 是查询类操作，多次调用结果可能不同
- 决策：在 Read 操作中处理，确保 Terraform state 的正确更新
