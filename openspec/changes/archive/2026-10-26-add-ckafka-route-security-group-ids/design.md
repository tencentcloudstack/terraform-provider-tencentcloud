## Context

`tencentcloud_ckafka_route` 资源（位于 `tencentcloud/services/ckafka/resource_tc_ckafka_route.go`）用于管理 ckafka 实例的路由。当前资源 schema 包含 `instance_id`、`vip_type`、`vpc_id`、`subnet_id`、`access_type`、`auth_flag`、`caller_appid`、`public_network`、`ip` 以及 computed 字段 `vip_list`、`broker_vip_list`。

资源使用复合 ID：`instanceId + FILED_SP + routeId`。Create 调用 `CreateRoute`，Read 调用 `DescribeRoute`（通过服务层 `DescribeCkafkaRouteById`），Delete 调用 `DeleteRoute`。当前 Update 方法中所有顶层字段均为 immutable（放入 `immutableArgs` 数组，变更即报错）。

腾讯云 ckafka SDK (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka/v20190819`, vendor) 已支持安全组相关的三个接口：
- `CreateRoute`：入参已包含 `SecurityGroupIds []*string`
- `DescribeSecurityGroupRoutes`：用于查询路由关联的安全组
- `ModifyRouteSecurityGroups`：用于修改路由关联的安全组

本次变更新增 `security_group_ids` 参数，使路由资源能够管理其关联的安全组。

### 云 API request/response struct 定义（从 vendor 摘录）

**CreateRouteRequest**（models.go 行 2807）— 新增入参 `SecurityGroupIds`：
```go
type CreateRouteRequest struct {
    *tchttp.BaseRequest
    InstanceId    *string  `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`
    VipType       *int64   `json:"VipType,omitnil,omitempty" name:"VipType"`
    VpcId         *string  `json:"VpcId,omitnil,omitempty" name:"VpcId"`
    SubnetId      *string  `json:"SubnetId,omitnil,omitempty" name:"SubnetId"`
    AccessType    *int64   `json:"AccessType,omitnil,omitempty" name:"AccessType"`
    AuthFlag      *int64   `json:"AuthFlag,omitnil,omitempty" name:"AuthFlag"`
    CallerAppid   *int64   `json:"CallerAppid,omitnil,omitempty" name:"CallerAppid"`
    PublicNetwork *int64   `json:"PublicNetwork,omitnil,omitempty" name:"PublicNetwork"`
    Ip            *string  `json:"Ip,omitnil,omitempty" name:"Ip"`
    Note          *string  `json:"Note,omitnil,omitempty" name:"Note"`
    // 关联安全组有序列表
    SecurityGroupIds []*string `json:"SecurityGroupIds,omitnil,omitempty" name:"SecurityGroupIds"`
    IpWhitelist      []*IpWhitelistDTO `json:"IpWhitelist,omitnil,omitempty" name:"IpWhitelist"`
}
```

**DescribeSecurityGroupRoutesRequest**（models.go 行 6766）：
```go
type DescribeSecurityGroupRoutesRequest struct {
    *tchttp.BaseRequest
    // 路由信息（InstanceId + RouteId）
    InstanceRoute *InstanceRoute `json:"InstanceRoute,omitnil,omitempty" name:"InstanceRoute"`
    // 过滤器
    Filters []*RouteFilter `json:"Filters,omitnil,omitempty" name:"Filters"`
    // 分页Offset,默认0
    Offset *int64 `json:"Offset,omitnil,omitempty" name:"Offset"`
    // 分页Limit,默认20
    Limit *int64 `json:"Limit,omitnil,omitempty" name:"Limit"`
    // 关键词,可根据实例id/实例名称/vip模糊搜索
    SearchWord *string `json:"SearchWord,omitnil,omitempty" name:"SearchWord"`
}
```

**InstanceRoute**（models.go 行 9360）：
```go
type InstanceRoute struct {
    // ckafka集群实例Id
    InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`
    // 路由Id
    RouteId *int64 `json:"RouteId,omitnil,omitempty" name:"RouteId"`
}
```

**DescribeSecurityGroupRoutesResponse**（models.go 行 6817）：
```go
type DescribeSecurityGroupRoutesResponse struct {
    *tchttp.BaseResponse
    Response *DescribeSecurityGroupRoutesResponseParams `json:"Response"`
}
type DescribeSecurityGroupRoutesResponseParams struct {
    // 返回的安全组路由信息结果对象
    Result *SecurityGroupRouteResp `json:"Result,omitnil,omitempty" name:"Result"`
    RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}
```

**SecurityGroupRouteResp**（models.go 行 12191）：
```go
type SecurityGroupRouteResp struct {
    // 符合条件的安全组路由信息总数
    TotalCount *int64 `json:"TotalCount,omitnil,omitempty" name:"TotalCount"`
    // 符合条件的安全组路由信息列表
    SecurityGroupRoutes []*SecurityGroupRoute `json:"SecurityGroupRoutes,omitnil,omitempty" name:"SecurityGroupRoutes"`
}
```

**SecurityGroupRoute**（models.go 行 12164）：
```go
type SecurityGroupRoute struct {
    // 路由信息
    InstanceRoute *InstanceRoute `json:"InstanceRoute,omitnil,omitempty" name:"InstanceRoute"`
    // 关联的安全组列表
    SecurityGroupIds []*string `json:"SecurityGroupIds,omitnil,omitempty" name:"SecurityGroupIds"`
    // ckafka集群实例名称
    InstanceName *string `json:"InstanceName,omitnil,omitempty" name:"InstanceName"`
    // 路由vpcId
    VpcId *string `json:"VpcId,omitnil,omitempty" name:"VpcId"`
    // 路由vip
    Vip *string `json:"Vip,omitnil,omitempty" name:"Vip"`
}
```

**ModifyRouteSecurityGroupsRequest**（models.go 行 10717）：
```go
type ModifyRouteSecurityGroupsRequest struct {
    *tchttp.BaseRequest
    // 实例路由
    InstanceRoute *InstanceRoute `json:"InstanceRoute,omitnil,omitempty" name:"InstanceRoute"`
    // 修改后的安全组有序列表。
    // 注意:不指定此参数或传空列表则代表解绑所有关联的安全组。
    SecurityGroupIds []*string `json:"SecurityGroupIds,omitnil,omitempty" name:"SecurityGroupIds"`
}
```

**ModifyRouteSecurityGroupsResponse**（models.go 行 10757）：
```go
type ModifyRouteSecurityGroupsResponse struct {
    *tchttp.BaseResponse
    Response *ModifyRouteSecurityGroupsResponseParams `json:"Response"`
}
type ModifyRouteSecurityGroupsResponseParams struct {
    // 返回结果
    Result *SecurityGroupRouteOperateResp `json:"Result,omitnil,omitempty" name:"Result"`
    RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}
```

**SecurityGroupRouteOperateResp**（models.go 行 12181）：
```go
type SecurityGroupRouteOperateResp struct {
    // 操作返回的code，0为正常，非0为错误
    ReturnCode *string `json:"ReturnCode,omitnil,omitempty" name:"ReturnCode"`
    // 操作返回的信息
    ReturnMessage *string `json:"ReturnMessage,omitnil,omitempty" name:"ReturnMessage"`
}
```

**SDK Client 方法**（vendor client.go）：
- `CreateRoute(request *CreateRouteRequest) (*CreateRouteResponse, error)` — 行 1489
- `DescribeSecurityGroupRoutes(request *DescribeSecurityGroupRoutesRequest) (*DescribeSecurityGroupRoutesResponse, error)` — 行 4643
- `DescribeSecurityGroupRoutesWithContext(ctx, request)` — 行 4657
- `ModifyRouteSecurityGroups(request *ModifyRouteSecurityGroupsRequest) (*ModifyRouteSecurityGroupsResponse, error)` — 行 6605
- `ModifyRouteSecurityGroupsWithContext(ctx, request)` — 行 6619

### 目标资源现有 schema 与新增参数关联结构

现有 schema 字段（`resource_tc_ckafka_route.go` 行 29-129）：
- `instance_id`（Required, TypeString）
- `vip_type`（Required, TypeInt）
- `vpc_id`（Optional+Computed, TypeString）
- `subnet_id`（Optional+Computed, TypeString）
- `access_type`（Optional+Computed, TypeInt）
- `auth_flag`（Optional, TypeInt）
- `caller_appid`（Optional, TypeInt）
- `public_network`（Optional, TypeInt）
- `ip`（Optional+Computed, TypeString）
- `vip_list`（Computed, TypeList）
- `broker_vip_list`（Computed, TypeList）

当前 Update 方法的 `immutableArgs` 数组为：
```go
immutableArgs := []string{"instance_id", "vip_type", "vpc_id", "subnet_id", "access_type", "auth_flag", "caller_appid", "public_network", "ip"}
```

新增参数 `security_group_ids`：
- 类型：`schema.TypeList`，元素类型 `schema.TypeString`
- 属性：`Optional`
- 不放入 `immutableArgs`，支持 Update

资源 ID 格式：`instanceId + FILED_SP + routeId`，Read/Update/Delete 中已通过 `strings.Split(d.Id(), FILED_SP)` 解析出 `instanceId` 与 `routeId`，可直接复用。

## Goals / Non-Goals

**Goals:**
- 使 `tencentcloud_ckafka_route` 资源支持创建时指定安全组、读取时同步安全组、更新时修改安全组。
- 通过 `CreateRoute`、`DescribeSecurityGroupRoutes`、`ModifyRouteSecurityGroups` 三个已存在的云 API 接口实现。
- 保持向后兼容：`security_group_ids` 为 Optional 字段，不指定时行为不变。

**Non-Goals:**
- 不修改路由其他已有字段的 schema 属性（如 Required/Optional）。
- 不引入新的异步等待逻辑（ModifyRouteSecurityGroups 为同步接口）。
- 不修改 `provider.go` / `provider.md` 中资源注册（资源已注册）。

## Decisions

### 决策 1：`security_group_ids` schema 类型使用 TypeList

**选择**：`schema.TypeList` + `schema.TypeString` 元素。

**理由**：云 API 中 `SecurityGroupIds` 为 `[]*string`，且文档说明为"有序列表"。使用 `TypeList` 可保留顺序信息，与云 API 语义一致。

### 决策 2：Read 方法中安全组数据来源

**选择**：新增调用 `DescribeSecurityGroupRoutes` 接口（而非复用 `DescribeRoute`）。

**理由**：现有 Read 方法通过 `DescribeRoute` 接口获取的 `Route` 对象中不包含安全组信息。`DescribeSecurityGroupRoutes` 接口专门返回 `SecurityGroupRoutes` 列表，其中每个元素包含 `InstanceRoute`（InstanceId + RouteId）与 `SecurityGroupIds`，可直接通过 instanceId + routeId 匹配到当前路由的安全组。使用 `InstanceRoute` 入参配合精确查询。

### 决策 3：Update 方法中安全组修改方式

**选择**：当 `d.HasChange("security_group_ids")` 时，调用 `ModifyRouteSecurityGroups` 接口，将 `security_group_ids` 从 `immutableArgs` 中移除。

**理由**：`ModifyRouteSecurityGroups` 接口接收 `InstanceRoute`（InstanceId + RouteId）和新的 `SecurityGroupIds` 列表，符合更新场景。接口为同步接口，无需轮询。传入空列表代表解绑所有安全组（与云 API 语义一致）。

### 决策 4：服务层封装方法

**选择**：在 `service_tencentcloud_ckafka.go` 中新增两个封装方法：
- `DescribeSecurityGroupRoutes(ctx, instanceId, routeId)` → 返回 `*ckafka.SecurityGroupRoute`（匹配的第一项）
- `ModifyRouteSecurityGroups(ctx, instanceId, routeId, securityGroupIds)` → 返回 error

**理由**：遵循现有代码风格（服务层封装 API 调用 + ratelimit.Check + 日志），资源层调用服务层方法。

### 决策 5：Create 方法中的处理

**选择**：在 `CreateRoute` 请求中，当用户指定了 `security_group_ids` 时，将 `[]*string` 传入 `request.SecurityGroupIds`。

**理由**：`CreateRoute` 接口入参已支持 `SecurityGroupIds`，创建时一次性关联安全组，避免创建后立即调用 Modify。

## Risks / Trade-offs

- **[Risk] DescribeSecurityGroupRoutes 返回多条记录需精确匹配** → **Mitigation**：在服务层封装中通过 `InstanceRoute.InstanceId` 与 `InstanceRoute.RouteId` 入参查询，并在结果中遍历匹配 `RouteId` 一致的路由项，确保拿到正确的安全组列表；若未匹配到则不设置该字段（保持空）。
- **[Risk] ModifyRouteSecurityGroups 传空列表代表解绑全部安全组** → **Mitigation**：这是云 API 的既定行为，用户显式传入空列表即代表解绑，符合预期，无需特殊处理。
- **[Risk] 向后兼容性** → **Mitigation**：`security_group_ids` 为 Optional 字段，不填写时 `request.SecurityGroupIds` 为 nil，云 API 行为与原先一致。