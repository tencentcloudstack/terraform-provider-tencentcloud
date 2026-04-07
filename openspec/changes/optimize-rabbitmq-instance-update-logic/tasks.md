## 1. 代码修改 - Update 函数增强

- [x] 1.1 修改 `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate` 函数的 `immutableArgs` 列表
  - 从 `immutableArgs` 列表中移除 `"auto_renew_flag"`、`"band_width"`、`"enable_public_access"` 三个参数
  - 保持其他不可变参数不变

- [x] 1.2 在 `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate` 函数中添加 `auto_renew_flag` 参数的修改逻辑
  - 检测 `d.HasChange("auto_renew_flag")`
  - 如果有变化，将新值设置到 `request.AutoRenewFlag`
  - 设置 `needUpdate = true`

- [x] 1.3 在 `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate` 函数中添加 `band_width` 参数的修改逻辑
  - 检测 `d.HasChange("band_width")`
  - 如果有变化，将新值设置到 `request.Bandwidth`
  - 设置 `needUpdate = true`

- [x] 1.4 在 `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate` 函数中添加 `enable_public_access` 参数的修改逻辑
  - 检测 `d.HasChange("enable_public_access")`
  - 如果有变化，将新值设置到 `request.EnablePublicAccess`
  - 设置 `needUpdate = true`

## 2. 代码修改 - 异步操作等待机制

- [x] 2.1 为 `band_width` 修改添加异步等待机制
  - 在 `ModifyRabbitMQVipInstance` API 调用成功后，添加等待逻辑
  - 使用 `resource.Retry()` 进行轮询
  - 调用 `DescribeRabbitMQVipInstances` API 检查 `PublicNetworkTps` 值是否更新为目标值
  - 设置超时时间为 `tccommon.WriteRetryTimeout * 2`

- [x] 2.2 为 `enable_public_access` 修改添加异步等待机制
  - 在 `ModifyRabbitMQVipInstance` API 调用成功后，添加等待逻辑
  - 使用 `resource.Retry()` 进行轮询
  - 调用 `DescribeRabbitMQVipInstances` API 检查 `PublicDataStreamStatus` 值是否更新为目标状态
  - 设置超时时间为 `tccommon.WriteRetryTimeout * 2`

- [x] 2.3 确保 `auto_renew_flag` 修改不添加等待机制
  - 验证 `auto_renew_flag` 修改后立即调用 Read 函数刷新状态
  - 不添加异步等待逻辑

## 3. 代码修改 - 错误处理和日志记录

- [x] 3.1 添加日志记录
  - 在修改 `auto_renew_flag` 时记录日志，显示新值
  - 在修改 `band_width` 时记录日志，显示新值
  - 在修改 `enable_public_access` 时记录日志，显示新值
  - 在 API 调用失败时记录详细的错误信息

- [x] 3.2 完善错误处理
  - 确保所有 API 调用错误都被正确捕获和返回
  - 在异步操作超时时返回清晰的错误信息
  - 保持现有的错误处理模式一致

## 4. 测试用例添加

- [x] 4.1 添加 `auto_renew_flag` 修改的测试用例
  - 测试将 `auto_renew_flag` 从 `false` 修改为 `true` 的场景
  - 测试将 `auto_renew_flag` 从 `true` 修改为 `false` 的场景
  - 验证修改后状态是否正确更新

- [x] 4.2 添加 `band_width` 修改的测试用例
  - 测试修改 `band_width` 值的场景
  - 验证修改后等待机制是否正常工作
  - 验证修改后状态是否正确更新

- [x] 4.3 添加 `enable_public_access` 修改的测试用例
  - 测试将 `enable_public_access` 从 `false` 修改为 `true` 的场景
  - 测试将 `enable_public_access` 从 `true` 修改为 `false` 的场景
  - 验证修改后等待机制是否正常工作
  - 验证修改后状态是否正确更新

- [x] 4.4 添加多参数同时修改的测试用例
  - 测试同时修改 `auto_renew_flag`、`band_width`、`enable_public_access` 的场景
  - 验证所有参数都正确更新

- [x] 4.5 添加不可变参数修改的测试用例
  - 测试修改不可变参数（如 `zone_ids`、`vpc_id`、`node_spec` 等）时返回错误
  - 验证错误信息清晰明确

## 5. 代码验证

- [x] 5.1 运行格式化检查
  - 执行 `go fmt` 格式化修改后的代码文件
  - 确保代码符合 Go 标准格式

- [x] 5.2 运行单元测试
  - 执行 `go test -v` 运行相关单元测试
  - 确保所有新增测试用例通过
  - 确保现有测试用例不受影响

- [x] 5.3 验证代码编译
  - 执行 `go build` 确保代码编译通过
  - 检查是否有编译错误或警告

## 6. 文档更新（如果需要）

- [x] 6.1 检查资源文档文件
  - 检查 `website/docs/r/tdmq_rabbitmq_vip_instance.md` 文件是否存在
  - 评估是否需要更新文档以反映新的可修改参数

- [x] 6.2 更新文档（如果需要）
  - 如果文档存在且需要更新，明确标注 `auto_renew_flag`、`band_width`、`enable_public_access` 为可修改参数
  - 在文档中添加示例说明如何修改这些参数
  - 注意：文档应通过 `make doc` 命令生成，禁止手动编辑

## 7. 集成测试（可选）

- [x] 7.1 执行集成测试（如果有环境）
  - 设置 `TF_ACC=1` 环境变量
  - 配置 `TENCENTCLOUD_SECRET_ID` 和 `TENCENTCLOUD_SECRET_KEY`
  - 运行完整的集成测试验证功能

- [x] 7.2 验证向后兼容性
  - 使用现有的 Terraform 配置进行测试
  - 确保 `terraform plan` 不显示不必要的变更
  - 确保现有资源正常工作

## 8. 最终验证

- [x] 8.1 检查所有修改的文件
  - 确认 `resource_tc_tdmq_rabbitmq_vip_instance.go` 文件修改正确
  - 确认 `resource_tc_tdmq_rabbitmq_vip_instance_test.go` 文件测试用例添加正确

- [x] 8.2 检查代码质量
  - 确保代码遵循项目规范
  - 确保错误处理和日志记录符合项目标准
  - 确保代码注释清晰完整

- [x] 8.3 准备提交
  - 检查是否所有任务都已完成
  - 确认代码可以通过所有测试
  - 准备提交代码变更
