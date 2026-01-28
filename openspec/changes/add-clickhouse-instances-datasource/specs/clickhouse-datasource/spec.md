# TCHouse-C Instances Data Source Specification

## ADDED Requirements

### Requirement: 支持查询 ClickHouse 实例列表
数据源 `tencentcloud_clickhouse_instances` 必须能够查询腾讯云 TCHouse-C (ClickHouse) 实例列表，支持多种过滤条件。

**Rationale**: 用户需要在 Terraform 配置中动态发现和引用已存在的 ClickHouse 实例，支持批量查询和资源规划。

#### Scenario: 查询所有实例
**Given** 用户配置了数据源但未指定任何过滤条件  
**When** 执行 `terraform plan` 或 `terraform apply`  
**Then** 数据源返回当前账号下所有 ClickHouse 实例的列表

**Acceptance Criteria**:
- 调用 `DescribeInstancesNew` API 获取实例列表
- 返回的 `instance_list` 包含所有实例
- 每个实例包含完整的基础信息字段

#### Scenario: 按实例 ID 精确查询
**Given** 用户指定了 `instance_id` 参数  
**When** 执行数据源查询  
**Then** 返回匹配该 ID 的实例（如果存在）

**Acceptance Criteria**:
- `instance_id` 映射到 API 的 `SearchInstanceId` 参数
- 返回结果最多包含一个实例
- 如果实例不存在，返回空列表

#### Scenario: 按实例名称模糊查询
**Given** 用户指定了 `instance_name` 参数  
**When** 执行数据源查询  
**Then** 返回名称匹配的所有实例

**Acceptance Criteria**:
- `instance_name` 映射到 API 的 `SearchInstanceName` 参数
- 支持模糊匹配
- 返回所有名称包含指定字符串的实例

#### Scenario: 按标签过滤查询
**Given** 用户指定了 `tags` 参数（map 类型）  
**When** 执行数据源查询  
**Then** 返回包含所有指定标签的实例

**Acceptance Criteria**:
- `tags` 转换为 API 的 `SearchTags` 参数
- 支持多标签同时过滤
- 只返回同时满足所有标签条件的实例

#### Scenario: 按 VIP 地址查询
**Given** 用户指定了 `vips` 参数（字符串列表）  
**When** 执行数据源查询  
**Then** 返回 VIP 匹配的实例

**Acceptance Criteria**:
- `vips` 映射到 API 的 `Vips` 参数
- 支持多个 VIP 地址查询
- 返回匹配任一 VIP 的实例

### Requirement: 完整的实例信息映射
数据源必须返回实例的完整详细信息，涵盖网络、计费、配置、状态等所有关键字段。

**Rationale**: 用户需要完整的实例信息用于资源规划、监控和管理决策。

#### Scenario: 返回基础信息字段
**Given** 查询到实例列表  
**When** 读取 `instance_list` 中的实例  
**Then** 每个实例包含以下基础字段

**Acceptance Criteria**:
- ✅ `instance_id` - 实例唯一标识
- ✅ `instance_name` - 实例名称
- ✅ `status` - 实例状态
- ✅ `status_desc` - 状态描述
- ✅ `version` - 版本号
- ✅ `region`, `zone` - 地域和可用区
- ✅ `region_id`, `region_desc`, `zone_desc` - 地域详细信息

#### Scenario: 返回网络信息字段
**Given** 查询到实例列表  
**When** 读取实例网络配置  
**Then** 包含以下网络字段

**Acceptance Criteria**:
- ✅ `vpc_id` - 私有网络 ID
- ✅ `subnet_id` - 子网 ID
- ✅ `access_info` - 访问信息
- ✅ `eip` - 弹性公网 IP
- ✅ `ch_proxy_vip` - CHProxy VIP 地址

#### Scenario: 返回计费信息字段
**Given** 查询到实例列表  
**When** 读取实例计费信息  
**Then** 包含以下计费字段

**Acceptance Criteria**:
- ✅ `pay_mode` - 付费模式
- ✅ `create_time` - 创建时间
- ✅ `expire_time` - 过期时间
- ✅ `renew_flag` - 续费标识

#### Scenario: 返回节点配置信息
**Given** 查询到实例列表  
**When** 读取实例节点配置  
**Then** 包含 `master_summary` 和 `common_summary` 嵌套对象

**Acceptance Criteria**:
- ✅ `master_summary` 包含 Master 节点规格信息
- ✅ `common_summary` 包含 Common 节点规格信息
- ✅ 两者都包含：`spec`, `node_size`, `core`, `memory`, `disk`, `disk_type`
- ✅ 包含 `attach_cbs_spec` 嵌套对象（云硬盘规格）
- ✅ 包含其他扩展字段：`sub_product_type`, `encrypt` 等

#### Scenario: 返回标签信息
**Given** 实例配置了标签  
**When** 查询实例列表  
**Then** 标签信息正确映射到 `tags` 字段

**Acceptance Criteria**:
- ✅ `tags` 为 TypeList 类型
- ✅ 每个标签包含 `tag_key` 和 `tag_value`
- ✅ 支持多标签

#### Scenario: 返回高可用和弹性信息
**Given** 查询到实例列表  
**When** 读取实例高可用配置  
**Then** 包含以下字段

**Acceptance Criteria**:
- ✅ `ha` - 高可用标识
- ✅ `ha_zk` - ZooKeeper 高可用
- ✅ `is_elastic` - 是否弹性实例
- ✅ `kind` - 实例类型

#### Scenario: 返回日志和存储信息
**Given** 查询到实例列表  
**When** 读取实例日志和存储配置  
**Then** 包含以下字段

**Acceptance Criteria**:
- ✅ `has_cls_topic` - 是否有 CLS 日志主题
- ✅ `cls_topic_id` - CLS 主题 ID
- ✅ `cls_log_set_id` - CLS 日志集 ID
- ✅ `cos_bucket_name` - COS 存储桶名称
- ✅ `can_attach_cbs`, `can_attach_cbs_lvm`, `can_attach_cos` - 存储能力标识

#### Scenario: 返回组件和升级信息
**Given** 查询到实例列表  
**When** 读取实例组件信息  
**Then** 包含以下字段

**Acceptance Criteria**:
- ✅ `components` - 组件列表（TypeList，包含 name 和 version）
- ✅ `upgrade_versions` - 可升级版本信息
- ✅ `enable_xml_config` - XML 配置能力标识

### Requirement: 数据类型转换与空值处理
数据源必须正确处理 API 返回的数据类型转换，安全处理空值和 nil 指针。

**Rationale**: 腾讯云 SDK 使用指针类型，需要安全解引用避免程序崩溃。

#### Scenario: 安全解引用指针字段
**Given** API 返回包含指针类型的字段  
**When** 转换为 Terraform schema  
**Then** 所有指针字段都经过 nil 检查

**Acceptance Criteria**:
- 使用辅助函数或条件判断检查指针非空
- nil 指针字段设置为零值或省略
- 不会因为 nil 指针导致 panic

#### Scenario: 嵌套对象的递归转换
**Given** API 返回包含嵌套对象（如 `MasterSummary`）  
**When** 转换为 Terraform schema  
**Then** 递归转换所有层级的字段

**Acceptance Criteria**:
- `MasterSummary` 和 `CommonSummary` 正确转换为 map
- `AttachCBSSpec` 嵌套对象正确转换
- `InstanceStateInfo` 嵌套对象正确转换
- 数组类型字段（如 `Tags`, `Components`）正确转换为 TypeList

#### Scenario: 布尔值和整数类型转换
**Given** API 返回布尔指针和整数指针  
**When** 转换为 Terraform schema  
**Then** 类型转换正确

**Acceptance Criteria**:
- `*bool` 正确转换为 schema.TypeBool
- `*int64` 正确转换为 schema.TypeInt
- `*string` 正确转换为 schema.TypeString
- 保持数据精度不丢失

### Requirement: 输出文件支持
数据源支持将查询结果输出到 JSON 文件，方便用户审查和分析。

**Rationale**: 用户可能需要将查询结果导出用于离线分析或审计。

#### Scenario: 输出结果到指定文件
**Given** 用户指定了 `result_output_file` 参数  
**When** 数据源查询完成  
**Then** 将结果序列化为 JSON 并写入文件

**Acceptance Criteria**:
- `result_output_file` 参数为 Optional String 类型
- 结果以可读的 JSON 格式输出
- 文件写入失败时给出明确错误提示
- 文件内容包含完整的 `instance_list` 数据

### Requirement: 错误处理与重试
数据源必须正确处理 API 错误，实现重试逻辑以应对临时性故障。

**Rationale**: 云 API 调用可能因为网络、限流等原因失败，需要重试机制。

#### Scenario: API 调用失败时重试
**Given** API 调用返回可重试错误（如限流）  
**When** 执行数据源查询  
**Then** 自动重试直到成功或超时

**Acceptance Criteria**:
- 使用 `resource.Retry` 包装 API 调用
- 设置合理的重试超时（使用 `tccommon.ReadRetryTimeout`）
- 对可重试错误使用 `tccommon.RetryError`
- 对不可重试错误使用 `resource.NonRetryableError`

#### Scenario: API 返回错误时记录日志
**Given** API 调用失败  
**When** 处理错误  
**Then** 记录详细的错误日志

**Acceptance Criteria**:
- 使用 `log.Printf` 记录错误信息
- 日志包含 LogId, API Action, 请求体, 错误原因
- 成功调用也记录 DEBUG 日志
- 日志格式符合项目规范

#### Scenario: 空结果的优雅处理
**Given** API 返回空实例列表  
**When** 设置 Terraform state  
**Then** 不报错，返回空列表

**Acceptance Criteria**:
- `instance_list` 设置为空数组 `[]`
- 不返回错误
- 用户可以正常使用数据源（结果为空）

### Requirement: 代码质量与规范
数据源代码必须符合项目规范，通过代码检查工具验证。

**Rationale**: 保持代码库一致性和可维护性。

#### Scenario: 代码通过格式化检查
**Given** 数据源代码实现完成  
**When** 运行 `go fmt`  
**Then** 代码格式符合 Go 标准

**Acceptance Criteria**:
- 运行 `go fmt` 后无修改
- 缩进、空格、换行符合规范

#### Scenario: 代码通过 linter 检查
**Given** 数据源代码实现完成  
**When** 运行 `make lint`  
**Then** 无 linter 错误或警告

**Acceptance Criteria**:
- 通过 golangci-lint 检查
- 通过 tfproviderlint 检查
- 无未使用的变量或导入
- 错误处理完整

#### Scenario: 遵循命名规范
**Given** 数据源代码  
**When** 审查代码  
**Then** 命名符合项目规范

**Acceptance Criteria**:
- 文件名: `data_source_tc_clickhouse_instances.go`
- 数据源名: `tencentcloud_clickhouse_instances`
- 函数名: `DataSourceTencentCloudClickhouseInstances`
- 辅助函数使用小驼峰命名（如 `flattenInstanceInfo`）
- 导入别名正确（tccommon, helper）

## Notes

### API 限制说明
- 接口频率限制：20 次/秒
- 分页默认步长：10
- 标签搜索：支持多标签，逻辑为 AND（必须同时满足）
- VIP 搜索：支持多 VIP，逻辑为 OR（满足任一即可）

### 实现注意事项
- SDK 已包含 `DescribeInstancesNew` 方法，无需修改 Service 层
- 现有 `CdwchService.DescribeInstancesNew()` 仅支持按 ID 查询，数据源应直接调用 SDK API 以支持更多过滤条件
- API 返回的 `InstanceInfo` 结构体字段非常多（60+），需要完整映射避免信息丢失
- 嵌套对象（如 `MasterSummary`）使用 TypeList 而非 TypeMap，以保持结构清晰
- 标签字段映射为 TypeList（而非 TypeMap），与其他资源保持一致

### 测试数据准备
- 测试前需要在腾讯云账号中创建至少一个 ClickHouse 实例
- 建议为测试实例添加标签以验证标签过滤功能
- 测试环境需要设置有效的 `TENCENTCLOUD_SECRET_ID` 和 `TENCENTCLOUD_SECRET_KEY`
