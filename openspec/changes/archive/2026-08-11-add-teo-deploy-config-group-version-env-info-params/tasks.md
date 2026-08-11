## 1. Service 层（service_tencentcloud_teo.go）

- [x] 1.1 在 `service_tencentcloud_teo.go` 新增 service 方法 `DescribeTeoEnvironmentsWithTotalCount(ctx context.Context, zoneId string) (totalCount uint64, envInfos []*teov20220901.EnvInfo, errRet error)`：构造 `NewDescribeEnvironmentsRequest()`，设置 `request.ZoneId = helper.String(zoneId)`，`ratelimit.Check(request.GetAction())`，调用 `me.client.UseTeoV20220901Client().DescribeEnvironments(request)`；校验 `response == nil || response.Response == nil` 后返回 `response.Response.TotalCount`（`uint64`）与 `response.Response.EnvInfos`；保留既有错误日志 defer。不修改既有 `DescribeTeoEnvironmentsByFilter`。
- [x] 1.2 校验新增 service 方法对 `DescribeEnvironments` 的字段访问与 vendor 中 `DescribeEnvironmentsResponseParams`（`TotalCount`、`EnvInfos`）一致，确保编译正确。

## 2. 资源 Schema（resource_tc_teo_deploy_config_group_version.go）

- [x] 2.1 在 `ResourceTencentCloudTeoDeployConfigGroupVersion` 的 `Schema` map 中新增 Computed 字段：`total_count`（`schema.TypeInt`）、`env_type`（`schema.TypeString`）、`scope`（`schema.TypeList`，`Elem: &schema.Schema{Type: schema.TypeString}`）、`env_create_time`（`schema.TypeString`）、`env_update_time`（`schema.TypeString`）。
- [x] 2.2 新增 `current_config_group_version_infos`（`schema.TypeSet`，`Computed: true`），`Elem` 为 `&schema.Resource{Schema: ...}`，其子字段均为 `schema.TypeString` 且 `Computed: true`：`version_id`、`version_number`、`source_version`、`group_type`、`group_id`、`description`、`status`、`create_time`。子字段命名与既有 `config_group_version_infos` 子字段保持风格一致，并补充 `source_version`。
- [x] 2.3 确保所有新增字段均为 `Computed: true`，不设 `Required`/`Optional`/`ForceNew`，保持向后兼容。

## 3. 资源 Read 方法（resource_tc_teo_deploy_config_group_version.go）

- [x] 3.1 在 `resourceTencentCloudTeoDeployConfigGroupVersionRead` 中，保留现有部署记录读取逻辑（`DescribeTeoDeployConfigVersionHistoryByFilter` 及其字段 set、空返回时 `d.SetId("")` 前打 `[CRUD]` 日志）不变。
- [x] 3.2 在部署记录读取成功（`len(respData) > 0`）后，使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 调用新增 service 方法 `service.DescribeTeoEnvironmentsWithTotalCount(ctx, zoneId)`；retry 块内仅调用接口并返回（`error` 用 `tccommon.RetryError` 包装），不做 set/定位。
- [x] 3.3 在 retry 块外：将返回的 `totalCount`（非 0 时）通过 `d.Set("total_count", totalCount)` 写入；遍历 `envInfos` 定位 `*item.EnvId == envId` 的元素。
- [x] 3.4 若未定位到目标环境（含 `envInfos` 为空）：打印 `log.Printf("[CRUD] teo deploy config group version env not found, zone_id=%s env_id=%s", zoneId, envId)`，跳过环境字段赋值，不 `d.SetId("")`，正常返回 nil。
- [x] 3.5 若定位到 `EnvInfo`：逐字段判 nil 后 set：`EnvType`→`env_type`；`Scope`（`[]*string`）转为 `[]interface{}` 后 `d.Set("scope", ...)`；`CreateTime`→`env_create_time`；`UpdateTime`→`env_update_time`。
- [x] 3.6 若 `CurrentConfigGroupVersionInfos` 非 nil：遍历构造 `[]map[string]interface{}`，每个子字段（`VersionId`/`VersionNumber`/`SourceVersion`/`GroupType`/`GroupId`/`Description`/`Status`/`CreateTime`）判 nil 后写入对应 key，最后 `d.Set("current_config_group_version_infos", ...)`。所有 `d.Set` 的 error 用 `_ =` 接收避免未使用变量。

## 4. 单元测试（resource_tc_teo_deploy_config_group_version_test.go）

- [x] 4.1 在既有测试文件中，使用 gomonkey mock `DescribeEnvironments` client 方法（与新增 service 方法对应），构造包含目标 `envId` 的 `EnvInfo`（含 `TotalCount`、`EnvType`、`Scope`、`CurrentConfigGroupVersionInfos` 各子字段、`CreateTime`、`UpdateTime`），验证 Read 后 state 中 `total_count`/`env_type`/`scope`/`current_config_group_version_infos`/`env_create_time`/`env_update_time` 被正确填充。
- [x] 4.2 补充场景：`EnvInfos` 中无目标 `envId` 时，Read 不报错、不清空 ID，且新增字段保持未设置。
- [x] 4.3 补充场景：`EnvInfo` 部分字段为 nil 时，对应字段被跳过 set 而不报错。
- [x] 4.4 确保新增测试用例可正确构建（不执行 `go test`，仅保证代码可编译）。

## 5. 文档

- [x] 5.1 更新 `tencentcloud/services/teo/resource_tc_teo_deploy_config_group_version.md`，在 Example Usage 中保持既有示例，并补充新增 Computed 字段说明由 `make doc` 自动生成 `Argument Reference`/`Attribute Reference`（不手动编写这两部分）。
- [x] 5.2 在收尾阶段通过 `make doc` 生成 `website/docs/r/teo_deploy_config_group_version.html.markdown`（不在本阶段手动编辑 website/ 目录）。

## 6. 验证（收尾阶段）

- [x] 6.1 执行 `gofmt` 格式化本次变更的 Go 文件（由 tfpacer-finalize skill 统一执行）。
- [x] 6.2 由收尾流程执行 `make doc` 生成文档，并生成 `.changelog/` 文件。
- [x] 6.3 由其他流程执行构建/检查（本阶段不执行 `go build`/`go vet`/`golint`）。
