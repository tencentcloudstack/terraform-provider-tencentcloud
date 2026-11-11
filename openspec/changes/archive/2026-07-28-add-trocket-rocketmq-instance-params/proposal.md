## Why

`tencentcloud_trocket_rocketmq_instance` 资源当前仅支持创建后付费（按量计费）实例，且无法指定部署可用区与最大主题数等关键创建参数。云 API `CreateInstance` 已提供 `PayMode`、`RenewFlag`、`TimeSpan`、`MaxTopicNum`、`ZoneIds` 入参，但 Terraform Schema 尚未暴露这些字段，导致用户无法通过 Terraform 创建预付费（包年包月）实例、配置自动续费、购买时长、自定义最大主题数及部署可用区。本次变更新增这五个参数以补齐云 API 能力。

## What Changes

- 在 `tencentcloud_trocket_rocketmq_instance` 资源 Schema 中新增五个 Optional 字段：
  - `pay_mode` (TypeInt, Optional, Default 0): 付费模式（0: 后付费；1: 预付费），默认值为 0。
  - `renew_flag` (TypeInt, Optional, Default 0): 预付费集群是否自动续费（0: 不自动续费；1: 自动续费），默认值为 0。
  - `time_span` (TypeInt, Optional, Default 1): 预付费集群的购买时长（单位：月），取值范围 1～60，默认值为 1。
  - `max_topic_num` (TypeInt, Optional): 最大可创建主题数。
  - `zone_ids` (TypeList of TypeInt, Optional): 部署可用区列表。
- Create 操作：将五个新增字段填充到 `CreateInstanceRequest` 中。
- Read 操作：从 `DescribeInstance` 响应中回填 `pay_mode`、`renew_flag`、`zone_ids`（注意云 API 返回的 `PayMode` 为字符串枚举 `POSTPAID`/`PREPAID`，需映射为 int 0/1；`renew_flag` 与 `zone_ids` 直接回填）。
- Update 操作：`max_topic_num` 和 `zone_ids` 支持通过 `ModifyInstance` 接口原地更新；`pay_mode`、`renew_flag`、`time_span` 为创建时参数，加入 `immutableArgs` 数组，变更时报错。
- 补充单元测试与资源文档。

## Capabilities

### New Capabilities
- `trocket-rocketmq-instance-params`: 为 `tencentcloud_trocket_rocketmq_instance` 资源新增付费模式、自动续费标志、购买时长、最大主题数、部署可用区列表等创建/更新参数，覆盖 Create、Read、Update 全链路。

### Modified Capabilities
<!-- 无：本次变更新增参数，不改变既有 spec 级别的需求行为。 -->

## Impact

- **Modified code**:
  - `tencentcloud/services/trocket/resource_tc_trocket_rocketmq_instance.go`: Schema 新增五个字段；Create/Read/Update 函数补充参数处理逻辑。
  - `tencentcloud/services/trocket/resource_tc_trocket_rocketmq_instance_test.go`: 补充新增参数的单元测试用例。
  - `tencentcloud/services/trocket/resource_tc_trocket_rocketmq_instance.md`: 更新示例用法文档。
- **APIs consumed**: `CreateInstance`、`DescribeInstance`、`ModifyInstance`（均已在 vendored SDK `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/trocket/v20230308/` 中）。
- **No breaking change**: 纯新增 Optional 字段，向后兼容；不修改既有 Schema 字段，不影响已有 state。
