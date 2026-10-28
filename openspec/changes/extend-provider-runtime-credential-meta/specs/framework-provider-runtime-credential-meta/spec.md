## ADDED Requirements

### Requirement: 在 `tencentcloud_provider_runtime` schema 中新增 4 个 Computed 字段

framework data source `tencentcloud_provider_runtime` MUST 在已有 `id` / `region` / `client_version` / `stack_mode` 4 个字段之外**追加**以下 4 个 Computed 字段：

| 字段 | 类型 | 来源 | 含义 |
| --- | --- | --- | --- |
| `protocol` | string | `*connectivity.TencentCloudClient.Protocol` | API 请求协议（`HTTP` / `HTTPS`） |
| `domain` | string | `client.Domain` | API 根域名 |
| `cos_domain` | string | `client.CosDomain` | COS 根域名 |
| `secret_id_present` | bool | `client.Credential != nil && client.Credential.SecretId != ""` | 仅指示 SecretId 字段是否非空 |

新增字段 MUST 全部为 Computed，**不得**为 Optional / Required。

新增字段的写入顺序 MUST 排在已有字段之后（schema 声明顺序、Read 写 state 的字段顺序、文档展示顺序均一致）。

#### Scenario: 通过 schema 输出验证 8 个字段全部存在

- **WHEN** 通过 `terraform providers schema -json` 获取 provider schema
- **THEN** `data_source_schemas["tencentcloud_provider_runtime"].block.attributes` 包含 8 个属性：`id` / `region` / `client_version` / `stack_mode` / `protocol` / `domain` / `cos_domain` / `secret_id_present`
- **AND** 全部为 Computed
- **AND** `secret_id_present` 类型为 `bool`，其余为 `string`

#### Scenario: HCL 中读取新字段

- **WHEN** 用户在 HCL 中声明 `data "tencentcloud_provider_runtime" "this" {}` 并 `output { value = data.tencentcloud_provider_runtime.this.secret_id_present }`
- **THEN** `terraform plan` 不发起任何 HTTPS 请求到 `*.tencentcloudapi.com` 或 `*.myqcloud.com`
- **AND** output 值在凭证已配置时为 `true`

### Requirement: `secret_id_present` 语义严格限定为"字段非空"

`secret_id_present` 字段 MUST 仅返回 `client.Credential.SecretId != ""` 的布尔结果，**不得**反映凭证是否能成功调用 API。

文档（website 与 .md 样例）MUST 以粗体或 NOTE 强调：该字段不代表凭证有效性。

#### Scenario: 凭证存在时返回 true

- **WHEN** SDKv2 provider Configure 阶段成功解析了 `TENCENTCLOUD_SECRET_ID` 环境变量或其他凭证来源，使 `client.Credential.SecretId` 非空
- **THEN** Read 写入的 state 中 `secret_id_present == true`

#### Scenario: 文档明确语义

- **WHEN** 检视 `website/docs/d/provider_runtime.html.markdown`
- **THEN** `secret_id_present` 字段说明部分包含粗体或 NOTE 形式的文字，明确"only indicates whether SecretId is configured, NOT whether the credential is valid"

### Requirement: 不暴露任何敏感字段

`tencentcloud_provider_runtime` 的 schema attributes MUST NOT 出现以下字段（无论以何种命名变体）：

- `secret_id` / `secretid` / `access_key`
- `secret_key` / `secretkey`
- `security_token` / `token` / `session_token`
- 上述任何字段的子串、前缀片段或掩码版本

Read 函数 MUST NOT 把上述任何凭证字符串写入 state。

#### Scenario: schema 不包含敏感字段名

- **WHEN** 通过 `terraform providers schema -json` 获取 schema
- **THEN** `tencentcloud_provider_runtime` 属性集合中不存在上述任一字段名

#### Scenario: state 中无凭证泄漏

- **WHEN** 在测试中读取该 data source 的 state map
- **THEN** map 中所有 string value 都不包含真实 SecretId / SecretKey / Token 子串

### Requirement: 已有 4 个字段保持向后兼容

本变更 MUST 保证已有的 `id` / `region` / `client_version` / `stack_mode` 4 个字段在升级前后的 schema、行为、typeName 完全一致。具体约束：

- 类型 MUST 保持不变（全部 string）
- Computed 标记 MUST 保持不变
- Read 行为 MUST 保持不变（`id` 仍等于 region；`stack_mode` 仍返回 `"sdkv2+framework"`）
- typeName MUST 保持 `tencentcloud_provider_runtime` 不变

#### Scenario: 已有字段在升级后值不变

- **WHEN** 同一份 HCL（仅引用 `region` / `client_version` / `stack_mode`）在升级前后分别 `terraform plan`
- **THEN** 这 3 个字段的值与上一版本一致
- **AND** plan 不出现关于这 3 个字段的 force replace

### Requirement: Acceptance test 覆盖新字段

`TestAccTencentCloudProviderRuntimeDataSource_basic` MUST 在原有断言基础上**新增**对 4 个新字段的断言：

- `protocol` 非空
- `domain` 非空
- `cos_domain` 非空
- `secret_id_present == "true"`（在 PreCheck 通过的环境下凭证必然存在）

测试 MUST 不调用任何云 API。

#### Scenario: TF_ACC=1 且凭证存在时所有断言通过

- **WHEN** 设置 `TF_ACC=1` 与 `TENCENTCLOUD_SECRET_ID/KEY`，运行 `go test ./tencentcloud/services/tcprovider/... -run TestAccTencentCloudProviderRuntimeDataSource_basic -v`
- **THEN** 测试通过
- **AND** 测试日志中不出现对 `*.tencentcloudapi.com` 的任何 HTTPS 请求

### Requirement: 文档同步覆盖新字段

provider MUST 提供 / 更新 `website/docs/d/provider_runtime.html.markdown`，使其覆盖**全部** 8 个字段（4 个原字段 + 4 个新字段）。

文档 MUST 包含一个最小可运行的 ```hcl 示例，至少引用一个新字段。

`secret_id_present` 字段说明 MUST 包含语义警告（详见上一条 Requirement）。

#### Scenario: website 文档字段完整

- **WHEN** 检视 `website/docs/d/provider_runtime.html.markdown`
- **THEN** 文件存在
- **AND** 包含 frontmatter（`subcategory` / `layout` / `page_title` / `description` 四项）
- **AND** 字段说明部分覆盖 `id` / `region` / `client_version` / `stack_mode` / `protocol` / `domain` / `cos_domain` / `secret_id_present` 共 8 项
- **AND** 至少包含一个 ```hcl 代码块示例
