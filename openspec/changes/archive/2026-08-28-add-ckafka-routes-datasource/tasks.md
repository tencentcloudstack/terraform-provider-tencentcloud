## 1. 服务层方法实现

- [x] 1.1 在 `tencentcloud/services/ckafka/service_tencentcloud_ckafka.go` 中新增 `DescribeCkafkaRouteByFilter` 方法，接收 `instanceId`、`routeId`、`mainRouteFlag` 参数，构造 `DescribeRouteRequest` 并调用 `DescribeRoute`，返回 `*ckafka.RouteResponse`
- [x] 1.2 方法中添加 `ratelimit.Check`、错误日志 defer、`response == nil` / `response.Response == nil` / `response.Response.Result == nil` 校验

## 2. 数据源核心功能实现

- [x] 2.1 创建 `tencentcloud/services/ckafka/data_source_tc_ckafka_routes.go` 文件
- [x] 2.2 实现 `DataSourceTencentCloudCkafkaRoutes()` Schema 定义：输入参数 `instance_id`(Required)、`route_id`(Optional,Int)、`main_route_flag`(Optional,Bool)、`result_output_file`(Optional)；输出 `routers`(TypeList)，元素含 `access_type`、`route_id`、`vip_type`、`vip_list`(vip/vport)、`domain`、`domain_port`、`delete_timestamp`、`subnet`、`broker_vip_list`(vip/vport)、`vpc_id`、`note`、`status`
- [x] 2.3 实现 `dataSourceTencentCloudCkafkaRoutesRead` 函数：`defer tccommon.LogElapsed` 与 `defer tccommon.InconsistentCheck`
- [x] 2.4 在 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 块内调用 `service.DescribeCkafkaRouteByFilter`，使用 `tccommon.RetryError` 包装错误
- [x] 2.5 retry 块内检查返回为空（`result == nil` / `len(result.Routers) == 0`）时返回 `NonRetryableError`，不执行 `d.SetId("")`，并在失败路径保留 `log.Printf("[DATASOURCE] read empty, skip SetId")` 日志
- [x] 2.6 遍历 `result.Routers`，逐字段判断 nil 后 `d.Set` 到 `routers` 列表（含嵌套 `vip_list`、`broker_vip_list` 的展开）
- [x] 2.7 使用 `helper.BuildToken()` 设置数据源 ID
- [x] 2.8 处理 `result_output_file` 输出

## 3. 集成到 Provider

- [x] 3.1 在 `tencentcloud/provider.go` 的 `DataSourcesMap` 中注册 `tencentcloud_ckafka_routes`，参考 `tencentcloud_ckafka_version` 注册位置

## 4. 编写文档

- [x] 4.1 创建 `tencentcloud/services/ckafka/data_source_tc_ckafka_routes.md`，一句话描述带上 CKafka 产品名，包含 Example Usage，不添加 Argument/Attribute Reference

## 5. 编写单元测试

- [x] 5.1 创建 `tencentcloud/services/ckafka/data_source_tc_ckafka_routes_test.go`，使用 gomonkey mock 云 API `DescribeRoute` 返回值
- [x] 5.2 编写测试用例验证 Read 函数对响应字段的映射逻辑（`routers` 列表展开、嵌套 `vip_list`/`broker_vip_list`、nil 字段跳过）

## 6. 验证（收尾阶段由 tfpacer-finalize 执行）

- [ ] 6.1 运行 `gofmt` 格式化代码（由 tfpacer-finalize skill 执行）
- [ ] 6.2 运行 `make doc` 生成 `website/docs/` 文档与 `provider.md`（由 tfpacer-finalize skill 执行）
- [ ] 6.3 在 `.changelog/` 下生成 changelog 文件（由 tfpacer-finalize skill 执行）
