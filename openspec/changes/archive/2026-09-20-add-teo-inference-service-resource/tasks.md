## 1. SDK 核查与 service 层

- [x] 1.1 核查 vendored SDK `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901/` 中已存在 `CreateInferenceService`、`DescribeInferenceServices`、`ModifyInferenceService`、`OperateInferenceService` 四个接口及对应结构体（`InferenceService`、`InferenceContainerConfig`、`InferenceContainerConfigForModify`、`InferenceResourceConfig`、`InferenceResourceConfigForModify`、`InferenceTCRRepositoryConfig`、`InferenceHardwareConfig`、`InferenceHardwareConfigForModify`、`InferenceAutoScalingConfig`、`InferenceScalingPolicy`、`InferenceScheduledScalingPolicy`、`InferenceScheduledScalingAction`、`InferenceScheduledScalingEffectiveRange`、`InferenceManualInstanceConfig`、`InferenceEnvironmentVariable`、`InferenceAffinityConfig`、`SessionIdAffinityConfig`）。
- [x] 1.2 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 新增 `DescribeTeoInferenceServiceById(ctx, zoneId, serviceId string) (*teo.InferenceService, error)`：在 for 循环外构建 `DescribeInferenceServicesRequest`，`Filters=[{Name:"service-id", Values:[serviceId]}]`，`ZoneId` 传入；循环内仅变 `Offset`，`Limit=200`（SDK 注释标注最大值）；每页 SDK 调用包裹 `resource.Retry(tccommon.ReadRetryTimeout, ...)`；严格相等 `*item.ServiceId == serviceId`；未找到返回 `(nil, nil)`；retry 失败打印 `[CRITAL]` 日志，资源名使用 `inference_service`。

## 2. 资源实现

- [x] 2.1 创建 `tencentcloud/services/teo/resource_tc_teo_inference_service.go`（单文件，参考 `tencentcloud_igtm_strategy` 风格）。文件开头不加注释。布局：package + imports → `ResourceTencentCloudTeoInferenceService()` schema → Create/Read/Update/Delete → build/flatten helper。
- [x] 2.2 schema 定义：`zone_id`(Required,ForceNew)、`name`(Required,ForceNew)、`listen_port`(Required,int)、`containers`(Required,TypeList,单元素嵌套: `image_type`/`tcr_repository_config`/`startup_command`/`environment_variables`)、`resource_config`(Required,TypeList,单元素嵌套: `scaling_mode`/`hardware_spec`/`hardware_spec_id`/`hardware_config`/`auto_scaling_config`/`manual_instance_config`/`concurrency`)、`affinity_config`(Optional,TypeList,嵌套: `switch`/`affinity_mode`/`session_id_affinity_config`)、`request_paths`(Optional,TypeList of string)、`description`(Optional)；Computed: `service_id`/`status`/`scaling_status`/`current_instance_count`/`inference_url`/`create_time`/`update_time`。multi-layer nested (scheduled_actions/effective_range/scaling_policies) per design D3。不带 `Timeouts` 块（接口为同步，无轮询）。
- [x] 2.3 实现 `resourceTencentCloudTeoInferenceServiceCreate`：顶部 `defer tccommon.LogElapsed(...)` + `defer tccommon.InconsistentCheck(...)`；用 `buildTeoInferenceContainers`(Create 变体) 等组装 `CreateInferenceServiceRequest`；SDK 调用包裹 `resource.Retry(tccommon.WriteRetryTimeout, ...)`，retry 块内仅 API 调用 + nil 防御(`result==nil||result.Response==nil` → `resource.NonRetryableError`)；retry 外先打印 `logId` 与 `d.Id()`，再校验 `response.Response.ServiceId` 非空(空则 `NonRetryableError`)；`d.SetId(strings.Join([]string{zoneId, serviceId}, tccommon.FILED_SP))`；返回 Read。
- [x] 2.4 实现 `resourceTencentCloudTeoInferenceServiceRead`：split `d.Id()` 校验 2 段；调 `DescribeTeoInferenceServiceById`；`nil` 时先 `log.Printf("[CRUD] inference_service id=%s", d.Id())` 再 `d.SetId("")` 返回；否则逐字段 nil 判断后 `_ = d.Set(...)`(input + computed)；`affinity_config` 不回填(Describe 响应无此字段)。
- [x] 2.5 实现 `resourceTencentCloudTeoInferenceServiceUpdate`：仅当可更新字段变更时调 `ModifyInferenceServiceWithContext`；用 `buildTeoInferenceContainersForModify` / `buildTeoInferenceResourceConfigForModify` 组装(ForModify 变体)；`hardware_spec`/`hardware_spec_id`/`gpu_num` 不写入 Modify 请求体；`zone_id`+`service_id` 来自 `d.Id()`；SDK 调用包裹 `resource.Retry(tccommon.WriteRetryTimeout, ...)`；成功后返回 Read。
- [x] 2.6 实现 `resourceTencentCloudTeoInferenceServiceDelete`：调 `OperateInferenceServiceWithContext`，`Operation="Delete"`，`ZoneId`+`ServiceId` 来自 `d.Id()`；包裹 `resource.Retry(tccommon.WriteRetryTimeout, ...)`；同步接口无轮询。
- [x] 2.7 Importer 使用 `schema.ImportStatePassthrough`；错误/日志描述统一用资源名 `inference_service`。

## 3. Provider 注册

- [x] 3.1 在 `tencentcloud/provider.go` 的 `ResourcesMap` 注册 `"tencentcloud_teo_inference_service": teo.ResourceTencentCloudTeoInferenceService()`，置于现有 teo 资源组相邻位置。

## 4. 文档示例

- [x] 4.1 创建 `tencentcloud/services/teo/resource_tc_teo_inference_service.md`：一句话描述提到 TEO（"Provides a resource to ..."）；Example Usage 用 HCL 展示主要入参(含 jsonencode 场景若有)；Import 部分说明使用复合 ID `zoneId#serviceId`；不手写 Argument/Attribute Reference（自动生成）；不手动改 `website/`。

## 5. 单元测试

- [x] 5.1 创建 `tencentcloud/services/teo/resource_tc_teo_inference_service_test.go`：使用 **gomonkey** mock 云 API(不使用 terraform 测试套件)；覆盖 Create(成功/空 ServiceId)、Read(存在/不存在)、Update(ForModify 路径)、Delete(Operation=Delete) 业务逻辑；mock `UseTeoV20220901Client().CreateInferenceServiceWithContext` 等方法；保证生成代码可正确构建。

## 6. 验证（由后续流程执行）

- [x] 6.1 确认 `tencentcloud/...` 可编译构建（构建检查由后续流程执行，本阶段不运行 `go build`/`go vet`/`golint`）。
- [x] 6.2 确认 gateway/字段映射正确性：Create 入参均在 `CreateInferenceServiceRequest` 中存在；Modify 入参均在 `ModifyInferenceServiceRequest` 中存在；Read 出参与 `InferenceService` 结构体一致。
- [x] 6.3 收尾阶段通过 tfpacer-finalize skill 执行 `gofmt`、`make doc`、生成 changelog（不在本阶段执行）。