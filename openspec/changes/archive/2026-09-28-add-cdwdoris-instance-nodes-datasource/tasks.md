## 1. Service 层实现

- [x] 1.1 在 `tencentcloud/services/cdwdoris/service_tencentcloud_cdwdoris.go` 新增方法 `DescribeCdwdorisInstanceNodesByFilter(ctx, param map[string]interface{}) (instanceNodes []*cdwdorisv20211228.InstanceNode, nodeRoles []*string, errRet error)`，构造 `DescribeInstanceNodesRequest`，按 `InstanceId`/`NodeRole`/`DisplayPolicy` 赋值，内部 Offset/Limit 循环分页（Limit=10）聚合全部 `InstanceNodesList`，同时返回 `NodeRoles`，含 `ratelimit.Check` 与 defer 错误日志，参照 `DescribeCdwdorisInstancesByFilter` 风格
- [x] 1.2 确认 service 方法通过 `me.client.UseCdwdorisV20211228Client().DescribeInstanceNodes(request)` 调用，并处理 `response == nil || response.Response == nil` 的边界

## 2. 数据源实现

- [x] 2.1 新建 `tencentcloud/services/cdwdoris/data_source_tc_cdwdoris_instance_nodes.go`，定义 `DataSourceTencentCloudCdwdorisInstanceNodes() *schema.Resource`，仅 `Read`，Schema 入参 `instance_id`(Required)、`node_role`(Optional)、`display_policy`(Optional)，Computed 列表 `instance_nodes_list`(15 字段)、`node_roles`，及 `result_output_file`(Optional)
- [x] 2.2 实现 `dataSourceTencentCloudCdwdorisInstanceNodesRead`：defer LogElapsed/InconsistentCheck，构造 paramMap，`resource.Retry(tccommon.ReadRetryTimeout, ...)` 调用 service 方法，retry 块内对空响应返回 NonRetryableError，外层失败路径保留 `log.Printf("[DATASOURCE] read empty, skip SetId")`，逐字段 nil 检查后 `d.Set`，`d.SetId(helper.BuildToken())`，支持 `result_output_file` 输出

## 3. Provider 注册

- [x] 3.1 在 `tencentcloud/provider.go` 的 dataSources map 中注册 `"tencentcloud_cdwdoris_instance_nodes": cdwdoris.DataSourceTencentCloudCdwdorisInstanceNodes(),`
- [x] 3.2 在 `tencentcloud/provider.md` 对应数据源列表中补充 `tencentcloud_cdwdoris_instance_nodes` 条目（由收尾阶段 `make doc` 生成，本阶段不手改 website/）

## 4. 文档

- [x] 4.1 新建 `tencentcloud/services/cdwdoris/data_source_tc_cdwdoris_instance_nodes.md`，包含一句话描述（提及 CDWDoris 产品）、Example Usage（按 instance_id 查询、含 node_role/display_policy 可选过滤示例）

## 5. 测试

- [x] 5.1 新建 `tencentcloud/services/cdwdoris/data_source_tc_cdwdoris_instance_nodes_test.go`，使用 gomonkey mock 云 API `DescribeInstanceNodes`，进业务逻辑单元测试，覆盖正常返回（含全字段）、空响应返回 NonRetryableError、node_role/display_policy 过滤入参传递等场景

## 6. 验证（收尾阶段执行）

- [ ] 6.1 执行 `gofmt` 格式化变更的 Go 代码（由 tfpacer-finalize skill 执行）
- [ ] 6.2 执行 `make doc` 生成 website/docs 文档与 provider.md 条目（由 tfpacer-finalize skill 执行）
- [ ] 6.3 在 `.changelog/` 下新增 changelog 文件（由 tfpacer-finalize skill 执行）