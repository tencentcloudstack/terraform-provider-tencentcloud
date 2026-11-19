## Context

云防火墙管理（FWM，Firewall Manager）产品提供「下发规则策略（Strategy）」能力。下发策略 = 规则组 × 执行区域（前区 pre / 后区 post）× 接收账号集合，是规则下发的载体；规则组本身由其他资源/接口创建，不在本资源范围内。

当前 Terraform Provider 未覆盖 FWM 产品。本次新增通用资源 `tencentcloud_fwm_strategy`（RESOURCE_KIND_GENERAL），调用以下 5 个 FWM 云 API 完成单条策略的 CRUD：

| TF 操作 | 云 API | 包路径 |
|---------|--------|--------|
| Create | CreateStrategy | `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611` |
| Read（主） | DescribeStrategy | 同上 |
| Read（备） | DescribeStrategies | 同上 |
| Update | ModifyStrategy | 同上 |
| Delete | DeleteStrategy | 同上 |

SDK 已存在于 `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611/`。client.go 中方法签名：
- `func (c *Client) CreateStrategy(request *CreateStrategyRequest) (response *CreateStrategyResponse, err error)`
- `func (c *Client) DescribeStrategy(request *DescribeStrategyRequest) (response *DescribeStrategyResponse, err error)` — 可能返回错误码 `ResourceNotFound` / `ResourceNotFound.NotFound`
- `func (c *Client) DescribeStrategies(request *DescribeStrategiesRequest) (response *DescribeStrategiesResponse, err error)`
- `func (c *Client) ModifyStrategy(request *ModifyStrategyRequest) (response *ModifyStrategyResponse, err error)`
- `func (c *Client) DeleteStrategy(request *DeleteStrategyRequest) (response *DeleteStrategyResponse, err error)`

API 版本：`v20250611`。describe action 由 `c.InitBaseRequest(..., "fwm", APIVersion, ...)` 指定。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_fwm_strategy` 资源，覆盖单条下发策略的 Create/Read/Update/Delete/Import
- 单实例 1:1 建模：Create 每次只发 1 个 `PreStrategy`（或 `PostStrategy`）元素，取 `StrategyIds[0]`
- 用 DescribeStrategy 做精确单查 Read，ResourceNotFound 时清空 state 不报错
- ModifyStrategy 增量更新：只传变更字段
- DeleteStrategy 幂等：不存在的 ID 静默成功
- 新增 `UseFwmV20250611Client()` 客户端访问方法
- 在 provider.go / provider.md 注册资源

**Non-Goals:**
- 不管理规则组（RuleGroup，由其他资源创建）
- 不实现策略下发（Dispatch）能力，只做策略本身的增删改查
- 不生成 waitReady/async 轮询（四个接口均为同步）
- 不处理多元素批量创建（多元素要求用户拆成多个资源实例）

## Decisions

### D1: 单实例 1:1 建模，资源 ID = 服务端 StrategyId

CreateStrategy 入参支持 `PreStrategy[]` 与 `PostStrategy[]`，每元素 `{GroupId,Sequence}` 生成一条策略。但本资源按 1:1 建模：**每次只发 1 个元素**。
- 默认传 `PreStrategy`（1 个元素，对应 `exec_area=pre`）；若用户配置 `exec_area=post` 则传 `PostStrategy`（1 个元素）。
- 取响应 `StrategyIds[0]` 作为 `d.SetId()`。
- 若用户在 `pre_strategy`/`post_strategy` 配置多个元素，provider 返回 error 提示拆成多个资源实例。
- 资源 ID 为 `StrategyId`（`fwms-xxxx`），**非复合 ID**，import 直接传 `StrategyId`。

### D2: Read 优先使用 DescribeStrategy 单实例查询

存在 `DescribeStrategy`（按 `StrategyId` 精确单查，不存在返回 `ResourceNotFound`），比列表 `DescribeStrategies` 更适合 Read/Import。采用它做主 Read 接口。
- 遇 `ResourceNotFound` 或查不到（`response.Response.Strategy == nil`）时，打印 `log.Printf("[CRUD] fwm_strategy id=%s", d.Id())` 后 `d.SetId("")` 移除 state 不报错。
- DescribeStrategy 响应由 `convertStrategyToStrategyResp` 组装，**不回传 `Product` 与 `ExecArea` 字段**，因此 `product`/`exec_area` 无法回读 → 标 `Required` + `ForceNew`，Read 时**不 set** 它们，保留 state 中用户配置值，不参与 diff。

### D3: product / sub_account_uin / exec_area 为 Required + ForceNew

- `product`：产品枚举 `enterprise_sg` / `cfw_vpc_acl` / `cfw_nat_acl` / `cfw_edge_acl`，必填，**不回读**（DescribeStrategy 不返回）。
- `exec_area`：执行区域 `pre` / `post`，必填，**不回读**（DescribeStrategy 不返回）。
- `sub_account_uin`：子账号 UIN，必填，**不回读**（DescribeStrategy 不返回）。
- 三者变更视为重建（ForceNew）。

> 说明：`sub_account_uin` 在 CreateStrategy 入参 struct 中并不直接体现为顶层字段（服务端由调用者身份/凭证上下文确定 SubAccountUin）。根据需求，schema 中保留 `sub_account_uin` 为 Required+ForceNew 字段并在 description 注明其语义，但因 API 入参无对应字段，Create/Modify 时**不发送**该字段（仅作为声明性参数与 state 保留），Read 时不回读。这符合需求「product/sub_account_uin/exec_area 标 required + forceNew」且不可回读的约束。

### D4: group_id / sequence / receive_account / receive_group 可变，走 Update

- `group_id`（string）：规则组 ID。Create 时写入 `PreStrategy/PostStrategy[0].GroupId`；Modify 时写入 `request.GroupId`。
- `sequence`（int）：优先级，从 1 开始。Create 时写入 `PreStrategy/PostStrategy[0].Sequence`；Modify 时写入 `request.Sequence`。
- `receive_account`（string list）：接收账号 UIN 数组。Create/Modify 时写入 `request.ReceiveAccount`（`[]*string`）。
- `receive_group`（string list）：接收账号组 ID 数组。Create/Modify 时写入 `request.ReceiveGroup`（`[]*string`）。

ModifyStrategy 是**增量语义**：字段对应的 API 指针为 `*int64`/`*string`/`[]*string`，nil 表示不改。Update 实现中只把发生 `d.HasChange` 的字段赋给 request 指针；未变更字段保持指针 nil（不改）。不要把所有字段无脑赋值。

### D5: 接收对象类型不对称的回读处理

入参 `ReceiveAccount` 是 `[]*string`（UIN 字符串数组）、`ReceiveGroup` 是`[]*string`（账号组 ID 数组）；响应 `ReceiveAccount` 是 `[]*ReceiveAccount`（对象数组 `{Uin,Nickname,ReceiverType,Members}`，`ReceiverType` 0=账号、1=账号组）。
- 回读时遍历响应 `ReceiveAccount`，按 `ReceiverType` 拆分：
  - `ReceiverType == 0`（账号）→ 收集 `Uin` 到 `receive_account` 集合
  - `ReceiverType == 1`（账号组）→ 收集 `Uin` 到 `receive_group` 集合
- 只比对 `Uin` 集合，不比 `Nickname`。

### D6: 只读字段（Computed）

从 DescribeStrategies 响应可读、入参不可设的字段，标 `Computed`：
`group_name`、`rule_count`、`rule_status`、`err_msg`、`error_type`、`create_time`、`update_time`、`create_by`、`create_name`、`update_by`、`update_name`。
（`strategy_id` 在 import 场景由 id 表示，仍标 Computed。）

### D7: Create 返回值校验

CreateStrategy 成功后，**必须**检查 `response.Response.StrategyIds` 是否为空、长度是否为 0、`StrategyIds[0]` 是否为空字符串。任一为空则返回 `NonRetryableError`，避免写入空 id 触发后续状态混乱。检查前先打印 `logId` 与 `d.Id()` 便于排障。

### D8: 四接口同步，无轮询

Create/Modify/Delete 成功即已落库，查询立即可读。不生成 waitReady/async。所有云 API 调用用 `resource.Retry(tccommon.ReadRetryTimeout/WriteRetryTimeout)` 包装，失败用 `tccommon.RetryError(e)` 包装返回。

### D9: FWM 客户端访问方法

参照 `UseIgtmV20231024Client()`（client.go:2248）模式新增 `UseFwmV20250611Client()`：
- import：`fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"`
- 连接字段：`fwmv20250611Conn *fwmv20250611.Client`
- 方法体：懒加载 + `fwmv20250611.NewClient(me.Credential, me.Region, cpf)` + `WithHttpTransport(&LogRoundTripper{})`

### D10: Schema 不暴露 PreStrategy/PostStrategy 嵌套数组

为符合 1:1 单实例建模，schema **不**创建 `pre_strategy`/`post_strategy` 嵌套结构列表。而是把 `{GroupId,Sequence}` 平铺为顶层 `group_id`（string）+ `sequence`（int），并用 `exec_area`（pre/post）决定 Create 时填入 `PreStrategy` 还是 `PostStrategy`。`group_id`/`sequence` 为可选。

## Risks / Trade-offs

- **[product/exec_area 不可回读]** → Mitigation: 标 ForceNew 且 Read 不 set，依赖 state 保留用户配置值。修改三者会重建资源，避免 state 漂移。
- **[sub_account_uin 无对应 API 入参字段]** → Mitigation: schema 保留为 Required+ForceNew 声明性参数，Create/Modify 不发送、Read 不回读，避免误传 API 不支持的字段。
- **[ReceiveAccount 入参/回读类型不对称]** → Mitigation: 按 ReceiverType 拆分回填两个集合，只比 Uin 不比 Nickname。
- **[DeleteStrategy 对不存在 ID 静默成功]** → Mitigation: destroy 天然幂等，无需吞错；Read 查不到即 SetId("")。
- **[并发保护：同 adminAppId+Product 有进行中下发任务时 Create/Modify 报错]** → Mitigation: 用 WriteRetryTimeout 重试而非直接失败。

## API Struct 关键字段定义（vendor 摘录，供实施阶段直接引用）

> 以下摘自 `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611/models.go`，实施阶段无需回读大文件。

### CreateStrategyRequest (models.go:886)
```go
type CreateStrategyRequest struct {
    *tchttp.BaseRequest
    Product        *string         `json:"Product,omitnil,omitempty"`        // 产品类型
    ReceiveAccount []*string       `json:"ReceiveAccount,omitnil,omitempty"` // 下发账号
    PreStrategy    []*StrategyReq  `json:"PreStrategy,omitnil,omitempty"`    // 前区规则组
    PostStrategy   []*StrategyReq  `json:"PostStrategy,omitnil,omitempty"`   // 后区规则组
    ReceiveGroup   []*string       `json:"ReceiveGroup,omitnil,omitempty"`   // 下发账号组
}
```
### CreateStrategyResponse (models.go:937)
```go
type CreateStrategyResponse struct {
    *tchttp.BaseResponse
    Response *CreateStrategyResponseParams `json:"Response"`
}
type CreateStrategyResponseParams struct {
    StrategyIds []*string `json:"StrategyIds,omitnil,omitempty"` // 策略Id列表
    RequestId   *string   `json:"RequestId,omitnil,omitempty"`
}
```

### DescribeStrategyRequest / Response (models.go:2636 / 2671) — 单实例查询（主 Read）
```go
type DescribeStrategyRequest struct {
    *tchttp.BaseRequest
    StrategyId *string `json:"StrategyId,omitnil,omitempty"` // 策略Id
}
type DescribeStrategyResponse struct {
    *tchttp.BaseResponse
    Response *DescribeStrategyResponseParams `json:"Response"`
}
type DescribeStrategyResponseParams struct {
    Strategy  *StrategyResp `json:"Strategy,omitnil,omitempty"` // 策略详情
    RequestId *string       `json:"RequestId,omitnil,omitempty"`
}
```

### DescribeStrategiesRequest / Response (models.go:2418 / 2472) — 列表查询（备用）
```go
type DescribeStrategiesRequest struct {
    *tchttp.BaseRequest
    Product  *string          `json:"Product,omitnil,omitempty"`  // 产品类型
    ExecArea *string          `json:"ExecArea,omitnil,omitempty"` // 执行区域 pre/post
    Filters  []*CommonFilter  `json:"Filters,omitnil,omitempty"`  // 筛选条件
    Limit    *int64           `json:"Limit,omitnil,omitempty"`     // 条数限制
    Offset   *int64           `json:"Offset,omitnil,omitempty"`    // 偏移量
}
type DescribeStrategiesResponse struct {
    *tchttp.BaseResponse
    Response *DescribeStrategiesResponseParams `json:"Response"`
}
type DescribeStrategiesResponseParams struct {
    TotalCount  *uint64          `json:"TotalCount,omitnil,omitempty"`
    Strategies  []*StrategyResp  `json:"Strategies,omitnil,omitempty"`
    RequestId   *string          `json:"RequestId,omitnil,omitempty"`
}
```

### ModifyStrategyRequest (models.go:3464) — 增量语义，指针 nil=不改
```go
type ModifyStrategyRequest struct {
    *tchttp.BaseRequest
    StrategyId     *string    `json:"StrategyId,omitnil,omitempty"`     // 策略Id
    ReceiveAccount []*string  `json:"ReceiveAccount,omitnil,omitempty"` // 下发规则接收账号（可选，增量）
    Sequence       *int64     `json:"Sequence,omitnil,omitempty"`        // 优先级（可选，增量）
    GroupId        *string    `json:"GroupId,omitnil,omitempty"`         // 规则组Id（可选，增量）
    ReceiveGroup   []*string  `json:"ReceiveGroup,omitnil,omitempty"`   // 下发规则接收账号组（可选，增量）
}
```

### DeleteStrategyRequest / Response (models.go:1331 / 1363)
```go
type DeleteStrategyRequest struct {
    *tchttp.BaseRequest
    StrategyIds []*string `json:"StrategyIds,omitnil,omitempty"` // 策略Id列表
}
type DeleteStrategyResponse struct {
    *tchttp.BaseResponse
    Response *DeleteStrategyResponseParams `json:"Response"`
}
type DeleteStrategyResponseParams struct {
    RequestId *string `json:"RequestId,omitnil,omitempty"`
}
```

### StrategyReq (models.go:4449) — Create 时的元素结构
```go
type StrategyReq struct {
    GroupId   *string `json:"GroupId,omitnil,omitempty"`   // 规则组Id
    Sequence  *int64  `json:"Sequence,omitnil,omitempty"`   // 优先级
}
```

### StrategyResp (models.go:4457) — Read 响应元素结构
```go
type StrategyResp struct {
    StrategyId     *string            `json:"StrategyId,omitnil,omitempty"`     // 策略Id
    GroupId        *string            `json:"GroupId,omitnil,omitempty"`         // 规则组Id
    GroupName      *string            `json:"GroupName,omitnil,omitempty"`       // 规则组名称（只读）
    RuleCount      *int64             `json:"RuleCount,omitnil,omitempty"`       // 规则数（只读）
    RuleStatus     *int64             `json:"RuleStatus,omitnil,omitempty"`      // 策略状态（只读）
    ReceiveAccount []*ReceiveAccount  `json:"ReceiveAccount,omitnil,omitempty"`  // 下发账号（对象数组）
    Sequence       *int64             `json:"Sequence,omitnil,omitempty"`         // 优先级
    ErrMsg         *string            `json:"ErrMsg,omitnil,omitempty"`           // 下发失败原因（只读）
    ErrorType      *string            `json:"ErrorType,omitnil,omitempty"`        // 下发失败原因类型（只读）
    UpdateTime     *string            `json:"UpdateTime,omitnil,omitempty"`       // 更新时间（只读）
    CreateBy      *string             `json:"CreateBy,omitnil,omitempty"`          // 创建人（只读）
    UpdateBy      *string             `json:"UpdateBy,omitnil,omitempty"`          // 更新人（只读）
    ExecArea      *string             `json:"ExecArea,omitnil,omitempty"`           // 执行区域（只读，且单查不回传）
    CreateName    *string             `json:"CreateName,omitnil,omitempty"`         // 创建人名称（只读）
    UpdateName    *string             `json:"UpdateName,omitnil,omitempty"`         // 更新人名称（只读）
    CreateTime    *string             `json:"CreateTime,omitnil,omitempty"`         // 创建时间（只读）
}
```

### ReceiveAccount (models.go:4068) — 响应下发账号对象
```go
type ReceiveAccount struct {
    Uin          *string       `json:"Uin,omitnil,omitempty"`          // 租户 uin
    Nickname     *string       `json:"Nickname,omitnil,omitempty"`     // 租户名称
    ReceiverType *int64        `json:"ReceiverType,omitnil,omitempty"` // 0=账号uin，1=账号组
    Members      []*MemberInfo `json:"Members,omitnil,omitempty"`      // ReceiverType=1 时返回账号列表
}
type MemberInfo struct {
    AppId    *string `json:"AppId,omitnil,omitempty"`
    Uin      *string `json:"Uin,omitnil,omitempty"`
    Nickname *string `json:"Nickname,omitnil,omitempty"`
    MemberId *string `json:"MemberId,omitnil,omitempty"`
    NodeName *string `json:"NodeName,omitnil,omitempty"`
}
```

### CommonFilter (models.go:321)
```go
type CommonFilter struct {
    Name        *string   `json:"Name,omitnil,omitempty"`        // 筛选字段名
    Values      []*string `json:"Values,omitnil,omitempty"`      // 筛选值列表
    OperatorType *int64   `json:"OperatorType,omitnil,omitempty"` // 1=精确匹配,7=in,9=模糊匹配
}
```

## 代码风格参照：tencentcloud_igtm_strategy（关键结构说明）

参考文件 `tencentcloud/services/igtm/resource_tc_igtm_strategy.go`，关键结构：

1. **Resource 定义骨架**：
```go
func ResourceTencentCloudIgtmStrategy() *schema.Resource {
    return &schema.Resource{
        Create: resourceTencentCloudIgtmStrategyCreate,
        Read:   resourceTencentCloudIgtmStrategyRead,
        Update: resourceTencentCloudIgtmStrategyUpdate,
        Delete: resourceTencentCloudIgtmStrategyDelete,
        Importer: &schema.ResourceImporter{
            State: schema.ImportStatePassthrough,
        },
        Schema: map[string]*schema.Schema{ ... },
    }
}
```

2. **Create 函数骨架**（注意 retry + nil 检查 + SetId 在 retry 之外）：
```go
func resourceTencentCloudIgtmStrategyCreate(d, meta) error {
    defer tccommon.LogElapsed("resource.tencentcloud_igtm_strategy.create")()
    defer tccommon.InconsistentCheck(d, meta)()
    var (
        logId = tccommon.GetLogId(tccommon.ContextNil)
        ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
        request  = xxx.NewCreateXxxRequest()
        response = xxx.NewCreateXxxResponse()
    )
    // ... 填充 request ...
    reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
        result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseXxxClient().CreateXxxWithContext(ctx, request)
        if e != nil {
            return tccommon.RetryError(e)
        }
        log.Printf("[DEBUG]%s api[%s] success ...", logId, request.GetAction(), ...)
        if result == nil || result.Response == nil {
            return resource.NonRetryableError(fmt.Errorf("Create xxx failed, Response is nil."))
        }
        response = result
        return nil
    })
    if reqErr != nil {
        log.Printf("[CRITAL]%s create xxx failed, reason:%+v", logId, reqErr)
        return reqErr
    }
    // 校验返回 id，SetId（retry 块外）
    ...
    d.SetId(strategyId)
    return resourceTencentCloudIgtmStrategyRead(d, meta)
}
```

3. **Read 函数骨架**（注意查不到先 log 再 SetId("")）：
```go
func resourceTencentCloudIgtmStrategyRead(d, meta) error {
    defer tccommon.LogElapsed("resource.tencentcloud_igtm_strategy.read")()
    defer tccommon.InconsistentCheck(d, meta)()
    var (
        logId   = tccommon.GetLogId(tccommon.ContextNil)
        ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
        service = XxxService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
    )
    respData, err := service.DescribeXxxById(ctx, id)
    if err != nil {
        return err
    }
    if respData == nil {
        log.Printf("[WARN]%s resource `tencentcloud_xxx` [%s] not found...", logId, d.Id())
        d.SetId("")
        return nil
    }
    // set 字段（先判 nil）
    ...
    return nil
}
```

4. **Update 函数骨架**（mutableArgs 判 HasChange + retry）：
```go
func resourceTencentCloudIgtmStrategyUpdate(d, meta) error {
    // ...
    needChange := false
    mutableArgs := []string{"field1", "field2"}
    for _, v := range mutableArgs {
        if d.HasChange(v) {
            needChange = true
            break
        }
    }
    if needChange {
        request := xxx.NewModifyXxxRequest()
        // 只填变更字段
        request.Id = &id
        reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError { ... })
        if reqErr != nil { return reqErr }
    }
    return resourceTencentCloudIgtmStrategyRead(d, meta)
}
```

5. **Delete 函数骨架**：
```go
func resourceTencentCloudIgtmStrategyDelete(d, meta) error {
    // ...
    request := xxx.NewDeleteXxxRequest()
    request.Id = &id
    reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
        result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseXxxClient().DeleteXxxWithContext(ctx, request)
        if e != nil { return tccommon.RetryError(e) }
        log.Printf("[DEBUG]...")
        if result == nil || result.Response == nil {
            return resource.NonRetryableError(fmt.Errorf("Delete xxx failed, Response is nil."))
        }
        return nil
    })
    if reqErr != nil {
        log.Printf("[CRITAL]%s delete xxx failed, reason:%+v", logId, reqErr)
        return reqErr
    }
    return nil
}
```

**注意差异**：igtm_strategy 用复合 ID（`instanceId#strategyId`），本资源用单一 `StrategyId`。igtm Read 用 service 层 `DescribeIgtmStrategyById`；本资源 Read 优先用 `DescribeStrategy`（SDK 内置），可直接在 Read 函数中调用或封装到 service 层。helper 包提供 `helper.String()`、`helper.IntInt64()`、`helper.Strings()` 等转换函数。

## Provider 注册参照

`provider.go` 顶部 import 块需新增：
```go
"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/fwm"
```
资源 map 中新增（参照 igtm_strategy 注册位置，line 2783）：
```go
"tencentcloud_fwm_strategy": fwm.ResourceTencentCloudFwmStrategy(),
```

## 客户端访问方法参照（connectivity/client.go）

参照 `UseIgtmV20231024Client()` (line 2248) 新增：
```go
// import
fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"
// struct 字段
fwmv20250611Conn *fwmv20250611.Client
// 方法
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
调用方式：`meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmV20250611Client()`