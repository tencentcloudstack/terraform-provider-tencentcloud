## ADDED Requirements

### Requirement: 资源注册

Provider SHALL 在 `tencentcloud/provider.go` 的 `ResourcesMap` 中注册 `tencentcloud_cvm_repair_task_control_operation` 资源，使其可被 Terraform 配置引用。

#### Scenario: 资源在 Provider 中可用
- **WHEN** 用户在 `.tf` 配置中声明 `resource "tencentcloud_cvm_repair_task_control_operation" "demo" {}`
- **THEN** Terraform 能识别该资源类型并加载对应 Schema

#### Scenario: 资源出现在文档列表
- **WHEN** 查看 `tencentcloud/provider.md` 的资源列表
- **THEN** 列表中包含 `tencentcloud_cvm_repair_task_control_operation`，按字母序位于 CVM 区域

### Requirement: 必填参数 Schema

资源 SHALL 定义以下必填参数，对应 `RepairTaskControl` API 的必需输入：

- `product` (String, Required, ForceNew): 产品类型，可选值 `CVM` / `CDH` / `CPM2.0`
- `instance_ids` (List(String), Required, ForceNew): 待操作的实例 ID 列表
- `task_id` (String, Required, ForceNew): 维修任务 ID
- `operate` (String, Required, ForceNew): 操作类型，当前云 API 仅支持 `AuthorizeRepair`

#### Scenario: 缺失必填参数
- **WHEN** 用户未提供 `task_id`
- **THEN** Terraform 在 plan 阶段报错：`"task_id" is required`

#### Scenario: 必填参数被修改触发重建
- **WHEN** 用户修改 `task_id` 的值并执行 `terraform apply`
- **THEN** Terraform 计划销毁原资源并重新创建（因为字段为 ForceNew）

### Requirement: 可选参数 Schema

资源 SHALL 定义以下可选参数：

- `order_auth_time` (String, Optional, ForceNew): 预约授权时间，格式 `YYYY-MM-DD HH:MM:SS`，由云端校验时间窗口
- `task_sub_method` (String, Optional, ForceNew): 附加授权处理策略；传入 `LossyLocal` 表示允许弃盘迁移

#### Scenario: 不指定可选参数
- **WHEN** 用户仅提供必填参数
- **THEN** 资源以"立即授权、非弃盘迁移"方式创建成功

#### Scenario: 指定预约授权时间
- **WHEN** 用户提供 `order_auth_time = "2030-01-01 12:00:00"`
- **THEN** API 请求中携带 `OrderAuthTime` 参数

#### Scenario: 启用弃盘迁移
- **WHEN** 用户提供 `task_sub_method = "LossyLocal"`
- **THEN** API 请求中携带 `TaskSubMethod=LossyLocal`

### Requirement: 创建（Create）行为

资源 Create 函数 SHALL：

1. 使用 `tccommon.LogElapsed` 与 `tccommon.InconsistentCheck` defer 打印日志与一致性检查
2. 构造 `cvm.NewRepairTaskControlRequest()`，将 Schema 字段映射到请求字段
3. 通过 `meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().RepairTaskControl(request)` 调用 API
4. 使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包装调用并使用 `tccommon.RetryError` 区分可重试错误
5. 成功后调用 `d.SetId(taskId)` 设置资源 ID 为 `task_id`
6. 调用 Read 函数返回

#### Scenario: 立即授权成功
- **WHEN** 用户配置 `product="CVM"`, `instance_ids=["ins-xxx"]`, `task_id="rep-xxx"`, `operate="AuthorizeRepair"` 并 `terraform apply`
- **THEN** Provider 调用 `RepairTaskControl` API，资源 ID 设置为 `rep-xxx`，apply 成功

#### Scenario: 网络瞬时故障重试
- **WHEN** API 调用返回可重试错误（如网络超时）
- **THEN** Provider 在 WriteRetryTimeout 内自动重试，直到成功或超时

#### Scenario: 任务非待授权状态
- **WHEN** 指定的 `task_id` 已不在"待授权"状态
- **THEN** Provider 透传云端错误信息给用户，apply 失败

### Requirement: 读取（Read）行为

资源 Read 函数 SHALL 是空操作（仅 defer 打印日志），返回 nil。这是因为 `RepairTaskControl` 没有对应的 Describe API，且操作为一次性副作用。

#### Scenario: terraform refresh
- **WHEN** 用户执行 `terraform refresh` 或 `terraform plan`
- **THEN** Provider 不会调用任何云 API，state 保持不变

### Requirement: 删除（Delete）行为

资源 Delete 函数 SHALL 是空操作（仅 defer 打印日志），返回 nil。因为云 API 不支持取消已授权的维修任务。

#### Scenario: terraform destroy
- **WHEN** 用户执行 `terraform destroy`
- **THEN** Provider 不调用任何云 API，仅从 state 中移除资源

### Requirement: 资源 ID 设置

资源 SHALL 使用 `task_id` 作为 Terraform 资源 ID，便于 `terraform import` 与 state 调试时识别。

#### Scenario: 创建后查看 state
- **WHEN** 用户在 apply 成功后执行 `terraform state show tencentcloud_cvm_repair_task_control_operation.demo`
- **THEN** state ID 等于 `task_id` 的值

### Requirement: 文档模板

实现 SHALL 在 `tencentcloud/services/cvm/resource_tc_cvm_repair_task_control_operation.md` 提供文档模板，至少包含：

- 资源功能描述
- 立即授权示例
- 预约授权示例
- 弃盘迁移示例
- Import 说明（不支持或说明用法）

#### Scenario: 文档模板存在
- **WHEN** 检查 `tencentcloud/services/cvm/` 目录
- **THEN** 存在 `resource_tc_cvm_repair_task_control_operation.md` 文件且包含上述各示例段

#### Scenario: 网站文档生成
- **WHEN** 执行 `make doc`
- **THEN** 生成 `website/docs/r/cvm_repair_task_control_operation.html.markdown`，格式正确，参数说明完整

### Requirement: 验收测试

实现 SHALL 提供 `resource_tc_cvm_repair_task_control_operation_test.go`，至少包含一个验收测试用例 `TestAccTencentCloudCvmRepairTaskControlOperationResource_basic`，验证基本授权流程。

测试用例需要：
- 使用 `acctest` 包提供的常量
- 使用 `terraform-plugin-sdk/v2` 测试框架
- 在 `PreCheck` 中校验 `TF_ACC` 环境变量
- 检查资源创建后 ID 非空

#### Scenario: 测试文件存在并可编译
- **WHEN** 运行 `go build ./...`
- **THEN** 测试文件编译通过，无语法或类型错误

#### Scenario: 测试用例命名规范
- **WHEN** 检查测试文件
- **THEN** 至少包含 `TestAccTencentCloudCvmRepairTaskControlOperationResource_basic` 函数

### Requirement: 错误处理

资源 SHALL 遵循 Provider 错误处理约定：

- 所有云 API 错误必须通过 `log.Printf("[CRITAL]...")` 记录到日志
- 重试错误使用 `tccommon.RetryError(e)` 包装
- 最终失败时返回原始 error 给 Terraform

#### Scenario: API 永久失败
- **WHEN** API 返回非可重试错误（如参数非法）
- **THEN** Provider 立即返回错误，不重试，且日志中有 `[CRITAL]` 记录
