## 1. Data Source Implementation

- [x] 1.1 新建 `tencentcloud/services/tke/data_source_tc_tke_addon_values.go`（package `tke`，文件首行不加注释），实现 `DataSourceTencentCloudTkeAddonValues()`
- [x] 1.2 定义 schema：`cluster_id`(Required, TypeString)、`addon_name`(Required, TypeString)、`values`(Computed, TypeString)、`default_values`(Computed, TypeString)、`result_output_file`(Optional, TypeString，放在 schema 末尾)
- [x] 1.3 实现 `dataSourceTencentCloudTkeAddonValuesRead`：`defer tccommon.LogElapsed("data_source.tencentcloud_tke_addon_values.read")()` + `defer tccommon.InconsistentCheck(d, meta)()`，构造 `logId`/`ctx`/`TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}`
- [x] 1.4 在 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 内调用已有服务层方法 `service.DescribeAddonValuesList(ctx, clusterId, addonName)`，错误用 `tccommon.RetryError(e)` 包装，retry 失败直接 `return reqErr`
- [x] 1.5 判 nil 后 set：`resp != nil && resp.Values != nil` 才 `d.Set("values", *resp.Values)`；`resp != nil && resp.DefaultValues != nil` 才 `d.Set("default_values", *resp.DefaultValues)`；不调用 `d.SetId("")`
- [x] 1.6 `d.SetId(strings.Join([]string{clusterId, addonName}, tccommon.FILED_SP))`，随后处理 `result_output_file`（`tccommon.WriteToFile(output.(string), d)` 并检查 error）

## 2. Provider Registration

- [x] 2.1 在 `tencentcloud/provider.go` 的 `DataSourcesMap` 中新增 `"tencentcloud_tke_addon_values": tke.DataSourceTencentCloudTkeAddonValues(),`
- [x] 2.2 在 `tencentcloud/provider.md` 的 `Tencent Kubernetes Engine(TKE)` 节下 `Data Source` 列表末尾（约 729 行 `tencentcloud_kubernetes_cluster_available_extra_args` 之后）新增 `tencentcloud_tke_addon_values`

## 3. Documentation

- [x] 3.1 新建 `tencentcloud/services/tke/data_source_tc_tke_addon_values.md`：首行 `Use this data source to query detailed information of TKE addon values`
- [x] 3.2 `Example Usage` 使用 HCL，`data "tencentcloud_tke_addon_values" "example"` 填写 `cluster_id`/`addon_name`，并用 `jsondecode(data.tencentcloud_tke_addon_values.example.values)` 演示 JSON 字符串消费
- [x] 3.3 确认 md 中不含 `Import` / `Argument Reference` / `Attribute Reference` 段落（由 `make doc` 自动生成，收尾阶段执行）

## 4. Unit Tests

- [x] 4.1 新建 `tencentcloud/services/tke/data_source_tc_tke_addon_values_test.go`（`package tke_test`），实现 `tccommon.ProviderMeta` 的 mock（`GetAPIV3Conn()` 返回 `*connectivity.TencentCloudClient`），参照 `tencentcloud/services/cam/data_source_tc_cam_accounts_test.go`
- [x] 4.2 用 gomonkey `ApplyMethodFunc(&tke.TkeService{}, "DescribeAddonValuesList", ...)` 编写正常返回用例，断言 `values` / `default_values` / `d.Id()` 为 `cluster_id#addon_name`
- [x] 4.3 编写 `Values` 为 nil 的用例，断言不报错且 `default_values` 正常
- [x] 4.4 编写服务层返回 error 的用例，断言 Read 返回 error
- [x] 4.5 编写 schema 断言用例，校验 5 个字段类型与 Required/Computed 标记
- [x] 4.6 在测试文件注释中说明需 `-gcflags="all=-l"` 运行（不在本阶段执行测试）

## 5. Verification

- [x] 5.1 对照 `design.md` 的「云 API 参考」核对 request/response 字段名拼写（`ClusterId`/`AddonName`/`Values`/`DefaultValues`）与 `tccommon.FILED_SP` 常量拼写
- [x] 5.2 确认未新增/修改服务层方法（`DescribeAddonValuesList` 已存在）
- [x] 5.3 由收尾阶段统一执行 `gofmt`、`make doc` 与 changelog 生成，本阶段不执行
