## Why

用户需要通过 Terraform 查询云拨测（CAT）的维度标签值，用于在拨测指标查询时获取可用的维度标签值列表。当前 Provider 已经支持 `tencentcloud_cat_metric_data` 数据源来查询拨测指标数据，但在构建维度标签过滤条件时，用户无法预先获知可用的维度标签值（如 host、errorInfo、area、operator、taskId 等），导致难以正确构建查询条件。

新增 `tencentcloud_cat_probe_metric_tag_values` 数据源可以解决这一问题，让用户能够在 Terraform 配置中动态查询并使用维度标签值。

## What Changes

- 新增 Data Source: `tencentcloud_cat_probe_metric_tag_values`
- 实现对 CAT API `DescribeProbeMetricTagValues` 接口的调用，查询云拨测维度标签值
- 支持以下输入参数（均为可选）：
  - `analyze_task_type`: 分析任务类型（网络质量/页面性能/端口性能/文件传输/音视频体验）
  - `key`: 维度标签值（host/errorInfo/area/operator/taskId）
  - `filter`: 过滤条件，支持正则匹配
  - `filters`: 过滤条件数组
  - `time_range`: 时间范围
- 返回属性：
  - `tag_value_set`: 标签值序列化后的字符串

## Capabilities

### New Capabilities
- `cat-probe-metric-tag-values-datasource`: 查询云拨测维度标签值的数据源，支持通过分析任务类型、维度标签、过滤条件和时间范围查询可用的维度标签值列表

### Modified Capabilities
<!-- 无需修改的现有能力 -->

## Impact

- **新增能力**: 云拨测维度标签值查询
- **受影响的服务**: CAT (tencentcloud/services/cat)
- **新增文件**:
  - `tencentcloud/services/cat/data_source_tc_cat_probe_metric_tag_values.go`
  - `tencentcloud/services/cat/data_source_tc_cat_probe_metric_tag_values.md`
  - `tencentcloud/services/cat/data_source_tc_cat_probe_metric_tag_values_test.go`
  - Provider 注册代码需要添加此 data source
  - service 层需要新增 `DescribeCatProbeMetricTagValuesByFilter` 方法
- **API 依赖**: 
  - CAT API v20180409: `DescribeProbeMetricTagValues`
  - 包名: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cat/v20180409`
- **兼容性**: 无破坏性变更，纯新增功能
