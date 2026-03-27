## Why

当前 RabbitMQ VIP 实例的 update 逻辑存在限制，将几乎所有关键配置参数标记为不可变字段（如 `node_spec`、`node_num`、`storage_size`、`enable_public_access` 等）。这些参数在实际使用场景中是经常需要调整的，比如扩容节点、升级规格、调整存储大小、开启/关闭公网访问等。当前的硬性限制导致用户无法通过 Terraform 平滑地调整这些参数，必须手动在控制台操作或删除重建实例，这增加了运维复杂度和风险。

## What Changes

- **支持节点规格变更（node_spec）**: 允许修改 RabbitMQ 实例的节点规格，支持水平扩展和垂直扩展
- **支持节点数量变更（node_num）**: 允许动态调整实例节点数量，支持扩容和缩容
- **支持存储大小变更（storage_size）**: 允许调整单节点存储容量
- **支持公网访问开关（enable_public_access）**: 允许开启或关闭公网访问功能
- **支持带宽调整（band_width）**: 允许修改公网带宽配置
- **优化 update 逻辑**: 移除不必要的不可变字段限制，使用 API 返回的错误来判断哪些参数不支持修改
- **增强错误处理**: 提供更友好的错误提示，明确告知用户哪些参数在当前状态不支持修改

## Capabilities

### New Capabilities
- `rabbitmq-instance-update`: 增强 RabbitMQ 实例更新能力，支持更多配置参数的动态调整

### Modified Capabilities

## Impact

- **受影响的代码**:
  - `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go`: Update 函数逻辑重构
- **API 依赖**: 腾讯云 TDMQ ModifyRabbitMQVipInstance API 已支持这些参数的修改
- **向后兼容性**: 修改后完全向后兼容，现有配置不受影响，只是移除了之前过于严格的限制
- **测试影响**: 需要更新或新增测试用例以覆盖新增的更新场景
