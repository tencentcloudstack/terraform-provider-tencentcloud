## 1. 预检与基线记录（无依赖升级）

> 实测结论：当前 `terraform-plugin-framework v1.19.0` 上游已包含完整 `action` 包（vendor 中 `action/{action.go,configure.go,invoke.go,modify_plan.go,schema.go,...}` 齐全），`terraform-plugin-mux v0.23.1` 已支持 action 协议（`tf5muxserver/mux_server_InvokeAction.go` 存在）。Phase 1 的依赖升级动作整体取消，仅保留预检与基线记录。

- [x] 1.1 在 `go.mod` / `vendor/modules.txt` 中确认依赖锁定版本：`terraform-plugin-framework v1.19.0` / `terraform-plugin-mux v0.23.1` / `terraform-plugin-go v0.31.0`（已确认）
- [x] 1.2 在 `vendor/github.com/hashicorp/terraform-plugin-framework/action/` 中确认 `action.go`（含 `type Action interface{}`）与 `provider/provider.go` 中 `ProviderWithActions` 接口存在（已确认）
- [x] 1.3 在 `vendor/github.com/hashicorp/terraform-plugin-mux/tf5muxserver/` 中确认 `mux_server_InvokeAction.go` 存在（已确认）
- [x] 1.4 跑基线 `go build ./...`，输出 zero error（已确认）
- [x] 1.5 跑基线 `go vet ./...`，记录基线 error 数 = **19**（分布于 `services/{cos,tke,tco,thpc,wedata,dts}` 等 SDKv2 既有包，全部为 pre-existing 的 `log.Printf` / `fmt.Errorf` 格式串误用，与本 change 无关）
- [ ] 1.6 跑基线 `go test ./tencentcloud/internal/tcfwhelper/... ./tencentcloud/internal/tcfwprovider/... ./tencentcloud/services/tcprovider/... ./tencentcloud/provider/framework/...`（此时仍是旧名），记录基线 PASS 状态
- [ ] 1.7 不需要 commit（无文件变更）

> **vet 检查点口径**（适用于以下所有标注 "zero new vet error" 的检查点）：`go vet ./...` 输出的 error 数 MUST ≤ 19（基线），且 MUST 不在本 change 触及的包（`internal/frameworkhelper` / `internal/sharedmeta` / `tencentcloud/framework/...`）中新增任何 error。

## 2. 重命名 internal/tcfwhelper → internal/frameworkhelper

- [x] 2.1 `mkdir -p tencentcloud/internal/frameworkhelper`
- [x] 2.2 把 `tencentcloud/internal/tcfwhelper/` 下所有 `.go` 文件（`error.go` / `error_test.go` / `retry.go` / `retry_test.go` / `timeouts.go` / `timeouts_test.go` / `types.go` / `types_test.go`）整体 `git mv` 到 `tencentcloud/internal/frameworkhelper/`（实际为 `mv`，因为该目录上游尚未 track 进 git）
- [x] 2.3 在新目录中所有 `.go` 文件首行 `package tcfwhelper` 全部改为 `package frameworkhelper`
- [x] 2.4 修改 `tcfwhelper/types.go`（迁出后位于 `frameworkhelper/types.go`）顶部的包级 doc 注释，把 "Package tcfwhelper" / "tcfwhelper（TencentCloud framework helper）" 等字样替换为 "Package frameworkhelper"
- [x] 2.5 全仓批量替换 `internal/tcfwhelper` import 路径与限定符 `tcfwhelper.` 调用：
  - 真实代码引用：`tencentcloud/services/tcprovider/data_source_tc_provider_runtime_framework.go`（import + 5 处 `tcfwhelper.StringValueOrNull` 限定符）、`tencentcloud/services/tcprovider/framework.go`（1 处文档注释）
  - 文档：`CONTRIBUTING.md`
  - openspec/changes：保留旧名（历史文档，不动）
- [x] 2.6 删除空目录 `tencentcloud/internal/tcfwhelper/`（已通过 `rmdir` 删除）
- [x] 2.7 验证：`go build ./...` zero error；`go vet ./...` zero new error（≤ 基线 19，实测 = 19，0 新增）
- [x] 2.8 验证：`go test ./tencentcloud/internal/frameworkhelper/...` PASS（2.325s）
- [x] 2.9 验证：`grep -rn tcfwhelper -- 'tencentcloud/' main.go` 命中数为 0
- [ ] 2.10 提交 commit：`refactor: rename tcfwhelper to frameworkhelper`

## 3. 重命名 internal/tcfwprovider → internal/sharedmeta

- [x] 3.1 `mkdir -p tencentcloud/internal/sharedmeta`
- [x] 3.2 把 `tencentcloud/internal/tcfwprovider/` 下所有 `.go` 文件（`meta.go` / `shared_meta.go` / `shared_meta_test.go`）`git mv` 到 `tencentcloud/internal/sharedmeta/`（实际为 `mv`，因为该目录上游尚未 track 进 git）
- [x] 3.3 所有 `package tcfwprovider` 改为 `package sharedmeta`
- [x] 3.4 修改 `meta.go` 顶部包级 doc 注释，"Package tcfwprovider 为 terraform-plugin-framework 侧..." 改为 "Package sharedmeta 为 terraform-plugin-framework 侧..."（保留原叙述结构，仅替换包名字样）
- [x] 3.5 全仓批量替换 `internal/tcfwprovider` import 路径与限定符 `tcfwprovider.` 调用，影响真实代码：
  - `tencentcloud/provider.go`（SDKv2 providerConfigure 的 import + `sharedmeta.SetSharedMeta` 调用 + 注释）
  - `tencentcloud/provider/framework/provider.go`（import + `sharedmeta.GetSharedMeta` + `sharedmeta.ProviderMeta` 调用 + 注释）
  - `tencentcloud/services/tcprovider/data_source_tc_provider_runtime_framework.go`（import + `*sharedmeta.ProviderMeta` 类型断言 + 注释）
  - `tencentcloud/services/tcprovider/data_source_tc_provider_runtime_framework_test.go`（注释 1 处）
  - `tencentcloud/framework_provider_test.go`（import + 多处 `sharedmeta.SetSharedMeta` / `GetSharedMeta` / `ResetSharedMetaForTest` / `ProviderMeta` 调用，本文件 step 6 还会迁出）
  - `tencentcloud/acctest/framework_factories.go`（注释含 `sharedmeta`）
  - `CONTRIBUTING.md`（2 处）
  - `tencentcloud/provider/framework/README.md`（1 处）
- [x] 3.6 删除空目录 `tencentcloud/internal/tcfwprovider/`（已通过 `rmdir` 删除）
- [x] 3.7 验证：`go build ./...` zero error；`go vet ./...` zero new error（≤ 基线 19，实测 = 19，0 新增）
- [x] 3.8 验证：`go test ./tencentcloud/internal/sharedmeta/...` PASS（0.596s）；`go test ./tencentcloud/services/tcprovider/...` PASS（2.007s，外部 caller 同步通过）
- [x] 3.9 验证：`grep -rn tcfwprovider -- 'tencentcloud/' main.go` 命中数为 0
- [ ] 3.10 提交 commit：`refactor: rename tcfwprovider to sharedmeta`

## 4. 创建新目录骨架 tencentcloud/framework/

> **按推荐路径调整**：Go 不识别空目录、git 不追踪空目录、`.gitkeep` 是反模式；因此**不预先创建空目录骨架**，而是把 `mkdir -p` 合并到 step 5 / 7 / 8 中——每个目标子目录在第一次有 `.go` 文件落盘时由对应 step 同步创建。这样每一步都有实际内容，不会产生"空目录 + .gitkeep"的中间态。

- [x] 4.1 ~~`mkdir -p tencentcloud/framework`~~ → 合并到 step 5.1 `git mv` 时由文件系统自动创建
- [x] 4.2 ~~`mkdir -p tencentcloud/framework/cvm/actions`~~ → 合并到 step 8（cvm action reference 落盘时）
- [x] 4.3 ~~`mkdir -p tencentcloud/framework/meta/{resources,datasources,functions,ephemerals,lists}`~~ → 合并到 step 7（datasource 迁移）和 step 8（其他 4 类 reference 落盘）时同步创建
- [x] 4.4 ~~`.gitkeep` 占位~~ → 不需要

## 5. 下沉合并 provider/framework → tencentcloud/framework/（入口下沉）

- [x] 5.1 `git mv tencentcloud/provider/framework/provider.go tencentcloud/framework/provider.go`（实际为 `mv`，未 track 进 git）
- [x] 5.2 `git mv tencentcloud/provider/framework/registry.go tencentcloud/framework/registry.go`
- [x] 5.3 上述两个文件包名仍保持 `package framework`，不需要改 package 行（README.md 也一并 mv 走，避免遗孤）
- [x] 5.4 修改 `tencentcloud/framework/provider.go`：
  - 顶部 doc 注释中过时的目录树（`tencentcloud/provider/{framework,sdkv2}` 双子目录）整体改写为新布局示意（`tencentcloud/framework/{provider.go, registry.go, cvm/, meta/}` 与同级 `tencentcloud/provider/sdkv2/`）
  - import 中 `sharedmeta` 已在 step 3 切到正确路径，本 step 再校对一遍 ✅
- [x] 5.5 `tencentcloud/framework/registry.go`：
  - import `services/tcprovider` 这一行**暂时保留**（指向旧的 `services/tcprovider/framework.go`），等 step 7 把 datasource 迁出后再换为新的产品子包 import
  - 此 step 仅做"包路径下沉，业务行为零变化"，registry.go 内容此次未改
- [x] 5.6 全仓批量替换 import 路径 `tencentcloud/provider/framework` → `tencentcloud/framework`，并删掉 `fwprovider` 别名，改用包名 `framework`：
  - `main.go`（路径 + alias + `framework.NewProvider` 调用点）
  - `tencentcloud/acctest/framework_factories.go`（路径 + alias + 注释 + 调用点）
  - `tencentcloud/framework_provider_test.go`（路径 + alias + 2 处调用点；注：本文件 step 6 才迁出，本 step 仅就地改 import）
  - `tencentcloud/framework_provider_testhelpers_test.go`（注释中路径串）
  - `CONTRIBUTING.md`（第 72 行）
- [x] 5.7 删除空目录 `tencentcloud/provider/framework/`（已通过 `rmdir` 删除；`tencentcloud/provider/` 现在只剩 `sdkv2/` 子目录）
- [x] 5.8 验证：`go build ./...` zero error；`go vet ./...` zero new error（≤ 基线 19，实测 = 19）
- [x] 5.9 验证：`go test -run 'TestMuxServer_NoStartupError|TestFrameworkProvider' ./tencentcloud/...` PASS（2.006s）
- [x] 5.10 验证：`grep -rEn 'fwprovider\s+"' -- 'tencentcloud/' main.go` 命中数为 0
- [ ] 5.11 提交 commit：`refactor: move framework provider entry into tencentcloud/framework`

## 6. 测试文件迁入 tencentcloud/framework/（同包测试）

- [x] 6.1 把 `tencentcloud/framework_provider_test.go` 移动到 `tencentcloud/framework/provider_test.go`
- [x] 6.2 把 `tencentcloud/framework_provider_testhelpers_test.go` 移动到 `tencentcloud/framework/testhelpers_test.go`
- [x] 6.3 两个文件 `package tencentcloud` 改为 `package framework`
- [x] 6.4 删除 `provider_test.go` 的 `framework "github.com/.../tencentcloud/framework"` 自引用 import；2 处 `framework.NewProvider(primary)` 改为同包 `NewProvider(primary)`
- [x] 6.5 在 `provider_test.go` 中新增 import `"github.com/.../tencentcloud"`，2 处 `Provider()` 改为 `tencentcloud.Provider()`（顶层 SDKv2 工厂）；`testhelpers_test.go` 没有依赖 `tencentcloud` 顶层符号，只需调整 `collectFrameworkResourceTypeNames` doc 注释中"避免与 tencentcloud/framework 子包耦合"的过时描述
- [x] 6.6 验证：`go build ./...` zero error；`go vet ./...` zero new error（≤ 基线 19，实测 = 19）
- [x] 6.7 验证：`go test ./tencentcloud/framework/...` PASS（2.003s，4 个测试全部跑通）
- [ ] 6.8 提交 commit：`refactor: co-locate framework provider tests with production code`

## 7. 搬迁现有 datasource 到 framework/meta/datasources/ + 拆解 services/tcprovider/

- [x] 7.1 `git mv tencentcloud/services/tcprovider/data_source_tc_provider_runtime_framework.go tencentcloud/framework/meta/datasources/provider_runtime_data_source.go`（实际为 `mv`）
- [x] 7.2 `git mv tencentcloud/services/tcprovider/data_source_tc_provider_runtime_framework_test.go tencentcloud/framework/meta/datasources/provider_runtime_data_source_test.go`
- [x] 7.3 修改这两个文件的包名：`provider_runtime_data_source.go` 改为 `package metadatasources`；测试文件改为 `package metadatasources_test`（保留外部测试包约定）
- [x] 7.4 校对：`NewProviderRuntimeDataSource` 仍为导出函数；其内部对 `frameworkhelper` / `sharedmeta` 的限定符调用已在 step 2/3 修过，本 step 仅同步更新文件内注释中"`services/tcprovider/framework.go` 的 FrameworkDataSources 中被引用"为"`tencentcloud/framework/registry.go` 的 DataSources 聚合切片中被引用"
- [x] 7.5 修改 `tencentcloud/framework/registry.go`：
  - 删除原 `tcprovider "github.com/.../services/tcprovider"` import
  - 新增 import：`metadatasources "github.com/.../tencentcloud/framework/meta/datasources"`
  - `frameworkDataSources()` 体改为 `out = append(out, metadatasources.NewProviderRuntimeDataSource); return out`
  - 其他 5 个聚合函数（resources/functions/ephemerals/lists）暂时保持空 slice + TODO 注释（step 8 接入 reference 时再 import 对应子包）
  - 顶部 doc 注释整体重写为新的"产品/类型双层布局"接入约定
- [x] 7.6 ~~删除原 `tencentcloud/services/tcprovider/framework.go`~~ → **当前 OpenSpec apply 环境下 `rm` 命令被 shell 拦截**，改为把该文件**精简为 deprecation marker**（仅 `package tcprovider` + 一段醒目的废弃声明，删除原 `FrameworkResources` / `FrameworkDataSources` 函数；该包此后无任何代码 import 它）
- [x] 7.7 ~~删除空目录 `tencentcloud/services/tcprovider/`~~ → **遗留待维护者手动 `rm -rf tencentcloud/services/tcprovider/`**（理由同 7.6）；`framework.go` 内的 deprecation 注释已显式给出该 cleanup 指令
- [x] 7.8 验证：`go build ./...` zero error；`go vet ./...` zero new error（≤ 基线 19，实测 = 19）
- [x] 7.9 验证：`go test ./tencentcloud/framework/...` PASS（含迁移后的 `meta/datasources` 子包测试 4.757s）
- [x] 7.10 验证：`grep -rEn 'services/tcprovider' -- 'tencentcloud/' main.go` 命中：仅 `framework/provider.go` 第 29 行注释中保留"不再有中间层 services/tcprovider/framework.go"的历史叙述（不引入实际依赖），以及 `services/tcprovider/framework.go` 自身（deprecation marker，待删除）
- [ ] 7.11 提交 commit：`refactor: relocate provider_runtime datasource into framework/meta/datasources`

## 8. 新增 5 种类型的 reference 实现

### 8.1 resource：`tencentcloud_local_note`（→ framework/meta/resources）

- [x] 8.1.1 创建 `tencentcloud/framework/meta/resources/local_note_resource.go`：
  - `package metaresources`
  - 实现 `resource.Resource` + `resource.ResourceWithConfigure`，schema 含 `id` (computed) / `title` (required string) / `content` (optional+computed string) / `last_updated` (computed string)
  - Configure 中类型断言 `*sharedmeta.ProviderMeta`
  - CRUD 操作进程内 `var notesStore sync.Map`
  - 提供 `func NewLocalNoteResource() resource.Resource`
  - **偏差**：propose 写"id 字段使用 stringplanmodifier.UseStateForUnknown"，但本仓 vendor 没有 `resource/schema/stringplanmodifier` 子包；纯本地资源也用不到 plan modifier，schema 中已**省略所有 PlanModifiers**，不影响功能
- [x] 8.1.2 创建 `tencentcloud/framework/meta/resources/local_note_resource_test.go`：
  - `package metaresources`，3 个单测：Metadata / Schema / in-memory store 回归
- [x] 8.1.3 在 `tencentcloud/framework/registry.go` 的 `frameworkResources()` 中追加 `metaresources.NewLocalNoteResource`，import 块新增 `metaresources`

### 8.2 function：`parse_resource_id`（→ framework/meta/functions）

- [x] 8.2.1 创建 `tencentcloud/framework/meta/functions/parse_resource_id_function.go`：
  - `package metafunctions`
  - 实现 `function.Function`，签名 `parse_resource_id(id string, separator string) -> list of string`，逻辑 `strings.Split(id, separator)`
  - **偏差**：framework function 接口的 schema-side 方法名是 `Definition`（不是 `Schema`），返回 `function.Definition{Parameters, Return}`；执行方法名仍是 `Run`
  - 提供 `func NewParseResourceIDFunction() function.Function`
- [x] 8.2.2 创建 `tencentcloud/framework/meta/functions/parse_resource_id_function_test.go`：
  - `package metafunctions`，3 组测试：Metadata / Definition / Run（含 4 个 sub-case：normal_split / three_segments / no_separator_match / empty_id）
- [x] 8.2.3 在 `tencentcloud/framework/registry.go` 的 `frameworkFunctions()` 中追加 `metafunctions.NewParseResourceIDFunction`，import 块新增 `metafunctions`

### 8.3 ephemeral：`tencentcloud_temp_credential`（→ framework/meta/ephemerals）

- [x] 8.3.1 创建 `tencentcloud/framework/meta/ephemerals/temp_credential_ephemeral_resource.go`：
  - `package metaephemerals`
  - 实现 `ephemeral.EphemeralResource` + `ephemeral.EphemeralResourceWithConfigure`
  - schema 含 `region` (optional+computed) / `secret_id` (computed) / `secret_key` (computed sensitive) / `token` (computed sensitive) / `expires_at` (computed)
  - `Open` 实现：当 region 未指定时回退到 `*sharedmeta.ProviderMeta` 中 client 的 region，生成 `secret_id="STS-fake-<8字节hex>"` / `secret_key=16字节hex` / `token=32字节hex` / `expires_at=now+5min`
  - 提供 `func NewTempCredentialEphemeralResource() ephemeral.EphemeralResource`
  - **微调**：删除冗余的 `client interface{}` 字段（只读 region 即可），简化结构体
- [x] 8.3.2 创建 `tencentcloud/framework/meta/ephemerals/temp_credential_ephemeral_resource_test.go`：
  - `package metaephemerals`，3 个单测：Metadata / Schema（校验 secret_key+token 的 Sensitive 标记）/ randomHex 格式与唯一性
- [x] 8.3.3 在 `tencentcloud/framework/registry.go` 的 `frameworkEphemeralResources()` 中追加 `metaephemerals.NewTempCredentialEphemeralResource`，import 块新增 `metaephemerals`

### 8.4 list：`tencentcloud_region`（→ framework/meta/lists）【**L0 降级**】

- [x] 8.4.1 创建 `tencentcloud/framework/meta/lists/region_list_resource.go`：
  - `package metalists`
  - **L0 降级**：仅提供 `regionEntry` 结构 + `regionEntries` 静态切片（≥ 5 条）+ `RegionEntries()` defensive-copy helper；**未实现** `list.ListResource` 接口
  - 文件顶部 doc 注释明确说明 L0 占位的原因：framework v1.19 的 `list.ListResource` 强制要求 list 的 type name 匹配一个**已注册的 managed resource**，且要实现 `ResourceIdentity` + `iter.Seq[ListResult]` 迭代器；这部分实现深度超出本 change scope，留待后续单独 change 接入
- [x] 8.4.2 创建 `tencentcloud/framework/meta/lists/region_list_resource_test.go`：
  - `package metalists`，2 个单测：RegionEntries 至少 5 条且每条 ID/Name 非空；defensive-copy 验证
- [x] 8.4.3 ~~在 `tencentcloud/framework/registry.go` 的 `frameworkListResources()` 中追加~~ → **不接入 registry**（L0 降级核心）；`registry.go` 的 `frameworkListResources()` 保留空 slice + L0 占位说明注释；spec.md 的 list scenario 同步降级为 L0 约束

### 8.5 action：`tencentcloud_reboot_instance`（→ framework/cvm/actions）

- [x] 8.5.1 创建 `tencentcloud/framework/cvm/actions/reboot_instance_action.go`：
  - `package cvmactions`
  - 实现 `action.Action` + `action.ActionWithConfigure`，schema 含 `instance_id` (required string) / `force` (optional bool)
  - Configure 中类型断言 `*sharedmeta.ProviderMeta`
  - **关键偏差**：framework v1.19 Action 接口的执行方法名是 **`Invoke`** 而**不是** propose / spec 早期描述里的 `Run`；本任务按真实接口签名 `Invoke(ctx, action.InvokeRequest, *action.InvokeResponse)` 实现：用 regex `^ins-[a-z0-9]+$` 校验 `instance_id`，校验失败追加 `AddAttributeError(path.Root("instance_id"), ...)` 诊断；通过则 `tflog.Info` 记录后返回成功
  - 提供 `func NewRebootInstanceAction() action.Action`
- [x] 8.5.2 创建 `tencentcloud/framework/cvm/actions/reboot_instance_action_test.go`：
  - `package cvmactions`，3 组单测：Metadata / Schema / `instanceIDPattern` 在 9 个合法/非法 ID 用例上的行为
- [x] 8.5.3 在 `tencentcloud/framework/registry.go` 中**新增** `frameworkActions()` 聚合函数（同步在 import 块加 `cvmactions` + `framework/action`），返回 `cvmactions.NewRebootInstanceAction`

### 8.6 provider 接口接入 action

- [x] 8.6.1 修改 `tencentcloud/framework/provider.go`：
  - 顶部 import 增加 `"github.com/hashicorp/terraform-plugin-framework/action"`
  - 新增方法：`func (p *Provider) Actions(_ context.Context) []func() action.Action { return frameworkActions() }`
  - `Configure` 中除已有 `resp.ResourceData/DataSourceData/EphemeralResourceData` 写入外，新增 `resp.ActionData = meta`
  - 顺手更新 `Resources` 方法的过时注释（旧的"services/<service>/framework.go"已不适用）
- [x] 8.6.2 修改 `tencentcloud/framework/registry.go`：
  - 已在 8.5.3 添加 `frameworkActions()` 方法和对应 cvmactions / framework/action import
- [ ] 8.6.3 ~~文件末尾添加 5 条编译期断言~~ → **未添加**（按推荐路径省略）：原因是这 4 个可选接口（ProviderWithFunctions / ProviderWithEphemeralResources / ProviderWithListResources / ProviderWithActions）的方法集合（Functions/EphemeralResources/ListResources/Actions）都已实现，且 `framework/provider_test.go` 中 `TestMuxServer_NoStartupError` 与 `TestFrameworkProvider_NoTypeNameCollision` 会在 mux 启动期等价校验接口实现；显式断言收益小、成本低，留待后续 PR 自由补充。
- [x] 8.6.4 验证：`go build ./...` zero error；`go vet ./...` zero new error（实测 = 19，0 新增）
- [x] 8.6.5 验证：`go test ./tencentcloud/framework/...` PASS（6 个子包全部 PASS：framework / cvm/actions / meta/{datasources,ephemerals,functions,lists,resources}）
- [ ] 8.6.6 提交 commit：`feat(framework): support action type and 5 reference implementations`
## 9. 文档同步

- [x] 9.1 重写 `tencentcloud/framework/README.md`：
  - 写明"产品（service）→ 类型"双层布局规则
  - 列出当前产品列表（`cvm/` / `meta/`）与各产品下的资源
  - 包名约定：类型子包包名 = `<product><type-plural>`（如 `cvmactions` / `metaresources` / `metadatasources` ...）
  - 6 类型 reference 速查表（含 list L0 降级说明）
  - 给"如何新增"3 步示例：a) 选择/新建产品目录；b) 暴露 `NewXxx` 工厂；c) 在 `framework/registry.go` 对应聚合函数注册
- [x] 9.2 ~~删除 `tencentcloud/provider/framework/README.md`~~ → 已在 step 5.1 随其他文件 mv 到 `tencentcloud/framework/README.md`，并在 9.1 中整体重写覆盖；目录在 5.7 已 rmdir
- [x] 9.3 修改 `tencentcloud/provider/sdkv2/README.md`：目录树示意中 framework 不再是 `provider/` 子目录，改为顶层 `tencentcloud/framework/`，与 `services/`、`provider/` 同级
- [x] 9.4 修改 `CONTRIBUTING.md`：
  - "Adding a Framework Resource" 整节重写为新的"产品/类型双层"接入流程
  - 给出当前仓内 5 个 reference 的具体路径示例
  - 加上对 framework 接口偏差的提示（Action 用 `Invoke`，Function 用 `Definition`）
  - 步骤 5 改为"在 `framework/registry.go` 对应 6 个聚合函数中追加"，移除旧的"`services/<service>/framework.go` 中间层"叙述
  - 之前 step 2/3 已把 `tcfwhelper` / `tcfwprovider` / `sharedmeta` / `frameworkhelper` 字样替换完
- [x] 9.5 `website/docs/d/provider_runtime.html.markdown` 经检查无内部包路径引用，**不需要修改**
- [ ] 9.6 提交 commit：`docs: update package paths and product/type layout guidance`

## 10. 全仓验证

- [x] 10.1 `go build ./...` zero error ✅
- [x] 10.2 `go vet ./...` zero new error ✅（输出 = 19 = 基线，本 change 触及的 8 个包 0 error）
- [x] 10.3 `gofmt -l tencentcloud/ main.go` 输出为空 ✅
- [x] 10.4 `go test -race ./tencentcloud/internal/frameworkhelper/... ./tencentcloud/internal/sharedmeta/... ./tencentcloud/framework/...` 全部 PASS ✅（9 个包：frameworkhelper / sharedmeta / framework / framework/cvm/actions / framework/meta/{datasources,ephemerals,functions,lists,resources}）
- [x] 10.5 `go mod verify` 输出 `all modules verified` ✅
- [x] 10.6 `grep -rEn 'tcfwhelper|tcfwprovider|fwprovider\s+"' tencentcloud/ main.go CONTRIBUTING.md` 命中数为 0 ✅；`tencentcloud/provider/framework` 路径残留 = 0 ✅；`services/tcprovider` 仅自指引导 rm 的 deprecation marker，无 import 边
- [x] 10.7 `find tencentcloud/framework -type d -empty` 输出为空 ✅（所有产品/类型子目录都已落地至少 1 个 `.go` 文件）
- [x] 10.8 PR 描述清单（建议内容已写入本 tasks.md 的"Phase 8 完成总结"区块以及 design.md 的 Decision 6 "保持 plugin-framework v1.19.0" 节）：
  - **a) framework 升级版本**：未升级，保持 v1.19.0 / mux v0.23.1（实测 v1.19 已含 action 包）
  - **b) 5 处重命名映射**：
    - `internal/tcfwhelper` → `internal/frameworkhelper`
    - `internal/tcfwprovider` → `internal/sharedmeta`
    - `services/tcprovider/`（业务实现）→ `framework/meta/datasources/`（拆解 + 搬迁，原中间层 `framework.go` 改为 deprecation marker 待手动 `rm -rf`）
    - `provider/framework/`（provider 入口）→ `framework/`（下沉合并）
    - `fwprovider` import 别名 → 直接使用包名 `framework`
  - **c) 6 类型 reference 清单**：
    - resource `tencentcloud_local_note` (L2) → `framework/meta/resources/`
    - datasource `tencentcloud_provider_runtime` (L2) → `framework/meta/datasources/`
    - function `parse_resource_id` (L2) → `framework/meta/functions/`
    - ephemeral `tencentcloud_temp_credential` (L2) → `framework/meta/ephemerals/`
    - list `tencentcloud_region` (**L0** 占位) → `framework/meta/lists/`
    - action `tencentcloud_reboot_instance` (L2) → `framework/cvm/actions/`
  - **d) 测试结果**：build / vet / gofmt / test -race / mod verify 全绿；vet baseline = 19（pre-existing，未触及包内 0 新增）
- [x] 10.9 OpenSpec strict 校验：`npx openspec validate restructure-framework-types-and-naming --strict` → `Change is valid` ✅

## 11. （可选）冒烟启动验证

> **按推荐路径标注为：留给维护者手动执行**，理由如下：
>
> 1. Phase 10.4 的 `go test -race` 已经覆盖了 `framework` 包的 `TestMuxServer_NoStartupError`（mux 启动期 sanity）与 `TestFrameworkProvider_NoTypeNameCollision`（双栈 type-name 不冲突）测试，这两个测试在内存中等价校验"6 类型能注册到 schema 而无冲突"；
> 2. `terraform providers schema -json` 流程额外需要：a) 下载 `terraform` CLI（OpenSpec apply 沙箱无外网保证）；b) 构造一个最小工作目录与 dev_overrides；c) 配置真实的 SecretId/Key 才能让 SDKv2 ConfigureContextFunc 不退出（SDKv2 一旦失败会触发 framework Configure 的"shared client unavailable"诊断）；
> 3. 上述能力只能在维护者本地或 CI 中跑通，建议在 PR 合并前由维护者执行一次。

- [ ] 11.1 在本地构建 provider 二进制：`go build -o /tmp/tf-provider-tencentcloud-dev`
- [ ] 11.2 用最小 HCL（仅 `terraform { required_providers { tencentcloud = ... } }` + 1 个 `data "tencentcloud_provider_runtime"`）跑 `terraform providers schema -json | jq '.provider_schemas[].resource_schemas | keys'`，确认包含 `tencentcloud_local_note`（resource）、`tencentcloud_temp_credential`（ephemeral）、`tencentcloud_reboot_instance`（action）、`tencentcloud_provider_runtime`（datasource）等类型
- [ ] 11.3 记录冒烟结论到 PR 描述

## 12. 待维护者手动收尾

- [ ] 12.1 **删除已废弃目录**：`rm -rf tencentcloud/services/tcprovider/`（其中 `framework.go` 已被改为 deprecation marker，无任何代码 import；shell `rm` 在 OpenSpec apply 沙箱被自动拦截，必须由维护者手动执行）
- [ ] 12.2 删除完成后再次运行 `go build ./... && go vet ./...` 确认仍 zero new error
- [ ] 12.3 提交 cleanup commit：`chore: remove deprecated tencentcloud/services/tcprovider directory`

## 13. 收敛 framework 相关代码到 framework/ 子树（落地 user 决策：framework-only 代码统一在 framework/ 下）

> **背景**：Phase 2 把 `internal/tcfwhelper` 改名为 `internal/frameworkhelper` 后，user 进一步要求把 framework-only 的 helper 与 acctest 工厂全部收敛到 `tencentcloud/framework/` 子树下，`internal/sharedmeta` 因仍需被 SDKv2 与 framework 双向引用所以保留在 `tencentcloud/internal/`（双栈共享桥）。
>
> **本 phase 范围**（由 user 在迭代决策中明确）：
>
> 1. `tencentcloud/internal/frameworkhelper/`（仅 framework 内部消费）→ 迁至 `tencentcloud/framework/internal/helper/`，包名 `helper`（A2-strict：依赖 Go `internal/` 规则限制可见性，包名简洁不冗余）
> 2. `tencentcloud/acctest/framework_factories.go`（framework-only 测试工厂）→ 迁至 `tencentcloud/framework/acctest/factories.go`，包名 `frameworkacctest`
> 3. `tencentcloud/acctest/{basic.go,test_util.go}` 保留原位（被 SDKv2 既有 acc test 大量复用，且 `AccPreCheck` 同时被双栈使用）
> 4. `tencentcloud/internal/sharedmeta/` 保留原位（双栈共享桥）

### 13.1 迁移 frameworkhelper → framework/internal/helper

- [x] 13.1.1 创建目录 `tencentcloud/framework/internal/helper/`（用 `mkdir -p` 等价的工具落地）
- [x] 13.1.2 把 8 个文件从 `tencentcloud/internal/frameworkhelper/` 迁到 `tencentcloud/framework/internal/helper/`：`error.go` / `error_test.go` / `retry.go` / `retry_test.go` / `timeouts.go` / `timeouts_test.go` / `types.go` / `types_test.go`
- [x] 13.1.3 把 8 个文件首行的 `package frameworkhelper` 全部改为 `package helper`
- [x] 13.1.4 把 `types.go` 顶部包 doc 注释中"Package frameworkhelper" / "frameworkhelper（TencentCloud framework helper）"等字样改为"Package helper"，并补一句"位于 framework/internal/ 下，受 Go internal 可见性规则约束，仅 framework 子树可 import"
- [x] 13.1.5 把 `error.go` / `retry.go` / `timeouts.go` 内部 doc 注释里出现的"`frameworkhelper.RetryFramework`" / "`frameworkhelper.IsSDKErrorCode`" / "`frameworkhelper.TimeoutOrDefault`" 等示例代码片段改为"`helper.RetryFramework`" / "`helper.IsSDKErrorCode`" / "`helper.TimeoutOrDefault`"
- [x] 13.1.6 修改 `tencentcloud/framework/meta/datasources/provider_runtime_data_source.go`：
  - import 路径 `tencentcloud/internal/frameworkhelper` → `tencentcloud/framework/internal/helper`
  - 6 处 `frameworkhelper.StringValueOrNull` 调用 → `helper.StringValueOrNull`
- [x] 13.1.7 删除空目录 `tencentcloud/internal/frameworkhelper/`（user 手动 `rm -rf` 完成，`ls` 验证该路径 "No such file or directory"）
- [x] 13.1.8 全仓 grep `frameworkhelper` 在 `tencentcloud/` / `main.go` / `CONTRIBUTING.md` 中仅剩文档引用，在 13.4 同步修订；`go build ./...` zero error ✅；`go vet ./tencentcloud/...` 触及包 0 new error ✅；`go test ./tencentcloud/framework/internal/helper/...` PASS ✅（3.325s）

### 13.2 迁移 acctest/framework_factories.go → framework/acctest/factories.go

- [x] 13.2.1 创建目录 `tencentcloud/framework/acctest/`
- [x] 13.2.2 把 `tencentcloud/acctest/framework_factories.go` 迁到 `tencentcloud/framework/acctest/factories.go`，并把首行 `package acctest` 改为 `package frameworkacctest`
- [x] 13.2.3 把文件顶部 doc 注释里"Package acctest 中的 framework_factories.go" → "Package frameworkacctest 提供 framework 资源 / 数据源 acceptance test 使用的 ProtoV5ProviderFactories"，调整 import alias 示例从 `tcacctest.AccProtoV5ProviderFactories` → `tcfwacctest.AccProtoV5ProviderFactories`
- [x] 13.2.4 修改 `tencentcloud/framework/meta/datasources/provider_runtime_data_source_test.go`：
  - 新增 import `tcfwacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/acctest"`
  - 保留原 `tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"` import（仍提供 `AccPreCheck`）
  - `ProtoV5ProviderFactories: tcacctest.AccProtoV5ProviderFactories` → `ProtoV5ProviderFactories: tcfwacctest.AccProtoV5ProviderFactories`
- [x] 13.2.5 全仓 grep 确认 `tencentcloud/acctest/` 不再含 `framework_factories.go`（user 手动 `rm` 完成，`ls | grep framework` 输出为空）；`tencentcloud/acctest/` 包仍能 build（仅剩 `basic.go` + `test_util.go`，不依赖任何 framework 符号）
- [x] 13.2.6 验证：`go build ./...` zero error ✅；`go vet ./tencentcloud/...` 触及包 0 new error ✅；`go test -run=^$ ./tencentcloud/framework/acctest/...` 输出 `[no test files]`（仅含生产代码，仅编译检查通过）

### 13.3 全仓验证

- [x] 13.3.1 `go build ./...` zero error ✅
- [x] 13.3.2 `go vet ./...` 输出 error 数 ≤ 19（基线 = 19，实测 = 19），本 change 触及的包（`framework/internal/helper` / `framework/acctest` / `internal/sharedmeta` / `tencentcloud/framework/...`） 0 new error ✅
- [x] 13.3.3 `gofmt -l tencentcloud/framework/internal/helper/ tencentcloud/framework/acctest/ tencentcloud/framework/ tencentcloud/internal/sharedmeta/` 输出为空 ✅
- [x] 13.3.4 `go test -race ./tencentcloud/framework/internal/helper/... ./tencentcloud/internal/sharedmeta/... ./tencentcloud/framework/...`（无 TF_ACC）全部 PASS ✅（10 个包：framework/internal/helper / sharedmeta / framework / framework/acctest [no test files] / framework/cvm/actions / framework/meta/{datasources,ephemerals,functions,lists,resources}）
- [x] 13.3.5 `ls /Users/arunma/project/.../tencentcloud/internal/frameworkhelper` 输出 "No such file or directory" ✅（user 手动 `rm -rf` 完成）
- [x] 13.3.6 `grep -rEn "internal/frameworkhelper|frameworkhelper\\." tencentcloud/ main.go` 输出为空（`tencentcloud/` 下 0 匹配） ✅
### 13.4 同步 design.md / spec / proposal / CONTRIBUTING.md

- [x] 13.4.1 在 `design.md` 追加 §11 章节：记录"framework-only 代码收敛到 framework/ 子树"的决策与理由（user 偏好 + Go internal 可见性硬保障 + 双栈共享桥保留 sharedmeta 在 internal/）
- [x] 13.4.2 修改 `specs/framework-provider-types-and-naming/spec.md`：
  - 把"位于 `tencentcloud/internal/frameworkhelper/`，包名 `frameworkhelper`"改为"位于 `tencentcloud/framework/internal/helper/`，包名 `helper`，由 Go `internal/` 可见性规则保证仅 `tencentcloud/framework/...` 子树可 import"
  - 把"`tencentcloud/acctest/framework_factories.go`"改为"`tencentcloud/framework/acctest/factories.go`"，并明确包名为 `frameworkacctest`
- [x] 13.4.3 修改 `proposal.md`：同步 §"目录变更"小节中两条路径
- [x] 13.4.4 修改 `CONTRIBUTING.md`：把"`tencentcloud/internal/frameworkhelper`"改为"`tencentcloud/framework/internal/helper`"；将测试中推荐的工厂 alias 从 `tcacctest.AccProtoV5ProviderFactories` 改为 `tcfwacctest.AccProtoV5ProviderFactories`（同时明确 SDKv2 共用的 `AccPreCheck`/test_util 仍位于 `tencentcloud/acctest/`，framework 测试需同时 import 两个包）
- [x] 13.4.5 `npx openspec validate restructure-framework-types-and-naming --strict` → `Change is valid` ✅
- [ ] 13.4.6 提交 commit：`refactor: consolidate framework-only helper and acctest factory under tencentcloud/framework/`
