## 1. Service 层实现

- [x] 1.1 在 `tencentcloud/services/dbdc/service_tencentcloud_dbdc.go` 中新增 `DescribeDBCustomZonesByFilter` 方法，签名为 `func (me *DbdcService) DescribeDBCustomZonesByFilter(ctx context.Context, param map[string]interface{}) (ret []*dbdcv20201029.ZoneInfo, totalCount int64, errRet error)`
- [x] 1.2 在方法内使用 `dbdcv20201029.NewDescribeDBCustomZonesRequest()` 构造请求，通过 `me.client.UseDbdcV20201029Client().DescribeDBCustomZones(request)` 发起调用（API 无入参，无需参数映射）
- [x] 1.3 在 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 内调用 API，调用失败时使用 `tccommon.RetryError(e)` 包装错误
- [x] 1.4 在 retry 块内严格检查返回空：`result == nil || result.Response == nil || result.Response.ZoneSet == nil` 时，打印 `log.Printf("[DATASOURCE] read empty, skip SetId")` 并返回 `resource.NonRetryableError`，不调用 `d.SetId("")`
- [x] 1.5 使用 `defer` 在错误时打印 `[CRITAL]` 日志，包含 logId、action、request body、错误信息
- [x] 1.6 单次调用返回全部数据，无需分页循环（API 无分页参数）

## 2. 数据源实现

- [x] 2.1 创建 `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_zones.go` 文件
- [x] 2.2 定义 `DataSourceTencentCloudDbdcDbCustomZones()` 函数返回 `*schema.Resource`，设置 `Read: dataSourceTencentCloudDbdcDbCustomZonesRead`
- [x] 2.3 定义 Schema：`result_output_file`（Optional, String）与 `zone_set`（Computed, TypeList，元素为 Resource，含 `zone`、`zone_state` 两个 Computed String 字段）
- [x] 2.4 实现 `dataSourceTencentCloudDbdcDbCustomZonesRead` 函数，添加 `defer tccommon.LogElapsed("data_source.tencentcloud_dbdc_db_custom_zones.read")()` 与 `defer tccommon.InconsistentCheck(d, meta)()`
- [x] 2.5 使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 调用 `service.DescribeDBCustomZonesByFilter(ctx, paramMap)`，失败时用 `tccommon.RetryError(e)` 包装
- [x] 2.6 遍历返回的 `[]*dbdcv20201029.ZoneInfo`，对每个 `Zone`、`ZoneState` 字段先判断 nil 再写入 map（遵循"先判 nil 再 set"要求）
- [x] 2.7 调用 `_ = d.Set("zone_set", zoneSetList)` 设置列表
- [x] 2.8 调用 `d.SetId(helper.BuildToken())` 设置数据源 id
- [x] 2.9 支持 `result_output_file`，当用户设置时调用 `tccommon.WriteToFile` 保存结果

## 3. Provider 注册

- [x] 3.1 在 `tencentcloud/provider.go` 的 DataSourcesMap 中注册 `tencentcloud_dbdc_db_custom_zones` 指向 `dbdc.DataSourceTencentCloudDbdcDbCustomZones()`，置于现有 dbdc 数据源注册项附近
- [x] 3.2 在 `tencentcloud/provider.md` 的数据源列表中添加 `tencentcloud_dbdc_db_custom_zones`，置于现有 dbdc 数据源项附近并保持字母序

## 4. 文档模板

- [x] 4.1 创建 `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_zones.md` 文件
- [x] 4.2 添加一句话描述并带云产品名称：`Use this data source to query available zones and their sale status from the TencentCloud DBDC product.`
- [x] 4.3 添加 Example Usage：`data "tencentcloud_dbdc_db_custom_zones" "example" {}`
- [x] 4.4 不手写 Argument Reference / Attribute Reference（由工具自动生成），不添加 Import 部分

## 5. 单元测试

- [x] 5.1 创建 `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_zones_test.go` 文件
- [x] 5.2 使用 gomonkey mock `DbdcService.DescribeDBCustomZonesByFilter` 方法，构造返回 `[]*dbdcv20201029.ZoneInfo`（含 `Zone`、`ZoneState` 字段值）
- [x] 5.3 编写测试用例验证 `zone_set` 列表字段映射正确（zone / zone_state）
- [x] 5.4 编写测试用例验证空返回场景下 Read 方法的行为（不调用 `d.SetId("")`，返回错误）
- [x] 5.5 确保测试代码可正确构建（禁止通过 `go test` 实际执行测试用例）

## 6. 代码正确性检查

- [x] 6.1 核对 `DescribeDBCustomZones` API 请求无入参，响应 `ZoneSet` 为 `[]*ZoneInfo`，`ZoneInfo` 含 `Zone`、`ZoneState` 字段，与 schema 字段一致
- [x] 6.2 检查所有函数返回的 error 均被处理（必定不出错的用 `_ =` 赋值）
- [x] 6.3 检查资源名称在日志/错误信息中统一使用 `dbdc_db_custom_zones`（小写蛇形命名），不使用模糊措辞
- [x] 6.4 检查数据源 Read 方法在 retry 失败路径不调用 `d.SetId("")`，保留 id 信息

## 7. 文档生成（收尾阶段）

- [ ] 7.1 在收尾阶段通过 `make doc` 命令生成 `website/docs/d/dbdc_db_custom_zones.html.markdown` 文档（禁止手动编写 website/ 下的文件）
- [ ] 7.2 在收尾阶段通过 `gofmt` 格式化变更的 Go 代码
- [ ] 7.3 在收尾阶段在 `.changelog/` 目录下创建 changelog 文件
