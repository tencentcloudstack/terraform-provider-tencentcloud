## 1. Schema 定义与 CRUD 函数

- [x] 1.1 创建 `tencentcloud/services/teo/resource_tc_teo_purge_task_operation.go`，定义 `ResourceTencentCloudTeoPurgeTaskOperation()` 函数，包含 schema 定义（zone_id、type、method、targets、cache_tag 输入字段和 job_id、tasks 计算字段），所有输入字段设置 ForceNew: true
- [x] 1.2 实现 Create 函数：调用 `CreatePurgeTask` API，使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 进行重试，检查返回的 JobId 是否为空
- [x] 1.3 实现 Create 函数中的异步轮询：调用 `DescribePurgeTasks` API 使用 job-id 过滤条件轮询任务状态，直到状态为终态（success/failed/timeout/canceled），使用 `tccommon.ReadRetryTimeout` 作为超时时间
- [x] 1.4 实现 Create 函数中的计算属性设置：轮询完成后，从 DescribePurgeTasks 响应中设置 tasks 计算属性，使用 `helper.BuildToken()` 设置资源 ID
- [x] 1.5 实现 Read 函数（空函数，返回 nil）
- [x] 1.6 实现 Delete 函数（空函数，返回 nil）

## 2. Provider 注册

- [x] 2.1 在 `tencentcloud/provider.go` 的 ResourcesMap 中添加 `tencentcloud_teo_purge_task` 资源注册
- [x] 2.2 在 `tencentcloud/provider.md` 中添加 TEO 产品的 `tencentcloud_teo_purge_task` 资源文档索引

## 3. 文档

- [x] 3.1 创建 `tencentcloud/services/teo/resource_tc_teo_purge_task_operation.md` 资源示例文档，包含一句话描述、Example Usage（purge_url 和 purge_cache_tag 两种场景）和 Import 说明

## 4. 单元测试

- [x] 4.1 创建 `tencentcloud/services/teo/resource_tc_teo_purge_task_operation_test.go`，使用 gomonkey mock 云 API 调用，编写 purge_url 类型创建成功的测试用例
- [x] 4.2 添加 purge_cache_tag 类型创建成功的测试用例
- [x] 4.3 添加 API 错误处理的测试用例
- [x] 4.4 使用 `go test -gcflags=all=-l` 运行单元测试并确保通过
