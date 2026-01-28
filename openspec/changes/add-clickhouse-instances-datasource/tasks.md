# Implementation Tasks

## 1. Schema 定义与数据源创建
- [x] 1.1 创建 `data_source_tc_clickhouse_instances.go` 文件
- [x] 1.2 定义数据源函数 `DataSourceTencentCloudClickhouseInstances()`
- [x] 1.3 定义输入参数 Schema
  - [x] `instance_id` (Optional, String) - 实例 ID 精确搜索
  - [x] `instance_name` (Optional, String) - 实例名称模糊搜索
  - [x] `tags` (Optional, Map) - 标签过滤
  - [x] `vips` (Optional, List of String) - VIP 地址列表
  - [x] `is_simple` (Optional, Bool) - 是否返回简化信息
  - [x] `result_output_file` (Optional, String) - 结果输出文件
- [x] 1.4 定义输出参数 Schema
  - [x] `instance_list` (Computed, List) - 实例列表主结构
- [x] 1.5 定义 `instance_list` 嵌套 Schema - 基础信息
  - [x] `instance_id`, `instance_name`, `status`, `status_desc`, `version`
  - [x] `region`, `zone`, `region_id`, `region_desc`, `zone_desc`
- [x] 1.6 定义 `instance_list` 嵌套 Schema - 网络信息
  - [x] `vpc_id`, `subnet_id`, `access_info`, `eip`, `ch_proxy_vip`
- [x] 1.7 定义 `instance_list` 嵌套 Schema - 计费信息
  - [x] `pay_mode`, `create_time`, `expire_time`, `renew_flag`
- [x] 1.8 定义 `instance_list` 嵌套 Schema - 节点配置
  - [x] `master_summary` (TypeList, 嵌套对象)
  - [x] `common_summary` (TypeList, 嵌套对象)
- [x] 1.9 定义 `master_summary` 和 `common_summary` 子 Schema
  - [x] `spec`, `node_size`, `core`, `memory`, `disk`, `disk_type`, `disk_desc`
  - [x] `attach_cbs_spec` (嵌套对象: disk_type, disk_size, disk_count)
  - [x] `sub_product_type`, `spec_core`, `spec_memory`, `disk_count`, `max_disk_size`, `encrypt`
- [x] 1.10 定义 `instance_list` 嵌套 Schema - 其他字段
  - [x] `ha`, `ha_zk`, `is_elastic`, `kind`, `monitor`
  - [x] `tags` (TypeList, TagKey/TagValue)
  - [x] `has_cls_topic`, `cls_topic_id`, `cls_log_set_id`
  - [x] `enable_xml_config`, `cos_bucket_name`
  - [x] `can_attach_cbs`, `can_attach_cbs_lvm`, `can_attach_cos`
  - [x] `components` (TypeList, Name/Version)
  - [x] `upgrade_versions`, `instance_state_info`, `flow_msg`

## 2. Read 函数实现
- [x] 2.1 实现 `dataSourceTencentCloudClickhouseInstancesRead()` 函数
- [x] 2.2 获取 LogId 和 Context
- [x] 2.3 初始化 CDWCH 客户端和 Service
- [x] 2.4 构建 API 请求参数
  - [x] 从 schema 读取 `instance_id` → `request.SearchInstanceId`
  - [x] 从 schema 读取 `instance_name` → `request.SearchInstanceName`
  - [x] 从 schema 读取 `vips` → `request.Vips`
  - [x] 从 schema 读取 `is_simple` → `request.IsSimple`
  - [x] 从 schema 读取 `tags` 并转换为 `request.SearchTags`
- [x] 2.5 调用 SDK API `DescribeInstancesNew`
  - [x] 使用 `resource.Retry` 包装调用
  - [x] 添加速率限制检查 `ratelimit.Check()`
  - [x] 添加错误处理和日志记录
- [x] 2.6 处理 API 响应
  - [x] 检查 response 非空
  - [x] 提取 `InstancesList`
- [x] 2.7 数据转换：API 响应 → Terraform Schema
  - [x] 遍历 `InstancesList`
  - [x] 转换每个实例的基础字段
  - [x] 转换嵌套对象 `MasterSummary`
  - [x] 转换嵌套对象 `CommonSummary`
  - [x] 转换 `Tags` 数组
  - [x] 转换 `Components` 数组
  - [x] 转换 `InstanceStateInfo` 对象
- [x] 2.8 设置 Terraform State
  - [x] `d.SetId()` - 使用时间戳或组合 ID
  - [x] `d.Set("instance_list", instanceList)`
- [x] 2.9 实现 `result_output_file` 输出
  - [x] 如果指定了输出文件，将结果序列化为 JSON
  - [x] 写入文件

## 3. 辅助函数
- [x] 3.1 实现 `flattenInstanceInfo()` - 将 SDK InstanceInfo 转换为 map
- [x] 3.2 实现 `flattenNodesSummary()` - 转换 MasterSummary/CommonSummary
- [x] 3.3 实现 `flattenAttachCBSSpec()` - 转换 AttachCBSSpec (已集成到 flattenNodesSummary 中)
- [x] 3.4 实现 `flattenTags()` - 转换 Tags 数组 (已集成到 flattenInstanceInfo 中)
- [x] 3.5 实现 `flattenComponents()` - 转换 Components 数组 (已集成到 flattenInstanceInfo 中)
- [x] 3.6 实现 `flattenInstanceStateInfo()` - 转换 InstanceStateInfo (已集成到 flattenInstanceInfo 中)
- [x] 3.7 实现安全的指针解引用辅助函数（处理 nil 值）- 通过 nil 检查实现

## 4. 测试用例
- [x] 4.1 创建 `data_source_tc_clickhouse_instances_test.go` 文件
- [x] 4.2 实现 `TestAccTencentCloudClickhouseInstancesDataSource_basic` 测试
  - [x] 测试不带任何过滤条件的查询
  - [x] 验证返回的实例列表非空
  - [x] 验证基础字段正确性
- [ ] 4.3 实现 `TestAccTencentCloudClickhouseInstancesDataSource_byId` 测试
  - [ ] 使用已知实例 ID 进行查询
  - [ ] 验证返回结果匹配
- [ ] 4.4 实现 `TestAccTencentCloudClickhouseInstancesDataSource_byName` 测试
  - [ ] 使用实例名称进行模糊查询
  - [ ] 验证返回结果包含匹配项
- [ ] 4.5 实现 `TestAccTencentCloudClickhouseInstancesDataSource_byTags` 测试
  - [ ] 使用标签进行过滤
  - [ ] 验证返回的实例都包含指定标签
- [x] 4.6 创建测试用的 Terraform 配置模板
  - [x] 基础查询配置
  - [ ] 带过滤条件的查询配置

## 5. 文档编写
- [x] 5.1 创建 `data_source_tc_clickhouse_instances.md` 文件
- [x] 5.2 编写资源描述和用途说明
- [x] 5.3 编写使用示例
  - [x] 基础查询示例
  - [x] 按 ID 查询示例
  - [x] 按名称查询示例
  - [x] 按标签查询示例
  - [x] 输出到文件示例
- [x] 5.4 编写参数说明（Argument Reference）
  - [x] 输入参数详细说明
  - [x] 可选/必选标注
  - [x] 参数类型和约束
- [x] 5.5 编写属性说明（Attributes Reference）
  - [x] `instance_list` 结构说明
  - [x] 所有输出字段的详细说明
- [x] 5.6 添加注意事项和最佳实践

## 6. 注册与集成
- [x] 6.1 修改 `provider.go`
- [x] 6.2 在数据源映射中注册数据源
  - [x] 添加 `"tencentcloud_clickhouse_instances": DataSourceTencentCloudClickhouseInstances()`

## 7. 代码质量
- [x] 7.1 运行 `go fmt` 格式化代码
- [x] 7.2 运行 `make lint` 检查代码质量
- [x] 7.3 修复所有 linter 警告（保留已存在的弃用警告）
- [x] 7.4 确保代码遵循项目规范
  - [x] 函数命名符合规范
  - [x] 导入别名正确（tccommon, helper 等）
  - [x] 日志记录完整
  - [x] 错误处理正确

## 8. OpenSpec 规范
- [ ] 8.1 创建 spec delta 文件 `specs/clickhouse-datasource/spec.md` (已存在)
- [ ] 8.2 定义 ADDED Requirements (已存在)
  - [ ] 数据源查询能力要求
  - [ ] 过滤条件支持要求
  - [ ] 数据转换要求
- [ ] 8.3 为每个 Requirement 添加 Scenario (已存在)
- [ ] 8.4 运行 `openspec validate add-clickhouse-instances-datasource --strict`
- [ ] 8.5 解决所有验证错误

## 9. 验收测试
- [ ] 9.1 设置测试环境变量
  - [ ] `TENCENTCLOUD_SECRET_ID`
  - [ ] `TENCENTCLOUD_SECRET_KEY`
  - [ ] `TF_ACC=1`
- [ ] 9.2 运行单元测试 `make test`
- [ ] 9.3 运行验收测试 `TF_ACC=1 go test -v -run TestAccTencentCloudClickhouseInstancesDataSource`
- [ ] 9.4 验证所有测试场景通过
- [ ] 9.5 测试真实 API 调用
  - [ ] 查询现有实例列表
  - [ ] 验证返回数据完整性
  - [ ] 验证过滤条件生效

## 10. 最终验证
- [x] 10.1 编译整个 provider `make build`
- [x] 10.2 验证没有编译错误
- [ ] 10.3 验证文档生成 `make doc`
- [ ] 10.4 创建示例 Terraform 配置并手动测试
- [ ] 10.5 验证 OpenSpec 最终状态
- [ ] 10.6 准备提交说明和变更日志
- [ ] 2.3 初始化 CDWCH 客户端和 Service
- [ ] 2.4 构建 API 请求参数
  - [ ] 从 schema 读取 `instance_id` → `request.SearchInstanceId`
  - [ ] 从 schema 读取 `instance_name` → `request.SearchInstanceName`
  - [ ] 从 schema 读取 `vips` → `request.Vips`
  - [ ] 从 schema 读取 `is_simple` → `request.IsSimple`
  - [ ] 从 schema 读取 `tags` 并转换为 `request.SearchTags`
- [ ] 2.5 调用 SDK API `DescribeInstancesNew`
  - [ ] 使用 `resource.Retry` 包装调用
  - [ ] 添加速率限制检查 `ratelimit.Check()`
  - [ ] 添加错误处理和日志记录
- [ ] 2.6 处理 API 响应
  - [ ] 检查 response 非空
  - [ ] 提取 `InstancesList`
- [ ] 2.7 数据转换：API 响应 → Terraform Schema
  - [ ] 遍历 `InstancesList`
  - [ ] 转换每个实例的基础字段
  - [ ] 转换嵌套对象 `MasterSummary`
  - [ ] 转换嵌套对象 `CommonSummary`
  - [ ] 转换 `Tags` 数组
  - [ ] 转换 `Components` 数组
  - [ ] 转换 `InstanceStateInfo` 对象
- [ ] 2.8 设置 Terraform State
  - [ ] `d.SetId()` - 使用时间戳或组合 ID
  - [ ] `d.Set("instance_list", instanceList)`
- [ ] 2.9 实现 `result_output_file` 输出
  - [ ] 如果指定了输出文件，将结果序列化为 JSON
  - [ ] 写入文件

## 3. 辅助函数
- [ ] 3.1 实现 `flattenInstanceInfo()` - 将 SDK InstanceInfo 转换为 map
- [ ] 3.2 实现 `flattenNodesSummary()` - 转换 MasterSummary/CommonSummary
- [ ] 3.3 实现 `flattenAttachCBSSpec()` - 转换 AttachCBSSpec
- [ ] 3.4 实现 `flattenTags()` - 转换 Tags 数组
- [ ] 3.5 实现 `flattenComponents()` - 转换 Components 数组
- [ ] 3.6 实现 `flattenInstanceStateInfo()` - 转换 InstanceStateInfo
- [ ] 3.7 实现安全的指针解引用辅助函数（处理 nil 值）

## 4. 测试用例
- [ ] 4.1 创建 `data_source_tc_clickhouse_instances_test.go` 文件
- [ ] 4.2 实现 `TestAccTencentCloudClickhouseInstancesDataSource_basic` 测试
  - [ ] 测试不带任何过滤条件的查询
  - [ ] 验证返回的实例列表非空
  - [ ] 验证基础字段正确性
- [ ] 4.3 实现 `TestAccTencentCloudClickhouseInstancesDataSource_byId` 测试
  - [ ] 使用已知实例 ID 进行查询
  - [ ] 验证返回结果匹配
- [ ] 4.4 实现 `TestAccTencentCloudClickhouseInstancesDataSource_byName` 测试
  - [ ] 使用实例名称进行模糊查询
  - [ ] 验证返回结果包含匹配项
- [ ] 4.5 实现 `TestAccTencentCloudClickhouseInstancesDataSource_byTags` 测试
  - [ ] 使用标签进行过滤
  - [ ] 验证返回的实例都包含指定标签
- [ ] 4.6 创建测试用的 Terraform 配置模板
  - [ ] 基础查询配置
  - [ ] 带过滤条件的查询配置

## 5. 文档编写
- [ ] 5.1 创建 `data_source_tc_clickhouse_instances.md` 文件
- [ ] 5.2 编写资源描述和用途说明
- [ ] 5.3 编写使用示例
  - [ ] 基础查询示例
  - [ ] 按 ID 查询示例
  - [ ] 按名称查询示例
  - [ ] 按标签查询示例
  - [ ] 输出到文件示例
- [ ] 5.4 编写参数说明（Argument Reference）
  - [ ] 输入参数详细说明
  - [ ] 可选/必选标注
  - [ ] 参数类型和约束
- [ ] 5.5 编写属性说明（Attributes Reference）
  - [ ] `instance_list` 结构说明
  - [ ] 所有输出字段的详细说明
- [ ] 5.6 添加注意事项和最佳实践

## 6. 注册与集成
- [ ] 6.1 修改 `extension_cdwch.go`
- [ ] 6.2 在 `GetResources()` 中注册数据源
  - [ ] 添加 `"tencentcloud_clickhouse_instances": DataSourceTencentCloudClickhouseInstances()`

## 7. 代码质量
- [ ] 7.1 运行 `go fmt` 格式化代码
- [ ] 7.2 运行 `make lint` 检查代码质量
- [ ] 7.3 修复所有 linter 警告
- [ ] 7.4 确保代码遵循项目规范
  - [ ] 函数命名符合规范
  - [ ] 导入别名正确（tccommon, helper 等）
  - [ ] 日志记录完整
  - [ ] 错误处理正确

## 8. OpenSpec 规范
- [ ] 8.1 创建 spec delta 文件 `specs/clickhouse-datasource/spec.md`
- [ ] 8.2 定义 ADDED Requirements
  - [ ] 数据源查询能力要求
  - [ ] 过滤条件支持要求
  - [ ] 数据转换要求
- [ ] 8.3 为每个 Requirement 添加 Scenario
- [ ] 8.4 运行 `openspec validate add-clickhouse-instances-datasource --strict`
- [ ] 8.5 解决所有验证错误

## 9. 验收测试
- [ ] 9.1 设置测试环境变量
  - [ ] `TENCENTCLOUD_SECRET_ID`
  - [ ] `TENCENTCLOUD_SECRET_KEY`
  - [ ] `TF_ACC=1`
- [ ] 9.2 运行单元测试 `make test`
- [ ] 9.3 运行验收测试 `TF_ACC=1 go test -v -run TestAccTencentCloudClickhouseInstancesDataSource`
- [ ] 9.4 验证所有测试场景通过
- [ ] 9.5 测试真实 API 调用
  - [ ] 查询现有实例列表
  - [ ] 验证返回数据完整性
  - [ ] 验证过滤条件生效

## 10. 最终验证
- [ ] 10.1 编译整个 provider `make build`
- [ ] 10.2 验证没有编译错误
- [ ] 10.3 验证文档生成 `make doc`
- [ ] 10.4 创建示例 Terraform 配置并手动测试
- [ ] 10.5 验证 OpenSpec 最终状态
- [ ] 10.6 准备提交说明和变更日志
