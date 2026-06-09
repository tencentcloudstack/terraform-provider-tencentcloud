# MongoDB availability_zone_list DiffSuppressFunc 优化

## 🎯 变更概述

为 MongoDB 资源的 `availability_zone_list` 字段添加 `DiffSuppressFunc`,解决因 API 返回顺序不一致导致的无意义 diff 问题。

**变更类型**: Enhancement (非破坏性)  
**影响资源**: 
- `tencentcloud_mongodb_instance`
- `tencentcloud_mongodb_sharding_instance`
- `tencentcloud_mongodb_readonly_instance`

---

## ❓ 为什么需要这个变更

### 问题背景

```hcl
# 用户配置
resource "tencentcloud_mongodb_instance" "example" {
  availability_zone_list = [
    "ap-guangzhou-3",
    "ap-guangzhou-4",
    "ap-guangzhou-6"
  ]
}
```

**问题现象**:
```bash
$ terraform plan

  # tencentcloud_mongodb_instance.example will be updated in-place
  ~ availability_zone_list = [
      - "ap-guangzhou-3",
      - "ap-guangzhou-4",
      - "ap-guangzhou-6",
      + "ap-guangzhou-4",  # API 返回顺序不同
      + "ap-guangzhou-6",
      + "ap-guangzhou-3",
    ]
```

**根本原因**:
- 创建时 API 要求**严格顺序**,必须使用 `TypeList`
- 读取时 API 返回的列表是**无序的**
- 导致每次 `terraform plan` 都显示不必要的 diff

---

## ✅ 解决方案

使用 Terraform 的 `DiffSuppressFunc` 特性,在比较时忽略列表元素的顺序差异:

```go
"availability_zone_list": {
    Type:             schema.TypeList,  // 保持 List 类型
    DiffSuppressFunc: mongodbAvailabilityZoneListDiffSuppress,  // 添加 diff 抑制
    // ... 其他配置
}
```

**核心逻辑**:
1. 获取配置和状态中的完整列表
2. 排序后比较内容是否相同
3. 相同则返回 `true` 抑制 diff
4. 不同则返回 `false` 显示 diff

---

## 🎨 技术方案

### 实现函数

```go
func mongodbAvailabilityZoneListDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
    // 获取完整列表
    oldList, newList := d.GetChange("availability_zone_list")
    
    // 转换为字符串切片
    oldZones := helper.InterfacesStrings(oldList.([]interface{}))
    newZones := helper.InterfacesStrings(newList.([]interface{}))
    
    // 长度检查
    if len(oldZones) != len(newZones) {
        return false
    }
    
    // 排序后比较
    sort.Strings(oldZones)
    sort.Strings(newZones)
    
    for i := range oldZones {
        if oldZones[i] != newZones[i] {
            return false
        }
    }
    
    return true // 内容相同,忽略顺序差异
}
```

### 测试用例

| 场景 | 旧值 | 新值 | 是否抑制 diff |
|------|------|------|--------------|
| 顺序不同 | `[a,b,c]` | `[c,a,b]` | ✅ 是 |
| 内容不同 | `[a,b]` | `[a,c]` | ❌ 否 |
| 长度不同 | `[a,b]` | `[a]` | ❌ 否 |
| 完全相同 | `[a,b]` | `[a,b]` | ✅ 是 |
| 都为空 | `[]` | `[]` | ✅ 是 |

---

## 📊 影响范围

### 修改的文件

1. **新增文件**:
   - `tencentcloud/services/mongodb/diff_suppress_funcs.go`
   - `tencentcloud/services/mongodb/diff_suppress_funcs_test.go`

2. **修改文件**:
   - `tencentcloud/services/mongodb/resource_tc_mongodb_instance.go`
   - `tencentcloud/services/mongodb/resource_tc_mongodb_sharding_instance.go`
   - `tencentcloud/services/mongodb/resource_tc_mongodb_readonly_instance.go`

### 用户影响

**正面影响**:
- ✅ 不再看到无意义的顺序 diff
- ✅ 减少误操作风险
- ✅ 提升用户体验

**破坏性评估**:
- ✅ **完全非破坏性**
- ✅ 不改变 API 行为
- ✅ 不需要状态迁移
- ✅ 不影响现有配置

---

## 🚀 实施步骤

### 快速上手

```bash
# 1. 创建 DiffSuppressFunc
vim tencentcloud/services/mongodb/diff_suppress_funcs.go

# 2. 创建测试
vim tencentcloud/services/mongodb/diff_suppress_funcs_test.go

# 3. 运行测试
go test -v -run TestMongodbAvailabilityZoneListDiffSuppress

# 4. 修改 resource schema
# 在 availability_zone_list 添加: DiffSuppressFunc: mongodbAvailabilityZoneListDiffSuppress

# 5. 格式化和检查
go fmt ./tencentcloud/services/mongodb/...
go vet ./tencentcloud/services/mongodb/...
```

详细步骤请查看 [tasks.md](./tasks.md)

---

## 🧪 验证方法

### 验证步骤

1. **创建测试实例**:
```hcl
resource "tencentcloud_mongodb_instance" "test" {
  availability_zone_list = [
    "ap-guangzhou-3",
    "ap-guangzhou-4", 
    "ap-guangzhou-6"
  ]
  # ... 其他配置
}
```

2. **初次应用**:
```bash
terraform apply
```

3. **验证不显示 diff**:
```bash
terraform plan
# 预期: No changes. Infrastructure is up-to-date.
```

4. **修改内容验证**:
```hcl
resource "tencentcloud_mongodb_instance" "test" {
  availability_zone_list = [
    "ap-guangzhou-3",
    "ap-guangzhou-5",  # 改变内容
    "ap-guangzhou-6"
  ]
}
```

```bash
terraform plan
# 预期: 显示 diff,需要更新
```

---

## ⚠️ 注意事项

### 重要提醒

1. **保持 List 类型**: 
   - ❌ 不要改为 `TypeSet`
   - ✅ 保持 `TypeList`
   - 原因: API 创建时需要顺序信息

2. **不修改 CRUD 逻辑**:
   - 只添加 DiffSuppressFunc
   - 不改变创建/读取/更新/删除逻辑
   - 保持 API 调用方式不变

3. **测试覆盖**:
   - 必须包含单元测试
   - 覆盖各种边界情况
   - 验证实际场景

### 常见问题

**Q: DiffSuppressFunc 会影响资源创建吗?**  
A: 不会。它只影响 diff 计算,不改变 API 调用。

**Q: 为什么不用 TypeSet?**  
A: 因为 API 创建时需要严格的顺序,Set 会丢失顺序信息。

**Q: 会导致状态不一致吗?**  
A: 不会。状态依然是 API 返回的原始值,只是 plan 时忽略顺序差异。

---

## 📚 相关文档

- [proposal.md](./proposal.md) - 完整的技术提案
- [tasks.md](./tasks.md) - 详细的实施任务清单
- [Terraform DiffSuppressFunc 文档](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema#DiffSuppressFunc)

---

## 📈 预期收益

| 指标 | 改进 |
|------|------|
| 无意义 diff | ✅ 消除 |
| 用户体验 | ⬆️ 提升 |
| 支持工单 | ⬇️ 减少 |
| 破坏性 | ✅ 无 |
| 实施难度 | 🟢 低 |
| 测试覆盖 | ✅ 完整 |

---

## 🎉 总结

这是一个**低风险、高收益**的优化变更:

✅ **优势**:
- 解决实际用户痛点
- 非破坏性变更
- 实现简单清晰
- 测试覆盖完整

✅ **最佳实践**:
- 使用 Terraform 推荐的 DiffSuppressFunc
- 保持 API 语义不变
- 充分的单元测试
- 清晰的文档说明

---

**状态**: 待实施  
**优先级**: High  
**预计时间**: 2 小时
