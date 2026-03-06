## 1. Schema 修改
- [ ] 1.1 从 `resource_tc_cfs_file_system.go:312` 的 `immutableArgs` 列表中移除 `"capacity"`
- [ ] 1.2 在 `Timeouts` 块中添加 `Update: schema.DefaultTimeout(30 * time.Minute)`

## 2. 服务层实现
- [ ] 2.1 在 `service_tencentcloud_cfs.go` 中添加 `ScaleUpFileSystem` 方法
- [ ] 2.2 实现 API 调用逻辑，包含请求构造和错误处理
- [ ] 2.3 添加日志记录（request/response body）

## 3. Update 函数扩容逻辑
- [ ] 3.1 在 `resourceTencentCloudCfsFileSystemUpdate` 中添加 `capacity` 变更检测
- [ ] 3.2 实现容量验证逻辑：
  - [ ] 3.2.1 验证新容量大于旧容量（不允许缩容）
  - [ ] 3.2.2 读取 `storage_type`，验证是否为 TB 或 TP
  - [ ] 3.2.3 验证扩容步长（TB: 20480, TP: 10240）
- [ ] 3.3 调用 `cfsService.ScaleUpFileSystem` 执行扩容
- [ ] 3.4 实现状态等待逻辑：
  - [ ] 3.4.1 使用 `resource.Retry(d.Timeout(schema.TimeoutUpdate), ...)` 等待
  - [ ] 3.4.2 轮询文件系统状态直到变为 `available`
  - [ ] 3.4.3 处理 timeout 和其他错误

## 4. 文档更新
- [ ] 4.1 更新 `website/docs/r/cfs_file_system.html.markdown`：
  - [ ] 4.1.1 修改 `capacity` 参数说明，标注支持修改（仅扩容）
  - [ ] 4.1.2 在 Timeouts 章节添加 `update` 说明（默认 30 分钟）
  - [ ] 4.1.3 添加扩容示例（包含 timeouts.update 配置）
  - [ ] 4.1.4 说明扩容是异步操作，会轮询 LifeCycleState 状态
- [ ] 4.2 更新 `tencentcloud/services/cfs/resource_tc_cfs_file_system.md` 源文档
- [ ] 4.3 运行 `make doc` 重新生成文档

## 5. 变更日志
- [ ] 5.1 创建 `.changelog/<PR_NUMBER>.txt` 文件
- [ ] 5.2 内容格式：`enhancement: resource/tencentcloud_cfs_file_system: support capacity update for Turbo file systems`

## 6. 代码质量检查
- [ ] 6.1 运行 `make fmt` 格式化代码
- [ ] 6.2 运行 `make lint` 检查代码质量
- [ ] 6.3 确保没有新增 linter 错误

## 7. 编译和测试验证
- [ ] 7.1 运行 `go build` 验证代码编译通过
- [ ] 7.2 运行单元测试（如适用）
- [ ] 7.3 手动验证扩容功能（可选，推荐在测试环境）

## 依赖关系
- 任务 1.2 依赖于之前的 `add-cfs-file-system-timeout-block` 变更已完成
- 任务 3 依赖于任务 1 和任务 2 完成
- 任务 4 依赖于所有代码修改完成

## 并行化建议
- 任务 1 和任务 2 可以并行执行
- 任务 3 需要等待任务 1 和任务 2 完成
- 任务 4 和任务 5 可以与任务 3 并行（在代码逻辑确定后）
- 任务 6 和任务 7 在所有开发任务完成后执行
