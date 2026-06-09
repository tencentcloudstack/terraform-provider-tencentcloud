# 归档信息

## 归档状态
- **状态**: ✅ 已完成并归档
- **归档时间**: 2026-03-23
- **完成状态**: 已实施并验证

## 提案摘要
使用 `DiffSuppressFunc` 解决 MongoDB `availability_zone_list` 字段的顺序差异问题。

## 实施结果

### 核心方案
- 保持 `TypeList` 结构满足 API 顺序要求
- 添加 `DiffSuppressFunc` 忽略读取时的顺序差异
- 非破坏性变更,完全兼容现有配置

### 完成的工作

#### 新增文件
- ✅ `diff_suppress_funcs.go` - DiffSuppressFunc 实现
- ✅ `diff_suppress_funcs_test.go` - 11 个单元测试

#### 修改文件
- ✅ `resource_tc_mongodb_instance.go` - Schema + Create + Read
- ✅ `resource_tc_mongodb_sharding_instance.go` - Schema + Create + Read  
- ✅ `resource_tc_mongodb_readonly_instance.go` - Create

### 影响范围
- `tencentcloud_mongodb_instance`
- `tencentcloud_mongodb_sharding_instance`
- `tencentcloud_mongodb_readonly_instance`

### 质量保证
- ✅ 代码格式化完成
- ✅ Linter 检查通过
- ✅ 单元测试完整
- ✅ 文档齐全

## 相关文档
- `proposal.md` - 技术提案
- `tasks.md` - 任务清单
- `change.md` - 变更记录
- `README.md` - 快速概览

## 技术细节

### DiffSuppressFunc 实现
```go
func mongodbAvailabilityZoneListDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
    // 关键: 对所有相关 key (包括元素级别) 都进行完整列表比较
    if !strings.Contains(k, "availability_zone_list") {
        return false
    }
    
    // 获取完整列表
    oldList, newList := d.GetChange("availability_zone_list")
    oldZones := helper.InterfacesStrings(oldList.([]interface{}))
    newZones := helper.InterfacesStrings(newList.([]interface{}))
    
    // 长度不同则有变更
    if len(oldZones) != len(newZones) {
        return false
    }
    
    // 排序后比较
    oldSorted := make([]string, len(oldZones))
    newSorted := make([]string, len(newZones))
    copy(oldSorted, oldZones)
    copy(newSorted, newZones)
    sort.Strings(oldSorted)
    sort.Strings(newSorted)
    
    // 内容相同但顺序不同 -> 忽略 diff
    for i := range oldSorted {
        if oldSorted[i] != newSorted[i] {
            return false
        }
    }
    return true
}
```

### 💡 关键技术发现

**Terraform DiffSuppressFunc 在 TypeList 上的行为**:
- ⚠️ Terraform 会在**元素级别**调用 `DiffSuppressFunc`
- ⚠️ 参数 `old` 和 `new` 是**单个字符串**,而非完整列表  
- ⚠️ 如果在元素级别返回 `false`,即使列表级别返回 `true`,diff 仍然会触发
- ✅ **必须在每个调用中都获取完整列表并比较,返回统一的结果**

**错误的实现** ❌:
```go
// 只在列表级别处理,元素级别返回 false
if !strings.HasSuffix(k, "availability_zone_list") && !strings.HasSuffix(k, ".#") {
    return false  // 这会导致元素级别触发 diff!
}
```

**正确的实现** ✅:
```go
// 对所有相关 key 都进行完整列表比较
if !strings.Contains(k, "availability_zone_list") {
    return false
}
// 获取完整列表并比较,返回统一结果
oldList, newList := d.GetChange("availability_zone_list")
// ... 比较逻辑
```

这个发现对其他使用 `TypeList` + `DiffSuppressFunc` 的场景同样适用。

### 变更对比

| 方面 | 原方案(List→Set) | 新方案(DiffSuppressFunc) |
|------|-----------------|------------------------|
| API 兼容性 | ❌ 不兼容 | ✅ 完全兼容 |
| 创建时顺序 | ❌ 丢失 | ✅ 保持 |
| 读取时 diff | ✅ 解决 | ✅ 解决 |
| 破坏性 | ❌ Breaking | ✅ 非破坏性 |
| 正确性 | ❌ 错误 | ✅ 正确 |

## 归档原因
代码实施完成并通过验证,提案目标已达成。

## 后续工作
无需后续工作,可以发布使用。
