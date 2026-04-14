## 1. 资源文件创建与 Schema 定义

- [x] 1.1 创建资源文件 `tencentcloud/services/teo/resource_tc_teo_alias_domain.go`，定义资源的基本结构和导入包
- [x] 1.2 在资源文件中定义 Schema，包括 zone_id（String, Required）、alias_name（String, Required）、target_name（String, Required）、paused（Bool, Optional, Computed）参数
- [x] 1.3 在 Schema 中添加 Timeouts 块，支持 create、update、delete 操作的超时配置

## 2. CRUD 函数实现

- [x] 2.1 实现 `resourceTencentCloudTeoAliasDomainCreate` 函数，调用 CreateAliasDomain API，支持创建时指定 paused 状态
- [x] 2.2 在 Create 函数中实现异步操作轮询机制，调用 DescribeAliasDomains API 直到资源创建完成
- [x] 2.3 实现 `resourceTencentCloudTeoAliasDomainRead` 函数，调用 DescribeAliasDomains API 读取资源状态，更新 Terraform state
- [x] 2.4 实现 `resourceTencentCloudTeoAliasDomainUpdate` 函数，检测参数变化并调用对应的 ModifyAliasDomain 或 ModifyAliasDomainStatus API
- [x] 2.5 在 Update 函数中实现多场景处理：单独更新 target_name、单独更新 paused、同时更新两者
- [x] 2.6 在 Update 函数中实现异步操作轮询机制，等待每个操作完成
- [x] 2.7 实现 `resourceTencentCloudTeoAliasDomainDelete` 函数，调用 DeleteAliasDomain API
- [x] 2.8 在 Delete 函数中实现异步操作轮询机制，调用 DescribeAliasDomains API 直到资源删除完成
- [x] 2.9 实现 `resourceTencentCloudTeoAliasDomainImporter` 函数，支持导入已有资源，使用 zone_id#alias_name 格式解析 ID

## 3. 辅助函数实现

- [x] 3.1 实现复合 ID 解析函数，将 zone_id#alias_name 格式的 ID 解析为 zone_id 和 alias_name
- [x] 3.2 实现 DescribeAliasDomains API 调用的封装函数，处理 zone_id 过滤和 alias_name 匹配
- [x] 3.3 实现轮询等待函数，使用 helper.Retry() 机制轮询 DescribeAliasDomains API 直到预期状态
- [x] 3.4 实现 API 响应字段到 Terraform Schema 的映射函数，处理 zone_id、alias_name、target_name、paused 字段

## 4. Service 层注册

- [x] 4.1 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中导入新创建的资源包
- [x] 4.2 在 `func init()` 或资源注册函数中添加 `resourceTencentCloudTeoAliasDomain()` 到 Provider 的 Resources 列表

## 5. 单元测试实现

- [ ] 5.1 创建测试文件 `tencentcloud/services/teo/resource_tc_teo_alias_domain_test.go`，定义测试套件和 mock 客户端
- [ ] 5.2 实现 Create 操作的单元测试，包括成功创建、创建时设置 paused=false、缺少必需参数的场景
- [ ] 5.3 实现 Read 操作的单元测试，包括成功读取、资源不存在的场景
- [ ] 5.4 实现 Update 操作的单元测试，包括更新 target_name、更新 paused 状态、同时更新两者的场景
- [ ] 5.5 实现 Delete 操作的单元测试，包括成功删除、删除已不存在资源的场景
- [ ] 5.6 实现 Import 操作的单元测试，包括成功导入、ID 格式错误的场景
- [ ] 5.7 实现异步操作轮询的单元测试，验证轮询机制正确调用 DescribeAliasDomains API

## 6. 资源示例文件创建

- [x] 6.1 创建资源示例文件 `tencentcloud/services/teo/resource_tc_teo_alias_domain.md`，包含完整的 Terraform 配置示例
- [x] 6.2 在示例文件中添加资源创建、更新、删除的使用示例，展示所有参数的使用方式

## 7. 代码验证

- [x] 7.1 运行 `gofmt` 格式化新增的代码文件，确保代码格式符合 Go 标准
- [x] 7.2 检查代码编译，确保没有语法错误和类型错误

## 8. 文档生成

- [x] 8.1 运行 `make doc` 命令，自动生成 `website/docs/r/teo_alias_domain.html.markdown` 文档
- [x] 8.2 验证生成的文档包含正确的参数说明、示例和引用信息
