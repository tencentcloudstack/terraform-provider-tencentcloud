# tdmysql-db-instance Specification

## Purpose
TBD - created by archiving change add-tdmysql-db-instance. Update Purpose after archive.
## Requirements
### Requirement: tdmysql 实例资源定义
系统 SHALL 提供 `tencentcloud_tdmysql_db_instance` 资源，用于管理腾讯云 TDSQL-C MySQL（tdmysql）实例的完整生命周期，包括创建、查询、修改实例名称与隔离（销毁）。

#### Scenario: 创建实例
- **WHEN** 用户定义 `tencentcloud_tdmysql_db_instance` 资源并执行 `terraform apply`
- **THEN** 系统调用 `CreateDBInstances` API 创建实例
- **AND** 使用 `tccommon.WriteRetryTimeout` 作为重试超时包裹 API 调用
- **AND** API 调用失败时使用 `tccommon.RetryError()` 包装错误并返回
- **AND** 成功后检查返回值：`Response` 不为 nil 且 `InstanceIds` 长度大于 0，否则返回 `NonRetryableError`
- **AND** 取返回的 `InstanceIds` 列表中的第一个实例 ID 作为 Terraform state id
- **AND** 记录 `logId` 与 `d.Id()` 便于排障
- **AND** 创建完成后调用 `DescribeDBInstanceDetail` 轮询实例详情，直到实例可读以保证后续状态收敛一致

#### Scenario: 读取实例状态
- **WHEN** 系统需要刷新 `tencentcloud_tdmysql_db_instance` 资源状态
- **THEN** 系统使用 `d.Id()` 作为 `InstanceId` 调用 `DescribeDBInstanceDetail` API
- **AND** API 调用使用 `tccommon.ReadRetryTimeout` 作为重试超时
- **AND** retry 块内仅执行接口调用，不在 retry 内执行 set 操作
- **AND** 设置 id 与回填 state 的操作放到 retry 块外、错误处理后
- **AND** 在调用 `setXX()` 设置字段前先判断 Response 中对应字段是否为 nil，为 nil 则不调用 set
- **AND** 若云端返回空（response / Response 为空），**先**打印 `log.Printf("[CRUD] tdmysql_db_instance id=%s", d.Id())` 保留现场，**再**执行 `d.SetId("")`
- **AND** 将 `DescribeDBInstanceDetail` 返回的所有非空字段回填到 Terraform state

#### Scenario: 修改实例名称
- **WHEN** 用户修改 `instance_name` 字段并执行 `terraform apply`
- **THEN** 系统调用 `ModifyInstanceName` API 修改实例名称
- **AND** 请求包含从 `d.Id()` 获取的 `InstanceId` 与新的 `InstanceName`
- **AND** 使用 `tccommon.WriteRetryTimeout` 作为重试超时
- **AND** 修改成功后调用 Read 方法刷新状态

#### Scenario: 不可变字段变更保护
- **WHEN** 用户修改除 `instance_name` 以外的顶层创建参数（如 zone、vpc_id、subnet_id、spec_code、disk、storage_node_num、replications、instance_count 等）
- **THEN** 系统返回错误，拒绝该变更
- **AND** 错误信息提示对应参数不可变

#### Scenario: 隔离（销毁）实例
- **WHEN** 用户执行 `terraform destroy` 或删除 `tencentcloud_tdmysql_db_instance` 资源
- **THEN** 系统调用 `IsolateDBInstance` API 隔离实例
- **AND** 请求的 `InstanceIds` 包含从 `d.Id()` 获取的实例 ID
- **AND** 使用 `tccommon.WriteRetryTimeout` 作为重试超时
- **AND** API 调用失败时使用 `tccommon.RetryError()` 包装错误并返回
- **AND** 成功后检查 `SuccessInstanceIds` 是否包含目标实例 ID，若未包含则返回错误

#### Scenario: 导入已存在实例
- **WHEN** 用户使用 `terraform import` 导入实例
- **THEN** 系统接受实例 ID 作为导入标识（RESOURCE_KIND_GENERAL 支持 import）
- **AND** 系统调用 Read 操作获取实例详情

### Requirement: 资源 Schema 定义
资源 Schema SHALL 包含创建实例所需的输入参数与查询返回的只读属性，并遵循指定的必填/可变约束。

#### Scenario: 创建输入参数
- **WHEN** 系统定义资源 Schema
- **THEN** Schema 包含 `CreateDBInstances` 接口的所有入参字段（snake_case 映射）：zone、vpc_id、subnet_id、spec_code、disk、storage_node_num、replications、instance_count、full_replications、create_version、instance_name、resource_tags、init_params、time_unit、time_span、storage_node_cpu、storage_node_mem、pay_mode、mc_num、vport、zones、auto_voucher、voucher_ids、instance_type、storage_type、az_mode、instance_mode、template_id、sql_mode、auto_scale_config、security_group_ids、user_name、password、encryption_enable
- **AND** 每个字段都有清晰的 Description 说明
- **AND** 复杂类型字段（resource_tags、init_params、auto_scale_config）使用 TypeList/嵌套 schema 表示

#### Scenario: 查询只读属性
- **WHEN** 系统定义资源 Schema
- **THEN** Schema 包含 `DescribeDBInstanceDetail` 接口返回的只读属性（Computed）：instance_id、vip、vport、status、create_time、update_time、char_set、node、region、status_desc、renew_flag、expire_at、isolated_at、disk_usage、binlog_status、standby_flag、binlog_type、timing_modify_instance_flag、columnar_node_cpu、columnar_node_mem、columnar_node_num、columnar_node_disk、columnar_node_storage_type、columnar_node_spec_code、columnar_vip、columnar_vport、is_support_columnar、instance_category、is_switch_full_replications_enable、dumper_vip、dumper_vport、template_name、analysis_mode、analysis_relation_infos、analysis_instance_info、maintenance_window、encryption_kms_region
- **AND** `CreateDBInstances` 返回的 instance_ids、flow_id 作为 Computed 属性

#### Scenario: 字段类型映射
- **WHEN** 系统定义资源 Schema
- **THEN** SDK 中 `*string` 字段映射为 `schema.TypeString`
- **AND** SDK 中 `*int64`/`*uint64` 字段映射为 `schema.TypeInt`
- **AND** SDK 中 `*bool` 字段映射为 `schema.TypeBool`
- **AND** SDK 中 `[]*string` 字段映射为 `schema.TypeList`（Elem 为 TypeString）
- **AND** SDK 中 `[]*ResourceTag` 字段映射为 `schema.TypeList`（嵌套 Resource，含 tag_key、tag_value）
- **AND** SDK 中 `[]*InstanceParam` 字段映射为 `schema.TypeList`（嵌套 Resource，含 param、value）
- **AND** SDK 中 `*AutoScalingConfig` 字段映射为 `schema.TypeList`（嵌套 Resource，max_items=1，含 range_min、range_max）
- **AND** SDK 中 `*float64` 字段（auto_scale_config 的 range_min/range_max）映射为 `schema.TypeFloat`

### Requirement: tdmysql 客户端访问方法
系统 SHALL 在 `tencentcloud/connectivity/client.go` 中提供 tdmysql SDK 客户端的访问方法。

#### Scenario: UseTdmysqlV20211122Client 方法
- **WHEN** 服务层或资源层需要访问 tdmysql SDK 客户端
- **THEN** 通过 `meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmysqlV20211122Client()` 获取 `*tdmysqlv20211122.Client`
- **AND** 该方法在 `TencentCloudClient` 结构体上定义，使用懒加载缓存（`tdmysqlv20211122Conn` 字段）
- **AND** 客户端使用 `NewClientProfile(300)` 创建并设置 `LogRoundTripper`

### Requirement: 服务层方法实现
系统 SHALL 在 `tencentcloud/services/tdmysql/service_tencentcloud_tdmysql.go` 中提供 tdmysql 实例管理的服务层辅助方法。

#### Scenario: DescribeTdmysqlDbInstanceById 方法
- **WHEN** 调用 `DescribeTdmysqlDbInstanceById(ctx, instanceId)`
- **THEN** 返回 `*tdmysqlv20211122.DescribeDBInstanceDetailResponseParams` 或 nil
- **AND** 构造请求设置 `InstanceId` 为入参
- **AND** 使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包裹调用
- **AND** retry 块内检查 `result == nil || result.Response == nil`，为空返回 `NonRetryableError`
- **AND** 失败时使用 `tccommon.RetryError()` 包装错误
- **AND** 使用 `ratelimit.Check(request.GetAction())` 进行限流检查

#### Scenario: IsolateTdmysqlDbInstance 方法
- **WHEN** 调用 `IsolateTdmysqlDbInstance(ctx, instanceId)`
- **THEN** 调用 `IsolateDBInstance` API 隔离实例
- **AND** 构造请求的 `InstanceIds` 包含入参 instanceId
- **AND** 使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包裹调用
- **AND** 返回错误信息或 nil

#### Scenario: DescribeTdmysqlFlow 方法
- **WHEN** 调用 `DescribeTdmysqlFlow(ctx, flowId)`
- **THEN** 调用 `DescribeFlow` API 查询异步任务流程状态
- **AND** 返回流程状态字符串（running/success/paused/failed）或错误
- **AND** 使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包裹调用

### Requirement: 异步创建轮询
系统 SHALL 在创建实例后轮询异步任务流程状态，确保实例创建完成后再返回。

#### Scenario: 创建后轮询流程状态
- **WHEN** `CreateDBInstances` 返回 `FlowId`
- **THEN** 系统使用 `resource.Retry` 或 `resource.StateChangeConf` 轮询 `DescribeFlow` 接口
- **AND** 当流程状态为 `success` 时认为创建完成
- **AND** 当流程状态为 `failed` 时返回错误
- **AND** 当流程状态为 `running` 时继续等待
- **AND** 轮询超时时间使用合理的超时配置

### Requirement: Provider 注册
系统 SHALL 在 `tencentcloud/provider.go` 中注册 `tencentcloud_tdmysql_db_instance` 资源。

#### Scenario: 资源注册
- **WHEN** Provider 初始化
- **THEN** `provider.go` 的 ResourcesMap 中包含 `"tencentcloud_tdmysql_db_instance": tdmysql.ResourceTencentCloudTdmysqlDbInstance()`
- **AND** 导入 `github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tdmysql` 包（别名 `tdmysql`）

#### Scenario: provider.md 声明
- **WHEN** 执行 `make doc` 生成文档
- **THEN** `tencentcloud/provider.md` 资源列表中包含 `tencentcloud_tdmysql_db_instance` 一行

### Requirement: 测试覆盖
系统 SHALL 提供 `resource_tc_tdmysql_db_instance_test.go` 单元测试，使用 gomonkey mock 云 API 进行业务逻辑测试。

#### Scenario: Create 单元测试
- **WHEN** 运行 Create 单元测试
- **THEN** 使用 gomonkey mock `UseTdmysqlV20211122Client` 返回模拟客户端
- **AND** mock `CreateDBInstancesWithContext` 返回包含 InstanceIds 的响应
- **AND** mock `DescribeFlowWithContext` 返回 success 状态
- **AND** mock `DescribeDBInstanceDetailWithContext` 返回实例详情
- **AND** 验证 `d.Id()` 与返回的实例 ID 一致

#### Scenario: Read 单元测试
- **WHEN** 运行 Read 单元测试
- **THEN** mock `DescribeDBInstanceDetailWithContext` 返回实例详情
- **AND** 验证 state 中各字段被正确回填

#### Scenario: Update 单元测试
- **WHEN** 运行 Update 单元测试（修改 instance_name）
- **THEN** mock `ModifyInstanceNameWithContext` 返回成功
- **AND** 验证名称更新成功

#### Scenario: Delete 单元测试
- **WHEN** 运行 Delete 单元测试
- **THEN** mock `IsolateDBInstanceWithContext` 返回包含目标实例 ID 的 SuccessInstanceIds
- **AND** 验证隔离操作成功

#### Scenario: 资源不存在的处理
- **WHEN** Read 时 `DescribeDBInstanceDetail` 返回空
- **THEN** 先打印 `[CRUD] tdmysql_db_instance id=<id>` 日志
- **AND** 设置 `d.SetId("")`
- **AND** 不返回错误

### Requirement: 文档完整性
系统 SHALL 提供完整的资源文档 `resource_tc_tdmysql_db_instance.md`。

#### Scenario: 文档内容
- **WHEN** 用户查看资源文档
- **THEN** 包含一句话描述，且描述中带上所属云产品名称（TDSQL-C MySQL）
- **AND** 格式为 "Provides a resource to ..."
- **AND** 包含 Example Usage 部分
- **AND** 包含 Import 部分（RESOURCE_KIND_GENERAL 支持 import）
- **AND** Import 示例说明使用实例 ID
- **AND** 不手动添加 `Argument Reference` 和 `Attribute Reference` 部分（由工具自动生成）
- **AND** 涉及 json 字符串的场景使用 `jsonencode()` 函数

