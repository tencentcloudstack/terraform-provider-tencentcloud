## Context

腾讯云防火墙管控面 SDK 包为 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611`（go.mod 已声明 `fwm v1.3.194`，vendor 目录已存在）。该包提供策略下发的三个接口：

- `CreateStrategy`（创建：下发前区/后区规则组，返回策略 Id 列表）
- `DescribeStrategy`（按单个策略 Id 查询详情）
- `DeleteStrategy`（按策略 Id 列表删除）

本 change 新增通用资源 `tencentcloud_fwm_strategy`，并将 `CreateStrategy` 的出参 `response.StrategyIds` 映射为 Terraform Computed 参数 `strategy_ids`。这是本次唯一的新增参数。

代码风格严格参考 `tencentcloud_igtm_strategy`（`tencentcloud/services/igtm/resource_tc_igtm_strategy.go`）：相同的 Create/Read/Update/Delete 结构、retry 包装、空指针保护、复合 ID、`immutableArgs` 处理 CRD-only 资源的 update。

### SDK Request/Response struct 定义（摘录自 vendor，供实施阶段直接引用）

**CreateStrategyRequest**（vendor/.../fwm/v20250611/models.go:869-903）：
```go
type CreateStrategyRequestParams struct {
    Product        *string        `json:"Product,omitnil,omitempty" name:"Product"`        // 产品类型
    ReceiveAccount []*string      `json:"ReceiveAccount,omitnil,omitempty" name:"ReceiveAccount"` // 下发账号
    PreStrategy    []*StrategyReq `json:"PreStrategy,omitnil,omitempty" name:"PreStrategy"`    // 前区规则组
    PostStrategy   []*StrategyReq `json:"PostStrategy,omitnil,omitempty" name:"PostStrategy"`   // 后区规则组
    ReceiveGroup   []*string      `json:"ReceiveGroup,omitnil,omitempty" name:"ReceiveGroup"`   // 下发账号组
}
// CreateStrategyRequest 内嵌 *tchttp.BaseRequest，字段同上
```

**CreateStrategyResponse**（models.go:929-940）：
```go
type CreateStrategyResponseParams struct {
    StrategyIds []*string `json:"StrategyIds,omitnil,omitempty" name:"StrategyIds"` // 策略Id列表
    RequestId   *string   `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}
type CreateStrategyResponse struct {
    *tchttp.BaseResponse
    Response *CreateStrategyResponseParams `json:"Response"`
}
```

**StrategyReq**（models.go:4310-4316）：
```go
type StrategyReq struct {
    GroupId  *string `json:"GroupId,omitnil,omitempty" name:"GroupId"`   // 规则组Id
    Sequence *int64  `json:"Sequence,omitnil,omitempty" name:"Sequence"` // 优先级
}
```

**StrategyResp**（models.go:4318-4366）—— DescribeStrategy 返回的单条策略详情：
```go
type StrategyResp struct {
    StrategyId     *string           // 策略Id
    GroupId        *string           // 规则组Id
    GroupName      *string           // 规则组名称
    RuleCount      *int64            // 规则数
    RuleStatus     *int64            // 策略状态
    ReceiveAccount []*ReceiveAccount // 下发账号
    Sequence       *int64            // 优先级
    ErrMsg         *string           // 下发失败原因
    ErrorType      *string           // 下发失败原因类型
    UpdateTime     *string
    CreateBy       *string
    UpdateBy       *string
    ExecArea       *string           // 执行区域
    CreateName     *string
    UpdateName     *string
    CreateTime     *string
}
```

**DescribeStrategyRequest / Response**（models.go:2536-2579）：
```go
type DescribeStrategyRequestParams struct {
    StrategyId *string `json:"StrategyId,omitnil,omitempty" name:"StrategyId"` // 策略Id
}
type DescribeStrategyResponseParams struct {
    Strategy  *StrategyResp `json:"Strategy,omitnil,omitempty" name:"Strategy"`
    RequestId *string       `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}
```

**DeleteStrategyRequest / Response**（models.go:1326-1366）：
```go
type DeleteStrategyRequestParams struct {
    StrategyIds []*string `json:"StrategyIds,omitnil,omitempty" name:"StrategyIds"` // 策略Id列表
}
type DeleteStrategyResponseParams struct {
    RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}
```

**Client 方法**（client.go）：`CreateStrategyWithContext`、`DescribeStrategyWithContext`、`DeleteStrategyWithContext` 均已存在；签名形如 `(ctx, request) (response, err)`。

### 现有 schema 与本资源的关联说明

本资源是全新资源，无"现有 schema"。新增 schema 入参严格对齐 `CreateStrategyRequest` 字段；新增 Computed 出参 `strategy_ids` 对齐 `CreateStrategyResponse.StrategyIds`。由于该资源由出参 `StrategyIds` 生成身份键，且无 ModifyStrategy 接口，按 CRD-only 惯例处理 Update（见 D3）。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_fwm_strategy` 资源，CRUD 对接 `CreateStrategy` / `DescribeStrategy` / `DeleteStrategy`。
- 将 `CreateStrategy` 出参 `StrategyIds` 映射为 Computed 参数 `strategy_ids`（本次唯一新增参数）。
- 入参字段与 `CreateStrategyRequest` 1:1（snake_case）：`product`(Required)、`receive_account`、`pre_strategy`、`post_strategy`、`receive_group`。
- 复合 ID 用 `tccommon.FILED_SP` 拼接 `StrategyIds` 列表；支持 import（import 时传联合 id）。
- 全部 SDK 调用 `resource.Retry` 包装 + `tccommon.RetryError`；retry 块内仅调用 API，SetId/Set 放块外。
- 在 `connectivity/client.go` 新增 `UseFwmV20250611Client()` 访问器；在 `provider.go` / `provider.md` 注册。
- 测试用 gomonkey mock 云 API（新增资源，不用 TF 测试套件）。

**Non-Goals:**
- 不实现配套数据源（`tencentcloud_fwm_strategies`）。
- 不实现 `ModifyStrategy`（SDK 无此接口）；Update 仅做 immutableArgs 校验，不真正调 API 改字段。
- 不修改任何既有资源/数据源/service 方法。
- 不在 schema 加 `product` 强枚举校验。

## Decisions

### D1 — 参数 `strategy_ids` 映射

| HCL 字段 | SDK 字段(JsonPath) | 类型 | Required | ForceNew | Computed |
|---|---|---|---|---|---|
| `strategy_ids` | `response.StrategyIds` | `TypeList` of `TypeString` | No | No | **Yes** |

`strategy_ids` 来自 `CreateStrategyResponse.Response.StrategyIds`（`[]*string`）。Create 成功后 `d.Set("strategy_ids", helper.StringsToInterface(...))`；Read 中通过遍历 `StrategyIds`（从 id 解析）逐条 `DescribeStrategy` 后不刷新该字段（StrategyIds 是创建产物，不会变），保留 Create 写入值。该字段是本次唯一新增参数变更。

### D2 — 资源 ID 采用 StrategyIds 联合

`CreateStrategy` 一次返回多个策略 Id（`StrategyIds` 列表）。资源 ID 用 `tccommon.FILED_SP`（`#`）将列表拼接为一个字符串：`strings.Join(strategyIds, tccommon.FILED_SP)`。Read/Update/Delete 中 `strings.Split(d.Id(), tccommon.FILED_SP)` 还原列表。

- Delete 入参 `StrategyIds` 直接用还原后的列表。
- Read：对列表中每个 id 调一次 `DescribeStrategy` 校验存在性；若全部查询为空，打印 `log.Printf("[CRUD] fwm_strategy id=%s", d.Id())` 后 `d.SetId("")`。
- import：`schema.ImportStatePassthrough`，文档示例说明需使用联合 id（`s1#s2#...`）。

**理由**：与项目"多字段联合 id 用 FILED_SP"惯例一致；`StrategyIds` 是创建产物且是 Delete/Describe 的身份键，作为 id 最自然。

### D3 — CRD-only 资源的 Update 处理

fwm SDK 无 `ModifyStrategy`。按项目惯例（CRD-only 资源）：仅 `product` 标 `ForceNew`（更换产品类型=重新创建），其余顶层入参 `receive_account` / `pre_strategy` / `post_strategy` / `receive_group` 加入 `immutableArgs`。Update 中遍历 `immutableArgs`，若 `d.HasChange(v)` 则返回 `fmt.Errorf("... is immutable, please recreate resource")`。`strategy_ids` 为 Computed，不加入 `immutableArgs`，不在 Update 中变更。

**理由**：无 Modify 接口意味着所有入参变更都无法同步到云端，必须重建；immutableArgs 让用户在尝试变更时得到明确报错而非静默忽略。

### D4 — 客户端访问器

在 `tencentcloud/connectivity/client.go` 新增字段 `fwmv20250611Conn *fwmv20250611.Client` 与方法 `UseFwmV20250611Client()`，实现严格参考 `UseIgtmV20231024Client`（client.go:2248-2257）：

```go
func (me *TencentCloudClient) UseFwmV20250611Client() *fwmv20250611.Client {
    if me.fwmv20250611Conn != nil {
        return me.fwmv20250611Conn
    }
    cpf := me.NewClientProfile(300)
    me.fwmv20250611Conn, _ = fwmv20250611.NewClient(me.Credential, me.Region, cpf)
    me.fwmv20250611Conn.WithHttpTransport(&LogRoundTripper{})
    return me.fwmv20250611Conn
}
```

并在 `TencentCloudClient` struct 与 import 别名（`fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"`）处补充。

### D5 — retry 与空指针保护

- Create retry 块内：调 `CreateStrategyWithContext`，错误 `tccommon.RetryError(e)`，并检查 `result == nil || result.Response == nil` 返回 `NonRetryableError`。
- 块外：检查 `response.Response.StrategyIds == nil || len(...) == 0` → 返回 `fmt.Errorf("StrategyIds is nil/empty")`（NonRetryableError 语义）。打印 `logId` 与 `d.Id()` 便于排障。`d.SetId(...)` 后再 `d.Set("strategy_ids", ...)`。
- Read retry 块内：调 `DescribeStrategyWithContext`，检查 `result == nil || result.Response == nil || result.Response.Strategy == nil` → 返回 `NonRetryableError`（不在此清空 id）。
- Delete retry 块内：调 `DeleteStrategyWithContext`，检查 `result == nil || result.Response == nil` → `NonRetryableError`。
- 错误信息统一使用资源名 `fwm_strategy`（小写蛇形）。

## Risks / Trade-offs

- [StrategyIds 为多值，复合 id 含 `#`] → 与其他含特殊字符的 id 风险相同；StrategyId 为纯数字字符串，不含 `#`，安全。
- [无 Modify 接口，Update 不真正同步] → immutableArgs 报错引导用户重建，避免 state 与云端漂移。
- [DescribeStrategy 单条查询，Read 需多次调用] → StrategyIds 通常数量有限（前后区规则组各若干条），可接受；若列表很长可后续优化为批量接口。
- [gomonkey mock 测试] → 新增资源用 mock，覆盖 Create/Read/Delete 主路径与空指针分支，确保构建可编译。