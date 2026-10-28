## Context

`tencentcloud_mongodb_instance` 是腾讯云 MongoDB 实例的通用资源（RESOURCE_KIND_GENERAL），管理实例整个生命周期。当前资源通过 `MongodbService.UpgradeInstance` 方法（封装 `ModifyDBInstanceSpec` 接口）支持 `memory`、`volume`、`node_num`、`add_node_list`、`remove_node_list`、`in_maintenance` 等参数的变配，但未暴露 `cpu` 参数。

腾讯云 MongoDB `ModifyDBInstanceSpec` 接口已支持 `Cpu` 入参（`*int64`，单位：C），允许在变更实例规格时指定变更后的 CPU 大小；`DescribeDBInstances` 返回的 `InstanceDetail.CpuNum`（`*uint64`）字段记录实例当前 CPU 核数。本次变更将 `cpu` 参数接入 Terraform 资源，补全实例规格变更能力。

约束：
- 资源 schema 已存在且已发布，必须保持向后兼容，只能新增 `Optional` + `Computed` 字段。
- `ModifyDBInstanceSpecRequest.Cpu` 为 `*int64` 类型；`InstanceDetail.CpuNum` 为 `*uint64` 类型，需注意类型转换。
- 变配操作可能返回订单 ID，需通过 `DescribeDBInstanceDeal` 轮询订单状态确认完成（现有逻辑已实现，本次复用）。

## Goals / Non-Goals

**Goals:**
- 为 `tencentcloud_mongodb_instance` 资源新增 `cpu` 参数（`Optional` + `Computed`，`TypeInt`）。
- 在 Update 流程中，当 `cpu` 发生变化时，将 `cpu` 传递给 `ModifyDBInstanceSpec` 接口的 `Cpu` 字段完成变配。
- 在 Read 流程中，从 `DescribeDBInstances` 返回的 `InstanceDetail.CpuNum` 回填 `cpu` 到 state。
- 补充基于 gomonkey 的单元测试，覆盖 `cpu` 变更触发变配与回填逻辑。

**Non-Goals:**
- 不修改 Create 流程：`CreateDBInstanceHour`/`CreateDBInstance` 创建时使用的是 `CpuCore` 字段（非本次需求范围内的 `Cpu`），且本次需求仅针对 `ModifyDBInstanceSpec` 的 `Cpu` 入参，不扩展到创建参数。
- 不修改 `mongos_cpu` 等已有分片集群相关 CPU 参数。
- 不修改已发布字段的 schema 行为（不新增 `ForceNew`，不改变 `Required`/`Optional` 语义）。
- 不执行 `make doc`、`gofmt`、`.changelog` 等收尾操作（由 tfpacer-finalize 阶段统一处理）。

## Decisions

### 决策 1：`cpu` 字段定义为 `Optional` + `Computed`，不设置 `ForceNew`

**理由**：`cpu` 变化通过 in-place update（调用 `ModifyDBInstanceSpec`）完成，无需重建实例。设为 `Computed` 使 Read 能从云上回填实际 CPU 值，避免用户未指定时 state 与云上不一致。与已有 `memory`（`Required`）、`volume`（`Required`）不同，`cpu` 为可选参数（云 API 注释说明为空值时默认取实例当前 CPU 大小），因此设为 `Optional` + `Computed` 更贴合 API 语义且向后兼容。

**备选方案**：将 `cpu` 设为 `Required`。否决，因为会破坏向后兼容（已有配置未设置 `cpu`）且与 API 语义不符。

### 决策 2：Update 触发条件扩展为包含 `cpu`

**理由**：现有触发条件为 `d.HasChange("memory") || d.HasChange("volume") || d.HasChange("node_num")`，新增 `d.HasChange("cpu")`。当用户单独修改 `cpu` 时也应触发变配。`cpu` 通过现有 `params map[string]interface{}` 传递到 `UpgradeInstance`，复用现有 `in_maintenance`/`node_num` 的传参模式，最小改动。

**备选方案**：为 `cpu` 单独创建一个 `ModifyDBInstanceSpec` 调用分支。否决，会引入重复的订单轮询与错误处理逻辑，违背 DRY，且云 API `ModifyDBInstanceSpec` 本身支持同时变更多个规格参数。

### 决策 3：服务层 `UpgradeInstance` 新增 `params["cpu"]` 映射到 `request.Cpu`

**理由**：`ModifyDBInstanceSpecRequest.Cpu` 为 `*int64`，使用 `helper.Int64(v.(int))` 构造指针。与现有 `params["in_maintenance"]` 等使用 `helper.IntUint64` 不同，需按字段实际类型 `int64` 使用 `helper.Int64`，避免类型不匹配。

### 决策 4：Read 回填使用 `instance.CpuNum`（`*uint64`）转为 `int`

**理由**：`DescribeDBInstances` 返回的 `InstanceDetail.CpuNum` 即实例 CPU 核数。在 Read 中先判断 `instance.CpuNum != nil` 再 `d.Set("cpu", int(*instance.CpuNum))`，遵循项目规范「调用 setXX 前先判断字段是否为 nil」。不将 `CpuNum` 加入 `CheckNil` 强校验列表，避免对未返回该字段的存量实例造成 Read 失败。

## Risks / Trade-offs

- **[风险] `CpuNum` 为 `*uint64`，`Cpu` 为 `*int64`，类型转换溢出** → 实际 CPU 核数为正小整数，不会溢出；转换使用 `int()` 显式转换，安全。
- **[风险] 用户单独修改 `cpu` 而不修改 `memory`/`volume` 时，云 API 是否允许** → 云 API 注释明确支持单独变更 CPU（"变更 mongod 或 mongos 的 CPU 与内存规格时，NodeNum 可以不配置"），且 `Cpu` 注释说明为空时取当前值，语义允许单独变更。具体可用规格由 `DescribeSpecInfo` 约束，由云 API 校验，Provider 不做额外校验。
- **[权衡] 未将 `CpuNum` 加入 `CheckNil` 强校验** → 避免对未返回该字段的存量实例造成 Read 失败；回填前单独判 nil，安全且向后兼容。
- **[风险] 既有 `UpgradeInstance` 复用同一订单轮询逻辑** → 不新增轮询逻辑，沿用现有 `DescribeDBInstanceDeal` 轮询，行为一致。
