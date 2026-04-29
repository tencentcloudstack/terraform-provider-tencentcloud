## 1. 数据源 Schema 定义

- [x] 1.1 在 `tencentcloud/services/cvm/` 目录下创建 `data_source_tc_cvm_repair_tasks.go` 文件
- [x] 1.2 定义 `DataSourceTencentCloudCvmRepairTasks()` 函数,返回 `*schema.Resource`
- [x] 1.3 在 Schema 中定义所有输入参数:product, task_status, task_type_ids, task_ids, instance_ids, aliases, start_date, end_date, order_field, order, result_output_file (不包含 limit 和 offset)
- [x] 1.4 在 Schema 中定义输出参数:repair_task_list (包含所有任务详情字段), total_count
- [x] 1.5 为数组类型参数(task_status, task_type_ids, task_ids, instance_ids, aliases)使用 TypeSet
- [x] 1.6 为所有字段添加 Description 说明

## 2. 数据源 Read 函数实现

- [x] 2.1 实现 `dataSourceTencentCloudCvmRepairTasksRead()` 函数
- [x] 2.2 从 ResourceData 中读取所有输入参数
- [x] 2.3 将 TypeSet 参数转换为 API 所需的格式(如 []*int64, []*string)
- [x] 2.4 构造 `cvm.DescribeTaskInfoRequest` 对象,映射所有参数
- [x] 2.5 使用 `helper.Retry()` 包装 API 调用,处理最终一致性
- [x] 2.6 添加 `defer tccommon.LogElapsed()` 记录执行时间
- [x] 2.7 添加 `defer tccommon.InconsistentCheck()` 处理数据一致性检查

## 3. 自动分页实现

- [x] 3.1 实现 for 循环进行自动分页,初始 Offset=0, Limit=100
- [x] 3.2 每次 API 调用后检查返回数据量,如果 >= 100 则继续下一页
- [x] 3.3 累加 Offset (每次 += 100),直到获取所有数据
- [x] 3.4 将所有分页结果合并到一个列表中
- [x] 3.5 处理分页过程中的错误,确保任何一次 API 调用失败都能正确返回错误

## 4. API 响应数据处理

- [x] 4.1 解析 API 返回的 `RepairTaskInfoSet` 列表
- [x] 4.2 将每个 `RepairTaskInfo` 转换为 map[string]interface{},映射所有字段
- [x] 4.3 处理可能为 nil 的字段,避免空指针错误
- [x] 4.4 将转换后的任务列表设置到 `repair_task_list` 输出参数
- [x] 4.5 设置 `total_count` 输出参数(等于实际返回的任务总数)
- [x] 4.6 生成唯一的资源 ID (使用参数哈希或时间戳)

## 5. 结果导出功能

- [x] 5.1 检查是否设置了 `result_output_file` 参数
- [x] 5.2 使用 `helper.WriteToFile()` 将查询结果导出为 JSON 文件
- [x] 5.3 确保导出的数据结构完整且格式正确

## 6. Provider 注册

- [x] 6.1 在 `tencentcloud/provider.go` 文件中找到数据源注册位置
- [x] 6.2 在 `DataSourcesMap` 中添加 `"tencentcloud_cvm_repair_tasks": cvm.DataSourceTencentCloudCvmRepairTasks()`

## 7. 示例文件创建

- [x] 7.1 在 `tencentcloud/services/cvm/` 目录下创建 `data_source_tc_cvm_repair_tasks.md` 文件
- [x] 7.2 编写基础查询示例(不带过滤条件)
- [x] 7.3 编写按任务状态过滤的示例
- [x] 7.4 编写按实例ID过滤的示例
- [x] 7.5 编写带排序的示例
- [x] 7.6 添加 Argument Reference 章节,列出所有输入参数及说明(不包含 limit/offset)
- [x] 7.7 添加 Attributes Reference 章节,列出所有输出属性及说明

## 8. 验收测试实现

- [x] 8.1 在 `tencentcloud/services/cvm/` 目录下创建 `data_source_tc_cvm_repair_tasks_test.go` 文件
- [x] 8.2 实现基础查询测试用例 `TestAccTencentCloudCvmRepairTasksDataSource_basic`
- [x] 8.3 实现带任务状态过滤的测试用例
- [x] 8.4 在测试用例中验证返回的数据结构正确性
- [x] 8.5 使用 `resource.TestCheckResourceAttrSet()` 验证关键字段存在
- [ ] 8.6 如果可能,测试自动分页功能(需要足够多的测试数据)

## 9. 文档生成

- [x] 9.1 运行 `make doc` 命令生成 `website/docs/d/cvm_repair_tasks.html.markdown` 文档
- [x] 9.2 验证生成的文档格式正确,包含所有参数说明和示例

## 10. 代码质量检查

- [x] 9.1 运行 `go fmt` 格式化代码
- [x] 9.2 运行 `go vet` 检查代码问题
- [x] 9.3 运行 golangci-lint 进行静态代码分析
- [x] 9.4 修复所有 linter 报告的问题

## 10. 功能验证

- [ ] 10.1 设置 TF_ACC=1, TENCENTCLOUD_SECRET_ID, TENCENTCLOUD_SECRET_KEY 环境变量
- [ ] 10.2 运行验收测试:`go test -v ./tencentcloud/services/cvm/data_source_tc_cvm_repair_tasks_test.go`
- [ ] 10.3 验证所有测试用例通过
- [ ] 10.4 手动运行 Terraform 配置,测试数据源在实际场景中的使用

## 11. 提交准备

- [ ] 11.1 在 `.changelog/` 目录下创建 changelog 文件,记录新增功能
- [ ] 11.2 检查所有修改的文件是否符合编码规范
- [ ] 11.3 确认没有引入破坏性变更
- [ ] 11.4 准备 PR 描述,说明新增功能和使用方法
