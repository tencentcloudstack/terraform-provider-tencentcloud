## 1. 资源文件创建

- [x] 1.1 创建 resource_tc_teo_identify_zone_operation.go 文件
- [x] 1.2 创建 resource_tc_teo_identify_zone_operation_test.go 测试文件

## 2. Schema 定义

- [x] 2.1 定义 Resource 函数，返回 terraform.Resource 结构
- [x] 2.2 定义 zone_name 参数（TypeString，Required）
- [x] 2.3 定义 domain 参数（TypeString，Optional）
- [x] 2.4 定义 ascription 嵌套对象（TypeList，Computed，MaxItems: 1）
  - [x] 2.4.1 定义 Subdomain 字段（TypeString，Computed）
  - [x] 2.4.2 定义 RecordType 字段（TypeString，Computed）
  - [x] 2.4.3 定义 RecordValue 字段（TypeString，Computed）
- [x] 2.5 定义 file_ascription 嵌套对象（TypeList，Computed，MaxItems: 1）
  - [x] 2.5.1 定义 IdentifyPath 字段（TypeString，Computed）
  - [x] 2.5.2 定义 IdentifyContent 字段（TypeString，Computed）

## 3. 创建函数实现

- [x] 3.1 实现 resourceTencentCloudTeoIdentifyZoneCreate 函数
- [x] 3.2 在 Create 函数中添加 defer tccommon.LogElapsed()
- [x] 3.3 在 Create 函数中添加 defer tccommon.InconsistentCheck()
- [x] 3.4 调用 serviceTeoIdentifyZone 函数执行云 API 调用
- [x] 3.5 处理 API 返回的 ascription 信息，填充到 state 中
- [x] 3.6 处理 API 返回的 file_ascription 信息，填充到 state 中
- [x] 3.7 使用 d.Set() 将所有输出字段设置到 state
- [x] 3.8 返回成功状态（nil diagnostics）

## 4. Service 层实现

- [x] 4.1 创建或更新 service_tencentcloud_teo.go 文件
- [x] 4.2 实现 TeoIdentifyZone 函数，接受请求参数
- [x] 4.3 构造 IdentifyZoneRequest 对象，设置 ZoneName 和 Domain
- [x] 4.4 调用 client.IdentifyZone() 方法
- [x] 4.5 处理 API 返回的 IdentifyZoneResponse
- [x] 4.6 提取 ascription 和 file_ascription 信息
- [x] 4.7 返回结构化的认证信息

## 5. 错误处理

- [x] 5.1 检查 zone_name 参数是否为空，返回错误提示
- [x] 5.2 捕获云 API 调用错误，通过 diag.Errorf() 返回
- [x] 5.3 处理资源不存在错误
- [x] 5.4 处理权限不足错误

## 6. 单元测试

- [x] 6.1 编写测试用例：成功获取站点认证配置
  - [x] 6.1.1 Mock IdentifyZone API 返回成功结果
  - [x] 6.1.2 验证返回的 ascription 信息正确
  - [x] 6.1.3 验证返回的 file_ascription 信息正确
- [x] 6.2 编写测试用例：成功获取子域名认证配置
  - [x] 6.2.1 Mock IdentifyZone API 返回带 domain 的结果
  - [x] 6.2.2 验证 domain 参数正确传递
- [x] 6.3 编写测试用例：缺少 zone_name 参数
  - [x] 6.3.1 验证返回参数缺失错误
- [x] 6.4 编写测试用例：云 API 调用失败
  - [x] 6.4.1 Mock IdentifyZone API 返回错误
  - [x] 6.4.2 验证错误信息正确返回

## 7. 代码验证

- [ ] 7.1 检查代码格式是否符合 Go 标准格式（收尾阶段由 gofmt 处理）
- [x] 7.2 确保所有导出的函数和类型都有注释

## 8. 文档生成

- [x] 8.1 确保资源代码中的文档注释完整
- [ ] 8.2 准备运行 make doc 生成文档（在收尾阶段执行）
