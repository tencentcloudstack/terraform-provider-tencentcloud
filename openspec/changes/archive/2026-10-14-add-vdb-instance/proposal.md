## Why

Terraform Provider for TencentCloud 目前不支持管理 VectorDB (VDB) 实例生命周期。用户无法通过 Terraform 创建、查询、升配/缩容、销毁 VDB 实例。为补齐云产品 `vdb` 的 Terraform 资源覆盖，需新增 `tencentcloud_vdb_instance` 通用资源 (RESOURCE_KIND_GENERAL)，支持实例的完整 CRUD 操作，满足 IaC 管理向量数据库实例的需求。

## What Changes

- 新增资源 `tencentcloud_vdb_instance`，管理 VDB 实例的创建、读取、更新、删除全生命周期。
- Create: 调用 `CreateInstance` API 创建实例，支持 VPC、子网、计费模式、规格、安全组、标签、组件(components)等参数；创建为异步操作，创建后轮询 Read 接口等待实例变为 `online` 状态。
- Read: 调用 `DescribeInstances` API 查询实例详情并回填 state，包括网络信息、规格、状态等；同时调用 `DescribeInstanceNodes` 查询节点信息。
- Update: 支持通过 `ScaleUpInstance`（cpu/memory/storage_size 垂直扩容）和 `ScaleOutInstance`（replica_num 水平扩容）更新规格；通过 `ModifyDBInstanceSecurityGroups` 更新安全组。不可变字段变更时拒绝更新。
- Delete: 调用 `IsolateInstance` 隔离实例；当 `force_delete` 为 true 时，继续调用 `DestroyInstances` 彻底销毁实例。
- 在 `tencentcloud/provider.go` 和 `tencentcloud/provider.md` 中注册新资源。
- 新增资源文档 `resource_tc_vdb_instance.md`（参考 `resource_tc_config_compliance_pack.md` 命名格式）。
- 新增单元测试 `resource_tc_vdb_instance_test.go`，使用 mock（gomonkey）方式对云 API 进行 mock，仅测试业务代码逻辑（参考 `resource_tc_config_compliance_pack_test.go` 命名格式）。

## Capabilities

### New Capabilities
- `vdb-instance`: 管理 VectorDB 实例的全生命周期，包括创建、查询、规格升降配（ScaleUp/ScaleOut）、安全组更新、隔离与销毁。

### Modified Capabilities
<!-- 无现有 capability 需要修改 -->

## Impact

- 新增文件:
  - `tencentcloud/services/vdb/resource_tc_vdb_instance.go` — 资源 CRUD 实现
  - `tencentcloud/services/vdb/resource_tc_vdb_instance.md` — 资源文档
  - `tencentcloud/services/vdb/resource_tc_vdb_instance_test.go` — 单元测试（mock 方式）
  - `tencentcloud/services/vdb/service_tencentcloud_vdb.go` — 服务层（封装 SDK 调用与轮询逻辑）
- 修改文件:
  - `tencentcloud/provider.go` — 注册 `tencentcloud_vdb_instance` 资源
  - `tencentcloud/provider.md` — 文档侧边栏新增资源条目
- 涉及云 API（包名 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vdb/v20230616`）:
  - `CreateInstance` — 创建实例（异步）
  - `DescribeInstances` — 查询实例列表
  - `ScaleUpInstance` — 垂直扩容
  - `ScaleOutInstance` — 水平扩容
  - `ModifyDBInstanceSecurityGroups` — 修改安全组
  - `IsolateInstance` — 隔离实例
  - `DestroyInstances` — 销毁实例
  - `DescribeInstanceNodes` — 查询实例节点
- 代码风格严格参考 `tencentcloud_igtm_monitor` 资源。