## 1. 资源代码实现

- [x] 1.1 创建 resource_tc_teo_dns_record_v10.go 文件，定义资源 Schema（zone_id, name, type, content, location, ttl, weight, priority, status, created_on, modified_on 等字段）
- [x] 1.2 实现 resourceTencentCloudTeoDnsRecordV10Create 函数，调用 CreateDnsRecord API 并轮询等待记录生效
- [x] 1.3 实现 resourceTencentCloudTeoDnsRecordV10Read 函数，调用 DescribeDnsRecords API 查询记录详情
- [x] 1.4 实现 resourceTencentCloudTeoDnsRecordV10Update 函数，调用 ModifyDnsRecords API 修改记录并轮询等待生效
- [x] 1.5 实现 resourceTencentCloudTeoDnsRecordV10Delete 函数，调用 DeleteDnsRecords API 删除记录
- [x] 1.6 添加 Timeouts 配置块到 Schema（Create: 20分钟, Read: 3分钟, Update: 20分钟, Delete: 20分钟）

## 2. 服务层实现

- [x] 2.1 在 service_tencentcloud_teo.go 中添加 DescribeDnsRecordById 方法，通过 record_id 查询单条 DNS 记录
- [x] 2.2 实现轮询逻辑，在 Create 和 Update 后等待记录生效

## 3. 单元测试

- [x] 3.1 创建 resource_tc_teo_dns_record_v10_test.go 文件
- [x] 3.2 实现 TestAccTencentCloudTeoDnsRecordV10Basic 测试用例，测试基本的创建和读取操作
- [x] 3.3 实现 TestAccTencentCloudTeoDnsRecordV10Update 测试用例，测试更新操作
- [x] 3.4 实现 TestAccTencentCloudTeoDnsRecordV10Import 测试用例，测试导入操作
- [x] 3.5 实现 TestAccTencentCloudTeoDnsRecordV10Delete 测试用例，测试删除操作

## 4. 资源注册

- [x] 4.1 在 tencentcloud/provider.go 中注册新资源 ResourceTencentCloudTeoDnsRecordV10
- [x] 4.2 更新 tencentcloud/services/teo/tencentcloud_teo_suite_test.go，添加测试配置（如需要）

## 5. 文档示例

- [x] 5.1 创建 resource_tc_teo_dns_record_v10.md 示例文件，包含资源使用示例
- [ ] 5.2 运行 `make doc` 命令自动生成 website/docs/r/teo_dns_record_v10.md 文档

## 6. 代码验证

- [x] 6.1 运行 go build 确保代码可编译
- [ ] 6.2 运行 gofmt 格式化代码（在收尾阶段执行）
- [ ] 6.3 运行单元测试确保测试通过（TF_ACC=1）
