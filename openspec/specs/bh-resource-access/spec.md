# bh-resource-access Specification

## Purpose

Defines the update-phase modification capability for the `web_access` and `client_access` fields of the `tencentcloud_bh_resource` resource. These fields were previously `ForceNew`, requiring instance destruction and rebuild to change. This specification covers removing `ForceNew`, calling the dedicated `EnableWebAccess`/`DisableWebAccess` and `EnableClientTcpAccess`/`DisableClientTcpAccess` APIs on change, and polling `DescribeResources` until the target field stabilizes.

## Requirements

### Requirement: web_access 字段支持在 update 阶段修改

`tencentcloud_bh_resource` 资源的 `web_access` 字段 SHALL 移除 `ForceNew` 属性，使得用户在不销毁重建堡垒机实例的前提下即可变更 web 访问开关。

在 `resourceTencentCloudBhResourceUpdate` 中，当 `web_access` 发生变更时，系统 SHALL 根据新值调用对应接口：
- 新值为 `1` 时调用 `EnableWebAccess`（入参 `ResourceId`）。
- 新值为 `0` 时调用 `DisableWebAccess`（入参 `ResourceId`）。

接口调用成功后，系统 SHALL 通过 `DescribeResources` 轮询 `WebAccess` 字段直至其稳定到目标值（Enable 等待 `== 1`，Disable 等待 `== 0`），轮询采用 `resource.Retry(tccommon.ReadRetryTimeout*7, ...)`。

`web_access` 在 Create 阶段的行为 SHALL 保持不变（仍通过 `DeployResource` 入参下发）。

#### Scenario: 从关闭切换为开启 web 访问
- **WHEN** 用户将 `tencentcloud_bh_resource` 的 `web_access` 从 `0` 修改为 `1` 并执行 `terraform apply`
- **THEN** 系统 SHALL 调用 `EnableWebAccess`（`ResourceId` 为当前实例 id），并在 `DescribeResources` 返回的 `WebAccess` 等于 `1` 后完成 update

#### Scenario: 从开启切换为关闭 web 访问
- **WHEN** 用户将 `tencentcloud_bh_resource` 的 `web_access` 从 `1` 修改为 `0` 并执行 `terraform apply`
- **THEN** 系统 SHALL 调用 `DisableWebAccess`（`ResourceId` 为当前实例 id），并在 `DescribeResources` 返回的 `WebAccess` 等于 `0` 后完成 update

#### Scenario: web_access 未变更时不触发开关接口
- **WHEN** 用户执行 `terraform apply` 且 `web_access` 未发生变更
- **THEN** 系统 SHALL NOT 调用 `EnableWebAccess` 或 `DisableWebAccess`

### Requirement: client_access 字段支持在 update 阶段修改

`tencentcloud_bh_resource` 资源的 `client_access` 字段 SHALL 移除 `ForceNew` 属性，使得用户在不销毁重建堡垒机实例的前提下即可变更客户端访问开关。

在 `resourceTencentCloudBhResourceUpdate` 中，当 `client_access` 发生变更时，系统 SHALL 根据新值调用对应接口：
- 新值为 `1` 时调用 `EnableClientTcpAccess`（入参 `ResourceId`）。
- 新值为 `0` 时调用 `DisableClientTcpAccess`（入参 `ResourceId`）。

接口调用成功后，系统 SHALL 通过 `DescribeResources` 轮询 `ClientAccess` 字段直至其稳定到目标值（Enable 等待 `== 1`，Disable 等待 `== 0`），轮询采用 `resource.Retry(tccommon.ReadRetryTimeout*7, ...)`。

`client_access` 在 Create 阶段的行为 SHALL 保持不变（仍通过 `DeployResource` 入参下发）。

#### Scenario: 从关闭切换为开启客户端访问
- **WHEN** 用户将 `tencentcloud_bh_resource` 的 `client_access` 从 `0` 修改为 `1` 并执行 `terraform apply`
- **THEN** 系统 SHALL 调用 `EnableClientTcpAccess`（`ResourceId` 为当前实例 id），并在 `DescribeResources` 返回的 `ClientAccess` 等于 `1` 后完成 update

#### Scenario: 从开启切换为关闭客户端访问
- **WHEN** 用户将 `tencentcloud_bh_resource` 的 `client_access` 从 `1` 修改为 `0` 并执行 `terraform apply`
- **THEN** 系统 SHALL 调用 `DisableClientTcpAccess`（`ResourceId` 为当前实例 id），并在 `DescribeResources` 返回的 `ClientAccess` 等于 `0` 后完成 update

#### Scenario: client_access 未变更时不触发开关接口
- **WHEN** 用户执行 `terraform apply` 且 `client_access` 未发生变更
- **THEN** 系统 SHALL NOT 调用 `EnableClientTcpAccess` 或 `DisableClientTcpAccess`

### Requirement: web_access 与 client_access 的变更各自独立判定

`web_access` 与 `client_access` 的变更处理 SHALL 各自通过独立的 `d.HasChange` 判定，互不影响；二者 SHALL NOT 被加入控制 `ModifyResource` 调用的 `mutableArgs` 集合，因为 `ModifyResource` 接口不接受这两个字段，变更 SHALL 通过各自的 Enable/Disable 开关接口完成。

在单次 `terraform apply` 中同时变更 `web_access` 和 `client_access` 时，系统 SHALL 串行执行各自的"调用接口 + 轮询等待"流程。

#### Scenario: 同时变更 web_access 与 client_access
- **WHEN** 用户在同一次 `terraform apply` 中将 `web_access` 改为 `1` 且将 `client_access` 改为 `0`
- **THEN** 系统 SHALL 串行执行：先完成 `web_access` 的 `EnableWebAccess` + 轮询等待，再完成 `client_access` 的 `DisableClientTcpAccess` + 轮询等待

### Requirement: Read 回填与空指针保护

`resourceTencentCloudBhResourceRead` 中对 `web_access`、`client_access` 的回填逻辑 SHALL 保持既有行为：仅当 `DescribeResources` 返回的对应字段非 nil 时调用 `d.Set`。

新增的 update 分支中，对 `DescribeResources` 返回值 SHALL 进行空指针保护：`result == nil || result.Response == nil || result.Response.ResourceSet == nil || len(result.Response.ResourceSet) == 0` 时返回 `resource.NonRetryableError`；待轮询字段为 nil 时返回 `resource.NonRetryableError`。

#### Scenario: 轮询期间 DescribeResources 返回空
- **WHEN** update 分支轮询 `DescribeResources` 时返回的 `Response`/`ResourceSet` 为空
- **THEN** 系统 SHALL 返回 `NonRetryableError` 终止重试，而不是无限等待
