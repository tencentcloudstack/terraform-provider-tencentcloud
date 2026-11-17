## 1. Schema 新增 security_group_ids 字段

- [x] 1.1 在 `tencentcloud/services/ckafka/resource_tc_ckafka_route.go` 的 `ResourceTencentCloudCkafkaRoute()` schema 中新增 `security_group_ids` 字段：类型 `schema.TypeList`，元素类型 `schema.TypeString`，属性 `Optional`，并添加 Description 说明"关联安全组有序列表"。

## 2. Service 层封装新增方法

- [x] 2.1 在 `tencentcloud/services/ckafka/service_tencentcloud_ckafka.go` 中新增 `DescribeSecurityGroupRoutes` 封装方法：接收 `ctx, instanceId, routeId`，构造 `DescribeSecurityGroupRoutesRequest`（设置 `InstanceRoute.InstanceId` 与 `InstanceRoute.RouteId`），调用 SDK `DescribeSecurityGroupRoutes`，在返回结果中遍历 `SecurityGroupRoutes` 列表匹配 `RouteId` 一致的项并返回；带 ratelimit.Check 与日志，使用 `tccommon.ReadRetryTimeout` 的 retry 包装。
- [x] 2.2 在 `tencentcloud/services/ckafka/service_tencentcloud_ckafka.go` 中新增 `ModifyRouteSecurityGroups` 封装方法：接收 `ctx, instanceId, routeId, securityGroupIds`，构造 `ModifyRouteSecurityGroupsRequest`（设置 `InstanceRoute` 与 `SecurityGroupIds`），调用 SDK `ModifyRouteSecurityGroups`，带 ratelimit.Check 与日志，使用 `tccommon.WriteRetryTimeout` 的 retry 包装。

## 3. CRUD 逻辑修改

- [x] 3.1 在 `resourceTencentCloudCkafkaRouteCreate` 中，当用户指定 `security_group_ids` 时，将 `d.Get("security_group_ids").([]interface{})` 转为 `[]*string` 并赋值给 `request.SecurityGroupIds`。
- [x] 3.2 在 `resourceTencentCloudCkafkaRouteRead` 中，调用服务层 `DescribeSecurityGroupRoutes` 方法查询路由关联的安全组列表，若返回结果非空且 `SecurityGroupIds` 非空，则 `d.Set("security_group_ids", ...)`；若为空则设置为空列表。注意在读取前先用 `log.Printf("[CRUD] ckafka_route id=%s", d.Id())` 保留现场。
- [x] 3.3 在 `resourceTencentCloudCkafkaRouteUpdate` 中，将 `security_group_ids` 从 `immutableArgs` 数组中移除；新增逻辑：当 `d.HasChange("security_group_ids")` 时，从 `d.Id()` 解析 instanceId 与 routeId，调用服务层 `ModifyRouteSecurityGroups` 方法进行修改，修改成功后执行 read 同步状态。

## 4. 单元测试补充

- [x] 4.1 在 `tencentcloud/services/ckafka/resource_tc_ckafka_route_test.go` 中补充 `security_group_ids` 字段的单元测试用例，使用 terraform 测试套件（现有资源修改），新增包含 `security_group_ids` 的 HCL 配置及 TestStep。

## 5. 文档与收尾

- [x] 5.1 更新 `tencentcloud/services/ckafka/resource_tc_ckafka_route.md` 示例，补充 `security_group_ids` 参数的使用示例。
- [x] 5.2 在收尾阶段通过 `make doc` 命令自动生成 `website/docs/` 下的 markdown 文档。