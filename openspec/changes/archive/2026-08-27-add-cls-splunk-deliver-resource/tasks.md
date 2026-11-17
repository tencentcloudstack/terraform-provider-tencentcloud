## 1. 资源 Schema 定义

- [x] 1.1 在 `tencentcloud/services/cls/resource_tc_cls_splunk_deliver.go` 中定义 `ResourceTencentCloudClsSplunkDeliver` 函数，创建 schema 包含所有顶层字段（`topic_id`、`name`、`enable`、`has_service_log`、`index_ack`、`source`、`source_type`、`index`、`channel`、`dsl_filter`、`task_id`）
- [x] 1.2 定义 `net_info` 嵌套 TypeList（MaxItems:1）schema，包含 `host`、`port`、`token`、`net_type`、`vpc_id`、`virtual_gateway_type`、`is_ssl`
- [x] 1.3 定义 `metadata_info` 嵌套 TypeList（MaxItems:1）schema，包含 `format`、`meta_fields`（TypeSet）、`enable_tag`、`tag_json_tiled`
- [x] 1.4 定义 `external_role` 嵌套 TypeList（MaxItems:1）schema，包含 `role_arn`、`external_id`

## 2. Create 操作实现

- [x] 2.1 实现 `resourceTencentCloudClsSplunkDeliverCreate` 函数：从 schema 中读取参数，构建 `CreateSplunkDeliverRequest`，调用 API
- [x] 2.2 在 Create 中使用 `tccommon.ReadRetryTimeout` 和 `tccommon.RetryError` 进行 retry 处理
- [x] 2.3 检查 Create 返回值是否为空/nil，若 `TaskId` 为空则返回 `NonRetryableError`
- [x] 2.4 使用 `task_id#topic_id` 格式设置资源 ID，然后调用 Read 刷新状态

## 3. Read 操作实现

- [x] 3.1 实现 `resourceTencentCloudClsSplunkDeliverRead` 函数：从 ID 解析 `task_id` 和 `topic_id`，构建 `DescribeSplunkDeliversRequest`
- [x] 3.2 使用 `Filters` 按 `taskId` 过滤，设置 `Limit` 为最大值 100
- [x] 3.3 检查响应是否为空，若 `Infos` 为空则记录日志后 `d.SetId("")`
- [x] 3.4 在 set 字段前检查 Response 中每个字段是否为 nil，nil 则不设置

## 4. Update 操作实现

- [x] 4.1 实现 `resourceTencentCloudClsSplunkDeliverUpdate` 函数：从 ID 解析参数，构建 `ModifySplunkDeliverRequest`
- [x] 4.2 使用 `d.HasChange()` 检测变更字段，仅将变更的字段设置到请求中
- [x] 4.3 调用 API 后调用 Read 刷新状态

## 5. Delete 操作实现

- [x] 5.1 实现 `resourceTencentCloudClsSplunkDeliverDelete` 函数：从 ID 解析 `task_id` 和 `topic_id`，构建 `DeleteSplunkDeliverRequest`
- [x] 5.2 使用 `tccommon.ReadRetryTimeout` 和 `tccommon.RetryError` 进行 retry 处理

## 6. Import 支持

- [x] 6.1 在 schema 中添加 `Importer` 配置，支持 `task_id#topic_id` 格式导入

## 7. Provider 注册

- [x] 7.1 在 `tencentcloud/provider.go` 中注册 `tencentcloud_cls_splunk_deliver` 资源
- [x] 7.2 在 `tencentcloud/provider.md` 中添加资源文档入口

## 8. 单元测试

- [x] 8.1 创建 `tencentcloud/services/cls/resource_tc_cls_splunk_deliver_test.go`，使用 gomonkey mock 云 API 进行单元测试

## 9. 文档

- [x] 9.1 创建 `tencentcloud/services/cls/resource_tc_cls_splunk_deliver.md`，包含资源描述、Example Usage 和 Import 示例