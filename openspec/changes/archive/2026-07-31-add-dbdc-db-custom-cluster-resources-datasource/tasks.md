## 1. Service Layer

- [x] 1.1 在 `tencentcloud/services/dbdc/service_tencentcloud_dbdc.go` 中新增 `DescribeDBCustomClusterResources(ctx context.Context, clusterId string) (ret *dbdcv20201029.DescribeDBCustomClusterResourcesResponseParams, errRet error)` 方法 — 封装 `DescribeDBCustomClusterResourcesWithContext` 调用，入参 `ClusterId`，使用 `resource.Retry(ReadRetryTimeout)` + `ratelimit.Check` 包装，retry 块内检查 `result == nil || result.Response == nil` 返回 `NonRetryableError`，风格参考 `DescribeDBCustomClusterById`

## 2. Schema 定义与数据源创建

- [x] 2.1 创建 `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_cluster_resources.go` 文件
- [x] 2.2 定义数据源函数 `DataSourceTencentCloudDbdcDbCustomClusterResources()`，仅含 Read 方法
- [x] 2.3 定义输入参数 Schema
  - `cluster_id` (Required, String) — 集群 ID，映射到 API 入参 `ClusterId`
  - `result_output_file` (Optional, String) — 结果输出文件
- [x] 2.4 定义输出参数 Schema（平铺到顶层，非嵌套 list）
  - `node_count` (Computed, TypeInt) — 工作节点总数
  - `capacity` (Computed, TypeList, MaxItems:1) — 资源物理总容量，嵌套 `cpu`(TypeFloat)/`memory`(TypeFloat)/`pods`(TypeInt)
  - `allocatable` (Computed, TypeList, MaxItems:1) — 可分配容量，嵌套 `cpu`/`memory`/`pods`
  - `requests` (Computed, TypeList, MaxItems:1) — 资源申请量，嵌套 `cpu`/`memory`/`pods`
  - `limits` (Computed, TypeList, MaxItems:1) — 资源上限，嵌套 `cpu`/`memory`/`pods`
  - `available` (Computed, TypeList, MaxItems:1) — 可用余量，嵌套 `cpu`/`memory`/`pods`

## 3. Read 函数实现

- [x] 3.1 实现 `dataSourceTencentCloudDbdcDbCustomClusterResourcesRead()` 函数，包含 `defer tccommon.LogElapsed()` 和 `defer tccommon.InconsistentCheck()`
- [x] 3.2 从 schema 读取 `cluster_id`，调用 service 层 `DescribeDBCustomClusterResources(ctx, clusterId)`
- [x] 3.3 使用 `resource.Retry(ReadRetryTimeout)` 包装调用，retry 块内检查 `response == nil` 或 `response.Response == nil` 返回 `NonRetryableError`，外层失败路径打印 `[DATASOURCE] read empty, skip SetId` 日志
- [x] 3.4 处理 API 响应，将 `NodeCount`、`Capacity`、`Allocatable`、`Requests`、`Limits`、`Available` 设置到 state，所有指针/嵌套对象做 nil 检查
- [x] 3.5 将 `MetaResource` 嵌套对象（Capacity/Allocatable/Requests/Limits/Available）转换为 TypeList map，含 `cpu`/`memory`/`pods` 字段，每个字段做 nil 检查
- [x] 3.6 使用 `helper.BuildToken()` 设置数据源 ID
- [x] 3.7 实现 `result_output_file` 输出，若指定则调用 `tccommon.WriteToFile` 写入结果

## 4. 测试文件

- [x] 4.1 创建 `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_cluster_resources_test.go` 文件
- [x] 4.2 使用 gomonkey mock 方法对云 API 进行 mock 处理，编写单元测试用例覆盖 Read 函数的业务逻辑（非 terraform 测试套件）
- [x] 4.3 测试场景：正常查询返回资源信息、API 返回空响应时返回错误、嵌套对象 nil 安全处理

## 5. 文档

- [x] 5.1 创建 `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_cluster_resources.md` 文件
- [x] 5.2 添加一句话描述（"Use this data source to query ..."，带 DBDC 产品名称）
- [x] 5.3 添加 Example Usage 部分（按 cluster_id 查询示例）
- [x] 5.4 不添加 Argument Reference 和 Attribute Reference 部分（由工具自动生成）

## 6. Provider 注册

- [x] 6.1 在 `tencentcloud/provider.go` 的 Data Source Map 中注册 `"tencentcloud_dbdc_db_custom_cluster_resources": dbdc.DataSourceTencentCloudDbdcDbCustomClusterResources()`
- [x] 6.2 在 `tencentcloud/provider.md` 的 DBDC Data Source 部分添加 `tencentcloud_dbdc_db_custom_cluster_resources` 声明

## 7. 代码质量

- [x] 7.1 确保代码遵循项目规范（函数命名、导入别名 tccommon/helper、日志记录、错误处理）
- [x] 7.2 确保所有函数返回的 error 被检查，不出错的用 `_ =` 赋值
- [x] 7.3 确保不添加 go build/vet/lint 等命令执行（收尾阶段 gofmt 除外）
