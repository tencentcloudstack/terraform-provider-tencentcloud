## Why

腾讯云 DTS（数据传输服务）同步任务的配置目前缺少专门的 Terraform 资源进行"动态修改"管理。现有资源 `tencentcloud_dts_sync_config` 通过 `ConfigureSyncJob` 接口进行初始化配置，但当任务处于运行中状态时，用户需要通过 `ModifySyncJobConfig` 接口动态调整同步对象（DynamicObjects）与同步选项（DynamicOptions），而这一能力在现有 Provider 中尚未覆盖。新增 `tencentcloud_dts_sync_job_config` 资源（RESOURCE_KIND_CONFIG）可以让用户以声明式方式管理同步任务的动态配置变更，满足运行态配置热更新的运维诉求。

## What Changes

- 新增 Terraform 资源 `tencentcloud_dts_sync_job_config`，对应代码文件 `resource_tc_tencentcloud_dts_sync_job_config.go`（位于 `tencentcloud/services/dts/`）。
- 资源类型为 RESOURCE_KIND_CONFIG，仅提供 Read/Update（RU）能力：
  - Read：调用 `DescribeSyncJobs` 接口读取同步任务配置信息。
  - Update：调用 `ModifySyncJobConfig` 接口动态修改同步对象（DynamicObjects）与同步选项（DynamicOptions）。
- 资源 id 使用 `job_id`（同步任务 ID）作为标识，Create 时设置 id 后直接复用 Update 逻辑。
- 在 `tencentcloud/provider.go` 与 `tencentcloud/provider.md` 中注册新资源。
- 新增对应的单元测试文件 `resource_tc_tencentcloud_dts_sync_job_config_test.go`（使用 gomonkey mock 云 API）。
- 新增对应的文档文件 `resource_tc_tencentcloud_dts_sync_job_config.md`。

## Capabilities

### New Capabilities
- `dts-sync-job-config`: 管理腾讯云 DTS 同步任务的动态配置（DynamicObjects / DynamicOptions），通过 ModifySyncJobConfig 接口实现运行态配置热更新，通过 DescribeSyncJobs 接口读取配置。

### Modified Capabilities
<!-- 无现有 spec 需要修改 -->

## Impact

- 新增代码文件：`tencentcloud/services/dts/resource_tc_tencentcloud_dts_sync_job_config.go` 及对应测试文件与文档文件。
- 修改注册：`tencentcloud/provider.go`、`tencentcloud/provider.md` 新增资源注册项。
- 服务层：在 `tencentcloud/services/dts/service_tencentcloud_dts.go` 中新增 `DescribeDtsSyncJobConfigById`（复用 DescribeSyncJobs）与 `ModifyDtsSyncJobConfig` 服务方法封装。
- 依赖云 API：`DescribeSyncJobs`、`ModifySyncJobConfig`（包 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206`），二者已在 vendor 中存在。
- 不影响现有资源与向后兼容。
