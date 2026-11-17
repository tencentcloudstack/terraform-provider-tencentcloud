## Why

`tencentcloud_mongodb_instance` 资源当前支持通过 `memory`、`volume`、`node_num` 等参数调整实例配置（调用 `ModifyDBInstanceSpec` 接口），但缺少对 CPU 核数的管理。腾讯云 MongoDB `ModifyDBInstanceSpec` 接口已支持 `Cpu` 入参，允许在变更实例规格时指定变更后的 CPU 大小（单位：C）。本次变更新增 `cpu` 参数，使用户能够通过 Terraform 管理 MongoDB 实例的 CPU 规格，补全实例配置变更能力。

## What Changes

- 在 `tencentcloud_mongodb_instance` 资源 schema 中新增 `cpu` 参数（`Optional` + `Computed`，`TypeInt`），用于指定实例变更后的 CPU 核数。
- 在资源 Update 流程中，将 `cpu` 加入 `ModifyDBInstanceSpec` 的触发条件，当 `cpu` 发生变化时（与 `memory`、`volume`、`node_num` 等一起或单独）调用 `ModifyDBInstanceSpec` 接口完成变配。
- 在 `MongodbService.UpgradeInstance` 服务方法中，将 `params["cpu"]` 映射到 `ModifyDBInstanceSpecRequest.Cpu` 字段（`*int64`）。
- 在资源 Read 流程中，通过 `DescribeDBInstances` 返回的 `InstanceDetail.CpuNum` 字段回填 `cpu` 参数到 state。
- 补充 `resource_tc_mongodb_instance_test.go` 的单元测试用例（使用 gomonkey mock 云 API）。
- 更新 `resource_tc_mongodb_instance.md` 资源文档。

## Capabilities

### New Capabilities
- `mongodb-instance-cpu`: 为 `tencentcloud_mongodb_instance` 资源新增 `cpu` 参数，支持在实例配置变更时设置 CPU 核数，并通过查询接口回填 CPU 状态。

### Modified Capabilities
<!-- 无需修改现有 spec 级别的需求 -->

## Impact

### 受影响的代码
- `tencentcloud/services/mongodb/resource_tc_mongodb_instance.go` - schema 新增 `cpu` 字段、Read 回填、Update 触发与传参
- `tencentcloud/services/mongodb/service_tencentcloud_mongodb.go` - `UpgradeInstance` 方法新增 `cpu` 参数映射到 `ModifyDBInstanceSpecRequest.Cpu`
- `tencentcloud/services/mongodb/resource_tc_mongodb_instance_test.go` - 新增/补充单元测试
- `tencentcloud/services/mongodb/resource_tc_mongodb_instance.md` - 文档更新

### 受影响的 API
- `ModifyDBInstanceSpec`（变配时设置 `Cpu` 入参，`*int64`）
- `DescribeDBInstances`（读取时从 `InstanceDetail.CpuNum` 回填，`*uint64`）

### 向后兼容性
- ✅ 完全向后兼容：仅新增 `Optional` + `Computed` 字段，不修改已有 schema，不影响现有 TF 配置和 state。
- ✅ 不设置 `ForceNew`，`cpu` 变化通过 in-place update 处理。

### 依赖关系
- 无新增依赖，复用现有 `ModifyDBInstanceSpec` 与 `DescribeDBInstances` 接口及现有 `UpgradeInstance` 服务方法。
