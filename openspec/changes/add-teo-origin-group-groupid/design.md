## Context

tencentcloud_teo_origin_group 资源当前支持源站组的创建、读取和更新操作，但在删除操作时缺少必要的 GroupId 参数。该参数是 DeleteOriginGroup API 的必填参数，用于标识要删除的源站组。当前实现中，该字段在资源 Schema 中未定义，导致删除功能无法正常工作。

## Goals / Non-Goals

**Goals:**

- 在 tencentcloud_teo_origin_group 资源的 Schema 中添加 `group_id` 字段（string 类型，Required）
- 确保 group_id 字段在 Create 操作中被正确设置和存储
- 确保 group_id 字段在 Read 操作中从 API 响应中正确读取
- 确保 group_id 字段在 Update 操作中正确处理
- 确保 group_id 字段在 Delete 操作中被正确传递给 DeleteOriginGroup API
- 更新单元测试和验收测试以覆盖新增字段
- 保持向后兼容性，不影响现有用户配置

**Non-Goals:**

- 不修改资源的其他字段或行为
- 不引入新的 API 或外部依赖
- 不改变资源的 ID 构成方式
- 不涉及数据迁移或状态转换

## Decisions

### Schema 字段设计

**决定：** 将 `group_id` 字段定义为 Required 字段

**理由：**
- DeleteOriginGroup API 要求 GroupId 为必填参数
- 源站组创建后，API 会返回 GroupId，应该立即存储到资源状态中
- 将其设为 Required 可以确保用户在依赖该字段时能够正确获取

**替代方案考虑：**
- 方案 A：将 group_id 设为 Optional - 不合适，因为 API 返回的值需要存储，且删除操作必需
- 方案 B：将 group_id 设为 Computed - 不合适，因为用户可能在某些场景下需要显式指定该值

### CRUD 操作集成

**决定：** 在所有 CRUD 操作中正确处理 group_id 字段

**Create 操作：**
- 在资源创建成功后，从 API 响应中提取 GroupId
- 使用 d.Set("group_id", groupId) 将值存储到状态中
- 返回时确保 group_id 已正确设置

**Read 操作：**
- 从 API 响应中读取 GroupId 字段
- 使用 d.Set("group_id", groupId) 更新资源状态
- 确保 Read 操作后 group_id 与 API 响应保持一致

**Update 操作：**
- 腾讯云 TEO 的 UpdateOriginGroup API 可能不返回 GroupId（需要确认）
- 如果 API 返回 GroupId，则更新状态
- 如果 API 不返回，则保持现有 group_id 不变

**Delete 操作：**
- 从资源状态中读取 group_id
- 将其作为必填参数传递给 DeleteOriginGroup API 的 GroupId 字段
- 确保 group_id 不为空

### 错误处理

**决定：** 添加必要的错误处理逻辑

**理由：**
- group_id 为删除操作的必填参数，如果缺失应该返回明确的错误
- 遵循 Terraform Provider 的错误处理最佳实践

**实现：**
- 在 Delete 函数中，检查 group_id 是否存在
- 如果缺失，返回 fmt.Errorf("group_id is required for deletion")

### 测试策略

**决定：** 更新单元测试和验收测试以覆盖 group_id 字段

**单元测试：**
- 添加测试用例验证 group_id 字段的 Schema 定义正确
- 添加测试用例验证 Create 操作正确设置 group_id
- 添加测试用例验证 Read 操作正确读取 group_id
- 添加测试用例验证 Delete 操作使用 group_id

**验收测试：**
- 在现有的验收测试中验证 group_id 字段的完整生命周期
- 确保测试覆盖创建、读取、更新、删除操作

## Risks / Trade-offs

### Risk 1: 向后兼容性风险

**风险：** 新增 Required 字段可能会影响现有的资源状态

**缓解措施：**
- group_id 字段虽然在 Schema 中定义为 Required，但实际上是 Computed 类型（由 API 返回）
- 对于已存在的资源状态，Terraform 会在下一次 Refresh 时自动填充该字段
- 用户不需要修改现有的 Terraform 配置文件

### Risk 2: API 兼容性风险

**风险：** Create 或 Read API 可能不返回 GroupId 字段

**缓解措施：**
- 在实现时先通过实际 API 调用验证响应结构
- 如果 API 确实不返回该字段，需要调整实现策略（例如使用其他方式获取）
- 预留错误处理逻辑，当 API 响应不符合预期时给出明确错误信息

### Risk 3: 测试覆盖不完整

**风险：** 测试可能无法完全覆盖所有场景

**缓解措施：**
- 单元测试覆盖每个 CRUD 操作中的 group_id 处理逻辑
- 验收测试覆盖完整的资源生命周期
- 手动测试验证删除操作的正确性

## Migration Plan

### 部署步骤

1. **代码实现**
   - 修改 resource_tencentcloud_teo_origin_group.go 的 Schema 定义
   - 更新 Create、Read、Update、Delete 函数的逻辑
   - 更新资源测试文件

2. **单元测试验证**
   - 运行 `go test` 验证单元测试通过

3. **验收测试验证**
   - 设置 TENCENTCLOUD_SECRET_ID 和 TENCENTCLOUD_SECRET_KEY 环境变量
   - 运行 `TF_ACC=1 go test` 验收测试
   - 确保删除操作正确使用 group_id

4. **文档更新**
   - 更新 website/docs/ 目录下的资源文档
   - 添加 group_id 字段的说明

### 回滚策略

如果出现问题，可以通过以下方式回滚：
- 回退代码修改
- 确保已创建的资源不受影响（因为 group_id 是新增字段）

## Open Questions

**Q1: UpdateOriginGroup API 是否返回 GroupId 字段？**

- 需要通过实际 API 调用或查看 API 文档确认
- 如果不返回，Update 操作中应该如何处理 group_id

**Q2: 资源 ID 与 group_id 的关系是什么？**

- 需要确认当前资源的 ID 构成
- group_id 是否应该作为资源 ID 的一部分
- 如果是 ID 的一部分，可能需要调整 ID 解析逻辑
