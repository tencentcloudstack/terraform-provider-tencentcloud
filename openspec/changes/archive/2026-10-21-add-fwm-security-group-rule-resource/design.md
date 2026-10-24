## Context

腾讯云防火墙新一代产品族 FWM（SDK 包 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611`）提供「企业安全组规则」能力。本资源 `tencentcloud_fwm_security_group_rule` 管理一条企业安全组规则的完整生命周期，规则隶属于一个已存在的规则组（Product=enterprise_sg，由独立资源 `fwm_rule_group` 创建，本资源不创建规则组）。

四个云 API 接口**全部同步**（无 FlowId/TaskId/JobId、无中间态），Create/Modify/Delete 成功即已生效，不做轮询等待。Key SDK 类型：

- `CreateSecurityGroupRuleRequest`：`GroupId *string` + `Rules []*SecurityGroupRule`（批量数组，本资源固定 1 元素）。响应 `CreateSecurityGroupRuleResponseParams.RuleIds []*string`（与入参 Rules 按位置一一对应，单实例取 `RuleIds[0]`）。
- `DescribeSecurityGroupRulesRequest`：`GroupId *string` + `Filters []*CommonFilter` + `Offset/Limit *int64`（列表查询，无单实例 Get）。响应 `DescribeSecurityGroupRulesResponseParams.Rules []*SecGroupRuleResp` + `TotalCount/AllTotalCount *int64`。本资源 Read 用 `Filters=[{Name:"RuleId",Values:[ruleId],OperatorType:1}]`（1=EQUAL 精确匹配）精准查单条，取 `Rules` 首元素。
- `ModifySecurityGroupRuleRequest`：`GroupId *string` + `Rule *SecurityGroupRule`（单条，全量覆盖语义）。响应仅 `RequestId`。
- `DeleteSecurityGroupRuleRequest`：`GroupId *string` + `RuleIds []*string`（批量，本资源传 1 元素）。响应仅 `RequestId`。RuleId 已不存在返回 `ResourceNotFound`，destroy 视为删除成功（幂等）。
- 共用入参结构 `SecurityGroupRule`：`IpVersion/SourceContent/SourceType/DestContent/DestType/Protocol/Port/ServiceTemplateId/RuleAction/Description/RuleId/Scope *string`、`OrderIndex/ProtocolPortType *int64`、`BelongMemberId *string`。
- 共用响应结构 `SecGroupRuleResp`：`OrderIndex *int64`、`RuleId/IpVersion/SourceId/TargetId/Protocol/Port/Detail/Region/ServiceTemplateId/SouInstanceName/SouPublicIp/SouPrivateIp/SouCidr/SouParameterName/InstanceName/PublicIp/PrivateIp/Cidr/ParameterName/ProtocolPortName/Scope *string`、`SourceType/TargetType/Strategy/Id/IsNew *int64`、`DnsParseCount *SgDnsParseCount`、`BelongMember *MemberInfo`。
- `CommonFilter`：`Name/Values *string` / `[]*string` / `OperatorType *int64`（1=精确,7=in,9=模糊）。
- `SgDnsParseCount`：`ValidCount/InvalidCount *int64`。`MemberInfo`：`AppId/Uin/Nickname/MemberId/NodeName *string`。

**关键映射差异（入参 string ↔ 响应 int，必须双向映射否则 plan 永久 diff）**：

| 入参(HCL 字符串) | 响应字段(int) | 映射关系 |
|---|---|---|
| `source_type` / `dest_type` | `SourceType` / `TargetType` | net=0, vpc=1, subnet=2, cvm=3, clb=4, eni=5, cdb=6, template=7, tag=8, region=9, resourcegroup=100；入参额外支持 addressbook/instance（无对应 int，回读时需处理） |
| `rule_action` | `Strategy` | accept=2, drop=1, log=0 |
| `description` | `Detail` | 直接对应 |
| `source_content` | `SourceId` | 直接对应 |
| `dest_content` | `TargetId` | 直接对应 |
| `belong_member_id` | `BelongMember.MemberId` | 仅非空才写入 |
| `port` / `protocol` / `scope` / `ip_version` | `Port` / `Protocol` / `Scope` / `IpVersion` | 同名直接对应 |

代码风格参照 `tencentcloud/services/igtm/resource_tc_igtm_strategy.go`（通用资源 CRUD 骨架、复合 ID `tccommon.FILED_SP`、`resource.Retry` 包装、`defer LogElapsed/InconsistentCheck`）与 `tencentcloud/services/igtm/service_tencentcloud_igtm.go`（service 层 `DescribeXxxById` 方法）。FWM 客户端注册参照 `tencentcloud/connectivity/client.go` 的 `UseIgtmV20231024Client` / `UseGa2V20250115Client` 模式（import + 字段 + accessor）。

## Goals / Non-Goals

**Goals:**

- Schema 平铺：入参字段定义在资源顶层，不嵌套 `rules`/`rule` 一层。Create 构造 1 元素 `Rules` 数组，Modify 构造 `Rule` 对象。
- 资源 ID 复合 `GroupId<tccommon.FILED_SP>RuleId`；import 支持 `GroupId:RuleId`。
- 入参/响应双向类型映射，Read 用 int 反查成字符串写入 state，避免 plan 永久 diff。
- 只读展示字段均为 Computed，不参与可比 diff。
- 四接口全同步，不生成轮询。
- Read 遇空（规则被外部删除）打印 `[CRUD]` 日志后 `d.SetId("")` 移除 state 不报错。
- Delete 幂等：遇 `ResourceNotFound` 视为删除成功。
- 新增 FWM 产品客户端注册（connectivity/client.go）与 Provider 注册（provider.go / provider.md）。
- 单元测试使用 gomonkey mock 云 API（新增资源用 mock，不用 TF 测试套件）。

**Non-Goals:**

- 不实现规则组资源（`fwm_rule_group`）—— 单独 change；本资源前提是规则组已存在且 Product=enterprise_sg。
- 不实现配套数据源（`tencentcloud_fwm_security_group_rules`）。
- 不回读 `Id`（int64 自增主键）作可比字段；`Id` 标 Computed（存在溢出风险，state/import 只用 `RuleId`）。
- 不做 `ProtocolPortType` 的强枚举校验。
- 不生成 `_extension.go` 文件。
- 不修改任何既有资源/数据源/service 方法。

## Decisions

### D1 — Schema 字段映射与约束

平铺 Schema（顶层 key → SDK 字段）：

| HCL key | SDK 入参(Create/Modify) | SDK 响应(Describe) | Type | Required | ForceNew | Computed | 说明 |
|---|---|---|---|---|---|---|---|
| `group_id` | `request.GroupId` | — | TypeString | Yes | No | No | 父规则组引用，四接口必填 |
| `ip_version` | `Rules/Rule.IpVersion` | `Rules.IpVersion` | TypeString | Yes | **Yes** | No | ipv4/ipv6 |
| `order_index` | `Rules/Rule.OrderIndex` | `Rules.OrderIndex` | TypeInt | Yes | No | No | 全局单一序列，>=1，传0报错 |
| `source_type` | `Rules/Rule.SourceType` | `Rules.SourceType`(int) | TypeString | Yes | **Yes** | No | 入参string，响应int反查 |
| `source_content` | `Rules/Rule.SourceContent` | `Rules.SourceId` | TypeString | Yes | No | No | 入参SourceContent↔响应SourceId |
| `dest_type` | `Rules/Rule.DestType` | `Rules.TargetType`(int) | TypeString | Yes | **Yes** | No | 入参DestType(int映射)，响应TargetType(int)反查 |
| `dest_content` | `Rules/Rule.DestContent` | `Rules.TargetId` | TypeString | Yes | No | No | 入参DestContent↔响应TargetId |
| `protocol` | `Rules/Rule.Protocol` | `Rules.Protocol` | TypeString | Yes | No | No | TCP/ANY/ICMP/UDP |
| `port` | `Rules/Rule.Port` | `Rules.Port` | TypeString | Yes | No | No | |
| `rule_action` | `Rules/Rule.RuleAction` | `Rules.Strategy`(int) | TypeString | Yes | No | No | accept/drop/log ↔ int |
| `scope` | `Rules/Rule.Scope` | `Rules.Scope` | TypeString | No | **Yes** | No | SG(缺省)/LH |
| `service_template_id` | `Rules/Rule.ServiceTemplateId` | `Rules.ServiceTemplateId` | TypeString | No | No | Yes | 端口模板ID |
| `description` | `Rules/Rule.Description` | `Rules.Detail` | TypeString | No | No | No | 入参Description↔响应Detail |
| `protocol_port_type` | `Rules/Rule.ProtocolPortType` | — | TypeInt | No | No | No | bool/int用GetOkExists |
| `belong_member_id` | `Rules/Rule.BelongMemberId` | `Rules.BelongMember.MemberId` | TypeString | No | No | No | 仅非空写入 |
| `rule_id` | — | `Rules.RuleId` | TypeString | — | — | Yes | 读 |
| `source_id` | — | `Rules.SourceId` | TypeString | — | — | Yes | 只读展示 |
| `target_id` | — | `Rules.TargetId` | TypeString | — | — | Yes | 只读展示 |
| `strategy` | — | `Rules.Strategy` | TypeInt | — | — | Yes | int原始值 |
| `detail` | — | `Rules.Detail` | TypeString | — | — | Yes | 只读展示 |
| `region` | — | `Rules.Region` | TypeString | — | — | Yes | 只读展示 |
| `protocol_port_name` | — | `Rules.ProtocolPortName` | TypeString | — | — | Yes | 只读展示 |
| `sou_instance_name` | — | `Rules.SouInstanceName` | TypeString | — | — | Yes | 只读展示 |
| `sou_public_ip` | — | `Rules.SouPublicIp` | TypeString | — | — | Yes | 只读展示 |
| `sou_private_ip` | — | `Rules.SouPrivateIp` | TypeString | — | — | Yes | 只读展示 |
| `sou_cidr` | — | `Rules.SouCidr` | TypeString | — | — | Yes | 只读展示 |
| `sou_parameter_name` | — | `Rules.SouParameterName` | TypeString | — | — | Yes | 只读展示 |
| `instance_name` | — | `Rules.InstanceName` | TypeString | — | — | Yes | 只读展示 |
| `public_ip` | — | `Rules.PublicIp` | TypeString | — | — | Yes | 只读展示 |
| `private_ip` | — | `Rules.PrivateIp` | TypeString | — | — | Yes | 只读展示 |
| `cidr` | — | `Rules.Cidr` | TypeString | — | — | Yes | 只读展示 |
| `parameter_name` | — | `Rules.ParameterName` | TypeString | — | — | Yes | 只读展示 |
| `dns_parse_count` | — | `Rules.DnsParseCount` | TypeList | — | — | Yes | 嵌套valid_count/invalid_count |
| `is_new` | — | `Rules.IsNew` | TypeInt | — | — | Yes | 只读展示 |
| `id` | — | `Rules.Id` | TypeInt | — | — | Yes | int64内部主键，有溢出风险，不用作state |
| `all_total_count` | — | `AllTotalCount` | TypeInt | — | — | Yes | 只读展示 |

`SourceContent`/`DestContent` 为入参字段，响应里对应 `SourceId`/`TargetId`：写入 state 时 `source_content` 写入参原值（用户配置保持不变，可比字段），同时 `source_id`/`target_id` 作只读展示。`source_type`/`dest_type`/`rule_action` 回读时用 int 反查成入参字符串写回（保证 plan 不 diff）。

### D2 — SourceType/TargetType int ↔ string 映射表

```
int 0  → net            string "net"          → int 0
int 1  → vpc            string "vpc"          → int 1
int 2  → subnet         string "subnet"       → int 2
int 3  → cvm            string "cvm"          → int 3
int 4  → clb            string "clb"          → int 4
int 5  → eni            string "eni"          → int 5
int 6  → cdb            string "cdb"          → int 6
int 7  → template       string "template"     → int 7
int 8  → tag            string "tag"          → int 8
int 9  → region         string "region"       → int 9
int 100→ resourcegroup  string "resourcegroup"→ int 100
```

入参 `source_type`/`dest_type` 额外支持 `addressbook`、`instance`（用户自定义要求第8条枚举），这些无对应 int 值；Create/Modify 时直接按字符串传入（SDK 字段为 `*string`，透传服务端）。Read 回读时若响应 int 不在映射表内，保留 int 字符串并记录 warning，避免回读 panic。`rule_action` 映射：accept=2, drop=1, log=0。

### D3 — Create：单实例批量接口取 RuleIds[0]

Create 构造 `Rules []*SecurityGroupRule` 固定 1 元素（所有入参字段填入），调用 `CreateSecurityGroupRuleWithContext`。响应 `RuleIds` 数组与入参位置一一对应，单实例取 `RuleIds[0]` 作 `RuleId`。`d.SetId(GroupId + FILED_SP + RuleId)`。检查 `response == nil || response.Response == nil || len(response.Response.RuleIds) == 0 || response.Response.RuleIds[0] == nil || *response.Response.RuleIds[0] == ""` → `NonRetryableError`（打印 logId 与 d.Id()）。Create 无需 SetId 前的等待（同步接口）。

### D4 — Read：列表查询精准查单条

Read 拆 `d.Id()` → `GroupId`、`RuleId`。构造 `DescribeSecurityGroupRulesRequest`：`GroupId`、`Filters=[{Name:"RuleId",Values:[ruleId],OperatorType:1}]`、`Limit` 取默认值（不传0/负数）。包 `resource.Retry(ReadRetryTimeout)`，retry 内检查空响应 → `NonRetryableError`（按通用资源 Read 模式，不直接 SetId）。retry 外检查 `len(Rules) == 0`：视为规则被外部删除，`log.Printf("[CRUD] fwm_security_group_rule id=%s", d.Id())` 后 `d.SetId("")` 返回 nil。取 `Rules[0]` 回写字段（int 反查字符串）。

注意：service 层 `DescribeFwmSecurityGroupRuleById` 返回 `*SecGroupRuleResp`（首元素），由资源 Read 处理空 → SetId。

### D5 — Update：全量覆盖 Modify

Update 用 `d.HasChange` 检测 `mutableArgs`（`order_index`/`source_content`/`dest_type`/`dest_content`/`protocol`/`port`/`service_template_id`/`rule_action`/`description`/`protocol_port_type`/`belong_member_id`）。命中则构造 `ModifySecurityGroupRuleRequest`：`GroupId`（从 id 拆）、`Rule *SecurityGroupRule`（含 `RuleId` + 所有可写字段全程重传，覆盖语义）。包 `resource.Retry(WriteRetryTimeout)`，检查空响应。`ip_version`/`source_type`/`target_type`/`scope` 为 ForceNew，不会进 Update。Group_id 不可改（不可修改 GroupId）。

### D6 — Delete：删除

Delete 拆 id → `GroupId`、`RuleId`。构造 `DeleteSecurityGroupRuleRequest`：`GroupId`、`RuleIds=[ruleId]`（1 元素）。包 `resource.Retry(WriteRetryTimeout)`，遇 `ResourceNotFound` 视为删除成功（幂等），返回 nil；其余错误 `tccommon.RetryError(e)` 重试。

### D7 — FWM 客户端注册

`tencentcloud/connectivity/client.go`：
- import 段新增 `fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"`。
- `TencentCloudClient` struct 新增字段 `fwmv20250611Conn *fwmv20250611.Client`。
- 新增访问器 `UseFwmV20250611Client() *fwmv20250611.Client`（懒初始化，与 `UseIgtmV20231024Client` 完全同款）。

`tencentcloud/provider.go`：import fwm 服务包 + 注册 `"tencentcloud_fwm_security_group_rule": fwm.ResourceTencentCloudFwmSecurityGroupRule()`。`tencentcloud/provider.md`：新增 Firewall(FWM) 段与资源名。

### D8 — 文档与测试

- `.md`：一句话描述带 FWM 产品名；Example Usage 用 jsonencode（本资源无 json 字段，用普通 HCL）；Import 段说明使用 `GroupId:RuleId` 复合 id；不含 Argument/Attribute Reference（工具生成）。
- `_test.go`：gomonkey mock `CreateSecurityGroupRuleWithContext`/`DescribeSecurityGroupRulesWithContext`/`ModifySecurityGroupRuleWithContext`/`DeleteSecurityGroupRuleWithContext`，测 Create/Read/Update/Delete 业务逻辑，包名 `fwm_test`。testValue：ip_version=ipv4、order_index=1、source_type=net、source_content=1.1.1.1、dest_type=net、dest_content=2.2.2.2、protocol=TCP、port=80、rule_action=accept、scope=SG。

### D9 — import 与复合 ID

import `schema.ImportStatePassthrough`，用户 import 时传 `GroupId:RuleId`。注意 import 用 `:` 分隔（用户输入），但本资源 `d.Id()` 内部存储用 `tccommon.FILED_SP` 分隔——参照 igtm_strategy（igtm md 文档 import 示例用 `#` 但 state 内 FILED_SP）。这里 import 字符串直接成为 `d.Id()`，因此需要在 Read 中兼容：用 `strings.Split(d.Id(), tccommon.FILED_SP)`，若失败再尝试 `:` 分割。为简化，文档 import 示例统一用 `tccommon.FILED_SP` 不便（用户不可见字符）。实际参照 igtm：文档写 `GroupId#RuleId`（igtm 用 `#`）。本资源统一约定 import 与内部存储都用 `tccommon.FILED_SP`，并在 md 文档 import 示例使用 `GroupId#RuleId` 形式（`#` 即 FILED_SP 实际值需确认）。

→ **决策**：待确认 `tccommon.FILED_SP` 实际值。若为 `#`，则 import 文档用 `#`；若为其他字符，文档 import 示例需与之一致。实施时先 grep `FILED_SP` 常量值，按实际值写 import 示例（参照 igtm_strategy 的 `terraform import tencentcloud_igtm_strategy.example gtm-uukztqtoaru#7556`）。

## Risks / Trade-offs

- **Risk**：`SourceType`/`TargetType` 入参支持字符串值（addressbook/instance）在响应端无对应 int，回读反查表查不到 → Mitigation：反查失败时保留 int 的字符串形式写回，并在日志记录 warning，不 panic；但此种情况 plan 可能出现 diff，需在 md 文档 Remark 说明这些类型回读展示为 int 值。
- **Risk**：企业安全组规则 OrderIndex 全局单一序列，同组内多条规则 OrderIndex 重复或乱序由服务端校验，Terraform 侧不保证 → Mitigation：md 文档说明用户需自行保证同组内 OrderIndex 严格递增且 >=1。
- **Risk**：`Id`（int64 自增主键）有溢出风险，不进 state → Mitigation：仅 `rule_id`（uuid）作 state 标识，`id` 标 Computed 只展示。
- **Trade-off**：Read 用列表查询 + Filter 精准查单条（无单实例 Get 接口），多一次网络往返 → 可接受，Filter 精确匹配 RuleId 命中唯一。
- **Trade-off**：不实现轮询（四接口同步） → 简化实现，符合接口语义；若后续接口改异步再补。