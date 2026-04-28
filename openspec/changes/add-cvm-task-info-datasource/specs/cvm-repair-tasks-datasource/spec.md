## ADDED Requirements

### Requirement: Data source registration
数据源 MUST 在 Provider 中注册为 `tencentcloud_cvm_repair_tasks`,使其可在 Terraform 配置中使用。

#### Scenario: Data source is accessible
- **WHEN** 用户在 Terraform 配置中使用 `data "tencentcloud_cvm_repair_tasks"`
- **THEN** Terraform 能够成功识别并初始化该数据源

### Requirement: Query all repair tasks
数据源 MUST 支持不带任何过滤条件的查询,返回默认数量的维修任务列表。

#### Scenario: Query without filters
- **WHEN** 用户未设置任何过滤参数
- **THEN** 数据源调用 API 时使用默认的 Limit=20 和 Offset=0
- **THEN** 返回最多 20 条维修任务记录

### Requirement: Filter by product type
数据源 MUST 支持通过 `product` 参数过滤产品类型。

#### Scenario: Filter CVM tasks
- **WHEN** 用户设置 `product = "CVM"`
- **THEN** 只返回云服务器(CVM)的维修任务

#### Scenario: Filter CDH tasks
- **WHEN** 用户设置 `product = "CDH"`
- **THEN** 只返回专用宿主机(CDH)的维修任务

#### Scenario: Filter CPM tasks
- **WHEN** 用户设置 `product = "CPM2.0"`
- **THEN** 只返回裸金属云服务器(CPM2.0)的维修任务

### Requirement: Filter by task status
数据源 MUST 支持通过 `task_status` 参数过滤任务状态,支持多个状态同时过滤。

#### Scenario: Filter pending authorization tasks
- **WHEN** 用户设置 `task_status = [1]`
- **THEN** 只返回状态为"待授权"的维修任务

#### Scenario: Filter multiple statuses
- **WHEN** 用户设置 `task_status = [1, 2, 4]`
- **THEN** 返回状态为"待授权"、"处理中"或"已预约"的维修任务

#### Scenario: Valid status values
- **WHEN** 用户设置 `task_status` 参数
- **THEN** 支持的值为: 1(待授权), 2(处理中), 3(已结束), 4(已预约), 5(已取消), 6(已避免)

### Requirement: Filter by task type
数据源 MUST 支持通过 `task_type_ids` 参数过滤任务类型,支持多个类型同时过滤。

#### Scenario: Filter instance hazard tasks
- **WHEN** 用户设置 `task_type_ids = [101]`
- **THEN** 只返回类型为"实例运行隐患"的维修任务

#### Scenario: Filter multiple task types
- **WHEN** 用户设置 `task_type_ids = [101, 102, 103]`
- **THEN** 返回"实例运行隐患"、"实例运行异常"或"实例硬盘异常"的维修任务

#### Scenario: Valid task type values
- **WHEN** 用户设置 `task_type_ids` 参数
- **THEN** 支持的值为: 101(实例运行隐患), 102(实例运行异常), 103(实例硬盘异常), 104(实例网络连接异常), 105(实例运行预警), 106(实例硬盘预警), 107(实例维护升级)

### Requirement: Filter by task IDs
数据源 MUST 支持通过 `task_ids` 参数根据任务ID精确查询,支持多个任务ID。

#### Scenario: Query specific task
- **WHEN** 用户设置 `task_ids = ["rep-12345678"]`
- **THEN** 只返回任务ID为 "rep-12345678" 的维修任务

#### Scenario: Query multiple tasks
- **WHEN** 用户设置 `task_ids = ["rep-12345678", "rep-87654321"]`
- **THEN** 返回这两个任务ID对应的维修任务

### Requirement: Filter by instance IDs
数据源 MUST 支持通过 `instance_ids` 参数根据实例ID查询,支持多个实例ID。

#### Scenario: Query tasks for specific instance
- **WHEN** 用户设置 `instance_ids = ["ins-12345678"]`
- **THEN** 返回实例ID为 "ins-12345678" 的所有维修任务

#### Scenario: Query tasks for multiple instances
- **WHEN** 用户设置 `instance_ids = ["ins-12345678", "ins-87654321"]`
- **THEN** 返回这两个实例的所有维修任务

### Requirement: Filter by instance names
数据源 MUST 支持通过 `aliases` 参数根据实例名称查询,支持多个实例名称。

#### Scenario: Query tasks by instance name
- **WHEN** 用户设置 `aliases = ["test-server-1"]`
- **THEN** 返回实例名称为 "test-server-1" 的所有维修任务

#### Scenario: Query tasks for multiple instance names
- **WHEN** 用户设置 `aliases = ["test-server-1", "prod-server-1"]`
- **THEN** 返回这两个实例名称对应的所有维修任务

### Requirement: Filter by time range
数据源 MUST 支持通过 `start_date` 和 `end_date` 参数根据任务创建时间过滤。

#### Scenario: Query tasks in date range
- **WHEN** 用户设置 `start_date = "2023-03-01 00:00:00"` 和 `end_date = "2023-04-01 00:00:00"`
- **THEN** 只返回创建时间在 2023-03-01 到 2023-04-01 之间的维修任务

#### Scenario: Query tasks from specific date
- **WHEN** 用户只设置 `start_date = "2023-03-01 00:00:00"`,不设置 `end_date`
- **THEN** 返回创建时间从 2023-03-01 到当前时刻的维修任务

#### Scenario: Query tasks until specific date
- **WHEN** 用户只设置 `end_date = "2023-04-01 00:00:00"`,不设置 `start_date`
- **THEN** 返回创建时间从当天 00:00:00 到 2023-04-01 的维修任务

#### Scenario: Date format validation
- **WHEN** 用户设置 `start_date` 或 `end_date` 参数
- **THEN** 格式 MUST 为 "YYYY-MM-DD hh:mm:ss"

### Requirement: Sorting control
数据源 MUST 支持通过 `order_field` 和 `order` 参数控制结果排序。

#### Scenario: Sort by create time ascending
- **WHEN** 用户设置 `order_field = "CreateTime"` 和 `order = 0`
- **THEN** 结果按创建时间升序排列

#### Scenario: Sort by auth time descending
- **WHEN** 用户设置 `order_field = "AuthTime"` 和 `order = 1`
- **THEN** 结果按授权时间降序排列

#### Scenario: Sort by end time
- **WHEN** 用户设置 `order_field = "EndTime"`
- **THEN** 结果按结束时间排列

#### Scenario: Default sorting
- **WHEN** 用户未设置排序参数
- **THEN** 结果按创建时间升序排列(默认行为)

#### Scenario: Valid order field values
- **WHEN** 用户设置 `order_field` 参数
- **THEN** 支持的值为: "CreateTime"(创建时间), "AuthTime"(授权时间), "EndTime"(结束时间)

### Requirement: Automatic pagination
数据源 MUST 内部实现自动分页,无需用户手动控制 limit 和 offset 参数,自动获取所有符合条件的数据。

#### Scenario: Automatically fetch all data
- **WHEN** 查询返回的数据总数超过单次请求的最大限制(100条)
- **THEN** 数据源自动进行多次 API 调用
- **THEN** 每次请求使用 Limit=100, 自动递增 Offset
- **THEN** 直到获取所有符合条件的数据

#### Scenario: Return complete result set
- **WHEN** 用户使用过滤条件查询
- **THEN** 返回所有符合条件的任务,不受单次 API 调用限制
- **THEN** `repair_task_list` 包含完整的任务列表
- **THEN** `total_count` 反映实际返回的任务总数

#### Scenario: Handle pagination automatically
- **WHEN** 第一次 API 调用返回 100 条数据
- **THEN** 数据源自动发起第二次请求,Offset=100
- **THEN** 继续请求直到返回数据量 < 100 或已获取全部数据

#### Scenario: Small result set
- **WHEN** 查询结果总数少于 100 条
- **THEN** 只需一次 API 调用
- **THEN** 返回所有结果

### Requirement: Return task details
数据源 MUST 返回完整的维修任务详细信息。

#### Scenario: Task basic information
- **WHEN** 查询成功返回任务列表
- **THEN** 每个任务 MUST 包含: task_id, instance_id, alias, task_type_id, task_status, product

#### Scenario: Task timing information
- **WHEN** 查询成功返回任务列表
- **THEN** 每个任务 MUST 包含: create_time, auth_time, end_time

#### Scenario: Task description information
- **WHEN** 查询成功返回任务列表
- **THEN** 每个任务 MUST 包含: task_detail, task_type_name, task_sub_type

#### Scenario: Task status information
- **WHEN** 查询成功返回任务列表
- **THEN** 每个任务 MUST 包含: device_status, operate_status, auth_type, auth_source

#### Scenario: Instance network information
- **WHEN** 查询成功返回任务列表
- **THEN** 每个任务 MUST 包含: zone, region, vpc_id, vpc_name, subnet_id, subnet_name, wan_ip, lan_ip

### Requirement: Return total count
数据源 MUST 返回符合条件的维修任务总数量。

#### Scenario: Total count with filters
- **WHEN** 用户使用过滤条件查询
- **THEN** `total_count` 返回符合过滤条件的任务总数(不受 limit 限制)

#### Scenario: Total count without filters
- **WHEN** 用户不使用过滤条件查询
- **THEN** `total_count` 返回用户账号下的所有维修任务总数

### Requirement: Export results to file
数据源 MUST 支持通过 `result_output_file` 参数将查询结果导出到 JSON 文件。

#### Scenario: Export to specified file
- **WHEN** 用户设置 `result_output_file = "/tmp/repair_tasks.json"`
- **THEN** 查询结果被写入到指定路径的 JSON 文件
- **THEN** 文件内容包含完整的任务列表数据

### Requirement: Error handling
数据源 MUST 正确处理 API 调用失败和参数错误的情况。

#### Scenario: API call failure
- **WHEN** API 调用返回错误(如网络问题、认证失败)
- **THEN** 数据源返回明确的错误信息给 Terraform
- **THEN** 错误信息包含原始 API 错误码和描述

#### Scenario: Invalid parameter
- **WHEN** 用户提供了无效的参数值(如错误的日期格式)
- **THEN** 数据源在调用 API 前进行验证
- **THEN** 返回清晰的参数错误提示

#### Scenario: Empty result
- **WHEN** API 返回空的任务列表
- **THEN** 数据源正常返回,`repair_task_list` 为空数组
- **THEN** `total_count` 为 0
- **THEN** 不报错

### Requirement: API retry for consistency
数据源 MUST 使用重试机制处理 API 最终一致性问题。

#### Scenario: Transient API error
- **WHEN** API 调用返回临时性错误(如限流、服务暂时不可用)
- **THEN** 数据源自动重试 API 调用
- **THEN** 在合理的重试次数内成功后返回结果

#### Scenario: Persistent API error
- **WHEN** API 调用多次重试后仍然失败
- **THEN** 数据源返回错误信息给 Terraform
- **THEN** 不进行无限重试

### Requirement: Documentation
数据源 MUST 提供完整的使用文档。

#### Scenario: Documentation includes examples
- **WHEN** 用户查看数据源文档
- **THEN** 文档包含至少 3 个使用示例:基础查询、带过滤条件查询、分页查询

#### Scenario: Documentation describes all parameters
- **WHEN** 用户查看数据源文档
- **THEN** 文档列出所有输入参数和输出属性
- **THEN** 每个参数都有清晰的描述和类型说明

#### Scenario: Documentation includes argument reference
- **WHEN** 用户查看数据源文档
- **THEN** 文档包含 "Argument Reference" 和 "Attributes Reference" 章节
- **THEN** 列出所有可选参数和计算属性

### Requirement: Acceptance testing
数据源 MUST 包含验收测试以确保功能正确性。

#### Scenario: Basic query test
- **WHEN** 运行验收测试
- **THEN** 包含测试用例验证不带过滤条件的基础查询

#### Scenario: Filter test
- **WHEN** 运行验收测试
- **THEN** 包含测试用例验证至少一种过滤条件(如按任务状态)

#### Scenario: Test with real API
- **WHEN** 设置 TF_ACC=1 环境变量
- **THEN** 验收测试使用真实的腾讯云 API
- **THEN** 测试需要有效的 TENCENTCLOUD_SECRET_ID 和 TENCENTCLOUD_SECRET_KEY
