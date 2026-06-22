## Context

`add-plugin-framework-muxing` 落地的 framework provider 链路（main.go → tf5muxserver → SDKv2 + framework）目前能正常工作，但其在仓内的物理布局是"先做出来再说"的状态：

- **过渡命名**：`internal/tcfwhelper`、`internal/tcfwprovider`、`services/tcprovider`、import 别名 `fwprovider`，都为了避开 SDKv2 既有同名包临时加了 `tcfw` / `fw` 前缀。
- **三处分散布局**：framework provider 入口在 `tencentcloud/provider/framework/`、registry 聚合点与第一个具体的 framework data source 在 `tencentcloud/services/tcprovider/`、provider 测试在 `tencentcloud/` 顶层目录。后续要新增 `resource` / `function` / `ephemeral` / `list` / `action` 时，开发者要在三个目录之间反复横跳。
- **缺产品维度**：现状只有 1 个 datasource 也就罢了，一旦 framework 资源开始覆盖多个云产品（CVM / VPC / CBS / ...），单层"按类型"分目录会再次撑爆。
- **缺示范**：除 1 个 datasource 外，其他 5 种类型在仓里没有任何 reference，新人无法照抄。
- **依赖状态预检结论**：当前 `terraform-plugin-framework v1.19.0` 上游**已包含完整 `action` 包**（该项在本 change 启动后的实测中发现，与处理 propose 阶段依靠 changelog 推断的结论相反）；`terraform-plugin-mux v0.23.1` 同步已支持 action 协议。因此**依赖升级不再是本 change 的必要件**；仅需锁定最低版本约束为当前锁定的 v1.19.0 / v0.23.1。

利益相关方：
- 后续在 framework 侧扩展产品能力的 PR 作者（需要稳定的目录约定 + 6 类型 reference）
- 阅读源码的外部贡献者（被 `tcfw*` 命名劝退、被三处分散布局绕晕）
- CI / 构建系统（需要兼容 framework 升级后的 mux 协议）

## Goals / Non-Goals

**Goals:**
- 全仓一次性完成 `tcfwhelper` / `tcfwprovider` / `services/tcprovider` / `provider/framework`（目录） / `fwprovider` 别名 5 处过渡命名/路径的重排，最终 framework 一切（入口 + 业务实现 + 测试）统一落到 `tencentcloud/framework/`，包名为 `framework`。
- 在 `tencentcloud/framework/` 下确立"产品（service）→ 类型"双层强制布局：第一层是产品文件夹（`cvm/` / `meta/` / 未来 `vpc/` 等），第二层是 plugin-framework 的 6 种资源类型子目录（`resources/` / `datasources/` / `functions/` / `ephemerals/` / `lists/` / `actions/`）。
- 为 6 种类型各落地 1 份 reference 实现（datasource 复用现有迁移；其余 5 种用本地 stub 实现，避免引入云 API 依赖）。
- 锁定 `terraform-plugin-framework` 最低版本为当前锁定的 v1.19.0（已含 `action` 包）；**不进行依赖升级**，`vendor/` 无变动。
- registry / provider 接口同步暴露 `Actions(ctx)`，并直接 import 各产品包工厂方法做聚合（不再有中间层 `framework.go`）。
- 测试文件迁入 `tencentcloud/framework/` 包内，与 provider 生产代码同包，方便共享未导出 helper。

**Non-Goals:**
- **不**改 Terraform 用户可见的任何 schema / state / type name。
- **不**新增任何真实云 API 调用 — reference 实现保持纯本地逻辑或 stub。
- **不**重写 SDKv2 那一侧的命名 / 布局（`services/<service>/` 仍按既有 `resource_tc_*` / `data_source_tc_*` 风格）。
- **不**触碰 `add-plugin-framework-muxing` / `extend-provider-runtime-credential-meta` 这两个未 archive 的 change 的内容（它们各自的 spec 增量保留原貌；本 change archive 时若它们已先 archive，则在 archive 阶段单独处理 delta 重叠）。
- **不**为 reference 实现写 acc test（`TF_ACC=1` 的真正验收测试不在本 scope）。

## Decisions

### Decision 1：包/目录最终命名

**选**：
- `internal/tcfwhelper` → `internal/frameworkhelper`
- `internal/tcfwprovider` → `internal/sharedmeta`
- `services/tcprovider`（业务实现） → 拆解后并入 `tencentcloud/framework/<product>/<type>/`
- `provider/framework`（provider 入口） → 下沉合并到 `tencentcloud/framework/`（与业务实现同目录）
- import 别名 `fwprovider` → 直接使用包名 `framework`

**理由**：
- `frameworkhelper` 字面表达"plugin-framework 侧的 helper"，去掉 `tcfw` 缩写歧义；位于 `internal/` 已经隐含"本仓内部"。
- `sharedmeta` 精确表达职责（在 SDKv2 与 framework 两侧共享 `*connectivity.TencentCloudClient` 的 meta 桥接），比 `tcfwprovider` 更贴切。
- `tencentcloud/framework/`（与 `tencentcloud/services/` 同级、包名 `framework`）一处目录承载 framework 的全部内容（入口 + 业务 + 测试）：阅读、跳转、导入路径都最短。
- 测试中保留 `framework` 包名，可直接 `import ".../tencentcloud/framework"` 后写 `framework.NewProvider(...)`，不再需要别名。

**已考虑的替代方案**：
- ❌ 保留 `tencentcloud/provider/framework/`（仅放入口） + 新增 `tencentcloud/framework/`（业务实现）：仍然两处分散，新增类型时心智负担没消除。
- ❌ `services/tcframework/` 子目录承载 framework：失去与 `services/`（SDKv2 专属）的语义区隔；`tc` 前缀又不必要。
- ❌ `internal/pluginhelper` / `internal/pf`：太泛，看不出是 plugin-framework 还是 plugin-sdk。
- ❌ 完全保留 `tcfw*`：用户已明确表态不喜欢。

### Decision 2：目录采用"产品 → 类型"双层布局

**选**：在 `tencentcloud/framework/` 下先按云产品分（`cvm/` / `meta/` / 未来 `vpc/` 等），再在每个产品下按 plugin-framework 6 大类型分子包：

```
tencentcloud/framework/
├── provider.go                # 入口（包名 framework）
├── registry.go                # 6 类型聚合（直接 import 各产品工厂）
├── provider_test.go
├── testhelpers_test.go
├── README.md
├── cvm/                       # 真实云产品：CVM
│   └── actions/
│       └── reboot_instance_action.go     # package cvmactions
├── meta/                      # 元产品：跨产品 / 不归属任何具体云产品
│   ├── resources/             # package metaresources
│   ├── datasources/           # package metadatasources
│   ├── functions/             # package metafunctions
│   ├── ephemerals/            # package metaephemerals
│   └── lists/                 # package metalists
└── ...                        # 未来按需新增 vpc/ / cbs/ / ...
```

**理由**：
- terraform-provider-aws、terraform-provider-google 等大型 provider 都以"按 service 分"为第一层（aws 是 `internal/service/<svc>/`，google 是 `services/<svc>/`），可读性 / 影响域 / 单测组织都最佳；
- plugin-framework 的 6 种资源类型是个一级概念（接口签名互不相同），用文件后缀区分需要读 import 才能确认，不如目录直接表达。第二层按类型分能进一步把"接同一种 framework 接口的实现"聚拢；
- 子包之间互不依赖，框架升级或 mux 协议变更时影响域更小；
- 子目录**按需创建**（不需要为没内容的类型预留空目录 / `.gitkeep`）。

**关于产品归属**：
- 严格规则：能明确归属真实云产品的，必须落到对应产品目录（如 reboot_instance → `cvm/`）；
- 跨产品 / 不归属任何产品的（local note、provider runtime、parse_resource_id、temp_credential、region），统一放 `meta/`；
- "meta" 命名灵感来自 "provider meta"，表达"provider 自身层面、跨产品"。

**已考虑的替代方案**：
- ❌ 单层"按类型"（`framework/resources/` / `framework/actions/` ...）：跨多产品后会撑爆单一目录。
- ❌ 单层"按产品"（`framework/cvm/<file>.go`）：失去类型维度，文件命名约定混乱。
- ❌ 反向"按类型 → 按产品"（`framework/resources/cvm/` / `framework/actions/cvm/`）：违反业界主流（AWS / Google），新人陡峭。
- ❌ 把 reference 全部塞到独立的 `examples/` 目录：违反"reference 必须接入 registry，可被 schema 检索到"（Decision 4）。

### Decision 3：包命名采用"产品前缀 + 类型名"消歧

**选**：类型子目录的 Go 包名 = `<product><type>`：
- `framework/cvm/actions/` → `package cvmactions`
- `framework/meta/resources/` → `package metaresources`
- `framework/meta/datasources/` → `package metadatasources`
- `framework/meta/functions/` → `package metafunctions`
- `framework/meta/ephemerals/` → `package metaephemerals`
- `framework/meta/lists/` → `package metalists`

**理由**：
- 多产品下的同类子包（如 `cvm/actions` 与 `vpc/actions`）若都用 `package actions`，registry.go 中 import 时必须给每个加 alias；用 `<product><type>` 命名后 alias 不再需要，import 块可读性最佳。
- 与 Go 社区"包名 = 调用方写的名字"哲学一致：调用方写 `cvmactions.NewRebootInstanceAction`、`metaresources.NewLocalNoteResource`，语义自解释。
- 产品文件夹本身（`cvm/` / `meta/`）**不放 `.go` 文件**，因此不需要 package 声明，也不会与子包撞名。

**已考虑的替代方案**：
- ❌ 类型子目录都用 `package resources` / `package actions`：必须 alias，import 块累赘。
- ❌ 包名直接用产品名（`package cvm`）：失去类型语义，且与 `tencentcloud/services/cvm`（SDKv2 侧）撞名。

### Decision 4：registry 直接聚合，取消中间层 `framework.go`

`tencentcloud/framework/registry.go` 内部直接 import 各产品包工厂，做 slice 聚合：

```go
package framework

import (
    cvmactions "github.com/.../tencentcloud/framework/cvm/actions"
    metadatasources "github.com/.../tencentcloud/framework/meta/datasources"
    metaresources "github.com/.../tencentcloud/framework/meta/resources"
    metafunctions "github.com/.../tencentcloud/framework/meta/functions"
    metaephemerals "github.com/.../tencentcloud/framework/meta/ephemerals"
    metalists "github.com/.../tencentcloud/framework/meta/lists"
)
```

provider.go 的 `Resources` / `DataSources` / `Functions` / `EphemeralResources` / `ListResources` / `Actions` 6 个方法直接 `return r.Resources()` / ... 形式，registry 是唯一聚合点。

**理由**：
- 入口（provider.go）已经在 framework 顶层；registry.go 也在同包，直接做最终聚合，消除原 `services/tcprovider/framework.go` 这一层无收益的转发。
- 未来要支持多产品时只需在 registry.go 的 import 块和 slice 中再加一组工厂，扩展点单一。

### Decision 5：reference 接入级别 = L2（注册到 registry）

每种类型的 reference 都通过 `registry.go` 注册进 provider，**真实出现在 provider schema 里**。

**理由**：
- 用户明确要求"生成每种类型的示例"，且未要求 example/ 隔离；
- L2 能覆盖最完整的链路（registry → provider.Configure 注入 meta → 资源 Configure → CRUD），后续开发者照抄能直接跑；
- 通过"业务名 + 仅本地逻辑"双重约束，避免污染用户真实数据。

**reference 实现的副作用控制**：

| 类型 | type name | 产品归属 | 实现策略 |
|---|---|---|---|
| resource | `tencentcloud_local_note` | meta | 本地 in-memory map，进程内持久；演示完整 CRUD + 状态 ID |
| datasource | `tencentcloud_provider_runtime` | meta | **搬迁现有实现**，行为零变化 |
| function | `parse_resource_id` | meta | 纯字符串拆分，无 IO |
| ephemeral | `tencentcloud_temp_credential` | meta | 从 `sharedmeta.GetSharedMeta()` 读 client 字段（如 region），生成一个 5 分钟过期的 fake token；不调真实 STS |
| list | `tencentcloud_region` | meta | 返回硬编码区域列表（与 `connectivity` 包里的常量保持一致），无 IO |
| action | `tencentcloud_reboot_instance` | cvm | 仅校验 `instance_id` 入参格式（regex `^ins-[a-z0-9]+$`），打日志，**不**调 CVM RebootInstances API |

### Decision 6：保持 plugin-framework 于 v1.19.0（不升级）

**选**：`terraform-plugin-framework` 保持在 `v1.19.0`，`terraform-plugin-mux` 保持在 `v0.23.1`；不执行 `go get` / `go mod tidy` / `go mod vendor`。

**理由**：
- 实测发现：`vendor/github.com/hashicorp/terraform-plugin-framework/action/` 在 v1.19.0 上游**已包含完整 `action` 包**（`action.go`、`configure.go`、`invoke.go`、`modify_plan.go`、`schema.go`、`schema/`、`validate_config.go`、`config_validator.go`、`deferred.go`、`metadata.go` 齐备），`provider/provider.go` 中 `ProviderWithActions` 接口存在。原 propose 阶段依靠 changelog 推断的"v1.19 不含 action 包"是错误的。
- 实测发现：`vendor/github.com/hashicorp/terraform-plugin-mux/tf5muxserver/mux_server_InvokeAction.go` 在 v0.23.1 已发布，表明当前 mux 版本已支持 action 协议。
- 依赖升级需要外网 GOPROXY、会产生庞大 `vendor/` diff、可能引入 v1.19→v1.20 隐式 break；既然当前状态已满足本 change 的全部诉求，升级纯属收益为零、风险为正。
- 本 change 的核心价值在于"命名规整 + 目录扁平化 + 产品分层 + 6 类型示范"，与 framework 的 minor 版本无必然耦合。

**已考虑的替代方案**：
- ❌ 升级到 v1.20.0（原 propose 方案）：发现现状已满足后不再需要；升级会带来外网依赖、vendor diff 股胀、隐式 break 三重额外风险。
- ❌ 不升级但恶意锁定为更高版本上限：无意义。

**后续升级路径**：如后续某个 change 明确需要 v1.20+ 才能使用的 API（例如 v1.20+ 新增的 schema validator），可在该 change 中独立完成升级；本 change 的 spec 如果在 archive 后被后续 change MODIFY，也可以在那个后续 change 中同步抬升最低版本约束。

### Decision 7：测试代码搬迁后的 package 选择

**选**：`framework_provider_test.go` / `framework_provider_testhelpers_test.go` 搬到 `tencentcloud/framework/` 后，使用 `package framework`（同生产代码包），而非 `package framework_test`。

**理由**：
- 测试需要直接构造内部辅助类型（如调用 `NewProvider(primary)` 的内部 helper、读取未导出的 schema 常量）；
- 与同包 `provider.go` 中已有的非导出 helper 共享更方便；
- 风险（循环 import）已通过现有代码验证 — 测试不再 import 上层 `tencentcloud` 包，仅 import `internal/sharedmeta` 和 mocks，链路清晰。

### Decision 8：搬迁顺序（避免中间态编译失败）

按"包级别原子搬迁"逐个落，每一步保证 `go build ./...` zero error且 `go vet ./...` 不新增错误（以 baseline = 19 为上限），便于二分回滚：

1. 依赖预检与基线记录（无依赖升级，包括记录 vet baseline = 19）
2. `internal/tcfwhelper` → `internal/frameworkhelper`（含 `package` 改名 + 全仓 import 替换 + 自身 test）
3. `internal/tcfwprovider` → `internal/sharedmeta`（同 2）
4. **新建** `tencentcloud/framework/` 目录骨架（暂时空，仅作为目标容器）
5. `tencentcloud/provider/framework/provider.go` + `registry.go` → `tencentcloud/framework/provider.go` + `registry.go`（下沉合并；包名仍 `framework`）；调整 import 路径与 `fwprovider` 别名移除；删除空目录 `tencentcloud/provider/framework/`
6. 测试文件迁出 → `tencentcloud/framework/provider_test.go` / `testhelpers_test.go`
7. `tencentcloud/services/tcprovider/data_source_tc_provider_runtime_framework.go` → `tencentcloud/framework/meta/datasources/provider_runtime_data_source.go`（包名 `metadatasources`）；`framework.go` 中间层删除并入 registry.go；删除空目录 `services/tcprovider/`
8. 新增 5 个 reference（按产品归属落到 `cvm/actions/` 或 `meta/<type>/`）+ registry 接入 action
9. 文档替换 + `vendor/` 刷新

## Risks / Trade-offs

| Risk | Mitigation |
|---|---|
| **plugin-framework v1.19.0 vendor 中 action 包隐式不稳定**（虽然齐备，但上游在 v1.19 阶段可能未完全 GA-tier） | reference action 文件实现保持最小（只 `Metadata` / `Schema` / `Configure` / `Run`），封装面 ≤ 100 行；若未来发现接口 break，在独立 commit 中完成升级到 ≥ v1.20.0，本 change 的其他产出不受影响 |
| **mux 协议在 v0.23.1 的 action 支持不稳定** | 同上；本 change 不跑 acc test，merge 后由现有 nightly acc 流水线兏底；若发生回归，可选择 revert reference action 或后续单独升级 mux || **大规模 import 路径替换出错**（手工漏改导致编译失败） | 使用 `goimports -local` 全仓刷 import；改完后 `go build ./...` 必须通过 |
| **`vendor/` 无变动** | 本 change 不升级依赖，`vendor/` 不变动，PR diff 只含仓内代码重组 |
| **reference 实现污染 schema**（用户误以为是真资源） | (a) 命名上用 `tencentcloud_local_note` / `tencentcloud_region` 等显式领域名；(b) 在 README 与 docs 显式标注"reference 实现，本地行为"；(c) action 资源的 `Run` 内只打日志，永远返回成功；(d) 后续若决定剥离，可只删 6 个 reference 文件 + registry 中对应 6 行注册 |
| **action 接口在不同 framework patch 版本签名仍可能微调** | reference action 文件实现保持最小（只 `Metadata` / `Schema` / `Configure` / `Run`），出问题时修改面 ≤ 100 行 |
| **`provider/framework` 目录被多处文档引用，删除后链接失效** | 文档替换在 step 9 同 commit 完成；PR 描述列出所有变更链接 |
| **`add-plugin-framework-muxing` 尚未 archive，spec delta 重叠** | 本 change 仅 ADDED 新 capability `framework-provider-types-and-naming`，**不**碰它的 spec；若它在本 change archive 之前先 archive，archive 阶段重新评估是否需要补 MODIFIED delta |

## Migration Plan

### 部署
- 全仓代码改完后（无依赖变动）：
  1. `go build ./...` zero error
  2. `go vet ./...` 输出 error 数 ≤ baseline 19，且本 change 触及的包（`framework/internal/helper` / `framework/acctest` / `internal/sharedmeta` / `tencentcloud/framework/...`）中 0 error
  3. `gofmt -l ./tencentcloud/framework/... ./tencentcloud/internal/sharedmeta/...` 输出为空
  4. `go test ./tencentcloud/framework/internal/helper/... ./tencentcloud/internal/sharedmeta/... ./tencentcloud/framework/...`（无 TF_ACC）全部 PASS

### Rollback
- 本 change 仅含代码重组，无依赖升级 commit。出现问题时：
  - 如仅 reference 实现（特别是 cvm action）出现错误：单独 revert 对应 reference 文件 + registry 中该一行注册即可
  - 全量回滚：`git revert <range>`，回到改名前状态；用户配置零影响

### 通知
- 在 PR 描述中显著标注 import 路径变化（影响仓内开发者，但**不**影响 Terraform 用户）
- 在 `CONTRIBUTING.md` 与 `tencentcloud/framework/README.md` 同步新约定

## Open Questions

- **Q1**（已解决）：`terraform-plugin-mux` 当前锁定版本是否已支持 framework v1.20 的 action 协议？
  - 实测结论：v0.23.1 已发布 `tf5muxserver/mux_server_InvokeAction.go`，支持 action 协议。不需升级。
- **Q2**（已解决）：Action 类型在 `tf5muxserver` 是否被支持？v6 mux 协议是 action 的稳定通道，可能需要把 main.go 的 mux 改为 v6。
  - 实测结论：`tf5muxserver` 在 v0.23.1 中已含 `mux_server_InvokeAction.go`，action 在 v5 mux 协议上可用，本 change 维持当前 v5 mux 不动。若后续 reference action 遇到协议问题再在独立 change 中升级到 v6 mux。
- **Q3**（已解决）：`vendor/` 全量刷新可能引入数千行 diff，是否需要拆 PR？
  - 本 change 不升级依赖，`vendor/` 不变动，PR diff 仅含仓内代码重组。
- **Q4**：`meta/` 这个产品名是否有更贴切的命名（例如 `provider/`、`platform/`、`core/`）？
  - 默认采用 `meta/`，如 review 中有强烈反对意见，apply 阶段可一次性 sed 替换（成本极低）。

## 11. Decision: framework-only 代码收敛到 framework/ 子树

### 背景

Phase 2 完成 `internal/tcfwhelper` → `internal/frameworkhelper` 重命名后，user 在迭代评审中进一步提出重要改进诉求：**framework-only 的 helper 与 acctest 工厂应该跟随业务实现一起收敛到 `tencentcloud/framework/` 子树下**，只有双栈必须共享的包（`sharedmeta`）才需要留在 `internal/`。

### 选项与判断

| 选项 | helper 路径 | 包名 | 可见性机制 |
|---|---|---|---|
| A1 | `tencentcloud/framework/helper/` | `helper` | 仅依靠约定（外部包仍能 import） |
| **A2 (采纳)** | `tencentcloud/framework/internal/helper/` | `helper` | Go `internal/` 可见性硬约束仅 framework 子树可 import |
| B | `tencentcloud/framework/internal/frameworkhelper/` | `frameworkhelper` | 同 A2 internal，但包名冗余（路径与包名连起来读会包含 `framework/internal/frameworkhelper`） |

**采纳 A2** 的决策理由：

1. **可见性硬约束**：Go 语言原生 `internal/` 规则能在编译器层面拒绝仓外任何位置的 import，是最强的封装能力，完美匹配 "framework-only" 语义；lint/CI 不需额外检查。
2. **包名简洁**：路径已隐含 `framework/internal`，包名再召 `frameworkhelper` 会在调用点写成 `frameworkhelper.RetryFramework` 中间级重复信息；采用 `helper.RetryFramework` 可读性更高，与项目原有面向 SDKv2 的 `tencentcloud/internal/helper`（路径完全不同）靠路径隔离不会冲突。
3. **调用点不需别名**：framework 子树内使用者仅 `tencentcloud/framework/meta/datasources/provider_runtime_data_source.go` 一个文件错过三个 `helper` 同名冲突场景的需要，本 change scope 内可不用 alias 直接调用；未来若某文件同时要 import 两个 `helper` 包，再在当场 alias（例 `fwhelper "tencentcloud/framework/internal/helper"`）。

### sharedmeta 保留在 internal/ 的理由

`internal/sharedmeta` **不**随 helper 迁入 framework 子树。原因：sharedmeta 要同时被 `tencentcloud/provider.go`（SDKv2 入口）与 `tencentcloud/framework/provider.go` 双向 import，是双栈共享桥；一旦迁入 `framework/internal/`，SDKv2 侧会被 Go `internal/` 规则拒绝访问，双栈架构会直接坍塔。

### acctest 拆分策略

原 `tencentcloud/acctest/` 同时含：
- SDKv2 + framework 共用的 `AccPreCheck` / `basic.go` / `test_util.go`
- framework-only 的 `framework_factories.go`（`AccProtoV5ProviderFactories`）

本 change 只迁走 framework-only 的那一份 → `tencentcloud/framework/acctest/factories.go`（包名 `frameworkacctest`）；共用部分保留原位。代价是 framework 资源的 acceptance test 需要同时 import 两个包（alias 惯例：`tcacctest` + `tcfwacctest`），Go 上完全合法且可读；另一选项（置换全部 `tencentcloud/acctest/` 到 framework 下）会强迫 SDKv2 既有测试代码也 import `tencentcloud/framework/...`，违反 "SDKv2 不依赖 framework 子树" 的纪律。

### 实施路径

1. 新建 `tencentcloud/framework/internal/helper/` 与 `tencentcloud/framework/acctest/` 两个目录；
2. 8 个 `frameworkhelper/*.go` 文件 → helper/（改包名 + doc 中示例代码中 `frameworkhelper.` 改为 `helper.`）；
3. `framework_factories.go` → `framework/acctest/factories.go`（包名 `frameworkacctest`）；
4. 商业代码唯一调用点 `tencentcloud/framework/meta/datasources/provider_runtime_data_source.go` 改 import 与 6 处符号引用；
5. 唯一测试调用点 `provider_runtime_data_source_test.go` 加 `tcfwacctest` import 且将 `ProtoV5ProviderFactories` 指向新包。
6. 旧 `tencentcloud/internal/frameworkhelper/`、`tencentcloud/acctest/framework_factories.go` 物理删除（sandbox 内 `rm` 被拦截，由维护者手动完成）。

### 验证报告

- `go build ./...`：zero error
- `go vet ./...`：error 数 = 19（= baseline），本 change 触及包 0 new error
- `gofmt -l tencentcloud/framework/internal/helper/ tencentcloud/framework/acctest/ tencentcloud/framework/ tencentcloud/internal/sharedmeta/`：输出为空
- `go test -race ./tencentcloud/framework/internal/helper/... ./tencentcloud/internal/sharedmeta/... ./tencentcloud/framework/...`：10 个包全 PASS（`framework/acctest` 包仅含 生产代码、无 test 文件，在输出中显示 `[no test files]`，符合预期）
- `grep -rE "internal/frameworkhelper|frameworkhelper\." tencentcloud/ main.go`：0 匹配