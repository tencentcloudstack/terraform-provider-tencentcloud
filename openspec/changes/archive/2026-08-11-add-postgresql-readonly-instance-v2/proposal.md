## Why

Terraform Provider for TencentCloud 当前缺少对 PostgreSQL 只读实例（tencentcloud_postgresql_readonly_instance_v2）的管理能力。用户需要在 Terraform 中创建、查询、隔离（销毁）PostgreSQL 只读实例，以支持读扩展和灾备场景。现在补齐该资源可让用户以基础设施即代码的方式管理只读实例的完整生命周期。

## What Changes

- 新增资源 `tencentcloud_postgresql_readonly_instance_v2`（RESOURCE_KIND_GENERAL），通过云 API `CreateReadOnlyDBInstance`、`DescribeDBInstanceAttribute`、`IsolateDBInstances` 实现完整的 CRUD 生命周期管理（增/查/改/删）。
- 新增资源文件 `tencentcloud/services/postgresql/resource_tc_postgresql_readonly_instance_v2.go`。
- 新增服务层逻辑文件 `tencentcloud/services/postgresql/service_tencentcloud_postgresql.go`（如已存在则在其中追加方法）。
- 新增单元测试文件 `tencentcloud/services/postgresql/resource_tc_postgresql_readonly_instance_v2_test.go`，使用 gomonkey mock 云 API 进行业务逻辑测试。
- 新增文档 `tencentcloud/services/postgresql/resource_tc_postgresql_readonly_instance_v2.md`（最终通过 `make doc` 生成 website 文档）。
- 在 `tencentcloud/provider.go` 与 `tencentcloud/provider.md` 中注册该资源。

## Capabilities

### New Capabilities
- `postgresql-readonly-instance-v2-resource`: 管理 PostgreSQL 只读实例（tencentcloud_postgresql_readonly_instance_v2）的创建、读取、更新、销毁生命周期。

### Modified Capabilities
<!-- 无既有能力的需求变更 -->

## Impact

- 新增资源代码：`tencentcloud/services/postgresql/resource_tc_postgresql_readonly_instance_v2.go`
- 新增/扩展服务层：`tencentcloud/services/postgresql/service_tencentcloud_postgresql.go`（封装 `CreateReadOnlyDBInstance`、`DescribeDBInstanceAttribute`、`IsolateDBInstances` 调用）
- 新增单元测试：`tencentcloud/services/postgresql/resource_tc_postgresql_readonly_instance_v2_test.go`
- 新增文档：`tencentcloud/services/postgresql/resource_tc_postgresql_readonly_instance_v2.md`
- 注册变更：`tencentcloud/provider.go`、`tencentcloud/provider.md`
- 依赖云 API（已在 vendor 中）：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312`
