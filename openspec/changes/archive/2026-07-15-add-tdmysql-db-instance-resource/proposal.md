## Why

TDSQL-C for MySQL（tdmysql）是腾讯云提供的云原生数据库产品。当前 Terraform Provider 中尚未支持通过 Terraform 管理 tdmysql 实例的完整生命周期，用户无法以基础设施即代码的方式创建、查询、修改和销毁 tdmysql 实例，只能依赖控制台或 API 手动操作。需要新增 `tencentcloud_tdmysql_db_instance` 通用资源，支持 tdmysql 实例的 CRUD 全生命周期管理。

## What Changes

- 新增 Terraform 通用资源 `tencentcloud_tdmysql_db_instance`，支持 tdmysql 实例的 CRUD 操作
  - **Create**: 调用 `CreateDBInstances` 接口批量创建实例（入参含 zone、vpc_id、subnet_id、spec_code、disk、storage_node_num、replications、instance_count、init_params 等），返回 `InstanceIds`（实例 ID 列表）和 `FlowId`（异步任务流程 ID）。创建后通过 `DescribeFlow` 轮询流程状态直至实例创建完成，再调用 `DescribeDBInstanceDetail` 回读实例详情
  - **Read**: 调用 `DescribeDBInstanceDetail` 接口查询实例详情，回写所有可读字段（含 instance_name、zone、vpc_id、subnet_id、vip、vport、status、disk、spec_code 等）
  - **Update**: 调用 `ModifyInstanceName` 接口修改实例名称（instance_name）；其余创建时传入的字段均为不可变字段（immutable），变更时返回 error
  - **Delete**: 调用 `IsolateDBInstance` 接口批量隔离实例
- 资源 ID 使用单个 `instance_id`（取自 `CreateDBInstances` 返回的 `InstanceIds` 列表中第一个元素）
- 在 `tencentcloud/connectivity/client.go` 中新增 `UseTdmysqlV20211122Client` 客户端辅助方法及对应的连接字段与 import
- 在 `tencentcloud/provider.go` 和 `tencentcloud/provider.md` 中注册新资源
- 生成对应的 `.md` 文档文件

## Capabilities

### New Capabilities
- `tdmysql-db-instance-resource`: 新增 TDSQL-C for MySQL 实例通用资源，支持实例的创建（含异步流程轮询）、查询、修改实例名称、隔离（删除）等完整生命周期管理

### Modified Capabilities
<!-- 无现有 capability 的需求变更 -->

## Impact

- 新增文件：`tencentcloud/services/tdmysql/resource_tc_tdmysql_db_instance.go`
- 新增文件：`tencentcloud/services/tdmysql/resource_tc_tdmysql_db_instance_test.go`
- 新增文件：`tencentcloud/services/tdmysql/resource_tc_tdmysql_db_instance.md`
- 修改文件：`tencentcloud/connectivity/client.go`（新增 tdmysql SDK import、`tdmysqlv20211122Conn` 字段、`UseTdmysqlV20211122Client()` 方法）
- 修改文件：`tencentcloud/provider.go`（注册新资源 `tencentcloud_tdmysql_db_instance`）
- 修改文件：`tencentcloud/provider.md`（添加资源文档索引）
- 依赖的云 API 接口：`CreateDBInstances`、`DescribeDBInstanceDetail`、`DescribeFlow`、`ModifyInstanceName`、`IsolateDBInstance`
- 依赖的 SDK 包：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmysql/v20211122`
