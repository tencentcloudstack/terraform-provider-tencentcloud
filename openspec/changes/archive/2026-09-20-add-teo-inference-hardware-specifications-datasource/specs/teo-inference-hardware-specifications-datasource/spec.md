## ADDED Requirements

### Requirement: 查询推理硬件规格列表
数据源 `tencentcloud_teo_inference_hardware_specifications` SHALL 通过调用 TEO API `DescribeInferenceHardwareSpecifications` 接口，按站点 ID 查询并返回该站点支持的推理硬件规格列表。

#### Scenario: 按站点 ID 成功查询硬件规格列表
- **WHEN** 用户配置 `zone_id` 为有效的 TEO 站点 ID 并执行 `terraform refresh`
- **THEN** 系统调用 `DescribeInferenceHardwareSpecifications` 接口（入参 `ZoneId` 为该站点 ID），将返回的 `HardwareSpecifications` 列表展开写入 `hardware_specifications` 字段，其中每个元素包含 `spec`、`hardware_spec_id`、`name`、`gpu_num`、`cpu_num`、`mem_size`、`gpu_mem_size`、`disk_size`、`allowed_gpu_nums` 字段

#### Scenario: 必填参数校验
- **WHEN** 用户未提供 `zone_id`
- **THEN** Terraform 在 plan 阶段报错，提示 `zone_id` 为必填参数

### Requirement: Read 重试与空返回处理
数据源 Read 方法 SHALL 使用 `tccommon.ReadRetryTimeout` 作为超时时间进行 retry 处理；当云 API 在 retry 块内返回空（`response == nil` 或 `Response.HardwareSpecifications` 长度为 0）时，SHALL 返回 `NonRetryableError` 而非直接清空 id。

#### Scenario: 云 API 返回空列表
- **WHEN** `DescribeInferenceHardwareSpecifications` 接口返回空列表
- **THEN** Read 方法在 retry 块内返回 `NonRetryableError`，让外层 retry 继续尝试，并在失败路径打印 `log.Printf("[DATASOURCE] read empty, skip SetId")` 提示

#### Scenario: 云 API 短暂波动后恢复
- **WHEN** 接口调用首次失败但重试后成功
- **THEN** 使用 `tccommon.RetryError()` 包装错误继续重试，最终成功时正常设置字段与 `d.SetId`

### Requirement: 字段类型映射
数据源 SHALL 按云 API 字段类型正确映射 Terraform schema 类型：字符串字段（`Spec`、`HardwareSpecId`、`Name`）使用 `schema.TypeString`；浮点数字段（`GPUNum`、`CPUNum`）使用 `schema.TypeFloat`；整数字段（`MemSize`、`GPUMemSize`、`DiskSize`）使用 `schema.TypeInt`；GPU 卡数列表（`AllowedGPUNums`）使用 `schema.TypeList` 且元素为 `schema.TypeFloat`。

#### Scenario: 浮点字段正确读取
- **WHEN** 云 API 返回 `GPUNum` 为 `1.0`、`CPUNum` 为 `8.0`
- **THEN** 数据源 `gpu_num`、`cpu_num` 字段分别为 `1.0`、`8.0`，类型为 float

#### Scenario: GPU 卡数列表正确读取
- **WHEN** 云 API 返回 `AllowedGPUNums` 为 `[1, 2, 4]`
- **THEN** 数据源 `allowed_gpu_nums` 字段为包含 `[1.0, 2.0, 4.0]` 的 float 列表

### Requirement: 数据源 ID 生成
数据源 SHALL 使用 `helper.DataResourceIdsHash(ids)` 基于返回的硬件规格 `HardwareSpecId` 列表生成哈希作为数据源 ID。

#### Scenario: 正常返回多个规格
- **WHEN** 接口返回多个硬件规格
- **THEN** 数据源 ID 为基于各 `HardwareSpecId` 拼接后的哈希值

### Requirement: Provider 注册
数据源 SHALL 在 `tencentcloud/provider.go` 中以名称 `tencentcloud_teo_inference_hardware_specifications` 注册，并在 `tencentcloud/provider.md` 中添加对应文档条目。

#### Scenario: 数据源可被 Terraform 识别
- **WHEN** Provider 构建完成
- **THEN** `terraform plan` 能识别 `data.tencentcloud_teo_inference_hardware_specifications` 数据源类型

### Requirement: 结果输出文件
数据源 SHALL 支持可选参数 `result_output_file`，当用户提供该参数时，将查询结果写入指定文件。

#### Scenario: 输出结果到文件
- **WHEN** 用户配置 `result_output_file` 为有效路径
- **THEN** 查询结果以 JSON 格式写入该文件

### Requirement: 单元测试
数据源 SHALL 包含单元测试文件 `data_source_tc_teo_inference_hardware_specifications_test.go`，使用 gomonkey mock 云 API 进行业务逻辑测试，不使用 Terraform 测试套件。

#### Scenario: mock 云 API 返回构造数据
- **WHEN** 运行单元测试
- **THEN** 通过 gomonkey mock `DescribeTeoInferenceHardwareSpecificationsByFilter` 返回构造的硬件规格数据，校验 flatten 后各字段值正确