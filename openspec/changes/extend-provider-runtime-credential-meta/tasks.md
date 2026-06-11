> 备注（与项目规则的差异）：
> 项目硬约束要求 `website/docs/` 文档由 `make doc` 生成，禁止手写。但本变更涉及的是 framework data source，现有 `gendoc` 工具仅识别 SDKv2 的 `DataSourcesMap`，**无法**为 framework data source 生成文档。本变更**临时**手写 / 编辑 website 文档；该例外的根本修复属于 stage D（gendoc 升级 / 切 tfplugindocs）。

## 1. 准备 / 前置确认

- [x] 1.1 `read_file` 当前 `tencentcloud/services/tcprovider/data_source_tc_provider_runtime_framework.go`，确认 `providerRuntimeModel` struct 与 `Schema` / `Read` 函数的真实最新行号与字段顺序
- [x] 1.2 `read_file` 当前 `tencentcloud/services/tcprovider/data_source_tc_provider_runtime_framework_test.go`，记录已有断言风格，新断言保持一致
- [x] 1.3 检查是否已有 `website/docs/d/provider_runtime.html.markdown`：文件已存在（1.63 KB），走 **1.3a** 在已有内容上**追加**字段说明

## 2. data source 实现扩展

- [x] 2.1 在 `providerRuntimeModel` struct 末尾追加 4 个字段：
  - `Protocol         types.String \`tfsdk:"protocol"\``
  - `Domain           types.String \`tfsdk:"domain"\``
  - `CosDomain        types.String \`tfsdk:"cos_domain"\``
  - `SecretIDPresent  types.Bool   \`tfsdk:"secret_id_present"\``

- [x] 2.2 在 `Schema` 函数 `Attributes` map 末尾追加 4 项：
  - [x] 2.2.1 `protocol`：`schema.StringAttribute{ Computed: true, Description: "..." }`
  - [x] 2.2.2 `domain`：同上
  - [x] 2.2.3 `cos_domain`：同上
  - [x] 2.2.4 `secret_id_present`：`schema.BoolAttribute{ Computed: true, Description: "..." }`，Description 中明确标注“仅指示是否配置，不代表凭证有效”

- [x] 2.3 在 `Read` 函数中：
  - [x] 2.3.1 从 `d.client` 读出 `Protocol` / `Domain` / `CosDomain`
  - [x] 2.3.2 计算 `secretIDPresent := d.client != nil && d.client.Credential != nil && d.client.Credential.SecretId != ""`
  - [x] 2.3.3 把 4 个字段写入 `state` 结构体（保持与已有字段相同的 `tcfwhelper.StringValueOrNull` / `types.BoolValue` 风格）
  - [x] 2.3.4 防御式：`d.client == nil` 时新字段全部走 null / false，不 panic

- [x] 2.4 不修改：`Metadata` / `Configure` / typeName / 已有 4 个字段的语义

## 3. 测试扩展

- [x] 3.1 在 `TestAccTencentCloudProviderRuntimeDataSource_basic` 的 `resource.TestCheckFunc` 列表末尾追加：
  - [x] 3.1.1 `resource.TestCheckResourceAttrSet(..., "protocol")`
  - [x] 3.1.2 `resource.TestCheckResourceAttrSet(..., "domain")`
  - [x] 3.1.3 `resource.TestCheckResourceAttrSet(..., "cos_domain")`
  - [x] 3.1.4 `resource.TestCheckResourceAttr(..., "secret_id_present", "true")`

- [x] 3.2 不修改：HCL config / PreCheck / Steps 数量

## 4. 文档扩展

- [x] 4.1 编辑 / 新建 `website/docs/d/provider_runtime.html.markdown`：
  - [x] 4.1.1 frontmatter 完整（`subcategory: "Provider Data Sources"` 或前序已用的子分类 / `layout` / `page_title` / `description`）
  - [x] 4.1.2 字段说明部分覆盖 8 个字段（4 已有 + 4 新增）
  - [x] 4.1.3 `secret_id_present` 说明使用 `**Note:**` 或 markdown 粗体标明“only indicates whether SecretId is configured, NOT whether the credential is valid”
  - [x] 4.1.4 至少 1 个 ```hcl 最小示例，输出至少 1 个新字段

## 5. 校验（代码修改后单独执行）

- [x] 5.1 `go build ./...` 通过
- [x] 5.2 `go vet ./tencentcloud/services/tcprovider/...` 通过
- [x] 5.3 `gofmt -l tencentcloud/services/tcprovider/` 输出为空
- [ ] 5.4 在设置 `TF_ACC=1` + `TENCENTCLOUD_SECRET_ID/KEY` 的环境下运行：`go test ./tencentcloud/services/tcprovider/... -run TestAccTencentCloudProviderRuntimeDataSource_basic -v`，测试通过  *(待用户在带凭证环境运行)*

## 6. mux 合并冲突检查

- [ ] 6.1 启动 provider 二进制，运行 `terraform providers schema -json`  *(需本地环境手动验证)*
- [ ] 6.2 确认 `data_source_schemas["tencentcloud_provider_runtime"].block.attributes` 包含 8 个字段（不多不少）
- [ ] 6.3 确认 SDKv2 已有的所有 data source schema 未被影响（无新增 / 删除 / 类型变化）

## 7. 收尾

- [x] 7.1 在 `.changelog/` 新建占位 `0000-extend-provider-runtime.txt`：`enhancement: data-source/tencentcloud_provider_runtime: Add protocol, domain, cos_domain, secret_id_present`（待 PR 号确定后由维护者重命名）
- [x] 7.2 自查 spec 6 条 Requirement 是否每条 Scenario 都有代码或测试覆盖
- [x] 7.3 `openspec status --change "extend-provider-runtime-credential-meta"` 输出 `isComplete: true`，准备 archive