## Context

腾讯云 Terraform Provider 现有 DTS 同步任务资源 `tencentcloud_dts_sync_job`、`tencentcloud_dts_sync_config` 等，其中 `tencentcloud_dts_sync_config` 通过 `ConfigureSyncJob` 接口在任务初始化阶段配置同步选项与对象。但任务运行过程中通过 `ModifySyncJobConfig` 接口进行的"动态配置修改"（DynamicObjects / DynamicOptions）尚无对应的 Terraform 资源管理。

本变更新增资源 `tencentcloud_dts_sync_job_config`，属于 RESOURCE_KIND_CONFIG 类型：

- 资源存在即配置存在，主要管理配置的读取与更新（RU 接口）。
- Read 通过 `DescribeSyncJobs` 读取任务详细配置。
- Update 通过 `ModifySyncJobConfig` 动态修改同步对象与同步选项。
- 资源 id = `job_id`（同步任务 ID，形如 `sync-werwfs23`）。

云 API 包路径：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206`，已在 vendor 中。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_dts_sync_job_config` 资源，提供 Read/Create/Update/Delete（Create 复用 Update，Delete 空实现）。
- Schema 覆盖 `ModifySyncJobConfig` 的 `DynamicObjects`（Objects 结构）与 `DynamicOptions`（DynamicOptions 结构）入参，以及 `DescribeSyncJobs` 返回的配置出参。
- 在 provider.go / provider.md 中注册新资源。
- 提供单元测试（gomonkey mock 云 API）与文档。

**Non-Goals:**
- 不替代现有 `tencentcloud_dts_sync_config`（ConfigureSyncJob）资源，二者职责不同。
- 不管理同步任务的创建/启停/隔离等生命周期操作（由 `tencentcloud_dts_sync_job` 及对应 operation 资源管理）。
- 不实现 `ModifySyncJobConfig` 入参以外的字段。

## Decisions

### 1. CRUD 模型：Create 复用 Update，Delete 空实现

参照同目录已有 CONFIG 资源 `resource_tc_dts_sync_config.go` 的实现模式：
- **Create**：读取 `job_id`，`d.SetId(jobId)`，然后直接 `return resourceTencentCloudDtsSyncJobConfigUpdate(d, meta)`。
- **Read**：调用 `DescribeDtsSyncJobConfigById`（封装 DescribeSyncJobs，request.JobId = jobId），取 `response.Response.JobList[0]`（SyncJobInfo），逐字段 set 回 schema。返回空（response/instance 为空）时先 `log.Printf("[CRUD] dts_sync_job_config id=%s", d.Id())` 再 `d.SetId("")`。
- **Update**：构建 `ModifySyncJobConfigRequest`（JobId 必填，DynamicObjects、DynamicOptions 可选），用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包裹调用，失败用 `tccommon.RetryError(e)` 包装。调用成功后无需异步轮询（ModifySyncJobConfig 为同步接口，响应仅 RequestId）。
- **Delete**：空实现（配置类资源，删除资源不清空云端配置）。

**理由**：与现有 CONFIG 资源保持一致，避免引入新不一致。ModifySyncJobConfig 响应仅含 RequestId，无任务标识返回，属同步接口，不需要异步轮询。

### 2. Schema 设计：DynamicObjects 与 DynamicOptions 两个顶层块

`ModifySyncJobConfig` 入参仅 3 个：`JobId`（必填）、`DynamicObjects`（可选）、`DynamicOptions`（可选）。因此顶层 Schema 为：

- `job_id`（Required, string, ForceNew）：同步任务 ID。
- `dynamic_objects`（Optional, TypeList, MaxItems:1）：对应 SDK `Objects` 结构。
- `dynamic_options`（Optional, TypeList, MaxItems:1）：对应 SDK `DynamicOptions` 结构。

**DynamicObjects**（Objects struct，vendor 行 6475）包含：`mode`(string)、`databases`([]*Database)、`advanced_objects`([]*string)、`online_ddl`(OnlineDDL{status})、`databases_op_filter`([]*DBOpFilter)。

**DynamicOptions**（DynamicOptions struct，vendor 行 4650）包含：`op_types`([]*string)、`ddl_options`([]*DdlOption)、`conflict_handle_type`(string)、`conflict_handle_option`(ConflictHandleOption)、`kafka_option`(KafkaOption)、`filter_begin_commit`(bool)、`filter_checkpoint`(bool)、`deal_of_exist_same_table`(string)、`start_position`(string)。

嵌套结构（Database / Table / Column / View / DdlOption / ConflictHandleOption / KafkaOption / TopicRule / KeyValuePairOption / OnlineDDL / DBOpFilter / OpFilter / TableFilter / ViewFilter）的 schema 组织严格参照 `resource_tc_dts_sync_config.go` 的既有 schema（objects / options 块），保持字段命名与类型一致，便于用户在两个资源间复用配置。

**字段必填性**：
- `job_id`：Required（ModifySyncJobConfig 入参 JobId 必填，标注为必填）。
- `dynamic_objects.online_ddl.status`：Required（需求标注必填）。
- `dynamic_options.op_types`：Required（需求标注必填）。
- `dynamic_options.ddl_options`：Required（需求标注必填，DynamicOptions 中 DdlOptions 必填）。
- `dynamic_objects.databases_op_filter.db_name`：Required。
- `dynamic_objects.databases_op_filter.tables.table_name`：Required。
- `dynamic_objects.databases_op_filter.views.view_name`：Required。
- 其余字段：Optional。

### 3. Read：从 SyncJobInfo 到 schema 的映射

`DescribeSyncJobs` 返回 `JobList []*SyncJobInfo`，取 `JobList[0]`。Read 阶段主要回填 `job_id`、`dynamic_objects`（来自 SyncJobInfo.Objects）、`dynamic_options`（回填 `op_types`、`conflict_handle_type`、`conflict_handle_option`、`ddl_options`、`kafka_option`、`filter_begin_commit`、`filter_checkpoint`、`deal_of_exist_same_table`、`start_position`，来自 SyncJobInfo.Options）。其它 SyncJobInfo 顶层只读字段（如 status、run_mode、src_info 等）按需求出参映射设为 Computed 顶层字段。

### 4. 服务层封装

在 `service_tencentcloud_dts.go` 新增：
- `DescribeDtsSyncJobConfigById(ctx, jobId) (*dts.SyncJobInfo, error)`：参照已有 `DescribeDtsSyncConfigById`（行 1082），request.JobId = jobId，返回 JobList[0]。
- `ModifyDtsSyncJobConfig(ctx, request) (*dts.ModifySyncJobConfigResponse, error)`：封装 ModifySyncJobConfig 调用，带 ratelimit.Check 与日志。

### 5. provider 注册与文档

- provider.go：在 `prov := tccommon.Provider{...}` 资源 map 中新增 `"tencentcloud_dts_sync_job_config": dts.ResourceTencentCloudDtsSyncJobConfig(),`，参照 `"tencentcloud_dts_sync_config"`（行 2370）相邻位置。
- provider.md：新增对应资源说明项（收尾阶段由 `make doc` 生成）。
- 文档 `resource_tc_tencentcloud_dts_sync_job_config.md`：一句话描述 + Example Usage + Import（CONFIG 类型支持 import，id 为 job_id）。

## Risks / Trade-offs

- **[DynamicOptions 字段语义与 Options 部分重叠]** → ModifySyncJobConfig 使用 DynamicOptions（仅含可热更新字段子集），Read 从 SyncJobInfo.Options 回填时需注意只回填 DynamicOptions 范围内的字段，避免与 sync_config 资源语义混淆。文档中明确说明本资源仅管理动态可变配置。
- **[ModifySyncJobConfig 为同步接口]** → 无需异步轮询，简化实现；但若云端实际生效存在延迟，Read 可能在更新后立即触发，需依赖 `defer tccommon.InconsistentCheck()` 兜底。
- **[Read 返回空判断]** → 严格遵守规则：response/instance 为空时先打印日志再 `d.SetId("")`。

## Appendix: 关键 vendor struct 字段定义（实施参考）

### DescribeSyncJobsRequest（vendor 行 4113）
入参：JobId *string, JobIds []*string, JobName *string, Order *string, OrderSeq *string, Offset *uint64, Limit *uint64, Status []*string, RunMode *string, JobType *string, PayMode *string, TagFilters []*TagFilter, SrcInfoPattern *string, DstInfoPattern *string。

### DescribeSyncJobsResponse（vendor 行 4247）
Response.JobList []*SyncJobInfo；Response.TotalCount *uint64；Response.RequestId *string。

### SyncJobInfo（vendor 行 8374）
JobId, JobName, PayMode, RunMode, ExpectRunTime, AllActions, Actions, Options *Options, Objects *Objects, Specification, ExpireTime, SrcRegion, SrcDatabaseType, SrcAccessType, SrcInfo *Endpoint, SrcNodeType, SrcInfos *SyncDBEndpointInfos, DstRegion, DstDatabaseType, DstAccessType, DstInfo *Endpoint, DstNodeType, DstInfos *SyncDBEndpointInfos, CreateTime, StartTime, Status, EndTime, Tags []*TagItem, Detail *SyncDetailInfo, TradeStatus, InstanceClass, AutoRenew, OfflineTime, OptObjStatus, AutoRetryTimeRangeMinutes, DumperResumeCtrl。

### ModifySyncJobConfigRequest（vendor 行 6316）
JobId *string（必填）, DynamicObjects *Objects（可选）, DynamicOptions *DynamicOptions（可选）。Response 仅 RequestId。

### DynamicOptions（vendor 行 4650）
OpTypes []*string（必填）, DdlOptions []*DdlOption（必填）, ConflictHandleType *string, ConflictHandleOption *ConflictHandleOption, KafkaOption *KafkaOption, FilterBeginCommit *bool, FilterCheckpoint *bool, DealOfExistSameTable *string, StartPosition *string。

### Objects（vendor 行 6475）
Mode *string, Databases []*Database, AdvancedObjects []*string, OnlineDDL *OnlineDDL, DatabasesOpFilter []*DBOpFilter。

### Database（vendor 行 1842）
DbName, NewDbName, DbMode, SchemaName, NewSchemaName, SchemaMode, TableMode, Tables []*Table, ViewMode, Views []*View, FunctionMode, Functions []*string, ProcedureMode, Procedures []*string, TriggerMode, Triggers []*string, EventMode, Events []*string。

### Table（vendor 行 8485）
TableName, NewTableName, FilterCondition, ColumnMode, Columns []*Column, TmpTables []*string, TableEditMode。

### Column（vendor 行 74）
ColumnName, NewColumnName。

### View（vendor 行 8605）
ViewName, NewViewName。

### DdlOption（vendor 行 1912）
DdlObject *string, DdlValue []*string。

### ConflictHandleOption（vendor 行 692）
ConditionColumn, ConditionOperator, ConditionOrderInSrcAndDst（均 *string）。

### KafkaOption（vendor 行 5128）
DataType *string, TopicType *string, DDLTopicName *string, TopicRules []*TopicRule, DataOption []*KeyValuePairOption。

### TopicRule（vendor 行 8546）
TopicName, PartitionType, DbMatchMode, DbName, TableMatchMode, TableName, Columns []*string。

### KeyValuePairOption（vendor 行 5145）
Key *string, Value *string。（DataOption 中的 key/value）

### OnlineDDL（vendor 行 6501）
Status *string（ON-启用, OFF-不启用）。

### DBOpFilter（vendor 行 1828）
DbName *string, OpFilter *OpFilter, Tables []*TableFilter, Views []*ViewFilter。

### OpFilter（vendor 行 6506）
OpTypes []*string, DdlOptions []*DdlOption。

### TableFilter（vendor 行 8508）
TableName *string, OpFilter *OpFilter。

### ViewFilter（vendor 行 8613）
ViewName *string, OpFilter *OpFilter。

### Endpoint（vendor 行 4679，src_info/dst_info 字段集）
Region, Role, DbKernel, InstanceId, Ip *string, Port *uint64, User, Password, DbName, VpcId, SubnetId, CvmInstanceId, UniqDcgId, UniqVpnGwId, CcnId, Supplier, EngineVersion, Account, AccountMode, AccountRole, RoleExternalId, TmpSecretId, TmpSecretKey, TmpToken, EncryptConn, DatabaseNetEnv, CcnOwnerUin, ChildInstanceId, ChildInstanceType, SetId。

### SyncDBEndpointInfos（vendor 行 8326，src_infos/dst_infos 字段集）
Region *string, AccessType *string, DatabaseType *string, Info []*Endpoint。

### SkipDdlOnError（vendor 行 7433，Options 下，读出参）
Tables []*string, ErrorCodes []*int64。

### RateLimitOption（vendor 行 6753，Options 下，读出参）
CurrentDumpThread, DefaultDumpThread *int64, CurrentDumpRps, DefaultDumpRps *int64, CurrentLoadThread, DefaultLoadThread *int64, CurrentLoadRps, DefaultLoadRps *int64, CurrentSinkerThread, DefaultSinkerThread *int64, HasUserSetRateLimit *string。

## Appendix: 参照资源结构（resource_tc_dts_sync_config.go）

### Schema 组织（参照文件：tencentcloud/services/dts/resource_tc_dts_sync_config.go）
- 顶层 `job_id`(Required, ForceNew)：行 27-31。
- `options` 块（行 45-200）：maxItems 1，包含 init_type、deal_of_exist_same_table、conflict_handle_type、add_additional_column、op_types(TypeSet)、conflict_handle_option(TypeList max1: condition_column/condition_operator/condition_order_in_src_and_dst)、ddl_options(TypeList: ddl_object/ddl_value TypeSet)、rate_limit_option(TypeList max1: 各 current/default 字段 + has_user_set_rate_limit)。
- `objects` 块（行 202-429）：maxItems 1，包含 mode、databases(TypeList: DbName/NewDbName/DbMode/SchemaName/NewSchemaName/SchemaMode/TableMode/Tables/ViewMode/Views/FunctionMode/Functions(TypeSet)/ProcedureMode/Procedures(TypeSet)/TriggerMode/Triggers(TypeSet)/EventMode/Events(TypeSet))、advanced_objects(TypeSet)、online_ddl(TypeList max1: status)。
- `src_info` / `dst_info` 块（行 457-743）：maxItems 1，包含 Endpoint 全字段集。
- `auto_retry_time_range_minutes`(Optional int)。

### CRUD 骨架
- Create（行 754-765）：set jobId → SetId → return Update。
- Read（行 767-）：调 service.DescribeDtsSyncConfigById；空则 log 且 SetId("")；逐字段 set。
- Update（行 1337-）：构建 request（含 retry 包裹 + InconsistentCheck 兜底）。
- Delete（行 1834-1839）：空实现。

### 服务层骨架（service_tencentcloud_dts.go 行 1082）
`DescribeDtsSyncConfigById(ctx, jobId)`：NewDescribeSyncJobsRequest → request.JobId = &jobId → DescribeSyncJobs → 返回 JobList[0]。