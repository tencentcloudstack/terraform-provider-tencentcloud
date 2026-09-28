## Context

BDRC（Business Disaster Recovery Center，业务灾难恢复中心）产品的 SDK 包 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bdrc/v20260330` 已 vendor 到仓库，但 Terraform Provider 中尚无任何 bdrc 服务目录与客户端访问器。安全组映射（SecurityGroupMapping）用于在站点对（SitePair）下建立生产端安全组与容灾端安全组之间的对应关系，是 BDRC 实例级容灾切换时安全组规则同步的基础配置。

当前涉及的云 API（已通过 vendor 目录核实）：

- `CreateSecurityGroupMapping`：入参 `SrcSecurityGroupId`、`TargetSecurityGroupId`、`SitePairId`；**响应仅含 `RequestId`，不返回新建映射的 ID**，且该接口为异步接口。
- `DescribeSecurityGroupMappings`：入参 `SitePairId`、`Filters`（`[]*FilterModel`，`FilterModel` 含 `Name`、`Values`）、`Offset`、`Limit`（注释标注最大值 500）、`Order`、`OrderField`；响应含 `TotalCount` 与 `SecurityGroupMappingSet []*SecurityGroupMapping`，其中 `SecurityGroupMapping` 含 `SecurityGroupMappingId`、`SitePairId`、`SourceSecurityGroupId`、`TargetSecurityGroupId`、`LifeState`。
- `DeleteSecurityGroupMapping`：入参 `SitePairId`、`SecurityGroupMappingIds`（`[]*string`，列表批量删除）；**无更新接口**。

参考资料：
- 代码风格模板：`tencentcloud/services/igtm/resource_tc_igtm_strategy.go`（通用资源 CRUD 样式参考）。
- 客户端访问器样式：`tencentcloud/connectivity/client.go` 中 `UseBhV20230418Client()` / `UseIgtmV20231024Client()`。
- 联合 ID 分隔符：`tccommon.FILED_SP`（`#`，定义于 `tencentcloud/common/common.go`）。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_bdrc_security_group_mapping` 通用资源，完整支持 Create / Read / Update（不可变约束）/ Delete / Import。
- 正确处理 `CreateSecurityGroupMapping` 异步且无返回 ID 的特性：Create 后轮询 `DescribeSecurityGroupMappings` 直到按 `src-security-group-id`+`target-security-group-id` 过滤命中新建映射，取其 `SecurityGroupMappingId`。
- 建立 bdrc 服务目录与客户端访问器，为后续 bdrc 资源/数据源接入打好基础。
- 资源 ID 采用联合 ID `site_pair_id#security_group_mapping_id`，支持 import 并在文档中说明联合 ID 用法。
- 补充 gomonkey mock 云 API 的单元测试（不使用 terraform 测试套件）。

**Non-Goals:**
- 不实现安全组映射的"更新"调用（云 API 无 Update 接口），Update 函数仅做不可变校验并返回 error。
- 不新增 bdrc 数据源（本次仅资源）。
- 不处理 `LifeState` 的异常态自愈（仅在 Read 中回填该字段）。
- 不修改任何既有资源的 schema 与 state。

## Decisions

### Decision 1: 客户端访问器命名与初始化
新增 `UseBdrcV20260330Client()` 方法及 `bdrcv20260330Conn *bdrcv20260330.Client` 字段，import 别名 `bdrcv20260330`。沿用 `UseBhV20230418Client` 样式：懒加载、`NewClientProfile(300)`、`WithHttpTransport(&LogRoundTripper{})`。

**理由**：与现有带版本号的客户端访问器（bhv20230418、igtmv20231024、ga2v20250115）保持一致，便于后续多版本共存。

### Decision 2: 联合 ID 设计
资源 ID = `site_pair_id + FILED_SP + security_group_mapping_id`。原因：
- `DescribeSecurityGroupMappings` 必须传 `SitePairId`，`DeleteSecurityGroupMapping` 必须传 `SitePairId` + `SecurityGroupMappingIds`，二者都需要 `site_pair_id`，无法仅凭 `SecurityGroupMappingId` 完成查删。
- 资源支持 import，文档中说明需使用联合 ID（`terraform import tencentcloud_bdrc_security_group_mapping.xxx sitePairId#securityGroupMappingId`）。

Read/Update/Delete 中统一使用 `strings.Split(d.Id(), tccommon.FILED_SP)` 解析，长度不为 2 时返回 `fmt.Errorf("id is broken,%s", d.Id())`。

### Decision 3: Create 异步轮询策略
`CreateSecurityGroupMapping` 无返回 ID 且为异步接口。Create 流程：
1. 用 `resource.Retry(tccommon.WriteRetryTimeout)` 调用 `CreateSecurityGroupMapping`。
2. 调用成功后，使用 `helper.Retry()`（最终一致性重试）轮询 `DescribeSecurityGroupMappings`：
   - 入参 `SitePairId` = schema 中的 `site_pair_id`；
   - `Filters` 设置 `Name="src-security-group-id"`、`Values=[srcSecurityGroupId]`（并叠加 `target-security-group-id` 过滤以精确定位），`Limit` 取注释最大值 500；
   - 当命中 `SecurityGroupMappingSet` 中存在 `SourceSecurityGroupId` 与 `TargetSecurityGroupId` 均匹配的记录时，取其 `SecurityGroupMappingId`，结束轮询。
3. 轮询成功后设置 `d.SetId(sitePairId + FILED_SP + securityGroupMappingId)`，再调用 Read 回填。

**理由**：符合"异步接口调用后调用 Read 接口轮询直到生效"的规范；用源/目标安全组 ID 双重过滤精确定位新建映射，避免同站点对下多条映射时误取。

### Decision 4: Read 查询单条映射
Read 不直接暴露列表型 schema 层（遵循"禁止创建列表型数据层 schema"规则），而是把 `SecurityGroupMappingSet` 第一条匹配项的字段平铺到资源顶层 schema。服务层 `BdrcService.DescribeSecurityGroupMappingById(ctx, sitePairId, securityGroupMappingId)` 内部：
- 构建 `DescribeSecurityGroupMappingsRequest`，`SitePairId` 传入，`Filters` 用 `Name="security-group-mapping-id"` 如不支持则改为遍历返回集合按 ID 过滤；稳妥做法：先按 `src-security-group-id` 无法定位（因为 Read 时未必知道源 ID，源 ID 来自 state），故改为不设 `Filters` 或仅设 `SitePairId`，`Limit=500`，然后在结果集中按 `SecurityGroupMappingId` 精确匹配。
- 在 `resource.Retry(tccommon.ReadRetryTimeout)` 中调用，返回匹配的 `*SecurityGroupMapping`。

**Read 空返回处理**：若 `response == nil || response.Response == nil || len(SecurityGroupMappingSet)==0` 或未匹配到目标 ID，**不要**直接 `d.SetId("")`，而是先 `log.Printf("[CRUD] bdrc security_group_mapping id=%s", d.Id())` 保留现场，再 `d.SetId("")`。

### Decision 5: Update 不可变约束
云 API 无安全组映射更新接口，属于"只有 CRD 接口"资源。`id` 字段设为 `ForceNew`（通过联合 ID 自然实现，schema 中不单独定义 id 参数）。Update 函数中：
```go
immutableArgs := []string{"src_security_group_id", "target_security_group_id", "site_pair_id"}
for _, v := range immutableArgs {
    if d.HasChange(v) {
        return fmt.Errorf("argument `%s` cannot be changed, please delete and recreate", v)
    }
}
return resourceTencentCloudBdrcSecurityGroupMappingRead(d, meta)
```

### Decision 6: Delete 实现
`DeleteSecurityGroupMapping` 入参 `SecurityGroupMappingIds` 为列表（`[]*string`），将单个映射 ID 包装成单元素列表传入。`SitePairId` 从 `d.Id()` 解析。使用 `resource.Retry(tccommon.WriteRetryTimeout)` 包装，失败用 `tccommon.RetryError(e)` 包装。

### Decision 7: Schema 字段定义
顶层 schema（参数均平铺，无列表嵌套层）：

| SchemaName | Type | 属性 | 来源 API 字段 | 说明 |
|---|---|---|---|---|
| `site_pair_id` | TypeString | Required, ForceNew | `request.SitePairId` / `response.SitePairId` | 站点对 ID |
| `src_security_group_id` | TypeString | Required, ForceNew | `request.SrcSecurityGroupId` | 生产端安全组 ID |
| `target_security_group_id` | TypeString | Required, ForceNew | `request.TargetSecurityGroupId` | 容灾端安全组 ID |
| `security_group_mapping_id` | TypeString | Computed | `response.SecurityGroupMappingId` | 安全组映射 ID（云分配） |
| `source_security_group_id` | TypeString | Computed | `response.SourceSecurityGroupId` | 生产端安全组 ID（回填，与 `src_security_group_id` 对应） |
| `target_security_group_id` | 已存在于入参，Computed 复用 | `response.TargetSecurityGroupId` | 容灾端安全组 ID（回填） |
| `life_state` | TypeString | Computed | `response.LifeState` | 生命状态（NORMAL 等） |

**注意**：Create 入参字段名为 `src_security_group_id`（对应 API `SrcSecurityGroupId`），而 Read 出参字段名为 `source_security_group_id`（对应 API `SourceSecurityGroupId`）。二者为同一概念但 API 命名不一致，因此 schema 中分别保留：`src_security_group_id` 作为入参（Required/ForceNew），`source_security_group_id` 作为 Computed 回填出参。Read 时同时回填两者，保证 state 与 API 出参一致。

`filters` / `order` / `order_field` 等仅属于 `DescribeSecurityGroupMappings` 的查询入参，不属于单个资源的管理参数，**不**纳入资源 schema（它们属于数据源范畴，本次不涉及）。

### Decision 8: 单元测试策略
新增资源测试使用 gomonkey mock 云 API，不使用 terraform 验收测试套件。mock 点：
- `UseBdrcV20260330Client` 返回的 client 上的 `CreateSecurityGroupMappingWithContext` / `DescribeSecurityGroupMappingsWithContext` / `DeleteSecurityGroupMappingWithContext`。
- 覆盖 Create（含轮询命中）、Read（命中/未命中）、Update（不可变报错）、Delete 成功路径。

## Risks / Trade-offs

- **[异步创建定位精度]** → 用 `src-security-group-id` + `target-security-group-id` 双重 Filter 精确定位新建映射，避免同站点对下重复配置误取；轮询超时由 `helper.Retry` 控制，超时返回明确错误便于人工介入。
- **[Create 无返回 ID]** → 通过轮询 Describe 获取 ID 作为联合 ID 一部分；若 Describe 始终查不到，Create 以重试耗尽失败，不会写入空 ID（遵循"Create 必须检查返回值，空则 NonRetryableError"的规范）。
- **[API 入参与出参命名不一致]**（`SrcSecurityGroupId` vs `SourceSecurityGroupId`）→ schema 中分别保留 `src_security_group_id`（入参）与 `source_security_group_id`（Computed 出参），并在文档中说明，避免用户混淆。
- **[无 Update 接口]** → 所有业务字段不可变，变更需删除重建；Update 函数返回明确 error 提示，用户体验略差但符合云 API 能力边界。
- **[新增客户端访问器]** → 仅新增，不修改既有访问器，零回归风险。
