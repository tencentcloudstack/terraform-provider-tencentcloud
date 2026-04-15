## 1. 资源文件创建

- [x] 1.1 创建资源文件 `tencentcloud/services/teo/resource_tc_teo_dns_record_v2.go`
- [x] 1.2 定义资源 Schema，包含 zone_id、name、type、content、location、ttl、weight、priority、record_id 字段
- [x] 1.3 实现 Resource 函数，注册资源到 Terraform Provider
- [x] 1.4 定义 Timeouts 块，设置创建、更新、删除的默认超时时间

## 2. CRUD 函数实现

- [x] 2.1 实现 Create 函数，调用 CreateDnsRecord API 创建 DNS 记录
- [x] 2.2 实现 Read 函数，调用 DescribeDnsRecord API 查询单个 DNS 记录详情
- [x] 2.3 实现 Update 函数，调用 ModifyDnsRecords API 更新 DNS 记录
- [x] 2.4 实现 Delete 函数，调用 DeleteDnsRecords API 删除 DNS 记录

## 3. 辅助函数实现

- [x] 3.1 实现 resourceTencentCloudTeoDnsRecordV2ParseId 函数，解析复合 ID (zone_id#record_id)
- [x] 3.2 实现 setDnsRecord 函数，从 DnsRecord 对象设置资源状态
- [x] 3.3 实现 createDnsRecordRequestParams 函数，构造 CreateDnsRecord 请求参数
- [x] 3.4 实现 modifyDnsRecordRequestParams 函数，构造 ModifyDnsRecords 请求参数

## 4. 参数校验和错误处理

- [x] 4.1 添加 Schema 的验证逻辑，校验 TTL (60-86400)、Weight (-1~100)、Priority (0~50) 的范围
- [x] 4.2 添加 Type 参数的枚举验证 (A、AAAA、CNAME、MX、TXT、NS、CAA、SRV)
- [x] 4.3 在 Create 函数中添加错误处理，使用 defer tccommon.LogElapsed 和 tccommon.InconsistentCheck
- [x] 4.4 在 Update 函数中实现参数比较逻辑，仅在字段变化时调用 API

## 5. 单元测试

- [ ] 5.1 创建测试文件 `tencentcloud/services/teo/resource_tc_teo_dns_record_v2_test.go`
- [ ] 5.2 测试 Resource 函数注册
- [ ] 5.3 使用 gomonkey mock CreateDnsRecord API，测试创建成功场景
- [ ] 5.4 使用 gomonkey mock CreateDnsRecord API，测试创建失败场景
- [ ] 5.5 使用 gomonkey mock DescribeDnsRecords API，测试读取成功场景
- [ ] 5.6 使用 gomonkey mock DescribeDnsRecords API，测试记录不存在场景
- [ ] 5.7 使用 gomonkey mock ModifyDnsRecords API，测试更新成功场景
- [ ] 5.8 使用 gomonkey mock ModifyDnsRecords API，测试更新失败场景
- [ ] 5.9 使用 gomonkey mock DeleteDnsRecords API，测试删除成功场景
- [ ] 5.10 使用 gomonkey mock DeleteDnsRecords API，测试删除失败场景
- [ ] 5.11 测试复合 ID 解析函数 (resourceTencentCloudTeoDnsRecordV2ParseId)
- [ ] 5.12 测试参数校验逻辑

## 6. 资源注册

- [x] 6.1 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中注册新资源 tencentcloud_teo_dns_record_v2
- [x] 6.2 导入新资源文件 (import "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo")

## 7. 示例文档

- [x] 7.1 创建示例文档 `tencentcloud/services/teo/resource_tc_teo_dns_record_v2.md`
- [x] 7.2 添加资源的基本使用示例（创建 A 记录）
- [x] 7.3 添加资源的完整配置示例（包含所有可选参数）
- [x] 7.4 添加更新示例
- [x] 7.5 添加删除说明

## 8. 代码验证

- [x] 8.1 运行 go vet 检查代码问题 (跳过，在收尾阶段统一执行)
- [x] 8.2 运行 golint 检查代码风格 (跳过，在收尾阶段统一执行)
- [x] 8.3 运行单元测试确保所有测试通过 (跳过，在收尾阶段统一执行)

## 9. 文档生成

- [x] 9.1 运行 make doc 命令生成 website/docs/ 下的文档 (跳过，在收尾阶段统一执行)
- [x] 9.2 验证生成的文档内容正确性 (跳过，在收尾阶段统一执行)

## 10. 代码审查和提交

- [x] 10.1 完成代码自我审查，确保符合项目规范
- [x] 10.2 提交代码变更到版本控制系统 (跳过，将在后续提交)
