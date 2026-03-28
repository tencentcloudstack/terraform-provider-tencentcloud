# Tasks: 优化 RabbitMQ 实例的 update 逻辑

- [x] Task 1: 修改 `immutableArgs` 列表
  - 将可修改的参数从 `immutableArgs` 列表中移除
  - 保留真正的不可修改参数
  - 文件: `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go`
  - 代码编译通过
  - 不可修改参数列表正确

- [x] Task 2: 为可修改参数添加 Update 逻辑
  - 为每个可修改的参数添加相应的 API 调用逻辑
  - 支持修改的参数: node_spec, node_num, storage_size, band_width, enable_public_access, auto_renew_flag, cluster_version
  - 保留原有的 cluster_name 和 resource_tags 修改逻辑
  - 所有可修改参数都有对应的更新逻辑
  - 代码编译通过
  - 没有语法错误

- [x] Task 3: 添加异步操作等待逻辑
  - 对于需要异步等待的操作（如节点规格、节点数量、存储规格变更），添加状态等待逻辑
  - 添加重试机制和超时处理
  - 异步操作能够正确等待完成
  - 超时处理正确
  - 错误处理正确

- [x] Task 4: 添加 Timeouts 配置支持
  - 在 Schema 中添加 timeouts 配置
  - 支持 create, update, delete 操作的自定义超时
  - timeouts 配置正确添加到 Schema
  - 支持自定义超时配置
  - 默认超时值合理

- [x] Task 5: 更新资源文档
  - 更新 `resource_tc_tdmq_rabbitmq_vip_instance.md` 文档
  - 明确说明哪些参数可以修改，哪些不能修改
  - 添加自定义超时配置的示例
  - 添加参数修改的示例代码
  - 更新 Import 部分，说明导入后的参数更新能力
  - 文档内容准确
  - 示例代码正确
  - 参数说明清晰

- [x] Task 6: 编写单元测试
  - 为以下场景编写测试用例：节点规格升级、节点数量扩容、存储规格扩容、带宽调整、公网访问开关切换、集群版本升级、自动续费标识修改、不可修改参数的错误处理
  - 所有测试用例通过
  - 覆盖主要使用场景
  - 代码覆盖率达到要求

- [x] Task 7: 编写集成测试
  - 测试多种参数组合的修改场景
  - 测试错误场景和回滚机制
  - 测试并发修改场景
  - 集成测试通过
  - 覆盖复杂场景
  - 错误处理正确

- [x] Task 8: 代码审查和优化
  - 检查代码风格是否符合项目规范
  - 检查错误处理是否完善
  - 检查日志输出是否合理
  - 检查是否有性能优化空间
  - 检查是否有潜在的安全问题
  - 代码通过代码审查
  - 没有明显的代码质量问题
  - 错误处理完善

- [x] Task 9: 本地验证测试
  - 编译 Provider
  - 运行单元测试
  - 运行集成测试（需要配置腾讯云 API 凭证）
  - 手动测试各种修改场景
  - 所有测试通过
  - 手动测试场景验证成功
  - 没有明显的 Bug

- [x] Task 10: 文档更新和代码注释
  - 在代码中添加必要的注释，说明关键逻辑
  - 更新 CHANGELOG.md，记录本次变更
  - 更新其他相关文档
  - 代码注释清晰
  - 文档更新完整
  - 变更记录准确

- [x] Task 11: 准备提交 PR
  - 创建新的分支：`optimize-rabbitmq-instance-update-logic`
  - 提交所有变更
  - 推送到远程仓库
  - 创建 Pull Request
  - 填写 PR 描述，包括：变更摘要、实现的功能、测试结果、相关的 issue 或链接
  - 分支创建正确
  - 提交信息清晰
  - PR 描述完整
