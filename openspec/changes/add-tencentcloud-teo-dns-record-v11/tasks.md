## 1. 研究云API接口

- [x] 1.1 阅读 TEO SDK 文档，了解 CreateDnsRecord、DescribeDnsRecords、ModifyDnsRecords、DeleteDnsRecords 接口的参数和返回值
- [x] 1.2 确认哪些参数是必填、哪些是可选、哪些是计算字段
- [x] 1.3 确认接口是否为异步操作，是否需要轮询确认状态
- [x] 1.4 分析 vendor 目录下的 TEO SDK 代码，确认接口的具体实现

## 2. 实现资源 Schema 定义

- [x] 2.1 创建 `tencentcloud/services/teo/resource_tc_teo_dns_record_v11.go` 文件
- [x] 2.2 定义资源 ID 格式为 `zone_id#record_id`
- [x] 2.3 定义必填字段：zone_id, domain, record_type, record_value
- [x] 2.4 定义可选字段：ttl, priority, weight 等（根据云API确认）
- [x] 2.5 定义计算字段：record_id, status, created_at, updated_at
- [x] 2.6 定义 Timeouts 块，设置 create/update/delete 的默认超时时间为 10 分钟

## 3. 实现 Create 函数

- [x] 3.1 实现 `resourceTencentCloudTeoDnsRecordV11Create` 函数
- [x] 3.2 构造 CreateDnsRecord API 请求参数
- [x] 3.3 调用 TEO SDK 的 CreateDnsRecord 接口
- [x] 3.4 使用 helper.Retry() 轮询调用 DescribeDnsRecords 接口，直到记录创建成功
- [x] 3.5 设置资源 ID 为 `zone_id#record_id` 格式
- [x] 3.6 添加错误处理：defer tccommon.LogElapsed(), defer tccommon.InconsistentCheck()

## 4. 实现 Read 函数

- [x] 4.1 实现 `resourceTencentCloudTeoDnsRecordV11Read` 函数
- [x] 4.2 解析资源 ID，提取 zone_id 和 record_id
- [x] 4.3 构造 DescribeDnsRecords API 请求参数
- [x] 4.4 调用 TEO SDK 的 DescribeDnsRecords 接口
- [x] 4.5 从返回结果中匹配对应的 DNS 记录
- [x] 4.6 更新 Terraform 状态，填充所有字段值
- [x] 4.7 添加错误处理：defer tccommon.LogElapsed(), defer tccommon.InconsistentCheck()
- [x] 4.8 处理资源不存在的情况（返回 nil 表示删除）

## 5. 实现 Update 函数

- [x] 5.1 实现 `resourceTencentCloudTeoDnsRecordV11Update` 函数
- [x] 5.2 构造 ModifyDnsRecords API 请求参数
- [x] 5.3 只更新有变化的字段
- [x] 5.4 调用 TEO SDK 的 ModifyDnsRecords 接口
- [x] 5.5 使用 helper.Retry() 轮询调用 DescribeDnsRecords 接口，直到更新生效
- [x] 5.6 添加错误处理：defer tccommon.LogElapsed(), defer tccommon.InconsistentCheck()

## 6. 实现 Delete 函数

- [x] 6.1 实现 `resourceTencentCloudTeoDnsRecordV11Delete` 函数
- [x] 6.2 解析资源 ID，提取 zone_id 和 record_id
- [x] 6.3 构造 DeleteDnsRecords API 请求参数
- [x] 6.4 调用 TEO SDK 的 DeleteDnsRecords 接口
- [x] 6.5 使用 helper.Retry() 轮询调用 DescribeDnsRecords 接口，直到记录删除成功
- [x] 6.6 添加错误处理：defer tccommon.LogElapsed(), defer tccommon.InconsistentCheck()
- [x] 6.7 处理资源已删除的情况（幂等操作，返回成功）

## 7. 实现 Service 层辅助函数

- [x] 7.1 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中添加辅助函数（如需要）
- [x] 7.2 实现 API 请求的构造函数
- [x] 7.3 实现 API 响应的解析函数
- [x] 7.4 实现 ID 解析和构造函数（parseResourceId, buildResourceId）

## 8. 注册资源

- [x] 8.1 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中注册新资源 `tencentcloud_teo_dns_record_v11`
- [x] 8.2 确认资源注册位置正确，与其他 TEO 资源放在一起

## 9. 编写单元测试

- [x] 9.1 创建 `tencentcloud/services/teo/resource_tc_teo_dns_record_v11_test.go` 文件
- [x] 9.2 使用 mock 方法模拟 TEO SDK API 调用
- [x] 9.3 编写 Create 操作的单元测试用例
- [x] 9.4 编写 Read 操作的单元测试用例
- [x] 9.5 编写 Update 操作的单元测试用例
- [x] 9.6 编写 Delete 操作的单元测试用例
- [x] 9.7 编写错误处理的单元测试用例
- [x] 9.8 确保所有测试用例都能通过

## 10. 更新资源示例文件

- [x] 10.1 创建 `tencentcloud/services/teo/resource_tc_teo_dns_record_v11.md` 文件
- [x] 10.2 编写资源的 Terraform 使用示例
- [x] 10.3 说明必填和可选参数的用法
- [x] 10.4 提供 Create、Update、Delete 的完整示例

## 11. 生成文档

- [x] 11.1 运行 `make doc` 命令生成 `website/docs/r/teo_dns_record_v11.html.markdown` 文档
- [x] 11.2 确认文档内容完整，包含所有参数说明
- [x] 11.3 确认文档包含使用示例

## 12. 验证

- [x] 12.1 执行 `gofmt` 格式化代码
- [x] 12.2 检查代码编译是否通过
- [x] 12.3 运行单元测试（不执行 TF_ACC 验收测试）
- [x] 12.4 检查生成的文档是否正确
