## Context

当前 Terraform Provider for TencentCloud 已经支持 TEO（TencentCloud EdgeOne）服务的多个资源，但缺少查询边缘 KV 数据的 Terraform Resource。TEO 服务的边缘 KV（Edge KV）提供分布式的键值对存储，允许在边缘节点快速存储和检索数据。

CAPI 接口 EdgeKVGet 提供了查询 KV 数据的能力，支持通过站点 ID（ZoneId）、命名空间（Namespace）和键名列表（Keys）批量查询键值对数据。该资源需要集成到现有的 Terraform Provider 架构中，遵循项目的代码组织模式和约定。

约束条件：
- 必须使用 Terraform Plugin SDK v2
- 必须保持向后兼容，不影响现有资源
- 必须实现完整的 CRUD 操作
- 必须包含单元测试和验收测试
- 必须提供文档示例

## Goals / Non-Goals

**Goals:**
- 创建完整的 Terraform Resource `tencentcloud_teo_edge_k_v_get`，支持通过 ZoneId、Namespace 和 Keys 查询 KV 数据
- 实现标准的 Create、Read、Update、Delete 操作函数，与 Terraform 状态管理机制集成
- 遵循 Terraform Provider for TencentCloud 的代码模式和约定
- 提供完整的测试覆盖（单元测试和验收测试）
- 编写清晰的文档和示例代码

**Non-Goals:**
- 不实现 KV 数据的创建、修改、删除操作（仅查询）
- 不涉及其他 TEO 资源的变更
- 不改变现有的 TEO 服务架构或 API 调用方式

## Decisions

### 1. 资源类型选择
**决策**: 将该功能实现为 Terraform Resource 而非 Data Source

**理由**:
- 用户需求明确要求实现 Resource 类型
- Resource 类型支持状态管理和刷新机制，适合需要持久化查询结果的使用场景
- 虽然 Data Source 更适合纯查询场景，但需求要求实现完整的 CRUD 操作，Resource 类型更符合要求

**替代方案**:
- Data Source: 更适合纯查询场景，但不支持 Update 和 Delete 操作

### 2. 资源 ID 构造
**决策**: 使用复合 ID 格式 `zoneId#namespace#keys[0]`（使用第一个键名作为 ID 的一部分）

**理由**:
- ZoneId 和 Namespace 是必填参数，可以作为 ID 的一部分
- Keys 是列表类型，但通常查询操作是幂等的，使用第一个键名可以保证唯一性
- 符合项目中使用 `#` 作为分隔符的惯例

**替代方案**:
- 使用整个 Keys 列表的 hash 作为 ID 部分: 更准确但增加复杂度
- 仅使用 ZoneId 和 Namespace: 无法区分不同的键名查询

### 3. Schema 设计
**决策**:
- ZoneId (string, Required): 站点 ID
- Namespace (string, Required): 命名空间名称
- Keys (list of string, Required): 键名列表
- Data (list of object, Computed): 键值对数据列表
  - Key (string, Computed): 键名
  - Value (string, Computed): 键值
  - Expiration (string, Computed): 过期时间

**理由**:
- 请求参数（ZoneId、Namespace、Keys）设置为 Required，因为都是必填项
- 响应参数（Data 及其子字段）设置为 Computed，因为由 API 返回
- 使用嵌套对象结构清晰表达键值对数据

**替代方案**:
- 将 Data 展平为三个单独的列表: 降低可读性

### 4. CRUD 操作实现策略
**决策**:
- **Create**: 调用 EdgeKVGet API，将查询结果保存到 Terraform 状态
- **Read**: 调用 EdgeKVGet API，更新 Terraform 状态（保持数据最新）
- **Update**: 调用 EdgeKVGet API（使用最新的 Keys），更新 Terraform 状态
- **Delete**: 从 Terraform 状态中删除资源（不调用 API，因为是查询操作）

**理由**:
- EdgeKVGet 是查询接口，不支持实际创建资源
- Create 和 Update 操作本质上是查询并保存结果
- Delete 操作仅清除状态，不影响云端数据
- 这种设计符合 Terraform 的 Resource 语义，同时适应查询型 API

**替代方案**:
- 仅实现 Read 操作: 不符合 Terraform Resource 的标准要求

### 5. 错误处理和重试
**决策**: 使用 `helper.Retry()` 和 `tccommon.InconsistentCheck()` 处理最终一致性

**理由**:
- TEO 服务是分布式系统，可能存在最终一致性
- 项目中已有成熟的重试机制，应该复用
- 使用 `tccommon.LogElapsed()` 记录操作耗时

**替代方案**:
- 不处理重试: 可能导致偶发性失败

### 6. 测试策略
**决策**:
- 单元测试: 测试 Schema 定义、状态管理逻辑
- 验收测试: 使用 TF_ACC=1 运行真实的 API 调用

**理由**:
- 单元测试快速验证核心逻辑
- 验收测试确保与真实 API 的集成正常工作
- 项目约定要求同时提供两种测试

**替代方案**:
- 仅提供单元测试: 无法验证 API 集成
- 仅提供验收测试: 运行成本高，调试困难

## Risks / Trade-offs

**风险 1**: Keys 列表顺序变更导致状态差异
- **缓解措施**: 在 ResourceData 中保存 Keys 列表的原始顺序，Read 时保持一致

**风险 2**: CAPI 接口返回的 Data 顺序与 Keys 顺序不一致
- **缓解措施**: 根据 Key 字段进行匹配，而不是依赖数组索引顺序

**风险 3**: 查询大量键名时可能超时（API 限制为 20 个键）
- **缓解措施**: 在 Schema 文档中明确说明限制，用户需要分批查询

**权衡**: 将查询功能实现为 Resource 会引入状态管理的复杂性，但用户需求明确要求这样做，且 Terraform Resource 提供更好的生命周期管理

## Migration Plan

**部署步骤**:
1. 创建资源文件 `resource_tencentcloud_teo_edge_k_v_get.go`
2. 创建测试文件 `resource_tencentcloud_teo_edge_k_v_get_test.go`
3. 运行单元测试验证核心逻辑
4. 运行验收测试验证 API 集成
5. 创建文档文件 `website/docs/r/teo_edge_k_v_get.html.markdown`
6. 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中注册资源

**回滚策略**:
- 如果发现问题，可以随时删除资源文件，不影响现有功能
- 保持向后兼容，不会破坏现有用户的 Terraform 配置

## Open Questions

（无）
