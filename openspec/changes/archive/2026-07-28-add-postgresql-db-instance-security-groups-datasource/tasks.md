## 1. 服务层实现

- [x] 1.1 在 `tencentcloud/services/postgresql/service_tencentcloud_postgresql.go` 中新增 `DescribePostgresqlDbInstanceSecurityGroups` 方法，封装 `DescribeDBInstanceSecurityGroups` API 调用，接收 `dbInstanceId`、`readOnlyGroupId` 参数
- [x] 1.2 在服务方法中构建 `DescribeDBInstanceSecurityGroupsRequest`，设置 `DBInstanceId`、`ReadOnlyGroupId`，调用 client 并返回 `[]*postgresql.SecurityGroup`

## 2. 数据源 Schema 与 Read 实现

- [x] 2.1 创建 `tencentcloud/services/postgresql/data_source_tc_postgresql_db_instance_security_groups.go` 文件
- [x] 2.2 定义数据源 Schema，包含输入参数 `db_instance_id`(Optional)、`read_only_group_id`(Optional)、`result_output_file`(Optional)
- [x] 2.3 定义输出参数 `security_group_set`(Computed, TypeList)，元素为 `schema.Resource`，平铺字段：`project_id`、`create_time`、`security_group_id`、`security_group_name`、`security_group_description`、`inbound`、`outbound`
- [x] 2.4 定义 `inbound`/`outbound` 的嵌套 `PolicyRule` schema，字段：`action`、`cidr_ip`、`port_range`、`ip_protocol`、`description`
- [x] 2.5 实现 `DataSourceTencentCloudPostgresqlDbInstanceSecurityGroups()` 函数返回 `*schema.Resource`
- [x] 2.6 实现 `dataSourceTencentCloudPostgresqlDbInstanceSecurityGroupsRead()` 函数，包含 `defer tccommon.LogElapsed()` 和 `defer tccommon.InconsistentCheck()`
- [x] 2.7 在 Read 函数中使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包装服务层调用，失败时用 `tccommon.RetryError()` 包装错误
- [x] 2.8 在 retry 块内检查 API 返回是否为空（`response == nil`/`response.Response == nil`/`len(SecurityGroupSet) == 0`），为空时返回 `NonRetryableError`，不直接 `d.SetId("")`
- [x] 2.9 在外层 retry 失败路径保留 `log.Printf("[DATASOURCE] read empty, skip SetId")` 提示
- [x] 2.10 实现 SecurityGroup 列表展开逻辑，遍历 `SecurityGroupSet`，在 set 每个字段前判断是否为 nil，为 nil 则跳过 set
- [x] 2.11 实现 `inbound`/`outbound` PolicyRule 列表的展开与 set 逻辑（同样做 nil 检查）
- [x] 2.12 使用 `helper.DataResourceIdsHash(ids)` 生成数据源 ID
- [x] 2.13 实现 `result_output_file` 功能，将结果以 JSON 格式保存到指定文件

## 3. Provider 注册

- [x] 3.1 在 `tencentcloud/provider.go` 的 `DataSourcesMap` 中注册 `tencentcloud_postgresql_db_instance_security_groups`，指向 `DataSourceTencentCloudPostgresqlDbInstanceSecurityGroups`
- [x] 3.2 在 `tencentcloud/provider.md` 的数据源列表中按字母序添加 `tencentcloud_postgresql_db_instance_security_groups`

## 4. 文档模板创建

- [x] 4.1 创建 `tencentcloud/services/postgresql/data_source_tc_postgresql_db_instance_security_groups.md` 文档模板
- [x] 4.2 添加一句话描述（RESOURCE_KIND_DATASOURCE 格式："Use this data source to query ..."，带上 postgres 产品名称）
- [x] 4.3 添加 Example Usage 示例（按 `db_instance_id` 查询、按 `read_only_group_id` 查询）
- [x] 4.4 不添加 `Argument Reference` 和 `Attribute Reference` 部分（由工具自动生成）

## 5. 单元测试创建

- [x] 5.1 创建 `tencentcloud/services/postgresql/data_source_tc_postgresql_db_instance_security_groups_test.go` 文件
- [x] 5.2 使用 gomonkey mock `PostgresqlService.DescribePostgresqlDbInstanceSecurityGroups` 方法
- [x] 5.3 实现按 `db_instance_id` 查询的单元测试用例，验证返回的 `security_group_set` 结构正确
- [x] 5.4 实现按 `read_only_group_id` 查询的单元测试用例
- [x] 5.5 实现包含 `inbound`/`outbound` 规则的测试用例，验证嵌套结构正确
- [x] 5.6 实现 `result_output_file` 功能测试用例
- [x] 5.7 确保测试代码在当前环境下可正确编译（不执行 go test）

## 6. 代码正确性检查

- [x] 6.1 检查新增参数在云 API 接口中的映射正确性：`db_instance_id` → `request.DBInstanceId`、`read_only_group_id` → `request.ReadOnlyGroupId`、`security_group_set` → `response.Response.SecurityGroupSet`
- [x] 6.2 检查所有函数返回的 error 是否被正确处理，不会出错函数用 `_ = func()` 赋值给 `_`
- [x] 6.3 检查 retry 块内仅包含 API 调用，set id 等成功操作放在 retry 块外
- [x] 6.4 检查使用资源名称（postgresql_db_instance_security_groups）而非模糊措辞进行日志打印
