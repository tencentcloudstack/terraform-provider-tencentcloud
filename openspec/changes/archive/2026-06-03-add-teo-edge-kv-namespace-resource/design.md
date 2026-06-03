## Context

TencentCloud TEO (EdgeOne) 提供 Edge KV 命名空间管理的云 API，包括 CreateEdgeKVNamespace、DescribeEdgeKVNamespaces、ModifyEdgeKVNamespace、DeleteEdgeKVNamespace 四个接口。这些接口已在 vendor 中可用（`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`）。当前 provider 已有 TEO 服务目录（`tencentcloud/services/teo/`），需要在其中新增资源文件。

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_teo_edge_k_v_namespace` 资源的完整 CRUD 生命周期
- 使用 `zone_id#namespace` 联合 ID，支持 terraform import
- 遵循项目现有代码风格（参考 `tencentcloud_igtm_strategy` 资源）
- 提供单元测试（使用 gomonkey mock 云 API）
- 提供资源文档（.md 文件）

**Non-Goals:**
- 不实现 Edge KV 键值对的 CRUD（仅管理命名空间）
- 不实现 datasource 类型资源
- 不新增 service 层文件（直接在资源文件中调用 API）

## Decisions

### 1. 资源 ID 设计：使用 zone_id + namespace 联合 ID

**选择**: 使用 `tccommon.FILED_SP`（即 `#`）作为分隔符拼接 `zone_id` 和 `namespace` 作为资源 ID。

**理由**: CreateEdgeKVNamespace 接口不返回独立的资源 ID，命名空间通过 zone_id + namespace 唯一标识。这与项目中其他类似资源的做法一致。

**替代方案**: 仅使用 namespace 作为 ID → 不可行，因为 namespace 在不同 zone 下可能重复。

### 2. ForceNew 字段设计

**选择**: `zone_id` 和 `namespace` 设为 ForceNew。

**理由**: ModifyEdgeKVNamespace 接口仅支持修改 `remark` 字段，`zone_id` 和 `namespace` 是定位资源的标识符，不可变更。

### 3. Update 方法设计

**选择**: Update 方法仅处理 `remark` 字段的变更，通过 ModifyEdgeKVNamespace 接口实现。

**理由**: 云 API ModifyEdgeKVNamespace 仅支持修改 remark。zone_id 和 namespace 已设为 ForceNew，变更时会触发 destroy + create。

### 4. Read 方法设计：使用 Filters 过滤

**选择**: 在 DescribeEdgeKVNamespaces 请求中使用 `namespace` 过滤条件精确查询目标命名空间，Limit 设为 1000（接口最大值）。

**理由**: DescribeEdgeKVNamespaces 是列表查询接口，通过 Filters 按 namespace 名称过滤可精确定位目标资源。设置最大 Limit 确保不会因分页遗漏。

### 5. 测试方案：gomonkey mock

**选择**: 使用 gomonkey 对云 API 调用进行 mock，编写纯业务逻辑单元测试。

**理由**: 项目要求新增资源使用 mock 方式进行单元测试，不依赖真实云环境。

## Risks / Trade-offs

- [DescribeEdgeKVNamespaces 使用模糊匹配] → 在 Read 方法中遍历返回结果，精确匹配 namespace 名称，避免模糊查询返回多个结果时取错数据。
- [CreateEdgeKVNamespace 不返回资源 ID] → 创建成功后直接使用入参 zone_id + namespace 组合为 ID，无需从响应中提取。
- [并发创建同名 namespace] → 依赖云 API 侧的唯一性校验，Terraform 层面不做额外处理。
