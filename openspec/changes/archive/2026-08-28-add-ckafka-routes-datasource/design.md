## Context

当前 Terraform Provider for TencentCloud 已支持 CKafka 路由资源的管理（`tencentcloud_ckafka_route`），涵盖创建、读取、更新、删除等操作，服务层 `service_tencentcloud_ckafka.go` 中已存在 `DescribeCkafkaRouteById` 方法用于按 `instanceId + routeId` 查询单条路由。

但是该方法是按单条路由 ID 查询的，且未暴露为数据源。用户需要一个 `RESOURCE_KIND_DATASOURCE` 类型数据源 `tencentcloud_ckafka_routes`，以便通过 `DescribeRoute` 接口查询实例下所有路由信息（或按 `routeId` / `main_route_flag` 过滤），并在 Terraform 配置中引用路由详情。

云 API `DescribeRoute` 为同步接口（查看路由信息），入参为 `InstanceId`（必填）、`RouteId`（可选，int64）、`MainRouteFlag`（可选，bool），出参为 `Response.Result`（`RouteResponse` 类型，内含 `Routers` 列表）。每个 `Route` 包含接入方式、路由 ID、网络类型、VIP 列表、域名、子网、VPC、状态等字段；`VipList` 与 `BrokerVipList` 均为 `VipEntity` 列表（`Vip` string、`Vport` string）。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_ckafka_routes` 数据源，封装 `DescribeRoute` 接口
- 支持按 `instance_id`（必填）、`route_id`（可选）、`main_route_flag`（可选）查询路由
- 将路由列表展开平铺到 schema 顶层，每个路由元素包含其全部字段（含嵌套的 `vip_list`、`broker_vip_list`）
- 正确处理 API 返回为空的情况：retry 块内返回 `NonRetryableError`，避免清空 state
- 在 `provider.go` 中注册新数据源
- 提供对应的 `.md` 文档与 `_test.go` 单元测试（使用 gomonkey mock 云 API）

**Non-Goals:**
- 不修改已有的 `tencentcloud_ckafka_route` 资源行为
- 不实现路由的创建/更新/删除（这些由资源承担）
- 不暴露分页参数（`DescribeRoute` 接口本身无分页字段）

## Decisions

### 1. 数据源命名与文件组织

采用 `tencentcloud_ckafka_routes`（复数），与现有数据源命名风格一致（如 `tencentcloud_ckafka_instances`、`tencentcloud_ckafka_acls`）。文件命名为 `data_source_tc_ckafka_routes.go`，遵循项目规范。

### 2. 服务层方法设计

现有 `DescribeCkafkaRouteById` 仅支持 `instanceId + routeId` 两个参数，缺少 `main_route_flag` 支持，且返回单条 `*ckafka.Route`。为满足数据源需求，新增服务层方法 `DescribeCkafkaRouteByFilter`，接收 `instanceId`、`routeId`、`mainRouteFlag` 参数，返回 `*ckafka.RouteResponse`（即 `Response.Result`），使数据源能够获取完整路由列表并支持全部三个入参。

**为什么不复用 `DescribeCkafkaRouteById`**：该方法签名固定为 `routeId int64`，无法表达"不传 routeId"和"传 mainRouteFlag"的语义；且返回单条 Route 会丢失列表语义。新增独立方法更清晰，也不影响现有资源逻辑。

### 3. Schema 结构

根据项目约束，数据源列表型数据禁止再嵌套一层 `xxx_set`/`xxx_list` 包裹结构，应将列表展开。因此：
- 顶层输出字段为 `routers`（`TypeList`），每个元素为 `schema.Resource`，包含平铺的路由字段
- `vip_list`、`broker_vip_list` 为 `TypeList`，元素为含 `vip`、`vport` 的 `schema.Resource`
- `result` 字段映射 `Response.Result`，由于 `RouteResponse` 仅含 `Routers` 一个字段，为避免冗余嵌套，不单独设置 `result` 包装层，直接使用 `routers` 列表承载所有路由信息

### 4. 数据源 ID 设置

数据源为只读查询，使用 `helper.BuildToken()` 生成随机 ID（参考 `tencentcloud_igtm_instance_list` 模式），而非使用 `instance_id`，因为同一实例可能被多次查询。

### 5. retry 与空响应处理

在 Read 函数的 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 块内调用服务层方法。若云 API 返回空（`response == nil` / `response.Response == nil` / `response.Response.Result == nil`），直接返回 `NonRetryableError`，不执行 `d.SetId("")`，避免因 API 短暂波动清空 state。在 retry 失败路径保留 `log.Printf("[DATASOURCE] read empty, skip SetId")` 日志。

### 6. 字段类型映射

- `route_id`（int64）→ `schema.TypeInt`
- `access_type`（int64）→ `schema.TypeInt`
- `vip_type`（int64）→ `schema.TypeInt`
- `domain_port`（int64）→ `schema.TypeInt`
- `status`（int64）→ `schema.TypeInt`
- `domain`、`delete_timestamp`、`subnet`、`vpc_id`、`note`（string）→ `schema.TypeString`
- `vip_list`、`broker_vip_list`（`[]*VipEntity`）→ `schema.TypeList`，元素含 `vip`（string）、`vport`（string）

### 7. 测试策略

按项目要求，新增资源使用 gomonkey mock 云 API 进行单元测试，不使用 Terraform 验收测试套件。mock `DescribeRoute` 的返回值，验证 Read 函数对响应的字段映射与 `d.Set` 逻辑。

## Risks / Trade-offs

- **[API 返回字段可能为 null]** → 部分字段（如 `domain`、`subnet`、`vpc_id`、`note`、`status`）云 API 注释标注可能返回 null。在 `d.Set` 前逐一判断 nil，避免空指针。
- **[复用现有服务层方法 vs 新增方法]** → 选择新增 `DescribeCkafkaRouteByFilter` 以完整支持三入参并返回列表，代价是多一个服务层方法，但保持了与现有 `DescribeCkafkaRouteById` 的隔离，降低回归风险。
- **[无分页]** → `DescribeRoute` 接口无分页参数，路由数量由实例规格决定，通常不会过大，不存在分页缺失风险。
