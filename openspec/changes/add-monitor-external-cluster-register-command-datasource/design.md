# Design Document

## Context
用户在使用 `tencentcloud_monitor_external_cluster` 资源创建外部集群后,需要获取集群的注册命令,以便在其 Kubernetes 集群中执行安装操作。这个注册命令包含了安装 Prometheus Agent 所需的脚本和配置信息。

## Goals / Non-Goals

### Goals
- 提供简单的数据源查询接口获取注册命令
- 确保与现有 `tencentcloud_monitor_external_cluster` 资源无缝集成
- 遵循项目编码规范和模式

### Non-Goals
- 不负责执行注册命令
- 不负责验证注册命令的有效性
- 不提供命令执行状态的跟踪

## Decisions

### Decision 1: 数据源 vs 资源属性
**选择**: 独立的数据源而非在 `tencentcloud_monitor_external_cluster` 资源中添加 computed 属性

**理由**:
1. **关注点分离**: 资源管理生命周期(CRUD),数据源负责查询
2. **灵活性**: 用户可能只需要查询已存在的集群命令,而不创建新资源
3. **一致性**: 与 Terraform 社区的最佳实践保持一致
4. **向后兼容**: 不影响现有资源的结构

**替代方案**: 在资源中添加 computed 字段
- 优点: 使用更简单,一次操作即可
- 缺点: 违反单一职责原则,增加资源复杂度

### Decision 2: ID 格式
**选择**: 使用复合 ID `{instanceId}#{clusterId}`

**理由**:
1. **唯一性**: 两个参数共同唯一标识一个查询
2. **一致性**: 与 `tencentcloud_monitor_external_cluster` 资源的 ID 格式保持一致
3. **可解析性**: 使用 `tccommon.FILED_SP` 分隔符便于解析

### Decision 3: 参考实现选择
**选择**: 使用 `data_source_tc_igtm_instance_list.go` 作为参考模板

**理由**:
1. **代码风格一致**: 该文件展示了标准的数据源实现模式
2. **完整示例**: 包含参数处理、API 调用、错误处理、结果映射等完整流程
3. **最佳实践**: 包含了 `result_output_file`、日志记录等标准功能
4. **用户指定**: 用户明确要求参考此文件

### Decision 4: 服务层方法设计
**选择**: 在 `MonitorService` 中添加独立的 `DescribeExternalClusterRegisterCommand` 方法

**设计**:
```go
func (me *MonitorService) DescribeExternalClusterRegisterCommand(
    ctx context.Context, 
    instanceId string, 
    clusterId string,
) (result *monitor.DescribeExternalClusterRegisterCommandResponse, errRet error)
```

**理由**:
1. **可重用性**: 方法可被其他代码复用
2. **测试性**: 独立方法易于单元测试
3. **一致性**: 与现有服务层方法风格保持一致

## Risks / Trade-offs

### Risk 1: API 响应字段不确定
- **风险**: 接口文档可能不完整,实际返回字段需要调试确认
- **缓解**: 在实现阶段进行 API 调试,确认所有返回字段

### Risk 2: 注册命令时效性
- **风险**: 注册命令可能包含时效性信息(如 token),数据源查询时可能已过期
- **缓解**: 在文档中说明这一点,建议用户及时使用查询到的命令

### Trade-off: 独立数据源 vs 集成到资源
- **选择**: 独立数据源
- **优势**: 更灵活,关注点分离
- **劣势**: 用户需要两个配置块(resource + data source)
- **判断**: 灵活性和一致性优先于便利性

## Migration Plan
不适用 - 这是新增功能,不涉及迁移。

## Open Questions
1. API 返回的具体字段结构是什么?
   - **解决方案**: 在实现阶段通过 API 调试确认
   
2. 注册命令是否包含敏感信息?
   - **解决方案**: 如果包含,需在 Schema 中标记 `Sensitive: true`

3. 是否需要支持批量查询多个集群?
   - **当前决策**: 不支持,保持简单
   - **理由**: 单个查询更符合 Terraform 模式,批量需求可通过 `count` 或 `for_each` 实现
