# Change: Add tencentcloud_ckafka_instances_v2 Datasource

## Why
当前的 `tencentcloud_ckafka_instances` datasource 已存在,但需要增加一个 v2 版本以更好地对齐 DescribeInstancesDetail API 接口的完整功能,并参考 `tencentcloud_igtm_instance_list` 的实现模式,提供更规范的过滤器机制和更完整的实例详情字段。

新的 datasource 将提供:
- 更标准的 filters 过滤器结构(name/value/fuzzy)
- 更完整的实例详情字段映射
- 更符合项目规范的代码结构

## What Changes
- 新增 `tencentcloud_ckafka_instances_v2` datasource
- 基于 DescribeInstancesDetail API 实现
- 参考 `tencentcloud_igtm_instance_list` 的代码模式和风格
- 支持标准化的 filters 过滤器
- 返回完整的实例详情列表

## Impact
- 受影响的 specs: 新增 `ckafka-instances-v2` capability
- 受影响的代码:
  - 新增文件: `tencentcloud/services/ckafka/data_source_tc_ckafka_instances_v2.go`
  - 修改文件: `tencentcloud/provider.go` (注册新 datasource)
  - 新增文件: `tencentcloud/services/ckafka/data_source_tc_ckafka_instances_v2.md` (文档)
- 无破坏性变更
- 不影响现有 `tencentcloud_ckafka_instances` datasource
