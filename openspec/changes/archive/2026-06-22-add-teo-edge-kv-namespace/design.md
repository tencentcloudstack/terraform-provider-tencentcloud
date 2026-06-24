## Context

TencentCloud TEO (Tencent EdgeOne) 提供 Edge KV 命名空间功能，用于在边缘节点管理 KV 存储。云 API 已在 vendor 中提供了完整的 CRUD 接口（CreateEdgeKVNamespace、DescribeEdgeKVNamespaces、ModifyEdgeKVNamespace、DeleteEdgeKVNamespace），均为同步接口。当前 provider 中已有 TEO 相关资源（如 `tencentcloud_teo_zone`），服务层代码已存在。

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_teo_edge_k_v_namespace` 资源的完整 CRUD 生命周期管理
- 遵循项目现有代码风格（参考 `tencentcloud_igtm_strategy` 资源）
- 使用 `zone_id#namespace` 联合 ID 标识资源
- 支持 import 功能
- 提供单元测试（使用 gomonkey mock 方式）

**Non-Goals:**
- 不实现 Edge KV 的键值对操作（仅管理命名空间）
- 不实现 datasource 类型资源
- 不暴露 DescribeEdgeKVNamespaces 的分页参数（Offset/Limit）给用户，内部使用 Filters 按 namespace 精确查询

## Decisions

1. **资源 ID 设计**: 使用 `zone_id#namespace` 作为联合 ID（通过 `tccommon.FILED_SP` 分隔），因为 Create 接口不返回独立 ID 字段，而 `zone_id` + `namespace` 在站点内唯一标识一个命名空间。

2. **Read 实现**: 使用 DescribeEdgeKVNamespaces 接口，通过 Filters 按 namespace 名称精确过滤（Fuzzy=false），Limit 设为最大值 1000。从返回列表中匹配目标 namespace。

3. **ForceNew 字段**: `zone_id` 和 `namespace` 设为 ForceNew，因为这两个字段构成资源标识，不可变更。`remark` 字段可通过 ModifyEdgeKVNamespace 接口更新。

4. **Schema 设计**:
   - `zone_id`: Required, ForceNew, String - 站点 ID
   - `namespace`: Required, ForceNew, String - 命名空间名称
   - `remark`: Optional, String - 命名空间描述
   - `capacity`: Computed, Int - KV 存储空间可用容量（只读）
   - `capacity_used`: Computed, Int - KV 存储空间已用容量（只读）
   - `created_on`: Computed, String - 创建时间（只读）
   - `modified_on`: Computed, String - 最后修改时间（只读）

5. **错误处理**: 所有 API 调用使用 `resource.Retry` + `tccommon.ReadRetryTimeout` 进行重试，失败时使用 `tccommon.RetryError()` 包装错误。

## Risks / Trade-offs

- [Risk] Create 接口不返回业务 ID → 使用 zone_id + namespace 组合作为 ID，依赖 namespace 在站点内唯一性（云 API 已保证）
- [Risk] DescribeEdgeKVNamespaces 返回列表而非单个资源 → 通过 Filters 精确匹配 namespace 名称，确保查询结果准确
- [Risk] namespace 字段不可修改 → 设为 ForceNew，修改时触发资源重建
