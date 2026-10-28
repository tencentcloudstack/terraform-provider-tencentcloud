## Why

用户需要查询 TencentCloud PostgreSQL（postgres）实例或只读组关联的安全组信息，用于安全组审计、网络配置核对及自动化编排。目前 Provider 缺少对应的数据源支持，用户只能通过控制台或 API 手动查询，无法在 Terraform 配置中直接引用安全组数据。

## What Changes

- 新增数据源 `tencentcloud_postgresql_db_instance_security_groups`，调用 postgres 的 `DescribeDBInstanceSecurityGroups` 接口查询实例/只读组关联的安全组
- 支持通过 `db_instance_id` 查询实例关联的安全组
- 支持通过 `read_only_group_id` 查询只读组关联的安全组
- 返回 `security_group_set` 列表，包含安全组 ID、名称、备注、项目 ID、创建时间及入站/出站规则

## Capabilities

### New Capabilities
- `postgresql-db-instance-security-groups-datasource`: 查询 PostgreSQL 实例/只读组关联安全组的数据源，支持按 `db_instance_id` 或 `read_only_group_id` 查询

### Modified Capabilities
<!-- 无现有功能需要修改 -->

## Impact

**新增文件:**
- `tencentcloud/services/postgresql/data_source_tc_postgresql_db_instance_security_groups.go` - 数据源实现
- `tencentcloud/services/postgresql/data_source_tc_postgresql_db_instance_security_groups_test.go` - 测试文件
- `tencentcloud/services/postgresql/data_source_tc_postgresql_db_instance_security_groups.md` - 文档模板

**修改文件:**
- `tencentcloud/provider.go` - 在 DataSourcesMap 中注册新数据源
- `tencentcloud/provider.md` - 添加数据源到文档列表

**依赖:**
- 已有的 tencentcloud-sdk-go postgres/v20170312 包（`DescribeDBInstanceSecurityGroups` 接口）
- 无需新增外部依赖
