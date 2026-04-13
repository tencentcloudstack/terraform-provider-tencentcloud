## Context

TEO (Tencent EdgeOne) 是腾讯云的边缘加速产品，提供站点配置管理功能。当前用户需要能够通过 Terraform Provider 导出站点配置，以便进行版本控制和配置备份。

当前状态：
- Terraform Provider 已支持 TEO 服务的其他资源
- 缺少导出站点配置的资源
- 用户无法通过 Terraform 导出站点配置

约束：
- 必须使用 Terraform Plugin SDK v2
- 必须保持与现有代码风格一致
- 必须支持完整的 CRUD 操作
- 必须提供完整的单元测试和验收测试
- 必须遵守 Terraform Provider 的最佳实践

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_teo_export_zone_config` 资源，支持导出 TEO 站点配置
- 提供完整的 Create、Read、Update、Delete 操作
- 生成符合规范的单元测试和验收测试代码
- 遵循现有代码结构和风格

**Non-Goals:**
- 修改现有的 TEO 资源
- 实现站点配置的导入功能
- 提供配置的自动备份或版本控制功能

## Decisions

1. **使用 Resource Schema 映射 CAPI 接口参数**
   - 直接根据 CAPI 接口的请求参数生成 Terraform Resource Schema
   - 保持参数的 Required/Optional 属性与 CAPI 接口定义一致
   - 理由：确保与云 API 的兼容性和一致性

2. **资源 ID 结构**
   - 使用复合 ID 格式：`zoneId#exportId` 或类似的唯一标识符
   - 理由：遵循项目现有模式，便于资源定位和管理

3. **异步操作处理**
   - 在 Schema 中声明 Timeouts 块
   - 在 CRUD 函数中使用 context 进行超时控制
   - 理由：满足异步操作的可靠性和可控制性要求

4. **错误处理**
   - 使用 `defer tccommon.LogElapsed()` 记录操作耗时
   - 使用 `defer tccommon.InconsistentCheck()` 进行一致性检查
   - 理由：遵循项目现有模式，便于问题排查和调试

5. **测试策略**
   - 单元测试：测试资源的基本功能
   - 验收测试：使用 `TF_ACC=1` 运行，需要真实的云环境凭证
   - 理由：确保代码质量和功能的正确性

## Risks / Trade-offs

1. **CAPI 接口变更风险**
   - 风险：CAPI 接口可能发生变更，导致 Terraform 资源不兼容
   - 缓解：定期检查接口变更，及时更新资源代码

2. **最终一致性延迟**
   - 风险：TEO 服务可能存在最终一致性延迟，导致 Read 操作返回过期数据
   - 缓解：使用 `helper.Retry()` 实现重试机制

3. **测试环境依赖**
   - 风险：验收测试需要真实的云环境凭证，可能在 CI/CD 环境中难以执行
   - 缓解：提供 Mock 测试或使用测试账号

## Open Questions

1. CAPI 接口的具体参数和返回值结构（需要在实施阶段获取）
2. 导出配置的具体格式和内容
3. 资源的更新操作是否支持（某些导出操作可能不支持更新）
