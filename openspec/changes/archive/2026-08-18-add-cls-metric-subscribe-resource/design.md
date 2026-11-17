## Context

CLS (日志服务) 提供指标订阅功能，可将日志主题中的指标数据采集并订阅到云监控。当前 Terraform Provider 缺少对该功能的管理，用户只能通过控制台或 API 手动操作。

本次新增 `tencentcloud_cls_metric_subscribe` 资源 (RESOURCE_KIND_GENERAL)，使用 CLS API v20201016 的 `CreateMetricSubscribe`、`DescribeMetricSubscribes`、`ModifyMetricSubscribe`、`DeleteMetricSubscribe` 四个接口实现完整 CRUD。

代码风格严格参考 `tencentcloud/services/igtm/resource_tc_igtm_strategy.go` 资源。

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_cls_metric_subscribe` 资源的完整 CRUD (Create/Read/Update/Delete)
- 支持 import 导入 (使用复合 ID `topicId#taskId`)
- 支持所有云 API 定义的参数: name、topic_id、namespace、metrics (嵌套)、instance_info (嵌套)、enable、task_id (computed)、status (computed)、create_time (computed)、update_time (computed)
- 在 provider.go 和 provider.md 中注册资源
- 生成对应的 .md 文档和单元测试文件

**Non-Goals:**
- 不实现数据源 (RESOURCE_KIND_DATASOURCE)
- 不修改任何现有资源的 schema
- 不涉及 website/ 目录文件的直接修改 (由 `make doc` 在收尾阶段生成)

## Decisions

### 1. 资源 ID 格式: 复合 ID `topicId#taskId`

**理由**: Create 接口出参为 `TaskId`，但 Read (DescribeMetricSubscribes) 接口入参需要 `TopicId` + `Filters` (按 taskId 过滤)。Delete 和 Modify 接口同时需要 `TopicId` 和 `TaskId`。因此使用复合 ID `topicId#taskId` (以 `tccommon.FILED_SP` 分隔)，在 Read/Update/Delete 中拆分获取两个值。

**import 说明**: 因为使用复合 ID，在 .md 文档 import 示例中需说明使用 `topicId#taskId` 格式。

### 2. Schema 字段设计

根据云 API 的请求/响应结构定义以下 schema:

**顶层字段:**
| 字段 | 类型 | Required | ForceNew | Computed | 说明 |
|------|------|----------|----------|----------|------|
| `name` | TypeString | Yes | No | No | 订阅任务名称 |
| `topic_id` | TypeString | Yes | Yes | No | 日志主题 ID (CRUD 均需此字段，ForceNew) |
| `namespace` | TypeString | Yes | No | No | 云产品命名空间 |
| `metrics` | TypeList | Yes | No | No | 指标配置信息 (嵌套) |
| `instance_info` | TypeList (MaxItems=1) | Yes | No | No | 实例配置信息 (嵌套) |
| `enable` | TypeInt | Optional | No | No | 任务开关: 1 暂停, 2 启用 |
| `task_id` | TypeString | Computed | No | Yes | 订阅任务 ID (创建后返回) |
| `status` | TypeInt | Computed | No | Yes | 运行状态: 0 创建中, 1 暂停, 2 运行中, 3 异常 |
| `create_time` | TypeInt | Computed | No | Yes | 创建时间 (秒级时间戳) |
| `update_time` | TypeInt | Computed | No | Yes | 更新时间 (秒级时间戳) |

**`metrics` 嵌套结构 (MetricConfig):**
| 字段 | 类型 | Required | 说明 |
|------|------|----------|------|
| `metric_name` | TypeString | Yes | 指标名称 |
| `periods` | TypeList(TypeInt) | Optional | 统计周期 (秒) |
| `metric_labels` | TypeList | Optional | 自定义指标标签 (嵌套) |

**`metrics.metric_labels` 嵌套结构 (MetricLabel):**
| 字段 | 类型 | Required | 说明 |
|------|------|----------|------|
| `key` | TypeString | Yes | 指标标签名称 |
| `value` | TypeString | Yes | 指标标签内容 |

**`instance_info` 嵌套结构 (InstanceConfig):**
| 字段 | 类型 | Required | 说明 |
|------|------|----------|------|
| `instance_dimension` | TypeList(TypeString) | Optional | 实例维度 |
| `instances` | TypeList | Optional | 实例值列表 (嵌套) |

**`instance_info.instances` 嵌套结构 (Instance):**
| 字段 | 类型 | Required | 说明 |
|------|------|----------|------|
| `values` | TypeList(TypeString) | Optional | 实例信息值列表 |

### 3. `topic_id` 设为 ForceNew

**理由**: `topic_id` 在 Create 中是必填入参，且 Read/Update/Delete 接口均需此字段。由于 `ModifyMetricSubscribe` 中 `TopicId` 是必填但语义为"指定要修改的任务所属主题"，不支持将任务迁移到其他主题。将 `topic_id` 设为 ForceNew，变更时重建资源。

### 4. Read 接口实现: 通过 DescribeMetricSubscribes + Filter 查询

**理由**: 云 API 没有提供单条查询接口 `DescribeMetricSubscribe`，只有列表查询接口 `DescribeMetricSubscribes`。Read 实现时使用 `TopicId` + `Filters` (Key=taskId, Values=[taskId]) 进行过滤查询，从返回的 `Datas` 列表中取第一条匹配记录。

分页参数 `Limit` 设为云 API 注释中的最大值 100，`Offset` 设为 0。

### 5. Update 实现: 使用 ModifyMetricSubscribe

Update 方法中，检查 `name`、`namespace`、metrics`、`instance_info`、`enable` 字段是否有变更。任一变更则调用 `ModifyMetricSubscribe` 接口 (传入 `TopicId` + `TaskId` + 变更字段)。

### 6. Service 层: 新增 `DescribeClsMetricSubscribeById`

在 `service_tencentcloud_cls.go` 中新增方法，封装 `DescribeMetricSubscribes` 调用逻辑 (带 `tccommon.ReadRetryTimeout` 重试)，通过 TopicId + taskId 过滤返回单条 `MetricSubscribeInfo`。

### 7. 错误处理与重试

- Create/Update/Delete: 使用 `tccommon.WriteRetryTimeout` 作为超时进行重试，失败时用 `tccommon.RetryError()` 包装错误
- Read: 使用 `tccommon.ReadRetryTimeout` 作为超时进行重试
- Create 完成后检查返回值 `TaskId` 是否为空，为空则返回 `NonRetryableError`
- Read 中若 API 返回空列表，先打印 `log.Printf("[CRUD] cls_metric_subscribe id=%s", d.Id())` 再 `d.SetId("")`

## Risks / Trade-offs

- **[Read 依赖列表查询接口]** → 通过 `Filters` 按 taskId 精确过滤，且 `Limit` 设为最大值 100 降低漏数据风险；若返回空列表则清理 state (SetId(""))
- **[嵌套结构较深]** → metrics 内含 metric_labels，instance_info 内含 instances，层级较深。严格按 igtm_strategy 风格逐层展开，保证 set/read 逻辑清晰
- **[enable 字段在 Create 中无对应入参]** → Create 接口 (`CreateMetricSubscribeRequest`) 不包含 `enable` 参数，新建任务默认状态由云 API 决定。`enable` 字段在 schema 中设为 Optional，仅在 Update 时通过 `ModifyMetricSubscribe` 设置
