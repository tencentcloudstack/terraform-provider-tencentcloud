## Context

Terraform Provider for TencentCloud 目前已支持多种数据库产品（如 CDB、CynosDB、MariaDB、PostgreSQL 等），但尚未支持 TDSQL-C for MySQL（tdmysql）产品。TDSQL-C 是腾讯云自研的云原生数据库，融合了传统数据库、云原生数据库和新一代分布式数据库的优势。用户需要通过 Terraform 管理 tdmysql 实例的完整生命周期（创建、查询、修改、销毁）。

当前代码库中资源遵循以下模式（参考 `tencentcloud/services/igtm/resource_tc_igtm_strategy.go`）：
- 资源文件位于 `tencentcloud/services/<service>/` 目录下，命名格式 `resource_tc_<service>_<name>.go`
- 通过 `tencentcloud-sdk-go` 调用云 API（按服务拆分包），客户端辅助方法定义在 `tencentcloud/connectivity/client.go`
- CRUD 操作均包含 `resource.Retry` 重试逻辑（写操作用 `tccommon.WriteRetryTimeout`，读操作用 `tccommon.ReadRetryTimeout`）
- 使用 `defer tccommon.LogElapsed()` 和 `defer tccommon.InconsistentCheck()` 进行错误处理
- 不可变字段变更时在 Update 方法中通过 immutableArgs 数组检查并返回 error

云 API 接口详情（基于 `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmysql/v20211122`）：
- `CreateDBInstances`：批量创建实例，入参含 Zone、VpcId、SubnetId、SpecCode、Disk、StorageNodeNum、Replications、InstanceCount、FullReplications、CreateVersion、InstanceName、ResourceTags、InitParams、TimeUnit、TimeSpan、StorageNodeCpu、StorageNodeMem、PayMode、MCNum、Vport、Zones、AutoVoucher、VoucherIds、InstanceType、StorageType、AZMode、InstanceMode、TemplateId、SQLMode、AutoScaleConfig、SecurityGroupIds、UserName、Password、EncryptionEnable；出参 InstanceIds（字符串数组）、FlowId（int64）。此为异步接口，返回 FlowId 用于追踪创建流程。
- `DescribeDBInstanceDetail`：查询实例详情，入参 InstanceId；出参含 InstanceName、Zone、VpcId、SubnetId、CreateVersion、Vip、Vport、Status、Disk、StorageNodeNum、InitParams、ResourceTags、CreateTime、UpdateTime、Replications、FullReplications、CharSet、Node、Region、SpecCode、InstanceId、StatusDesc、StorageNodeCpu、StorageNodeMem、RenewFlag、PayMode、ExpireAt、IsolatedAt、InstanceType、StorageType、Zones、DiskUsage、BinlogStatus、AZMode、StandbyFlag、BinlogType、TimingModifyInstanceFlag、ColumnarNodeCpu、ColumnarNodeMem、ColumnarNodeNum、ColumnarNodeDisk、ColumnarNodeStorageType、ColumnarNodeSpecCode、ColumnarVip、ColumnarVport、IsSupportColumnar、InstanceCategory、SQLMode、IsSwitchFullReplicationsEnable、InstanceMode、DumperVip、DumperVport、AutoScaleConfig、TemplateId、TemplateName、AnalysisMode、AnalysisRelationInfos、AnalysisInstanceInfo、MaintenanceWindow、EncryptionEnable、EncryptionKmsRegion。
- `DescribeFlow`：查询异步任务流程状态，入参 FlowId；出参 Status（running/success/paused/failed）。用于轮询 CreateDBInstances 的异步创建结果。
- `ModifyInstanceName`：修改实例名称，入参 InstanceId、InstanceName。同步接口。
- `IsolateDBInstance`：批量隔离实例，入参 InstanceIds（字符串数组）；出参 SuccessInstanceIds、FailedInstanceIds。

实例状态值（来自 InstanceInfo.Status 注释）：creating、created、initializing、running、modifying、isolating、isolated、destroying、destroyed。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_tdmysql_db_instance` 资源，支持 tdmysql 实例完整的 CRUD 生命周期管理
- 资源 schema 设计与云 API 参数对齐，正确映射 Create/Read/Update/Delete 各接口的入参与出参
- 支持 `CreateDBInstances` 异步接口：调用后通过 `DescribeFlow` 轮询 FlowId 直至流程成功（Status=success），再调用 Read 回读实例详情
- 支持资源导入（Import）
- 在 `tencentcloud/connectivity/client.go` 中新增 tdmysql SDK 客户端辅助方法
- 遵循现有资源代码风格和最佳实践（参考 igtm_strategy）
- 编写单元测试（gomonkey mock 方式）验证业务逻辑
- 生成资源文档 `.md` 文件

**Non-Goals:**
- 不新增 datasource（仅新增 resource）
- 不修改任何现有资源的 schema 或行为
- 不支持实例规格变配（磁盘扩容、节点变更等，云 API 暂未在本次需求范围内提供对应修改接口）
- 不处理实例销毁（DestroyInstances）操作，Delete 仅调用 IsolateDBInstance 进行隔离
- 不处理跨服务的依赖关系

## Decisions

### 1. 资源 ID 设计
**决策**：使用单个 `instance_id` 作为资源 ID（取自 `CreateDBInstances` 返回的 `InstanceIds` 列表中第一个元素）

**理由**：
- `CreateDBInstances` 是批量创建接口，返回 `InstanceIds` 字符串数组，但 Terraform 资源对应单个实例更符合使用习惯
- 通过设置 `instance_count` 入参为 1（或由用户指定，取首个实例 ID），使一个 Terraform 资源对应一个 tdmysql 实例
- `DescribeDBInstanceDetail` 入参为单个 `InstanceId`，与单 ID 设计匹配
- `ModifyInstanceName` 入参为单个 `InstanceId`，与单 ID 设计匹配
- `IsolateDBInstance` 入参为 `InstanceIds` 数组，Delete 时将单个 `instance_id` 包装为单元素数组传入

**备选方案**：使用 `instance_ids`（逗号拼接）作为 ID 以支持批量场景。未采用，因为批量场景下 Read/Update/Delete 的语义复杂，且与单实例 Terraform 资源模型不一致。

### 2. Schema 字段设计
**决策**：
- **Create 入参字段（Required/Optional）**：zone（Required, ForceNew）、vpc_id（Required, ForceNew）、subnet_id（Required, ForceNew）、spec_code（Required, ForceNew）、disk（Required, ForceNew）、storage_node_num（Required, ForceNew）、replications（Required, ForceNew）、instance_name（Required）、instance_count（Optional, ForceNew, 默认 1）、full_replications（Optional, ForceNew）、create_version（Optional, ForceNew）、resource_tags（Optional, ForceNew）、init_params（Optional, ForceNew, TypeList）、time_unit（Optional, ForceNew）、time_span（Optional, ForceNew）、storage_node_cpu（Optional, ForceNew）、storage_node_mem（Optional, ForceNew）、pay_mode（Optional, ForceNew）、mc_num（Optional, ForceNew）、vport（Optional, ForceNew）、zones（Optional, ForceNew, TypeList）、auto_voucher（Optional, ForceNew）、voucher_ids（Optional, ForceNew, TypeList）、instance_type（Optional, ForceNew）、storage_type（Optional, ForceNew）、az_mode（Optional, ForceNew）、instance_mode（Optional, ForceNew）、template_id（Optional, ForceNew）、sql_mode（Optional, ForceNew）、auto_scale_config（Optional, ForceNew, TypeList）、security_group_ids（Optional, ForceNew, TypeList）、user_name（Optional, ForceNew）、password（Optional, ForceNew, Sensitive）、encryption_enable（Optional, ForceNew）
- **Create 出参字段（Computed）**：instance_ids（Computed, TypeList of TypeString）、flow_id（Computed）
- **Read 出参字段（Computed）**：instance_id（Computed）、vip（Computed）、vport（Computed）、status（Computed）、create_time（Computed）、update_time（Computed）、char_set（Computed）、node（Computed, TypeList）、region（Computed）、status_desc（Computed）、renew_flag（Computed）、expire_at（Computed）、isolated_at（Computed）、zones（Computed）、disk_usage（Computed）、binlog_status（Computed）、standby_flag（Computed）、binlog_type（Computed）、timing_modify_instance_flag（Computed）、columnar_node_cpu（Computed）、columnar_node_mem（Computed）、columnar_node_num（Computed）、columnar_node_disk（Computed）、columnar_node_storage_type（Computed）、columnar_node_spec_code（Computed）、columnar_vip（Computed）、columnar_vport（Computed）、is_support_columnar（Computed）、instance_category（Computed）、is_switch_full_replications_enable（Computed）、dumper_vip（Computed）、dumper_vport（Computed）、template_name（Computed）、analysis_mode（Computed）、analysis_relation_infos（Computed, TypeList）、analysis_instance_info（Computed, TypeList）、maintenance_window（Computed, TypeList）、encryption_kms_region（Computed）

**理由**：
- 严格依据云 API 入参/出参映射关系定义 schema
- `instance_name` 在 Create 和 ModifyInstanceName 两个接口中都存在，设为 Required 且非 ForceNew，以便通过 Update 修改
- 其余 Create 入参字段无对应修改接口，设为 ForceNew，在 Update 的 immutableArgs 检查中返回 error
- `instance_count` 设为 Optional 且 ForceNew，默认 1，确保一个 Terraform 资源对应一个实例（取首个 instance_id 作为资源 ID）
- Read 出参字段设为 Computed，由 `DescribeDBInstanceDetail` 回写
- 子结构（init_params、node、auto_scale_config、analysis_relation_infos、analysis_instance_info、maintenance_window）使用 TypeList 平铺字段

### 3. Create 逻辑设计
**决策**：Create 操作调用 `CreateDBInstances`，由于该接口为异步接口（返回 FlowId），调用成功后需调用 `DescribeFlow` 轮询 FlowId 直至 Status=success，再调用 `DescribeDBInstanceDetail` 回读实例详情。

**理由**：
- `CreateDBInstances` 返回 FlowId，表明创建是异步流程
- 根据需求说明，异步接口调用后需调用 Read 接口轮询直到接口生效
- 通过 `DescribeFlow` 查询流程状态（running/success/paused/failed），当 Status=success 时表示实例创建完成
- 流程成功后从 `InstanceIds` 取首个元素作为 `instance_id`，设置资源 ID
- 若流程状态为 failed 或 paused，返回 error

**轮询实现**：使用 `resource.Retry` 包裹 `DescribeFlow` 调用，当 Status!=success 时返回 `resource.RetryableError` 继续重试，超时后返回 error。

### 4. Read 逻辑设计
**决策**：Read 操作调用 `DescribeDBInstanceDetail`，以 `instance_id`（即 d.Id()）作为入参查询实例详情。

**理由**：
- `DescribeDBInstanceDetail` 入参为单个 `InstanceId`，与单 ID 设计匹配
- 若返回 response/Response 为空，先打印 `log.Printf("[CRUD] tdmysql_db_instance id=%s", d.Id())` 保留现场，再 `d.SetId("")`
- 设置字段前判断 Response 中对应字段是否为 nil，nil 则不调用 setXX()
- 使用 `tccommon.ReadRetryTimeout` 作为 retry 超时

### 5. Update 逻辑设计
**决策**：Update 操作仅支持修改 `instance_name`（调用 `ModifyInstanceName`），其余顶层字段均加入 immutableArgs 数组，若检测到变更则返回 error。

**理由**：
- 云 API 仅提供 `ModifyInstanceName` 修改实例名称，未提供修改规格/磁盘/节点数等接口
- 根据 CODEBUDDY 规则，对于非 DATASOURCE 资源，若只有部分修改接口，则将无修改接口的字段加入 immutableArgs，变更时报错
- `instance_name` 在 ModifyInstanceName 入参中存在，可在 Update 中修改
- `instance_id` 字段为 ForceNew，不放入 immutableArgs

### 6. Delete 逻辑设计
**决策**：Delete 操作调用 `IsolateDBInstance`，将 `instance_id`（d.Id()）包装为单元素数组作为 `InstanceIds` 入参。

**理由**：
- `IsolateDBInstance` 入参为 `InstanceIds` 字符串数组，将单个 instance_id 包装为单元素数组
- 隔离是 tdmysql 实例的销毁前置操作（实例生命周期：running → isolating → isolated）
- 使用 `tccommon.WriteRetryTimeout` 作为 retry 超时
- 调用成功后检查 `SuccessInstanceIds` 是否包含目标 instance_id

### 7. 客户端辅助方法设计
**决策**：在 `tencentcloud/connectivity/client.go` 中新增 tdmysql SDK import、`tdmysqlv20211122Conn` 结构体字段和 `UseTdmysqlV20211122Client()` 方法。

**理由**：
- 现有资源均通过 `meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseXxxClient()` 获取 SDK 客户端
- tdmysql SDK 包已存在于 vendor 目录（`vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmysql/v20211122`），但 client.go 中尚未引入
- 参考现有 `UseTdcpgClient`/`UseIgtmV20231024Client` 模式实现

### 8. 测试策略
**决策**：使用 gomonkey mock 方式编写单元测试。

**理由**：
- 新增资源应使用 mock（gomonkey）方法对云 API 进行 mock 处理
- 仅进行业务代码逻辑的单元测试
- 使用 `go test -gcflags=all=-l` 运行测试
- 测试覆盖 Create（含异步轮询）、Read、Update、Delete 操作

## Risks / Trade-offs

- **[Risk] CreateDBInstances 是批量创建接口，资源 ID 取首个实例** → 当 instance_count > 1 时，Terraform 资源仅管理首个实例，其余实例不在 Terraform 管理范围内。通过设置 instance_count 默认为 1 并在文档中说明来缓解。
- **[Risk] 异步创建流程轮询可能超时** → 使用 `resource.Retry` 包裹 `DescribeFlow` 调用，重试超时由系统默认控制。若超时则返回 error，用户可重新 apply。
- **[Risk] 实例规格变配不支持** → 除 instance_name 外所有字段均 ForceNew 或 immutable，用户修改规格需重建资源。这是云 API 能力限制，非设计缺陷。
- **[Risk] Delete 仅隔离不销毁** → IsolateDBInstance 将实例置为 isolated 状态，未真正销毁。符合云产品最佳实践（隔离后可恢复或定时销毁）。
- **[Trade-off] 单实例资源模型 vs 批量实例模型** → 选择单实例模型，牺牲了批量管理能力，换取了与 Read/Update/Delete 接口语义的一致性和 Terraform 状态管理的简洁性。
