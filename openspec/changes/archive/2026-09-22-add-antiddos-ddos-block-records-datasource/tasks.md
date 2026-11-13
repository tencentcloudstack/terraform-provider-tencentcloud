## 1. Connectivity 层接入 v20250903 antiddos client

- [x] 1.1 在 `tencentcloud/connectivity/client.go` import 区新增 `antiddosv20250903 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos/v20250903"`（参照 igtm 的 `igtmv20231024` 别名模式）
- [x] 1.2 在 `TencentCloudClient` 结构体中 `antiddosConn` 字段（client.go:191）附近新增字段 `antiddosV20250903Conn *antiddosv20250903.Client`
- [x] 1.3 在 `UseAntiddosClient()` 方法（client.go:1236）之后新增 `UseAntiddosV20250903Client()` 方法，返回 `*antiddosv20250903.Client`，复用 `me.NewModelProfile(300)` + `antiddosv20250903.NewClient` + `WithHttpTransport(&LogRoundTripper{})` 模式

## 2. Service 层方法

- [x] 2.1 在 `tencentcloud/services/antiddos/service_tencentcloud_antiddos.go` 新增方法 `DescribeAntiddosDDoSBlockRecordsByFilter(ctx, param) (blockRecords []*antiddosv20250903.DDoSBlockRecord, unblockQuotaInfo *antiddosv20250903.DDoSUnblockQuota, errRet error)`
- [x] 2.2 在方法内使用 `antiddosv20250903.NewDescribeDDoSBlockRecordsRequest()`，从 paramMap 解析 StartTime/EndTime/Filters（Filters 转为 `[]*antiddosv20250903.Filter`，注意 v20250903 Filter 字段名为 `Name`）
- [x] 2.3 实现分页循环，limit=100（API 注释最大值），offset 递增；retry 必须放在 for 循环内部（`resource.Retry(tccommon.ReadRetryTimeout, ...)` + `ratelimit.Check` + `me.client.UseAntiddosV20250903Client().DescribeDDoSBlockRecords`）
- [x] 2.4 空响应处理：`result == nil || result.Response == nil` 时返回 `resource.NonRetryableError`，不调用 SetId；累积 BlockRecords，`len < int(limit)` 时退出；循环外保存 UnblockQuotaInfo

## 3. 数据源文件

- [x] 3.1 创建 `tencentcloud/services/antiddos/data_source_tc_antiddos_ddos_block_records.go`，入口函数 `DataSourceTencentCloudAntiddosDDoSBlockRecords() *schema.Resource`，仅设 Read
- [x] 3.2 定义 Schema：`start_time`(Required,string)、`end_time`(Required,string)、`filters`(Optional,TypeList，含 name Required/values Required TypeSet)、`block_records`(Computed,TypeList，含 resource/block_time/status)、`unblock_quota_info`(Computed,TypeList，含 total_quota/used_quota/quota_start_time/quota_end_time)、`result_output_file`(Optional,string)
- [x] 3.3 实现 Read 函数：参照 igtm_instance_list，组 paramMap（StartTime/EndTime 必填、Filters 可选）→ `resource.Retry(ReadRetryTimeout, ...)` 调 service 方法 → 映射响应（每字段先判 nil）→ `d.SetId(helper.BuildToken())` → WriteToFile
- [x] 3.4 import 区使用 `antiddosv20250903 "...antiddos/v20250903"`、`tccommon`、`helper`、`resource`、`schema`

## 4. Provider 注册

- [x] 4.1 在 `tencentcloud/provider.go` 的 DataSourcesMap 中 antiddos 数据源区（约 provider.go:1377-1381）新增 `"tencentcloud_antiddos_ddos_block_records": antiddos.DataSourceTencentCloudAntiddosDDoSblockRecords()`

## 5. 文档

- [x] 5.1 创建 `tencentcloud/services/antiddos/data_source_tc_antiddos_ddos_block_records.md`：一句话描述（带 AntiDDos 云产品名，"Use this data source to query ..."）+ Example Usage（含 start_time/end_time/filters 示例），不写 Argument/Attribute Reference（由 make doc 生成）

## 6. 单元测试

- [x] 6.1 创建 `tencentcloud/services/antiddos/data_source_tc_antiddos_ddos_block_records_test.go`，使用 gomonkey mock `UseAntiddosV20250903Client().DescribeDDoSBlockRecords`，测试 Read 函数对响应字段的映射、空响应处理、分页累积逻辑（不使用 TF_ACC 测试套件）

## 7. 验证（收尾阶段执行）

- [ ] 7.1 gofmt 格式化所有新增/修改的 Go 文件（由 tfpacer-finalize skill 执行）
- [ ] 7.2 通过 make doc 生成 website/docs/d/ 文档（由 tfpacer-finalize skill 执行）
- [ ] 7.3 在 .changelog/ 下生成 changelog 文件（由 tfpacer-finalize skill 执行）