## Context

云防火墙 FWM 的 VPC 边界 ACL 规则 4 个接口（CreateVpcAclRule / DescribeVpcAclRules / ModifyVpcAclRule / DeleteVpcAclRule）发布于 fwm v20250611 SDK，适合做成 Terraform CRUD 资源。资源命名遵守 provider 现有惯例 `tencentcloud_<service>_<resource>`，归属到 `tencentcloud/services/fwm/`（fwm 是全新服务目录）。

本资源是「单条 VPC 边界 ACL 规则」。规则组是父容器（由独立资源 fwm_rule_group 创建），不在本资源内创建。资源 ID 为复合结构 `GroupId#RuleId`。

关键业务语义：
- CreateVpcAclRule 支持 Rules 数组批量创建但本资源为单实例：Rules 固定 1 个元素。RuleIds 返回数组，与入参 Rules 严格按位置一一对应，单实例取 RuleIds[0]。
- 四个接口全部同步（无任务 ID、无中间态），不生成 waitReady/async。
- ModifyVpcAclRule 是全量覆盖（不是增量），复用与 Create 完全相同的字段校验，Update 时所有可写字段全部必填且必须完整回传。
- DescribeVpcAclRules 是列表查询（无单实例 Get），用 Filters 精准查单条。
- VPC 独有：EdgeId 与 FwGroupId 决定规则生效于哪些 VPC，均必填，且建后不可改（forceNew）。
- 入参出参字段名完全一致，唯一需要映射的是 OrderIndex ↔ 响应 Sequence（同一字段两个名字，回读用 Sequence 填 order_index）。另有入参 BelongMemberId（string）↔ 响应 BelongMember（对象），回读取 BelongMember.MemberId。

参考资源：
- `tencentcloud_igtm_strategy`（代码风格参考，复合 ID `FILED_SP` 分隔、`mutableArgs` 模式、retry 包装、nil 安全回读）
- `tencentcloud_cfw_vpc_policy`（业务模板参考，扁平化规则字段 schema、Create 用 Rules 数组单元素、全量覆盖 Update）
- `tencentcloud_cfw_cluster_nat_fw_switch`（测试模板参考，gomonkey mock 单元测试）

## Goals / Non-Goals

**Goals:**
- Schema 字段名与 CreateVpcAclRule 入参（VpcAclRule 结构）1:1 映射（snake_case 化），只读展示字段映射到 Computed。
- 资源 ID = `GroupId#RuleId`（使用 `tccommon.FILED_SP`，值为 `#`）。
- Read 通过 `DescribeVpcAclRules` + `Filters: [{Name: "RuleId", Values: [ruleId], OperatorType: 1}]` 精确查询单个实例；分页 Limit 取 0（默认 100），Offset 取 0。
- Create / Read / Update / Delete 的全部 SDK 调用均用 `resource.Retry(WriteRetryTimeout/ReadRetryTimeout, ...)` 包装。
- 全部接口返回值取值都做空指针保护。
- 代码风格严格对齐 `tencentcloud_igtm_strategy`（变量声明顺序、defer LogElapsed + InconsistentCheck 模式、复合 ID 拆分）。
- 单元测试使用 gomonkey mock 云 API，不使用 terraform 验收测试套件。

**Non-Goals:**
- 不实现对应的 data source —— 本 change 只交付 resource。
- 不在本资源内创建规则组（规则组由独立资源 fwm_rule_group 管理）。
- 不实现异步等待逻辑（四个接口全部同步）。
- 不修改任何已存在资源/数据源/service 方法。

## Decisions

### D1 — Schema 字段映射

| HCL 字段 | SDK 入参字段（VpcAclRule） | SDK 出参字段（VpcAclRuleResp） | 类型 | 必填 | ForceNew | Computed | 说明 |
|---|---|---|---|---|---|---|---|
| `group_id` | `request.GroupId` | — | TypeString | Yes | Yes | No | 规则组ID，父资源引用，四接口均必填，不可改 |
| `edge_id` | `Rules.EdgeId` | `Rules.EdgeId` | TypeString | Yes | Yes | No | 边界防火墙ID：ALL 或 cfws-xxx |
| `fw_group_id` | `Rules.FwGroupId` | `Rules.FwGroupId` | TypeString | Yes | Yes | No | 防火墙实例ID：ALL、cfwg-xxx 或 ccn-xxx |
| `ip_version` | — | `Rules.IpVersion` | TypeString | Yes | Yes | No | IP版本：ipv4 或 ipv6。注意：Create 入参 VpcAclRule 无此字段，仅 Describe 出参有；Create 时不需要传（服务端按规则组推断），但 schema 标 required + forceNew 用于区分与回读 |
| `source_type` | `Rules.SourceType` | `Rules.SourceType` | TypeString | Yes | Yes | No | 源类型：ip/url/template/instance/tag |
| `source_content` | `Rules.SourceContent` | `Rules.SourceContent` | TypeString | Yes | Yes | No | 源地址内容 |
| `dest_type` | `Rules.DestType` | `Rules.DestType` | TypeString | Yes | Yes | No | 目的类型：ip/url/template/instance/tag |
| `dest_content` | `Rules.DestContent` | `Rules.DestContent` | TypeString | Yes | Yes | No | 目的地址内容 |
| `protocol` | `Rules.Protocol` | `Rules.Protocol` | TypeString | Yes | No | No | 协议：TCP/UDP/ICMP/ANY 等 |
| `rule_action` | `Rules.RuleAction` | `Rules.RuleAction` | TypeString | Yes | No | No | 动作：accept/drop/log（小写） |
| `order_index` | `Rules.OrderIndex` | `Rules.Sequence` | TypeInt | Yes | No | No | 优先级（从1开始）。入参 OrderIndex，出参 Sequence，回读用 Sequence 填 order_index |
| `description` | `Rules.Description` | `Rules.Description` | TypeString | Yes | No | No | 规则描述，必填且≤30字符 |
| `port` | `Rules.Port` | `Rules.Port` | TypeString | No | No | No | 端口（ICMP 协议时为空） |
| `param_template_id` | `Rules.ParamTemplateId` | `Rules.ParamTemplateId` | TypeString | No | No | No | 端口模板ID |
| `belong_member_id` | `Rules.BelongMemberId` | `Rules.BelongMember.MemberId` | TypeString | No | No | No | 规则归属成员账号ID。入参为 string，出参为 BelongMember 对象，回读取 BelongMember.MemberId 仅非空才写入 |
| `rule_id` | `Rules.RuleId` | `Rules.RuleId` | TypeString | — | — | Yes | 规则ID，Create 不传，Modify 必填，回读设置 |
| `sequence` | — | `Rules.Sequence` | TypeInt | — | — | Yes | 优先级（只读回显，与 order_index 同源） |
| `source_name` | — | `Rules.SourceName` | TypeString | — | — | Yes | 源资产名称 |
| `dest_name` | — | `Rules.DestName` | TypeString | — | — | Yes | 目的资产名称 |
| `param_template_name` | — | `Rules.ParamTemplateName` | TypeString | — | — | Yes | 端口模板名称 |
| `fw_group_name` | — | `Rules.FwGroupName` | TypeString | — | — | Yes | 防火墙组名称 |
| `invalid` | — | `Rules.Invalid` | TypeInt | — | — | Yes | 规则是否失效：0有效 1失效 |
| `app_id` | — | `Rules.BelongMember.AppId` | TypeString | — | — | Yes | 成员 AppId（从 BelongMember 展开） |
| `uin` | — | `Rules.BelongMember.Uin` | TypeString | — | — | Yes | 成员 Uin |
| `nickname` | — | `Rules.BelongMember.Nickname` | TypeString | — | — | Yes | 成员昵称 |
| `member_id` | — | `Rules.BelongMember.MemberId` | TypeString | — | — | Yes | 成员 ID（从 BelongMember 展开） |
| `node_name` | — | `Rules.BelongMember.NodeName` | TypeString | — | — | Yes | 所属部门 |
| `create_time` | — | `Rules.CreateTime` | TypeString | — | — | Yes | 规则创建时间 |
| `update_time` | — | `Rules.UpdateTime` | TypeString | — | — | Yes | 规则修改时间 |
| `dns_parse_cnt` | — | `Rules.DnsParseCnt` | TypeInt | — | — | Yes | 域名数 |
| `hit_count` | — | `Rules.HitCount` | TypeInt | — | — | Yes | 命中次数 |
| `last_hit_time` | — | `Rules.LastHitTime` | TypeString | — | — | Yes | 最近命中时间 |
| `all_total_count` | — | `Response.AllTotalCount` | TypeInt | — | — | Yes | 不过滤的总数 |

**关于 `ip_version` 的说明**：Create 入参 VpcAclRule 结构体中没有 IpVersion 字段（服务端按规则组配置推断 IP 版本），但 Describe 出参 VpcAclRuleResp 有 IpVersion。用户自定义要求标 required + forceNew。实现时 Create 不传 IpVersion（SDK 结构体无此字段），仅作为 schema 字段存在并在 Read 时回读。由于 forceNew，变更 ip_version 会触发重建。这是符合业务语义的：VPC 边界规则用 IpVersion 区分 ipv4/ipv6，不同版本的规则属于不同语义实例。

**关于扁平化 schema**：用户自定义要求要求规则字段平铺到顶层（不嵌套 `rules` 列表层），这与 cfw_vpc_policy 的扁平化模式一致。Create 时从顶层字段组装成 `Rules` 数组（1 元素），Update 时组装成 `Rule` 单对象。

### D2 — ForceNew 集合

以下字段变更视为重建（forceNew）：`group_id`、`edge_id`、`fw_group_id`、`ip_version`、`source_type`、`source_content`、`dest_type`、`dest_content`。

理由：这些字段决定规则的「身份/生效范围/地址语义」，变更等同规则迁移（用户自定义要求第 8、18 条）。其余可写字段（`protocol`、`rule_action`、`order_index`、`description`、`port`、`param_template_id`、`belong_member_id`）为原地更新。

### D3 — Update 全量覆盖

ModifyVpcAclRule 是全量覆盖，复用与 Create 完全相同的字段校验（checkVpcAclUserInput），因此 `source_content`、`source_type`、`dest_content`、`dest_type`、`protocol`、`rule_action`、`description`、`edge_id`、`fw_group_id` 这些字段在 Update 时全部必填。

但由于 edge_id、fw_group_id、source_type、source_content、dest_type、dest_content 是 forceNew，实际不会进入 Update 路径。Update 只处理原地更新字段：`protocol`、`rule_action`、`order_index`、`description`、`port`、`param_template_id`、`belong_member_id`。

Update 实现：从 state 完整读取所有可写字段（包括 forceNew 字段的当前值，因为 Modify 要求全量），组装成 `Rule` 单对象（含 `RuleId`），传 `GroupId`，调用 ModifyVpcAclRule。使用 `mutableArgs` 模式检测变更，只要任一原地更新字段变更就发起 Modify 调用，全量回传所有字段。

### D4 — Create/Modify 结构不对称

Create 传 `Rules` 数组（1 元素，`[]*VpcAclRule`），Modify 传 `Rule` 单对象（`*VpcAclRule`）。这是接口既有设计，不统一。资源 schema 里两者都映射到同一组扁平字段，由生成器分别组装。

### D5 — Read 查询路径

无单实例 Get 接口。用 `DescribeVpcAclRules`：
- `request.GroupId = groupId`
- `request.Filters = [{Name: "RuleId", Values: [ruleId], OperatorType: 1}]`（1 = EQUAL 精确匹配）
- `request.Limit = 0`（取默认 100）
- `request.Offset = 0`

取 `response.Response.Rules[0]`。若 Rules 为空，说明规则被外部删除，`d.SetId("")` 移除 state。

回读字段映射要点：
- `order_index` ← `Rules[0].Sequence`（不是 OrderIndex，出参无 OrderIndex）
- `belong_member_id` ← `Rules[0].BelongMember.MemberId`（BelongMember 是 MemberInfo 对象，仅非空才写）
- `app_id`/`uin`/`nickname`/`member_id`/`node_name` ← `Rules[0].BelongMember.*`（从对象展开，仅非空才写）
- `all_total_count` ← `response.Response.AllTotalCount`

### D6 — Delete 幂等

DeleteVpcAclRule 是批量接口，单个删除时 `RuleIds` 传 1 个元素。RuleId 已不存在返回 ResourceNotFound，destroy 视为删除成功（幂等）：在 service 层 `DeleteFwmVpcAclRuleById` 中捕获 ResourceNotFound 错误码，返回 nil。

### D7 — 复合 ID 格式

`d.SetId(strings.Join([]string{groupId, ruleId}, tccommon.FILED_SP))`，即 `GroupId#RuleId`。

在 Read/Update/Delete 中拆分：`idSplit := strings.Split(d.Id(), tccommon.FILED_SP)`，校验 `len(idSplit) == 2`，否则返回 `fmt.Errorf("id is broken,%s", d.Id())`。

import 支持 `terraform import tencentcloud_fwm_vpc_acl_rule.example GroupId#RuleId`。

### D8 — 不生成 _extension.go 文件

按用户自定义要求第 10 条，若非必须不生成 extension 文件。本资源无需 extension。

## SDK 结构定义参考（从 vendor 摘录）

### client.go 访问器签名（line 624-648）

```go
func (c *Client) CreateVpcAclRule(request *CreateVpcAclRuleRequest) (response *CreateVpcAclRuleResponse, err error) {
    return c.CreateVpcAclRuleWithContext(context.Background(), request)
}
func (c *Client) CreateVpcAclRuleWithContext(ctx context.Context, request *CreateVpcAclRuleRequest) (response *CreateVpcAclRuleResponse, err error)
func (c *Client) DescribeVpcAclRules(request *DescribeVpcAclRulesRequest) (response *DescribeVpcAclRulesResponse, err error)
func (c *Client) DescribeVpcAclRulesWithContext(ctx context.Context, request *DescribeVpcAclRulesRequest) (response *DescribeVpcAclRulesResponse, err error)
func (c *Client) ModifyVpcAclRule(request *ModifyVpcAclRuleRequest) (response *ModifyVpcAclRuleResponse, err error)
func (c *Client) ModifyVpcAclRuleWithContext(ctx context.Context, request *ModifyVpcAclRuleRequest) (response *ModifyVpcAclRuleResponse, err error)
func (c *Client) DeleteVpcAclRule(request *DeleteVpcAclRuleRequest) (response *DeleteVpcAclRuleResponse, err error)
func (c *Client) DeleteVpcAclRuleWithContext(ctx context.Context, request *DeleteVpcAclRuleRequest) (response *DeleteVpcAclRuleResponse, err error)
```

### models.go — CreateVpcAclRuleRequest (line 1033-1041)

```go
type CreateVpcAclRuleRequest struct {
	*tchttp.BaseRequest
	// 规则组ID
	GroupId *string `json:"GroupId,omitnil,omitempty" name:"GroupId"`
	// 规则列表
	Rules []*VpcAclRule `json:"Rules,omitnil,omitempty" name:"Rules"`
}
```

### models.go — CreateVpcAclRuleResponseParams (line 1064-1070)

```go
type CreateVpcAclRuleResponseParams struct {
	// 规则Id列表
	RuleIds []*string `json:"RuleIds,omitnil,omitempty" name:"RuleIds"`
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}
type CreateVpcAclRuleResponse struct {
	*tchttp.BaseResponse
	Response *CreateVpcAclRuleResponseParams `json:"Response"`
}
```

### models.go — VpcAclRule (line 4507-4549，Create/Modify 入参用)

```go
type VpcAclRule struct {
	SourceContent  *string `json:"SourceContent,omitnil,omitempty"`
	SourceType     *string `json:"SourceType,omitnil,omitempty"`
	DestContent    *string `json:"DestContent,omitnil,omitempty"`
	DestType       *string `json:"DestType,omitnil,omitempty"`
	Protocol       *string `json:"Protocol,omitnil,omitempty"`
	RuleAction     *string `json:"RuleAction,omitnil,omitempty"`
	OrderIndex     *int64  `json:"OrderIndex,omitnil,omitempty"`   // 优先级（从1开始）
	EdgeId         *string `json:"EdgeId,omitnil,omitempty"`       // ALL 或 CFWS-xxx
	FwGroupId      *string `json:"FwGroupId,omitnil,omitempty"`    // ALL/ccn-xxx/cfwg-xxx
	RuleId         *string `json:"RuleId,omitnil,omitempty"`       // 修改时必填
	Port           *string `json:"Port,omitnil,omitempty"`         // ICMP 时为空
	Description    *string `json:"Description,omitnil,omitempty"`
	ParamTemplateId *string `json:"ParamTemplateId,omitnil,omitempty"`
	BelongMemberId *string `json:"BelongMemberId,omitnil,omitempty"`
}
```

注意：VpcAclRule 入参结构体**无 IpVersion 字段**。

### models.go — VpcAclRuleResp (line 4551-4633，Describe 出参用)

```go
type VpcAclRuleResp struct {
	RuleId           *string      `json:"RuleId,omitnil,omitempty"`
	Sequence         *int64       `json:"Sequence,omitnil,omitempty"`         // 优先级（对应入参 OrderIndex）
	IpVersion        *string      `json:"IpVersion,omitnil,omitempty"`        // ipv4 或 ipv6
	SourceContent    *string      `json:"SourceContent,omitnil,omitempty"`
	SourceType       *string      `json:"SourceType,omitnil,omitempty"`
	SourceName       *string      `json:"SourceName,omitnil,omitempty"`       // 可 null
	DestContent      *string      `json:"DestContent,omitnil,omitempty"`
	DestType         *string      `json:"DestType,omitnil,omitempty"`
	DestName         *string      `json:"DestName,omitnil,omitempty"`          // 可 null
	Protocol         *string      `json:"Protocol,omitnil,omitempty"`
	Port             *string      `json:"Port,omitnil,omitempty"`              // 可 null
	ParamTemplateId  *string      `json:"ParamTemplateId,omitnil,omitempty"`  // 可 null
	ParamTemplateName *string     `json:"ParamTemplateName,omitnil,omitempty"`// 可 null
	RuleAction       *string      `json:"RuleAction,omitnil,omitempty"`
	Description      *string      `json:"Description,omitnil,omitempty"`      // 可 null
	EdgeId           *string      `json:"EdgeId,omitnil,omitempty"`
	FwGroupId        *string      `json:"FwGroupId,omitnil,omitempty"`
	Invalid          *int64       `json:"Invalid,omitnil,omitempty"`          // 0有效 1失效
	BelongMember     *MemberInfo  `json:"BelongMember,omitnil,omitempty"`     // 可 null，对象
	CreateTime       *string      `json:"CreateTime,omitnil,omitempty"`
	UpdateTime       *string      `json:"UpdateTime,omitnil,omitempty"`
	DnsParseCnt      *int64       `json:"DnsParseCnt,omitnil,omitempty"`
	FwGroupName      *string      `json:"FwGroupName,omitnil,omitempty"`
	HitCount         *int64       `json:"HitCount,omitnil,omitempty"`
	LastHitTime      *string      `json:"LastHitTime,omitnil,omitempty"`      // 从未命中为 "-"
}
```

### models.go — MemberInfo (line 3049-3064)

```go
type MemberInfo struct {
	AppId    *string `json:"AppId,omitnil,omitempty"`
	Uin      *string `json:"Uin,omitnil,omitempty"`
	Nickname *string `json:"Nickname,omitnil,omitempty"`
	MemberId *string `json:"MemberId,omitnil,omitempty"`   // 个人场景为空
	NodeName *string `json:"NodeName,omitnil,omitempty"`
}
```

### models.go — ModifyVpcAclRuleRequest (line 3605-3613)

```go
type ModifyVpcAclRuleRequest struct {
	*tchttp.BaseRequest
	GroupId *string     `json:"GroupId,omitnil,omitempty"`
	Rule    *VpcAclRule `json:"Rule,omitnil,omitempty"`   // 注意：单对象，非数组
}
type ModifyVpcAclRuleResponseParams struct {
	RequestId *string `json:"RequestId,omitnil,omitempty"`
}
```

### models.go — DescribeVpcAclRulesRequest (line 2702-2716)

```go
type DescribeVpcAclRulesRequest struct {
	*tchttp.BaseRequest
	GroupId *string          `json:"GroupId,omitnil,omitempty"`
	Filters []*CommonFilter  `json:"Filters,omitnil,omitempty"`
	Offset  *int64           `json:"Offset,omitnil,omitempty"`
	Limit   *int64           `json:"Limit,omitnil,omitempty"`
}
type DescribeVpcAclRulesResponseParams struct {
	Rules        []*VpcAclRuleResp `json:"Rules,omitnil,omitempty"`
	TotalCount   *int64            `json:"TotalCount,omitnil,omitempty"`
	AllTotalCount *int64           `json:"AllTotalCount,omitnil,omitempty"`
	RequestId    *string           `json:"RequestId,omitnil,omitempty"`
}
```

### models.go — CommonFilter (line 321-333)

```go
type CommonFilter struct {
	Name         *string   `json:"Name,omitnil,omitempty"`         // 可 null
	Values       []*string `json:"Values,omitnil,omitempty"`       // 可 null
	OperatorType *int64    `json:"OperatorType,omitnil,omitempty"` // 1=精确匹配 7=in 9=模糊
}
```

### models.go — DeleteVpcAclRuleRequest (line 1388-1396)

```go
type DeleteVpcAclRuleRequest struct {
	*tchttp.BaseRequest
	GroupId *string   `json:"GroupId,omitnil,omitempty"`
	RuleIds []*string `json:"RuleIds,omitnil,omitempty"`
}
type DeleteVpcAclRuleResponseParams struct {
	RequestId *string `json:"RequestId,omitnil,omitempty"`
}
```

## 代码风格参照说明

### tencentcloud_igtm_strategy（复合 ID + CRUD 骨架参考）

- 资源结构：`Create/Read/Update/Delete` 回调 + `Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough}`。
- 复合 ID：`d.SetId(strings.Join([]string{instanceId, strategyId}, tccommon.FILED_SP))`，Read/Update/Delete 中 `idSplit := strings.Split(d.Id(), tccommon.FILED_SP)` + `len(idSplit) != 2` 校验。
- Create：`defer tccommon.LogElapsed("resource.tencentcloud_igtm_strategy.create")()` + `defer tccommon.InconsistentCheck(d, meta)()`；retry 块内校验 `result == nil || result.Response == nil` → NonRetryableError；拿到 id 后 `d.SetId(...)`，最后 `return resourceTencentCloudIgtmStrategyRead(d, meta)`。
- Read：service 层 `DescribeIgtmStrategyById(ctx, ...)`，返回 nil 则 `log.Printf("[WARN]...")` + `d.SetId("")` + return nil；每个字段 `if respData.X != nil { _ = d.Set("x", respData.X) }`。
- Update：`mutableArgs` 数组 + `d.HasChange(v)` 检测；变更则全量组装 request 调 Modify，retry 块包装；最后 `return resourceTencentCloudIgtmStrategyRead(d, meta)`。
- Delete：拆分 ID，service 层方法包装。
- service 层：`NewIgtmService(client)` + `IgtmService{client}`，方法内 `resource.Retry(tccommon.ReadRetryTimeout, ...)` + `ratelimit.Check(request.GetAction())` + `tccommon.RetryError(e)`。

### tencentcloud_cfw_vpc_policy（业务模板参考）

- 扁平化 schema：规则字段全部平铺到顶层（source_content、dest_content、protocol 等），不嵌套 rules 列表层。
- Create：`vpcRuleItem := cfw.VpcRuleItem{}` 组装单元素，`request.Rules = append(request.Rules, &vpcRuleItem)`，取 `response.Response.RuleUuids[0]`。
- Update：全量覆盖，从 state 读取所有字段组装 `vpcRuleItem`，`request.Rules = append(request.Rules, &vpcRuleItem)`。
- Read：字段逐一 `if vpcPolicy.X != nil { _ = d.Set("x", vpcPolicy.X) }`。

### tencentcloud_cfw_cluster_nat_fw_switch（gomonkey mock 测试参考）

- `package cfw_test`。
- mockMeta 结构体实现 `tccommon.ProviderMeta`，`GetAPIV3Conn()` 返回 mock client。
- `patches := gomonkey.NewPatches()` + `defer patches.Reset()`。
- `patches.ApplyMethodReturn(mockClient, "UseCfwV20190904Client", cfwClient)` mock client 访问器。
- `patches.ApplyMethodFunc(cfwClient, "CreateXxxWithContext", func(...) {...})` mock 具体 API。
- `schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{...})` 构造 ResourceData。
- `res.Create(d, meta)` / `res.Read(d, meta)` / `res.Update(d, meta)` / `res.Delete(d, meta)` 调用。
- `assert.NoError(t, err)` + `assert.Equal(t, expectedId, d.Id())` 断言。

## connectivity/client.go 改动说明

需在 `tencentcloud/connectivity/client.go` 新增：

1. import（按字母序插入）：
```go
fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"
```

2. struct 字段（line 268 附近，与其他 Conn 字段并列）：
```go
fwmv20250611Conn            *fwmv20250611.Client
```

3. 访问器方法（参照 UseIgtmV20231024Client 模式）：
```go
// UseFwmV20250611Client return FWM client for service
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

## provider.go 改动说明

1. import（按字母序插入，约 line 84 附近）：
```go
"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/fwm"
```

2. 资源注册映射（与其他资源并列，按字母序插入）：
```go
"tencentcloud_fwm_vpc_acl_rule": fwm.ResourceTencentCloudFwmVpcAclRule(),
```