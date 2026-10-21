## Context

SSL 证书（SSL Certificate，产品代号 ssl）的"上传更新（证书 ID 不变）"能力由云 API `UploadUpdateCertificateInstance` 提供，可上传新的证书公私钥内容并联动更新已关联的云资源（如 clb）。当前 Terraform Provider 的 ssl 服务下仅有证书部署/回滚等 operation 类资源和若干数据源，缺少以"配置即资源"方式托管上传更新与查询更新记录的能力。

本资源类型为 `RESOURCE_KIND_CONFIG`：描述一个资源的配置，资源存在则配置存在，主要管理配置的更新和读取。资源文件名称格式 `resource_tc_<Product>_<资源名>_config.go`，即 `tencentcloud/services/ssl/resource_tc_ssl_upload_certificate_instance_config.go`。

云 API 包：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205`。客户端获取方式与现有 ssl 资源一致：`meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSSLCertificateClient()`。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_ssl_upload_certificate_instance` 资源（CONFIG 类型），CRU 接口分别为：Create/Update → `UploadUpdateCertificateInstance`，Read → `DescribeHostUploadUpdateRecord`（按 `OldCertificateId` 查询更新记录列表），并增加 `DescribeHostUploadUpdateRecordDetail`（按 `DeployRecordId` 查询更新记录详情）。
- 异步：`UploadUpdateCertificateInstance` 为异步接口，拿到 `DeployRecordId` 后先 `d.SetId()`，再轮询 Read 直到生效。
- 资源 id 采用 `DeployRecordId`。

**Non-Goals:**
- 不改造已有 ssl 数据源/资源与 schema。
- 不暴露 `limit/offset` 给用户（Read 内部按云 API 注释最大值分页拉取）。

## Decisions

### 1. 资源 id 选择与 schema 组织
- id = `DeployRecordId`（uint64 → 字符串）。
- 用户可配置入参：`old_certificate_id`（必填）、`resource_types`（必填）、`certificate_public_key`（必填）、`certificate_private_key`（必填，Sensitive）、`resource_types_regions`（可选 list）、`certificate_id`（可选）。
- 出参（Computed）：`deploy_record_id`、`deploy_status`、`update_sync_progress`（list，含 `resource_type`/`status` 及子 list `update_sync_progress_regions`，子 list 含 `region`/`total_count`/`offset_count`/`status`）。
- 查询更新记录列表出参：`deploy_record_list`（list，含 `id`/`old_cert_id`/`resource_types`/`status`/`create_time`/`update_time`）。

### 2. Create/Update 流程（异步）
配置型资源 Create 与 Update 均调用 `UploadUpdateCertificateInstance`，把 schema 入参映射到 request：
- `request.OldCertificateId`（必填）
- `request.ResourceTypes`（必填，TypeSet → []*string）
- `request.CertificatePublicKey`（必填）
- `request.CertificatePrivateKey`（必填，Sensitive）
- `request.ResourceTypesRegions`（可选 list，每项含 `ResourceType` + `Regions[]*string`）
- `request.CertificateId`（可选）

调用外包 `resource.Retry(tccommon.WriteRetryTimeout, ...)`，retry 块内只调接口、判空返回 `RetryableError`。成功后 **retry 外** 取 `DeployRecordId`，先 `d.SetId(helper.UInt64ToStr(deployRecordId))`，再保留 `deploy_status`/`update_sync_progress`，最后调 Read 轮询生效。`DeployStatus` 含义：1 创建成功、0 表示当前存在更新中的任务（参考 response 字段注释）。

### 3. Read 流程
- 以 `d.Id()` 作为 `DeployRecordId`，同时资源入参 `old_certificate_id` 作为 `DescribeHostUploadUpdateRecord` 的 `OldCertificateId` 过滤。
- 调用 `DescribeHostUploadUpdateRecord`（外包 `resource.Retry(tccommon.ReadRetryTimeout, ...)`）：request 中 `Limit` 取云 API 注释最大值 `200`，`Offset` 内部从 0 分页累加直到取完。
- 判空：若 `response == nil || response.Response == nil || len(response.Response.DeployRecordList) == 0`，先 `log.Printf("[CRUD] ssl_upload_certificate_instance id=%s", d.Id())` 再 `d.SetId("")`。
- 从 `DeployRecordList` 中按 id 匹配到本资源记录后，回填 `deploy_record_list`（展平到顶层 schema，不套第二层 list）。
- 增加调用 `DescribeHostUploadUpdateRecordDetail`：以 `DeployRecordId`（= id）查询托管记录详情（`DeployRecordDetail` 列表），`Limit` 取注释最大值 `200`，用于补充详情数据。其结构为 `DeployRecord`（含 `TotalCount`/`SuccessTotalCount`/`FailedTotalCount`/`RunningTotalCount`/`Status`/`CreateTime`/`RecordDetailList` 等），可在 design 中记录以便实施参考，但本资源 schema 以 `DescribeHostUploadUpdateRecord` 出参为主。
- 回填字段前判断 nil。

### 4. Delete 流程
配置型资源删除：调用 `UploadUpdateCertificateRecordRollback`（request 需 `DeployRecordId`，由 `d.Id()` 解析得到 int64），回滚更新记录使资源可被销毁。外包 `resource.Retry(tccommon.WriteRetryTimeout, ...)`，失败用 `RetryError` 包装。

### 5. 代码风格参照
参照 `tencentcloud/services/igtm/resource_tc_igtm_strategy.go` 的 schema/CRUD 骨架与 `tencentcloud/services/ssl/resource_tc_ssl_update_certificate_instance_operation.go` 的 ssl 客户端与异步等待模式：
- import：`tccommon`、`helper`、`ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"`、`schema`、`resource`、`log`、`fmt`。
- `defer tccommon.LogElapsed(...)` + `defer tccommon.InconsistentCheck(d, meta)()`。
- Config 资源通常无需 Importer；但 CONFIG 类型属可 import 的三种之一，参照 igtm 增加 Importer（State passthrough）。

## Risks / Trade-offs

- [异步任务尚未生效] → 拿到 id 后立即 SetId 再轮询 Read，即使最终失败 tfstate 仍保留 id，可 destroy 清理。
- [更新记录列表为空导致 id 被清空] → 先打印日志保留 id 现场再 SetId("")。
- [DeployStatus=0 表示存在更新中任务] → 在 Create retry 中将其视为未完成继续重试，直到返回有效 DeployRecordId（>0）。

## vendor API 字段沉淀（实施参考，无需再回读 vendor 大文件）

### UploadUpdateCertificateInstanceRequest（models.go:8539）
- `OldCertificateId *string` — 一键更新的旧证书 ID
- `ResourceTypes []*string` — 需要部署的资源类型，取值小写如 `clb`
- `CertificatePublicKey *string` — 公钥证书
- `CertificatePrivateKey *string` — 私钥证书
- `ResourceTypesRegions []*ResourceTypeRegions` — 云资源部署地域列表
- `CertificateId *string` — 证书 ID

### ResourceTypeRegions（models.go:6975）
- `ResourceType *string` — 云资源类型（clb/waf/apigateway/cos/tke/tse/tcb）
- `Regions []*string` — 地域列表

### UploadUpdateCertificateInstanceResponse（models.go:8601，Response = ResponseParams）
- `DeployRecordId *uint64` — 云资源更新任务 ID，0 表示进行中；>0 表示创建成功
- `DeployStatus *int64` — 1 创建成功；0 当前存在更新中的任务
- `UpdateSyncProgress []*UpdateSyncProgress` — 异步创建任务进度详情（可能为 nil）

### UpdateSyncProgress（models.go:8236）
- `ResourceType *string`
- `UpdateSyncProgressRegions []*UpdateSyncProgressRegion`
- `Status *int64` — 0 待处理 / 1 已处理 / 3 处理中

### UpdateSyncProgressRegion（models.go:8250）
- `Region *string`
- `TotalCount *int64`
- `OffsetCount *int64`
- `Status *int64`

### DescribeHostUploadUpdateRecordRequest（models.go:5468）
- `Offset *uint64` — 分页偏移
- `Limit *uint64` — 每页数量，默认 10，取注释最大值 200
- `OldCertificateId *string` — 原证书 ID

### DescribeHostUploadUpdateRecordResponse（models.go:5514，Response = ResponseParams）
- `TotalCount *uint64`
- `DeployRecordList []*UploadUpdateRecordInfo`

### UploadUpdateRecordInfo（models.go:8732）
- `Id *uint64`
- `OldCertId *string`
- `ResourceTypes []*string`
- `Status *uint64`
- `CreateTime *string`
- `UpdateTime *string`

### DescribeHostUploadUpdateRecordDetailRequest（models.go:5397）
- `DeployRecordId *int64` — 托管记录 ID
- `Limit *int64` — 每页数量，默认 10，最大 200
- `Offset *int64` — 偏移量，默认 0

### DescribeHostUploadUpdateRecordDetailResponse（models.go:5440，Response = ResponseParams）
- `DeployRecordDetail []*DeployRecord`

### DeployRecord（models.go:2652）
- `TotalCount/SuccessTotalCount/FailedTotalCount/RunningTotalCount/PendingTotalCount *int64`
- `Type *int64` — 0 部署 / 1 回滚
- `Status *int64` — 0 等待 / 1 成功 / 2 失败 / 3 部署中 / 4 回滚成功 / 5 回滚失败
- `RecordDetailList []*DeployRecordList`
- `CreateTime *string`

### UploadUpdateCertificateRecordRollbackRequest（models.go:8684）
- `DeployRecordId *int64` — 待回滚记录 ID

### UploadUpdateCertificateRecordRollbackResponse（models.go:8716）
- 仅 `RequestId`（无业务字段）