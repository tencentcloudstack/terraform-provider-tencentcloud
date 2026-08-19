## 1. Service Layer

- [x] 1.1 在 `tencentcloud/services/cls/service_tencentcloud_cls.go` 中新增 `DescribeClsMetricSubscribeById` 方法，封装 `DescribeMetricSubscribes` API 调用（带 `tccommon.ReadRetryTimeout` 重试），通过 TopicId + Filter(taskId) 查询单条 `MetricSubscribeInfo`

## 2. Resource Schema & CRUD

- [x] 2.1 创建 `tencentcloud/services/cls/resource_tc_cls_metric_subscribe.go`，定义 `ResourceTencentCloudClsMetricSubscribe()` schema，包含顶层字段 `name`、`topic_id`(ForceNew)、`namespace`、`metrics`(嵌套)、`instance_info`(嵌套)、`enable`、以及 computed 字段 `task_id`、`status`、`create_time`、`update_time`
- [x] 2.2 定义 `metrics` 嵌套 schema（`metric_name`、`periods`、`metric_labels` 嵌套含 `key`/`value`）
- [x] 2.3 定义 `instance_info` 嵌套 schema（`instance_dimension`、`instances` 嵌套含 `values`），MaxItems=1
- [x] 2.4 实现 `resourceTencentCloudClsMetricSubscribeCreate`：构建 `CreateMetricSubscribeRequest`，设置 Name/TopicId/Namespace/Metrics/InstanceInfo，调用 API（WriteRetryTimeout 重试），检查返回 TaskId 非空，设置复合 ID `topicId#taskId`，调用 Read 刷新
- [x] 2.5 实现 `resourceTencentCloudClsMetricSubscribeRead`：拆分复合 ID 获取 topicId/taskId，调用 `DescribeClsMetricSubscribeById`，判空后打印日志并 SetId("")，将返回的 MetricSubscribeInfo 字段 set 到 state（含 computed 字段）
- [x] 2.6 实现 `resourceTencentCloudClsMetricSubscribeUpdate`：拆分复合 ID，检查 name/namespace/metrics/instance_info/enable 变更，构建 `ModifyMetricSubscribeRequest`（含 TopicId/TaskId），调用 API（WriteRetryTimeout 重试），调用 Read 刷新
- [x] 2.7 实现 `resourceTencentCloudClsMetricSubscribeDelete`：拆分复合 ID，构建 `DeleteMetricSubscribeRequest`（含 TaskId/TopicId），调用 API（WriteRetryTimeout 重试）
- [x] 2.8 在资源定义中注册 Importer（`schema.ImportStatePassthrough`）

## 3. Provider Registration

- [x] 3.1 在 `tencentcloud/provider.go` 的 resources map 中注册 `tencentcloud_cls_metric_subscribe` → `cls.ResourceTencentCloudClsMetricSubscribe()`
- [x] 3.2 在 `tencentcloud/provider.md` 中添加 `tencentcloud_cls_metric_subscribe` 资源条目

## 4. Documentation

- [x] 4.1 创建 `tencentcloud/services/cls/resource_tc_cls_metric_subscribe.md`，包含一句话描述（带 CLS 产品名）、Example Usage（含 metrics 和 instance_info 嵌套块示例）、Import 部分（说明使用复合 ID `topicId#taskId`）

## 5. Unit Tests

- [x] 5.1 创建 `tencentcloud/services/cls/resource_tc_cls_metric_subscribe_test.go`，使用 gomonkey mock 云 API，补充 Create/Read/Update/Delete 业务逻辑的单元测试用例（不使用 terraform 测试套件）

## 6. Validation

- [x] 6.1 验证代码编译通过（go build 由其他流程执行，仅检查代码正确性）
- [x] 6.2 检查所有函数返回的 error 是否已正确处理
