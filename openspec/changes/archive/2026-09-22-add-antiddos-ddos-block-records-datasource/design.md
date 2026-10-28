## Context

本次变更为 Terraform Provider for TencentCloud 新增 AntiDDoS（DDoS 防护）数据源 `tencentcloud_antiddos_ddos_block_records`，用于查询 DDoS 封堵解封记录列表及解封配额信息，对应云 API 接口 `DescribeDDoSBlockRecords`。

该接口位于 SDK 包 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos/v20250903`，而现有 antiddos 服务使用的 `UseAntiddosClient()` 仅返回 v20200309 的 client（`tencentcloud/connectivity/client.go:1236`），v20200309 包中不存在 `DescribeDDoSBlockRecords`。因此需在 connectivity 层新增 v20250903 client 访问器。

数据源类型为 RESOURCE_KIND_DATASOURCE，只需实现 R（Read）接口，代码风格参照 `tencentcloud_igtm_instance_list` 数据源实现（含分页、retry 在分页循环内部、filters 输入、列表输出）。

## Goals / Non-Goals

**Goals:**
- 提供查询 AntiDDoS 封堵解封记录的 Terraform 数据源。
- 正确接入 v20250903 antiddos SDK 包，不破坏现有 v20200309 antiddos client。
- 实现自动分页、retry 在分页循环内部、空响应不清理 state 的健壮读取逻辑。
- 沉淀云 API struct 关键字段定义与参考实现结构，使实施阶段无需回读 vendor 大文件。

**Non-Goals:**
- 不实现封堵/解封操作（仅查询，无 CUD）。
- 不修改现有 antiddos 资源/数据源行为。
- 不新增 limit/offset 用户可见参数（内部自动分页）。

## Decisions

### Decision 1: 新增 `UseAntiddosV20250903Client()` 而非替换现有 client

现有 `antiddos import ... v20200309` 别名 `antiddos` 和字段 `antiddosConn *antiddos.Client` 已被大量 antiddos 服务代码使用。直接替换会造成大范围改动且破坏向后兼容。

**决定**：在 `tencentcloud/connectivity/client.go` 中：
- 新增 import：`antiddosv20250903 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos/v20250903"`（参照 igtm 使用 `igtmv20231024` 别名 + `UseIgtmV20231024Client()` 的模式）。
- 在 `TencentCloudClient` 结构体 `antiddosConn` 字段附近新增 `antiddosV20250903Conn *antiddosv20250903.Client`。
- 在 `UseAntiddosClient()` 方法之后新增 `UseAntiddosV20250903Client()` 方法，模式与 `UseAntiddosClient()` 一致（`me.NewModelProfile(300)` + `antiddosv20250903.NewClient` + `WithHttpTransport(&LogRoundTripper{})`）。

服务层调用：`me.client.UseAntiddosV20250903Client().DescribeDDoSBlockRecords(request)`。

### Decision 2: Schema 设计 — filters 嵌套块 + 平铺输出字段

参照 `data_source_tc_igtm_instance_list.go` 的 filters 模式：
- 入参：`start_time`（Required, string）、`end_time`（Required, string）、`filters`（Optional, TypeList，每项含 `name` Required string、`values` Required TypeSet string）。
- 出参（computed）：
  - `block_records`（TypeList）：每项含 `resource`、`block_time`、`status`（均 string, computed）。
  - `unblock_quota_info`（TypeList，单元素）：含 `total_quota`、`used_quota`（TypeInt）、`quota_start_time`、`quota_end_time`（string）。
  - `result_output_file`（Optional, string）。

注意：响应中 `UnblockQuotaInfo` 是单个对象（非数组），在 Terraform schema 中以 `TypeList`（最多 1 个元素）表达，便于 `d.Set`。

### Decision 3: Read 函数结构 — 参照 igtm_instance_list

Read 函数骨架：
1. `defer tccommon.LogElapsed(...)` + `defer tccommon.InconsistentCheck(d, meta)()`。
2. 构造 `logId` / `ctx`，`service := AntiddosService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}`。
3. 组装 `paramMap`：`StartTime`、`EndTime`（必填）、`Filters`（可选，转换为 `[]*antiddosv20250903.Filter`）。
4. `resource.Retry(tccommon.ReadRetryTimeout, ...)` 调用 service 方法，获取 `blockRecords []*DDoSBlockRecord`、`unblockQuotaInfo *DDoSUnblockQuota`。
5. 映射到 schema（每个字段先判 nil 再 set），`d.SetId(helper.BuildToken())`，处理 `result_output_file`。

### Decision 4: Service 层方法 — 分页 + retry 在循环内部

方法签名：
```go
func (me *AntiddosService) DescribeAntiddosDDoSBlockRecordsByFilter(ctx context.Context, param map[string]interface{}) (blockRecords []*antiddosv20250903.DDoSBlockRecord, unblockQuotaInfo *antiddosv20250903.DDoSUnblockQuota, errRet error)
```

实现要点（参照 `DescribeIgtmInstanceListByFilter`，service_tencentcloud_igtm.go:398-460）：
- `request = antiddosv20250903.NewDescribeDDoSBlockRecordsRequest()`。
- 从 paramMap 解析 `StartTime`、`EndTime`、`Filters`。
- 分页循环：`limit uint64 = 100`（API 注释标注最大值 100），`offset` 递增。
- **retry 在 for 循环内部**：`resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError { ratelimit.Check(request.GetAction()); result, e := me.client.UseAntiddosV20250903Client().DescribeDDoSBlockRecords(request); ... })`。
- 空响应处理：当 `result == nil || result.Response == nil` 时返回 `resource.NonRetryableError(...)`（不 SetId）。
- 累积 `result.Response.BlockRecords`，每页 fetch 后判断 `len(BlockRecords) < int(limit)` 退出。
- 用循环外变量保存 `unblockQuotaInfo = response.Response.UnblockQuotaInfo`（每个分页响应均带同一配额信息，取最后一页即可，或首次非空即取）。

### Decision 5: 测试 — gomonkey mock 云 API

按项目规范，新增资源的数据源测试使用 gomonkey mock 云 API（不使用 TF_ACC 测试套件），只测业务逻辑。对 `UseAntiddosV20250903Client().DescribeDDoSBlockRecords` 进行 mock，验证 Read 函数对响应字段的映射、空响应处理、分页累积逻辑。

## Risks / Trade-offs

- **[v20250903 与 v20200309 共存]** → 风险低。两个 SDK 包独立，import 别名区分（`antiddos` vs `antiddosv20250903`），connectivity 层各自独立 client 字段与方法，互不影响。
- **[UnblockQuotaInfo 作为 TypeList 表达单对象]** → 若 API 返回 nil，不 set 即可，schema computed 字段不 set 不会报错。可接受。
- **[分页累积大结果集内存]** → 封堵记录为历史查询数据，单次查询时间范围 ≤31 天（API 约束），结果量有限，内存可接受。
- **[空响应不清理 state]** → 按 DATASOURCE 规范要求，返回 NonRetryableError 让重试耗尽失败，避免误清 id。属增强健壮性的正确做法。

## Migration Plan

纯新增数据源，无需迁移。回滚仅需移除 provider 注册与新增文件。

## Reference: Cloud API Struct 关键字段定义（摘自 vendor .../antiddos/v20250903/models.go）

以下为实施阶段所需的 struct 字段定义，无需回读 vendor 大文件。

### DDoSBlockRecord (models.go:23-32)
```go
type DDoSBlockRecord struct {
    // 被封堵的资源，公网 IP，示例：117.175.94.231
    Resource *string `json:"Resource,omitnil,omitempty" name:"Resource"`
    // 被封堵的时间
    BlockTime *string `json:"BlockTime,omitnil,omitempty" name:"BlockTime"`
    // 封堵解封状态。枚举：Blocked(已封堵)、Unblocking(解封中)、Unblocked(已解封)
    Status *string `json:"Status,omitnil,omitempty" name:"Status"`
}
```

### DDoSUnblockQuota (models.go:34-46)  ← 注意：出参 JsonPath `UnblockQuotaInfo` 对应此 struct
```go
type DDoSUnblockQuota struct {
    // 解封次数配额总数
    TotalQuota *uint64 `json:"TotalQuota,omitnil,omitempty" name:"TotalQuota"`
    // 已使用的配额总数
    UsedQuota *uint64 `json:"UsedQuota,omitnil,omitempty" name:"UsedQuota"`
    // 配额生效的起始时间
    QuotaStartTime *string `json:"QuotaStartTime,omitnil,omitempty" name:"QuotaStartTime"`
    // 配额生效的结束时间
    QuotaEndTime *string `json:"QuotaEndTime,omitnil,omitempty" name:"QuotaEndTime"`
}
```

### DescribeDDoSBlockRecordsRequest (models.go:66-83)
```go
type DescribeDDoSBlockRecordsRequest struct {
    *tchttp.BaseRequest
    // 查询的起始时间。参数格式：2026-02-04T11:30:00+08:00
    StartTime *string `json:"StartTime,omitnil,omitempty" name:"StartTime"`
    // 查询的结束时间。(EndTime - StartTime) 需 <= 31 天。参数格式：2026-03-04T11:30:00+08:00
    EndTime *string `json:"EndTime,omitnil,omitempty" name:"EndTime"`
    // 过滤条件，Filters.Values 上限为 20。不填写返回当前 appid 下所有被封堵过的资源列表。
    //   Resource: 按被封堵的 IP 或资源六段式过滤
    //   Status:   按封堵状态过滤 (Blocked/Unblocking/Unblocked)
    Filters []*Filter `json:"Filters,omitnil,omitempty" name:"Filters"`
    // 分页查询限制数，最大值为 100。默认值 20
    Limit *uint64 `json:"Limit,omitnil,omitempty" name:"Limit"`
    // 分页查询偏移量。默认值 0
    Offset *uint64 `json:"Offset,omitnil,omitempty" name:"Offset"`
}
```

### DescribeDDoSBlockRecordsResponse / ResponseParams (models.go:108-126)
```go
type DescribeDDoSBlockRecordsResponseParams struct {
    // 封堵解封记录总数
    TotalCount *uint64 `json:"TotalCount,omitnil,omitempty" name:"TotalCount"`
    // 封堵解封记录
    BlockRecords []*DDoSBlockRecord `json:"BlockRecords,omitnil,omitempty" name:"BlockRecords"`
    // 解封次数配额信息
    UnblockQuotaInfo *DDoSUnblockQuota `json:"UnblockQuotaInfo,omitnil,omitempty" name:"UnblockQuotaInfo"`
    // 唯一请求 ID
    RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeDDoSBlockRecordsResponse struct {
    *tchttp.BaseResponse
    Response *DescribeDDoSBlockRecordsResponseParams `json:"Response"`
}
```

### Filter (models.go:139-145)
```go
type Filter struct {
    // 需要过滤的字段
    Name *string `json:"Name,omitnil,omitempty" name:"Name"`
    // 字段的过滤值
    Values []*string `json:"Values,omitnil,omitempty" name:"Values"`
}
```
注意：v20250903 的 `Filter` 结构体字段名为 `Name`（不是 v20200309 的 `Key`），赋值时使用 `filter.Name = helper.String(v)`。

### client.go 构造函数 (client.go:48-71)
```go
func NewDescribeDDoSBlockRecordsRequest() (request *DescribeDDoSBlockRecordsRequest)
func NewDescribeDDoSBlockRecordsResponse() (response *DescribeDDoSBlockRecordsResponse)
func (c *Client) DescribeDDoSBlockRecords(request *DescribeDDoSBlockRecordsRequest) (response *DescribeDDoSBlockRecordsResponse, err error)
```
API 描述：查询封堵解封记录和解封配额信息。

## Reference: 代码风格参照 — tencentcloud_igtm_instance_list

### 数据源文件 data_source_tc_igtm_instance_list.go
- 包名与 import：`package igtm`，import `igtmv20231024 "...igtm/v20231024"`、`tccommon`、`helper`、`resource`、`schema`。
- 入口函数 `DataSourceTencentCloudIgtmInstanceList() *schema.Resource`，仅设 `Read`，Schema 含 filters(TypeList) + 列表输出(TypeList, computed) + `result_output_file`(Optional)。
- Read 函数：组 paramMap → `resource.Retry(ReadRetryTimeout, ...)` 调 service 方法 → 映射响应各字段（每个先判 nil）→ `d.SetId(helper.BuildToken())` → WriteToFile。

### service 层 DescribeIgtmInstanceListByFilter (service_tencentcloud_igtm.go:398-460)
- 分页：`limit uint64 = 100`，offset 递增的 for 循环。
- retry 在 for 循环内部：`resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError { ratelimit.Check(...); result, e := me.client.UseIgtmV20231024Client().DescribeInstanceList(request); ... })`。
- 空响应：`if result == nil || result.Response == nil || result.Response.InstanceSet == nil { return resource.NonRetryableError(...) }`。
- 累积：`ret = append(ret, response.Response.InstanceSet...)`，`len(...) < int(limit)` 时 break。

### 现有 antiddos 数据源风格
- `data_source_tc_antiddos_bgp_biz_trend.go` 使用 `antiddos "...v20200309"`、`UseAntiddosClient()`、 paramMap + `resource.Retry(ReadRetryTimeout, ...)` 模式。本次新增数据源沿用相同骨架，但 import 与 client 访问器改用 v20250903。