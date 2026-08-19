## 1. 数据源 Schema 定义

- [x] 1.1 在 `tencentcloud/services/teo/` 目录下创建 `data_source_tc_teo_origin_group_health_status.go` 文件
- [x] 1.2 定义 `DataSourceTencentCloudTeoOriginGroupHealthStatus()` 函数，返回 `*schema.Resource`
- [x] 1.3 在 Schema 中定义必填输入参数：`zone_id`（TypeString, Required）、`lb_instance_id`（TypeString, Required）
- [x] 1.4 在 Schema 中定义可选输入参数：`origin_group_ids`（TypeList of TypeString, Optional）、`result_output_file`（TypeString, Optional）
- [x] 1.5 在 Schema 中定义输出参数：`origin_group_health_status_list`（TypeList, Computed），包含嵌套结构：
  - `origin_group_id`（TypeString, Computed）
  - `origin_health_status`（TypeList, Computed）含 `origin`、`healthy`
  - `check_region_health_status`（TypeList, Computed）含 `region`、`healthy`、`origin_health_status`（含 `origin`、`healthy`）
- [x] 1.6 为所有字段添加 Description 说明

## 2. 数据源 Read 函数实现

- [x] 2.1 实现 `dataSourceTencentCloudTeoOriginGroupHealthStatusRead()` 函数
- [x] 2.2 添加 `defer tccommon.LogElapsed("data_source.tencentcloud_teo_origin_group_health_status.read")()` 记录执行时间
- [x] 2.3 添加 `defer tccommon.InconsistentCheck(d, meta)()` 处理数据一致性检查
- [x] 2.4 从 ResourceData 中读取 `zone_id`、`lb_instance_id`、`origin_group_ids` 参数
- [x] 2.5 构造 `teov20220901.DescribeOriginGroupHealthStatusRequest` 对象，映射所有参数
- [x] 2.6 使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包装 API 调用，处理最终一致性
- [x] 2.7 在 retry 块内检查 API 返回是否为空（`response == nil` 或 `response.Response == nil` 或 `len(response.Response.OriginGroupHealthStatusList) == 0`），若为空返回 `NonRetryableError`
- [x] 2.8 在外层 retry 失败路径上保留 `log.Printf("[DATASOURCE] read empty, skip SetId")` 提示
- [x] 2.9 调用接口失败时使用 `tccommon.RetryError()` 包装错误

## 3. API 响应数据处理

- [x] 3.1 解析 API 返回的 `OriginGroupHealthStatusList` 列表
- [x] 3.2 将每个 `OriginGroupHealthStatusDetail` 转换为 map[string]interface{}，映射 `origin_group_id` 字段
- [x] 3.3 映射 `origin_health_status` 列表，每个元素包含 `origin`、`healthy` 字段
- [x] 3.4 映射 `check_region_health_status` 列表，每个元素包含 `region`、`healthy`、`origin_health_status` 字段
- [x] 3.5 处理可能为 nil 的字段，跳过 setXX() 调用，避免空指针错误
- [x] 3.6 将转换后的列表设置到 `origin_group_health_status_list` 输出参数
- [x] 3.7 生成唯一的资源 ID（使用 `zone_id`、`lb_instance_id`、`origin_group_ids` 组合哈希）
- [x] 3.8 使用 `d.SetId()` 设置资源 ID（放在 retry 块外、retry 错误处理后）

## 4. 结果导出功能

- [x] 4.1 检查是否设置了 `result_output_file` 参数
- [x] 4.2 使用 `helper.WriteToFile()` 将查询结果导出为 JSON 文件
- [x] 4.3 确保导出的数据结构完整且格式正确

## 5. Provider 注册

- [x] 5.1 在 `tencentcloud/provider.go` 文件中找到数据源注册位置
- [x] 5.2 在 `DataSourcesMap` 中添加 `"tencentcloud_teo_origin_group_health_status": teo.DataSourceTencentCloudTeoOriginGroupHealthStatus()`
- [x] 5.3 在 `tencentcloud/provider.md` 中添加数据源注册说明

## 6. 文档创建

- [x] 6.1 在 `tencentcloud/services/teo/` 目录下创建 `data_source_tc_teo_origin_group_health_status.md` 文件
- [x] 6.2 编写一句话描述（带上 TEO 产品名称），格式为 "Use this data source to query ..."
- [x] 6.3 编写 Example Usage 示例（包含必填参数和可选参数两种场景）
- [x] 6.4 不添加 Argument Reference 和 Attribute Reference 章节（由工具自动生成）

## 7. 单元测试实现

- [x] 7.1 在 `tencentcloud/services/teo/` 目录下创建 `data_source_tc_teo_origin_group_health_status_test.go` 文件
- [x] 7.2 使用 gomonkey 对云 API 进行 mock 处理，不使用 terraform 测试套件
- [x] 7.3 编写测试用例验证必填参数查询场景
- [x] 7.4 编写测试用例验证可选 `origin_group_ids` 过滤场景
- [x] 7.5 编写测试用例验证空结果处理场景
- [x] 7.6 确保所有函数返回的 error 都被检查，必定不出错的用 `_ = func()` 赋值给 `_`
