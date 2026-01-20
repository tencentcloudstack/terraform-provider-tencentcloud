## 1. 实现数据源核心功能
- [x] 1.1 创建 `data_source_tc_ckafka_version.go` 文件
- [x] 1.2 实现 Schema 定义，包含输入参数和输出参数
- [x] 1.3 实现 `dataSourceTencentCloudCkafkaVersionRead` 函数
- [x] 1.4 调用 CKafka SDK 的 `DescribeCkafkaVersionWithContext` 方法
- [x] 1.5 处理 API 响应并设置 Terraform 状态

## 2. 集成到 Provider
- [x] 2.1 在 CKafka 服务的 provider 注册中添加新数据源
- [x] 2.2 确保数据源名称为 `tencentcloud_ckafka_version`

## 3. 编写测试
- [x] 3.1 创建 `data_source_tc_ckafka_version_test.go` 文件
- [x] 3.2 实现基本功能验收测试 `TestAccTencentCloudCkafkaVersionDataSource_basic`
- [x] 3.3 测试用例应验证返回的版本信息字段
- [x] 3.4 确保测试使用真实的 CKafka 实例 ID

## 4. 编写文档
- [x] 4.1 创建 `data_source_tc_ckafka_version.md` 文档文件
- [x] 4.2 包含参数说明、属性说明和使用示例
- [ ] 4.3 运行 `make doc` 生成最终文档

## 5. 代码质量保证
- [ ] 5.1 运行 `make fmt` 格式化代码
- [ ] 5.2 运行 `make lint` 检查代码质量
- [ ] 5.3 运行 `make test` 执行单元测试
- [ ] 5.4 运行 `TF_ACC=1 make testacc TEST=./tencentcloud/services/ckafka` 执行验证测试

## 6. 验证和文档
- [ ] 6.1 手动测试数据源功能
- [ ] 6.2 验证输出格式符合预期
- [ ] 6.3 确保错误处理正确（如实例不存在时）
- [ ] 6.4 检查生成的文档格式正确