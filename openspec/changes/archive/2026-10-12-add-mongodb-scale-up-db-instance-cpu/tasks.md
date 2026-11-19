## 1. 服务层实现

- [x] 1.1 在 `tencentcloud/services/mongodb/service_tencentcloud_mongodb.go` 新增 `ScaleUpDBInstanceCpu` 封装方法：构造 `mongodb.NewScaleUpDBInstanceCpuRequest()`，设置 `InstanceId`/`ExtraCpu`，在 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 内调用 `me.client.UseMongodbClient().ScaleUpDBInstanceCpu`，错误用 `tccommon.RetryError(e)` 包装，返回 `FlowId(*int64)` 与 error
- [x] 1.2 新增 `ScaleDownDBInstanceCpu` 封装方法：构造 `mongodb.NewScaleDownDBInstanceCpuRequest()`，设置 `InstanceId`，retry 调用，返回 `FlowId(*int64)` 与 error
- [x] 1.3 新增 `CheckDBInstanceElasticCpuScalable` 封装方法：构造 `mongodb.NewCheckDBInstanceElasticCpuScalableRequest()`，设置 `InstanceId`，在 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 内调用，返回 `*mongodb.CheckDBInstanceElasticCpuScalableResponse` 与 error
- [x] 1.4 新增弹性 CPU 扩容生效轮询方法（如 `WaitMongodbScaleUpDbInstanceCpuComplete`）：在 `resource.Retry(3*tccommon.ReadRetryTimeout, ...)` 内轮询 `CheckDBInstanceElasticCpuScalable`，依据入参 `expectScaled bool` 判定——`IsLocked==true` 或 `IsScaled != expectScaled` 时返回 `resource.RetryableError`，API 错误用 `tccommon.RetryError(e)` 包装

## 2. 资源 Schema 定义

- [x] 2.1 创建 `tencentcloud/services/mongodb/resource_tc_mongodb_scale_up_db_instance_cpu.go`，定义 `ResourceTencentCloudMongodbScaleUpDbInstanceCpu()` 返回 `*schema.Resource`，含 Create/Read/Update/Delete + Importer(`schema.ImportStatePassthrough`)
- [x] 2.2 定义 schema 字段：`instance_id`(TypeString, Optional, ForceNew)、`extra_cpu`(TypeInt, Optional)；computed 输出：`flow_id`(TypeInt)、`scalable`(TypeBool)、`is_scaled`(TypeBool)、`is_locked`(TypeBool)、`reason`(TypeString)、`max_extra_cpu`(TypeInt)、`trigger_type`(TypeInt)、`scale_up_time`(TypeString)
- [x] 2.3 为资源添加 `Timeouts` 块（Create/Update/Delete），用于异步操作等待

## 3. 资源 CRUD 实现

- [x] 3.1 实现 `resourceTencentCloudMongodbScaleUpDbInstanceCpuCreate`：读取 `instance_id`/`extra_cpu` → 调服务层 ScaleUp（retry）→ 检查 `response == nil || response.Response == nil` 返回 NonRetryableError → 打印 logId 与 d.Id() → `d.SetId(instanceId)` → 设置 `flow_id`（FlowId 判 nil）→ 轮询 Wait 方法直到 `IsScaled == true` → 调 Read 回填
- [x] 3.2 实现 `resourceTencentCloudMongodbScaleUpDbInstanceCpuRead`：从 `d.Id()` 取 instanceId → 调服务层 Check → 若 `response == nil || response.Response == nil` 先 `log.Printf("[CRUD] tencentcloud_mongodb_scale_up_db_instance_cpu id=%s", d.Id())` 再 `d.SetId("")` → 逐字段判 nil 后 `d.Set`（scalable/is_scaled/is_locked/reason/max_extra_cpu/extra_cpu/trigger_type/scale_up_time）并 set `instance_id`
- [x] 3.3 实现 `resourceTencentCloudMongodbScaleUpDbInstanceCpuUpdate`：`instance_id` 为 ForceNew 无需处理；`d.HasChange("extra_cpu")` 时调服务层 ScaleUp（retry）→ 设置新 `flow_id`（判 nil）→ 轮询直到 `IsScaled == true` → 调 Read 回填
- [x] 3.4 实现 `resourceTencentCloudMongodbScaleUpDbInstanceCpuDelete`：从 `d.Id()` 取 instanceId → 调服务层 ScaleDown（retry）→ 轮询直到 `IsScaled == false`；若 ScaleDown 返回"未处于扩容"类错误则记录日志后返回 nil（避免 destroy 卡死）

## 4. Provider 注册

- [x] 4.1 在 `tencentcloud/provider.go` 的 ResourcesMap 中新增 `"tencentcloud_mongodb_scale_up_db_instance_cpu": mongodb.ResourceTencentCloudMongodbScaleUpDbInstanceCpu()`（放在已有 mongodb 资源注册项附近）

## 5. 单元测试（gomonkey）

- [x] 5.1 创建 `tencentcloud/services/mongodb/resource_tc_mongodb_scale_up_db_instance_cpu_test.go`，使用 gomonkey mock `UseMongodbClient().ScaleUpDBInstanceCpu`/`ScaleDownDBInstanceCpu`/`CheckDBInstanceElasticCpuScalable`
- [x] 5.2 编写 Create 单测：验证调用 ScaleUp 后 `d.SetId` 被调用、`flow_id` 被设置、轮询直到 `IsScaled==true`
- [x] 5.3 编写 Read 单测：mock Check 返回各字段，验证 state 回填与空响应时 `d.SetId("")` 行为
- [x] 5.4 编写 Update 单测：`extra_cpu` 变化触发 ScaleUp 重新调用
- [x] 5.5 编写 Delete 单测：验证调用 ScaleDown 并轮询直到 `IsScaled==false`；及未扩容错误容错分支

## 6. 文档与示例

- [x] 6.1 创建 `tencentcloud/services/mongodb/resource_tc_mongodb_scale_up_db_instance_cpu.md`：一句话描述（带上 MongoDB 产品名）、Example Usage（含 jsonencode 场景按需）、Import 部分（说明使用 instance id）
- [x] 6.2 在 `tencentcloud/provider.md` 中新增该资源的声明条目

## 7. 验证任务

- [x] 7.1 运行 `gofmt` 格式化代码（收尾阶段由 tfpacer-finalize 统一执行，或本地格式化）
- [x] 7.2 运行 `make doc` 生成 website/docs/ 文档（收尾阶段由 tfpacer-finalize 执行）
- [x] 7.3 编译验证：确保新增/修改的 go 文件可正确构建（由后续流程执行，不在本阶段手动 go build）
- [x] 7.4 校验 schema 参数与云 API CRUD 接口参数一致：ScaleUp 入参(InstanceId/ExtraCpu)、ScaleDown 入参(InstanceId)、Check 入参(InstanceId) 均在对应云 API 中存在