## ADDED Requirements

### Requirement: Resource schema includes security_group_ids parameter

`tencentcloud_ckafka_route` 资源 SHALL 新增一个 `security_group_ids` 参数，类型为字符串列表（`TypeList` of `TypeString`），属性为 `Optional`，用于指定路由关联的安全组有序列表。

#### Scenario: security_group_ids in schema

- **WHEN** `ResourceTencentCloudCkafkaRoute()` 被调用
- **THEN** 返回的 schema 中 MUST 包含 `security_group_ids` 字段，类型为 `schema.TypeList`，元素类型为 `schema.TypeString`，属性为 `Optional`

### Requirement: Create route with security group ids

资源 Create 方法 SHALL 在调用 `CreateRoute` 接口时，将用户指定的 `security_group_ids` 转换为 `[]*string` 传入 `request.SecurityGroupIds`。

#### Scenario: Create with security_group_ids specified

- **WHEN** 用户在 Terraform 配置中指定了 `security_group_ids`
- **THEN** Create 方法 MUST 将该列表转换为 `[]*string` 并赋值给 `request.SecurityGroupIds`，然后调用 `CreateRoute`

#### Scenario: Create without security_group_ids

- **WHEN** 用户未指定 `security_group_ids`
- **THEN** Create 方法 MUST 不设置 `request.SecurityGroupIds`（保持 nil），`CreateRoute` 调用行为与原先一致

### Requirement: Read route security group ids

资源 Read 方法 SHALL 调用 `DescribeSecurityGroupRoutes` 接口查询路由关联的安全组列表，并将结果 set 到 `security_group_ids` 字段。

#### Scenario: Read sets security_group_ids from DescribeSecurityGroupRoutes

- **WHEN** Read 方法被执行
- **THEN** 方法 MUST 调用 `DescribeSecurityGroupRoutes` 接口（通过 InstanceRoute 入参传入 instanceId 与 routeId），从返回结果中匹配当前路由的安全组列表，并 set 到 `d.Set("security_group_ids", ...)`

#### Scenario: Read with no security group association

- **WHEN** `DescribeSecurityGroupRoutes` 返回的结果中未匹配到当前路由的安全组信息或安全组列表为空
- **THEN** 方法 MUST 将 `security_group_ids` 设置为空列表，不报错

### Requirement: Update route security group ids

资源 Update 方法 SHALL 支持修改 `security_group_ids` 参数，当该参数发生变更时调用 `ModifyRouteSecurityGroups` 接口。`security_group_ids` MUST 从 `immutableArgs` 数组中排除。

#### Scenario: Update security_group_ids on change

- **WHEN** `security_group_ids` 在 Terraform 配置中发生变更（`d.HasChange("security_group_ids")` 为 true）
- **THEN** Update 方法 MUST 调用 `ModifyRouteSecurityGroups` 接口，传入 `InstanceRoute`（instanceId + routeId）与新的 `SecurityGroupIds` 列表

#### Scenario: Update with other immutable args changed

- **WHEN** `immutableArgs` 中的字段（不含 `security_group_ids`）发生变更
- **THEN** Update 方法 MUST 返回错误，提示该参数不可修改

#### Scenario: Update with empty security_group_ids

- **WHEN** 用户将 `security_group_ids` 变更为空列表
- **THEN** 调用 `ModifyRouteSecurityGroups` 时传入空列表，解绑所有关联的安全组

### Requirement: Service layer encapsulates security group APIs

服务层 `CkafkaService` SHALL 新增两个封装方法用于安全组路由的查询与修改。

#### Scenario: DescribeSecurityGroupRoutes encapsulated in service

- **WHEN** 资源 Read 需要查询路由安全组
- **THEN** 服务层 MUST 提供方法接收 instanceId 与 routeId，调用 `DescribeSecurityGroupRoutes` 接口，返回匹配的 `*ckafka.SecurityGroupRoute`

#### Scenario: ModifyRouteSecurityGroups encapsulated in service

- **WHEN** 资源 Update 需要修改路由安全组
- **THEN** 服务层 MUST 提供方法接收 instanceId、routeId 与 securityGroupIds，调用 `ModifyRouteSecurityGroups` 接口，返回 error