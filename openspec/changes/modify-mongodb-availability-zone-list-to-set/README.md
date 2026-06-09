# 提案: 修改 MongoDB availability_zone_list 字段为 Set 类型

> **状态**: 📝 Proposed (待审核)  
> **类型**: 🔨 Breaking Change  
> **优先级**: 🔵 Medium  
> **预计工时**: 7-12 小时

## 📋 快速概览

将 `tencentcloud_mongodb_instance` 和 `tencentcloud_mongodb_sharding_instance` 资源中的 `availability_zone_list` 字段从 `TypeList` 改为 `TypeSet`,以更准确地表达其无序集合的语义,并避免因顺序变化导致的误判 diff。

## 🎯 核心变更

| 资源 | 字段 | 变更前 | 变更后 |
|------|------|--------|--------|
| `tencentcloud_mongodb_instance` | `availability_zone_list` | `schema.TypeList` | `schema.TypeSet` |
| `tencentcloud_mongodb_sharding_instance` | `availability_zone_list` | `schema.TypeList` | `schema.TypeSet` |

## ✨ 优势

- ✅ 避免因可用区顺序变化产生的误判 diff
- ✅ 自动去重,防止配置错误
- ✅ 更准确的语义表达(无序集合)
- ✅ 符合 Terraform 最佳实践
- ✅ 与项目中其他类似字段保持一致

## ⚠️ 破坏性影响

这是一个**破坏性变更**,会影响现有 Terraform 状态文件:

1. 用户升级 Provider 后首次执行 `terraform plan` 会看到状态变更
2. 状态文件格式会从 List 转换为 Set
3. 建议在主版本更新(如 v2.0.0)中发布

## 📂 受影响的文件

### 代码文件
- `tencentcloud/services/mongodb/resource_tc_mongodb_instance.go`
- `tencentcloud/services/mongodb/resource_tc_mongodb_sharding_instance.go`
- `tencentcloud/services/mongodb/resource_tc_mongodb_readonly_instance.go`

### 文档文件
- `tencentcloud/services/mongodb/resource_tc_mongodb_instance.md`
- `tencentcloud/services/mongodb/resource_tc_mongodb_sharding_instance.md`
- `CHANGELOG.md`

### 测试文件
- 现有测试用例无需修改(兼容)
- 建议添加 Set 特性测试用例

## 📖 详细文档

- [提案详情 (proposal.md)](./proposal.md) - 完整的变更说明和影响分析
- [任务清单 (tasks.md)](./tasks.md) - 详细的实现步骤和检查清单

## 🚀 实施步骤

### Phase 1: 代码修改
1. 修改三个资源文件的 Schema 定义
2. 更新数据读取/写入逻辑
3. 添加去重逻辑
4. 执行 `go fmt` 格式化

### Phase 2: 测试验证
1. 编译验证
2. 运行单元测试
3. 运行验收测试
4. 状态迁移测试

### Phase 3: 文档更新
1. 更新资源文档
2. 生成网站文档
3. 更新 CHANGELOG
4. 创建迁移指南

### Phase 4: 发布
1. 确定版本号(建议主版本升级)
2. 准备发布说明
3. 标注 BREAKING CHANGE
4. 通知社区

## 📝 示例代码变更

### Schema 定义变更
```diff
 "availability_zone_list": {
-    Type:     schema.TypeList,
+    Type:     schema.TypeSet,
     Optional: true,
     Computed: true,
     Elem: &schema.Schema{
         Type: schema.TypeString,
     },
     Description: "...",
 }
```

### 数据读取变更
```diff
 if v, ok := d.GetOk("availability_zone_list"); ok {
-    availabilityZoneList := helper.InterfacesStringsPoint(v.([]interface{}))
+    availabilityZoneList := helper.InterfacesStringsPoint(v.(*schema.Set).List())
     value.FieldByName("AvailabilityZoneList").Set(reflect.ValueOf(availabilityZoneList))
 }
```

### 数据写入变更(添加去重)
```diff
 if len(replicateSets) > 0 {
     var hiddenZone string
-    availabilityZoneList := make([]string, 0, 3)
+    availabilityZoneMap := make(map[string]bool)
     for _, replicate := range replicateSets[0].Nodes {
         itemZone := *replicate.Zone
         if *replicate.Hidden {
             hiddenZone = itemZone
         }
-        availabilityZoneList = append(availabilityZoneList, itemZone)
+        availabilityZoneMap[itemZone] = true
+    }
+    
+    availabilityZoneList := make([]string, 0, len(availabilityZoneMap))
+    for zone := range availabilityZoneMap {
+        availabilityZoneList = append(availabilityZoneList, zone)
     }
     _ = d.Set("hidden_zone", hiddenZone)
     _ = d.Set("availability_zone_list", availabilityZoneList)
 }
```

## 🔍 迁移示例

用户配置文件保持不变:
```hcl
resource "tencentcloud_mongodb_instance" "example" {
  instance_name          = "example"
  memory                 = 4
  volume                 = 100
  engine_version         = "MONGO_40_WT"
  machine_type           = "HIO10G"
  available_zone         = "ap-guangzhou-3"
  availability_zone_list = ["ap-guangzhou-3", "ap-guangzhou-4", "ap-guangzhou-6"]
  # ... 其他配置
}
```

升级 Provider 后执行 `terraform plan` 可能会看到:
```
  ~ resource "tencentcloud_mongodb_instance" "example" {
      ~ availability_zone_list = [
          - "ap-guangzhou-3",
          - "ap-guangzhou-4",
          - "ap-guangzhou-6",
        ] -> (known after apply)
      # ... 其他字段无变化
    }
```

这是正常的状态类型转换,执行 `terraform apply` 后不会有实际资源变更。

## ❓ FAQ

**Q: 为什么要做这个变更?**  
A: List 类型会因可用区顺序变化产生误判 diff,Set 类型更符合语义且避免此问题。

**Q: 这会影响我的现有配置吗?**  
A: HCL 配置语法不变,但首次升级会触发状态格式转换。

**Q: 是否必须立即升级?**  
A: 不是必须的,但升级后可以避免误判 diff 问题。

**Q: 如何回滚?**  
A: 可以降级到旧版本 Provider,状态文件会自动适配。

**Q: 何时发布?**  
A: 建议在下一个主版本更新(如 v2.0.0)中发布。

## 👥 审核人

- [ ] MongoDB 资源维护者
- [ ] Provider 架构负责人
- [ ] 技术文档团队

## 🔗 相关链接

- [Terraform Schema 类型文档](https://developer.hashicorp.com/terraform/plugin/sdkv2/schemas/schema-types)
- [项目贡献指南](../../CONTRIBUTING.md)
- [代码审查清单](../../docs/code-review-checklist.md)

---

**创建时间**: 2026-03-18  
**创建人**: Terraform Provider 开发团队  
**提案ID**: `modify-mongodb-availability-zone-list-to-set`
