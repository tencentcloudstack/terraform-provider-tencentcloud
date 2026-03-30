## 1. 研究和准备

- [x] 1.1 阅读 tencentcloud_teo_l7_acc_rule 资源现有代码，了解当前 Update 函数的实现
- [x] 1.2 查阅 TEO API 文档，确认 ImportZoneConfig API 的响应结构和 TaskId 字段位置
- [x] 1.3 确定任务状态查询 API 接口（如 DescribeZoneConfigRollBackTasks 或类似接口）
- [x] 1.4 确认任务状态的枚举值（成功、失败、进行中等状态的定义）
- [x] 1.5 查看 provider 中其他资源的异步等待实现示例（如 CDN、DNSPod 等服务）

## 2. 实现任务等待逻辑

- [x] 2.1 在 tencentcloud/services/teo/service_tencentcloud_teo.go 中添加任务状态查询的辅助函数
- [x] 2.2 在 tencentcloud/services/teo/resource_tencentcloud_teo_l7_acc_rule.go 的 Update 函数中，从 ImportZoneConfig 响应提取 TaskId
- [x] 2.3 实现基于 helper.Retry() 的任务状态轮询逻辑
- [x] 2.4 添加任务完成状态的判断逻辑（成功返回，失败返回错误）
- [x] 2.5 添加超时处理逻辑，使用资源的 Timeout 配置或默认超时时间
- [x] 2.6 添加详细的日志记录（开始等待、轮询进度、完成状态等）
- [x] 2.7 添加错误处理，包括 API 调用失败、任务失败等场景

## 3. 添加和更新测试

- [x] 3.1 在 resource_tencentcloud_teo_l7_acc_rule_test.go 中添加任务等待逻辑的单元测试
- [x] 3.2 添加任务成功完成场景的测试用例
- [x] 3.3 添加任务失败场景的测试用例
- [x] 3.4 添加任务超时场景的测试用例
- [x] 3.5 确保现有测试用例仍然通过，验证向后兼容性

## 4. 更新文档

- [x] 4.1 检查并更新 tencentcloud/services/teo/resource_tencentcloud_teo_l7_acc_rule.md 示例文件
- [x] 4.2 运行 `make doc` 命令生成 website/docs/r/teo_l7_acc_rule.md 文档
- [x] 4.3 验证生成的文档包含异步操作的说明

## 5. 代码验证和清理

- [x] 5.1 运行 `go build` 验证代码编译通过
- [x] 5.2 运行 `go vet` 进行静态代码检查
- [x] 5.3 运行 `gofmt` 确保代码格式符合规范
- [x] 5.4 运行 `TF_ACC=1 go test ./tencentcloud/services/teo/... -v -run TestAccTencentCloudTeoL7AccRule` 执行验收测试
- [x] 5.5 检查日志输出，确保日志级别和内容合理
- [x] 5.6 检查错误消息，确保包含 TaskId 等关键信息

## 6. 代码审查和最终验证

- [x] 6.1 自我审查代码，检查是否遵循项目的编码规范
- [x] 6.2 验证不破坏向后兼容性（不修改 schema，不影响现有功能）
- [x] 6.3 确保所有测试通过，包括单元测试和验收测试
- [x] 6.4 确认文档已正确更新
