## 1. Schema 定义与 CRUD 函数

- [x] 1.1 创建资源文件 `tencentcloud/services/sqlserver/resource_tc_sqlserver_db_instance_ssl_config.go`，定义 Schema（instance_id Required+ForceNew, type Required, wait_switch/is_kms/key_id/key_region Optional, encryption/ssl_validity_period/ssl_validity Computed），实现 Create/Read/Update/Delete 函数，支持 Import
- [x] 1.2 实现 Create 函数：设置 ID 为 instance_id，然后调用 Update 函数
- [x] 1.3 实现 Read 函数：调用 `SqlserverService.DescribeSqlserverInstanceSslById` 读取 SSLConfig，设置 encryption、ssl_validity_period、ssl_validity 等 Computed 字段
- [x] 1.4 实现 Update 函数：调用 ModifyDBInstanceSSL（带 retry），获取 FlowId，轮询 DescribeFlowStatus 直到 Status == 0（SQLSERVER_TASK_SUCCESS），然后调用 Read
- [x] 1.5 实现 Delete 函数：空操作（CONFIG 类型资源）

## 2. Provider 注册

- [x] 2.1 在 `tencentcloud/provider.go` 中注册 `tencentcloud_sqlserver_db_instance_ssl_config` 资源
- [x] 2.2 在 `tencentcloud/provider.md` 中添加资源描述

## 3. 资源文档

- [x] 3.1 创建 `tencentcloud/services/sqlserver/resource_tc_sqlserver_db_instance_ssl_config.md` 文件，包含一句话描述、Example Usage、Import 说明

## 4. 单元测试

- [x] 4.1 创建 `tencentcloud/services/sqlserver/resource_tc_sqlserver_db_instance_ssl_config_test.go`，使用 gomonkey mock 方式编写 Create/Read/Update 单元测试

## 5. 验证

- [x] 5.1 使用 `go test -gcflags=all=-l` 运行单元测试确保通过
