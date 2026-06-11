## Why

`add-plugin-framework-muxing` 已经把 `terraform-plugin-framework` 接入 mux server，并在仓里落地了**第一个** framework data source（`tencentcloud_provider_runtime`）以及一组共享基础设施。但这套基础设施在落地时为了避免与项目原有 `tccommon` / `services/tcXxx` 命名冲突，临时使用了 `tcfwhelper` / `tcfwprovider` / `services/tcprovider` 这些**带 `fw` 字样的过渡命名**，且把 framework registry 的入口、framework provider 的入口、第一个 data source 分散在 `tencentcloud/`（顶层 test）、`tencentcloud/provider/framework/`（provider 入口）、`tencentcloud/services/tcprovider/`（registry + 业务实现）三处。

随着接下来要在 framework 侧继续新增 `resource` / `function` / `ephemeral` / `list` / `action` 五种新类型，**当前命名与目录布局已经成为可读性 / 可维护性的明显瓶颈**：
- 包名 `tcfwhelper` / `tcfwprovider` / 别名 `fwprovider` 包含 `fw` 字样，对新读者不友好；
- framework provider 入口在 `provider/framework/`，业务实现却在 `services/tcprovider/`，分裂在两个目录树，新增类型时心智负担高；
- 没有"产品（service）"维度的分层，未来 framework 资源跨多个云产品（CVM、VPC、CBS……）时会再次撜爆；
- 没有任何 `resource` / `function` / `ephemeral` / `list` / `action` 的 reference 实现，后续开发者无从参考"如何在本仓接入 X 类型"。

关于依赖版本：本 change 启动后实测发现，当前 `terraform-plugin-framework v1.19.0` 上游**已包含完整 `action` 包**（vendor 中 `action.go` / `provider.ProviderWithActions` 接口齐备），且 `terraform-plugin-mux v0.23.1` 已发布 `mux_server_InvokeAction.go` 支持 action 协议。因此本 change **不进行依赖升级**，仅锁定最低版本约束为当前锁定的 `v1.19.0` / `v0.23.1`（详见 design.md Decision 6）。

本次改动一次性把"命名规整 + 目录扁平化合并 + 引入产品分层 + 6 类型示范"做完，为后续在 framework 侧持续扩张产品能力打好地基。
## What Changes

### 命名（全仓重命名）
- **BREAKING**（仅内部包路径，不影响 Terraform 用户配置）：
  - `tencentcloud/internal/tcfwhelper/` → `tencentcloud/framework/internal/helper/`（中间经过一轮过渡名 `tencentcloud/internal/frameworkhelper/`，最终收敛到 framework 子树并受 Go `internal/` 可见性约束，包名由 `frameworkhelper` 改为 `helper`）
  - `tencentcloud/internal/tcfwprovider/` → `tencentcloud/internal/sharedmeta/`（双栈共享桥，**保留在 internal/**，需被 SDKv2 与 framework 两侧 import）
  - `tencentcloud/services/tcprovider/`（业务实现） → 拆解后并入 `tencentcloud/framework/<product>/<type>/`
  - `tencentcloud/provider/framework/`（provider 入口） → **下沉合并** 到 `tencentcloud/framework/`（与下方业务实现同目录，包名 `framework`）
  - `tencentcloud/acctest/framework_factories.go`（framework-only 测试工厂） → `tencentcloud/framework/acctest/factories.go`（包名 `frameworkacctest`；SDKv2 共用的 `AccPreCheck` / `test_util.go` 仍保留在 `tencentcloud/acctest/`）
  - import 别名 `fwprovider` → 直接使用包名 `framework`
- 文件后缀 `_framework.go` 统一去掉（迁入 `tencentcloud/framework/` 后通过目录隔离已经表达"framework 风味"）

### 目录布局（一次性合并 + 产品/类型双层）

最终顶层布局：

```
tencentcloud/
├── services/                          # SDKv2 资源/数据源（既有，保持不变）
│   ├── cvm/
│   ├── vpc/
│   └── ...
├── framework/                         # plugin-framework 一切：入口 + 业务实现 + framework-only helper/acctest
│   ├── provider.go                    # framework Provider Schema/Configure（包名 framework）
│   ├── registry.go                    # 6 类型聚合器，直接 import 各产品子包工厂
│   ├── provider_test.go               # 同包测试（从 tencentcloud/framework_provider_test.go 迁入）
│   ├── testhelpers_test.go            # 同包测试 helper（从 tencentcloud/framework_provider_testhelpers_test.go 迁入）
│   ├── README.md
│   ├── internal/                      # framework 子树私有（受 Go internal/ 可见性约束）
│   │   └── helper/                    # 包名 helper：types/retry/error/timeouts 工具
│   ├── acctest/                       # 包名 frameworkacctest：framework 资源 acceptance test 用的 ProtoV5 工厂
│   │   └── factories.go
│   ├── cvm/                           # 产品：CVM
│   │   └── actions/
│   │       ├── reboot_instance_action.go
│   │       └── reboot_instance_action_test.go
│   └── meta/                          # 元产品：不归属任何具体云产品的 reference
│       ├── resources/
│       │   ├── local_note_resource.go
│       │   └── local_note_resource_test.go
│       ├── datasources/
│       │   ├── provider_runtime_data_source.go        # 从 services/tcprovider 迁移
│       │   └── provider_runtime_data_source_test.go
│       ├── functions/
│       │   ├── parse_resource_id_function.go
│       │   └── parse_resource_id_function_test.go
│       ├── ephemerals/
│       │   ├── temp_credential_ephemeral_resource.go
│       │   └── temp_credential_ephemeral_resource_test.go
│       └── lists/
│           ├── region_list_resource.go
│           └── region_list_resource_test.go
├── provider/
│   └── sdkv2/                         # SDKv2 provider 入口（占位/既有）
├── acctest/                           # SDKv2 + 双栈共用 acceptance test helper（AccPreCheck 等）
└── internal/
    └── sharedmeta/                    # 双栈共享 *connectivity.TencentCloudClient 的桥（必须留在 internal/，被 SDKv2 与 framework 双向 import）
```

### 包命名约定

- `tencentcloud/framework/` 的 Go 包名为 **`framework`**（目录与包名一致）。
  - 调用方写法：`framework.NewProvider(primary)`（替代旧 `fwprovider.NewProvider(...)` 与 `tcprovider.FrameworkResources(...)` 两处）
- 产品目录（`cvm/` / `meta/`）本身不放 `.go` 文件，仅作为命名空间容器；不需要包声明
- 类型子目录的包名直接采用目录名，并以**产品前缀消歧**（避免不同产品下同名类型子包撞名）：
  - `framework/cvm/actions/` → `package cvmactions`
  - `framework/meta/resources/` → `package metaresources`
  - `framework/meta/datasources/` → `package metadatasources`
  - `framework/meta/functions/` → `package metafunctions`
  - `framework/meta/ephemerals/` → `package metaephemerals`
  - `framework/meta/lists/` → `package metalists`
- 这样 `registry.go` 中可以同时 import 多个产品的同类子包而不需要 alias

### 产品归属约定

本次 6 个 reference 的产品归属：

| Reference | 类型 | 归属产品 | 理由 |
|---|---|---|---|
| `tencentcloud_reboot_instance` | action | `cvm/` | 语义上属于 CVM |
| `tencentcloud_local_note` | resource | `meta/` | 纯本地 in-memory，不归属云产品 |
| `tencentcloud_provider_runtime` | datasource | `meta/` | 暴露的是 provider 自身运行时元数据 |
| `parse_resource_id` | function | `meta/` | provider-defined function，无产品归属 |
| `tencentcloud_temp_credential` | ephemeral | `meta/` | 跨产品共享凭证 |
| `tencentcloud_region` | list | `meta/` | provider 级区域列表 |

约定：**未来新增 framework 资源若能明确归属真实云产品，必须落到对应产品目录**（如 `framework/vpc/resources/...`）；**只有跨产品 / 不归属任何产品的资源**才允许进 `framework/meta/`。

### `tencentcloud/framework/` 内部 registry 形式

`registry.go` 直接 import 各产品包的工厂方法做 slice 聚合（**取消**原 `services/tcprovider/framework.go` 中间层）：

```go
package framework

import (
    cvmactions     "github.com/.../tencentcloud/framework/cvm/actions"
    metadatasources "github.com/.../tencentcloud/framework/meta/datasources"
    metaresources   "github.com/.../tencentcloud/framework/meta/resources"
    metafunctions   "github.com/.../tencentcloud/framework/meta/functions"
    metaephemerals  "github.com/.../tencentcloud/framework/meta/ephemerals"
    metalists       "github.com/.../tencentcloud/framework/meta/lists"
)

func (p *Provider) Resources(_ context.Context) []func() resource.Resource {
    return []func() resource.Resource{ metaresources.NewLocalNoteResource }
}
func (p *Provider) DataSources(_ context.Context) []func() datasource.DataSource {
    return []func() datasource.DataSource{ metadatasources.NewProviderRuntimeDataSource }
}
func (p *Provider) Functions(_ context.Context) []func() function.Function {
    return []func() function.Function{ metafunctions.NewParseResourceIDFunction }
}
func (p *Provider) EphemeralResources(_ context.Context) []func() ephemeral.EphemeralResource {
    return []func() ephemeral.EphemeralResource{ metaephemerals.NewTempCredentialEphemeralResource }
}
func (p *Provider) ListResources(_ context.Context) []func() list.ListResource {
    return []func() list.ListResource{ metalists.NewRegionListResource }
}
func (p *Provider) Actions(_ context.Context) []func() action.Action {
    return []func() action.Action{ cvmactions.NewRebootInstanceAction }
}
```

### 6 种类型各落地一份 reference 实现（按业务命名，不带 `example_` 前缀）
- `framework/meta/datasources/provider_runtime_data_source.go`（**搬迁**已有实现）
- `framework/meta/resources/local_note_resource.go`（轻量 stub：纯本地 in-memory 资源）
- `framework/meta/functions/parse_resource_id_function.go`（演示 plugin-framework Function：拆 `instanceId#userId` → list）
- `framework/meta/ephemerals/temp_credential_ephemeral_resource.go`（演示 EphemeralResource：本地构造短期凭证）
- `framework/meta/lists/region_list_resource.go`（演示 ListResource：本地静态区域列表）
- `framework/cvm/actions/reboot_instance_action.go`（演示 Action：仅 stub 校验，不调 CVM API）

> 6 个 reference 都是 L2 级别：**接入 registry**，可以被 `terraform providers schema -json` 看到；但所有"业务副作用"或被替换为本地实现（local note / region list 静态返回），或被替换为只校验/打日志的 stub（reboot action），保证不会污染真实用户数据，也无需联网即可 unit test。

### 测试搬迁
- `tencentcloud/framework_provider_test.go` → `tencentcloud/framework/provider_test.go`
- `tencentcloud/framework_provider_testhelpers_test.go` → `tencentcloud/framework/testhelpers_test.go`
- 测试中的 `fwprovider` 别名替换为同包直接调用（`NewProvider(...)`）

### Registry 支持 Action
- `tencentcloud/framework/registry.go` 新增 `Actions()` 聚合方法
- `tencentcloud/framework/provider.go` 新增 `Actions(ctx)` provider 接口实现 + Configure 中 `resp.ActionData = meta`
- 不再有独立 `framework.go` 聚合点（取消中间层）

### 依赖升级
- **不升级**：`terraform-plugin-framework` 保持 `v1.19.0`（实测该版本已包含完整 `action` 包）；`terraform-plugin-mux` 保持 `v0.23.1`（已支持 action 协议）。
- 同步锁定最低版本约束：`terraform-plugin-framework >= v1.19.0` 且 vendor 中 MUST 存在 `action/` 包。
- `vendor/` **不需重新刷新**（无依赖变动）。

### 文档
- `tencentcloud/framework/README.md`：写明"产品/类型"双层布局、命名约定（包名前缀消歧）、"如何新增第 N 种类型资源"步骤
- 旧 `tencentcloud/provider/framework/README.md`：删除（目录已不复存在）
- `tencentcloud/provider/sdkv2/README.md`：目录树示意更新（framework 不再是 `provider/` 子目录）
- `CONTRIBUTING.md`：所有 `tcfwhelper` / `tcfwprovider` / `services/tcprovider` / `provider/framework` 文案替换
- 现有 openspec/changes 中**正在进行**的 change（如 `add-plugin-framework-muxing` 若状态非 done）保留旧名，**仅当其状态为 done 后**由本 change 的 archive 步骤一并刷新

## Capabilities

### New Capabilities
- `framework-provider-types-and-naming`：定义 framework provider 侧的**包命名约定**（`framework/internal/helper`、`framework/acctest`、`internal/sharedmeta`、顶层 `tencentcloud/framework/`，Go 包名 `framework`）、**入口与业务实现合并**到同一目录、**产品/类型双层目录**的强制布局、**registry 直接聚合各产品工厂**的契约，以及**最低 plugin-framework 版本**约束（`>= v1.19.0`，且 vendor 中 MUST 包含 `action/` 包以支持 `action` 类型注册）。

### Modified Capabilities
*（无）* — 上一个 change `add-plugin-framework-muxing` 落地的 spec 尚未 archive 进 `openspec/specs/`，因此本次不存在需要 modify 的已存档 capability；若该 change 在本 change archive 之前先行 archive，本 change 会在 archive 阶段补充对其 spec 的 delta；否则两个 change 的 spec 各自独立。

## Impact

### 受影响代码
- **真实代码引用替换**（重命名 + 路径合并）：
  - `tencentcloud/provider.go`（SDKv2 providerConfigure 中 `tcfwprovider.SetSharedMeta` 调用 + import）
  - `tencentcloud/framework/acctest/factories.go`（原 `tencentcloud/acctest/framework_factories.go` 迁入 framework 子树，包名 `frameworkacctest`）
  - `tencentcloud/framework_provider_test.go` / `framework_provider_testhelpers_test.go`（迁移路径 + import）
  - `main.go`（import 别名 `fwprovider` → `framework`）
- **目录搬迁**：
  - `internal/tcfwhelper/` 下 `.go` 文件 → `tencentcloud/framework/internal/helper/`（两步完成：先 `package tcfwhelper` → `package frameworkhelper` 且迁至 `internal/frameworkhelper/`，再迁至 `framework/internal/helper/` 且包名改为 `helper`；最终只保留 `framework/internal/helper/` 一份）
  - `internal/tcfwprovider/` 下 `.go` 文件 → `internal/sharedmeta/`（同时 `package tcfwprovider` → `package sharedmeta`）
  - `services/tcprovider/framework.go`（registry 聚合点） → 内容合并进 `tencentcloud/framework/registry.go`，原文件**删除**
  - `services/tcprovider/data_source_tc_provider_runtime_framework.go` → `tencentcloud/framework/meta/datasources/provider_runtime_data_source.go`（包名 `metadatasources`）
  - `services/tcprovider/data_source_tc_provider_runtime_framework_test.go` → `tencentcloud/framework/meta/datasources/provider_runtime_data_source_test.go`
  - `provider/framework/provider.go` → `tencentcloud/framework/provider.go`（包名仍 `framework`）
  - `provider/framework/registry.go` → `tencentcloud/framework/registry.go`（包名仍 `framework`，import 路径全改）
  - `tencentcloud/acctest/framework_factories.go` → `tencentcloud/framework/acctest/factories.go`（包名 `frameworkacctest`；SDKv2 共用的 `basic.go` / `test_util.go` 不动）
  - 旧 `services/tcprovider/`、`provider/framework/`、`internal/frameworkhelper/` 三个目录最终**全部删除**
- **文档替换**：`README.md`、`CONTRIBUTING.md`、`tencentcloud/provider/sdkv2/README.md`、`tencentcloud/framework/README.md`（新建）、`website/docs/d/provider_runtime.html.markdown`

### API / 行为
- **对 Terraform 终端用户：零影响**。所有改动都是 Go 包内部路径与文件组织；HCL 配置、provider schema、resource/datasource type name (`tencentcloud_*`)、state shape 全部保持不变。
- **对仓内开发者**：重命名后的 import 路径必须更新；新增 framework 资源时直接在 `tencentcloud/framework/<product>/<type>/` 下落地，不再进 `tencentcloud/services/`，也不再进 `tencentcloud/provider/framework/`；提交 PR 时 `go build ./...` 必须 zero error，`go vet ./...` 必须不新增错误（以 baseline = 19 为上限，且本 change 触及的包中 0 error）。

### 依赖
- `github.com/hashicorp/terraform-plugin-framework`：保持 `v1.19.0`（实测已包含完整 `action` 包，无需升级）
- `github.com/hashicorp/terraform-plugin-mux`：保持 `v0.23.1`（已发布 `tf5muxserver/mux_server_InvokeAction.go`，支持 action 协议）
- `github.com/hashicorp/terraform-plugin-go`：保持 `v0.31.0`
- `vendor/` 无变动

### 测试
- `go test ./tencentcloud/framework/internal/helper/...`
- `go test ./tencentcloud/internal/sharedmeta/...`
- `go test ./tencentcloud/framework/...`（含 6 类型 reference 的 unit test 与 provider 同包测试）
- `go vet ./...` 输出 error 数 ≤ baseline 19，且本 change 触及的包 0 error；`gofmt -l ./tencentcloud/framework/... ./tencentcloud/internal/sharedmeta/...` 输出为空
- 不需要 acc test（reference 实现都不打云 API）
