## 1. 资源代码实现

- [x] 1.1 创建 resource_tc_teo_alias_domain.go 文件
- [x] 1.2 定义资源 Schema，包括 zone_id（Required）、alias_name（Required）、target_name（Required）、paused（Optional）参数
- [x] 1.3 在 Schema 中定义 Timeouts 块，支持 Create、Update、Delete 的自定义超时配置
- [x] 1.4 实现 resourceTencentcloudTeoAliasDomainCreate 函数，调用 CreateAliasDomain 云 API
- [x] 1.5 在 Create 函数中添加轮询逻辑，使用 helper.Retry() 调用 DescribeAliasDomains API 直到创建成功
- [x] 1.6 实现 resourceTencentcloudTeoAliasDomainRead 函数，调用 DescribeAliasDomains 云 API 查询资源信息
- [x] 1.7 在 Read 函数中实现资源 ID 解析，从 zone_id#alias_name 格式中解析出 zone_id 和 alias_name
- [x] 1.8 实现 resourceTencentcloudTeoAliasDomainUpdate 函数，处理 target_name 和 paused 参数的更新
- [x] 1.9 在 Update 函数中实现 target_name 更新逻辑，调用 ModifyAliasDomain 云 API 并轮询
- [x] 1.10 在 Update 函数中实现 paused 状态更新逻辑，调用 ModifyAliasDomainStatus 云 API 并轮询
- [x] 1.11 实现 resourceTencentcloudTeoAliasDomainDelete 函数，调用 DeleteAliasDomain 云 API
- [x] 1.12 在 Delete 函数中添加轮询逻辑，使用 helper.Retry() 调用 DescribeAliasDomains API 直到资源被删除
- [x] 1.13 实现 resourceTencentcloudTeoAliasDomainImporter 函数，支持使用 terraform import 导入资源
- [x] 1.14 添加错误处理和日志记录，使用 defer tccommon.LogElapsed() 和 defer tccommon.InconsistentCheck()
- [x] 1.15 注册资源到 teo 服务的资源列表中

## 2. 测试代码实现

- [x] 2.1 创建 resource_tc_teo_alias_domain_test.go 文件
- [x] 2.2 实现 TestAccTencentcloudTeoAliasDomain_basic 测试用例，测试基本 CRUD 操作
- [x] 2.3 在测试用例中覆盖创建资源场景，验证 zone_id、alias_name、target_name 参数
- [x] 2.4 在测试用例中覆盖更新资源场景，验证 target_name 更新
- [x] 2.5 在测试用例中覆盖暂停/启用状态更新场景，验证 paused 参数
- [x] 2.6 在测试 case 中覆盖删除资源场景，验证资源被正确删除
- [x] 2.7 实现 TestAccTencentcloudTeoAliasDomain_import 测试用例，测试资源导入功能
- [ ] 2.8 添加 mock 测试或使用真实的云 API 进行验收测试（需要 TF_ACC=1 环境变量）

## 3. 文档更新

- [x] 3.1 创建 resource_tc_teo_alias_domain.md 示例文件
- [x] 3.2 在示例文件中添加资源使用示例，包括完整的 Terraform 配置
- [x] 3.3 在示例文件中说明资源 ID 格式（zone_id#alias_name）
- [x] 3.4 在示例文件中说明 Timeouts 配置方法
- [x] 3.5 在示例文件中说明导入资源的方法
- [ ] 3.6 运行 make doc 命令生成 website/docs/ 下的文档

## 4. 代码验证

- [x] 4.1 使用 gofmt 格式化新增的 Go 代码
- [ ] 4.2 运行单元测试验证代码逻辑
- [ ] 4.3 运行验收测试（TF_ACC=1）验证资源功能（如果环境支持）
- [ ] 4.4 检查代码是否符合 provider 编码规范和最佳实践
