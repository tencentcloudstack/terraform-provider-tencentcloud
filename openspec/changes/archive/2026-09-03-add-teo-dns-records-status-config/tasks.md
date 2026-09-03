## 1. Schema 定义

- [x] 1.1 在 `tencentcloud/services/teo/resource_tc_teo_dns_records_status_config.go` 中定义 `ResourceTencentCloudTeoDnsRecordsStatus()`，声明 Create/Read/Update/Delete 函数与 Importer（`schema.ImportStatePassthrough`）
- [x] 1.2 定义 schema 字段（仅包含 Modify 接口 `ModifyDnsRecordsStatus` 的参数字段）：`zone_id`（Required, ForceNew, TypeString）、`records_to_enable`（Optional, TypeList of TypeString）、`records_to_disable`（Optional, TypeList of TypeString）
- [x] 1.3 不暴露 `filters`、`sort_by`、`sort_order`、`match`（Describe 查询参数）、`dns_records`（响应数据）、`limit`/`offset` 分页参数
- [x] 1.4 不添加 `Timeouts` 块（无需异步轮询，与其他 config 资源一致）

## 2. Create 函数（复用 Update 逻辑）

- [x] 2.1 实现 `resourceTencentCloudTeoDnsRecordsStatusCreate`，使用 `defer tccommon.LogElapsed()` 和 `defer tccommon.InconsistentCheck()`
- [x] 2.2 从 `d.GetOk("zone_id")` 获取 `zoneId`
- [x] 2.3 调用 `d.SetId(zoneId)` 设置资源 ID
- [x] 2.4 直接 `return resourceTencentCloudTeoDnsRecordsStatusUpdate(d, meta)` 复用 Update 逻辑（与 `tencentcloud_teo_ddos_protection_config` 等其他 config 资源一致）

## 3. Read 函数

- [x] 3.1 实现 `resourceTencentCloudTeoDnsRecordsStatusRead`，构造 `teov20220901.NewDescribeDnsRecordsRequest()`
- [x] 3.2 仅填充 `ZoneId`，调用 `DescribeDnsRecords` 查询 DNS 记录列表
- [x] 3.3 使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包装 `DescribeDnsRecordsWithContext` 调用，失败用 `tccommon.RetryError(e)` 包装；retry 块内不执行 set ID 等操作
- [x] 3.4 若云 API 返回空（response nil / response.Response nil / len(DnsRecords)==0），先 `log.Printf("[CRUD] teo_dns_records_status id=%s", d.Id())` 保留现场，再 `d.SetId("")`
- [x] 3.5 Read 不回填任何列表数据到 schema（schema 不含 `dns_records` 字段）

## 4. Update 函数

- [x] 4.1 实现 `resourceTencentCloudTeoDnsRecordsStatusUpdate`，使用 `defer tccommon.LogElapsed()` 和 `defer tccommon.InconsistentCheck()`
- [x] 4.2 当 `records_to_enable` 或 `records_to_disable` 发生变化（`d.HasChange`）时，构造 `teov20220901.NewModifyDnsRecordsStatusRequest()`，填充 `ZoneId`、`RecordsToEnable`（若有，使用 `helper.InterfacesStrings` 转换）、`RecordsToDisable`（若有，使用 `helper.InterfacesStrings` 转换）
- [x] 4.3 使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包装 `ModifyDnsRecordsStatusWithContext` 调用，失败用 `tccommon.RetryError(e)` 包装；retry 块内仅调用接口
- [x] 4.4 Update 函数末尾调用 `resourceTencentCloudTeoDnsRecordsStatusRead(d, meta)` 回写最新状态，不额外轮询

## 5. Delete 函数

- [x] 5.1 实现 `resourceTencentCloudTeoDnsRecordsStatusDelete` 为 no-op，返回 nil，不调用云 API，不重置记录状态

## 6. Provider 注册

- [x] 6.1 在 `tencentcloud/provider.go` 的资源 map 中添加 `"tencentcloud_teo_dns_records_status": teo.ResourceTencentCloudTeoDnsRecordsStatus()`
- [x] 6.2 在 `tencentcloud/provider.md` 中添加对应资源注册条目

## 7. 文档

- [x] 7.1 新增 `tencentcloud/services/teo/resource_tc_teo_dns_records_status.md`，包含一句话描述（提及 TEO/EdgeOne）、Example Usage（展示 `zone_id`、`records_to_enable`/`records_to_disable`）、Import 部分（说明使用 `zone_id` 作为 ID）
- [x] 7.2 不添加 `Argument Reference` 和 `Attribute Reference` 部分（由工具自动生成）

## 8. 单元测试

- [x] 8.1 新增 `tencentcloud/services/teo/resource_tc_teo_dns_records_status_config_test.go`，使用 gomonkey mock 云 API（不使用 Terraform 测试套件）
- [x] 8.2 新增 Create 测试用例（含 `records_to_enable` 和 `records_to_disable` 两个分支），Create 经由 Update 调用 `ModifyDnsRecordsStatusWithContext` 和 `DescribeDnsRecordsWithContext`，验证 Modify 被调用且参数正确
- [x] 8.3 新增 Create 无 records 测试用例，Update 检测到无变化时跳过 Modify，仅调用 Read，验证 Modify 未被调用
- [x] 8.4 新增 Read 测试用例，mock `DescribeDnsRecordsWithContext` 返回 DNS 记录，验证资源保留在 state；返回空时验证 `d.SetId("")`
- [x] 8.5 新增 Update 测试用例，mock `ModifyDnsRecordsStatusWithContext` 和 `DescribeDnsRecordsWithContext`，验证 Modify 被调用且 Read 回填
- [x] 8.6 新增 API 错误测试用例（Create/Update），验证错误正确返回
- [x] 8.7 新增 Schema 测试用例，验证 schema 仅包含 `zone_id`、`records_to_enable`、`records_to_disable`，不包含 `filters`、`sort_by`、`sort_order`、`match`、`dns_records`、`limit`、`offset`
- [x] 8.8 保证生成的代码在当前环境下可正确构建执行（不执行 `go test`，仅保证代码正确性）

## 9. 验证（收尾阶段执行）

- [ ] 9.1 由收尾阶段 tfpacer-finalize skill 执行 `gofmt` 格式化
- [ ] 9.2 由收尾阶段 tfpacer-finalize skill 执行 `make doc` 生成 `website/docs/` 文档
- [ ] 9.3 由收尾阶段 tfpacer-finalize skill 生成 `.changelog/` 文件
