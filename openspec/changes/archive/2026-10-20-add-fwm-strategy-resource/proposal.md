## Why

腾讯云防火墙（FWM，Edge / Nat / VPC 防火墙统一管控面）提供 `CreateStrategy` 接口，用于一次性下发"前区规则组 + 后区规则组"策略，并返回本次创建生成的**策略 Id 列表**（`StrategyIds`）。当前 Terraform Provider 尚未暴露 `fwm_strategy` 资源，用户只能在控制台手工下发策略，无法做基础设施代码化；且创建返回的 `StrategyIds` 是后续 `DescribeStrategy` / `DeleteStrategy`（按策略 Id 列表操作）的入参身份键，必须被管理。

本 change 新增 `RESOURCE_KIND_GENERAL` 资源 `tencentcloud_fwm_strategy`，并**将该新建资源的 Create 出参 `response.StrategyIds` 映射为 Terraform 资源参数 `strategy_ids`**（Computed，由云 API 返回），作为本次唯一的新增参数变更。

## What Changes

- 新增通用资源 `tencentcloud_fwm_strategy`（Create / Read / Update / Delete 全生命周期），代码风格严格参考 `tencentcloud_igtm_strategy`。
- 新增资源参数 `strategy_ids`（Computed，`TypeList` of `TypeString`）：映射 `CreateStrategy` 出参 `response.StrategyIds`。Create 完成后 `d.Set("strategy_ids", ...)`，并在 Read 中刷新。
  - 资源 ID 采用 `StrategyIds` 用 `tccommon.FILED_SP` 拼接（与"多 id 联合"惯例一致）；由于 `StrategyIds` 是创建产物，ID 由出参生成，支持 import（import 时使用联合 id `s1#s2#...`）。
- Schema 入参严格按 `CreateStrategyRequest` 字段 1:1 映射（snake_case 化）：`product`、`receive_account`、`pre_strategy`、`post_strategy`、`receive_group`，其中 `product` 必填，其余按业务语义。
  - `pre_strategy` / `post_strategy` 为 `StrategyReq` 列表（`group_id`、`sequence`）。
- Update：因 `CreateStrategy` 的入参字段均无独立 Modify 接口对应字段（fwm SDK 仅有 `CreateStrategy` / `DescribeStrategy` / `DeleteStrategy`，无 `ModifyStrategy`），按"CRD 资源"惯例：仅 `product` 标 `ForceNew`，其余顶层字段加入 `immutableArgs`，update 中检测到变化则返回 error。`strategy_ids` 为 Computed，不在 immutableArgs 中。
- CRUD 调用全部用 `resource.Retry(tccommon.ReadRetryTimeout / WriteRetryTimeout, ...)` 包装，错误用 `tccommon.RetryError(e)` 包装。retry 块内只做 API 调用，SetId / Set 等放到 retry 块外。
- 在 `tencentcloud/connectivity/client.go` 新增 `UseFwmV20250611Client()` 访问器（参考 `UseIgtmV20231024Client`）。
- 在 `tencentcloud/provider.go` 注册新资源 `tencentcloud_fwm_strategy`，并在 `tencentcloud/provider.md` 的 FireWall(FWM) 段追加资源名。
- 新增 `tencentcloud/services/fwm/resource_tc_fwm_strategy.go` + `.md` + `_test.go`（测试用 gomonkey mock 云 API，不用 TF 测试套件）。

## Capabilities

### New Capabilities
- `fwm-strategy-resource`: 新增 `tencentcloud_fwm_strategy` 资源，封装 `CreateStrategy` / `DescribeStrategy` / `DeleteStrategy`，其中 `CreateStrategy` 出参 `StrategyIds` 映射为 Computed 参数 `strategy_ids`。

### Modified Capabilities
<!-- 不修改任何已有 capability 的 requirement。 -->

## Impact

- 新文件：
  - `tencentcloud/services/fwm/resource_tc_fwm_strategy.go`
  - `tencentcloud/services/fwm/resource_tc_fwm_strategy.md`
  - `tencentcloud/services/fwm/resource_tc_fwm_strategy_test.go`
  - `tencentcloud/services/fwm/service_tencentcloud_fwm.go`（可选，若需 service 层封装 Describe）
  - `website/docs/r/fwm_strategy.html.markdown`（由 `make doc` 生成）
- 既有文件改动：
  - `tencentcloud/connectivity/client.go`：新增 `UseFwmV20250611Client()` 访问器与 `fwmv20250611Conn` 字段。
  - `tencentcloud/provider.go`：注册 `"tencentcloud_fwm_strategy": fwm.ResourceTencentCloudFwmStrategy()`。
  - `tencentcloud/provider.md`：在 FWM 段追加 `tencentcloud_fwm_strategy`。
- 依赖：现有 `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611` 已包含 `CreateStrategy` / `DescribeStrategy` / `DeleteStrategy` 及对应 Request/Response/model，无需升级 SDK（go.mod 已存在 `fwm v1.3.194`）。
- 不修改任何既有资源/数据源 schema，对现有用户零影响。