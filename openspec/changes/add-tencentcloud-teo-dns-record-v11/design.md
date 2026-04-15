## Context

当前 Terraform Provider for TencentCloud 已经支持多个云产品的资源管理，但对于 TEO（TencentCloud EdgeOne）产品的 DNS 记录管理还处于空白状态。TEO 是腾讯云的边缘加速和安全服务，DNS 记录是其核心功能之一，用户需要能够通过 Terraform 声明式地管理这些记录。

TEO SDK 提供了完整的 DNS 记录 CRUD 接口：
- CreateDnsRecord: 创建单个 DNS 记录
- DescribeDnsRecords: 查询 DNS 记录列表
- ModifyDnsRecords: 批量修改 DNS 记录
- DeleteDnsRecords: 批量删除 DNS 记录

根据项目规范，新增资源需要遵循以下模式：
- 文件命名: `resource_tc_teo_<name>.go`
- 位置: `tencentcloud/services/teo/`
- 测试文件: `resource_tc_teo_<name>_test.go`
- 文档文件: `website/docs/r/teo_<name>.html.markdown`

## Goals / Non-Goals

**Goals:**
- 实现完整的 `tencentcloud_teo_dns_record_v11` 资源 CRUD 操作
- 支持所有必要的 DNS 记录参数（域名、记录类型、记录值、TTL 等）
- 提供准确的错误处理和状态管理
- 确保异步操作的幂等性和一致性
- 提供完整的单元测试和文档

**Non-Goals:**
- 实现数据源（DataSource）功能（不在本次需求中）
- 实现高级 DNS 管理功能（如批量导入导出）
- 修改现有 TEO 资源的行为

## Decisions

### 1. 资源 ID 格式
**决策**: 使用复合 ID 格式 `<zone_id>#<record_id>`

**理由**:
- DNS 记录属于某个 DNS Zone，需要 Zone ID 来唯一标识
- 记录 ID 在 Zone 内唯一，但全局可能重复
- 复合 ID 符合项目现有模式（如 CVM 等资源使用实例 ID 组合）

**替代方案考虑**:
- 单一 ID: 不可行，Zone 是必需参数
- 带分隔符的其他格式: `zone_id:record_id` 或 `zone_id|record_id`，但 `#` 是项目标准分隔符

### 2. 异步操作处理
**决策**: 对于所有 CRUD 操作，使用 `helper.Retry()` 实现最终一致性重试，并在 Schema 中声明 Timeouts 块

**理由**:
- TEO 的 DNS 记录操作通常是异步的，需要轮询确认状态
- 项目已有成熟的重试机制，最终一致性模式已被广泛采用
- Timeouts 块让用户可以控制超时时间

**实现细节**:
- Create 操作完成后，调用 Read 接口轮询直到记录存在
- Update/Delete 操作完成后，调用 Read 接口轮询直到变更生效
- 默认超时时间设置为 10 分钟

### 3. Schema 字段映射
**决策**: 根据 TEO SDK 的 API 参数直接映射到 Terraform Schema，使用 PascalCase 命名

**理由**:
- SDK API 参数已经过良好设计，直接映射可以避免混淆
- 符合 Terraform Provider 的命名规范
- 简化维护，API 变更时更容易同步

**注意事项**:
- 必填字段在 Schema 中设置为 Required
- 可选字段设置为 Optional，并提供合理的默认值
- 计算字段设置为 Computed

### 4. 错误处理
**决策**: 使用 `defer tccommon.LogElapsed()` 和 `defer tccommon.InconsistentCheck()` 进行统一的错误处理和日志记录

**理由**:
- 项目的标准错误处理模式，已被所有资源采用
- 提供详细的调试信息和性能监控
- 自动处理不一致情况，提高可靠性

## Risks / Trade-offs

### Risk 1: API 变更导致资源失效
**风险**: TEO SDK API 可能在未来版本中发生变更，导致现有资源失效

**缓解措施**:
- 使用 vendor 模式管理 SDK 版本，确保依赖稳定性
- 在代码中添加版本兼容性检查
- 定期更新 SDK 并进行兼容性测试

### Risk 2: 异步操作超时
**风险**: 某些情况下 DNS 记录操作可能需要较长时间，导致超时

**缓解措施**:
- 提供可配置的 Timeouts 块，用户可根据实际情况调整
- 在文档中说明可能的超时情况和解决方案
- 使用合理的默认超时时间（10 分钟）

### Risk 3: 批量操作的原子性
**风险**: ModifyDnsRecords 和 DeleteDnsRecords 是批量操作，部分成功时如何处理

**缓解措施**:
- SDK 批量操作通常返回成功/失败列表，需要解析结果
- 对于 Terraform，每次只操作一个记录，避免批量操作的原子性问题
- 如果批量操作失败，抛出错误让用户重试

### Trade-off 1: 实现复杂度 vs 功能完整性
**权衡**: 简单实现可能导致某些高级功能不支持

**决策**: 优先支持最常用的功能，保持代码简洁
- 基础 CRUD 操作完整实现
- 高级功能（如批量导入）可在未来版本中通过数据源或其他资源实现

### Trade-off 2: 性能 vs 一致性
**权衡**: 频繁的轮询会增加 API 调用次数和延迟

**决策**: 平衡性能和一致性
- Create/Delete: 立即轮询，确保操作生效
- Update: 使用合理的重试间隔（如 5 秒）
- Read: 使用缓存，减少不必要的 API 调用

## Migration Plan

本变更不涉及数据迁移，是新增资源，不修改现有资源。

部署步骤：
1. 实现资源代码和测试
2. 运行单元测试和集成测试
3. 生成并更新文档
4. 合并到主分支

回滚策略：
- 如果发现问题，可以暂时禁用该资源（不注册到 Provider）
- 不影响现有资源和用户配置

## Open Questions

无
