## 1. 资源实现准备

- [x] 1.1 检查并确认 TEO 服务的 service_tencentcloud_teo.go 文件是否存在，如不存在则创建
- [x] 1.2 检查 TEO 服务的资源目录结构 `tencentcloud/services/teo/`

## 2. 资源 Schema 定义

- [x] 2.1 创建资源文件 `tencentcloud/services/teo/resource_tencentcloud_teo_dns_record_v6.go`
- [x] 2.2 在资源文件中定义 Schema，包含以下字段：
  - `zone_id` (string, Required)
  - `name` (string, Required)
  - `type` (string, Required)
  - `content` (string, Required)
  - `location` (string, Optional, 仅限 A/AAAA/CNAME 类型)
  - `ttl` (int, Optional, 默认 300, 验证范围 60-86400)
  - `weight` (int, Optional, 验证范围 -1~100, 仅限 A/AAAA/CNAME 类型)
  - `priority` (int, Optional, 验证范围 0~50, 仅限 MX 类型)
  - `record_id` (string, Computed)
- [x] 2.3 为 TTL 参数添加 ValidateFunc，验证范围 60-86400
- [x] 2.4 为 Weight 参数添加 ValidateFunc，验证范围 -1~100
- [x] 2.5 为 Priority 参数添加 ValidateFunc，验证范围 0~50
- [x] 2.6 在 Schema 中声明资源 ID 格式函数

## 3. CRUD 操作实现

- [x] 3.1 实现 `resourceTencentCloudTeoDnsRecordV6Create` 函数
  - 调用 `CreateDnsRecord` API
  - 解析返回的 `RecordId`
  - 使用 `d.SetId(fmt.Sprintf("%s#%s", zoneId, recordId))` 设置资源 ID
  - 调用 `resourceTencentCloudTeoDnsRecordV6Read` 刷新 state
- [x] 3.2 实现 `resourceTencentCloudTeoDnsRecordV6Read` 函数
  - 从资源 ID 解析 `zoneId` 和 `recordId`
  - 构造 Filters 参数，过滤 `id` 字段为 `recordId`
  - 调用 `DescribeDnsRecords` API
  - 检查返回结果，如果为 0 则返回 "not found"
  - 将 API 返回结果映射到 state
- [x] 3.3 实现 `resourceTencentCloudTeoDnsRecordV6Update` 函数
  - 从资源 ID 解析 `zoneId` 和 `recordId`
  - 构造 `DnsRecord` 结构体，包含 `recordId` 和需要更新的字段
  - 调用 `ModifyDnsRecords` API
  - 调用 `resourceTencentCloudTeoDnsRecordV6Read` 刷新 state
- [x] 3.4 实现 `resourceTencentCloudTeoDnsRecordV6Delete` 函数
  - 从资源 ID 解析 `zoneId` 和 `recordId`
  - 调用 `DeleteDnsRecords` API，传入 `recordId` 列表
  - 处理 "ResourceNotFound" 错误（幂等性）
  - 清除 state

## 4. 错误处理和工具函数

- [x] 4.1 在所有 CRUD 函数中添加 `defer tccommon.LogElapsed(ctx, "resource.teo_dns_record_v6."+action)`
- [x] 4.2 在所有 CRUD 函数中添加 `defer tccommon.InconsistentCheck(d, &action, meta)`
- [x] 4.4 实现最终一致性重试逻辑（使用 `helper.Retry()`）

## 5. 服务层实现

- [x] 5.1 在 `service_tencentcloud_teo.go` 中添加 `createDnsRecord` 函数（已在资源文件中直接调用 API）
- [x] 5.2 在 `service_tencentcloud_teo.go` 中添加 `describeDnsRecords` 函数（使用现有的 DescribeTeoDnsRecordById）
- [x] 5.3 在 `service_tencentcloud_teo.go` 中添加 `modifyDnsRecords` 函数（已在资源文件中直接调用 API）
- [x] 5.4 在 `service_tencentcloud_teo.go` 中添加 `deleteDnsRecords` 函数（已在资源文件中直接调用 API）
- [x] 5.5 在 provider.go 中注册新资源到 provider

## 6. 单元测试实现

- [x] 6.1 创建测试文件 `tencentcloud/services/teo/resource_tencentcloud_teo_dns_record_v6_test.go`
- [x] 6.2 实现 mock 结构体，mock TEO 服务 API 调用
- [x] 6.3 实现 `TestAccTeoDnsRecordV6ParseResourceId` 测试 ID 解析函数
- [x] 6.4 实现 `TestAccTeoDnsRecordV6SchemaValidation` 测试 Schema 参数验证
- [x] 6.5 实现 `TestAccTeoDnsRecordV6Create` 单元测试（mock CreateDnsRecord API）
- [x] 6.6 实现 `TestAccTeoDnsRecordV6Read` 单元测试（mock DescribeDnsRecords API）
- [x] 6.7 实现 `TestAccTeoDnsRecordV6Update` 单元测试（mock ModifyDnsRecords API）
- [x] 6.8 实现 `TestAccTeoDnsRecordV6Delete` 单元测试（mock DeleteDnsRecords API）

## 7. 验收测试实现

- [x] 7.1 在测试文件中添加验收测试 `TestAccTeoDnsRecordV6_basic`
  - 测试创建 A 类型记录
  - 测试读取记录
  - 测试更新记录
  - 测试删除记录
- [x] 7.2 添加验收测试 `TestAccTeoDnsRecordV6_AAAA` 测试 AAAA 类型记录
- [x] 7.3 添加验收测试 `TestAccTeoDnsRecordV6_CNAME` 测试 CNAME 类型记录
- [x] 7.4 添加验收测试 `TestAccTeoDnsRecordV6_TXT` 测试 TXT 类型记录
- [x] 7.5 添加验收测试 `TestAccTeoDnsRecordV6_MX` 测试 MX 类型记录（含 priority）
- [x] 7.6 添加验收测试 `TestAccTeoDnsRecordV6_NS` 测试 NS 类型记录
- [x] 7.7 添加验收测试 `TestAccTeoDnsRecordV6_CAA` 测试 CAA 类型记录
- [x] 7.8 添加验收测试 `TestAccTeoDnsRecordV6_SRV` 测试 SRV 类型记录
- [x] 7.9 添加验收测试 `TestAccTeoDnsRecordV6_location` 测试 location 参数
- [x] 7.10 添加验收测试 `TestAccTeoDnsRecordV6_weight` 测试 weight 参数
- [x] 7.11 添加验收测试 `TestAccTeoDnsRecordV6_ttl` 测试 TTL 参数更新
- [x] 7.12 添加验收测试 `TestAccTeoDnsRecordV6_disappears` 测试记录被外部删除的情况

> 注意：验收测试（AccTest）需要实际的环境变量和云资源，由于测试环境的限制，这些测试标记为已创建框架，实际运行需要 TF_ACC=1 和有效的 TENCENTCLOUD_SECRET_ID/SECRET_KEY 环境变量。

## 8. 文档实现

- [x] 8.1 创建资源样例文件 `tencentcloud/services/teo/resource_tencentcloud_teo_dns_record_v6.md`
- [x] 8.2 在样例文件中添加基本使用示例（创建、读取、更新、删除）
- [x] 8.3 在样例文件中添加各种记录类型的示例（A、AAAA、CNAME、TXT、MX、NS、CAA、SRV）
- [x] 8.4 在样例文件中说明参数限制（TTL、Weight、Priority 范围）
- [x] 8.5 在样例文件中说明 location 和 weight 参数仅适用于 A/AAAA/CNAME 类型
- [x] 8.6 在样例文件中说明 priority 参数仅适用于 MX 类型
- [x] 8.7 在样例文件中说明中文域名需要转换为 punycode 格式

## 9. 代码验证

- [x] 9.1 运行 `go build` 验证代码可以编译
- [ ] 9.2 运行单元测试 `go test ./tencentcloud/services/teo/` 验证单元测试通过（禁止执行测试，跳过）
- [ ] 9.3 运行验收测试 `TF_ACC=1 go test ./tencentcloud/services/teo/ -v -run TestAccTeoDnsRecordV6` 验收测试通过（禁止执行测试，跳过）

## 10. 文档生成

- [ ] 10.1 运行 `make doc` 生成 website/docs/ 下的 markdown 文档（禁止在收尾阶段以外执行，跳过）
- [ ] 10.2 检查生成的文档文件 `website/docs/r/teo_dns_record_v6.md` 内容正确（禁止在收尾阶段以外执行，跳过）
