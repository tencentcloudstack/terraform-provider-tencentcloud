## 1. 服务层实现

- [x] 1.1 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中新增方法 `DescribeTeoAvailableOriginACLFamilyByFilter(ctx context.Context, param map[string]interface{}) (ret []*teov20220901.OriginACLFamilyInfo, totalCount *int64, errRet error)`
  - 构造 `teov20220901.NewDescribeAvailableOriginACLFamilyRequest()`
  - 解析 paramMap：`ZoneId`→request.ZoneId（*string）；`Filters`→request.Filters（[]*teov20220901.Filter）
  - `ratelimit.Check(request.GetAction())`
  - 分页循环：`offset=0`、`limit=100`（API 标注最大值），每页在 `resource.Retry(tccommon.ReadRetryTimeout,...)` 内调用 `me.client.UseTeoClient().DescribeAvailableOriginACLFamily(request)`，错误用 `tccommon.RetryError(e)` 包装
  - `response==nil || response.Response==nil` 时返回 `resource.NonRetryableError`
  - `len(OriginACLFamilyInfos)<1` break；本页长度<limit break；累加 ret；记录 totalCount；offset+=limit
  - defer 错误日志

## 2. 数据源实现

- [x] 2.1 新增 `tencentcloud/services/teo/data_source_tc_teo_available_origin_acl_family.go`
  - `DataSourceTencentCloudTeoAvailableOriginAclFamily()` 返回 `*schema.Resource`，仅 `Read`
  - schema 输入：`zone_id`(Required,TypeString)、`filters`(Optional,TypeList,name Required/values Required TypeSet)、`result_output_file`(Optional,TypeString)
  - schema 输出：`origin_acl_family_infos`(Computed,TypeList)，元素含 `version`(TypeString)、`active_time`(TypeString)、`entire_addresses`(TypeList 单元素，含 `i_pv4`/`i_pv6` TypeSet)、`origin_acl_family`(TypeString)
- [x] 2.2 实现 `dataSourceTencentCloudTeoAvailableOriginAclRead(d, meta)`
  - `defer tccommon.LogElapsed(...)` + `defer tccommon.InconsistentCheck(d, meta)`
  - 构造 paramMap（ZoneId、Filters，参考 data_source_tc_teo_zones.go 解析 Filter）
  - `resource.Retry(tccommon.ReadRetryTimeout,...)` 调用服务方法；错误用 `tccommon.RetryError(e)` 包装
  - retry 失败路径保留 `log.Printf("[DATASOURCE] read empty, skip SetId")` 提示（当服务方法返回空切片非 error 时正常流转）
  - 遍历结果映射到 `origin_acl_family_infos`，每个字段 set 前判 nil；`entire_addresses` 用单元素列表 set
  - `d.SetId(helper.BuildToken())`
  - `result_output_file` 通过 `tccommon.WriteToFile` 输出

## 3. 文档与注册

- [x] 3.1 新增 `tencentcloud/services/teo/data_source_tc_teo_available_origin_acl_family.md`：一句话描述（带 TEO 产品名）+ Example Usage（按 zone_id 查询、按 filters 过滤）；不写 Argument Reference/Attribute Reference；不写 Import（DATASOURCE 不需要）
- [x] 3.2 在 `tencentcloud/provider.go` 的 `DataSourcesMap` 中按字母序追加 `"tencentcloud_teo_available_origin_acl_family": teo.DataSourceTencentCloudTeoAvailableOriginAclFamily()`
- [x] 3.3 在 `tencentcloud/provider.md` 数据源列表的 TencentCloud EdgeOne(TEO) → Data Source 段追加 `tencentcloud_teo_available_origin_acl_family`

## 4. 测试

- [x] 4.1 新增 `tencentcloud/services/teo/data_source_tc_teo_available_origin_acl_family_test.go`，使用 gomonkey mock 云 API 进行业务逻辑单元测试（不使用 terraform 验收测试套件）
  - mock `TeoService.DescribeTeoAvailableOriginACLFamilyByFilter` 返回构造的 `[]*teov20220901.OriginACLFamilyInfo`
  - 验证 Read 函数对返回数据的解析、set 逻辑（origin_acl_family_infos 列表、entire_addresses 子结构）
  - 确保函数返回的 error 都被检查；必不出错函数用 `_ = func()` 忽略 error
- [x] 4.2 校验生成的 Go 代码在当前环境下可正确构建执行（不执行 go build/go test，仅保证代码正确性）

## 5. 收尾（由 tfpacer-finalize skill 执行，本阶段不操作）

- [ ] 5.1 gofmt 格式化变更的 Go 代码
- [ ] 5.2 通过 `make doc` 生成 website/docs/ 文档
- [ ] 5.3 生成 .changelog 文件