## 1. 服务层实现 - 添加 OrganizationService 方法

- [ ] 1.1 在 `tencentcloud/services/tco/service_tencentcloud_organization.go` 中添加 `AddOrganizationOrgShareUnitNodeById` 方法
- [ ] 1.2 在 `AddOrganizationOrgShareUnitNodeById` 中调用 `organization.NewAddShareUnitNodeRequest()` 创建请求
- [ ] 1.3 设置请求参数: `request.UnitId` 和 `request.NodeId`
- [ ] 1.4 使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包装 API 调用
- [ ] 1.5 在 retry 函数中调用 `ratelimit.Check(request.GetAction())` 进行限流检查
- [ ] 1.6 调用 `me.client.UseOrganizationClient().AddShareUnitNode(request)` 执行 API
- [ ] 1.7 添加成功和失败的日志记录(使用 logId 和 ToJsonString())

- [ ] 1.8 在 `tencentcloud/services/tco/service_tencentcloud_organization.go` 中添加 `DeleteOrganizationOrgShareUnitNodeById` 方法
- [ ] 1.9 在 `DeleteOrganizationOrgShareUnitNodeById` 中调用 `organization.NewDeleteShareUnitNodeRequest()` 创建请求
- [ ] 1.10 设置请求参数并使用 retry 包装 API 调用
- [ ] 1.11 调用 `me.client.UseOrganizationClient().DeleteShareUnitNode(request)` 执行 API
- [ ] 1.12 添加日志记录

- [ ] 1.13 在 `tencentcloud/services/tco/service_tencentcloud_organization.go` 中添加 `DescribeOrganizationOrgShareUnitNodeById` 方法
- [ ] 1.14 在 `DescribeOrganizationOrgShareUnitNodeById` 中调用 `organization.NewDescribeShareUnitNodesRequest()` 创建请求
- [ ] 1.15 设置请求参数: `UnitId`, `SearchKey`(使用 NodeId), `Limit`, `Offset`
- [ ] 1.16 实现分页查询逻辑(循环调用直到找到匹配的 NodeId 或遍历完所有结果)
- [ ] 1.17 返回匹配的 `*organization.ShareUnitNode` 或 nil

- [ ] 1.18 在 `tencentcloud/services/tco/service_tencentcloud_organization.go` 中添加 `DescribeOrganizationOrgShareUnitNodesByFilter` 方法
- [ ] 1.19 在 `DescribeOrganizationOrgShareUnitNodesByFilter` 中从 param map 中提取参数(UnitId, Offset, Limit, SearchKey)
- [ ] 1.20 实现分页查询逻辑(循环调用直到获取所有结果)
- [ ] 1.21 返回 `[]*organization.ShareUnitNode` 数组

## 2. 资源实现 - resource_tc_organization_org_share_unit_node.go

- [ ] 2.1 创建文件 `tencentcloud/services/tco/resource_tc_organization_org_share_unit_node.go`
- [ ] 2.2 添加 package 声明和必要的 imports
- [ ] 2.3 实现 `ResourceTencentCloudOrganizationOrgShareUnitNode()` 函数返回 `*schema.Resource`

- [ ] 2.4 在 Resource Schema 中定义 `unit_id` 字段(TypeString, Required, ForceNew, Description: "共享单元ID")
- [ ] 2.5 在 Resource Schema 中定义 `node_id` 字段(TypeInt, Required, ForceNew, Description: "组织部门ID")
- [ ] 2.6 添加 Importer 配置: `State: schema.ImportStatePassthrough`

- [ ] 2.7 实现 `resourceTencentCloudOrganizationOrgShareUnitNodeCreate` 函数
- [ ] 2.8 在 Create 函数开始处添加 `defer tccommon.LogElapsed()` 和 `defer tccommon.InconsistentCheck()`
- [ ] 2.9 从 `d.GetOk()` 中获取 unit_id 和 node_id
- [ ] 2.10 调用服务层的 `AddOrganizationOrgShareUnitNodeById` 方法
- [ ] 2.11 设置复合 ID: `d.SetId(strings.Join([]string{unitId, strconv.FormatInt(nodeId, 10)}, tccommon.FILED_SP))`
- [ ] 2.12 调用 Read 函数更新 state

- [ ] 2.13 实现 `resourceTencentCloudOrganizationOrgShareUnitNodeRead` 函数
- [ ] 2.14 在 Read 函数开始处添加 defer 语句
- [ ] 2.15 使用 `strings.Split(d.Id(), tccommon.FILED_SP)` 解析复合 ID
- [ ] 2.16 验证 ID 格式(长度必须为 2),不符合则返回错误
- [ ] 2.17 使用 `strconv.ParseInt` 将 NodeId 字符串转换为 int64
- [ ] 2.18 调用服务层的 `DescribeOrganizationOrgShareUnitNodeById` 方法
- [ ] 2.19 如果返回 nil,设置 `d.SetId("")` 并打印 WARN 日志
- [ ] 2.20 使用 `d.Set()` 设置 unit_id 和 node_id 到 state

- [ ] 2.21 实现 `resourceTencentCloudOrganizationOrgShareUnitNodeDelete` 函数
- [ ] 2.22 在 Delete 函数开始处添加 defer 语句
- [ ] 2.23 解析复合 ID 并转换 NodeId
- [ ] 2.24 调用服务层的 `DeleteOrganizationOrgShareUnitNodeById` 方法
- [ ] 2.25 处理删除错误(特别是 FailedOperation.ShareNodeNotExist)

## 3. 数据源实现 - data_source_tc_organization_org_share_unit_nodes.go

- [ ] 3.1 创建文件 `tencentcloud/services/tco/data_source_tc_organization_org_share_unit_nodes.go`
- [ ] 3.2 添加 package 声明和必要的 imports
- [ ] 3.3 实现 `DataSourceTencentCloudOrganizationOrgShareUnitNodes()` 函数返回 `*schema.Resource`

- [ ] 3.4 在 Data Source Schema 中定义 `unit_id` 字段(TypeString, Required, Description: "共享单元ID")
- [ ] 3.5 定义 `offset` 字段(TypeInt, Optional, Description: "偏移量,默认为0")
- [ ] 3.6 定义 `limit` 字段(TypeInt, Optional, Description: "限制数目,取值范围1-50,默认为10")
- [ ] 3.7 定义 `search_key` 字段(TypeString, Optional, Description: "搜索关键字,支持部门ID搜索")
- [ ] 3.8 定义 `items` 字段(TypeList, Computed, Elem 为 Resource 类型)
- [ ] 3.9 在 items 的 Elem Schema 中定义 `share_node_id`(TypeInt, Required)
- [ ] 3.10 在 items 的 Elem Schema 中定义 `create_time`(TypeString, Required)
- [ ] 3.11 定义 `result_output_file` 字段(TypeString, Optional)

- [ ] 3.12 实现 `dataSourceTencentCloudOrganizationOrgShareUnitNodesRead` 函数
- [ ] 3.13 在 Read 函数开始处添加 defer 语句
- [ ] 3.14 从 d.GetOk() 中获取 unit_id, offset, limit, search_key
- [ ] 3.15 构造 param map 传递给服务层方法
- [ ] 3.16 调用服务层的 `DescribeOrganizationOrgShareUnitNodesByFilter` 方法
- [ ] 3.17 遍历返回的节点列表,构造 tmpList([]interface{})
- [ ] 3.18 每个节点构造 map[string]interface{} 包含 share_node_id 和 create_time
- [ ] 3.19 使用 `d.Set("items", tmpList)` 设置 items
- [ ] 3.20 如果提供了 result_output_file,使用 `tccommon.WriteToFile()` 写入文件
- [ ] 3.21 使用 `d.SetId(dataSourceIdHash(idParam))` 生成唯一 ID

## 4. Provider 注册

- [ ] 4.1 打开 `tencentcloud/provider.go` 文件
- [ ] 4.2 在 ResourcesMap 中添加资源注册: `"tencentcloud_organization_org_share_unit_node": tco.ResourceTencentCloudOrganizationOrgShareUnitNode()`
- [ ] 4.3 在 DataSourcesMap 中添加数据源注册: `"tencentcloud_organization_org_share_unit_nodes": tco.DataSourceTencentCloudOrganizationOrgShareUnitNodes()`

## 5. 资源测试实现

- [ ] 5.1 创建文件 `tencentcloud/services/tco/resource_tc_organization_org_share_unit_node_test.go`
- [ ] 5.2 实现 `TestAccTencentCloudOrganizationOrgShareUnitNodeResource_basic` 测试函数
- [ ] 5.3 在测试中使用 `resource.Test()` 框架
- [ ] 5.4 定义 TestStep 1: 创建资源并验证存在
- [ ] 5.5 在 TestStep 1 的 Config 中使用测试用的 unit_id 和 node_id
- [ ] 5.6 在 TestStep 1 的 Check 中使用 `resource.TestCheckResourceAttr()` 验证 unit_id 和 node_id
- [ ] 5.7 定义 TestStep 2: 测试 Import 功能
- [ ] 5.8 在 TestStep 2 中设置 `ImportState: true`, `ImportStateVerify: true`
- [ ] 5.9 添加 PreCheck 函数验证必要的环境变量

## 6. 数据源测试实现

- [ ] 6.1 创建文件 `tencentcloud/services/tco/data_source_tc_organization_org_share_unit_nodes_test.go`
- [ ] 6.2 实现 `TestAccTencentCloudOrganizationOrgShareUnitNodesDataSource_basic` 测试函数
- [ ] 6.3 在测试中先创建资源,然后查询数据源
- [ ] 6.4 定义 TestStep 查询数据源并验证返回结果
- [ ] 6.5 使用 `resource.TestCheckResourceAttrSet()` 验证 items 字段存在
- [ ] 6.6 测试 search_key 过滤功能
- [ ] 6.7 测试分页参数(offset, limit)

## 7. 资源文档

- [ ] 7.1 创建文件 `tencentcloud/services/tco/resource_tc_organization_org_share_unit_node.md`
- [ ] 7.2 添加资源描述: "Provides a resource to create an organization org share unit node"
- [ ] 7.3 添加 Example Usage 部分,包含完整的 Terraform 配置示例
- [ ] 7.4 在示例中展示必需的 unit_id 和 node_id 参数
- [ ] 7.5 添加 Import 部分,说明导入格式: `terraform import tencentcloud_organization_org_share_unit_node.foo {unit_id}#{node_id}`
- [ ] 7.6 添加 Argument Reference 部分,列出所有参数及其说明
- [ ] 7.7 标记 unit_id 和 node_id 为 Required 和 ForceNew
- [ ] 7.8 添加 Attributes Reference 部分(如果有计算字段)

## 8. 数据源文档

- [ ] 8.1 创建文件 `tencentcloud/services/tco/data_source_tc_organization_org_share_unit_nodes.md`
- [ ] 8.2 添加数据源描述: "Use this data source to query organization org share unit nodes"
- [ ] 8.3 添加 Example Usage 部分,包含完整的 Terraform 配置示例
- [ ] 8.4 在示例中展示 unit_id(必需)和 search_key(可选)的使用
- [ ] 8.5 展示如何访问返回的 items 列表
- [ ] 8.6 添加 Argument Reference 部分,列出所有输入参数
- [ ] 8.7 添加 Attributes Reference 部分,说明返回的 items 结构
- [ ] 8.8 描述 items 中的 share_node_id 和 create_time 字段

## 9. 生成文档

- [ ] 9.1 运行 `make doc` 命令生成 website 文档
- [ ] 9.2 验证生成的 `website/docs/r/organization_org_share_unit_node.html.markdown` 文件
- [ ] 9.3 验证生成的 `website/docs/d/organization_org_share_unit_nodes.html.markdown` 文件
- [ ] 9.4 检查生成文档的格式和内容是否正确

## 10. 代码验证

- [ ] 10.1 运行 `go build ./tencentcloud/services/tco/...` 验证编译通过
- [ ] 10.2 运行 `go vet ./tencentcloud/services/tco/...` 检查代码问题
- [ ] 10.3 运行 `gofmt -l tencentcloud/services/tco/` 检查代码格式
- [ ] 10.4 如果有格式问题,运行 `gofmt -w` 修复

## 11. 运行测试

- [ ] 11.1 设置测试所需的环境变量(TENCENTCLOUD_SECRET_ID, TENCENTCLOUD_SECRET_KEY)
- [ ] 11.2 运行资源测试: `TF_ACC=1 go test ./tencentcloud/services/tco -run TestAccTencentCloudOrganizationOrgShareUnitNodeResource_basic -v`
- [ ] 11.3 运行数据源测试: `TF_ACC=1 go test ./tencentcloud/services/tco -run TestAccTencentCloudOrganizationOrgShareUnitNodesDataSource_basic -v`
- [ ] 11.4 验证所有测试通过

## 12. 最终检查

- [ ] 12.1 确认所有新文件都已创建并包含必要的代码
- [ ] 12.2 确认 provider.go 中已注册新资源和数据源
- [ ] 12.3 确认资源支持 CRUD 和 Import 操作
- [ ] 12.4 确认数据源支持分页和搜索
- [ ] 12.5 确认文档完整且格式正确
- [ ] 12.6 确认测试覆盖主要功能场景
- [ ] 12.7 确认代码遵循项目规范(错误处理、日志记录、ID 格式等)
