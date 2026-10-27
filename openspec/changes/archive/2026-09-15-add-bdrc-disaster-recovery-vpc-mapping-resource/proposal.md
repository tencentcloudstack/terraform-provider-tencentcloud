## Why

腾讯云多活容灾（BDRC，Backup & Disaster Recovery Center）已发布"容灾站点对 VPC 网络映射"能力：用户在搭建双活/容灾架构时，需要在源端 VPC/子网与目标端 VPC/子网之间建立网络映射，容灾切换后业务流量才能无缝切换到目标端。该能力已通过云 API（`CreateDisasterRecoveryVpcMapping` / `DescribeVpcMappings` / `DeleteDisasterRecoveryVpcMapping`）发布，但 Terraform Provider 尚未支持 bdrc 产品（`tencentcloud/services/bdrc/` 目录不存在），用户无法用 IaC 方式管理容灾 VPC 映射。

本 change 新增通用型（RESOURCE_KIND_GENERAL）资源 `tencentcloud_bdrc_disaster_recovery_vpc_mapping`，让用户在 HCL 里声明"某站点对下的一条 VPC 映射规则"，由 Provider 调用 BDRC 云 API 完成映射的创建、查询与删除，管理资源整个生命周期。

## What Changes

- 新增通用型资源 `tencentcloud_bdrc_disaster_recovery_vpc_mapping`，对应文件 `tencentcloud/services/bdrc/resource_tc_bdrc_disaster_recovery_vpc_mapping.go`。
- 该资源仅有 CRD 接口（云 API 未提供 Update/Modify 类接口），因此：
  - 仅 `id` 相关的标识语义之外的全部业务字段（`site_pair_id`、`source_vpc_id`、`source_subnet_id`、`target_vpc_id`、`target_subnet_id`）均设为 `ForceNew: true`；
  - 不实现 Update 调用云 API 的逻辑（无可用接口），所有变更通过销毁重建完成。
- Create 接口 `CreateDisasterRecoveryVpcMapping` 的响应**不返回映射 ID**（仅 RequestId），因此 Create 成功后需调用 `DescribeVpcMappings` 按源端 VPC/子网过滤查询，拿到新建映射的主键 `Id`（uint64）后再 `d.SetId(...)`。
- 资源 ID 采用联合 ID：`sitePairId` + `vpcMappingId`，以 `tccommon.FILED_SP` 为分隔符；Read/Delete 时从 `d.Id()` 拆出两段分别填充请求参数。
- Describe 出参中映射列表已按本项目规范**平铺展开**到资源 schema 顶层（不引入 `vpc_mapping_set` 嵌套层）：`id`、`site_pair_id`、`source_vpc`、`source_subnet`、`target_vpc`、`target_subnet`、`status`、`life_state`。
- Delete 调用 `DeleteDisasterRecoveryVpcMapping`，入参 `VpcMappingIds` 为 `[]*uint64`，仅传该资源对应的单个映射 ID。
- 所有云 API 调用（Create/Describe/Delete）均以 `tccommon.ReadRetryTimeout` / `tccommon.WriteRetryTimeout` 为超时时间包裹 retry，失败时用 `tccommon.RetryError()` 包装。
- 新增 `tencentcloud/connectivity/client.go` 中 bdrc 客户端方法 `UseBdrcV20260330Client()` 及对应连接字段（当前 connectivity 层尚无 bdrc 客户端）。
- 在 `tencentcloud/provider.go` 与 `tencentcloud/provider.md` 注册新资源。
- 新增资源文档 `tencentcloud/services/bdrc/resource_tc_bdrc_disaster_recovery_vpc_mapping.md`（含一句话描述、Example Usage、Import 说明），`website/docs/` 下文档由收尾阶段 `make doc` 统一生成。
- 新增单元测试 `tencentcloud/services/bdrc/resource_tc_bdrc_disaster_recovery_vpc_mapping_test.go`，使用 gomonkey mock 云 API，仅测试业务代码逻辑（不使用 terraform 测试套件）。

## Capabilities

### New Capabilities
- `bdrc-disaster-recovery-vpc-mapping-resource`: 新增 `tencentcloud_bdrc_disaster_recovery_vpc_mapping` 资源的 schema、CRUD 行为（Create=创建后回查取 ID / Read=按联合 ID 查询 / Update=不支持（CRD-only，全部字段 ForceNew）/ Delete=按映射 ID 删除）、ID 约定（`sitePairId#vpcMappingId`）、字段约束、空值保护、retry 规范、文档与测试规范，作为 bdrc 容灾 VPC 映射能力的 IaC 入口。

### Modified Capabilities

<!-- 不修改任何已有 capability 的 requirement，仅在 provider.go / provider.md / connectivity/client.go 增加注册与客户端连接代码。 -->

## Impact

- 新文件：
  - `tencentcloud/services/bdrc/resource_tc_bdrc_disaster_recovery_vpc_mapping.go`
  - `tencentcloud/services/bdrc/resource_tc_bdrc_disaster_recovery_vpc_mapping.md`
  - `tencentcloud/services/bdrc/resource_tc_bdrc_disaster_recovery_vpc_mapping_test.go`
  - `website/docs/r/bdrc_disaster_recovery_vpc_mapping.html.markdown`（由收尾阶段 `make doc` 生成）
- 既有文件改动：
  - `tencentcloud/connectivity/client.go`：新增 `bdrcV20260330Conn *bdrcv20260330.Client` 字段与 `UseBdrcV20260330Client()` 方法（参考 `UseIgtmV20231024Client` 实现）。
  - `tencentcloud/provider.go`：新增 `"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/bdrc"` import，并在资源注册 map 中追加 `"tencentcloud_bdrc_disaster_recovery_vpc_mapping": bdrc.ResourceTencentCloudBdrcDisasterRecoveryVpcMapping()`。
  - `tencentcloud/provider.md`：新增 BDRC 产品段并登记资源名，使 gendoc 索引可识别。
- 依赖：SDK `tencentcloud-sdk-go/tencentcloud/bdrc/v20260330` 已 vendor 到 `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bdrc/`，包含 `CreateDisasterRecoveryVpcMapping`、`DescribeVpcMappings`、`DeleteDisasterRecoveryVpcMapping` 及对应 model，无需额外升级 SDK。
- 不修改任何既有资源/数据源 schema，对现有用户零影响。
