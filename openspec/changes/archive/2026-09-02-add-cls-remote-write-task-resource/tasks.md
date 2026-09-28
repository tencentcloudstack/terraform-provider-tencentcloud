## 1. Schema 定义

- [x] 1.1 创建 `tencentcloud/services/cls/resource_tc_cls_remote_write_task.go`，定义 `ResourceTencentCloudClsRemoteWriteTask()` 资源函数，参考 `tencentcloud_igtm_strategy` 代码风格
- [x] 1.2 定义 schema 必填字段：`topic_id`(string, Required), `name`(string, Required), `target`(string, Required), `remote_write_url`(string, Required), `auth_type`(int, Required), `net_type`(int, Required)
- [x] 1.3 定义 schema 可选字段：`vpc_id`(string, Optional), `virtual_gateway_type`(int, Optional), `enable`(int, Optional, Computed)
- [x] 1.4 定义 `auth_info` 嵌套块（TypeList, Optional, MaxItems:1），含 `username`(string, Optional), `password`(string, Optional), `token`(string, Optional)
- [x] 1.5 定义只读 Computed 字段：`task_id`(string, Computed), `status`(int, Computed), `create_time`(string, Computed), `update_time`(string, Computed), `logset_id`(string, Computed)

## 2. Create 操作

- [x] 2.1 实现 `resourceTencentCloudClsRemoteWriteTaskCreate()` 函数，从 schema 读取参数构建 `CreateRemoteWriteTaskRequest`
- [x] 2.2 在 Create 中设置 `AuthInfo` 结构体（从 `auth_info` 块读取 username/password/token）
- [x] 2.3 使用 `tccommon.ReadRetryTimeout` 和 retry 包装 `CreateRemoteWriteTask` API 调用，错误用 `tccommon.RetryError()` 包装
- [x] 2.4 检查返回的 `Response.TaskId` 是否为空，为空则打印 logId 和 d.Id() 后返回 `NonRetryableError`
- [x] 2.5 设置复合 ID `d.SetId(taskId + tccommon.FILED_SP + topicId)`，并调用 Read 填充完整 state

## 3. Read 操作

- [x] 3.1 实现 `resourceTencentCloudClsRemoteWriteTaskRead()` 函数，从 `d.Id()` 解析出 `task_id` 和 `topic_id`
- [x] 3.2 构建 `DescribeRemoteWriteTasksRequest`，设置 `Filters`（Key=taskId, Values=[task_id]），Limit 设为 API 最大值 100
- [x] 3.3 使用 `tccommon.ReadRetryTimeout` 和 retry 包装 `DescribeRemoteWriteTasks` API 调用
- [x] 3.4 检查返回是否为空（response==nil 或 Infos 为空），先打印 `log.Printf("[CRUD] tencentcloud_cls_remote_write_task id=%s", d.Id())`，再 `d.SetId("")`
- [x] 3.5 从 `Infos[0]` 提取任务详情，对每个字段做 nil 检查后再 `d.Set()`，包括 `auth_info` 块的 nil 检查

## 4. Update 操作

- [x] 4.1 实现 `resourceTencentCloudClsRemoteWriteTaskUpdate()` 函数，从 `d.Id()` 解析 `task_id` 和 `topic_id`
- [x] 4.2 构建 `ModifyRemoteWriteTaskRequest`，设置 TaskId、TopicId 及所有变更的可选参数（name、net_type、vpc_id、target、remote_write_url、auth_type、enable、virtual_gateway_type）
- [x] 4.3 从 `auth_info` 块构建 `RemoteWriteAuthInfo` 结构体并设置到请求
- [x] 4.4 使用 `tccommon.ReadRetryTimeout` 和 retry 包装 `ModifyRemoteWriteTask` API 调用
- [x] 4.5 Modify 成功后调用 Read 刷新 state

## 5. Delete 操作

- [x] 5.1 实现 `resourceTencentCloudClsRemoteWriteTaskDelete()` 函数，从 `d.Id()` 解析 `task_id` 和 `topic_id`
- [x] 5.2 构建 `DeleteRemoteWriteTaskRequest`，设置 TaskId 和 TopicId
- [x] 5.3 使用 `tccommon.ReadRetryTimeout` 和 retry 包装 `DeleteRemoteWriteTask` API 调用
- [x] 5.4 处理删除后资源不存在的场景（避免重复删除报错）

## 6. Provider 注册

- [x] 6.1 在 `tencentcloud/provider.go` 中注册 `tencentcloud_cls_remote_write_task` 资源
- [x] 6.2 在 `tencentcloud/provider.md` 中添加资源文档条目

## 7. 资源文档

- [x] 7.1 创建 `tencentcloud/services/cls/resource_tc_cls_remote_write_task.md`，包含一句话描述（带上 CLS 产品名称）
- [x] 7.2 添加 Example Usage 部分（含 auth_info 嵌套块的 HCL 示例，json 字符串场景用 jsonencode()）
- [x] 7.3 添加 Import 部分，说明使用复合 ID 格式 `task_id#topic_id`

## 8. 单元测试

- [x] 8.1 创建 `tencentcloud/services/cls/resource_tc_cls_remote_write_task_test.go`，使用 gomonkey mock 云 API
- [x] 8.2 编写 Create 操作的单元测试（mock `CreateRemoteWriteTask` 和 `DescribeRemoteWriteTasks`）
- [x] 8.3 编写 Update 操作的单元测试（mock `ModifyRemoteWriteTask` 和 `DescribeRemoteWriteTasks`）
- [x] 8.4 编写 Delete 操作的单元测试（mock `DeleteRemoteWriteTask`）
- [x] 8.5 确保所有 error 返回值都被检查，必定不出错的用 `_ = func()` 赋值

## 9. 验证

- [x] 9.1 检查所有 CRUD 接口参数映射的正确性：Create 参数在 CreateRemoteWriteTask 中存在，Update 参数在 ModifyRemoteWriteTask 中存在
- [x] 9.2 确认代码遵循项目规范（资源名称使用小写蛇形命名、retry 块外设置 id、nil 检查后再 set）
