## 1. 调研和验证

- [x] 1.1 查阅腾讯云 TDMQ RabbitMQ API 文档，确认哪些字段支持修改
- [x] 1.2 在测试环境中验证 node_spec 字段的修改操作
- [x] 1.3 在测试环境中验证 node_num 字段的修改操作
- [x] 1.4 在测试环境中验证 storage_size 字段的修改操作
- [x] 1.5 在测试环境中验证 band_width 字段的修改操作
- [x] 1.6 在测试环境中验证 enable_public_access 字段的修改操作
- [x] 1.7 根据验证结果更新 immutableArgs 列表，移除可修改的字段

## 2. 代码实现

- [x] 2.1 修改 resource_tc_tdmq_rabbitmq_vip_instance.go 中的 resourceTencentCloudTdmqRabbitmqVipInstanceUpdate 函数
- [x] 2.2 从 immutableArgs 列表中移除 node_spec 字段，添加相应的 update 逻辑
- [x] 2.3 从 immutableArgs 列表中移除 node_num 字段，添加相应的 update 逻辑
- [x] 2.4 从 immutableArgs 列表中移除 storage_size 字段，添加相应的 update 逻辑
- [x] 2.5 从 immutableArgs 列表中移除 band_width 字段，添加相应的 update 逻辑
- [x] 2.6 从 immutableArgs 列表中移除 enable_public_access 字段，添加相应的 update 逻辑
- [x] 2.7 增强错误处理，为每个可修改字段提供清晰的错误提示
- [x] 2.8 优化 update 逻辑，确保向后兼容性
- [x] 2.9 执行 go fmt 格式化代码

## 3. 测试和验证

- [x] 3.1 为 node_spec 字段的 update 操作编写单元测试
- [x] 3.2 为 node_num 字段的 update 操作编写单元测试
- [x] 3.3 为 storage_size 字段的 update 操作编写单元测试
- [x] 3.4 为 band_width 字段的 update 操作编写单元测试
- [x] 3.5 为 enable_public_access 字段的 update 操作编写单元测试
- [x] 3.6 为错误处理场景编写单元测试
- [x] 3.7 运行所有单元测试，确保测试通过
- [x] 3.8 在测试环境中验证所有修改操作的正确性
- [x] 3.9 验证错误处理和状态更新机制

## 4. 文档更新

- [x] 4.1 更新 resource_tc_tdmq_rabbitmq_vip_instance.md 示例文件，说明可修改的字段
- [x] 4.2 更新 CHANGELOG，记录本次变更
- [x] 4.3 验证文档更新的正确性

## 5. 代码验证

- [x] 5.1 执行 go build 验证代码编译正确性
- [x] 5.2 执行单元测试，确保所有测试通过
- [x] 5.3 执行 go vet 进行静态检查
- [x] 5.4 验证所有代码格式正确（go fmt）
