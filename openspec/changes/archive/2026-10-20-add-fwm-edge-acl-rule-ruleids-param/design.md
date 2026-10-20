## Context

`tencentcloud_fwm_edge_acl_rule` 是管理腾讯云防火墙（fwm）互联网边界 ACL 规则的 RESOURCE_KIND_GENERAL 资源，对应云 API 包 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611`。

创建规则调用 `CreateEdgeAclRule` 接口，该接口会返回新创建规则对应的规则 ID 列表（`response.RuleIds`），但当前资源 schema 中没有对应字段来承接该出参，导致创建出来的规则 ID 无法落入 Terraform state、也无法被配置间引用。

本次变更仅为该资源新增一个 **出参** `RuleIds`，不涉及入参变更、不修改既有字段。

**云 API vendor 摘录**（来源：`vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611/models.go`）

`CreateEdgeAclRuleRequest`（入参，不受本次变更影响）：
```go
// models.go:471-487
type CreateEdgeAclRuleRequestParams struct {
	// <p>规则组ID</p>
	GroupId *string `json:"GroupId,omitnil,omitempty" name:"GroupId"`
	// <p>规则列表</p>
	Rules []*EdgeAclRuleInfo `json:"Rules,omitnil,omitempty" name:"Rules"`
}
type CreateEdgeAclRuleRequest struct {
	*tchttp.BaseRequest
	GroupId *string `json:"GroupId,omitnil,omitempty" name:"GroupId"`
	Rules   []*EdgeAclRuleInfo `json:"Rules,omitnil,omitempty" name:"Rules"`
}
```

`CreateEdgeAclRuleResponse`（出参，本次新增字段来源）：
```go
// models.go:510-521
type CreateEdgeAclRuleResponseParams struct {
	// <p>规则Id列表</p>
	RuleIds []*string `json:"RuleIds,omitnil,omitempty" name:"RuleIds"`
	// 唯一请求 ID ...
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}
type CreateEdgeAclRuleResponse struct {
	*tchttp.BaseResponse
	Response *CreateEdgeAclRuleResponseParams `json:"Response"`
}
```

即出参路径为 `response.Response.RuleIds`，类型为 `[]*string`。

**资源现状说明**

当前仓库中尚不存在 `tencentcloud/services/fwm/` 目录与 `resource_tc_fwm_edge_acl_rule.go` 文件（vendor 中的 fwm SDK 已通过 go.mod / tools.go 引入，但 provider 尚未注册该资源）。因此本次"新增参数"实际伴随该资源文件的首次落地：实施阶段需创建 `tencentcloud/services/fwm/resource_tc_fwm_edge_acl_rule.go`，按 RESOURCE_KIND_GENERAL 完整 CRUD 实现，并在其中包含 `RuleIds` 出参字段；同时在 `tencentcloud/provider.go` / `tencentcloud/provider.md` 注册该资源。

**关联 CRUD 接口**（同一 fwm SDK，用于实施阶段参考）：
- 创建：`CreateEdgeAclRule`（client.go:193）→ 返回 `RuleIds []*string`
- 查询：`DescribeEdgeAclRules`（client.go:1038 附近）→ `DescribeEdgeAclRulesResponseParams.Rules []*EdgeAclRuleResp`，每条含 `RuleId`
- 删除：`DeleteEdgeAclRule`（client.go:727）→ 入参 `GroupId` + `RuleIds []*string`
- 更新：fwm SDK 无独立 `ModifyEdgeAclRule` 接口，更新通过对规则组内规则先删后建实现（实施阶段以实际 SDK 接口为准）

## Goals / Non-Goals

**Goals:**
- 在 `tencentcloud_fwm_edge_acl_rule` 资源 schema 中新增 `RuleIds` 出参字段（`schema.TypeList` of `schema.TypeString`，`Computed: true`）。
- 在资源 Create 方法中，调用 `CreateEdgeAclRule` 成功后将 `response.Response.RuleIds` 回填到 state。
- 在资源 Read 方法中保持 `RuleIds` 与云端一致（基于 `DescribeEdgeAclRules` 返回的规则列表回填）。

**Non-Goals:**
- 不修改 `CreateEdgeAclRule` 的任何入参。
- 不调整既有字段的 schema 行为（不引入 ForceNew / 不可变约束变更）。
- 不处理 NAT / VPC 等其他类型 ACL 规则资源，仅限 edge ACL。

## Decisions

1. **`RuleIds` 作为 Computed 出参，而非用户可配置入参**
   - 依据：需求明确该字段为 `CreateEdgeAclRule` 的 **出参**（`response.RuleIds`），云 API 创建接口入参中并无同名 `RuleIds`（入参为 `GroupId` + `Rules []*EdgeAclRuleInfo`）。
   - 因此 schema 中 `RuleIds` 设置为 `Computed: true`、不设 `Optional`/`Required`、不设 `ForceNew`，用户不能配置，仅由 Create/Read 回填。
   - 替代方案（设为 Optional 入参）被否定，因为入参侧无对应字段，强行写入会被 SDK 校验拒绝。

2. **类型映射：`[]*string` → `schema.TypeList, Elem: schema.TypeString`**
   - 云 API `RuleIds` 类型为 `[]*string`，Terraform Plugin SDK v2 中使用 `schema.TypeList` + `&schema.Schema{Type: schema.TypeString}` 表达字符串列表，与 provider 内既有写法一致。

3. **ID 策略**
   - 资源 id 采用 `GroupId` + tccommon `FILED_SP` 分隔符 + 首条 `RuleId` 组合（或按规则组维度），具体由实施阶段在 CRUD 中统一；`RuleIds` 仅作为出参落 state，不参与 id 拼接。

4. **Create 失败保护**
   - 按 tencentcloud 代码规范：调用 `CreateEdgeAclRule` 后检查 `response`/`response.Response`/`RuleIds` 是否为空；为空则返回 `NonRetryableError`，避免写入空 id。设置 `d.SetId(...)` 后再回填 `RuleIds`。

## Risks / Trade-offs

- [Create 接口返回多条规则 ID] `CreateEdgeAclRule` 入参 `Rules` 是列表，返回 `RuleIds` 也是列表；资源以单条规则维度建模时，`RuleIds` 可能含多个元素 → 在 Read 时以 `DescribeEdgeAclRules` 实际返回为准回填，Create 时直接落库返回值，不截断。
- [向后兼容] 本次仅新增 Computed 字段 → 现有 state 升级后字段为空直到下一次 Read 刷新，无破坏性，已通过 Computed 属性保证。