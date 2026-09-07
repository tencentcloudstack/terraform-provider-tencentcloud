## Context

`tencentcloud_teo_security_policy_config` 资源封装了 TEO Web 防护的安全策略（`SecurityPolicy`），其中 `exception_rules`（对应 `ExceptionRules`）用于配置例外规则，使指定请求跳过安全防护模块的扫描。该资源已有同级字段 `web_security_modules_for_exception`（对应 `WebSecurityModulesForException`），用于在 `SkipScope = WebSecurityModules` 时指定**模块级**例外目标（托管规则、速率限制、自定义规则、自适应频控、Bot 管理等）。

TEO 云 API 在 `teov20220901.ExceptionRule` 结构体中新增了 `WebSecuritySubmodulesForException []*string` 字段，用于在 `SkipScope = WebSecuritySubmodules` 时指定**子模块级**例外目标，提供更细粒度的控制（例如托管规则模块下的"规则集"与"高频扫描防护"、HTTP DDoS 防护下的"自适应频控"/"智能客户端过滤"/"流量盗刷防护"、高级 Bot 管理下的多个子功能、基础 Bot 管理下的"AI 爬虫处置"/"人机校验页"等）。

当前状态：
- `resource_tc_teo_security_policy_config.go` 的 schema 在 `security_policy.exception_rules.rules` 下已定义 `web_security_modules_for_exception`（TypeSet）、`managed_rules_for_exception`、`managed_rule_groups_for_exception`、`request_fields_for_exception` 等字段，但缺少 `web_security_submodules_for_exception`。
- Read 逻辑（约 3845-3920 行）遍历 `respData.ExceptionRules.Rules` 时，对每个 rule 读取 `WebSecurityModulesForException` 但未读取 `WebSecuritySubmodulesForException`。
- Create/Update 逻辑（约 5203-5292 行）遍历 `exception_rules.rules` 配置时，处理 `web_security_modules_for_exception`（从 `*schema.Set` 转 `[]*string`）但未处理 `web_security_submodules_for_exception`。
- SDK 字段已确认存在：`teov20220901.ExceptionRule.WebSecuritySubmodulesForException` 类型为 `[]*string`，注释说明仅当 `SkipScope` 为 `WebSecuritySubmodules` 时有效。

约束：
- 向后兼容：必须保持现有 TF 配置与 state 不变，新增字段为纯加法的 Optional 字段。
- 复用现有模式：`WebSecuritySubmodulesForException` 与 `WebSecurityModulesForException` 的数据结构完全一致（均为 `[]*string`），应严格参照后者的 schema（TypeSet）、Read（nil 检查后直接赋值）与 Update（`*schema.Set` 遍历转指针）实现模式。
- 不引入新的 vendor 依赖，不修改 CRD 接口的其他行为。

## Goals / Non-Goals

**Goals:**
- 让用户能够通过 Terraform 声明式配置例外规则的子模块级跳过目标（`WebSecuritySubmodulesForException`）。
- 保持与现有 `web_security_modules_for_exception` 完全一致的处理模式，降低维护成本。
- 向后兼容：现有配置与 state 升级后不产生 plan diff。
- 通过单元测试覆盖新参数的 Create/Read/Update 路径。

**Non-Goals:**
- 不修改 `web_security_modules_for_exception` 及其他已有例外规则字段的现有行为。
- 不调整 `skip_scope` 字段的取值约束（`WebSecuritySubmodules` 作为合法取值由云 API 决定，provider 不在 schema 层做枚举校验）。
- 不新增独立的 `tencentcloud_teo_security_policy_config_submodule` 资源。
- 不改变资源 Timeouts 或 Provider 注册逻辑。

## Decisions

### Decision 1: 新增字段为 TypeSet（与 `web_security_modules_for_exception` 一致）

**选择**：在 `security_policy.exception_rules.rules` 块下新增 `web_security_submodules_for_exception`（`schema.TypeSet`, `Optional: true`, `Elem: &schema.Schema{Type: schema.TypeString}`），紧邻 `web_security_modules_for_exception`。

**备选**：使用 `schema.TypeList`。

**理由**：
- 与同级、同语义的 `web_security_modules_for_exception` 保持一致，后者为 TypeSet，且子模块集合本身是无序的，TypeSet 的集合语义与去重行为更贴合。
- TypeSet 在 Update 时通过 `.*schema.Set` 取值，复用已有遍历模式，降低实现与审查成本。

### Decision 2: Read ���辑——nil 检查后直接赋值

**选择**：在 Read 遍历 `respData.ExceptionRules.Rules` 的循环中，紧接 `WebSecurityModulesForException` 的处理之后，新增：

```go
if rules.WebSecuritySubmodulesForException != nil {
    rulesMap["web_security_submodules_for_exception"] = rules.WebSecuritySubmodulesForException
}
```

**理由**：
- 与 `WebSecurityModulesForException` 的读取模式完全一致（SDK 返回 `[]*string`，直接写入 map，由 SDK 的 `d.Set` 完成 flatten）。
- 遵循 provider 规范：仅在 Response 字段非 nil 时调用 set，避免 nil 覆盖已有配置。
- 遵循规范：若云 API 返回空，保留现场日志后再处理（本字段为列表，nil 时直接跳过 set 即可）。

### Decision 3: Update/Create 逻辑——TypeSet 遍历转 `[]*string`

**选择**：在 Create/Update 遍历 `exception_rules.rules` 配置的循环中，紧接 `web_security_modules_for_exception` 的处理之后，新增：

```go
if v, ok := rulesMap["web_security_submodules_for_exception"]; ok {
    webSecuritySubmodulesForExceptionSet := v.(*schema.Set).List()
    for i := range webSecuritySubmodulesForExceptionSet {
        if webSecuritySubmodulesForExceptionSet[i] != nil {
            webSecuritySubmodulesForException := webSecuritySubmodulesForExceptionSet[i].(string)
            exceptionRule.WebSecuritySubmodulesForException = append(exceptionRule.WebSecuritySubmodulesForException, &webSecuritySubmodulesForException)
        }
    }
}
```

**理由**：
- 与 `web_security_modules_for_exception` 的扩展模式逐字对齐，保证可读性与可维护性。
- 对集合元素做非 nil 判断后取 `string`，再取地址追加到 `[]*string`，避免空指针。
- Create 与 Update 共用同一构造逻辑（该资源 Update 即调用 Create 后续逻辑或共享函数），无需重复实现。

### Decision 4: 单元测试使用 gomonkey mock 云 API

**选择**：在 `resource_tc_teo_security_policy_config_test.go` 中新增测试用例，使用 gomonkey mock `ModifySecurityPolicyWithContext` 与 `DescribeSecurityPolicyWithContext`（或对应 service 层封装），构造包含 `web_security_submodules_for_exception` 的配置，验证 Create/Read/Update 路径中该字段被正确填充与回读。

**备选**：使用 Terraform 验收测试套件（`TF_ACC=1`）。

**理由**：
- 遵循项目约束：新增参数不使用 terraform 测试套件，而使用 gomonkey mock 云 API 进行业务逻辑单测，避免依赖真实云资源与凭证。
- 覆盖关键断言：请求中 `ExceptionRule.WebSecuritySubmodulesForException` 非空且元素正确；Read 回读后 state 中 `web_security_submodules_for_exception` 与预期一致。

### Decision 5: 文档通过收尾阶段 `make doc` 生成

**选择**：仅在 `resource_tc_teo_security_policy_config.md` 中更新 Example Usage（展示 `web_security_submodules_for_exception`），`website/docs/` 下的文档由收尾阶段 `make doc` 自动生成。

**理由**：
- 遵循项目硬约束：禁止直接新增/修改 `website/` 目录下文件，统一由 `make doc` 生成。

## Risks / Trade-offs

- **Risk**：存量 state 中不存在 `web_security_submodules_for_exception`，升级后首次 plan 可能因 Read 回填新字段而显示一次 from-None 的 diff → **Mitigation**：字段为纯 Optional，Read 仅在云端返回非 nil 时写入；存量资源若云端未配置子模块例外则不会产生 diff，若云端已配置则 plan diff 属于合理的状态收敛。
- **Risk**：用户配置 `skip_scope = WebSecuritySubmodules` 但未填 `web_security_submodules_for_exception`，或反之 → **Mitigation**：provider 不在 schema 层强制关联校验，由云 API 返回错误并透传给用户，与现有 `web_security_modules_for_exception` 一致。
- **Trade-off**：新增字段与现有 `web_security_modules_for_exception` 名字高度相似，用户易混淆 → 通过 schema Description 明确两者的区别（模块 vs 子模块），并给出子模块合法取值示例。

## Migration Plan

- 新增字段为纯加法（Optional），无 state 迁移需求。
- 存量资源：Terraform state 中无该字段，升级后 `terraform plan` 对未在 HCL 配置该字段的资源不会产生 diff（Optional + Read 仅在云端返回非 nil 时回填）。
- 文档更新：在 `resource_tc_teo_security_policy_config.md` 的例外规则示例中补充 `web_security_submodules_for_exception` 字段。
- 回滚：若需回退，移除 schema 中字段及对应 Read/Update 分支即可，state 中已有的值会被忽略，无破坏性影响。

## Open Questions

- 无
