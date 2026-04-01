## Context

当前 `tencentcloud_vpc_end_point` 资源支持基本的 VPC 端点创建和管理功能，但缺少安全组绑定、标签管理和 IP 地址类型配置等关键功能。这些功能在腾讯云 VPC API 中已经完全支持，但 Terraform Provider 尚未实现。

VPC 端点是腾讯云提供的一种用于在 VPC 内部访问私有服务的网络服务。安全组绑定可以增强端点的安全性，标签管理便于资源分类和成本管理，IP 地址类型配置允许用户根据网络环境选择 IPv4 或 IPv6。

文件位置：
- 资源定义：`tencentcloud/services/vpc/resource_tc_vpc_end_point.go`
- 单元测试：`tencentcloud/services/vpc/resource_tc_vpc_end_point_test.go`
- 资源文档：`website/docs/r/vpc_end_point.html.md`

## Goals / Non-Goals

**Goals:**
- 为 `tencentcloud_vpc_end_point` 资源新增 3 个可选字段：SecurityGroupId、Tags、IpAddressType
- 更新 CRUD 操作以支持这些新字段
- 确保向后兼容性，不破坏现有用户配置
- 添加完整的单元测试和验收测试覆盖
- 更新资源文档和示例

**Non-Goals:**
- 不修改现有字段的行为或类型
- 不涉及其他 VPC 资源的修改
- 不添加新的 API 调用或依赖
- 不实现标签的批量操作（如批量打标、解标等）

## Decisions

### 1. Schema 字段定义

**Decision**: 将 3 个字段都定义为 Optional 类型，确保向后兼容性。

**Rationale**:
- 根据 API 文档，这些参数在 CreateVpcEndPoint 中都是可选的
- 添加 Optional 字段不会影响现有 Terraform 配置
- 遵循 Terraform Provider 的最佳实践

**Alternatives considered**:
- 将 IpAddressType 定义为 Computed 可选字段：未采纳，因为该值应该由用户显式配置，而不是从云 API 读取后设置默认值

### 2. Tags 字段结构

**Decision**: 使用 Terraform Plugin SDK 的 `schema.TypeSet` 和 `schema.Resource` 定义 Tags 字段，包含 Key（必填）和 Value（可选）子字段。

**Rationale**:
- 标签的 Key 是唯一的，使用 TypeSet 可以自动处理去重
- 符合腾讯云 Provider 中其他资源标签字段的通用模式
- 支持标签的增删改操作，与 Update API 兼容

**Alternatives considered**:
- 使用 TypeList：未采纳，因为 TypeList 允许重复项，不符合标签的唯一性约束

### 3. 更新策略

**Decision**: 支持所有新字段的更新操作，通过 UpdateVpcEndPointAttribute API 实现。

**Rationale**:
- VPC 端点支持动态更新属性，无需重建资源
- 提供更好的用户体验，避免不必要资源销毁和重建
- 符合 Terraform "update in place" 的最佳实践

**Alternatives considered**:
- 强制重建（ForceNew）：未采纳，因为这些属性可以通过 API 直接更新，无需重建

### 4. 默认值处理

**Decision**: IpAddressType 的默认值 "Ipv4" 在用户未指定时由云 API 处理，不在 Provider 层设置默认值。

**Rationale**:
- 云 API 已经处理了默认值逻辑
- 避免在 Provider 层维护额外的默认值逻辑
- Read 操作时从云 API 读取真实值，确保一致性

**Alternatives considered**:
- 在 schema 中设置 Default: "Ipv4"：未采纳，因为这可能导致与云 API 的默认行为不一致

### 5. 测试策略

**Decision**: 添加单元测试覆盖新字段的 CRUD 操作，使用 Acceptance Tests 进行端到端验证。

**Rationale**:
- 单元测试可以快速验证核心逻辑
- Acceptance Tests 可以验证与真实云 API 的集成
- 符合 Terraform Provider 的测试标准

**Alternatives considered**:
- 仅添加单元测试：未采纳，因为无法验证与云 API 的集成
- 仅添加 Acceptance Tests：未采纳，因为测试耗时较长，不适合频繁运行

## Risks / Trade-offs

### Risk 1: API 兼容性问题
[Risk] 如果腾讯云 API 的行为与文档描述不一致，可能导致功能异常

**Mitigation**: 
- 在开发过程中进行实际 API 调用验证
- 使用 Acceptance Tests 验证所有场景
- 如果发现问题，及时调整实现或与云产品团队沟通

### Risk 2: 标签同步延迟
[Risk] 标签操作可能存在同步延迟，导致 Read 操作获取到不一致的数据

**Mitigation**:
- 使用 tccommon.Retry() 实现最终一致性重试
- 在 Read 操作中添加适当的重试逻辑
- 在文档中说明可能的延迟

### Risk 3: 状态漂移
[Risk] 如果用户直接通过云控制台修改了这些属性，可能导致 Terraform 状态漂移

**Mitigation**:
- 在 Read 操作中完整读取所有属性
- 用户可以通过 `terraform refresh` 同步状态
- 在文档中建议用户尽量通过 Terraform 管理资源

### Risk 4: 向后兼容性
[Risk] 如果未来云 API 修改了这些字段的含义或行为，可能破坏兼容性

**Mitigation**:
- 使用 Optional 字段，确保不影响现有配置
- 密切关注云 API 的变更公告
- 在发布前进行充分的回归测试

## Migration Plan

### 部署步骤

1. **代码变更**: 实现所有新字段和更新逻辑
2. **单元测试**: 验证核心逻辑的正确性
3. **验收测试**: 运行 Acceptance Tests 验证与云 API 的集成
4. **文档更新**: 更新资源文档和示例
5. **代码审查**: 团队成员进行代码审查
6. **发布**: 发布新的 Provider 版本

### 回滚策略

如果发布后发现问题，可以通过以下步骤回滚：

1. 立即停止推荐使用新版本
2. 发布补丁版本，移除或禁用新字段
3. 使用旧版本创建的资源继续使用旧版本管理
4. 已经使用新字段的用户需要手动迁移到旧版本

## Open Questions

1. **问题**: UpdateVpcEndPointAttribute API 是否支持同时更新所有 3 个字段？

   **状态**: 需要验证 API 文档和实际行为

2. **问题**: 标签的更新是增量更新还是全量替换？

   **状态**: 需要验证 API 文档和实际行为

3. **问题**: 如果资源已经绑定了安全组，是否允许更换安全组？

   **状态**: 需要验证 API 文档和实际行为