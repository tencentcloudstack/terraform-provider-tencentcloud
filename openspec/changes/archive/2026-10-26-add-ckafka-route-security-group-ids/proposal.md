## Why

`tencentcloud_ckafka_route` 资源当前不支持关联安全组，用户无法通过 Terraform 在创建 ckafka 路由时指定安全组，也无法在后续修改路由关联的安全组。需要新增 `security_group_ids` 参数，使路由资源能够管理其关联的安全组列表，与云 API 的能力对齐。

## What Changes

- 在 `tencentcloud_ckafka_route` 资源 schema 中新增 Optional 参数 `security_group_ids`（类型为 `TypeList` of `TypeString`），用于指定路由关联的安全组有序列表。
- 在资源 Create 方法中，将 `security_group_ids` 传入 `CreateRoute` 接口的 `SecurityGroupIds` 入参。
- 在资源 Read 方法中，新增调用 `DescribeSecurityGroupRoutes` 接口查询路由关联的安全组列表，并将结果 set 到 `security_group_ids` 字段。
- 在资源 Update 方法中，新增调用 `ModifyRouteSecurityGroups` 接口，当 `security_group_ids` 发生变更时对路由关联的安全组进行修改；同时将其从 `immutableArgs` 中排除。
- 在资源对应的服务层 `service_tencentcloud_ckafka.go` 中新增 `DescribeSecurityGroupRoutes` 与 `ModifyRouteSecurityGroups` 的封装方法。
- 在资源单元测试文件 `resource_tc_ckafka_route_test.go` 中补充 `security_group_ids` 字段相关的测试用例。

## Capabilities

### New Capabilities

- `ckafka-route-security-groups`: 为 `tencentcloud_ckafka_route` 资源提供安全组关联管理能力，涵盖创建时指定、读取时同步、更新时修改安全组列表。

### Modified Capabilities

### Impact

- 代码文件：
  - `tencentcloud/services/ckafka/resource_tc_ckafka_route.go`：新增 schema 字段及 CRUD 逻辑。
  - `tencentcloud/services/ckafka/resource_tc_ckafka_route_test.go`：补充单元测试用例。
  - `tencentcloud/services/ckafka/service_tencentcloud_ckafka.go`：新增服务层封装方法。
  - `tencentcloud/services/ckafka/resource_tc_ckafka_route.md`：文档示例更新（由收尾阶段 `make doc` 生成）。
- 云 API 依赖（均已在 vendor 中存在）：
  - `CreateRoute`（`request.SecurityGroupIds []*string`）
  - `DescribeSecurityGroupRoutes`（response 中 `Result.SecurityGroupRoutes[].SecurityGroupIds`）
  - `ModifyRouteSecurityGroups`（`request.InstanceRoute` + `request.SecurityGroupIds []*string`）
- 兼容性：`security_group_ids` 为 Optional 字段，不填写时行为与原先一致，向后兼容。