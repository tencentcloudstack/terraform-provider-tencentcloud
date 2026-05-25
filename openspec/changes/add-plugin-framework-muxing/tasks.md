## 1. 依赖与 vendor 准备

- [x] 1.1 核对 `go.mod` 已有 `terraform-plugin-framework v1.19.0` / `terraform-plugin-go v0.31.0` / `terraform-plugin-mux v0.23.1`；本次不引入 `framework-validators` 与 `framework-timeouts`（timeouts 自行用原生 framework schema 实现，validators 留待首个真实使用的资源时再引入），避免触发大范围 vendor diff
- [x] 1.2 因未新增依赖，无需 `go mod vendor`；本步骤记为已完成
- [ ] 1.3 验证 `go build ./...` 与 `go vet ./...` 全部通过（在第 9 章统一执行）

## 2. 共享 ProviderData 与 client 桥接

- [x] 2.1 新建 `tencentcloud/internal/tcfwprovider/meta.go`，定义 `type ProviderMeta struct { Client *connectivity.TencentCloudClient }`（命名贴近 SDKv2 这边"meta"的说法，与 `tccommon`、`apiV3Conn` 等现有概念保持一致）
- [x] 2.2 在 `tencentcloud/internal/tcfwprovider/shared_meta.go` 中提供进程级 `atomic.Pointer[connectivity.TencentCloudClient]`，导出 `SetSharedMeta`/`GetSharedMeta`/`ResetSharedMetaForTest`
- [x] 2.3 修改 SDKv2 `providerConfigure`（位于 `tencentcloud/provider.go` line 3066），在 `return &tcClient, nil` 之前调用 `tcfwprovider.SetSharedMeta(tcClient.apiV3Conn)`
- [x] 2.4 完成 `tcfwprovider/shared_meta_test.go`：`go test -race ./tencentcloud/internal/tcfwprovider/...` 全部通过

## 3. FrameworkProvider runtime 实现

- [x] 3.1 在 `tencentcloud/framework_provider.go` 的 `Schema` 中镜像 SDKv2 全部字段：顶层 attribute（`secret_id`/`secret_key`/`security_token`/`region`/`protocol`/`domain`/`cos_domain`/`enable_pod_oidc`/`shared_credentials_dir`/`profile`/`cam_role_name`/`allowed_account_ids`/`forbidden_account_ids`），以及 block（`assume_role`/`assume_role_with_saml`/`assume_role_with_web_identity`/`mfa_certification`）；敏感字段标 `Sensitive`
- [x] 3.2 `Configure` 从 `tcfwprovider.GetSharedMeta()` 读取 client；nil 追加 Error 诊断；非 nil 时写入 `resp.ResourceData/DataSourceData/EphemeralResourceData`
- [x] 3.3 抽出 `tencentcloud/framework_registry.go` 提供 5 个聚合函数（`frameworkResources`/`frameworkDataSources`/`frameworkFunctions`/`frameworkEphemeralResources`/`frameworkListResources`），本变更阶段返回空切片，后续按服务追加
- [x] 3.4 `framework_provider_test.go` + `framework_provider_testhelpers_test.go`：mux 启动验证、两栈资源/数据源名交集校验、shared meta 准备类试验，均通过

## 4. 公共工具包 tcfwhelper

- [x] 4.1 `tencentcloud/internal/tcfwhelper/types.go`：`StringValueOrNull`/`StringPointerValueOrNull`/`Int64ValueOrNull`/`Int64ValueFromUint`/`BoolValueOrNull` + 反向的 `*PointerFromValue`，覆盖 nil/Null/Unknown 三种边界
- [x] 4.2 `retry.go`：`RetryFramework[T any](ctx, timeout, retryable, fn) (T, error)` 复用 `helper/resource.RetryContext` 语义；提供 `IsTimeoutError`
- [x] 4.3 `error.go`：`WrapSDKError(summary, err) diag.Diagnostic` 从 `*sdkErrors.TencentCloudSDKError` 提取 Code/Message/RequestId；提供 `IsSDKErrorCode(err, codes...)` 供 retryable 判定使用
- [x] 4.4 `timeouts.go`：用原生 framework `SingleNestedBlock` 自实现 timeouts 块 + `TimeoutsModel`（`tfsdk:"create/read/update/delete"`） + `ParseTimeout`/`TimeoutOrDefault` 解析帮手
- [x] 4.5 `types_test.go` / `retry_test.go` / `error_test.go` / `timeouts_test.go` 表驱动单测全部通过（`go test -race`）

## 5. mux 入口与启动校验

- [x] 5.1 确认 `main.go` mux 顺序：SDKv2 在前、framework 在后，本变更不动
- [x] 5.2 `GNUmakefile` 新增 `check-mux` target：fmtcheck 后调用 `go test -count=1 -run 'TestMuxServer_NoStartupError|TestFrameworkProvider_NoTypeNameCollision' ./tencentcloud/`，本地 `make check-mux` 通过
- [x] 5.3 两栈资源/数据源名交集校验已集成在 `TestFrameworkProvider_NoTypeNameCollision`（与 3.4 共用）
- [x] 5.4 `.github/workflows/unit-tests.yml` 在 `unit` 步骤后追加 `mux compatibility` 步骤 `run: make check-mux`

## 6. 示范 framework 资源（PoC）

- [x] 6.1 选定“零存量、零云 API 调用”的诊断型数据源 `tencentcloud_provider_runtime`（读取当前 region/client_version/stack_mode，用于调试双栈状态）作为首个 framework PoC。选型理由：与 SDKv2 零冲突、不需要真账号即可跑 acceptance test、完整覆盖 framework 数据源全生命周期
- [x] 6.2 新建 `tencentcloud/services/tcprovider/framework.go`，导出 `FrameworkResources()` / `FrameworkDataSources()`
- [x] 6.3 新建 `tencentcloud/services/tcprovider/data_source_tc_provider_runtime_framework.go`：实现 `Metadata`/`Schema`/`Configure`/`Read`，使用 `tcfwhelper.StringValueOrNull`、`tcfwprovider.ProviderMeta` 类型断言
- [x] 6.4 `tencentcloud/framework_registry.go` 中 `out = append(out, tcprovider.FrameworkResources()...)` / `tcprovider.FrameworkDataSources()`亮出接入
- [x] 6.5 `tencentcloud/acctest/framework_factories.go` 提供 `AccProtoV5ProviderFactories`；示范数据源的 `data_source_tc_provider_runtime_framework_test.go` 使用该工厂
- [x] 6.6 手写 `website/docs/d/provider_runtime.html.markdown`（`make doc` 现有 gendoc 不识 framework schema，本变更不改造生成器，由 task 7 调研担当）

## 7. 文档生成器适配

- [x] 7.1 调研结果：`gendoc/main.go` 是项目自研生成器，强依赖 `terraform-plugin-sdk/v2/helper/schema`，仅遍历 SDKv2 `ResourcesMap` / `DataSourcesMap`，无法识别 framework schema。详见 `openspec/changes/add-plugin-framework-muxing/docgen-research.md`
- [x] 7.2 决议：本变更采用混合方案 C（不动 gendoc，框架资源手写文档）；待 framework 资源 ≥ 5 个或出现复杂嵌套时，独立 change 升级 `tfplugindocs`。不在本变更范围内改造 `make doc`。
- [x] 7.3 手写 `website/docs/d/provider_runtime.html.markdown` 覆盖示范数据源的字段、示例与 Attributes Reference；结论记录在 `docgen-research.md`

## 8. CI / 流程文档

- [x] 8.1 项目原本无 `CONTRIBUTING.md`，新增一份：涵盖双栈选型规则、新资源开发步骤、硬规则（两栈不同名、凭证只有一份、vendor 同步）、`make` 目标说明、PR checklist
- [x] 8.2 项目原本无 PR 模板，新增 `.github/PULL_REQUEST_TEMPLATE.md`，含勾选项“确认资源类型名未在另一栈重复注册”
- [x] 8.3 新增 `.changelog/add-plugin-framework-muxing.txt`（合并时交维护者重命名为 PR 号），包含 `release-note:enhancement` 与 `release-note:new-data-source`
- [x] 8.4 README 在 "Requirements" 之前插入 "Architecture" 段落，描述双栈与 mux，并指向 CONTRIBUTING.md
## 9. 验证

- [x] 9.1 `go build ./...` 通过（本地 BUILD_EXIT=0）
- [x] 9.2 `make fmt` 本地通过；gofmt 无 diff。`make lint` 本地环境缺 `tfproviderlint`，且本地 golangci-lint 是用 go1.24 构建（项目为 1.25.8）无法跑；交由 CI 运行
- [x] 9.3 `go test` 覆盖本变更所有新增包：
      - `./tencentcloud/internal/tcfwprovider/...` PASS
      - `./tencentcloud/internal/tcfwhelper/...` PASS
      - `./tencentcloud/services/tcprovider/...` PASS
      - `./tencentcloud/` 中 `TestProvider|TestMuxServer|TestFrameworkProvider` PASS
- [x] 9.4 `make check-mux` 通过
- [x] 9.5 `make doc` 不在本变更范围内（`gendoc` 仅识 SDKv2 schema，示范数据源采用手写文档；详见 docgen-research.md）
- [ ] 9.6 真账号 acceptance test：示范资源不调用云 API，仅需 `TF_ACC=1 SECRET_ID/KEY` 即可跑；不产生费用。**请你在本地或 CI 有凭证环境下手动执行：**
        ```
        TF_ACC=1 go test -count=1 -timeout=10m \
          -run TestAccTencentCloudProviderRuntimeDataSource_basic \
          ./tencentcloud/services/tcprovider/...
        ```
- [ ] 9.7 真实 workspace 回归测试：需要生产环境，本地无法执行；请在有现有 SDKv2 资源的 workspace 上跑 `terraform plan` 确认 plan 为空
