# Design: tencentcloud_ckafka_instances_v2 Datasource

## Context
需要实现一个新的 CKafka 实例查询 datasource,以提供更标准化的过滤器机制和更完整的字段映射。参考 `tencentcloud_igtm_instance_list` 的实现模式,确保代码风格和结构的一致性。

## Goals / Non-Goals

### Goals
- 实现 `tencentcloud_ckafka_instances_v2` datasource
- 使用标准化的 filters 过滤器结构
- 完整映射 DescribeInstancesDetail API 的所有字段
- 遵循项目代码规范和最佳实践
- 参考 `data_source_tc_igtm_instance_list.go` 的代码模式

### Non-Goals
- 不修改现有的 `tencentcloud_ckafka_instances` datasource
- 不新增或修改 SDK 接口
- 不添加任何破坏性变更
- 本阶段不实现全面的单元测试(可选实现)

## Decisions

### Decision 1: 使用 filters 参数结构
参考 igtm_instance_list,采用标准的 filters 参数结构:
```hcl
filters {
  name  = "InstanceId"
  value = ["ckafka-xxx", "ckafka-yyy"]
}
```

**Why**: 
- 与项目中其他 datasource 保持一致
- 提供灵活的过滤能力
- 符合 Terraform 最佳实践

**Alternatives considered**:
- 直接使用 instance_ids 参数: 不够灵活,无法支持多种过滤条件
- 使用 search_word 参数: 不够精确,查询结果不可预测

### Decision 2: 命名为 _v2 后缀
使用 `tencentcloud_ckafka_instances_v2` 作为资源名称。

**Why**:
- 避免与现有 `tencentcloud_ckafka_instances` 冲突
- 明确表示这是新版本的实现
- 允许两个版本并存,用户可以逐步迁移

**Alternatives considered**:
- 替换现有 datasource: 可能破坏现有用户配置
- 使用不同名称(如 instances_detail): 不够直观,命名不一致

### Decision 3: 完整字段映射
映射 DescribeInstancesDetail API 返回的所有有用字段,包括:
- 基本信息(InstanceId, InstanceName, Status)
- 网络配置(VpcId, SubnetId, ZoneId, VipList)
- 资源规格(Bandwidth, DiskSize, InstanceType)
- 容量信息(MaxTopicNumber, TopicNum, PartitionNumber)
- 计费信息(CreateTime, ExpireTime, RenewFlag)
- 健康状态(Healthy, HealthyMessage)
- 其他元数据(Version, Tags, Features)

**Why**:
- 提供完整的实例信息
- 减少用户需要调用其他 API 的次数
- 满足各种查询需求

### Decision 4: 参考 igtm_instance_list 实现模式
严格参考 `data_source_tc_igtm_instance_list.go` 的代码结构和风格。

**Why**:
- 保持代码库一致性
- 复用经过验证的模式
- 减少维护成本

**Key patterns to follow**:
- 使用 `d.GetOk()` 获取参数
- 使用 `resource.Retry()` 实现重试
- 使用 `helper.BuildToken()` 生成资源 ID
- 正确的错误处理和日志记录
- nil 检查和类型转换

## Implementation Details

### Schema Structure
```go
Schema: map[string]*schema.Schema{
    "filters": {
        Type:     schema.TypeList,
        Optional: true,
        Elem: &schema.Resource{
            Schema: map[string]*schema.Schema{
                "name":  { Type: schema.TypeString, Required: true },
                "value": { Type: schema.TypeSet, Required: true, Elem: &schema.Schema{Type: schema.TypeString} },
            },
        },
    },
    "instance_list": {
        Type:     schema.TypeList,
        Computed: true,
        Elem: &schema.Resource{
            Schema: map[string]*schema.Schema{
                "instance_id":   { Type: schema.TypeString, Computed: true },
                "instance_name": { Type: schema.TypeString, Computed: true },
                // ... 其他字段
            },
        },
    },
    "result_output_file": { Type: schema.TypeString, Optional: true },
}
```

### API Integration
```go
// 构建 filters 参数
if v, ok := d.GetOk("filters"); ok {
    filtersSet := v.([]interface{})
    tmpSet := make([]*ckafka.Filter, 0, len(filtersSet))
    for _, item := range filtersSet {
        filtersMap := item.(map[string]interface{})
        filter := ckafka.Filter{}
        if v, ok := filtersMap["name"].(string); ok && v != "" {
            filter.Name = helper.String(v)
        }
        if v, ok := filtersMap["value"]; ok {
            valueSet := v.(*schema.Set).List()
            for i := range valueSet {
                value := valueSet[i].(string)
                filter.Values = append(filter.Values, helper.String(value))
            }
        }
        tmpSet = append(tmpSet, &filter)
    }
    request.Filters = tmpSet
}

// 调用 API
reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
    response, e := ckafkaService.client.UseCkafkaClient().DescribeInstancesDetail(request)
    if e != nil {
        return tccommon.RetryError(e)
    }
    respData = response.Response.Result.InstanceList
    return nil
})
```

### Data Mapping Pattern
```go
instanceList := make([]map[string]interface{}, 0, len(respData))
for _, instance := range respData {
    instanceMap := map[string]interface{}{}
    
    if instance.InstanceId != nil {
        instanceMap["instance_id"] = instance.InstanceId
    }
    if instance.InstanceName != nil {
        instanceMap["instance_name"] = instance.InstanceName
    }
    // ... 映射其他字段
    
    // 处理嵌套对象
    if instance.VipList != nil {
        vipList := make([]map[string]interface{}, 0, len(instance.VipList))
        for _, vip := range instance.VipList {
            vipMap := map[string]interface{}{}
            if vip.Vip != nil {
                vipMap["vip"] = vip.Vip
            }
            if vip.Vport != nil {
                vipMap["vport"] = vip.Vport
            }
            vipList = append(vipList, vipMap)
        }
        instanceMap["vip_list"] = vipList
    }
    
    instanceList = append(instanceList, instanceMap)
}
_ = d.Set("instance_list", instanceList)
```

## Risks / Trade-offs

### Risk: 字段过多导致复杂性增加
**Mitigation**: 
- 所有字段都是 Computed,不影响用户配置
- 用户只需关注自己需要的字段
- 文档清晰说明每个字段的含义

### Risk: API 变更导致不兼容
**Mitigation**:
- 使用 nil 检查保护所有字段访问
- SDK 版本锁定在 go.mod 中
- 定期更新 SDK 并测试

### Trade-off: v2 后缀可能造成命名混淆
**Acceptance**:
- 清晰的文档说明两个版本的差异
- 推荐新用户使用 v2 版本
- 现有用户可以继续使用旧版本

## Migration Plan
不需要迁移计划,因为:
- 这是新增的 datasource,不影响现有资源
- 与 `tencentcloud_ckafka_instances` 并存
- 用户可以根据需要选择使用哪个版本

## Open Questions
- [ ] 是否需要添加分页支持?(API 支持 Offset/Limit)
  - 初步决定: 不在第一版实现,保持简单
- [ ] 是否需要支持 fuzzy 参数?
  - 初步决定: filters 暂不支持 fuzzy,保持与 API 原生行为一致
- [ ] 是否需要完整的测试覆盖?
  - 初步决定: 基本验收测试可选,优先保证代码质量
