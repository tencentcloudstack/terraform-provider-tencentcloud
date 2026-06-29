## 1. 硬改写旧 change `restructure-framework-types-and-naming` 的 spec.md

- [x] 1.1 修改 `openspec/changes/restructure-framework-types-and-naming/specs/framework-provider-types-and-naming/spec.md` 中 `Requirement: Framework Product/Type Two-Level Layout`：将"产品 → 类型双层目录布局"整段改写为"产品单层目录布局"，删除所有 `<product>/<type>/<file>.go` 路径表述，统一为 `<product>/<file>.go`。
- [x] 1.2 同文件 `Requirement: Framework Product/Type Two-Level Layout` 的 Scenarios：删除 "类型子目录 SHALL 按需创建"、"类型子目录列表（resources/datasources/...）" 相关表述；新增 "类型子目录 SHALL NOT 创建" 的反向硬约束。
- [x] 1.3 同文件该 Requirement 中文件名规则：将 `*_resource.go` / `*_data_source.go` / `*_function.go` / `*_ephemeral_resource.go` / `*_list_resource.go` / `*_action.go` 后缀法表述，整段替换为 `<type>_tc_<product>_<name>.go` 前缀法表述（`resource_tc_` / `data_source_tc_` / `function_tc_` / `ephemeral_tc_` / `list_tc_` / `action_tc_`）。
- [x] 1.4 同文件 `Requirement: Framework Sub-package Naming Convention`：将"包名 = `<product><type>` 拼接"整段改写为"包名 = `<product>` 直接同名"；删除 `metaresources` / `metadatasources` / `metafunctions` / `metaephemerals` / `metalists` / `cvmactions` 6 个示例，替换为 `meta` / `cvm` / `ssm` 等示例；删除"避免不同产品下同名类型子包在 registry 中需要 alias"的旧理由（同名包靠 import path 区分天然不需要 alias）。
- [x] 1.5 同文件 `Requirement: Framework Registry Direct Aggregation`：将"import 各产品类型子包（如 cvmactions / metaresources）"改为"import 各产品包（如 cvm / meta / ssm）"；保留"无中间聚合层"硬约束不变。
- [x] 1.6 同文件 `Requirement: Reference Implementations Per Type`：6 个 Scenario 的检查路径全部由 `tencentcloud/framework/meta/<type>/` 改为 `tencentcloud/framework/meta/`（list 与 action 同步）；建议文件名同步更新为 `resource_tc_meta_local_note.go` / `data_source_tc_meta_provider_runtime.go` / `function_tc_meta_parse_resource_id.go` / `ephemeral_tc_meta_temp_credential.go` / `list_tc_meta_region.go` / `action_tc_cvm_reboot_instance.go`。
- [x] 1.7 同文件 `Requirement: Framework Provider Action Support` 与 `Requirement: Plugin Framework Minimum Version` 与 `Requirement: Framework Test Co-location` 与 `Requirement: No Import Alias Required`：通读，把任何残留的"双层"或"后缀法"措辞同步刷新（这几条主要约束不变，仅做文字一致性）。
- [x] 1.8 同文件文件级搜索：`grep -nE 'metaresources|metadatasources|metafunctions|metaephemerals|metalists|cvmactions|<product><type>|<product>/<type>/|_ephemeral_resource\.go|_data_source\.go|_list_resource\.go' specs/framework-provider-types-and-naming/spec.md`，验证命中条数为 0。

## 2. 硬改写旧 change `restructure-framework-types-and-naming` 的 design.md

- [x] 2.1 修改 `openspec/changes/restructure-framework-types-and-naming/design.md` 顶部"目录树示例"图示：将 `meta/{resources/, datasources/, functions/, ephemerals/, lists/}` 的 5 个子目录全部抹平到 `meta/` 单层；同时把每个文件名换成 `<type>_tc_meta_<name>.go` 前缀法。
- [x] 2.2 同文件"包名拼接"小节：删除"`framework/meta/resources/` → `package metaresources`"等 5 条映射；改写为"`framework/meta/` → `package meta`"等 1 条映射，并补充"同名包靠 import path 区分、无需 alias"的简短解释。
- [x] 2.3 同文件"反方案对比"小节：保留"反向按类型→按产品"被否的论述；新增一段"双层方案被否"的说明，指向新 change `restructure-framework-single-level-layout` 的 design.md 决策 D1（无需展开重复理由）。
- [x] 2.4 同文件 `services/tcprovider/` 搬迁路径表述：将 `framework/meta/datasources/provider_runtime_data_source.go` 全部改为 `framework/meta/data_source_tc_meta_provider_runtime.go`（影响约 5 处出现位置）。
- [x] 2.5 同文件文件级搜索：`grep -nE 'metaresources|metadatasources|metafunctions|metaephemerals|metalists|cvmactions|/resources/|/datasources/|/functions/|/ephemerals/|/lists/|/actions/|_ephemeral_resource\.go|_data_source\.go' design.md`，验证命中条数为 0（仅限 framework 子树相关上下文；如有 services 侧 SDKv2 文件名命中可保留）。

## 3. 硬改写旧 change `restructure-framework-types-and-naming` 的 tasks.md

- [x] 3.1 修改 `openspec/changes/restructure-framework-types-and-naming/tasks.md` "## 7." 节（搬迁现有 datasource）：路径 `framework/meta/datasources/provider_runtime_data_source.go` 全部改写为 `framework/meta/data_source_tc_meta_provider_runtime.go`；包名指令 `package metadatasources` 改为 `package meta`；测试包 `package metadatasources_test` 改为 `package meta_test`。
- [x] 3.2 同文件 "## 8." 节（创建 5 类 reference）：6 个 reference 文件路径与名称全部同步刷新到新规范——`framework/meta/resources/local_note_resource.go` → `framework/meta/resource_tc_meta_local_note.go`；`framework/meta/ephemerals/temp_credential_ephemeral_resource.go` → `framework/meta/ephemeral_tc_meta_temp_credential.go`；`framework/meta/lists/region_list_resource.go` → `framework/meta/list_tc_meta_region.go`；新增 `framework/meta/function_tc_meta_parse_resource_id.go` 与 `framework/cvm/action_tc_cvm_reboot_instance.go`（如旧 tasks 已有相应行则改写，无则补齐）。
- [x] 3.3 同文件 "## 13." 节（registry import 调整）：`framework/meta/datasources/provider_runtime_data_source.go` 改为 `framework/meta/data_source_tc_meta_provider_runtime.go`；registry import path 由 `framework/meta/datasources` 改为 `framework/meta`；符号引用由 `metadatasources.NewXxx` 改为 `meta.NewXxx`。
- [x] 3.4 同文件全文件级搜索：`grep -nE 'metaresources|metadatasources|metafunctions|metaephemerals|metalists|cvmactions|framework/meta/(resources|datasources|functions|ephemerals|lists)/|framework/cvm/actions/|_ephemeral_resource\.go|_data_source\.go|_list_resource\.go|_action\.go' tasks.md`，验证命中条数为 0。

## 4. 硬改写旧 change `restructure-framework-types-and-naming` 的 proposal.md

- [x] 4.1 修改 `openspec/changes/restructure-framework-types-and-naming/proposal.md`：将 What Changes / Capabilities / Impact 中关于"双层目录"、"`<product><type>` 拼接包名"、"`*_xxx_resource.go` 后缀法"的措辞统一改写为单层 / 同名包 / 前缀法；如有 6 个 reference 的预期文件路径示例，同步使用新文件名。
- [x] 4.2 同文件文件级搜索：`grep -nE 'metaresources|metadatasources|metafunctions|metaephemerals|metalists|cvmactions|/resources/|/datasources/|/functions/|/ephemerals/|/lists/|_ephemeral_resource\.go|_data_source\.go|_list_resource\.go' proposal.md`，验证命中条数为 0。

## 5. 跨文档一致性回归

- [x] 5.1 在仓内执行 `grep -rEn 'metaresources|metadatasources|metafunctions|metaephemerals|metalists|cvmactions' -- 'openspec/changes/restructure-framework-types-and-naming/' 'openspec/changes/restructure-framework-single-level-layout/'`，验证两个 change 目录内零命中。
- [x] 5.2 在仓内执行 `grep -rEn 'framework/(meta|cvm|ssm|vpc|cbs)/(resources|datasources|functions|ephemerals|lists|actions)/' -- 'openspec/changes/restructure-framework-types-and-naming/' 'openspec/changes/restructure-framework-single-level-layout/' 'tencentcloud/'`，验证零命中（双层路径在仓内任何位置都不应出现）。
- [x] 5.3 在仓内执行 `grep -rEn '_ephemeral_resource\.go|_list_resource\.go' -- 'openspec/changes/restructure-framework-types-and-naming/' 'openspec/changes/restructure-framework-single-level-layout/'`，验证两个 change 目录内零命中（仓代码侧本来就没有这些文件，仅校验文档措辞一致）。
- [x] 5.4 通读新 change 的 proposal.md / design.md / spec.md：确认"取代旧 change"语义在三份文档中至少各出现一次显式表述，且文件名 / 包名 / 路径示例之间互相一致（同一个例子在三份文档里必须给同样的字面值）。

## 6. OpenSpec 验证

- [x] 6.1 运行 `openspec validate restructure-framework-single-level-layout`，确认无 error。
- [x] 6.2 运行 `openspec validate restructure-framework-types-and-naming`，确认硬改写后仍无 error（spec / design / tasks / proposal 互相一致）。
- [x] 6.3 运行 `openspec status --change restructure-framework-single-level-layout`，确认 4 个 artifact 状态均为 done。

## 7. 后续 apply 阶段（不在本 change 范围、仅做记录）

- [ ] 7.1 （后续独立 apply）按新规范创建 `tencentcloud/framework/meta/` 目录与 5 个 reference 实现：`resource_tc_meta_local_note.go` / `data_source_tc_meta_provider_runtime.go` / `function_tc_meta_parse_resource_id.go` / `ephemeral_tc_meta_temp_credential.go` / `list_tc_meta_region.go`（L0 占位）。
- [ ] 7.2 （后续独立 apply）按新规范创建 `tencentcloud/framework/cvm/action_tc_cvm_reboot_instance.go`。
- [ ] 7.3 （后续独立 apply）按新规范创建 `tencentcloud/framework/ssm/ephemeral_tc_ssm_secret_version.go` 及其测试与文档（用户的真实业务诉求）。
- [ ] 7.4 （后续独立 apply）`tencentcloud/framework/registry.go` 按 `Resources / DataSources / Functions / EphemeralResources / ListResources / Actions` 6 方法直接 import `framework/meta` 与 `framework/cvm` / `framework/ssm` 的工厂，无 alias。
