## Why

腾讯云 TDSQL-C MySQL（tdmysql，TcaplusDB cluster / Cyberspace MySQL）目前未在 Terraform Provider 中提供实例资源的声明式管理能力。用户只能通过控制台或 API 手动创建、隔离实例，无法将其纳入基础设施即代码（IaC）流水线，导致环境不可复现、配置易漂移。

腾讯云 tdmysql（v20211122）已提供 `CreateDBInstances`（批量创建实例）、`DescribeDBInstanceDetail`（查询实例详情）、`ModifyInstanceName`（修改实例名称）、`IsolateDBInstance`（批量隔离实例）四个接口，覆盖了一个实例资源从创建、查询、改名到销毁（隔离）的完整生命周期。通过新增 `tencentcloud_tdmysql_db_instance` 资源，用户可以以 Terraform 声明式方式管理 tdmysql 实例。

## What Changes

新增 Terraform RESOURCE_KIND_GENERAL 资源 `tencentcloud_tdmysql_db_instance`，支持完整的 CRUD 操作：

- **Create**: 调用 `CreateDBInstances` 批量创建实例。创建为异步流程，接口返回 `InstanceIds`（创建出的实例 ID 列表）与 `FlowId`；资源以返回的第一个实例 ID 作为 Terraform state id。创建完成后通过 `DescribeDBInstanceDetail` 轮询实例详情直到实例可读，保证后续 Read/状态收敛一致。
- **Read**: 调用 `DescribeDBInstanceDetail` 查询实例详情，将云端实际状态回填到 Terraform state。若实例不存在则清空 id。
- **Update**: 当 `instance_name` 发生变化时，调用 `ModifyInstanceName` 修改实例名称。除 `instance_name` 外的顶层创建参数不可变（Immutable，修改将报错）。
- **Delete**: 调用 `IsolateDBInstance` 隔离实例（tdmysql 的销毁语义为隔离）。

### 新增文件
- `tencentcloud/services/tdmysql/resource_tc_tdmysql_db_instance.go` - 资源实现
- `tencentcloud/services/tdmysql/resource_tc_tdmysql_db_instance_test.go` - 单元测试（使用 gomonkey mock 云 API）
- `tencentcloud/services/tdmysql/resource_tc_tdmysql_db_instance.md` - 资源文档
- `tencentcloud/services/tdmysql/service_tencentcloud_tdmysql.go` - 服务层方法（查询实例详情、隔离实例、查询流程状态）

### 修改文件
- `tencentcloud/provider.go` - 注册新资源 `tencentcloud_tdmysql_db_instance`，导入 tdmysql 服务包
- `tencentcloud/provider.md` - 在资源列表中声明 `tencentcloud_tdmysql_db_instance`
- `tencentcloud/connectivity/client.go` - 新增 `UseTdmysqlV20211122Client()` 方法及 `tdmysqlv20211122Conn` 字段、对应 import

### 资源 Schema
```hcl
resource "tencentcloud_tdmysql_db_instance" "example" {
  zone              = "ap-guangzhou-3"
  vpc_id            = "vpc-xxxxxxxx"
  subnet_id         = "subnet-xxxxxxxx"
  spec_code         = "spec-code"
  disk              = 100
  storage_node_num  = 2
  replications      = 3
  instance_count    = 1
  instance_name     = "tf-tdmysql-example"
  pay_mode          = "0"
  instance_type     = "separate"
  storage_type      = "CLOUD_HSSD"
  # ... 其他可选参数
}
```

### 资源 ID 格式
使用单个实例 ID 作为 Terraform state id（`CreateDBInstances` 返回 `InstanceIds` 列表，取第一个）。资源支持 import（RESOURCE_KIND_GENERAL），import 时使用实例 ID。

## Capabilities

### New Capabilities
- `tdmysql-db-instance`: 管理 tdmysql（TDSQL-C MySQL）实例资源的完整生命周期（创建、查询、改名、隔离），对应 `tencentcloud_tdmysql_db_instance` 资源。

### Modified Capabilities
（无，本次为全新资源，不修改既有规范。）

## Impact

### 受影响的代码
- `tencentcloud/services/tdmysql/` - 新增服务目录及全部资源实现、服务层、测试、文档
- `tencentcloud/connectivity/client.go` - 新增 tdmysql SDK 客户端访问方法
- `tencentcloud/provider.go` - 资源注册与包导入
- `tencentcloud/provider.md` - 资源名声明（用于 `make doc` 文档生成）

### 受影响的 API
- `CreateDBInstances`（tdmysql v20211122）- 批量创建实例（异步，返回 FlowId）
- `DescribeDBInstanceDetail`（tdmysql v20211122）- 查询实例详情
- `DescribeFlow`（tdmysql v20211122）- 查询异步任务流程状态（用于创建后轮询）
- `ModifyInstanceName`（tdmysql v20211122）- 修改实例名称
- `IsolateDBInstance`（tdmysql v20211122）- 批量隔离实例

### 向后兼容性
- ✅ 完全向后兼容，新增资源，不修改任何既有资源/数据源/接口。

### 依赖关系
- 依赖已 vendored 的 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmysql/v20211122` 包（已在 vendor 目录中）。
