## Context

仓内 plugin-framework 改造分两步走：

1. **第 1 步**：`add-plugin-framework-muxing` 已 archive，落地了 `tf5muxserver` 脚手架与 framework 入口。
2. **第 2 步**：`restructure-framework-types-and-naming` 起草了 framework 业务实现的目录/包/命名规范——但代码侧 `tencentcloud/framework/meta/` 树**实际未落地**（spec/tasks 标记完成而文件系统空），处于"提案态"。

旧规范（第 2 步）选择的关键约束：

- **目录**：`tencentcloud/framework/<product>/<type>/`（双层）。
- **包名**：`<product><type>` 拼接（如 `metaresources` / `metaephemerals`），目的"避免 registry 中 import alias"。
- **文件名**：类型后缀法（`*_resource.go` / `*_ephemeral_resource.go` / `*_action.go` 等）。

新规范（本 change）回归与 SDKv2 现行 `tencentcloud/services/<product>/` 完全一致的心智：

- **目录**：`tencentcloud/framework/<product>/`（单层）。
- **包名**：`<product>` 直接同名（如 `meta` / `ssm` / `cvm`），靠 import path 区分，依然不需要 alias。
- **文件名**：类型前缀法（`<type>_tc_<product>_<name>.go`）。

约束方：本仓 framework 子树目前**只有一个真实非 reference 业务诉求**——`tencentcloud_ssm_secret_version`（ephemeral 类型）。reference 实现 6 个类型还没落地，是改规范的最佳窗口期。

## Goals / Non-Goals

**Goals:**

- 让 framework 业务实现的目录/包/命名与 SDKv2 现行规则**完全对齐**，消除"两套心智"。
- 用一份 ADDED Requirements 完整定义新规范，覆盖旧 change 中所有"双层 / 拼接包名 / 后缀法"条款。
- 同步硬改写旧 change 目录下的全部 Markdown 文档，避免后续读到的人被旧规范误导。
- 为后续 `tencentcloud_ssm_secret_version` 等真实 framework reference 的落地铺平路径，apply 阶段一次到位。

**Non-Goals:**

- **不改任何运行时代码**：旧 change 的代码搬迁实际未发生，本 change 也不做搬迁动作；reference 实现由后续独立 apply 阶段按新规范一次创建。
- **不改 mux / framework provider 入口**：`provider.go` / `registry.go` / `provider_test.go` 顶层 5 个文件（package 为 `framework`）位置和 package 名不变。
- **不动 helper / sharedmeta / 依赖版本**：`framework/internal/helper/`（包名 `helper`，受 `internal/` 可见性约束）与 `tencentcloud/internal/sharedmeta/` 的设计沿用旧 change 的决策。
- **不改 SDKv2 侧 `tencentcloud/services/<product>/` 任何文件**。
- **不动 `openspec/specs/`**：旧 change 的 spec 还没 archive 进去，新规范作为 New Capability 直接登场即可；本 change 也不做 archive。

## Decisions

### D1：目录布局——单层 vs 双层

**决策**：选**单层** `tencentcloud/framework/<product>/`。

**理由**：

| 维度 | 双层 `<product>/<type>/` | 单层 `<product>/`（本决策） |
|---|---|---|
| 与 SDKv2 心智 | 不一致（services 是单层） | 一致 |
| 文件密度 | 每个产品下 6 个子目录、空目录留空 | 单一目录、有什么放什么 |
| import 路径 | `framework/meta/resources` | `framework/meta` |
| 同产品跨类型代码协作 | 跨子目录 | 同包，可共享 unexported helper |
| 可视化清单 | "type 子目录"是冗余信息（文件名前缀已表达） | 文件名一眼能扫到所有类型 |

**Alternatives 考虑过：**

- **AWS 风格按服务再按类型**：上游 AWS provider 也用过双层，但近版本逐步收敛；且 AWS 一个服务下文件量级远大于本 provider 的 framework 子树（reference + 个位数 ephemeral），双层的"分组"价值不显著。
- **保留双层但去掉拼接包名**：例如 `framework/meta/resources/` 包名直接叫 `resources`——会引入跨产品同名包，registry import 必须 alias，反而更糟。

### D2：Go 包名——`<product>` 直接同名 vs 拼接

**决策**：选**直接同名**：`framework/ssm/` 包名 `ssm`，`framework/meta/` 包名 `meta`。

**理由**：

- Go 不要求"全局唯一包名"，只要求**同一文件内** import path 的 last segment 不冲突；即使存在 `services/ssm` 和 `framework/ssm`，只要不在同一文件 import，就不需要 alias。
- registry 的 import 块只引 `framework/<product>/...`，不会同时引 `services/<product>/...`，**天然不需要 alias**。
- 对人类可读性最好（与目录名一致、与 SDKv2 一致、与 AWS framework 子树一致）。

**Alternatives 考虑过：**

| 选项 | 评价 |
|---|---|
| `package fwssm` 前缀 | 多一套规则、可读性差，被否 |
| `package ssmfw` 后缀 | 同上 |
| `package <product><type>` 拼接 | 旧规范，已 D1 一并废弃 |

**验证**：万一未来某文件**真的**同时需要 `services/ssm` 和 `framework/ssm`（极罕见，几乎只在跨层桥接代码出现），此时**该文件**自行加 alias（如 `fwssm "github.com/.../tencentcloud/framework/ssm"`）即可，**不影响全局规范**。

### D3：文件名——类型前缀法 vs 后缀法

**决策**：选**类型前缀法**：`<type>_tc_<product>_<name>.go`。

**理由**：

- 与 SDKv2 现行 `resource_tc_<product>_<name>.go` / `data_source_tc_<product>_<name>.go` **一字不差**对齐，新人零迁移成本。
- 在 `ls` / IDE 文件树排序中，**同类型自动聚簇**（resource 一起、ephemeral 一起），优于后缀法（按业务名字典序混排）。
- gendoc 等仓内工具依赖文件名前缀做分类的，沿用同一套约定即可。

**类型前缀全集（与 SDKv2 + framework 6 类型对齐）：**

| 类型 | 文件前缀 |
|---|---|
| resource | `resource_tc_` |
| datasource | `data_source_tc_` |
| function | `function_tc_` |
| ephemeral | `ephemeral_tc_` |
| list | `list_tc_` |
| action | `action_tc_` |

测试文件统一加 `_test.go` 尾缀，doc 文件统一同名 `.md`（gendoc 约定）。

### D4：旧 change 文档处理——硬改写 vs 软勘误

**决策**：选**硬改写**。

**理由**：

- 旧 change 的代码搬迁实际未落地，"双层 / 拼接 / 后缀"只存在于文字。
- 旧 change `tasks.md` 中 `[x]` 标记的"reference 实现要落地"业务意图**不变**，只是落到的路径/包名/文件名换了——硬改写是把旧文档**升级到新规范**，不是回退。
- 软勘误（保留旧文字 + 加 Superseded 横幅）会让"两份对立的规范文字"长期共存，正是用户提到的"防止后续误解"想避免的局面。

**保留**：旧 change 目录的存在本身、`[x]` 标记、CLI 状态——这些是历史决策的痕迹，不影响新规范的清晰度。

### D5：本 change 走 New Capability 路径

**决策**：在 `specs/framework-provider-layout-and-naming/spec.md` 用 `## ADDED Requirements` 一次性给出完整新规范。

**理由**：旧 change 的 spec **尚未 archive 进 `openspec/specs/`**（grep 已确认 `openspec/specs/framework-provider-types-and-naming/` 不存在），从 OpenSpec 角度它仍在 changes 阶段，没有"已发布"的旧 capability 供本 change 做 MODIFIED delta。新 capability 直接 ADDED 是最干净的形式。

## Risks / Trade-offs

- **[R1] 旧 change 仍然存在于 `openspec/changes/` 目录** → 本 change 硬改写其文档；同时本 change 的 proposal/design 明示"取代"关系；后续如果有人执行 `openspec archive` 系列命令，应优先 archive 本 change，旧 change 因 spec 已被本 change 完整覆盖、archive 时直接丢弃即可（archive 阶段再决定）。
- **[R2] "framework/ssm" 与 "services/ssm" 同包名导致 grep 误匹配** → 风险低：所有跨层调用都走完整 import path，grep 时用 `framework/<product>` 或 `services/<product>` 路径片段过滤即可，无歧义。
- **[R3] 文件名前缀 `ephemeral_tc_` / `function_tc_` / `list_tc_` / `action_tc_` 是新发明** → 不是发明：`resource_tc_` / `data_source_tc_` 是 SDKv2 既有约定，本 change 把同一套规则延伸到 plugin-framework 新增的 4 类。延伸方式与既有约定**完全同构**（前缀 = 类型词、连接符 = `_tc_`、尾段 = `<product>_<name>`），不引入新心智。
- **[R4] 旧 change 文档硬改写后，未来 git blame 会指向本 change 而非原作者** → 可接受：本 change proposal 已显式记录了"取代"语义，blame 能追溯回 PR 的上下文；且 archive 阶段会沉淀完整变更记录。

## Migration Plan

1. 本 change apply 阶段（独立任务）按新规范创建 `tencentcloud/framework/<product>/<file>.go` 的 reference 实现与首个真实 ephemeral（`framework/ssm/ephemeral_tc_ssm_secret_version.go`）。
2. 旧 change 目录内的文档由本 change tasks 阶段**一次性**改写完成，不分批、不保留旧措辞。
3. 回滚策略：本 change 仅改文档，回滚等价于 `git revert` 本 change 的 commits；不会破坏任何运行时行为（因为旧 change 的代码本来就没落地）。

## Open Questions

无。决策 D1–D5 已全部锁定。
