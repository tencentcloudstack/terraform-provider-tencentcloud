## Why

tencentcloud_teo_origin_group 资源的 Create 函数存在缺失参数传递问题。虽然 schema 中已定义了 `host_header` 参数，CreateOriginGroup API 也支持该参数，但在创建资源时该参数未被传递给 API，导致用户无法在创建时指定回源 Host Header，必须在创建后通过 Update 操作才能设置该参数。这造成了用户体验问题，因为：
1. 用户需要两次操作（create + update）才能正确配置资源
2. 违背了声明式配置的原则
3. 与 Update 和 Read 操作的行为不一致（Update 和 Read 都正确处理了该参数）

## What Changes

- **修复 Create 函数中缺失的 host_header 参数传递逻辑**
  - 在 resourceTencentCloudTeoOriginGroupCreate 函数中添加对 `host_header` 参数的处理
  - 将用户配置的 `host_header` 值传递给 CreateOriginGroupRequest.HostHeader 字段
  - 保持与现有 Update 和 Read 函数的一致性

- **新增验收测试用例**
  - 验证创建时指定 host_header 参数能正确传递到 API
  - 确保创建后能正确读取回 host_header 值

## Capabilities

### New Capabilities
- `teo-origin-group-hostheader-support`: 为 tencentcloud_teo_origin_group 资源提供在创建时配置 host_header 参数的能力

### Modified Capabilities
- 无（仅修复实现错误，不改变 spec 级别的行为要求）

## Impact

- **影响文件**：
  - tencentcloud/services/teo/resource_tc_teo_origin_group.go
  - tencentcloud/services/teo/resource_tc_teo_origin_group_test.go（新增测试用例）

- **影响行为**：
  - 用户现在可以在创建 origin_group 时直接指定 host_header 参数
  - 不破坏现有配置，向后兼容
  - 不修改 schema，保持 API 稳定性

- **风险**：低风险
  - 仅在 Create 函数中添加缺失的参数传递逻辑
  - 不修改现有参数定义或删除任何功能
  - 符合现有的 Update 和 Read 函数实现模式
