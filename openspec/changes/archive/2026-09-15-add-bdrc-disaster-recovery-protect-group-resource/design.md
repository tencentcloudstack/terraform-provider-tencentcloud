## Context

BDRC（Business Disaster Recovery Center，业务灾备中心）是腾讯云容灾产品，`v20260330` 版本 SDK 提供了容灾保护组的完整 CRUD 接口：

- `CreateDisasterRecoveryProtectGroup`：入参 `SitePairId` / `ProtectGroupType` / `RecoveryPointObjective`（必填）、`ProtectGroupName` / `DataDirection`（可选）；出参 `ProtectGroupId`。
- `DescribeDisasterRecoveryProtectGroups`：入参 `ProtectGroupType`（必填）、`ProtectGroupIds` / `Filters` / `Order` / `OrderField` / `Offset` / `Limit`（可选）；出参 `TotalCount` + `ProtectGroupSet []*ProtectGroup`。`ProtectGroup` 结构含 25 个字段（含嵌套的 `ProtectedResourceStatusSet`）。
- `ModifyProtectGroupAttribute`：入参 `ProtectGroupId`（必填）、`ProtectGroupName`（可选）——**仅支持改名**。
- `DeleteDisasterRecoveryProtectGroups`：入参 `ProtectGroups []*string`（必填，ID 列表）。

本 change 是 BDRC 在 Terraform Provider 中的**首次接入**，因此除资源本身外，还需新增客户端方法（`UseBdrcV20260330Client`）与服务层（`BdrcService`）。

接口语义关键点：

- 创建接口返回 `ProtectGroupId`，单字段即可作为资源 ID，无需复合分隔符。
- 修改接口仅暴露 `ProtectGroupName`，意味着 `SitePairId` / `ProtectGroupType` / `RecoveryPointObjective` / `DataDirection` 一旦创建即不可变，必须 `ForceNew`，并在 Update 中对这几个字段做 `immutableArgs` 校验。
- 查询接口是列表型，按 `ProtectGroupIds` 过滤可定位单个保护组；返回的 `ProtectGroupSet` 是数组，取第一项即可（因为按单个 ID 查询）。
- 查询接口返回的 `RecoveryPointObjective` 单位是**秒**，而创建接口入参单位是**分钟**；为避免单位转换歧义，schema 中 `recovery_point_objective` 类型为 `TypeInt`，创建时直接传入，Read 时直接回写（按 SDK 实际返回值回填，不做人为换算）。

参考：本 change 的代码风格严格对齐 `tencentcloud_igtm_strategy` —— 一个标准通用型资源，在 Provider 内已有完整 CRUD 样板（Create 带 retry + 空指针保护、Read 用 service 层查询、Update 用 `immutableArgs` 校验可变字段、Delete 带 retry）。

## Goals / Non-Goals

**Goals:**

- Schema 入参字段名与 `CreateDisasterRecoveryProtectGroup` 接口入参 1:1 映射（snake_case 化）：`site_pair_id` / `protect_group_type` / `recovery_point_objective` / `protect_group_name` / `data_direction`。
- Schema 只读 Computed 字段覆盖 `DescribeDisasterRecoveryProtectGroups` 返回的 `ProtectGroup` 全部字段（展开平铺到顶层，不引入 `protect_group_set` 嵌套层）。
- 资源 ID = `ProtectGroupId`，单字段，支持 import（`schema.ImportStatePassthrough`）。
- `site_pair_id` / `protect_group_type` / `recovery_point_objective` / `data_direction` 标记 `ForceNew`；Update 仅处理 `protect_group_name` 变更，其余顶层入参字段列入 `immutableArgs`，检测到变更返回 error。
- 新增 `UseBdrcV20260330Client()` 客户端方法 + `bdrcv20260330Conn` 字段到 `tencentcloud/connectivity/client.go`。
- 新增 `BdrcService` 服务层，封装 `DescribeDisasterRecoveryProtectGroupById(ctx, protectGroupId, protectGroupType)` 查询方法，内部按 `Limit=100` 分页直至找到目标 ID。
- 全部 SDK 调用包裹 `resource.Retry(...)`（Read 用 `tccommon.ReadRetryTimeout`，Create/Update/Delete 用 `tccommon.WriteRetryTimeout`），错误经 `tccommon.RetryError(e)` 包装。
- 全部接口返回值做空指针保护（`result == nil || result.Response == nil`）。
- 资源 doc 命名 `resource_tc_bdrc_disaster_recovery_protect_group.md`，测试 `resource_tc_bdrc_disaster_recovery_protect_group_test.go`（gomonkey mock，不用 terraform 验收套件）。

**Non-Goals:**

- 不实现配套数据源 `tencentcloud_bdrc_disaster_recovery_protect_groups`（本次需求仅资源）。
- 不实现保护组下复制对 / 保护资源的管理（独立资源，后续 PR）。
- 不在 schema 对 `protect_group_type` / `data_direction` 做强枚举校验：API 文档枚举值后续可能扩展，不收紧。
- 不做 `recovery_point_objective` 的单位换算：创建入参与查询出参可能单位不同（分钟 vs 秒），schema 直接透传 SDK 值，不做人为转换，避免引入隐性 bug。
- 不修改任何既有资源/数据源/service 方法。

## Decisions

### D1 — Schema 字段映射（入参 + ForceNew）

| HCL 字段 | SDK 字段（Create） | 类型 | 必填 | ForceNew | Computed | Sensitive |
|---|---|---|---|---|---|---|
| `site_pair_id` | `SitePairId` | TypeString | Yes | **Yes** | No | No |
| `protect_group_type` | `ProtectGroupType` | TypeString | Yes | **Yes** | No | No |
| `recovery_point_objective` | `RecoveryPointObjective` | TypeInt | Yes | **Yes** | No | No |
| `protect_group_name` | `ProtectGroupName` | TypeString | No | No | No | No |
| `data_direction` | `DataDirection` | TypeString | No | **Yes** | No | No |

**理由**：
- `ModifyProtectGroupAttribute` 仅暴露 `ProtectGroupId` + `ProtectGroupName`，其余 4 个入参创建后不可变，必须 `ForceNew`。
- `protect_group_name` 可在 Update 中通过 `ModifyProtectGroupAttribute` 修改，故 No-ForceNew。
- `recovery_point_objective` 在 SDK 中是 `*int64`，HCL 用 `TypeInt`（按 igtm / 其他同类资源惯例），转换时用 `helper.Int64(v.(int))`。

### D2 — Schema 只读 Computed 字段（来自 Describe 出参，展开平铺）

| HCL 字段 | SDK 字段（ProtectGroup） | 类型 | Computed |
|---|---|---|---|
| `app_id` | `AppId` | TypeInt | Yes |
| `site_pair_name` | `SitePairName` | TypeString | Yes |
| `source_region` | `SourceRegion` | TypeString | Yes |
| `source_zone` | `SourceZone` | TypeString | Yes |
| `source_vpc` | `SourceVpc` | TypeString | Yes |
| `target_region` | `TargetRegion` | TypeString | Yes |
| `target_zone` | `TargetZone` | TypeString | Yes |
| `target_vpc` | `TargetVpc` | TypeString | Yes |
| `copy_type` | `CopyType` | TypeString | Yes |
| `disaster_recovery_type` | `DisasterRecoveryType` | TypeString | Yes |
| `peer_cloud_name` | `PeerCloudName` | TypeString | Yes |
| `create_from` | `CreateFrom` | TypeString | Yes |
| `life_state` | `LifeState` | TypeString | Yes |
| `account_uin` | `AccountUin` | TypeString | Yes |
| `sub_account_uin` | `SubAccountUin` | TypeString | Yes |
| `create_time` | `CreateTime` | TypeString | Yes |
| `modify_time` | `ModifyTime` | TypeString | Yes |
| `bind_protected_resource_count` | `BindProtectedResourceCount` | TypeInt | Yes |
| `error_recovery_point_objective_count` | `ErrorRecoveryPointObjectiveCount` | TypeInt | Yes |
| `protected_resource_status_set` | `ProtectedResourceStatusSet` | TypeList | Yes |

其中 `protected_resource_status_set` 是 `TypeList`，`Elem` 为 `schema.Resource`，内含 `status`（TypeString，Computed）和 `count`（TypeInt，Computed），对应 SDK 的 `ProtectedResourceStatus` 结构。

**理由**：按项目硬约束"资源参数 schema 中禁止创建'该资源列表型数据'这一层嵌套"，`ProtectGroupSet` 数组本身不暴露为 schema 字段，而是把数组中元素（单个 ProtectGroup）的字段平铺到顶层。因为按单个 ID 查询时 `ProtectGroupSet` 只有一项，直接取 `ProtectGroupSet[0]` 回填。

### D3 — 资源 ID = ProtectGroupId（单字段，支持 import）

```go
d.SetId(*response.Response.ProtectGroupId)
```

并声明 `Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough}`。import 时用户传入 `ProtectGroupId`，Read 用 `DescribeDisasterRecoveryProtectGroupById` 查询。由于查询接口需要 `ProtectGroupType`（必填），import 后首次 Read 在不知道 type 的情况下，服务层方法通过 `ProtectGroupIds` 过滤查询（type 传空会有 API 报错——需在 service 层处理：若 type 为空，尝试用空字符串或省略，由 API 返回；若 API 强制要求 type，则 import 不支持——见 Risks）。

**决策**：查询接口 `ProtectGroupType` 是必填，import 时仅有 ID 没有 type。因此 service 层 `DescribeDisasterRecoveryProtectGroupById` 在 `protectGroupType` 为空时，不传 `ProtectGroupType`（SDK 字段为 `*string`，nil 即不传），由后端按 `ProtectGroupIds` 过滤。若后端强制要求 type，则 import 会失败并提示用户通过 `terraform import` 时无法使用——但根据 SDK 字段定义（`*string` 带 `omitnil`），传 nil 不会报错，后端应支持仅按 ID 查询。Read 回填后会写入 `protect_group_type`，后续刷新即可正常。

### D4 — Create 实现

严格参考 `resourceTencentCloudIgtmStrategyCreate`：

1. `defer tccommon.LogElapsed(...)` + `defer tccommon.InconsistentCheck(...)`。
2. 构建 `var (logId, ctx, request, response)` 块。
3. 从 `d.GetOk(...)` 填充 5 个入参字段。
4. `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包装 `CreateDisasterRecoveryProtectGroupWithContext`；retry 内校验 `result == nil || result.Response == nil` → `NonRetryableError`；成功赋值 `response = result`。
5. retry 失败路径 `log.Printf("[CRITAL]%s create bdrc disaster_recovery_protect_group failed, reason:%+v", logId, reqErr)`。
6. 校验 `response.Response.ProtectGroupId == nil` → `fmt.Errorf("ProtectGroupId is nil.")`（打印 logId + d.Id 便于排障）。
7. `d.SetId(*response.Response.ProtectGroupId)`。
8. `return resourceTencentCloudBdrcDisasterRecoveryProtectGroupRead(d, meta)`。

### D5 — Read 实现（service 层 + 展开 + 空响应保护）

1. `defer tccommon.LogElapsed(...)` + `defer tccommon.InconsistentCheck(...)`。
2. 构建 `var (logId, ctx, service = BdrcService{client: ...})`。
3. `protectGroupId := d.Id()`；`protectGroupType` 从 `d.Get("protect_group_type")` 取（首次 Read 后已有值；import 首次为空）。
4. `respData, err := service.DescribeDisasterRecoveryProtectGroupById(ctx, protectGroupId, protectGroupType)`，错误经 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包装在 service 层内。
5. 若 `respData == nil`：先 `log.Printf("[CRUD] bdrc disaster_recovery_protect_group id=%s", d.Id())` 保留现场，再 `d.SetId("")`，`return nil`。
6. 逐字段判 nil 后 `d.Set(...)` 回填全部 Computed 字段 + 入参字段（`site_pair_id` / `protect_group_type` / `recovery_point_objective` / `protect_group_name` / `data_direction` 也从响应回填，保证 state 与后端一致）。
7. `protected_resource_status_set` 转 `[]map[string]interface{}` 回填。

**service 层 `DescribeDisasterRecoveryProtectGroupById` 设计**：
- 构建 `DescribeDisasterRecoveryProtectGroupsRequest`，设置 `ProtectGroupIds = []*string{&protectGroupId}`。
- 若 `protectGroupType != ""`，设置 `ProtectGroupType`。
- `Limit = helper.Int64(100)`（云 API 最大值），`Offset` 从 0 开始循环递增，直到 `len(ProtectGroupSet) < Limit` 或找到目标 ID。
- 在 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 内调用接口，遍历 `ProtectGroupSet` 匹配 `ProtectGroupId`，找到即返回该项指针；全部页查完未找到返回 `nil, nil`。
- retry 内校验 `result == nil || result.Response == nil` → `NonRetryableError`。

### D6 — Update 实现（immutableArgs 校验 + 仅改名）

严格参考 `resourceTencentCloudIgtmStrategyUpdate` 的 `mutableArgs` 模式，但本资源用 `immutableArgs`（因为可变字段只有一个）：

1. `defer` + `var (logId, ctx)`。
2. `immutableArgs := []string{"site_pair_id", "protect_group_type", "recovery_point_objective", "data_direction"}`。
3. 遍历 `immutableArgs`，若 `d.HasChange(v)` → `return fmt.Errorf("bdrc disaster_recovery_protect_group `%s` is immutable, cannot be updated, please recreate.", v)`。
4. 若 `d.HasChange("protect_group_name")`：
   - 构建 `ModifyProtectGroupAttributeRequest`，`ProtectGroupId = helper.String(d.Id())`，`ProtectGroupName` 从 `d.GetOk` 取。
   - `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包装调用；校验 `result == nil || result.Response == nil`。
   - 失败路径 `log.Printf("[CRITAL]%s update bdrc disaster_recovery_protect_group failed, reason:%+v", logId, reqErr)`。
5. `return resourceTencentCloudBdrcDisasterRecoveryProtectGroupRead(d, meta)`。

### D7 — Delete 实现

1. `defer` + `var (logId, ctx, request = NewDeleteDisasterRecoveryProtectGroupsRequest())`。
2. `protectGroupId := d.Id()`；`request.ProtectGroups = []*string{&protectGroupId}`。
3. `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包装调用；校验 `result == nil || result.Response == nil`。
4. 失败路径 `log.Printf("[CRITAL]%s delete bdrc disaster_recovery_protect_group failed, reason:%+v", logId, reqErr)`。
5. `return nil`。

### D8 — 客户端方法（connectivity/client.go）

在 `tencentcloud/connectivity/client.go` 新增：
- import：`bdrcv20260330 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bdrc/v20260330"`。
- struct 字段：`bdrcv20260330Conn *bdrcv20260330.Client`。
- 方法（仿 `UseIgtmV20231024Client`）：
```go
func (me *TencentCloudClient) UseBdrcV20260330Client() *bdrcv20260330.Client {
    if me.bdrcv20260330Conn != nil {
        return me.bdrcv20260330Conn
    }
    cpf := me.NewClientProfile(300)
    me.bdrcv20260330Conn, _ = bdrcv20260330.NewClient(me.Credential, me.Region, cpf)
    me.bdrcv20260330Conn.WithHttpTransport(&LogRoundTripper{})
    return me.bdrcv20260330Conn
}
```

### D9 — Provider 注册 + provider.md

- `tencentcloud/provider.go`：新增 `"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/bdrc"` import；在资源注册段新增 `"tencentcloud_bdrc_disaster_recovery_protect_group": bdrc.ResourceTencentCloudBdrcDisasterRecoveryProtectGroup()`。
- `tencentcloud/provider.md`：新增 `Business Disaster Recovery Center(BDRC)` 段（Data Source 留空或省略，Resource 段追加 `tencentcloud_bdrc_disaster_recovery_protect_group`），确保 gendoc 能扫描生成 website doc。

### D10 — 测试策略（gomonkey mock）

新增资源使用 gomonkey mock 云 API，不使用 terraform 验收测试套件：
- `resource_tc_bdrc_disaster_recovery_protect_group_test.go`，package `bdrc_test`。
- mock `UseBdrcV20260330Client` 返回的 client 的四个方法（`CreateDisasterRecoveryProtectGroupWithContext` / `DescribeDisasterRecoveryProtectGroupsWithContext` / `ModifyProtectGroupAttributeWithContext` / `DeleteDisasterRecoveryProtectGroupsWithContext`）。
- 测试 Create → Read → Update（改名）→ Delete 全流程的业务逻辑。

## Risks / Trade-offs

- **Risk**: `DescribeDisasterRecoveryProtectGroups` 的 `ProtectGroupType` 是必填，import 时仅有 ID 没有 type → **Mitigation**: service 层在 type 为空时省略该字段（SDK `*string` nil 即不传），依赖后端按 `ProtectGroupIds` 过滤；若后端强制要求 type 则 import 失败，文档中说明 import 需确保资源仍存在。Read 回填 type 后后续刷新正常。
- **Risk**: `RecoveryPointObjective` 创建入参单位是分钟、查询出参单位是秒，若直接透传可能导致 plan 永远 diff → **Mitigation**: schema `recovery_point_objective` 标 `ForceNew`，创建后不可变；Read 回填的是后端返回值（秒），若与 HCL 声明值（分钟）数值不同会触发 diff。但因 `ForceNew` 不会触发 in-place update，只会提示重建。文档 NOTE 中说明单位差异，建议用户以 API 文档为准。这是可接受的权衡，避免人为换算引入隐性 bug。
- **Trade-off**: 不实现配套数据源。本次需求仅资源，数据源后续单独 PR。
- **Trade-off**: 不对 `protect_group_type` / `data_direction` 做强枚举校验。API 枚举值后续可能扩展，强校验会限制未来可用性。
