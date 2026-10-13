## Why

MongoDB 实例目前无法通过 Terraform 手动开启 CPU 弹性扩容能力。腾讯云 MongoDB 提供了 `ScaleUpDBInstanceCpu`（手动开启 CPU 弹性扩容）与 `ScaleDownDBInstanceCpu`（关闭弹性 CPU 扩容并回缩）两个接口，以及 `CheckDBInstanceElasticCpuScalable` 查询接口，但 provider 中尚未封装这些能力。运维人员需要以声明式方式管理 MongoDB 实例的弹性 CPU 扩容状态，以便应对业务高峰、并在低峰时回缩，从而兼顾性能与成本。

## What Changes

- 新增 Terraform 通用资源 `tencentcloud_mongodb_scale_up_db_instance_cpu`（RESOURCE_KIND_GENERAL），支持完整 CRUD：
  - Create：调用 `ScaleUpDBInstanceCpu` 手动开启 CPU 弹性扩容；接口为异步（返回 `FlowId`），需轮询 `CheckDBInstanceElasticCpuScalable` 直到 `IsScaled` 为 true 生效；拿到资源 id 后、异步轮询之前先 `d.SetId()`。
  - Read：调用 `CheckDBInstanceElasticCpuScalable` 查询当前弹性 CPU 扩容信息；若返回空则打印现场日志后再 `d.SetId("")`。
  - Update：`extra_cpu` 变化时调用 `ScaleUpDBInstanceCpu` 重新触发扩容（实例 id 为 ForceNew，不可变）。
  - Delete：调用 `ScaleDownDBInstanceCpu` 关闭弹性 CPU 扩容并回缩；接口为异步（返回 `FlowId`），轮询直到 `IsScaled` 为 false。
- 资源参数：
  - 输入：`instance_id`（可选，ForceNew）、`extra_cpu`（可选）
  - 输出（Computed）：`flow_id`、`scalable`、`is_scaled`、`is_locked`、`reason`、`max_extra_cpu`、`trigger_type`、`scale_up_time`
- 资源 id：`{instanceId}`（用实例 id 唯一标识一个实例的弹性 CPU 扩容状态）
- 新增文件：
  - `tencentcloud/services/mongodb/resource_tc_mongodb_scale_up_db_instance_cpu.go`
  - `tencentcloud/services/mongodb/resource_tc_mongodb_scale_up_db_instance_cpu_test.go`（gomonkey 单测）
  - `tencentcloud/services/mongodb/resource_tc_mongodb_scale_up_db_instance_cpu.md`
- 修改文件：
  - `tencentcloud/provider.go`：注册 `tencentcloud_mongodb_scale_up_db_instance_cpu`
  - `tencentcloud/services/mongodb/service_tencentcloud_mongodb.go`：新增服务层方法（ScaleUp/ScaleDown/Check 三接口的封装 + 轮询等待逻辑）
- 该资源支持 Import，import 时使用实例 id。

## Capabilities

### New Capabilities
- `mongodb-scale-up-db-instance-cpu`: 管理 MongoDB 实例的弹性 CPU 扩容能力，包括手动开启扩容（Create/Update）、查询扩容状态（Read）、关闭扩容并回缩（Delete）

### Modified Capabilities
<!-- 无 -->

## Impact

- **代码**：新增 mongodb 服务下的资源实现、测试与示例；provider.go 注册新资源；service 层扩展三接口封装。
- **API**：使用 `ScaleUpDBInstanceCpu`、`ScaleDownDBInstanceCpu`、`CheckDBInstanceElasticCpuScalable` 三个 mongodb v20190725 接口，均已存在于 vendor 中。
- **依赖**：无新增外部依赖，复用现有 SDK 与 helper/tccommon 包。
- **向后兼容**：纯新增资源，不影响现有资源与 state。