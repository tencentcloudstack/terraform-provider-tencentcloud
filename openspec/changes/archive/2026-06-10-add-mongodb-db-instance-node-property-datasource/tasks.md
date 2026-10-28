## 1. 数据源实现

- [x] 1.1 创建 `tencentcloud/services/mongodb/data_source_tc_mongodb_db_instance_node_property.go`，实现数据源 schema 定义（入参：`instance_id`、`node_ids`、`roles`、`only_hidden`、`priority`、`votes`、`tags`；出参：`mongos`、`replicate_sets`）
- [x] 1.2 实现数据源 Read 函数，调用云 API `DescribeDBInstanceNodeProperty`，使用 `tccommon.ReadRetryTimeout` 超时和 retry 重试逻辑
- [x] 1.3 实现 `mongos` 字段的映射逻辑，将 `NodeProperty` 结构体字段（`zone`、`node_name`、`address`、`wan_service_address`、`role`、`hidden`、`status`、`slave_delay`、`priority`、`votes`、`tags`、`replicate_set_id`）映射到 schema
- [x] 1.4 实现 `replicate_sets` 字段的映射逻辑，将 `ReplicateSetInfo.Nodes` 中的 `NodeProperty` 列表映射到 schema

## 2. Provider 注册

- [x] 2.1 在 `tencentcloud/provider.go` 中注册数据源 `tencentcloud_mongodb_db_instance_node_property`
- [x] 2.2 在 `tencentcloud/provider.md` 中添加数据源引用

## 3. 文档

- [x] 3.1 创建 `tencentcloud/services/mongodb/data_source_tc_mongodb_db_instance_node_property.md`，包含一句话描述和 Example Usage

## 4. 测试

- [x] 4.1 创建 `tencentcloud/services/mongodb/data_source_tc_mongodb_db_instance_node_property_test.go`，使用 gomonkey mock 云 API 实现单元测试，覆盖正常查询、空结果、错误处理等场景
- [x] 4.2 使用 `go test -gcflags=all=-l` 运行单元测试，确保测试通过
