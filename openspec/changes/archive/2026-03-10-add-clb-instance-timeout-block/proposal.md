# Change: 为 tencentcloud_clb_instance 添加 Timeout 块支持

## Why

当前 `tencentcloud_clb_instance` 资源的 create 和 update 方法中存在大量异步任务等待（`waitForTaskFinish` / `waitForTaskFinishGetID`），这些等待使用硬编码的超时值（`tccommon.WriteRetryTimeout` 和 `4*tccommon.ReadRetryTimeout`），用户无法根据实际场景自定义超时时间。

**异步等待点统计**：
- **Create**: 6 处 `waitForTaskFinish` 调用（创建CLB、设置安全组、设置日志、修改 target_region_info、设置 delete_protect、关联 endpoint）
- **Update**: 6 处 `waitForTaskFinish` 调用（修改 SLA、修改属性、设置安全组、设置日志、修改 project、EIP 解绑/绑定等待）

用户在复杂网络环境或大规模部署中可能遇到超时失败，需要能通过 `timeouts` 块自定义超时。

## What Changes

- 在 `resource_tc_clb_instance.go` 的 `ResourceTencentCloudClbInstance()` schema 中添加 `Timeouts` 字段，声明 Create 和 Update 的默认超时
- 修改 `resourceTencentCloudClbInstanceCreate` 函数中所有 `resource.Retry` 和 `waitForTaskFinish` 调用，使用 `d.Timeout(schema.TimeoutCreate)` 替换硬编码超时
- 修改 `resourceTencentCloudClbInstanceUpdate` 函数中所有 `resource.Retry` 和 `waitForTaskFinish` 调用，使用 `d.Timeout(schema.TimeoutUpdate)` 替换硬编码超时
- 更新文档 `resource_tc_clb_instance.md` 添加 Timeouts 说明

## Impact

- Affected specs: `clb-instance-resource`
- Affected code:
  - `tencentcloud/services/clb/resource_tc_clb_instance.go` — 添加 Timeouts 字段 + 修改 Create/Update 超时逻辑
  - `tencentcloud/services/clb/service_tencentcloud_clb.go` — `waitForTaskFinish` / `waitForTaskFinishGetID` 需新增接受 timeout 参数的版本
  - `tencentcloud/services/clb/resource_tc_clb_instance.md` — 文档更新
- 非破坏性变更：默认值保持向后兼容，用户可选配置
