## Context

当前 `tencentcloud_clb_instances` 数据源实现于早期版本,只返回了 CLB 实例的基础信息字段。随着腾讯云 CLB 产品能力的演进,`DescribeLoadBalancers` API 的 `LoadBalancer` 结构体新增了大量字段,包括:

**当前已支持的字段**(约 25 个):
- 基础信息: `clb_id`, `clb_name`, `network_type`, `status`, `create_time`, `status_time`, `project_id`
- 网络配置: `vpc_id`, `subnet_id`, `clb_vips`, `security_groups`
- 区域信息: `zone_id`, `zone`, `zone_name`, `zone_region`, `local_zone`, `zones`
- 跨域绑定: `target_region_info_region`, `target_region_info_vpc_id`
- IP 相关: `address_ip_version`, `vip_isp`, `numerical_vpc_id`
- 网络属性: `internet_charge_type`, `internet_bandwidth_max_out`
- 标签和集群: `tags`, `cluster_id`

**API 返回但未支持的字段**(约 40 个):
- 类型标识: `Forward` (区分传统型/应用型)
- 域名相关: `Domain`, `LoadBalancerDomain`
- 高防与隔离: `OpenBgp`, `Snat`, `Isolation`, `IsolatedTime`
- 计费相关: `ChargeType`, `PrepaidAttributes`, `ExpireTime`
- 日志配置: `LogSetId`, `LogTopicId`, `HealthLogSetId`, `HealthLogTopicId`
- IPv6 配置: `AddressIPv6`, `IPv6Mode`, `MixIpTarget`
- 高级特性: `SlaType` (性能容量型规格), `SnatPro`, `SnatIps`
- 封禁状态: `IsBlock`, `IsBlockTime`, `LocalBgp`
- 集群与标签: `ClusterIds` (多集群), `ClusterTag`, `NfvInfo`
- 独占集群: `ExclusiveCluster`, `Exclusive`
- 配置与扩展: `ConfigId`, `LoadBalancerPassToTarget`, `ExtraInfo`, `IsDDos`
- 属性标志: `AttributeFlags` (DeleteProtect, WAF等20+ 种属性)
- 网络出口: `Egress`
- Endpoint 关联: `AssociateEndpoint`
- 后端服务数: `TargetCount`
- 可用区亲和: `AvailableZoneAffinityInfo`
- 备可用区: `BackupZoneSet`
- 废弃字段: `Log` (已标记 Deprecated)
- 私有网络: `AnycastZone`

## Goals / Non-Goals

**Goals:**
- 补全所有 `DescribeLoadBalancers` API 返回的 LoadBalancer 字段到数据源 schema
- 保持 100% 向后兼容,不影响现有用户配置
- 正确处理所有嵌套对象和 nil 值
- 提供清晰的文档说明每个字段的含义和使用场景

**Non-Goals:**
- 不修改任何过滤参数 (input schema 保持不变)
- 不改变数据源的查询逻辑
- 不添加新的服务层方法
- 不支持已废弃的 `Log` 字段 (SDK 已标记 Deprecated)

## Decisions

### Decision 1: 字段命名约定

**选择**: 使用蛇形命名法(snake_case),字段名尽量与 SDK 字段名保持一致但转换为小写下划线格式

**理由**:
- 符合 Terraform 社区惯例
- 与现有字段命名风格一致
- 便于用户从 API 文档查找对应字段

**示例映射**:
- SDK: `ChargeType` → Terraform: `charge_type`
- SDK: `SlaType` → Terraform: `sla_type`  
- SDK: `PrepaidAttributes` → Terraform: `prepaid_attributes` (嵌套对象)
- SDK: `AttributeFlags` → Terraform: `attribute_flags` (字符串数组)

### Decision 2: 嵌套对象处理策略

**选择**: 所有嵌套对象平铺为独立的 Computed 字段,使用前缀区分

**理由**:
- 保持与现有实现一致 (如 `target_region_info_*`)
- 避免引入复杂的嵌套 schema 定义
- 更符合 Terraform 数据源的使用习惯
- 减少 nil 检查的复杂度

**对比方案**:
- **方案 A** (采用): 平铺字段
  ```hcl
  prepaid_period = 12
  prepaid_renew_flag = "NOTIFY_AND_AUTO_RENEW"
  ```
- **方案 B** (未采用): 嵌套对象
  ```hcl
  prepaid_attributes = {
    period = 12
    renew_flag = "NOTIFY_AND_AUTO_RENEW"
  }
  ```
  缺点: 增加 schema 复杂度,不符合现有风格

### Decision 3: 复杂结构体的处理

**选择**: 对于复杂嵌套结构(如 `ExclusiveCluster`, `AvailableZoneAffinityInfo`),使用 JSON 字符串序列化

**理由**:
- 这些字段使用频率低
- 结构复杂,完全平铺会导致字段爆炸
- 用户可以通过 `jsondecode()` 函数解析

**示例**:
```hcl
exclusive_cluster = jsonencode({
  "L4Clusters": [...],
  "L7Clusters": [...],
  "ClassicalCluster": {...}
})
```

### Decision 4: 数组字段类型选择

**选择**: 
- 简单字符串数组使用 `schema.TypeList` + `schema.TypeString`
- 对象数组使用 `schema.TypeList` + `schema.TypeMap` 或 JSON 序列化

**理由**:
- 与现有 `clb_vips`, `security_groups`, `zones` 实现一致
- SDK helper 函数 `helper.StringsInterfaces()` 可直接使用

### Decision 5: Nil 值处理

**选择**: 所有可能为 nil 的字段在映射前检查,避免 panic

**模式**:
```go
if clbInstance.ChargeType != nil {
    mapping["charge_type"] = *clbInstance.ChargeType
}
if clbInstance.PrepaidAttributes != nil {
    mapping["prepaid_period"] = *clbInstance.PrepaidAttributes.Period
    mapping["prepaid_renew_flag"] = *clbInstance.PrepaidAttributes.RenewFlag
}
```

### Decision 6: 废弃字段处理

**选择**: 不添加 SDK 中标记为 `Deprecated` 的字段 (`Log`)

**理由**:
- 避免引入即将下线的 API 字段
- SDK 文档明确建议使用 `LogSetId` 和 `LogTopicId` 替代

## Risks / Trade-offs

### Risk 1: Schema 字段过多导致文档可读性下降
**缓解**: 
- 在文档中按功能分组(基础信息、网络配置、计费、日志、高级特性)
- 为每个字段提供清晰的描述和注意事项
- 标注哪些字段可能返回 null

### Risk 2: 新字段的 nil 处理遗漏导致运行时 panic
**缓解**:
- 所有指针字段访问前必须检查 nil
- 在测试用例中覆盖 nil 字段场景
- 参考现有代码的 nil 检查模式

### Risk 3: 嵌套对象序列化后不便于用户使用
**缓解**:
- 仅对极少使用的复杂结构使用 JSON 序列化
- 文档中提供 `jsondecode()` 使用示例
- 主要字段仍使用平铺方式

### Risk 4: 字段过多导致性能影响
**影响**: 极小,数据源查询本身已经获取完整的 LoadBalancer 对象,只是映射到 state 的开销
**缓解**: 无需特殊处理

## Implementation Notes

### 字段分组

**基础信息**(已支持 + 新增):
- 已有: `clb_id`, `clb_name`, `network_type`, `status`, `create_time`, `status_time`, `project_id`
- 新增: `forward`, `domain`, `load_balancer_domain`

**网络配置**(已支持 + 新增):
- 已有: `vpc_id`, `subnet_id`, `clb_vips`, `numerical_vpc_id`, `address_ip_version`, `vip_isp`
- 新增: `address_ipv6`, `ipv6_mode`, `mix_ip_target`, `anycast_zone`, `egress`, `local_bgp`

**计费相关**(新增):
- `charge_type`, `expire_time`
- `prepaid_period`, `prepaid_renew_flag`, `prepaid_cur_instance_deadline` (平铺 PrepaidAttributes)

**日志配置**(新增):
- `log_set_id`, `log_topic_id`, `health_log_set_id`, `health_log_topic_id`

**高防与隔离**(新增):
- `open_bgp`, `snat`, `snat_pro`, `snat_ips` (JSON)
- `isolation`, `isolated_time`, `is_block`, `is_block_time`

**性能与规格**(新增):
- `sla_type`, `exclusive`, `target_count`

**集群相关**(已支持 + 新增):
- 已有: `cluster_id`, `zones`
- 新增: `cluster_ids`, `cluster_tag`, `nfv_info`

**高级配置**(新增):
- `config_id`, `load_balancer_pass_to_target`, `is_ddos`
- `attribute_flags` (字符串数组)
- `associate_endpoint`
- `exclusive_cluster` (JSON)
- `extra_info` (JSON)

**可用区配置**(已支持 + 新增):
- 已有: `zone_id`, `zone`, `zone_name`, `zone_region`, `local_zone`, `zones`
- 新增: `backup_zone_set` (对象数组), `available_zone_affinity_info` (JSON)

### 代码修改范围

1. **Schema 定义**: `DataSourceTencentCloudClbInstances()` - 添加约 40 个新字段
2. **数据映射**: `dataSourceTencentCloudClbInstancesRead()` - 添加对应的字段赋值逻辑
3. **文档模板**: `data_source_tc_clb_instances.md` - 添加字段说明
4. **生成文档**: `website/docs/d/clb_instances.html.markdown` - 运行 `make doc`

### 测试策略

- 复用现有测试用例,验证向后兼容
- 在测试环境创建不同配置的 CLB 实例(公网/内网、预付费/后付费、不同规格)
- 验证所有新字段正确返回且处理 nil 值
