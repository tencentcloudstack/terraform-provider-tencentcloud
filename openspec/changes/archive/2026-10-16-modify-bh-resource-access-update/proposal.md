## Why

`tencentcloud_bh_resource` 资源当前的 `web_access` 与 `client_access` 字段被标记为 `ForceNew: true`，用户一旦需要变更这两个访问开关就必须销毁并重建整个堡垒机实例（而堡垒机资源本身不支持删除，需工单处理），导致这两个字段在实际使用中无法调整。腾讯云 bh v20230418 SDK 已提供独立的 `EnableWebAccess`/`DisableWebAccess` 与 `EnableClientTcpAccess`/`DisableClientTcpAccess` 接口，且 `DescribeResources` 返回的 `WebAccess`/`ClientAccess` 字段含有"开通中/关闭中"的中间态，可以在 update 流程中通过轮询等待生效。因此需要移除两者的 ForceNew 属性，并参照同资源中已有的 `intranet_access`/`external_access` 写法，在 update 模块中支持这两个字段的修改。

## What Changes

- 移除 `tencentcloud_bh_resource` schema 中 `web_access` 字段的 `ForceNew: true` 属性，使其支持在 update 阶段修改。
- 移除 `tencentcloud_bh_resource` schema 中 `client_access` 字段的 `ForceNew: true` 属性，使其支持在 update 阶段修改。
- 在 `resourceTencentCloudBhResourceUpdate` 中新增 `web_access` 字段的变更处理逻辑：当值变为 `1` 时调用 `EnableWebAccess`，当值变为 `0` 时调用 `DisableWebAccess`；调用成功后通过 `DescribeResources` 轮询 `WebAccess` 字段等待生效（参照 `intranet_access`/`external_access` 的 wait 逻辑）。
- 在 `resourceTencentCloudBhResourceUpdate` 中新增 `client_access` 字段的变更处理逻辑：当值变为 `1` 时调用 `EnableClientTcpAccess`，当值变为 `0` 时调用 `DisableClientTcpAccess`；调用成功后通过 `DescribeResources` 轮询 `ClientAccess` 字段等待生效（参照 `intranet_access`/`external_access` 的 wait 逻辑）。
- 在 `resource_tc_bh_resource_test.go` 中补充 `web_access` 与 `client_access` 字段 update 变更的单元测试用例。

## Capabilities

### New Capabilities
- `bh-resource-access`: 堡垒机资源 `tencentcloud_bh_resource` 的 web 访问与客户端访问开关在 update 阶段的修改能力，覆盖 `web_access`/`client_access` 字段的 ForceNew 移除、Enable/Disable 接口调用及 DescribeResources 轮询等待生效逻辑。

### Modified Capabilities
<!-- 无既有 spec 需要修改 -->

## Impact

- **受影响代码**:
  - `tencentcloud/services/bh/resource_tc_bh_resource.go`：修改 `web_access`、`client_access` 两个 schema 字段（移除 ForceNew），并在 `resourceTencentCloudBhResourceUpdate` 中新增对应分支处理逻辑。
  - `tencentcloud/services/bh/resource_tc_bh_resource_test.go`：补充 update 变更相关单元测试用例。
  - `tencentcloud/services/bh/resource_tc_bh_resource.md`：更新文档示例，体现两字段支持修改。
- **依赖的云 API**（均已在 vendor SDK `bh/v20230418` 中存在）:
  - `EnableWebAccess` / `DisableWebAccess`（入参仅 `ResourceId`）
  - `EnableClientTcpAccess` / `DisableClientTcpAccess`（入参仅 `ResourceId`）
  - `DescribeResources`（轮询等待，返回 `WebAccess`/`ClientAccess` 中间态：0 关闭 / 1 开启 / 2 开通中 / 3 关闭中）
- **兼容性**: 字段类型与含义不变，仅移除 ForceNew，属于向后兼容的能力增强；既有的 `intranet_access`/`external_access` update 写法已验证该模式可行。