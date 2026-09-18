## Context

腾讯云 MongoDB 提供 CPU 弹性扩容/回缩功能，通过 `ScaleUpDBInstanceCpu`（开启弹性扩容）和 `ScaleDownDBInstanceCpu`（关闭弹性扩容回缩）两个云 API 管理。当前 Terraform Provider 缺少对 MongoDB 实例 CPU 弹性扩容配置的声明式管理能力。

本资源为 RESOURCE_KIND_CONFIG 类型，管理的是已存在 MongoDB 实例的 CPU 弹性扩容配置，而非独立的云资源。配置资源的生命周期为：创建（开启弹性扩容）→ 读取（查询实例状态）→ 更新（调整额外 CPU 核数）→ 删除（关闭弹性扩容回缩）。

### 参考实现
- `tencentcloud_mongodb_instance_ssl`（同产品 config 资源）
- `tencentcloud_igtm_strategy`（代码风格参考）

### 涉及云 API
| API | 包路径 | 用途 |
|-----|--------|------|
| ScaleUpDBInstanceCpu | mongodb/v20190725 | 开启 CPU 弹性扩容 |
| ScaleDownDBInstanceCpu | mongodb/v20190725 | 关闭弹性扩容并回缩 |
| DescribeDBInstances | mongodb/v20190725 | 查询实例详情（Read 验证用） |

## Goals / Non-Goals

**Goals:**
- 提供 `tencentcloud_mongodb_instance_cpu_scale` 配置型资源，管理 MongoDB 实例的 CPU 弹性扩容
- 支持通过 `extra_cpu` 字段设置额外 CPU 核数，开启弹性扩容
- 支持移除 `extra_cpu` 或设置 `extra_cpu=0` 来关闭弹性扩容
- 支持 Import 已有配置
- 遵循项目中 config 资源的标准实现模式

**Non-Goals:**
- 不管理 MongoDB 实例本身（由 `tencentcloud_mongodb_instance` 资源管理）
- 不暴露实例详细信息作为 terraform 输出（配置资源范围限定在 CPU 弹性扩容相关）
- 不支持批量操作多个实例

## Decisions

### 决策 1：资源 ID 使用 instance_id

使用实例 ID（如 `cmgo-xxxxxxxx`）作为资源的唯一标识，与 mongodb ssl config 资源一致。

**理由**：CPU 弹性扩容配置与实例一一对应，实例 ID 天然唯一。简单、直观，便于 Import。

### 决策 2：使用 extra_cpu 字段的有无来判断弹性扩容开关

`extra_cpu` 为 Optional 字段。当用户设置 `extra_cpu=N`（N>0）时，视为开启弹性扩容；当用户不设置 `extra_cpu` 或移除它时，视为关闭弹性扩容。

**理由**：`ScaleUpDBInstanceCpu` API 需要 `ExtraCpu` 参数，`ScaleDownDBInstanceCpu` 仅需 `InstanceId`。用 `extra_cpu` 的有无区分两种操作，语义清晰。

### 决策 3：Read 使用 DescribeDBInstances 验证实例存在

`DescribeDBInstances` 通过 `InstanceIds` 参数查询指定实例，返回 `InstanceDetails` 列表。Read 方法中使用此 API 验证实例是否存在，但不将实例详细字段映射回 terraform schema。

**理由**：`DescribeDBInstances` 返回的 `InstanceDetail` 中只有 `CpuNum`（基础 CPU 核数），没有 `extra_cpu`（弹性扩容核数），无法读取当前的弹性扩容配置。因此 Read 仅做存在性验证，配置值由 terraform state 维护。

### 决策 4：Update 通过判断 extra_cpu 变化选择调用 ScaleUp 或 ScaleDown

- 若 `extra_cpu` 从无到有或值发生变更：调用 `ScaleUpDBInstanceCpu`
- 若 `extra_cpu` 被移除（从有到无）：调用 `ScaleDownDBInstanceCpu`
- 若 `extra_cpu` 不变且实例存在：无需操作

**理由**：两个 API 分别对应开启和关闭操作，Update 方法需要根据参数变化路由到正确的 API。

### 决策 5：异步操作处理

`ScaleUpDBInstanceCpu` 和 `ScaleDownDBInstanceCpu` 为异步接口，返回 `FlowId`。在 Create 和 Update 调用后，需要通过 Read 轮询等待操作生效。

**理由**：遵循项目规范——"对于标注为异步接口的接口，调用完后要调用 Read 接口轮询直到接口生效"。

### 决策 6：CRUD 方法与 API 映射

| Terraform 方法 | 调用的云 API | 说明 |
|----------------|-------------|------|
| Create | ScaleUpDBInstanceCpu | 传入 instance_id + extra_cpu，开启弹性扩容 |
| Read | DescribeDBInstances | 通过 instance_id 查询验证实例存在 |
| Update (extra_cpu 变更) | ScaleUpDBInstanceCpu | 调整额外 CPU 核数 |
| Update (移除 extra_cpu) | ScaleDownDBInstanceCpu | 关闭弹性扩容 |
| Delete | ScaleDownDBInstanceCpu | 关闭弹性扩容回缩 |

## Risks / Trade-offs

### 风险 1：无法从云 API 读取当前 extra_cpu 值
**描述**：`DescribeDBInstances` 不返回弹性扩容的额外 CPU 核数，导致 Read 无法验证 state 中 `extra_cpu` 与实际云侧配置是否一致。

**缓解措施**：Read 验证实例存在性即可，`extra_cpu` 依赖 terraform state 维护。若用户在控制台手动修改了弹性扩容配置，terraform 会在下一次 plan 时检测到 drift（但仅限通过 `instance_id` 匹配，`extra_cpu` 的实际 drift 无法检测）。

### 风险 2：ScaleDownDBInstanceCpu 仅需 InstanceId
**描述**：关闭弹性扩容时不需要传入原始 `extra_cpu` 值，API 会自动将 CPU 回缩到原始规格。

**缓解措施**：Delete 和 Update（移除 extra_cpu 时）直接调用 ScaleDownDBInstanceCpu，无需额外处理。

### 风险 3：并发修改冲突
**描述**：如果多个 terraform 配置或控制台同时操作同一实例的 CPU 弹性扩容，可能导致不一致。

**缓解措施**：`ScaleUpDBInstanceCpu` 和 `ScaleDownDBInstanceCpu` API 可能返回 `InvalidParameterValue.LockFailed` 错误码，terraform 的重试机制会捕获并重试处理。