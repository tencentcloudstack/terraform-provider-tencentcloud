## Why

前序 change `add-plugin-framework-muxing` 已经落地了 framework 侧的第一个 data source `tencentcloud_provider_runtime`，覆盖 `id` / `region` / `client_version` / `stack_mode` 四个字段，端到端验证了 SDKv2 → tcfwprovider → framework Configure → data source Configure 的链路。

但仍有用户视角的运行时元信息**缺失**，导致排障体验不完整：

1. **协议**（`HTTP` vs `HTTPS`）：私有云 / 国际站 / 自定义代理场景下用户经常配错
2. **API 根域名**（`domain` / `cos_domain`）：私有云通过自定义 endpoint 接入时，无法在 plan 阶段确认 provider 实际使用的 endpoint
3. **凭证存在性**：用户经常在 CI 流水线里因 `TENCENTCLOUD_SECRET_ID` 注入失败而困惑，目前没有非敏感的方式确认凭证是否注入到 provider

继续新增一个 `tencentcloud_provider_meta` 会与 `provider_runtime` 高度重叠（都暴露 `region`，都用于自检），徒增维护成本。**正确做法是扩展现有 data source**。

## What Changes

- 在已有 `tencentcloud_provider_runtime` data source 的 schema 中**新增 4 个 Computed 属性**：
  - `protocol` (string)：API 协议，来自 `client.Protocol`
  - `domain` (string)：API 根域名，来自 `client.Domain`
  - `cos_domain` (string)：COS 根域名，来自 `client.CosDomain`
  - `secret_id_present` (bool)：仅指示 `client.Credential.SecretId != ""`，**永远不暴露**真实 SecretId / SecretKey / Token
- 同步更新 Read 函数：把上述字段从 `*connectivity.TencentCloudClient` 读出来写入 state
- 同步更新 acceptance test：在 `TestAccTencentCloudProviderRuntimeDataSource_basic` 中新增字段断言
- 同步更新 / 新建 website 文档：`website/docs/d/provider_runtime.html.markdown`（如不存在则新建）

不引入任何 BREAKING：

- 现有 `id` / `region` / `client_version` / `stack_mode` 4 个字段语义、类型、Computed 属性不变
- 仅追加 Computed 字段，符合 terraform schema 向后兼容规则
- typeName `tencentcloud_provider_runtime` 不变
- 无任何 SDKv2 资源 / 数据源被影响

## Capabilities

### New Capabilities

- `framework-provider-runtime-credential-meta`：在 framework data source `tencentcloud_provider_runtime` 中新增凭证存在性 / 协议 / 域名等运行时元信息字段（前序 change `add-plugin-framework-muxing` 尚未 archive，因此这里以新 capability 的形式声明本次新增字段的 Requirement，待两个 change 都 archive 后由项目维护者合并整理）

### Modified Capabilities

无（前序 spec 还未沉淀到 `openspec/specs/`，无法做 MODIFIED Requirements 引用）。

## Impact

- **Affected code**
  - 修改：`tencentcloud/services/tcprovider/data_source_tc_provider_runtime_framework.go`
    - `providerRuntimeModel` struct 新增 4 个字段
    - `Schema` 函数 attributes map 新增 4 项
    - `Read` 函数填充 4 个新字段
  - 修改：`tencentcloud/services/tcprovider/data_source_tc_provider_runtime_framework_test.go`
    - 新增 4 个字段的断言
  - 新建（或修改）：`website/docs/d/provider_runtime.html.markdown`
    - 完整文档：4 个原字段 + 4 个新字段，明确标注 `secret_id_present` 不代表凭证有效性

- **Affected APIs**
  - 用户已经在用 `data.tencentcloud_provider_runtime.this.region` 等表达式不受影响
  - 用户可以新写 `data.tencentcloud_provider_runtime.this.secret_id_present` 等

- **Affected dependencies**
  - 无新依赖

- **Affected systems**
  - mux server schema 合并：新增字段不会引发冲突（`tencentcloud_provider_runtime` 在 SDKv2 端不存在）

- **Affected docs**
  - 1 份 website 文档新增 / 修改

- **风险与回滚**
  - 风险：极低。仅追加 Computed 字段
  - 回滚：把 .go / _test.go / website 文档中新增的内容删除即可，无 schema / state 影响
