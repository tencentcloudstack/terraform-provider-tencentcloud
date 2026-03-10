# Tasks: 为 tencentcloud_clb_instance 添加 Timeout 块支持

## Phase 1: Schema 修改
- [x] 1.1 在 `resource_tc_clb_instance.go` 的 `ResourceTencentCloudClbInstance()` 中添加 `Timeouts` 字段，声明 Create（10分钟）和 Update（10分钟）默认超时
- [x] 1.2 确保 `time` 包已导入

## Phase 2: Create 函数超时替换
- [x] 2.1 修改 `resourceTencentCloudClbInstanceCreate` 中所有 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 为 `resource.Retry(d.Timeout(schema.TimeoutCreate), ...)`
- [x] 2.2 修改 `waitForTaskFinishGetID` 和 `waitForTaskFinish` 调用，传入 `d.Timeout(schema.TimeoutCreate)` 作为超时参数（需新增带 timeout 参数的函数或直接内联替换）

## Phase 3: Update 函数超时替换
- [x] 3.1 修改 `resourceTencentCloudClbInstanceUpdate` 中所有 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 为 `resource.Retry(d.Timeout(schema.TimeoutUpdate), ...)`
- [x] 3.2 修改 `waitForTaskFinish` 调用，传入 `d.Timeout(schema.TimeoutUpdate)` 作为超时参数

## Phase 4: 代码质量与文档
- [x] 4.1 运行 `gofmt` 格式化代码，编译通过
- [x] 4.2 更新 `resource_tc_clb_instance.md` 文档，添加 Timeouts 说明
