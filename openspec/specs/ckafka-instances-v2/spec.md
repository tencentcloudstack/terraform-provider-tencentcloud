# CKafka Instances V2 Datasource Specification

## Requirements

### Requirement: CKafka Instances V2 Query
系统应当提供 `tencentcloud_ckafka_instances_v2` datasource,允许用户通过标准化的过滤器查询 CKafka 实例列表详情。

#### Scenario: Query instances with filters
- **GIVEN** 用户配置了 datasource `tencentcloud_ckafka_instances_v2`
- **WHEN** 用户指定 filters 参数(例如 name="InstanceId", value=["ckafka-xxx"])
- **THEN** 系统应当调用 DescribeInstancesDetail API
- **AND** 返回符合过滤条件的实例列表
- **AND** 每个实例包含完整的详情字段(InstanceId, InstanceName, Vip, Vport, Status, Bandwidth, DiskSize, ZoneId, VpcId, SubnetId 等)

#### Scenario: Query all instances without filters
- **GIVEN** 用户配置了 datasource `tencentcloud_ckafka_instances_v2`
- **WHEN** 用户不指定任何 filters
- **THEN** 系统应当返回账号下所有 CKafka 实例
- **AND** 实例列表按照 API 默认顺序返回

#### Scenario: Export results to file
- **GIVEN** 用户指定了 result_output_file 参数
- **WHEN** 查询成功完成
- **THEN** 系统应当将查询结果写入指定的文件
- **AND** 文件格式为 JSON

### Requirement: Filter Structure
系统应当支持标准化的 filters 参数结构,遵循 Terraform 最佳实践。

#### Scenario: Filter by instance ID
- **GIVEN** 用户指定 filters 参数
- **AND** filter.name = "InstanceId"
- **AND** filter.value = ["ckafka-xxx", "ckafka-yyy"]
- **WHEN** 执行查询
- **THEN** 系统应当仅返回指定 ID 的实例

#### Scenario: Filter by VPC ID
- **GIVEN** 用户指定 filters 参数
- **AND** filter.name = "VpcId"
- **AND** filter.value = ["vpc-xxx"]
- **WHEN** 执行查询
- **THEN** 系统应当返回指定 VPC 内的所有实例

#### Scenario: Filter by subnet ID
- **GIVEN** 用户指定 filters 参数
- **AND** filter.name = "SubNetId"
- **AND** filter.value = ["subnet-xxx"]
- **WHEN** 执行查询
- **THEN** 系统应当返回指定子网内的所有实例

#### Scenario: Multiple filters
- **GIVEN** 用户指定多个 filters
- **WHEN** 执行查询
- **THEN** 系统应当应用所有过滤条件(AND 逻辑)
- **AND** 返回同时满足所有条件的实例

### Requirement: Instance Detail Fields
系统应当返回 DescribeInstancesDetail API 提供的完整实例详情字段。

#### Scenario: Basic instance information
- **GIVEN** 查询返回实例列表
- **THEN** 每个实例必须包含以下基本字段:
  - instance_id: 实例 ID
  - instance_name: 实例名称
  - vip: 实例访问 VIP
  - vport: 实例访问端口
  - status: 实例状态(Integer)

#### Scenario: Network configuration
- **GIVEN** 查询返回实例列表
- **THEN** 每个实例必须包含网络配置字段:
  - vpc_id: 私有网络 ID
  - subnet_id: 子网 ID
  - zone_id: 可用区 ID
  - zone_ids: 跨可用区列表
  - vip_list: 虚拟 IP 列表(嵌套对象)

#### Scenario: Resource specifications
- **GIVEN** 查询返回实例列表
- **THEN** 每个实例必须包含资源规格字段:
  - bandwidth: 实例带宽(Mbps)
  - disk_size: 磁盘大小(GB)
  - disk_type: 磁盘类型
  - instance_type: 实例类型
  - cluster_type: 集群类型

#### Scenario: Capacity information
- **GIVEN** 查询返回实例列表
- **THEN** 每个实例必须包含容量信息字段:
  - max_topic_number: 最大 Topic 数
  - max_partition_number: 最大分区数
  - topic_num: 当前 Topic 数量
  - partition_number: 当前分区数

#### Scenario: Billing information
- **GIVEN** 查询返回实例列表
- **THEN** 每个实例必须包含计费信息字段:
  - create_time: 创建时间(Unix 时间戳)
  - expire_time: 过期时间(Unix 时间戳)
  - renew_flag: 续费标识
  - public_network: 公网带宽
  - public_network_charge_type: 公网计费模式

#### Scenario: Health status
- **GIVEN** 查询返回实例列表
- **THEN** 每个实例必须包含健康状态字段:
  - healthy: 健康状态(Integer)
  - healthy_message: 健康状态信息

#### Scenario: Additional metadata
- **GIVEN** 查询返回实例列表
- **THEN** 每个实例可以包含额外的元数据字段:
  - version: Kafka 版本号
  - features: 功能特性列表
  - tags: 标签列表(嵌套对象)
  - is_internal: 是否内部客户
  - cvm: 售卖类型
  - rebalance_time: 升级配置时间

### Requirement: Code Implementation Standards
系统实现必须遵循项目代码规范和参考实现模式。

#### Scenario: Reference implementation pattern
- **GIVEN** 实现 CKafka instances v2 datasource
- **THEN** 代码结构必须参考 `data_source_tc_igtm_instance_list.go`
- **AND** 使用相同的代码组织模式
- **AND** 遵循相同的命名约定

#### Scenario: Parameter retrieval pattern
- **GIVEN** 从 ResourceData 获取参数
- **THEN** 必须使用 `d.GetOk()` 模式
- **AND** 格式为: `if v, ok := d.GetOk("param_name"); ok { ... }`
- **AND** 不使用直接的 `d.Get()` 调用

#### Scenario: Service layer integration
- **GIVEN** 实现数据查询逻辑
- **THEN** 必须通过 CkafkaService 调用 API
- **AND** 使用 `resource.Retry()` 实现重试逻辑
- **AND** 设置超时为 `tccommon.ReadRetryTimeout`

#### Scenario: Error handling
- **GIVEN** API 调用可能失败
- **THEN** 必须使用 `tccommon.RetryError()` 包装错误
- **AND** 添加适当的日志记录
- **AND** 使用 `defer tccommon.LogElapsed()` 记录耗时

#### Scenario: Resource ID generation
- **GIVEN** datasource 查询成功
- **THEN** 必须使用 `helper.BuildToken()` 或 `helper.DataResourceIdsHash()` 生成资源 ID
- **AND** 调用 `d.SetId()` 设置资源 ID

#### Scenario: Data mapping
- **GIVEN** API 返回实例列表
- **THEN** 必须检查字段是否为 nil 后再访问
- **AND** 使用 `helper.String()`, `helper.Int()` 等辅助函数
- **AND** 正确处理嵌套对象(如 vip_list, tags)

### Requirement: Documentation
系统必须提供完整的使用文档。

#### Scenario: Basic usage example
- **GIVEN** 用户查阅文档
- **THEN** 文档必须包含基本使用示例
- **AND** 示例展示如何查询所有实例
- **AND** 示例展示如何使用 filters 过滤

#### Scenario: Parameter documentation
- **GIVEN** 用户查阅参数说明
- **THEN** 文档必须列出所有输入参数
- **AND** 每个参数包含类型、是否必需、默认值、描述
- **AND** filters 参数说明支持的 filter.name 值

#### Scenario: Attribute documentation
- **GIVEN** 用户查阅返回属性
- **THEN** 文档必须列出所有输出属性
- **AND** 每个属性包含类型和描述
- **AND** 嵌套对象展示完整结构
