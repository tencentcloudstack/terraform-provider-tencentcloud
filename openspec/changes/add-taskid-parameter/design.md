## Context

**Background:**
tencentcloud_teo_l7_acc_rule 资源用于管理 TEO 的七层访问规则。该资源的 Create/Update 操作通过 ImportZoneConfig API 实现配置导入。最近 ImportZoneConfig API 在响应中新增了 `TaskId` 字段，用于表示导入配置的任务 ID。

**Current State:**
- tencentcloud_teo_l7_acc_rule 资源已存在，实现了基本的 CRUD 操作
- 资源 schema 定义在 `tencentcloud/services/teo/resource_tencentcloud_teo_l7_acc_rule.go`
- Read 函数调用 ImportZoneConfig API 读取配置状态
- 当前未读取和映射 TaskId 字段

**Constraints:**
- 必须保持向后兼容，不能破坏现有 TF 配置和 state
- 新增字段必须是 Optional，不能是 Required
- 使用现有的 tencentcloud-sdk-go TEO 服务包，不引入新依赖

## Goals / Non-Goals

**Goals:**
- 在 tencentcloud_teo_l7_acc_rule 资源的 schema 中添加 `task_id` 字段（Computed + Optional）
- 在 Read 函数中从 ImportZoneConfig API 响应中读取并映射 TaskId 到 task_id
- 确保新字段不影响现有的 Create/Update/Delete 操作
- 添加相应的单元测试和验收测试验证字段读取逻辑

**Non-Goals:**
- 不修改 TaskId 在 API 层面的行为（仅读取，不发送）
- 不添加对 DescribeZoneConfigImportResult API 的支持（TaskId 仅用于用户查询）
- 不修改其他 TEO 资源

## Decisions

### 字段 Schema 定义
**Decision:** 使用 `schema.TypeString` 配合 `Computed: true, Optional: true`

**Rationale:**
- `Computed: true` 表示该字段由 API 返回填充，用户无需提供
- `Optional: true` 确保向后兼容，现有 state 升级时不会报错
- 不使用 `ForceNew` 因为该字段仅用于读取

**Alternatives Considered:**
- 仅使用 `Computed: true`：会导致现有 state 升级时出现 schema 错误
- 使用 `Optional: true` 不加 `Computed`：用户可能误填该字段

### 字段映射时机
**Decision:** 仅在 Read 函数中映射 TaskId

**Rationale:**
- TaskId 仅在 ImportZoneConfig API 响应中返回，Create/Update 不需要该字段
- 避免在 Create/Update 中发送无意义的请求参数
- 保持代码简洁，减少不必要的逻辑

**Alternatives Considered:**
- 在 Create/Update 中也读取并填充：增加复杂度，没有实际价值

### 测试策略
**Decision:** 添加单元测试和验收测试

**Rationale:**
- 单元测试验证字段映射逻辑
- 验收测试验证实际 API 调用返回的字段
- 确保后续改动不会破坏该功能

**Alternatives Considered:**
- 仅添加验收测试：无法快速验证边界情况
- 不添加测试：违反测试覆盖率要求

## Risks / Trade-offs

### Risk 1: API 字段返回 null 或空值
**Mitigation:** 使用 `d.Set()` 直接设置，允许空值。Terraform SDK 会处理 nil 情况

### Risk 2: TaskId 字段在未来版本中被移除
**Mitigation:** 字段为 Optional，即使 API 不返回该字段也不会导致资源读取失败

### Risk 3: 向后兼容性问题
**Mitigation:** 字段定义为 Computed + Optional，现有用户不会受到影响。新用户可以忽略该字段

### Trade-off: 仅支持读取不提供查询功能
用户需要自行调用 DescribeZoneConfigImportResult API 查询任务结果，但这是合理的设计，因为 task_id 仅作为信息返回

## Migration Plan

该变更不涉及数据库迁移或 state schema 破坏性变更：

1. 用户升级 provider 版本后，state 会自动包含新的 `task_id` 字段（为空）
2. 下次执行 `terraform refresh` 或 `terraform apply` 时，provider 会读取并填充 `task_id`
3. 无需手动迁移操作
4. 如需回滚，用户可降级 provider 版本（不会影响现有配置）

## Open Questions

无 outstanding issues。变更范围明确，实现路径清晰。
