## Why

腾讯云 TEO Load Balancer 的健康检查策略新增了 `ProbeCluster`（探测集群区域）字段，允许 HTTP 专用型（V2）实例选择发起健康探测的集群区域范围，以降低探测请求量。同时顶层 `type` 字段新增了 `HTTP_V2` 取值。当前 `tencentcloud_teo_load_balancer` 资源尚未暴露该能力，用户无法通过 Terraform 声明式配置探测集群区域，也无法使用 V2 版本实例并配合该字段。腾讯云 SDK 已在 `HealthChecker` struct 中新增 `ProbeCluster *string` 字段，并在 CRUD 接口中可用，具备接入条件。

## What Changes

- 在 `tencentcloud_teo_load_balancer` 资源的 `health_checker` 嵌套块 schema 中新增可选字段 `probe_cluster`（TypeString, Optional），用于指定发起健康探测的集群区域范围。合法取值：`global`、`mainland_china`、`overseas`，默认为 `global`。仅 `HTTP_V2` 类型实例支持设置。
- 在 `buildTeoHealthChecker` 函数中补充 `probe_cluster` 字段的写入逻辑（schema → API request）。
- 在 `flattenTeoHealthChecker` 函数中补充 `probe_cluster` 字段的读取逻辑（API response → schema）。
- 更新顶层 `type` 字段的 Description，新增 `HTTP_V2` 取值说明，使其与最新云 API 文档一致（HTTP 专用型 V2 支持 HTTP 专用型和通用型源站组，仅被站点加速服务引用，支持选择探测区域；HTTP 专用型 V1 不支持选择探测区域；GENERAL 通用型不支持选择探测区域）。
- 在单元测试 `resource_tc_teo_load_balancer_test.go` 中补充覆盖 `probe_cluster` 字段的测试用例。

非破坏性：`probe_cluster` 为新增 Optional 字段，未配置时不影响已有 state 与 TF 配置；`type` 字段描述更新不影响 schema 行为。

## Capabilities

### New Capabilities
<!-- 无新增 capability，本次仅修改现有 capability -->

### Modified Capabilities
- `teo-load-balancer-resource`: 在 `health_checker` 嵌套块中新增可选字段 `probe_cluster`（探测集群区域）；更新 `type` 字段描述以包含 `HTTP_V2` 取值。

## Impact

- 代码：
  - `tencentcloud/services/teo/resource_tc_teo_load_balancer.go`（schema 新增 `probe_cluster`、更新 `type` 描述；`buildTeoHealthChecker` / `flattenTeoHealthChecker` 补充读写逻辑）
  - `tencentcloud/services/teo/resource_tc_teo_load_balancer_test.go`（新增 `probe_cluster` 字段测试用例）
- 依赖：使用已 vendored 的 `tencentcloud-sdk-go` 中 `teov20220901.HealthChecker.ProbeCluster` 字段，无需变更 vendor。
- 向后兼容：新增 Optional 字段，不影响已有 state 与 TF 配置；`type` 描述更新无行为变化。
- 文档：由收尾阶段 `make doc` 自动生成 `website/docs/` 下文档（禁止手改）。
