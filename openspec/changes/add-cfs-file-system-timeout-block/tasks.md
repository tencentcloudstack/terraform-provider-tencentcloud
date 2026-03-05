# Tasks: 为 tencentcloud_cfs_file_system 添加 Timeout 块支持

## 实施任务清单

### 任务 1: 添加 Timeouts 字段到 Schema
- [x] 在 `resource_tc_cfs_file_system.go` 中添加 Timeouts 字段
- [x] 确保 `time` 包已导入
- [x] 只添加 Create timeout (20 分钟)

### 任务 2: 修改 Create 函数使用 timeout
- [x] 修改第 196 行的状态等待 retry 使用 `d.Timeout(schema.TimeoutCreate)`
- [x] 保持第 172-189 行的 API 调用 retry 不变

### 任务 3: 代码格式化和编译验证
- [x] 运行 `go fmt`
- [x] 编译验证通过

### 任务 4: 运行测试
- [x] 编译测试通过
- [x] 无新增 lint 错误

### 任务 5: 更新文档
- [x] 在 `website/docs/r/cfs_file_system.html.markdown` 添加 Timeouts 说明

### 任务 6: 创建变更日志
- [x] 在 `.changelog/` 目录创建变更日志文件

### 任务 7: 更新任务清单
- [x] 标记所有任务为完成
