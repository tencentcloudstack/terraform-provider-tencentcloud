## Context

本变更新增腾讯云云防火墙（CFW / 产品代号 `fwm`）互联网边界（Edge）ACL 单条规则的 Terraform 通用 CRUD 资源 `tencentcloud_fwm_edge_acl_rule`。

SDK 位于 `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611/`（namespace `v20250611`），已 vendored 就绪。当前 Provider 中无任何 `fwm` 服务包、无 `UseFwmV20250611Client` 访问器、无 fwm 资源注册，本变更将三者一并引入。

四个接口均为同步接口（无 TaskId、无中间态），Create/Modify/Delete 成功即生效，**不生成 waitReady/异步轮询逻辑**。

### 业务语义关键约束（接口名无法体现）
1. 本资源管理「单条 Edge ACL 规则」。规则组（Group）是父容器，由独立资源（如 `fwm_rule_group`）创建，**不在本资源内创建**。CreateEdgeAclRule 前必须已存在规则组且其 Product=cfw_edge_acl。
2. CreateEdgeAclRule 支持 `Rules` 数组批量创建，但本资源为单实例：`Rules` 固定 1 个元素；返回 `RuleIds` 数组，**取首元素**作为 `RuleId`。服务端按入参顺序遍历插入，顺序严格对应。
3. DescribeEdgeAclRules 的 `Direction` 是**必填**（0=出向，1=入向），非可选过滤项。精准查单条用 `GroupId` + `Filters=[{Name:"RuleId",Values:[ruleId],OperatorType:1}]`（1=EQUAL）。该接口是列表查询（非单实例 Get），取 `Rules` 首元素。
4. ModifyEdgeAclRule 按 `Rule.RuleId` 单条更新，语义为**全量覆盖**（服务端先查旧记录再以完整新值写入），必须带所有可写字段。不可修改 GroupId。
5. DeleteEdgeAclRule 是批量接口，单删时 `RuleIds` 只传 1 个元素。RuleId 不存在返回 `ResourceNotFound`，destroy 视为删除成功（幂等）。

### 云 API struct 关键字段定义（vendor 摘录）

**`CreateEdgeAclRuleRequestParams`** (models.go:471):
```go
GroupId *string                       // 规则组ID
Rules   []*EdgeAclRuleInfo            // 规则列表，本资源固定 1 个元素
```
**`CreateEdgeAclRuleResponseParams`** (models.go:510):
```go
RuleIds  []*string                    // 规则Id列表，取首元素为 RuleId
RequestId *string
```

**`EdgeAclRuleInfo`**（Create/Modify 入参 Rule 结构，models.go:2835）:
```go
RuleId         *string   // 修改时必填
OrderIndex     *int64    // 执行顺序，数字越小优先级越高
Direction      *int64    // 0-出向 1-入向
SourceContent  *string
SourceType     *string   // ip/domain/template/instance/tag...
TargetContent  *string
TargetType     *string
Port           *string
Protocol       *string   // TCP/UDP/ICMP/ANY
RuleAction     *string   // accept/drop/log
Description    *string
Scope          *string   // serial/side/all
BelongMemberId *string   // instance/tag 等场景必填
ParamTemplateId *string
```

**`EdgeAclRuleResp`**（Describe 出参 Rules 元素，models.go:2879）— 注意字段名与入参错位:
```go
RuleId          *string
Sequence        *int64    // ← 入参为 OrderIndex（错位！回读填 order_index）
Direction        *int64
SourceContent   *string
SourceType      *string
SourceName      *string   // 只读
TargetContent   *string
TargetType      *string
TargetName      *string   // 只读
Port            *string
Protocol        *string
RuleAction      *string   // 出参已转小写字符串，与入参一致
Description     *string
Scope           *string   // 出参为 all/serial/side 字符串，与入参一致
CountryName     *string   // 只读
CityName        *string   // 只读
ParamTemplateId *string
ParamTemplateName *string // 只读
Invalid         *int64    // 只读
BelongMember    *MemberInfo // ← 入参为 BelongMemberId 字符串（错位！回读取 .MemberId）
CountryCode     *uint64
CityCode        *uint64
IsRegion        *uint64
CloudCode       *string
IsCloud         *uint64
InstanceName    *string   // 只读
CountryKey      *string
CityKey         *string
CreateTime      *string   // 只读
UpdateTime      *string   // 只读
DnsParseCnt     *int64
HitCount        *int64
LastHitTime     *string
```

**`MemberInfo`** (models.go:3049):
```go
AppId    *string
Uin      *string
Nickname *string
MemberId *string   // ← 回读 belong_member_id 取此字段
NodeName *string
```

**`CommonFilter`** (models.go:321):
```go
Name        *string
Values      []*string
OperatorType *int64  // 1=精确匹配(EQUAL)，7=in，9=模糊匹配
```

**`DescribeEdgeAclRulesRequestParams`** (models.go:1441):
```go
GroupId    *string                 // 必填
Direction  *int64                  // 必填 0/1
Filters    []*CommonFilter        // 用 RuleId 精准过滤
Limit      *int64                  // 默认100，最大1000 → 取 1000
Offset     *int64
Order      *string                 // asc/desc
By         *string                 // Sequence/RuleId...
```
**`DescribeEdgeAclRulesResponseParams`** (models.go:1515):
```go
TotalCount    *int64
Rules         []*EdgeAclRuleResp   // 取首元素
AllTotalCount *int64
RequestId     *string
```

**`ModifyEdgeAclRuleRequestParams`** (models.go:3067):
```go
GroupId *string
Rule    *EdgeAclRuleInfo          // 单条，必须含 RuleId
```

**`DeleteEdgeAclRuleRequestParams`** (models.go:1089):
```go
GroupId *string
RuleIds []*string                 // 单删传 1 个元素
```
DeleteEdgeAclRuleResponse 仅含 RequestId。

错误码（errors.go）: `RESOURCENOTFOUND = "ResourceNotFound"`, `RESOURCENOTFOUND_NOTFOUND = "ResourceNotFound.NotFound"`。

### 代码风格参照文件：`tencentcloud_igtm_strategy`
参照 `tencentcloud/services/igtm/resource_tc_igtm_strategy.go` 的结构与命名风格：
- 资源函数 `ResourceTencentCloudFwmEdgeAclRule()` 返回 `*schema.Resource`，含 Create/Read/Update/Delete + Importer(`schema.ImportStatePassthrough`)。
- CRUD 函数签名 `func resourceTencentCloudFwmEdgeAclRuleXxx(d *schema.ResourceData, meta interface{}) error`，开头 `defer tccommon.LogElapsed("...")` + `defer tccommon.InconsistentCheck(d, meta)`。
- 局部变量取 `logId := tccommon.GetLogId(tccommon.ContextNil)`, `ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)`。
- 复合 ID 用 `strings.Split(d.Id(), tccommon.FILED_SP)`，校验 `len == 2`，否则 `fmt.Errorf("id is broken,%s", d.Id())`。
- SDK 调用包裹 `resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {...})`，错误经 `tccommon.RetryError(e)` 包装。
- `request := fwmv20250611.NewCreateEdgeAclRuleRequest()` 风格构造请求；指针字段用 `helper.String(v)`, `helper.Int64(v)` 等设置。
- import 别名: `fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"`, `tccommon ".../tencentcloud/common"`, `helper ".../tencentcloud/internal/helper"`。

参照 `tencentcloud/services/igtm/service_tencentcloud_igtm.go`:
- `NewFwmService(client *connectivity.TencentCloudClient) FwmService`，`type FwmService struct { client *connectivity.TencentCloudClient }`。
- `DescribeFwmEdgeAclRuleById(ctx, groupId, ruleId, direction)` 包裹 `resource.Retry(tccommon.ReadRetryTimeout, ...)`，未找到返回 `(nil, nil)`。

参照 `tencentcloud/connectivity/client.go` `UseIgtmV20231024Client` (line 2247): lazy-init cached `*fwmv20250611.Client`，`me.NewClientProfile(300)` + `WithHttpTransport(&LogRoundTripper{})`。需新增 import 别名与私有字段 `fwmv20250611Conn *fwmv20250611.Client`。

### 单元测试风格参照：`resource_tc_teo_create_cls_index_operation_test.go`
本资源为新增资源，测试**不使用 terraform 测试套件**，使用 gomonkey mock 云 API：
- `package fwm_test`，`import "github.com/agiledragon/gomonkey/v2"`, `"github.com/stretchr/testify/assert"`, `fwmv20250611 ".../fwm/v20250611"`。
- 定义 `mockMeta` 实现 `tccommon.ProviderMeta`，`newMockMeta()` 返回 `&mockMeta{client: &connectivity.TencentCloudClient{}}`，`ptrString(s)` 辅助。
- `patches.ApplyMethodReturn(newMockMeta().client, "UseFwmV20250611Client", fwmClient)` mock client；`patches.ApplyMethodFunc(fwmClient, "CreateEdgeAclRule", ...)` mock 各 SDK 方法返回。
- `schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{...})` 构造 ResourceData 调用 CRUD。
- 运行：`go test ./tencentcloud/services/fwm/ -run "..." -v -count=1 -gcflags="all=-l"`。

## Goals / Non-Goals

**Goals:**
- 资源 `tencentcloud_fwm_edge_acl_rule` 覆盖单条 Edge ACL 规则 CRUD 生命周期，支持 `GroupId#RuleId` 两段式导入。
- Schema 一一映射 Create/Modify 入参 + Describe 只读字段，并显式处理 OrderIndex↔Sequence、BelongMemberId↔BelongMember.MemberId 错位。
- Delete 幂等处理 `ResourceNotFound`（规则已不存在视为删除成功）。
- 所有 SDK 调用包裹 retry + 防御性空值检查；Create 返回空 id 时返回 NonRetryableError。
- 新增 `UseFwmV20250611Client` 访问器、`FwmService` 服务层、provider 注册与文档、gomonkey 单元测试。

**Non-Goals:**
- 不创建/管理规则组（Group）— 由独立资源负责。
- 不支持批量创建（Rules 固定 1 个）。
- 不暴露 Describe 的 Order/By/排序参数给用户（内部精准过滤用 RuleId）。
- 不生成异步轮询逻辑（接口均为同步）。

## Decisions

### D1: 复合 ID = `GroupId#RuleId`，使用 `tccommon.FILED_SP`
`ModifyEdgeAclRule`/`DeleteEdgeAclRule`/`DescribeEdgeAclRules` 均需 `GroupId`，而 Create 仅返回 `RuleId`，故将两者复合入资源 ID。Create 成功后 `d.SetId(strings.Join([]string{groupId, ruleId}, tccommon.FILED_SP))`。Read/Update/Delete 拆分 `d.Id()` 校验 2 段。导入通过 `schema.ImportStatePassthrough` 复用 Read（用户传 `GroupId:RuleId` 或 `GroupId#RuleId`，Read 内做兼容拆分）。

### D2: 字段错位显式映射（避免 import 后永久 diff）
| 用途 | 入参 (EdgeAclRuleInfo) | 出参 (EdgeAclRuleResp) | schema 字段 |
|---|---|---|---|
| 执行顺序 | `OrderIndex` | `Sequence` | `order_index` |
| 归属成员 | `BelongMemberId`（string）| `BelongMember.MemberId` | `belong_member_id` |
| 生效范围 | `Scope` | `Scope`（同名字符串）| `scope` |
| 规则动作 | `RuleAction` | `RuleAction`（小写字符串）| `rule_action` |

写入(Crate/Modify): 用 `OrderIndex`/`BelongMemberId`/`Scope`/`RuleAction`。回读(Read): `order_index` ← `Sequence`；`belong_member_id` ← `BelongMember.MemberId`（非空才 set）；`scope`/`rule_action` 直接透传出参字符串，**不做整数/"0" 二次映射**。

### D3: Describe 走列表查询 + RuleId 过滤
无单实例 Get 接口。Read 构造 `DescribeEdgeAclRulesRequest`：`GroupId`、`Direction`（来自 schema，必填），`Filters` = 单个 `CommonFilter{Name:"RuleId", Values:[ruleId], OperatorType:1}`，`Limit=1000`（API 最大值，单条查询足够）。取 `Rules` 首元素。若 `Rules` 为空 → 打印 `[CRUD]` log → `d.SetId("")`。

### D4: ForceNew 精确范围
- `direction`：Required + ForceNew（Modify 不支持改方向）。
- `source_type` / `target_type`：ForceNew。
- `group_id`：Required，**不 ForceNew**（仅作父引用，Update 时透传 GroupId；变更它语义上等于换父组，服务端不支持，但不在 schema 层强制重建，避免误删；文档提示需配套调整）。
- `order_index`：Required，可变（原地更新）。
- 其余规则属性可变。

### D5: bool/int 可选字段用 `d.GetOkExists`
`direction`、`order_index`（int）等用 `d.GetOkExists(...)` 取值，避免 0 被当未设置。字符串字段用 `d.GetOk`。

### D6: Delete 幂等 + Modify 全量覆盖
- Delete：`RuleIds` 传单元素；retry 内若错误码含 `ResourceNotFound` 返回 nil（幂等），否则 `tccommon.RetryError(e)`。
- Update：检测到任意可变字段 HasChange 即构造 `ModifyEdgeAclRuleRequest.Rule`（含 `RuleId` + 全部可写字段，全量覆盖语义），retry 调用，成功后 Read。

### D7: Create 空值防护
Create retry 内检查 `result == nil || result.Response == nil` → NonRetryableError；retry 外检查 `RuleIds` 为空或 `RuleIds[0] == ""` → 打印 logId/d.Id 后返回 `NonRetryableError`。

## Risks / Trade-offs

- **[Risk] Describe 为列表查询，单条匹配靠 RuleId 过滤** → Mitigation: 严格用 `OperatorType:1`(EQUAL) 且 `Limit=1000`，取首元素；若匹配多条（理论不会，RuleId 唯一）以首元素为准并在 log 中提示。
- **[Risk] Modify 全量覆盖，遗漏字段可能影响规则** → Mitigation: Update 始终从 schema 读取全部可写字段拼装请求，不依赖 diff 只传变更项。
- **[Risk] group_id 不 ForceNew 但服务端不可改** → Mitigation: 文档说明 group_id 变更需重建资源；schema 不强制以与 IGTM 父引用风格一致。
- **[Trade-off] 不暴露 Order/By 排序参数** → 单实例资源不需排序，内部固定精准过滤，降低用户认知负担。