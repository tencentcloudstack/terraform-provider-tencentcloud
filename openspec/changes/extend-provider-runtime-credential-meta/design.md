## Context

前序 change `add-plugin-framework-muxing`（status: in-progress, 37/40）已经把 framework provider 接入 mux，并在 `tencentcloud/services/tcprovider/` 下落地了**第一个** framework data source `tencentcloud_provider_runtime`，schema 字段为：

| 字段 | 类型 | 来源 | 用途 |
| --- | --- | --- | --- |
| `id` | string | region 同值（synthetic） | 兼容老工具 |
| `region` | string | `client.Region` | 自检当前 region |
| `client_version` | string | `connectivity.GetReqClientVersion()` | 调试用 X-TC-RequestClient |
| `stack_mode` | string | 常量 `"sdkv2+framework"` | 标识双栈架构 |

实现位于：
- `tencentcloud/services/tcprovider/data_source_tc_provider_runtime_framework.go`
- `tencentcloud/services/tcprovider/data_source_tc_provider_runtime_framework_test.go`
- `tencentcloud/services/tcprovider/framework.go`（在 `FrameworkDataSources()` 中通过 `NewProviderRuntimeDataSource` 注册）

`*connectivity.TencentCloudClient` 上还有公开字段尚未被该 data source 暴露：`Protocol` / `Domain` / `CosDomain` / `Credential.SecretId`。这些字段在私有云 / 国际站 / CI 注入凭证等场景下都是排障刚需，需要补齐。

约束：
- 不能引入新的外部 API 调用（保持纯本地、零网络）
- 不能暴露任何敏感字段（SecretId / SecretKey / Token 不能进 schema 输出）
- 不能改变现有 4 个字段的类型 / Computed / 含义（向后兼容）
- 不能改变 typeName `tencentcloud_provider_runtime`

利益相关方：provider 维护者（需要 schema 字段稳定）、运维 / CI（需要 plan 阶段自检凭证 / endpoint）、用户（HCL 中可读 provider 元信息）。

## Goals / Non-Goals

**Goals:**
- 在 `tencentcloud_provider_runtime` 上**新增** 4 个 Computed 字段：`protocol` / `domain` / `cos_domain` / `secret_id_present`
- 字段值全部从已存在的 `*connectivity.TencentCloudClient` 读取，**不发起任何外部 API 调用**
- `secret_id_present` 用 bool 而不是直接暴露 SecretId，零敏感泄漏风险
- 同步补齐 acceptance test 与 website 文档
- 字段命名与 SDKv2 provider schema 保持一致（`protocol` / `domain` / `cos_domain` 都是 SDKv2 已有的 provider 配置项名）

**Non-Goals:**
- ❌ 不暴露 `app_id` / `owner_uin`：需要异步调用 GetUserAppId / DescribeOwnerUin，违背"零外部 API 调用"目标
- ❌ 不新增 `tencentcloud_provider_meta` 等并列 data source（避免与 `provider_runtime` 职责重叠）
- ❌ 不修改 `provider_runtime` 已有 4 个字段（向后兼容）
- ❌ 不改造 gendoc 让其识别 framework data source（属于 stage D 范围）
- ❌ 不引入 bool 之外的方式表达凭证存在性（如不暴露 SecretId 掩码）

## Decisions

### D1：扩展 `tencentcloud_provider_runtime`，**不**新建 data source

**为什么：**
- 前序 change 已落地的 `provider_runtime` 已经是"运行时元信息"语义，新增字段是合理补充，不引入第二个语义边界
- 用户视角下"我的 provider 现在跑在哪个 region / 用什么协议 / 凭证有没有注入"是同一类问题，应在同一个 data source 里回答

**备选与否决：**
- ~~新建 `tencentcloud_provider_meta`~~：与 `provider_runtime` 高度重叠（都有 `region`），徒增维护成本与用户认知负担
- ~~新建 `tencentcloud_provider_credentials`~~：单独把凭证元信息拆开 → 同样的语义碎片化，且名字易让人误以为暴露真实凭证

### D2：新增字段集 = `protocol` / `domain` / `cos_domain` / `secret_id_present`

**为什么：**
- 这 4 个是 `*connectivity.TencentCloudClient` 上**剩余的、可安全暴露的、非敏感的**公开字段
- 字段命名与 SDKv2 provider schema 同名（`protocol` / `domain` / `cos_domain` 都是已有的 provider 配置项），用户认知一致
- `secret_id_present` 是 bool 而非掩码字符串，**完全无法**反推真实凭证

**备选与否决：**
- ~~暴露 `secret_id` 前 4 位+后 4 位掩码~~：即便掩码，也存在用户误认为这是"真实凭证片段"导致误传播的风险；bool 已经足以表达"是否配置"
- ~~暴露 `region_endpoint`~~：region 已在；endpoint 可由 region+domain 推导，不必单独存

### D3：`secret_id_present` 的语义严格限定为"字段非空"

**为什么：**
- 如果让该字段反映"凭证有效性"，则需要调用 GetUserAppId 等 API → 违背 Non-Goals
- "字段非空" 的判定可以在 Read 中本地完成，零网络

**文档约束：**
- website 文档与 .md 样例都 MUST 用粗体明确："**This field only indicates whether SecretId is configured, NOT whether the credential is valid.**"

### D4：字段与现有字段的写入顺序保持稳定

**为什么：**
- terraform 在 plan 输出 / state diff 中按字段声明顺序展示
- 把新字段插到 schema attributes map 末尾，避免现有用户在 plan 输出中看到不必要的 diff 噪音
- 测试断言同样追加到末尾

### D5：测试只新增字段断言，不调用云 API

**为什么：**
- 前序 change 已经验证了基础链路；本 change 只增量验证 4 个新字段
- 新增字段断言：`protocol`、`domain`、`cos_domain` 非空；`secret_id_present == "true"`
- 沿用前序 `acctest.AccPreCheck(t)` 前置，未设置环境变量则 skip

## Risks / Trade-offs

- **[R1] 用户误把 `secret_id_present == true` 理解为"凭证可用"**
  - → Mitigation：website 文档与 .md 样例都用粗体警告

- **[R2] 现有用户的 state 在升级后产生 diff 噪音**
  - → Trade-off：data source 的 state 在每次 plan 都重新 Read，新增 Computed 字段不会触发 force replace；用户首次升级后会看到一次性的 state schema 变更，但不会触发实际操作

- **[R3] 字段命名与 SDKv2 provider schema 同名导致歧义**
  - → Trade-off：实际语义完全一致（同一个 client 字段），同名反而降低用户学习成本；文档统一描述，明确这是 provider 当前实际使用的值（可能与 HCL 配置不同——例如来自环境变量）

- **[R4] 后续若要新增更多元信息字段（如 endpoint / appId）会再次膨胀该 data source**
  - → Trade-off：可接受。`provider_runtime` 是"运行时元信息总入口"的设计本就预期会增长；阈值控制在"不调用任何外部 API 且非敏感"

## Migration Plan

部署：
1. PR 合入后随 next patch 版本（如 1.81.x）发布
2. 现有用户 `data "tencentcloud_provider_runtime" "this" {}` 的引用零修改
3. 用户可选择性在 HCL 里读新字段

回滚：
- 把 `.go` / `_test.go` / website 文档中新增的内容删除即可，schema 字段是 Computed only-add，不会破坏任何 state

## Open Questions

- 暂无。所有方案已闭环。
