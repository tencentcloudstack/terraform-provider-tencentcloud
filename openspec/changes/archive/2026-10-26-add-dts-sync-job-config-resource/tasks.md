## 1. Service 层封装

- [x] 1.1 在 `tencentcloud/services/dts/service_tencentcloud_dts.go` 新增 `DescribeDtsSyncJobConfigById(ctx, jobId) (*dts.SyncJobInfo, error)`，封装 `DescribeSyncJobs`（request.JobId = jobId，返回 JobList[0]），参照已有 `DescribeDtsSyncConfigById`
- [x] 1.2 在 `tencentcloud/services/dts/service_tencentcloud_dts.go` 新增 `ModifyDtsSyncJobConfig(ctx, request) (*dts.ModifySyncJobConfigResponse, error)`，封装 `ModifySyncJobConfig` 调用（含 ratelimit.Check 与日志）

## 2. 资源 Schema 与 CRUD

- [x] 2.1 创建 `tencentcloud/services/dts/resource_tc_tencentcloud_dts_sync_job_config.go`，定义 `ResourceTencentCloudDtsSyncJobConfig()` 的 Schema：`job_id`(Required, ForceNew)、`dynamic_objects`(Optional, TypeList, MaxItems:1)、`dynamic_options`(Optional, TypeList, MaxItems:1)，以及 DescribeSyncJobs 出参对应的 Computed 字段
- [x] 2.2 实现 `dynamic_objects` 块嵌套 schema（Objects：mode/databases/advanced_objects/online_ddl{status Required}/databases_op_filter，参照 resource_tc_dts_sync_config.go 的 objects 块）
- [x] 2.3 实现 `dynamic_options` 块嵌套 schema（DynamicOptions：op_types[Required]/ddl_options[Required]/conflict_handle_type/conflict_handle_option/kafka_option/filter_begin_commit/filter_checkpoint/deal_of_exist_same_table/start_position，参照 resource_tc_dts_sync_config.go 的 options 块对应字段）
- [x] 2.4 实现 Computed 输出字段 schema（DescribeSyncJobs 出参：job_name/pay_mode/run_mode/expect_run_time/all_actions/actions/options/objects/specification/expire_time/src_region/.../status/tags/detail/trade_status 等，参照 SyncJobInfo struct）
- [x] 2.5 实现 `resourceTencentCloudDtsSyncJobConfigCreate`：读取 job_id → d.SetId(jobId) → return Update
- [x] 2.6 实现 `resourceTencentCloudDtsSyncJobConfigRead`：调 DescribeDtsSyncJobConfigById；空则 log.Printf("[CRUD] dts_sync_job_config id=%s", d.Id()) 再 d.SetId("")；非空逐字段 set（set 前判断 nil）
- [x] 2.7 实现 `resourceTencentCloudDtsSyncJobConfigUpdate`：构建 ModifySyncJobConfigRequest（JobId 必填、DynamicObjects/DynamicOptions 可选），用 resource.Retry(tccommon.WriteRetryTimeout,...) 包裹调用，失败 tccommon.RetryError(e)；成功后 return Read
- [x] 2.8 实现 `resourceTencentCloudDtsSyncJobConfigDelete`：空实现（配置类资源）

## 3. Provider 注册

- [x] 3.1 在 `tencentcloud/provider.go` 资源 map 中新增 `"tencentcloud_dts_sync_job_config": dts.ResourceTencentCloudDtsSyncJobConfig(),`（参照 `tencentcloud_dts_sync_config` 相邻位置）

## 4. 文档

- [x] 4.1 创建 `tencentcloud/services/dts/resource_tc_tencentcloud_dts_sync_job_config.md`：一句话描述（带云产品 DTS）+ Example Usage（含 dynamic_objects/dynamic_options）+ Import 部分（CONFIG 类型，import id 为 job_id）；不添加 Argument/Attribute Reference
- [x] 4.2 在 `tencentcloud/provider.md` 中新增 `tencentcloud_dts_sync_job_config` 资源说明项（收尾阶段由 make doc 生成，本任务仅在 provider.go 注册后由收尾流程处理）

## 5. 单元测试

- [x] 5.1 创建 `tencentcloud/services/dts/resource_tc_tencentcloud_dts_sync_job_config_test.go`，使用 gomonkey mock `DescribeSyncJobs` 与 `ModifySyncJobConfig` 云 API，覆盖 Create/Read/Update 业务逻辑（不使用 terraform 测试套件）

## 6. 验证与收尾

- [x] 6.1 代码正确性自查：核对 ModifySyncJobConfig 入参（JobId/DynamicObjects/DynamicOptions）与 DescribeSyncJobs 出参字段映射，确保 CRUD 各方法使用的字段在对应云 API 中存在
- [ ] 6.2 收尾阶段执行 tfpacer-finalize skill：gofmt 格式化、make doc 生成文档与 changelog（不在本阶段手动执行）
