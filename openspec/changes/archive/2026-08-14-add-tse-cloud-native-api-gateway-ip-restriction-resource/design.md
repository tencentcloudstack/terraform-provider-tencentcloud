## Context

TSE（微服务引擎）云原生 API 网关支持对路由或服务绑定 IP 访问控制插件（白名单/黑名单）。相关云 API 位于 `tencentcloud-sdk-go/tencentcloud/tse/v20201207`：

- `CreateOrModifyCloudNativeAPIGatewayIPRestriction`：创建或编辑访问控制（upsert 语义，再次调用同一目标会覆盖旧配置）。
- `DescribeCloudNativeAPIGatewayIPRestriction`：查询访问控制，出参为 `Result`（`DescribeKongIpRestrictionResult`，含 SourceType/SourceId/Enabled/RestrictionType/AddressList）。
- `DeleteCloudNativeAPIGatewayIPRestriction`：删除访问控制。

关键约束：
- 三个接口均无独立资源 ID 返回，资源由 `GatewayId` + `SourceType` + `SourceId` 三元组唯一确定。
- Create/Modify 接口为 upsert：同一目标重复调用即覆盖，Create 与 Update 共用同一调用。
- Delete 接口入参仅需 GatewayId/SourceType/SourceId。
- Describe 出参 `Result` 中不含 GatewayId/SourceType（这两个由请求入参确定），故 Read 时这两个字段直接从复合 ID 拆分回填。

参考资源：本 change 代码风格严格对齐 `tencentcloud_igtm_strategy`（复合 ID + Create/Read/Update/Delete 标准实现）与 `tencentcloud_tse_cngw_certificate`（TSE 包内已有复合 ID 样板，使用 `UseTseClient()`）。

## Goals / Non-Goals

**Goals:**
- Schema 字段名与 `CreateOrModifyCloudNativeAPIGatewayIPRestriction` 接口入参 1:1 映射（snake_case 化）。
- 资源 ID = `gateway_id#source_type#source_id`，使用 `tccommon.FILED_SP`（"#"）分隔，与三元组身份键对齐。
- Create 与 Update 共用 `CreateOrModifyCloudNativeAPIGatewayIPRestrictionWithContext` 调用（API 为 upsert）。
- Read 调用 `DescribeCloudNativeAPIGatewayIPRestrictionWithContext`，空返回时按规则打印日志后清空 ID。
- Delete 调用 `DeleteCloudNativeAPIGatewayIPRestrictionWithContext`。
- 全部 SDK 调用包裹 `resource.Retry(tccommon.WriteRetryTimeout/ReadRetryTimeout, ...)`，错误经 `tccommon.RetryError` 包装。
- 所有接口返回值做空指针保护。
- 资源支持 import（ImportStatePassthrough）。
- 资源 doc 与单元测试（gomonkey mock 云 API）符合项目规范。

**Non-Goals:**
- 不实现配套数据源（无 List 接口）。
- 不在 schema 加 `source_type`/`restriction_type` 的强枚举校验：当前值集合可能扩展，不收紧。
- 不修改任何既有资源/数据源/service 方法。

## Decisions

### D1 — Schema 字段映射与 ForceNew 策略

| HCL 字段 | SDK 字段 | 类型 | 必填 | ForceNew | Computed |
|---|---|---|---|---|---|
| `gateway_id` | `GatewayId` | TypeString | Required | **Yes** | No |
| `source_type` | `SourceType` | TypeString | Required | **Yes** | No |
| `source_id` | `SourceId` | TypeString | Required | **Yes** | No |
| `enabled` | `Enabled` | TypeBool | Required | No | No |
| `restriction_type` | `RestrictionType` | TypeString | Required | No | No |
| `address_list` | `AddressList` | TypeList of TypeString | Required | No | No |

**理由**：
- `gateway_id`/`source_type`/`source_id` 三者共同唯一确定一个访问控制目标；改变任一即等于"换一个目标绑定 IP 限制"，必须 ForceNew（重建）。
- `enabled`/`restriction_type`/`address_list` 为同一目标上的可变配置，走 Update 覆盖式更新，No-ForceNew。
- `address_list` 用 `TypeList`（有序）而非 `TypeSet`（无序），因为 IP 白/黑名单顺序对用户有可读性意义，且 API 返回数组顺序应与配置一致。元素为 `TypeString`（cidr 或 ip）。

### D2 — 复合 ID

资源 ID = `strings.Join([]string{gateway_id, source_type, source_id}, tccommon.FILED_SP)`，即 `gateway_id#source_type#source_id`。

Read/Update/Delete 均从 `d.Id()` 按 `tccommon.FILED_SP` 拆分三段，校验 `len == 3`，否则返回 `"id is broken"` 错误。

**理由**：三元组是唯一身份键，无独立 ID；复用项目内 `igtm_strategy`/`cngw_certificate` 的复合 ID 模式。

### D3 — Create 走 "SetId + Update" 双跳

参考 `tencentcloud_waf_owasp_rule_status_config` 与 `bh_bind_device_account_kubeconfig` 的写法，Create 仅：
1. 读取 `gateway_id`/`source_type`/`source_id`，`d.SetId(...)`。
2. 转调 `resourceTencentCloudTseCloudNativeAPIGatewayIPRestrictionUpdate(d, meta)`。

Update 内含真正的 `CreateOrModifyCloudNativeAPIGatewayIPRestrictionWithContext` 调用。

**理由**：Create 与 Update 共用同一 SDK 接口（upsert），避免调用代码两份重复，与已有 upsert 型资源风格一致。

### D4 — Read 刷新 state

Read 调用 `DescribeCloudNativeAPIGatewayIPRestrictionWithContext`，入参 GatewayId/SourceType/SourceId 取自复合 ID 拆分。出参 `response.Response.Result`：
- 若 `result == nil || result.Response == nil || result.Response.Result == nil`，打印 `[CRUD]` 现场日志后 `d.SetId("")`，返回 nil。
- 否则按字段是否为 nil 回填 `enabled`/`restriction_type`/`address_list`；`gateway_id`/`source_type`/`source_id` 直接从复合 ID 回填。

### D5 — Delete 调用

Delete 调用 `DeleteCloudNativeAPIGatewayIPRestrictionWithContext`，入参 GatewayId/SourceType/SourceId 取自复合 ID 拆分。包裹 `resource.Retry(tccommon.WriteRetryTimeout, ...)`，校验返回值非空。

### D6 — 重试与空指针保护

- 所有 SDK 调用包裹 `resource.Retry(...)`，错误经 `tccommon.RetryError(e)` 转换。
- Create/Update/Delete 使用 `tccommon.WriteRetryTimeout`；Read 使用 `tccommon.ReadRetryTimeout`。
- retry 块内仅放 SDK 调用；`d.SetId`/字段回填/`d.Set` 等成功操作放在 retry 块外。
- retry 块内校验 `result == nil || result.Response == nil` → `resource.NonRetryableError`。
- 失败路径打印 `[CRITAL]` 日志（沿用项目既有拼写习惯）。

### D7 — Provider 注册与文档

- `tencentcloud/provider.go`：在 tse 资源注册段追加一行注册。
- `tencentcloud/provider.md`：在 TSE Resource 段追加资源名，供 gendoc 索引。
- `resource_tc_tse_cloud_native_api_gateway_ip_restriction.md`：一句话描述（带 TSE 产品名）+ Example Usage + Import（说明使用复合 ID）。
- `make doc` 生成 website 文档（收尾阶段执行）。

### D8 — 单元测试

- 新增资源使用 gomonkey mock 云 API，不使用 terraform 测试套件。
- mock `UseTseClient` 返回 `&tse.Client{}`，再 mock 三个 `...WithContext` 方法。
- 测试覆盖 Create、Read、Update、Delete 四个回调。

## Risks / Trade-offs

- **Risk**: `source_type` 的取值当前为 `route`/`service`，未来可能扩展 → 不做强枚举校验，避免破坏向前兼容。文档中以描述说明合法值。
- **Risk**: 并发管理同一复合 ID 的资源会发生 upsert 覆盖竞争 → 文档说明"同一网关的同一路由/服务同时只能由一份 HCL 管理 IP 限制"。
- **Trade-off**: `address_list` 用 TypeList 而非 TypeSet，顺序敏感；若用户在 HCL 中打乱顺序会触发不必要的 diff。但对白/黑名单场景，可读性优先，且 API 返回顺序可预期，可接受。
