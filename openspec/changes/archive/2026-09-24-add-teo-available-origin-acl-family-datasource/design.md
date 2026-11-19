## Context

EdgeOne (TEO) 源站防护（Origin ACL）通过"回源 IP 网段控制域"（OriginACLFamily）下发 ACL 规则。控制域包含标准控制域（gaz/mlc/emc）和精简控制域（plat-gaz/plat-mlc/plat-emc）等；不同控制域提供的回源 IP 网段数量与使用限制不同。在配置 `tencentcloud_teo_origin_acl` 资源时，`origin_acl_family` 字段需取自可用控制域列表。

当前 Provider 没有查询该列表的数据源。本次新增 `tencentcloud_teo_available_origin_acl_family` 数据源，封装云 API `DescribeAvailableOriginACLFamily`（teo v20220901），按站点 ID 查询可配置控制域及其版本、回源 IP 网段信息。

参考资源模式：`tencentcloud_igtm_instance_list`（DATASOURCE 类型，filters 输入 + list 输出 + 内部自动分页）。

## Goals / Non-Goals

**Goals:**
- 提供可查询 EdgeOne 源站防护可用控制域列表的 Terraform 数据源。
- 支持按 `zone_id` 查询；可选 `filters`（name/values）过滤。
- 内部自动分页拉取全部数据，不向用户暴露 Offset/Limit。
- 输出列表中每个元素的子结构（entire_addresses）平铺，符合规范化要求。
- 沉淀实施所需参考信息（vendor struct 定义、参考代码结构），避免实施阶段反复回读源文件。

**Non-Goals:**
- 不实现控制域的创建/修改/删除（云 API 仅为查询接口，控制域由平台维护）。
- 不暴露分页参数给用户。
- 不修改任何现有资源或数据源的 schema。

## Decisions

### D1: 数据源 schema 设计（参考 igtm_instance_list + teo_origin_acl）

数据源采用 DATASOURCE 标准结构，仅 `Read`。

输入参数：
- `zone_id`（Required, TypeString）：站点 ID。
- `filters`（Optional, TypeList）：过滤条件，元素结构含 `name`（Required, TypeString）、`values`（Required, TypeSet）。
- `result_output_file`（Optional, TypeString）：结果输出文件。

输出参数：
- `origin_acl_family_infos`（Computed, TypeList）：控制域信息列表，每个元素含：
  - `version`（Computed, TypeString）
  - `active_time`（Computed, TypeString）
  - `entire_addresses`（Computed, TypeList，最多一个元素）：含 `i_pv4`（TypeSet）与 `i_pv6`（TypeSet）
  - `origin_acl_family`（Computed, TypeString）

**说明**：`entire_addresses` 在 SDK 中类型为 `*Addresses`（结构体指针，非数组），按数据源惯例用 `TypeList` 包裹为单元素列表再 set。

**备选**：将 `i_pv4`/`i_pv6` 平铺到 `origin_acl_family_infos` 下而不保留 `entire_addresses` 层。**否决**：proposal 接口映射明确要求保留 `entire_addresses` 层（对应 `OriginACLFamilyInfo.EntireAddresses`），且与现有 `data_source_tc_teo_origin_acl.go` 中 `entire_addresses` 的处理一致。

### D2: 服务层方法设计（参考 DescribeTeoZonesByFilter 分页模式）

新增服务层方法：`DescribeTeoAvailableOriginACLFamilyByFilter(ctx, param) (ret []*teov20220901.OriginACLFamilyInfo, totalCount *int64, errRet error)`

实现要点：
- 内部 `for` 循环分页，`limit = 100`（API 标注的最大值），不暴露给用户。
- **重试逻辑放在分页循环内部**：每页调用用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包裹。
- 调用 `me.client.UseTeoClient().DescribeAvailableOriginACLFamily(request)`。
- response.Response.OriginACLFamilyInfos 为空时 break。
- 本页返回数 < limit 时 break。
- 累加 ret，记录 totalCount。
- defer 中记录错误日志。

**为什么重试在循环内**：与 `DescribeTeoZonesByFilter` 一致；re-auth/rate-limit 类瞬时错误应针对单次调用重试，而非整体重试导致已拉取页丢失。

### D3: Read 函数严格按 DATASOURCE 规范实现

- retry 块内调用服务方法；若服务方法返回空（response==nil / OriginACLFamilyInfos 长度 0），**不直接 d.SetId("")**，而是返回 `NonRetryableError`，让外层 retry 继续尝试，避免云 API 短暂波动清空 state。
- 在外层 retry 失败路径上保留 `log.Printf("[DATASOURCE] read empty, skip SetId")` 提示（当确实无数据时，服务方法返回空切片而非 error；此时正常走到 set 列表并 `d.SetId(helper.BuildToken())`）。
- 设置字段前先判 nil，nil 不调用 set。
- 不在 retry 块内执行 set id 等成功操作。
- `d.SetId(helper.BuildToken())`（数据源无真实 id，使用 token）。

**关于空数据处理**：`DescribeAvailableOriginACLFamily` 是查询接口，站点下确实可能无控制域（返回空列表）。此时不应视为错误。因此服务层在 response 为 nil 或列表空时返回空切片（非 error）；数据源 Read 层 retry 仅捕获真实 API 错误。对于"response.Response 为 nil"这种异常空，服务层返回 `NonRetryableError`（参考 DescribeTeoOriginAclByFilter 中 `result == nil || result.Response == nil` 的处理）。

### D4: Filter 结构体映射

`request.Filters` 类型为 `[]*teo.Filter`（teo v20220901 的 `Filter` 结构体，字段 `Name *string`、`Values []*string`）。解析方式参考 `data_source_tc_teo_zones.go` 中 `AdvancedFilter` 的解析逻辑，但改用 `teov20220901.Filter`。

### D5: 文档与注册

- `data_source_tc_teo_available_origin_acl_family.md`：一句话描述（带 TEO 产品名）+ Example Usage + Import（无，DATASOURCE 不需要）。不写 Argument Reference / Attribute Reference（由 `make doc` 生成）。格式参考 `data_source_tc_teo_origin_acl.md`。
- `provider.go`：在 `DataSourcesMap` 中追加注册，按字母序插入 teo 数据源段。
- `provider.md`：在数据源列表的 teo 段追加条目。

## Risks / Trade-offs

- **[空数据被误判为失败]** → 区分"列表为空"（正常，返回空切片）与"response 为 nil"（异常，返回 NonRetryableError），避免数据源对正常空结果报错。
- **[分页遗漏数据]** → 使用 API 标注最大 limit 100，并在本页长度 < limit 时终止，确保完整拉取。
- **[Filter 字段名混淆]** → 本数据源 Filter 的字段为 `Name`/`Values`（注意 igtm 用 `ResourceFilter` 的 `Name`/`Value`，不同服务 Filter 结构不同），须按 teo 的 `Filter` 类型构造。
- **[无 import 支持]** → DATASOURCE 类型资源不支持 import，文档不写 Import 部分。

## 参考信息（实施阶段直接使用，无需回读 vendor 大文件）

### 1. 云 API 接口与 request/response struct 关键字段

**包**：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`

**接口**：`DescribeAvailableOriginACLFamily`（非异步接口，仅查询）
- 客户端获取：`me.client.UseTeoClient().DescribeAvailableOriginACLFamily(request)`（`UseTeoClient` 返回 `*teo.Client`）
- 请求构造：`teov20220901.NewDescribeAvailableOriginACLFamilyRequest()`
- 响应构造：`teov20220901.NewDescribeAvailableOriginACLFamilyResponse()`

**Request struct**（vendor/models.go L8722-8736，字段同 Params L8708-8720）：
```go
type DescribeAvailableOriginACLFamilyRequest struct {
    *tchttp.BaseRequest
    ZoneId  *string  `json:"ZoneId,omitnil,omitempty"`   // 站点ID
    Filters []*Filter `json:"Filters,omitnil,omitempty"` // 过滤条件，Values 上限 20
    Offset  *uint64  `json:"Offset,omitnil,omitempty"`   // 分页偏移，默认0
    Limit   *uint64  `json:"Limit,omitnil,omitempty"`    // 分页限制，默认20，最大100
}
```

**Filter struct**（vendor/models.go L16970-16976）：
```go
type Filter struct {
    Name   *string  `json:"Name,omitnil,omitempty"`
    Values []*string `json:"Values,omitnil,omitempty"`
}
```
**支持的过滤字段**：`OriginACLFamily`（按控制域过滤，取值如 gaz/mlc/emc/plat-gaz/plat-mlc/plat-emc 等）。

**Response struct**（vendor/models.go L8761-8775）：
```go
type DescribeAvailableOriginACLFamilyResponseParams struct {
    TotalCount          *int64                `json:"TotalCount"`          // 总数
    OriginACLFamilyInfos []*OriginACLFamilyInfo `json:"OriginACLFamilyInfos"` // 控制域信息列表
    RequestId           *string               `json:"RequestId"`
}
type DescribeAvailableOriginACLFamilyResponse struct {
    *tchttp.BaseResponse
    Response *DescribeAvailableOriginACLFamilyResponseParams `json:"Response"`
}
```

**OriginACLFamilyInfo struct**（vendor/models.go L23626-23657）：
```go
type OriginACLFamilyInfo struct {
    Version        *string    `json:"Version"`        // 源站防护版本号（如 gaz-xxxxx）
    ActiveTime    *string    `json:"ActiveTime"`     // 版本生效时间 ISO 8601 UTC+8
    EntireAddresses *Addresses `json:"EntireAddresses"` // 回源 IP 网段详情（结构体指针）
    OriginACLFamily *string   `json:"OriginACLFamily"` // 控制域（gaz/mlc/emc/plat-*）
}
```

**Addresses struct**（vendor/models.go L380-386）：
```go
type Addresses struct {
    IPv4 []*string `json:"IPv4"` // IPv4 网段列表
    IPv6 []*string `json:"IPv6"` // IPv6 网段列表
}
```

### 2. 代码风格参照文件关键结构

**参照 1：data_source_tc_teo_zones.go（teo 数据源 + 分页 + AdvancedFilter 解析）**
- 位置：`tencentcloud/services/teo/data_source_tc_teo_zones.go`
- 关键结构：
  - `DataSourceTencentCloudTeoZones()` 返回 `*schema.Resource`，仅 `Read`。
  - `filters` schema：TypeList/Optional，elem Resource 含 `name`(Required)/`values`(Required, TypeSet)/`fuzzy`(Optional)。
  - Read 函数解析 filters 为 `[]*teov20220901.AdvancedFilter`，放入 `paramMap["Filters"]`。
  - `resource.Retry(tccommon.ReadRetryTimeout, ...)` 调 `service.DescribeTeoZonesByFilter`。
  - 遍历 respData 映射到 list，set 到 computed 列表字段。
  - `d.SetId(helper.DataResourceIdsHash(zoneIds))`。
  - `result_output_file` 通过 `tccommon.WriteToFile` 输出。

**参照 2：service_tencentcloud_teo.go 中 DescribeTeoZonesByFilter（分页服务方法）**
- 位置：`tencentcloud/services/teo/service_teo.go` L1739-1799
- 关键结构：
  ```go
  func (me *TeoService) DescribeTeoZonesByFilter(ctx context.Context, param map[string]interface{}) (ret []*teov20221701.Zone, errRet error) {
      request := teov20220901.NewDescribeZonesRequest()
      defer func() { if errRet != nil { log.Printf("[CRITAL]...") } }()
      for k, v := range param { /* 解析 Filters/Order/Direction */ }
      ratelimit.Check(request.GetAction())
      var offset int64 = 0
      var limit  int64 = 100
      for {
          request.Offset = &offset
          request.Limit  = &limit
          response := teo.NewDescribeZonesResponse()
          err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
              ratelimit.Check(request.GetAction())
              result, e := me.client.UseTeoClient().DescribeZones(request)
              if e != nil { return tccommon.RetryError(e) }
              response = result
              return nil
          })
          if err != nil { errRet = err; return }
          if response == nil || response.Response == nil || len(response.Response.Zones) < 1 { break }
          ret = append(ret, response.Response.Zones...)
          if len(response.Response.Zones) < int(limit) { break }
          offset += limit
      }
      return
  }
  ```
- 注意：服务层分页循环内已有 retry，数据源 Read 层不再重复 retry（仅捕获服务方法返回的 error）。本数据源参考此模式：服务层方法内部含分页 + retry，Read 层用 `resource.Retry` 包裹服务方法调用处理瞬时错误（与服务层 retry 不冲突，因 Read 层 retry 只会在服务方法整体返回 error 时重试）。

**参照 3：data_source_tc_teo_origin_acl.go（teo 数据源，set 时判 nil）**
- 位置：`tencentcloud/services/teo/data_source_tc_teo_origin_acl.go`
- 关键结构：每个字段 set 前判 nil；`entire_addresses` 用 `[]interface{}{entireAddressesMap}` 单元素列表 set；`i_pv4`/`i_pv6` 直接 set `[]*string`。

**参照 4：data_source_tc_igtm_instance_list.go（DATASOURCE 通用结构）**
- 位置：`tencentcloud/services/igtm/data_source_tc_igtm_instance_list.go`
- 关键结构：`defer tccommon.LogElapsed(...)` + `defer tccommon.InconsistentCheck(...)`；`d.SetId(helper.BuildToken())`；filters 解析；`paramMap` 传参。

### 3. provider.go / provider.md 注册位置

- `provider.go`：DataSourcesMap 中 teo 数据源段（当前在 L962-978 附近），新增 `"tencentcloud_teo_available_origin_acl_family": teo.DataSourceTencentCloudTeoAvailableOriginAclFamily()`。
- `provider.md`：数据源列表的 `### TencentCloud EdgeOne(TEO)` -> `Data Source` 段（当前在 L1598-1604 附近）追加 `tencentcloud_teo_available_origin_acl_family`。
- `service_tencentcloud_teo.go` 已导入：`teov20220901 "..."` 与 `teo "..."` 均指向同包；ratelimit、tccommon、helper 已导入。UseTeoClient 返回 `*teo.Client`，其方法 `DescribeAvailableOriginACLFamily` 存在。

### 4. 测试模式（gomonkey mock）

按规范，新增的 terraform 数据源使用 gomonkey 对云 API 进行 mock，仅做业务逻辑单元测试，不使用 terraform 验收测试套件。测试文件需导入 gomonkey，mock `TeoService` 的 `DescribeTeoAvailableOriginACLFamilyByFilter` 方法或 `connectivity.TencentCloudClient.UseTeoClient`，验证 Read 函数对返回数据的解析与 set 逻辑。