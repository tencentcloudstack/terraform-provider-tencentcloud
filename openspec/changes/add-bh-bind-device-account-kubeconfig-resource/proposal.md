## Why

腾讯云堡垒机（BH）支持为已有的"容器账号"绑定 kubeconfig 凭据，使运维人员可以通过堡垒机访问 Kubernetes 集群。该能力在腾讯云控制台和云 API（`BindDeviceAccountKubeconfig`）已经发布，但 Terraform Provider 尚未暴露，导致用户必须在控制台手动绑定，无法做基础设施代码化（IaC）。

本 change 新增配置型资源 `tencentcloud_bh_bind_device_account_kubeconfig`，让用户在 HCL 里声明"为某个容器账号绑定 kubeconfig"，由 Provider 调用 `BindDeviceAccountKubeconfig` 完成绑定与变更。

## What Changes

- 新增配置型（config）资源 `tencentcloud_bh_bind_device_account_kubeconfig`。
- 新增 Schema 字段，**严格按 BindDeviceAccountKubeconfig 接口入参 1:1 映射**：`id`（容器账号Id，uint64）、`kubeconfig`（凭据，敏感）、`manage_dimension`（托管维度，1=集群）。
- 资源 ID 设为 HCL `id` 字段的字符串形式（`fmt.Sprintf("%d", id)`），与接口 `Id` 入参对齐。
- `id` 字段 `ForceNew: true`：换账号 Id = 换绑定对象，必须重建。
- `kubeconfig` 和 `manage_dimension` 可变更：Update 重新调 `BindDeviceAccountKubeconfig`，API 语义为"覆盖式绑定"。
- `kubeconfig` 字段 `Sensitive: true`，避免凭据出现在 plan 输出。
- Read 直接 `return nil`（业务侧暂无对应查询接口），不刷新 state；与 `tencentcloud_waf_owasp_rule_status_config` 中"配置型资源 + 无独立查询"的处理方式一致。
- Delete 直接 `return nil`（API 未提供独立解绑/UnbindDeviceAccountKubeconfig 接口）。
- 所有 SDK 调用必须用 `resource.Retry(...)` 包装。
- 新增 `.md` 资源文档（命名与 `resource_tc_config_compliance_pack.md` 同款）+ 验收测试 `_test.go`（命名与 `resource_tc_config_compliance_pack_test.go` 同款）。
- 在 `tencentcloud/provider.go` 注册新资源 `tencentcloud_bh_bind_device_account_kubeconfig`，并在 `tencentcloud/provider.md` 的 Bastion Host(BH) 段追加资源名。

## Capabilities

### New Capabilities

- `bh-bind-device-account-kubeconfig-resource`: 新增 `tencentcloud_bh_bind_device_account_kubeconfig` 资源的 schema、CRUD（Create=Bind / Read=noop / Update=Re-Bind / Delete=noop）行为、ID 约定、字段约束、文档与测试规范，作为堡垒机容器账号 kubeconfig 绑定能力的 IaC 入口。

### Modified Capabilities

<!-- 不修改任何已有 capability 的 requirement，仅在 provider.go / provider.md 增加注册条目。 -->

## Impact

- 新文件：
  - `tencentcloud/services/bh/resource_tc_bh_bind_device_account_kubeconfig.go`
  - `tencentcloud/services/bh/resource_tc_bh_bind_device_account_kubeconfig.md`
  - `tencentcloud/services/bh/resource_tc_bh_bind_device_account_kubeconfig_test.go`
  - `website/docs/r/bh_bind_device_account_kubeconfig.html.markdown`（由 `make doc` 生成）
- 既有文件改动：
  - `tencentcloud/provider.go`：在 bh 资源注册段追加 `"tencentcloud_bh_bind_device_account_kubeconfig": bh.ResourceTencentCloudBhBindDeviceAccountKubeconfig()` 一行。
  - `tencentcloud/provider.md`：在 Bastion Host(BH) 的 Resource 段追加 `tencentcloud_bh_bind_device_account_kubeconfig`，使 gendoc 索引可识别。
- 依赖：现有 SDK `tencentcloud-sdk-go/tencentcloud/bh/v20230418` 已包含 `BindDeviceAccountKubeconfig` 与 `BindDeviceAccountKubeconfigWithContext` 接口，且对应 model `BindDeviceAccountKubeconfigRequest/Response`，无需升级 SDK。
- 不修改任何既有资源/数据源 schema，对现有用户零影响。
