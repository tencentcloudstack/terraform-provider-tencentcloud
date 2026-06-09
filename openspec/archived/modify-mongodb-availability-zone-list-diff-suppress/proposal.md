# MongoDB availability_zone_list 字段 DiffSuppressFunc 优化

## 📋 变更概述

**变更类型**: 增强 (Enhancement)  
**提案日期**: 2026-03-23  
**提案人**: AI Assistant  
**优先级**: High  
**影响范围**: MongoDB 资源的状态管理

---

## 🎯 问题描述

### 当前问题

`tencentcloud_mongodb_instance` 等资源的 `availability_zone_list` 字段存在以下矛盾:

1. **创建时要求**:
   - API 要求严格的可用区顺序
   - 必须使用 `TypeList` 保持顺序
   
2. **读取时问题**:
   - 查询接口返回的可用区列表是**无序的**
   - 导致即使实际值相同,仅因顺序不同就产生 diff
   - 用户体验差,每次 `terraform plan` 都显示不必要的变更

### 错误的解决方案

之前尝试将 `availability_zone_list` 从 `TypeList` 改为 `TypeSet`:
- ❌ 破坏了创建时的顺序要求
- ❌ 可能导致 API 调用失败
- ❌ 属于破坏性变更

### 正确的解决方案

保持 `TypeList` 结构,但添加 `DiffSuppressFunc` 来忽略顺序差异:
- ✅ 保持创建时的顺序要求
- ✅ 读取时忽略顺序差异
- ✅ 非破坏性变更
- ✅ 用户体验优化

---

## 🔧 技术方案

### 方案设计

使用 Terraform Schema 的 `DiffSuppressFunc` 特性,在比较新旧值时忽略列表元素的顺序差异。

### 实现原理

```go
// DiffSuppressFunc 会在以下情况被调用:
// 1. terraform plan 时比较配置与状态
// 2. terraform apply 前的最终确认
// 3. 任何需要判断资源是否变更的场景

func mongodbAvailabilityZoneListDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
    // k: 字段路径,如 "availability_zone_list.0", "availability_zone_list.1"
    // old: 状态中的旧值
    // new: 配置中的新值
    // d: 资源数据对象
    
    // 如果不是根字段(是子元素),直接返回 false
    if k != "availability_zone_list" && !strings.HasSuffix(k, ".#") {
        return false
    }
    
    // 获取完整的列表进行比较
    oldList, newList := d.GetChange("availability_zone_list")
    
    // 转换为字符串切片
    oldZones := helper.InterfacesStrings(oldList.([]interface{}))
    newZones := helper.InterfacesStrings(newList.([]interface{}))
    
    // 长度不同,肯定有变化
    if len(oldZones) != len(newZones) {
        return false
    }
    
    // 排序后比较内容是否相同
    sort.Strings(oldZones)
    sort.Strings(newZones)
    
    for i := range oldZones {
        if oldZones[i] != newZones[i] {
            return false // 内容不同
        }
    }
    
    return true // 内容相同,仅顺序不同,忽略 diff
}
```

### Schema 修改

```go
"availability_zone_list": {
    Type:     schema.TypeList,  // 保持 List 类型
    Optional: true,
    Computed: true,
    DiffSuppressFunc: mongodbAvailabilityZoneListDiffSuppress,  // 添加 diff 抑制函数
    Elem: &schema.Schema{
        Type: schema.TypeString,
    },
    Description: "...",
},
```

---

## 📊 影响分析

### 影响的资源

| 资源 | 字段 | 变更类型 |
|------|------|----------|
| `tencentcloud_mongodb_instance` | `availability_zone_list` | 添加 DiffSuppressFunc |
| `tencentcloud_mongodb_sharding_instance` | `availability_zone_list` | 添加 DiffSuppressFunc |
| `tencentcloud_mongodb_readonly_instance` | `availability_zone_list` | 添加 DiffSuppressFunc |

### 用户影响

**正面影响**:
- ✅ 不再显示仅顺序不同的无意义 diff
- ✅ 用户体验提升
- ✅ 减少误操作风险

**破坏性评估**:
- ✅ **非破坏性变更**
- ✅ 不影响现有配置
- ✅ 不改变 API 调用行为
- ✅ 不需要状态迁移

---

## 🎯 实施计划

### 阶段 1: 准备工作
1. 创建 DiffSuppressFunc 辅助函数
2. 编写单元测试验证逻辑
3. 准备测试用例

### 阶段 2: 代码修改
1. 修改 `resource_tc_mongodb_instance.go`
2. 修改 `resource_tc_mongodb_sharding_instance.go`
3. 修改 `resource_tc_mongodb_readonly_instance.go`
4. 代码格式化和 lint 检查

### 阶段 3: 测试验证
1. 单元测试验证 DiffSuppressFunc 逻辑
2. 集成测试验证实际场景
3. 边界条件测试

### 阶段 4: 文档更新
1. 更新资源文档
2. 添加示例说明
3. 更新 CHANGELOG

---

## 🧪 测试场景

### 场景 1: 顺序不同但内容相同
```hcl
# 配置
availability_zone_list = ["ap-guangzhou-3", "ap-guangzhou-4", "ap-guangzhou-6"]

# 状态(API 返回)
availability_zone_list = ["ap-guangzhou-4", "ap-guangzhou-6", "ap-guangzhou-3"]

# 预期: 不显示 diff
```

### 场景 2: 内容不同
```hcl
# 配置
availability_zone_list = ["ap-guangzhou-3", "ap-guangzhou-4"]

# 状态
availability_zone_list = ["ap-guangzhou-3", "ap-guangzhou-6"]

# 预期: 显示 diff,需要更新
```

### 场景 3: 长度不同
```hcl
# 配置
availability_zone_list = ["ap-guangzhou-3", "ap-guangzhou-4"]

# 状态
availability_zone_list = ["ap-guangzhou-3"]

# 预期: 显示 diff,需要更新
```

### 场景 4: 空列表
```hcl
# 配置
availability_zone_list = []

# 状态
availability_zone_list = []

# 预期: 不显示 diff
```

---

## 📝 代码变更清单

### 新增文件
- `tencentcloud/services/mongodb/diff_suppress_funcs.go` - 存放 DiffSuppressFunc 函数

### 修改文件
1. `resource_tc_mongodb_instance.go` - 添加 DiffSuppressFunc
2. `resource_tc_mongodb_sharding_instance.go` - 添加 DiffSuppressFunc
3. `resource_tc_mongodb_readonly_instance.go` - 添加 DiffSuppressFunc

### 单元测试文件
- `diff_suppress_funcs_test.go` - 测试 DiffSuppressFunc 逻辑

---

## ✅ 验收标准

1. **功能正确性**:
   - ✅ 顺序不同但内容相同时不显示 diff
   - ✅ 内容不同时正常显示 diff
   - ✅ 创建时顺序保持不变

2. **代码质量**:
   - ✅ 通过所有单元测试
   - ✅ 通过 linter 检查
   - ✅ 代码格式符合规范

3. **文档完整**:
   - ✅ 代码注释清晰
   - ✅ 测试覆盖充分
   - ✅ 变更文档完整

---

## 🎓 技术要点

### Terraform DiffSuppressFunc

```go
type DiffSuppressFunc func(k, old, new string, d *schema.ResourceData) bool
```

**参数说明**:
- `k`: 字段路径(如 "availability_zone_list.0")
- `old`: 状态中的值
- `new`: 配置中的值
- `d`: ResourceData 对象

**返回值**:
- `true`: 忽略此 diff
- `false`: 保留此 diff

**调用时机**:
- terraform plan
- terraform apply 前
- 状态刷新后

**最佳实践**:
1. 只在根字段级别比较完整列表
2. 避免在子元素级别重复比较
3. 考虑边界条件(nil, empty list)
4. 添加详细的单元测试

---

## 📚 参考资料

- [Terraform Schema DiffSuppressFunc](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema#DiffSuppressFunc)
- [Terraform Provider Best Practices](https://www.terraform.io/docs/extend/best-practices/index.html)
- 类似实现参考: AWS Provider 的 tags 处理

---

## 🔄 与之前方案的对比

| 对比项 | 之前方案(List→Set) | 新方案(DiffSuppressFunc) |
|-------|-------------------|------------------------|
| 创建时顺序 | ❌ 丢失顺序信息 | ✅ 保持顺序 |
| API 兼容性 | ❌ 可能失败 | ✅ 完全兼容 |
| 破坏性 | ❌ Breaking Change | ✅ 非破坏性 |
| 状态迁移 | ❌ 需要 | ✅ 不需要 |
| 实现复杂度 | 简单 | 中等 |
| 用户体验 | ✅ 解决 diff 问题 | ✅ 解决 diff 问题 |

---

## ⚠️ 风险与限制

### 潜在风险
1. **DiffSuppressFunc 调用时机**: 需要确保在正确的时机调用
2. **性能影响**: 每次 plan 都会调用,但影响可忽略(列表通常很短)
3. **测试覆盖**: 需要充分的测试确保各种场景正确

### 限制条件
1. 仅解决顺序差异问题,不改变 List 语义
2. 用户仍然需要按 API 要求的顺序创建资源
3. 不影响其他字段的 diff 行为

---

## 🎉 预期收益

1. **用户体验提升**:
   - 消除无意义的 diff 警告
   - 减少用户困惑
   - 降低误操作风险

2. **维护性提升**:
   - 减少用户支持工单
   - 代码逻辑清晰
   - 符合 Terraform 最佳实践

3. **兼容性保证**:
   - 非破坏性变更
   - 无需版本迁移
   - 平滑升级

---

**提案状态**: 待审核  
**预计实施时间**: 2 小时  
**风险等级**: 低
