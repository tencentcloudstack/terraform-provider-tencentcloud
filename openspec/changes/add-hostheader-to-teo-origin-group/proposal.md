## Why

tencentcloud_teo_origin_group 资源当前缺少在创建时设置 HostHeader 参数的支持。虽然该参数已在 Schema 中定义（第 109-113 行）并在 Update 和 Read 逻辑中正确处理，但在 CreateOriginGroup API 调用时未传递该参数，导致用户无法在资源创建时指定 HostHeader，只能通过后续的 update 操作来设置。

## What Changes

- 在 CreateOriginGroupRequest 构建过程中添加 HostHeader 参数的处理
- 保持现有 Update 和 Read 逻辑不变
- 不修改 Schema 定义（已存在）

## Capabilities

### New Capabilities
- `teo-origin-group-hostheader`: Enable HostHeader parameter support during tencentcloud_teo_origin_group resource creation

### Modified Capabilities
- (无现有 capability 需要修改)

## Impact

- 受影响文件：tencentcloud/services/teo/resource_tc_teo_origin_group.go
- 具体修改：在 resourceTencentCloudTeoOriginGroupCreate 函数中添加 HostHeader 参数读取和设置逻辑（约在第 220 行之后）
- 不影响其他资源或数据源
- 无 API 兼容性影响（仅补全缺失的参数支持）
