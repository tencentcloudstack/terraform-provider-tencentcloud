## 1. 资源实现

- [x] 1.1 创建资源文件 `tencentcloud/services/teo/resource_tc_teo_dns_record_10.go`
- [x] 1.2 定义资源 Schema，包括必需参数（ZoneId, Name, Type, Content）和可选参数（Location, TTL, Weight, Priority）
- [x] 1.3 定义 computed 属性（RecordId, Status, CreatedOn）
- [x] 1.4 在 Schema 中添加 Timeouts 块以支持异步操作
- [x] 1.5 实现 Create 函数，调用 CreateDnsRecord API
- [x] 1.6 在 Create 函数中实现幂等性检查，使用 DescribeDnsRecords 检查记录是否已存在
- [x] 1.7 在 Create 函数中实现轮询逻辑，等待 DNS 记录生效
- [x] 1.8 实现 Read 函数，调用 DescribeDnsRecords API
- [x] 1.9 在 Read 函数中实现从复合 ID "zoneId#recordId" 中提取 ZoneId 和 RecordId
- [x] 1.10 在 Read 函数中实现处理记录不存在的情况
- [x] 1.11 实现 Update 函数，调用 ModifyDnsRecords API
- [x] 1.12 在 Update 函数中实现 no-op 检测，避免不必要的 API 调用
- [x] 1.13 在 Update 函数中实现轮询逻辑，等待 DNS 记录更新生效
- [x] 1.14 实现 Delete 函数，调用 DeleteDnsRecords API
- [x] 1.15 实现资源 ID 生成逻辑，使用 "zoneId#recordId" 格式
- [x] 1.16 在所有 CRUD 函数中添加错误处理和日志记录
- [x] 1.17 在 service_tencentcloud_teo.go 中注册新资源

## 2. 测试实现

- [x] 2.1 创建测试文件 `tencentcloud/services/teo/resource_tc_teo_dns_record_10_test.go`
- [x] 2.2 实现 CreateDnsRecord API 的 mock
- [x] 2.3 实现 DescribeDnsRecords API 的 mock
- [x] 2.4 实现 ModifyDnsRecords API 的 mock
- [x] 2.5 实现 DeleteDnsRecords API 的 mock
- [x] 2.6 编写 Create 函数的单元测试，包括成功创建和幂等性场景
- [x] 2.7 编写 Read 函数的单元测试，包括成功读取和记录不存在场景
- [x] 2.8 编写 Update 函数的单元测试，包括成功更新和 no-op 场景
- [x] 2.9 编写 Delete 函数的单元测试
- [x] 2.10 编写资源 ID 生成和解析的单元测试
- [x] 2.11 编写错误处理的单元测试

## 3. 文档和样例

- [x] 3.1 创建资源样例文件 `tencentcloud/services/teo/resource_tc_teo_dns_record_10.md`
- [x] 3.2 编写基本使用示例，包括创建 A 记录
- [x] 3.3 编写 MX 记录示例，展示 Priority 参数的使用
- [x] 3.4 编写使用可选参数（Location, TTL, Weight）的示例
- [x] 3.5 在样例文件中添加资源参数说明
- [x] 3.6 运行 `make doc` 命令生成 `website/docs/r/teo_dns_record_10.md` 文档

## 4. 验证和测试

- [ ] 4.1 运行单元测试验证功能正确性
- [ ] 4.2 运行 go vet 检查代码问题
- [ ] 4.3 验证资源可以正常导入和注册
- [ ] 4.4 检查生成的文档格式正确
- [ ] 4.5 验证复合 ID 格式正确且可解析
